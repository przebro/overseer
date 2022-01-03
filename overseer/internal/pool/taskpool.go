package pool

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/config"
	"github.com/przebro/overseer/overseer/internal/events"
	"github.com/przebro/overseer/overseer/internal/unique"
)

//ActiveTaskPool - Holds tasks that are currently processed.
type ActiveTaskPool struct {
	config              config.ActivePoolConfiguration
	dispatcher          events.Dispatcher
	tasks               *Store
	log                 logger.AppLogger
	isProcActive        bool
	processing          chan *activeTask
	enforcedTasks       map[unique.TaskOrderID]bool
	lock                sync.RWMutex
	shutdown            chan struct{}
	activate            chan bool
	done                <-chan struct{}
	activeDefinitionRWC ActiveDefinitionReadWriterRemover
}

//TaskViewer - Provides a view for an active tasks in pool
type TaskViewer interface {
	Detail(unique.TaskOrderID) (events.TaskDetailResultMsg, error)
	List(filter string) []events.TaskInfoResultMsg
}

//NewTaskPool - creates new task pool
func NewTaskPool(dispatcher events.Dispatcher,
	cfg config.ActivePoolConfiguration,
	provider *datastore.Provider,
	isProcActive bool,
	log logger.AppLogger,
	activeDefinitionRWC ActiveDefinitionReadWriterRemover) (*ActiveTaskPool, error) {

	var store *Store
	var err error

	if store, err = NewStore(cfg.Collection, log, cfg.SyncTime, provider, activeDefinitionRWC); err != nil {
		return nil, err
	}

	pool := &ActiveTaskPool{
		tasks:               store,
		dispatcher:          dispatcher,
		config:              cfg,
		isProcActive:        isProcActive,
		log:                 log,
		processing:          make(chan *activeTask, 8),
		enforcedTasks:       map[unique.TaskOrderID]bool{},
		lock:                sync.RWMutex{},
		activate:            make(chan bool),
		shutdown:            make(chan struct{}),
		activeDefinitionRWC: activeDefinitionRWC,
	}

	if dispatcher != nil {
		dispatcher.Subscribe(events.RouteTimeOut, pool)
	}

	if pool.isProcActive {
		pool.log.Info("Starting in ACTIVE mode")
	} else {
		pool.log.Info("Starting in QUIESCE mode")
	}

	return pool, nil
}

func (pool *ActiveTaskPool) cleanupCompletedTasks() int {

	pool.log.Info("Cleanup tasks")
	var numDeleted = 0

	pool.tasks.ForEach(func(k unique.TaskOrderID, v *activeTask) {

		cleanDate := date.AddDays(v.OrderDate(), v.Retention())
		if v.State() == TaskStateEndedOk && date.IsBeforeCurrent(cleanDate, date.CurrentOdate()) {
			delete(pool.tasks.store, v.OrderID())

			numDeleted++
		}
	})

	pool.log.Info(fmt.Sprintf("cleanup comlpete. %d tasks deleted.", numDeleted))
	return numDeleted
}

func (pool *ActiveTaskPool) enforceTask(taskID unique.TaskOrderID) {
	defer pool.lock.Unlock()
	pool.lock.Lock()
	pool.enforcedTasks[taskID] = true
}

func (pool *ActiveTaskPool) isEnforced(taskID unique.TaskOrderID) bool {
	defer pool.lock.Unlock()
	pool.lock.Lock()

	enforced := pool.enforcedTasks[taskID]
	delete(pool.enforcedTasks, taskID)

	return enforced
}

func (pool *ActiveTaskPool) addTask(orderID unique.TaskOrderID, t *activeTask) {

	pool.tasks.add(orderID, t)
}

//task - Returns an active task with given id or error if the task was not found.
func (pool *ActiveTaskPool) task(orderID unique.TaskOrderID) (*activeTask, error) {

	var err error = nil

	task, exists := pool.tasks.get(orderID)
	if !exists {
		return nil, errors.New("task does not exists")
	}

	return task, err

}

func (pool *ActiveTaskPool) cycleTasks(t time.Time) {

	tsart := time.Now()
	routines := 8

	tchannel := make(chan *activeTask, pool.tasks.len())
	wg := sync.WaitGroup{}
	wg.Add(routines)

	for x := 0; x < routines; x++ {
		go pool.processTaskState(tchannel, &wg, t)
	}

	pool.tasks.Over(func(k unique.TaskOrderID, v *activeTask) { tchannel <- v })

	close(tchannel)
	wg.Wait()
	pool.log.Info(time.Since(tsart))
}

func (pool *ActiveTaskPool) processTaskState(ch <-chan *activeTask, wg *sync.WaitGroup, t time.Time) {

	for task := range ch {

		executionState := getProcessState(task.State(), task.IsHeld())
		if executionState == nil {
			continue
		}

		exCtx := &TaskExecutionContext{
			odate:        date.CurrentOdate(),
			task:         task,
			time:         t,
			maxRc:        pool.config.MaxOkReturnCode,
			state:        executionState,
			dispatcher:   pool.dispatcher,
			log:          pool.log,
			isEnforced:   pool.isEnforced(task.OrderID()),
			isInTime:     false,
			scheduleTime: pool.config.NewDayProc,
		}
		for exCtx.state.processState(exCtx) {
		}

		n, g, _ := task.GetInfo()
		pool.log.Debug(n, ":", g, " Task state:", task.State(), task.OrderID())
	}
	wg.Done()

}

//Detail - Gets task details
func (pool *ActiveTaskPool) Detail(orderID unique.TaskOrderID) (events.TaskDetailResultMsg, error) {

	result := events.TaskDetailResultMsg{}
	var t *activeTask
	var exists bool

	if t, exists = pool.tasks.get(orderID); !exists {
		return result, ErrUnableFindTask
	}

	result.Name, result.Group, result.Description = t.GetInfo()
	result.Odate = t.OrderDate()
	result.TaskID = t.OrderID()
	result.State = int32(t.State())
	result.Confirmed = t.Confirmed()
	result.EndTime = t.EndTime().Format("2006-01-02 15:04:05")
	result.StartTime = t.StartTime().Format("2006-01-02 15:04:05")
	result.Held = t.IsHeld()
	result.RunNumber = int32(t.RunNumber())
	result.WaitingInfo = t.WaitingInfo()
	result.Worker = t.WorkerName()

	if c := t.CycleData(); c.IsCyclic {

		result.TaskCycleMsg = events.TaskCycleMsg{
			IsCyclic:    c.IsCyclic,
			NextRun:     c.NextRun,
			RunFrom:     c.RunFrom,
			MaxRun:      c.MaxRun,
			RunInterval: c.RunInterval,
		}
	}
	result.From, result.To = func(f, t types.HourMinTime) (string, string) { return string(f), string(t) }(t.TimeSpan())

	result.Tickets = make([]struct {
		Name      string
		Odate     date.Odate
		Fulfilled bool
	}, 0)

	for _, ticket := range t.Collected() {
		result.Tickets = append(result.Tickets, struct {
			Name      string
			Odate     date.Odate
			Fulfilled bool
		}{Name: ticket.name, Odate: date.Odate(ticket.odate), Fulfilled: ticket.fulfilled})
	}

	return result, nil

}

//List - Filters and Lists tasks in pool
func (pool *ActiveTaskPool) List(filter string) []events.TaskInfoResultMsg {

	result := make([]events.TaskInfoResultMsg, 0)

	pool.tasks.Over(func(k unique.TaskOrderID, v *activeTask) {
		n, g, _ := v.GetInfo()
		data := events.TaskInfoResultMsg{Group: g,
			Name:        n,
			Odate:       v.OrderDate(),
			TaskID:      v.OrderID(),
			State:       int32(v.State()),
			WaitingInfo: v.WaitingInfo(),
			RunNumber:   v.RunNumber(),
			Held:        v.IsHeld(),
			Confirmed:   v.Confirmed(),
		}
		result = append(result, data)
	})

	sort.Sort(taskInfoSorter{result})

	return result
}

//Process - receive notification from dispatcher
func (pool *ActiveTaskPool) Process(receiver events.EventReceiver, routename events.RouteName, msg events.DispatchedMessage) {

	switch routename {
	case events.RouteTimeOut:
		{
			pool.log.Debug("task action message, route:", events.RouteTimeOut, "id:", msg.MsgID())

			msgdata, istype := msg.Message().(events.RouteTimeOutMsgFormat)
			if !istype {
				er := events.ErrUnrecognizedMsgFormat
				pool.log.Error(er)
				events.ResponseToReceiver(receiver, er)
				break
			}
			pool.log.Debug("Process time events")
			pool.ProcessTimeEvent(msgdata)
		}
	default:
		{
			err := events.ErrInvalidRouteName
			pool.log.Error(err)
			events.ResponseToReceiver(receiver, err)
		}
	}
}

//ProcessTimeEvent - entry point for processing tasks
func (pool *ActiveTaskPool) ProcessTimeEvent(data events.RouteTimeOutMsgFormat) {

	t := time.Date(data.Year, time.Month(data.Month), data.Day, data.Hour, data.Min, data.Sec, 0, time.Local)

	pool.cycleTasks(t)
}

//Start - starts tasks processing
func (pool *ActiveTaskPool) Start() error {

	defer pool.lock.Unlock()
	pool.lock.Lock()

	if pool.done == nil {
		pool.log.Info("starting store watch:", pool.isProcActive)
		pool.done = pool.tasks.watch(pool.activate, pool.shutdown)
	}

	pool.activate <- pool.isProcActive
	<-pool.done

	return nil
}

//Shutdown - shutdowns task pool
func (pool *ActiveTaskPool) Shutdown() error {

	defer pool.lock.Unlock()
	pool.lock.Lock()
	pool.isProcActive = false
	pool.shutdown <- struct{}{}
	<-pool.done

	return nil
}

//Resume - resumes tasks processing
func (pool *ActiveTaskPool) Resume() error {

	defer pool.lock.Unlock()
	pool.lock.Lock()

	pool.isProcActive = true
	pool.activate <- pool.isProcActive
	<-pool.done

	return nil
}

//Quiesce - puts taskpool into sleep mode
func (pool *ActiveTaskPool) Quiesce() error {

	defer pool.lock.Unlock()
	pool.lock.Lock()

	pool.isProcActive = false
	pool.activate <- pool.isProcActive

	<-pool.done

	return nil
}
