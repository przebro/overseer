package pool

import (
	"sync"
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/przebro/overseer/overseer/internal/unique"
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
	IsCyclic    bool
	NextRun     types.HourMinTime
	RunFrom     string
	MaxRun      int
	RunInterval int
}

type activeTask struct {
	taskdef.TaskDefinition
	reference  string
	holded     bool
	confirmed  bool
	orderID    unique.TaskOrderID
	orderDate  date.Odate
	tickets    []taskInTicket
	executions []taskExecution
	cycle      taskCycle
	waiting    string
	lock       sync.RWMutex
	collected  []taskInTicket
}

//ActiveDefinitionReader - reads definitions
type ActiveDefinitionReader interface {
	GetActiveDefinition(refID string) (taskdef.TaskDefinition, error)
}

//ActiveDefinitionWriter - writes definitions
type ActiveDefinitionWriter interface {
	WriteActiveDefinition(def taskdef.TaskDefinition, id unique.MsgID) error
}

//ActiveDefinitionRemover - removes definitions
type ActiveDefinitionRemover interface {
	RemoveActiveDefinition(id string) error
}

//ActiveDefinitionReadWriter - groups definition reader and writer
type ActiveDefinitionReadWriter interface {
	ActiveDefinitionReader
	ActiveDefinitionWriter
}

//ActiveDefinitionReadWriter - groups definition reader, writer and remover
type ActiveDefinitionReadWriterRemover interface {
	ActiveDefinitionReader
	ActiveDefinitionWriter
	ActiveDefinitionRemover
}

func newActiveTask(orderID unique.TaskOrderID, odate date.Odate, definition taskdef.TaskDefinition, refID unique.MsgID) *activeTask {

	tickets := make([]taskInTicket, 0)

	for _, e := range definition.TicketsIn() {

		realOdat := calcRealOdate(odate, e.Odate, definition.Calendar())
		tickets = append(tickets, taskInTicket{odate: string(realOdat), name: e.Name, fulfilled: false})
	}

	isconfirmed := !definition.Confirm()

	cdata := definition.Cyclic()
	cycle := taskCycle{IsCyclic: cdata.IsCycle, MaxRun: cdata.MaxRuns, RunFrom: string(cdata.RunFrom), RunInterval: cdata.TimeInterval}
	if cycle.IsCyclic {
		cycle.NextRun = types.Now()
	}

	task := &activeTask{orderID: orderID,
		TaskDefinition: definition,
		reference:      refID.Hex(),
		orderDate:      odate,
		executions:     []taskExecution{{ExecutionID: unique.NewID().Hex(), state: TaskStateWaiting}},
		tickets:        tickets,
		cycle:          cycle,
		lock:           sync.RWMutex{},
		confirmed:      isconfirmed,
	}

	return task
}

//fromModel - creates an active task from model
func fromModel(model activeTaskModel, rdr ActiveDefinitionReader) (*activeTask, error) {

	def, err := rdr.GetActiveDefinition(model.Reference)
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
		NextRun: types.HourMinTime(model.Cycle.NextRun),
		MaxRun:  model.Cycle.MaxRun,
		RunFrom: model.Cycle.RunFrom,
	}

	task := &activeTask{
		orderID:        unique.TaskOrderID(model.OrderID),
		reference:      model.Reference,
		TaskDefinition: def,
		orderDate:      model.OrderDate,
		executions:     execs,
		tickets:        tickets,
		lock:           sync.RWMutex{},
		confirmed:      model.Confirmed,
		holded:         model.Holded,
		cycle:          cycle,
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
	return task.getRunNumber()
}

func (task *activeTask) getRunNumber() int32 {
	execInfo := task.executions[len(task.executions)-1]
	if execInfo.Start.IsZero() {
		return int32(len(task.executions) - 1)
	}

	return int32(len(task.executions))
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

	if task.confirmed {
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

func (task *activeTask) prepareNextCycle() bool {

	defer task.lock.Unlock()
	task.lock.Lock()

	if !task.cycle.IsCyclic {
		return false
	}

	if int(task.getRunNumber()) >= task.cycle.MaxRun {
		return false
	}

	var tm time.Time

	switch task.cycle.RunFrom {
	case "start":
		{

			tm = task.executions[len(task.executions)-1].Start.Add(time.Duration(task.cycle.RunInterval) * time.Minute)
		}
	case "end":
		{
			tm = task.executions[len(task.executions)-1].End.Add(time.Duration(task.cycle.RunInterval) * time.Minute)
		}
	case "schedule":
		{
			//:TODO for now it acts like from end
			tm = task.executions[len(task.executions)-1].End.Add(time.Duration(task.cycle.RunInterval) * time.Minute)
		}
	}

	task.cycle.NextRun = types.FromTime(tm)

	return true
}

func (task *activeTask) getModel() activeTaskModel {
	defer task.lock.RUnlock()
	task.lock.RLock()

	cycle := taskCycleModel{
		IsCyclic: task.cycle.IsCyclic,
		NextRun:  string(task.cycle.NextRun),
		MaxRun:   task.cycle.MaxRun,
		RunFrom:  task.cycle.RunFrom,
	}

	n, g, _ := task.GetInfo()

	t := activeTaskModel{
		Name:       n,
		Group:      g,
		Reference:  task.reference,
		Holded:     task.holded,
		Confirmed:  task.confirmed,
		OrderID:    string(task.orderID),
		OrderDate:  task.orderDate,
		Tickets:    []taskInTicketModel{},
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
