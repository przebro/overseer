package pool

import (
	"overseer/common/types"
	"overseer/common/types/date"
	"overseer/overseer/internal/taskdef"
	"overseer/overseer/internal/unique"
	"sync"
	"time"
)

//TaskState - current state of a task
type TaskState int32

//Possible states of an active task
const (
	//Task waits for a time window, tickets
	TaskStateWaiting  TaskState = 1
	TaskStateStarting TaskState = 2
	//Time and tickets prerequisites were met, task was sent to worker. Task may wait for confirm or flag or being executed
	TaskStateExecuting  TaskState = 3
	TaskStateEndedOk    TaskState = 4
	TaskStateEndedNotOk TaskState = 5
	TaskStateHold       TaskState = 6
)

type taskInTicket struct {
	name      string
	odate     string
	fulfilled bool
}
type taskExecution struct {
	ExecutionID string
	Worker      string
	Start       time.Time
	End         time.Time
	state       TaskState
}

type taskCycle struct {
	IsCyclic         bool
	NextRun          types.HourMinTime
	RunFrom          string
	MaxRun           int
	CurrentRunNumber int
	RunInterval      int
}

type activeTask struct {
	taskdef.TaskDefinition
	holded     bool
	confirmed  bool
	orderID    unique.TaskOrderID
	orderDate  date.Odate
	tickets    []taskInTicket
	runNumber  int32
	executions []taskExecution
	cycle      taskCycle
	waiting    string
	lock       sync.RWMutex
	collected  []taskInTicket
}

func newActiveTask(orderID unique.TaskOrderID, odate date.Odate, definition taskdef.TaskDefinition) *activeTask {

	tickets := make([]taskInTicket, 0)

	for _, e := range definition.TicketsIn() {

		realOdat := calcRealOdate(odate, e.Odate, definition.Calendar())
		tickets = append(tickets, taskInTicket{odate: string(realOdat), name: e.Name, fulfilled: false})
	}

	isconfirmed := func() bool {
		if definition.Confirm() {
			return false
		}
		return true
	}()

	cdata := definition.Cyclic()
	cycle := taskCycle{IsCyclic: cdata.IsCycle, MaxRun: cdata.MaxRuns, CurrentRunNumber: 0, RunFrom: string(cdata.RunFrom), RunInterval: cdata.TimeInterval}
	if cycle.IsCyclic {
		cycle.NextRun = types.Now()
	}

	task := &activeTask{orderID: orderID,
		TaskDefinition: definition,
		orderDate:      odate,
		runNumber:      0,
		executions:     []taskExecution{{ExecutionID: unique.NewID().Hex(), state: TaskStateWaiting}},
		tickets:        tickets,
		cycle:          cycle,
		lock:           sync.RWMutex{},
		confirmed:      isconfirmed,
	}

	return task
}

//fromModel - creates an active task from model
func fromModel(model activeTaskModel) (*activeTask, error) {

	def, err := taskdef.FromString(string(model.Definition))
	if err != nil {
		return nil, err
	}

	tickets := []taskInTicket{}
	if model.Tickets != nil {
		for _, n := range model.Tickets {
			tickets = append(tickets, taskInTicket{name: n.Name, odate: n.Odate, fulfilled: n.Fulfilled})
		}
	}

	execs := []taskExecution{}
	for _, n := range model.Executions {
		execs = append(execs, taskExecution{ExecutionID: n.ID, Start: n.StartTime, End: n.EndTime, Worker: n.Worker, state: n.State})
	}

	cycle := taskCycle{IsCyclic: model.Cycle.IsCyclic,
		NextRun:          types.HourMinTime(model.Cycle.NextRun),
		MaxRun:           model.Cycle.MaxRun,
		CurrentRunNumber: model.Cycle.Current,
		RunFrom:          model.Cycle.RunFrom,
	}

	task := &activeTask{
		orderID:        unique.TaskOrderID(model.OrderID),
		TaskDefinition: def,
		orderDate:      model.OrderDate,
		executions:     execs,
		tickets:        tickets,
		lock:           sync.RWMutex{},
		confirmed:      model.Confirmed,
		holded:         model.Holded,
		cycle:          cycle,
		waiting:        model.Waiting,
		runNumber:      model.RunNumber,
		collected:      []taskInTicket{},
	}

	return task, nil
}

func (task *activeTask) State() TaskState {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.executions[len(task.executions)-1].state
}
func (task *activeTask) SetState(state TaskState) {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.executions[len(task.executions)-1].state = state

}
func (task *activeTask) SetRunNumber() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.runNumber++

}
func (task *activeTask) SetExecutionID() {
	defer task.lock.Unlock()
	task.lock.Lock()
	id := unique.NewID().Hex()
	task.executions = append(task.executions, taskExecution{ExecutionID: id, state: TaskStateWaiting})
}

func (task *activeTask) Tickets() []taskInTicket {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.tickets
}

func (task *activeTask) Collected() []taskInTicket {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.collected
}

func (task *activeTask) OrderID() unique.TaskOrderID {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.orderID
}

func (task *activeTask) OrderDate() date.Odate {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.orderDate
}
func (task *activeTask) RunNumber() int32 {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.runNumber
}
func (task *activeTask) CurrentExecutionID() string {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.executions[len(task.executions)-1].ExecutionID
}
func (task *activeTask) Confirmed() bool {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.confirmed
}

func (task *activeTask) SetConfirm() bool {
	defer task.lock.RUnlock()
	task.lock.RLock()

	if task.confirmed == true {
		return false
	}

	task.confirmed = true
	return true
}

func (task *activeTask) SetStartTime() time.Time {
	defer task.lock.Unlock()
	task.lock.Lock()
	stime := time.Now()
	task.executions[len(task.executions)-1].Start = stime
	return stime
}
func (task *activeTask) SetEndTime() time.Time {
	defer task.lock.Unlock()
	task.lock.Lock()
	etime := time.Now()
	task.executions[len(task.executions)-1].End = etime
	return etime
}
func (task *activeTask) StartTime() time.Time {
	return task.executions[len(task.executions)-1].Start
}
func (task *activeTask) EndTime() time.Time {
	return task.executions[len(task.executions)-1].End
}

func (task *activeTask) SetWorkerName(name string) {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.executions[len(task.executions)-1].Worker = name
}
func (task *activeTask) WorkerName() string {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.executions[len(task.executions)-1].Worker
}
func (task *activeTask) WaitingInfo() string {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.waiting
}

func (task *activeTask) SetWaitingInfo(info string) {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.waiting = info
}

func (task *activeTask) Hold() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.holded = true
}
func (task *activeTask) Free() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.holded = false
}
func (task *activeTask) IsHeld() bool {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.holded
}

func (task *activeTask) IsCyclic() bool {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.cycle.IsCyclic
}

func (task *activeTask) CycleData() taskCycle {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.cycle
}

func (task *activeTask) getModel() activeTaskModel {
	defer task.lock.RUnlock()
	task.lock.RLock()
	def, _ := taskdef.SerializeDefinition(task.TaskDefinition)

	cycle := taskCycleModel{
		IsCyclic: task.cycle.IsCyclic,
		NextRun:  string(task.cycle.NextRun),
		MaxRun:   task.cycle.MaxRun,
		Current:  task.cycle.CurrentRunNumber,
		RunFrom:  task.cycle.RunFrom,
	}

	t := activeTaskModel{
		Definition: []byte(def),
		Holded:     task.holded,
		Confirmed:  task.confirmed,
		OrderID:    string(task.orderID),
		OrderDate:  task.orderDate,
		Tickets:    []taskInTicketModel{},
		RunNumber:  task.runNumber,
		Waiting:    task.waiting,
		Executions: []taskExecutionModel{},
		Cycle:      cycle,
	}
	for _, n := range task.tickets {
		t.Tickets = append(t.Tickets, taskInTicketModel{Name: n.name, Odate: n.odate, Fulfilled: n.fulfilled})
	}

	for _, n := range task.executions {
		t.Executions = append(t.Executions, taskExecutionModel{ID: n.ExecutionID, Worker: n.Worker, StartTime: n.Start, EndTime: n.End, State: n.state})
	}

	return t
}
