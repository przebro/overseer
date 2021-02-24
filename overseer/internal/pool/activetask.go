package pool

import (
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
type activeTask struct {
	taskdef.TaskDefinition
	holded     bool
	confirmed  bool
	orderID    unique.TaskOrderID
	orderDate  date.Odate
	tickets    []taskInTicket
	runNumber  int32
	executions []taskExecution
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

	task := &activeTask{orderID: orderID,
		TaskDefinition: definition,
		orderDate:      odate,
		runNumber:      0,
		executions:     []taskExecution{{ExecutionID: unique.NewID().Hex(), state: TaskStateWaiting}},
		tickets:        tickets,
		lock:           sync.RWMutex{},
		confirmed:      isconfirmed,
	}

	return task
}

//FromModel - creates an active task from model
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

	task := &activeTask{
		orderID:        unique.TaskOrderID(model.OrderID),
		TaskDefinition: def,
		orderDate:      model.OrderDate,
		executions:     execs,
		tickets:        tickets,
		lock:           sync.RWMutex{},
		confirmed:      model.Confirmed,
		holded:         model.Holded,
		waiting:        model.Waiting,
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
	task.executions[len(task.executions)-1].state = TaskStateHold
}
func (task *activeTask) Free() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.holded = false
	task.executions[len(task.executions)-1].state = TaskStateWaiting
}
func (task *activeTask) IsHeld() bool {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.holded
}

func (task *activeTask) getModel() activeTaskModel {
	defer task.lock.RUnlock()
	task.lock.RLock()
	def, _ := taskdef.SerializeDefinition(task.TaskDefinition)
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
	}
	for _, n := range task.tickets {
		t.Tickets = append(t.Tickets, taskInTicketModel{Name: n.name, Odate: n.odate, Fulfilled: n.fulfilled})
	}

	for _, n := range task.executions {
		t.Executions = append(t.Executions, taskExecutionModel{ID: n.ExecutionID, Worker: n.Worker, StartTime: n.Start, EndTime: n.End, State: n.state})
	}

	return t
}
