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

type activeTask struct {
	taskdef.TaskDefinition
	state     TaskState
	holded    bool
	confirmed bool
	orderID   unique.TaskOrderID
	orderDate date.Odate
	tickets   []taskInTicket
	runNumber int32
	worker    string
	waiting   string
	startTime time.Time
	endTime   time.Time
	output    []string
	lock      sync.RWMutex
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
		tickets:        tickets,
		state:          TaskStateWaiting,
		lock:           sync.RWMutex{},
		confirmed:      isconfirmed,
		output:         make([]string, 0),
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

	task := &activeTask{
		orderID:        unique.TaskOrderID(model.OrderID),
		TaskDefinition: def,
		orderDate:      model.OrderDate,
		runNumber:      model.RunNumber,
		tickets:        tickets,
		state:          model.State,
		startTime:      model.StartTime,
		endTime:        model.EndTime,
		lock:           sync.RWMutex{},
		confirmed:      model.Confirmed,
		holded:         model.Holded,
		output:         model.Output,
		waiting:        model.Waiting,
		worker:         model.Worker,
	}

	return task, nil

}

func (task *activeTask) State() TaskState {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.state
}
func (task *activeTask) SetState(state TaskState) {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.state = state
}
func (task *activeTask) SetRunNumber() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.runNumber++

}

func (task *activeTask) Tickets() []taskInTicket {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.tickets
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

func (task *activeTask) SetStartTime() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.startTime = time.Now()
}
func (task *activeTask) SetEndTime() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.endTime = time.Now()
}
func (task *activeTask) StartTime() time.Time {
	return task.startTime
}
func (task *activeTask) EndTime() time.Time {
	return task.endTime
}

func (task *activeTask) SetWorkerName(name string) {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.worker = name
}
func (task *activeTask) WorkerName() string {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.worker
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

func (task *activeTask) Output() []string {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.output

}
func (task *activeTask) AddOutput(data []string) {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.output = append(task.output, data...)
}

func (task *activeTask) Hold() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.holded = true
	task.state = TaskStateHold
}
func (task *activeTask) Free() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.holded = false
	task.state = TaskStateWaiting
}
func (task *activeTask) IsHeld() bool {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.holded

}

func (task *activeTask) getModel() activeTaskModel {
	defer task.lock.RUnlock()
	task.lock.RLock()
	def, _ := taskdef.WriteDefinitionFile(task.TaskDefinition)
	t := activeTaskModel{
		Definition: []byte(def),
		State:      task.state,
		Holded:     task.holded,
		Confirmed:  task.confirmed,
		OrderID:    string(task.orderID),
		OrderDate:  task.orderDate,
		Tickets:    []taskInTicketModel{},
		RunNumber:  task.runNumber,
		Worker:     task.worker,
		Waiting:    task.waiting,
		StartTime:  task.startTime,
		EndTime:    task.endTime,
		Output:     make([]string, len(task.output)),
	}

	copy(t.Output, task.output)
	for _, n := range task.tickets {
		t.Tickets = append(t.Tickets, taskInTicketModel{Name: n.name, Odate: n.odate, Fulfilled: n.fulfilled})
	}

	return t
}
