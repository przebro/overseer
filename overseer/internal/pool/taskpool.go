package pool

import (
	"context"
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/config"
	"github.com/przebro/overseer/overseer/internal/events"
	"github.com/przebro/overseer/overseer/internal/pool/activetask"
	"github.com/przebro/overseer/overseer/internal/pool/calc"
	"github.com/przebro/overseer/overseer/internal/pool/models"
	"github.com/przebro/overseer/overseer/internal/pool/readers"
	"github.com/przebro/overseer/overseer/internal/pool/states"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type ResourceManager interface {
	CheckTickets([]types.CollectedTicketModel) []types.CollectedTicketModel
	ProcessTicketAction([]types.TicketActionModel) bool
	ProcessAcquireFlag([]types.FlagModel) (bool, []string)
	ProcessReleaseFlag([]string) (bool, []string)
}

type WorkManager interface {
	Push(ctx context.Context, t types.TaskDescription, vars types.EnvironmentVariableList) (types.WorkerTaskStatus, error)
	Status(ctx context.Context, t types.WorkDescription) types.TaskExecutionStatus
}

// TaskViewer - Provides a view for an active tasks in pool
type TaskViewer interface {
	Detail(unique.TaskOrderID) (events.TaskDetailResultMsg, error)
	List(filter string) []events.TaskInfoResultMsg
}

// ActiveTaskPool - Holds tasks that are currently processed.
type ActiveTaskPool struct {
	config              config.ActivePoolConfiguration
	tasks               *Store
	log                 zerolog.Logger
	isProcActive        bool
	processing          chan *activetask.TaskInstance
	enforcedTasks       map[unique.TaskOrderID]bool
	lock                sync.RWMutex
	shutdown            chan struct{}
	activate            chan bool
	done                <-chan struct{}
	activeDefinitionRWC readers.ActiveDefinitionReadWriterRemover
	resourceManager     ResourceManager
	workerManager       WorkManager
	journal             readers.JournalWriter
}

// NewTaskPool - creates new task pool
func NewTaskPool(
	cfg config.ActivePoolConfiguration,
	provider *datastore.Provider,
	isProcActive bool,
	activeDefinitionRWC readers.ActiveDefinitionReadWriterRemover,
	wm WorkManager,
	rm ResourceManager,
	jw readers.JournalWriter,
) (*ActiveTaskPool, error) {

	var store *Store
	var err error

	lg := log.With().Str("component", "pool").Logger()

	if store, err = NewStore(lg, cfg.SyncTime, provider, activeDefinitionRWC); err != nil {
		return nil, err
	}

	pool := &ActiveTaskPool{
		tasks:               store,
		config:              cfg,
		isProcActive:        isProcActive,
		log:                 lg,
		processing:          make(chan *activetask.TaskInstance, 8),
		enforcedTasks:       map[unique.TaskOrderID]bool{},
		lock:                sync.RWMutex{},
		activate:            make(chan bool),
		shutdown:            make(chan struct{}),
		activeDefinitionRWC: activeDefinitionRWC,
		resourceManager:     rm,
		workerManager:       wm,
		journal:             jw,
	}

	if pool.isProcActive {
		pool.log.Info().Msg("Starting in ACTIVE mode")
	} else {
		pool.log.Info().Msg("Starting in QUIESCE mode")
	}

	return pool, nil
}

func (pool *ActiveTaskPool) CleanupCompletedTasks() int {

	pool.log.Info().Msg("Cleanup tasks")
	var numDeleted = 0

	pool.tasks.ForEach(func(k unique.TaskOrderID, v *activetask.TaskInstance) {

		cleanDate := date.AddDays(v.OrderDate(), 7) //:TODO: make it configurable
		if v.State() == models.TaskStateEndedOk && date.IsBeforeCurrent(cleanDate, date.CurrentOdate()) {
			delete(pool.tasks.store, v.OrderID())

			numDeleted++
		}
	})

	pool.log.Info().Int("deleted", numDeleted).Msg("cleanup comlpete")
	return numDeleted
}

func (pool *ActiveTaskPool) enforceTask(taskID unique.TaskOrderID) {
	pool.log.Info().Str("task_id", string(taskID)).Msg("Enforce task")
	defer pool.lock.Unlock()
	pool.lock.Lock()
	pool.enforcedTasks[taskID] = true
}

func (pool *ActiveTaskPool) isEnforced(taskID unique.TaskOrderID) bool {
	defer pool.lock.Unlock()
	pool.lock.Lock()
	pool.log.Debug().Str("task_id", string(taskID)).Bool("enforced", pool.enforcedTasks[taskID]).Msg("checking isEnforced")
	enforced := pool.enforcedTasks[taskID]
	delete(pool.enforcedTasks, taskID)

	return enforced
}

func (pool *ActiveTaskPool) addTask(orderID unique.TaskOrderID, t *activetask.TaskInstance) {

	pool.tasks.add(orderID, t)
}

// task - Returns an active task with given id or error if the task was not found.
func (pool *ActiveTaskPool) task(orderID unique.TaskOrderID) (*activetask.TaskInstance, error) {

	var err error = nil

	task, exists := pool.tasks.get(orderID)
	if !exists {
		return nil, errors.New("task does not exists")
	}

	return task, err

}

func (pool *ActiveTaskPool) cycleTasks(t time.Time) {

	//tsart := time.Now()
	routines := 8

	tchannel := make(chan *activetask.TaskInstance, pool.tasks.len())
	wg := sync.WaitGroup{}
	wg.Add(routines)

	for x := 0; x < routines; x++ {
		go pool.processTaskState(tchannel, &wg, t)
	}

	pool.tasks.Over(func(k unique.TaskOrderID, v *activetask.TaskInstance) { tchannel <- v })

	close(tchannel)
	wg.Wait()

	//pool.log.Info().Int("total", pool.tasks.len()).Dur("duration", time.Since(tsart)).Msg("completed")

}

func (pool *ActiveTaskPool) processTaskState(ch <-chan *activetask.TaskInstance, wg *sync.WaitGroup, t time.Time) {

	for task := range ch {

		executionState := pool.getProcessState(task.State(), task.IsHeld())
		if executionState == nil {
			continue
		}

		exCtx := &states.TaskExecutionContext{
			Odate:      date.CurrentOdate(),
			Task:       task,
			Time:       t,
			MaxRc:      pool.config.MaxOkReturnCode,
			State:      executionState,
			Log:        pool.log,
			IsEnforced: pool.isEnforced(task.OrderID()),
			IsInTime:   false,
			Rmanager:   pool.resourceManager,
			Wmanager:   pool.workerManager,
			Journal:    pool.journal,
		}
		for exCtx.State.ProcessState(exCtx) {
		}

		n, g := task.Definition.Name, task.Definition.Group

		pool.log.Debug().Str("task_id", string(task.OrderID())).Str("group", g).Str("name", n).Str("state", task.State().String()).Msg("state")

	}
	wg.Done()

}

// Detail - Gets task details
func (pool *ActiveTaskPool) Detail(orderID unique.TaskOrderID) (events.TaskDetailResultMsg, error) {

	result := events.TaskDetailResultMsg{}
	var t *activetask.TaskInstance
	var exists bool

	if t, exists = pool.tasks.get(orderID); !exists {
		return result, ErrUnableFindTask
	}

	result.Name, result.Group, result.Description = t.Definition.Name, t.Definition.Group, t.Definition.Description
	result.Odate = t.OrderDate()
	result.TaskID = t.OrderID()
	result.State = int32(t.State())
	result.Confirmed = t.Confirmed()
	result.EndTime = t.EndTime().Format("2006-01-02 15:04:05")
	result.StartTime = t.StartTime().Format("2006-01-02 15:04:05")
	result.Held = t.IsHeld()
	result.RunNumber = int32(t.RunNumber())
	result.Worker = t.WorkerName()

	if c := t.Cyclic; c.IsCycle {

		result.TaskCycleMsg = events.TaskCycleMsg{
			IsCyclic:    c.IsCycle,
			NextRun:     "", //c.NextRun,
			RunFrom:     string(c.RunFrom),
			MaxRun:      c.MaxRuns,
			RunInterval: c.TimeInterval,
		}
	}
	result.From, result.To = func(f, t types.HourMinTime) (string, string) { return string(f), string(t) }(t.TimeSpan())

	result.Tickets = make([]struct {
		Name      string
		Odate     date.Odate
		Fulfilled bool
	}, 0)

	tickets := make([]types.CollectedTicketModel, 0)

	for _, e := range t.InTickets {

		realOdat := calc.CalcRealOdate(t.OrderDate(), e.Odate, t.Schedule)
		tickets = append(tickets, types.CollectedTicketModel{Odate: realOdat, Name: e.Name, Exists: false})
	}

	tickets = pool.resourceManager.CheckTickets(tickets)

	for _, ticket := range tickets {
		result.Tickets = append(result.Tickets, struct {
			Name      string
			Odate     date.Odate
			Fulfilled bool
		}{Name: ticket.Name, Odate: ticket.Odate, Fulfilled: ticket.Exists})
	}

	return result, nil

}

// List - Filters and Lists tasks in pool
func (pool *ActiveTaskPool) List(filter string) []events.TaskInfoResultMsg {

	result := make([]events.TaskInfoResultMsg, 0)

	pool.tasks.Over(func(k unique.TaskOrderID, v *activetask.TaskInstance) {
		n, g, _ := v.GetInfo()
		data := events.TaskInfoResultMsg{Group: g,
			Name:      n,
			Odate:     v.OrderDate(),
			TaskID:    v.OrderID(),
			State:     int32(v.State()),
			RunNumber: v.RunNumber(),
			Held:      v.IsHeld(),
			Confirmed: v.Confirmed(),
		}
		result = append(result, data)
	})

	sort.Sort(taskInfoSorter{result})

	return result
}

func (pool *ActiveTaskPool) NewDayProc() types.HourMinTime {
	return pool.config.NewDayProc
}

func (pool *ActiveTaskPool) ProcessTimeEvent(t time.Time) {
	pool.cycleTasks(t)
}

// Start - starts tasks processing
func (pool *ActiveTaskPool) Start() error {

	defer pool.lock.Unlock()
	pool.lock.Lock()

	if pool.done == nil {
		pool.log.Info().Msg("starting store watch")
		pool.done = pool.tasks.watch(pool.activate, pool.shutdown)
	}

	pool.activate <- pool.isProcActive
	<-pool.done

	return nil
}

// Shutdown - shutdowns task pool
func (pool *ActiveTaskPool) Shutdown() error {

	defer pool.lock.Unlock()
	pool.lock.Lock()
	pool.isProcActive = false
	pool.shutdown <- struct{}{}
	<-pool.done

	return nil
}

// Resume - resumes tasks processing
func (pool *ActiveTaskPool) Resume() error {

	defer pool.lock.Unlock()
	pool.lock.Lock()

	pool.isProcActive = true
	pool.activate <- pool.isProcActive
	<-pool.done

	return nil
}

// Quiesce - puts taskpool into sleep mode
func (pool *ActiveTaskPool) Quiesce() error {

	defer pool.lock.Unlock()
	pool.lock.Lock()

	pool.isProcActive = false
	pool.activate <- pool.isProcActive

	<-pool.done

	return nil
}

func (pool *ActiveTaskPool) getProcessState(state models.TaskState, isHeld bool) states.TaskExecutionState {

	if isHeld {
		return nil
	}

	if state == models.TaskStateWaiting {
		return &states.OstateConfirm{}
	}
	if state == models.TaskStateExecuting {
		return &states.OstateExecuting{}
	}
	//Any other case means that task should not be processed.
	return nil

}
