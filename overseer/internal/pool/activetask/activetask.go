package activetask

import (
	"sync"
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/overseer/internal/pool/calc"
	"github.com/przebro/overseer/overseer/internal/pool/models"
	"github.com/przebro/overseer/overseer/internal/pool/readers"
	"github.com/przebro/overseer/overseer/internal/taskdef"
)

type TaskInstance struct {
	*taskdef.TaskDefinition
	reference  string
	holded     bool
	confirmed  bool
	orderID    unique.TaskOrderID
	orderDate  date.Odate
	executions []models.TaskExecution
	lock       sync.RWMutex
	nextRun    types.HourMinTime
}

func NewActiveTask(orderID unique.TaskOrderID, odate date.Odate, definition *taskdef.TaskDefinition, refID unique.MsgID) *TaskInstance {

	isconfirmed := !definition.Confirm

	cdata := definition.Cyclic

	var nr types.HourMinTime

	if cdata.IsCycle {
		nr = types.Now()
	}

	task := &TaskInstance{orderID: orderID,
		TaskDefinition: definition,
		reference:      refID.Hex(),
		orderDate:      odate,
		executions:     []models.TaskExecution{{ExecutionID: unique.NewID().Hex(), State: models.TaskStateWaiting}},
		lock:           sync.RWMutex{},
		confirmed:      isconfirmed,
		nextRun:        nr,
	}

	return task
}

// fromModel - creates an active task from model
func FromModel(model ActiveTaskModel, rdr readers.ActiveDefinitionReader) (*TaskInstance, error) {

	def, err := rdr.GetActiveDefinition(model.Reference)
	if err != nil {
		return nil, err
	}

	execs := []models.TaskExecution{}
	for _, n := range model.Executions {
		execs = append(execs, models.TaskExecution{ExecutionID: n.ID, Start: n.StartTime, End: n.EndTime, Worker: n.Worker, State: n.State})
	}

	cycle := models.TaskCycle{IsCyclic: model.Cycle.IsCyclic,
		NextRun: types.HourMinTime(model.Cycle.NextRun),
		MaxRun:  model.Cycle.MaxRun,
		RunFrom: model.Cycle.RunFrom,
	}

	task := &TaskInstance{
		orderID:        unique.TaskOrderID(model.OrderID),
		reference:      model.Reference,
		TaskDefinition: def,
		orderDate:      model.OrderDate,
		executions:     execs,
		lock:           sync.RWMutex{},
		confirmed:      model.Confirmed,
		holded:         model.Holded,
		nextRun:        cycle.NextRun,
	}

	return task, nil
}

func (task *TaskInstance) State() models.TaskState {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.executions[len(task.executions)-1].State
}
func (task *TaskInstance) SetState(state models.TaskState) {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.executions[len(task.executions)-1].State = state

}

func (task *TaskInstance) SetExecutionID() {
	defer task.lock.Unlock()
	task.lock.Lock()
	id := unique.NewID().Hex()
	task.executions = append(task.executions, models.TaskExecution{ExecutionID: id, State: models.TaskStateWaiting})
}

func (task *TaskInstance) Tickets(odate date.Odate) []models.TaskInTicket {
	defer task.lock.RUnlock()
	task.lock.RLock()

	tickets := make([]models.TaskInTicket, 0)

	for _, e := range task.InTickets {

		realOdat := calc.CalcRealOdate(odate, e.Odate, task.Schedule)
		tickets = append(tickets, models.TaskInTicket{Odate: string(realOdat), Name: e.Name, Fulfilled: false})
	}

	return tickets
}

func (task *TaskInstance) OrderID() unique.TaskOrderID {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.orderID
}

func (task *TaskInstance) OrderDate() date.Odate {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.orderDate
}
func (task *TaskInstance) RunNumber() int32 {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.getRunNumber()
}

func (task *TaskInstance) getRunNumber() int32 {
	execInfo := task.executions[len(task.executions)-1]
	if execInfo.Start.IsZero() {
		return int32(len(task.executions) - 1)
	}

	return int32(len(task.executions))
}

func (task *TaskInstance) ExecutionID() string {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.executions[len(task.executions)-1].ExecutionID
}

func (task *TaskInstance) Confirmed() bool {
	defer task.lock.RUnlock()
	task.lock.RLock()
	return task.confirmed
}

func (task *TaskInstance) SetConfirm() bool {
	defer task.lock.RUnlock()
	task.lock.RLock()

	if task.confirmed {
		return false
	}

	task.confirmed = true
	return true
}

func (task *TaskInstance) SetStartTime() time.Time {
	defer task.lock.Unlock()
	task.lock.Lock()
	stime := time.Now()
	task.executions[len(task.executions)-1].Start = stime
	return stime
}
func (task *TaskInstance) SetEndTime() time.Time {
	defer task.lock.Unlock()
	task.lock.Lock()
	etime := time.Now()
	task.executions[len(task.executions)-1].End = etime
	return etime
}
func (task *TaskInstance) StartTime() time.Time {
	return task.executions[len(task.executions)-1].Start
}
func (task *TaskInstance) EndTime() time.Time {
	return task.executions[len(task.executions)-1].End
}

func (task *TaskInstance) SetWorkerName(name string) {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.executions[len(task.executions)-1].Worker = name
}
func (task *TaskInstance) WorkerName() string {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.executions[len(task.executions)-1].Worker
}

func (task *TaskInstance) Hold() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.holded = true
}
func (task *TaskInstance) Free() {
	defer task.lock.Unlock()
	task.lock.Lock()
	task.holded = false
}
func (task *TaskInstance) IsHeld() bool {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.holded
}

func (task *TaskInstance) IsCyclic() bool {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.Cyclic.IsCycle
}

func (task *TaskInstance) NextRun() types.HourMinTime {
	defer task.lock.RUnlock()
	task.lock.RLock()

	return task.nextRun
}

func (task *TaskInstance) PrepareNextCycle() bool {
	return task.prepareNextCycle()
}
func (task *TaskInstance) prepareNextCycle() bool {

	defer task.lock.Unlock()
	task.lock.Lock()

	if !task.Cyclic.IsCycle {
		return false
	}

	if int(task.getRunNumber()) >= task.Cyclic.MaxRuns {
		return false
	}

	var tm time.Time

	switch task.Cyclic.RunFrom {
	case "start":
		{

			tm = task.executions[len(task.executions)-1].Start.Add(time.Duration(task.Cyclic.TimeInterval) * time.Minute)
		}
	case "end":
		{
			tm = task.executions[len(task.executions)-1].End.Add(time.Duration(task.Cyclic.TimeInterval) * time.Minute)
		}
	case "schedule":
		{
			//:TODO for now it acts like from end
			tm = task.executions[len(task.executions)-1].End.Add(time.Duration(task.Cyclic.TimeInterval) * time.Minute)
		}
	}

	task.nextRun = types.FromTime(tm)

	return true
}

func (task *TaskInstance) GetModel() ActiveTaskModel {
	defer task.lock.RUnlock()
	task.lock.RLock()

	cycle := taskCycleModel{
		IsCyclic: task.Cyclic.IsCycle,
		NextRun:  string(task.nextRun),
		MaxRun:   task.Cyclic.MaxRuns,
		RunFrom:  string(task.Cyclic.RunFrom),
	}

	n, g := task.Definition.Name, task.Definition.Group

	t := ActiveTaskModel{
		Name:       n,
		Group:      g,
		Reference:  task.reference,
		Holded:     task.holded,
		Confirmed:  task.confirmed,
		OrderID:    string(task.orderID),
		OrderDate:  task.orderDate,
		Tickets:    []taskInTicketModel{},
		Executions: []taskExecutionModel{},
		Cycle:      cycle,
	}
	// for _, tt := range task.tickets {
	// 	t.Tickets = append(t.Tickets, taskInTicketModel{Name: tt.Name, Odate: tt.Odate, Fulfilled: tt.Fulfilled})
	// }

	for _, n := range task.executions {
		t.Executions = append(t.Executions, taskExecutionModel{ID: n.ExecutionID, Worker: n.Worker, StartTime: n.Start, EndTime: n.End, State: n.State})
	}

	return t
}
