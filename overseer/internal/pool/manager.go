package pool

import (
	"errors"
	"fmt"
	"overseer/common/logger"
	"overseer/common/types/date"
	"overseer/datastore"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/taskdef"
	"overseer/overseer/internal/unique"
	"overseer/overseer/taskdata"
	"strings"
	"sync"
)

var (
	//ErrUnableFindDef - a definition was not found
	ErrUnableFindDef = errors.New("unable to find definition")
	//ErrInvalidStatus - invalid status
	ErrInvalidStatus = errors.New("Invalid status")
)

//ActiveTaskPoolManager - Manages tasks in Active task pool
type ActiveTaskPoolManager struct {
	lock     sync.RWMutex
	log      logger.AppLogger
	tdm      taskdef.TaskDefinitionManager
	pool     *ActiveTaskPool
	sequence SequenceGenerator
}

//NewActiveTaskPoolManager - Creates a new Managager
func NewActiveTaskPoolManager(dispatcher events.Dispatcher,
	tdm taskdef.TaskDefinitionManager,
	pool *ActiveTaskPool,
	provider *datastore.Provider) (*ActiveTaskPoolManager, error) {

	var err error
	manager := &ActiveTaskPoolManager{}
	manager.lock = sync.RWMutex{}
	manager.tdm = tdm
	manager.pool = pool
	manager.log = logger.Get()

	if manager.sequence, err = NewSequenceGenerator("sequence", provider); err != nil {
		return nil, err
	}

	if dispatcher != nil {
		dispatcher.Subscribe(events.RouteTaskAct, manager)
		dispatcher.Subscribe(events.RouteChangeTaskState, manager)
	}

	return manager, nil
}

//Order - Orders a new task, this method checks if all precoditions are met before it adds a new task
func (manager *ActiveTaskPoolManager) Order(task taskdata.GroupNameData, odate date.Odate) (string, error) {

	var err error
	var definition taskdef.TaskDefinition
	if definition, err = manager.tdm.GetTask(task); err != nil {
		return "", ErrUnableFindDef
	}

	orderID, descr := manager.orderDefinition(definition, odate)
	if descr != "" {
		return descr, errors.New("failed order task")
	}

	return string(orderID), nil
}

//Force - forces a new task, this method does not check for precondtions
func (manager *ActiveTaskPoolManager) Force(task taskdata.GroupNameData, odate date.Odate) (string, error) {

	var err error
	var definition taskdef.TaskDefinition
	if definition, err = manager.tdm.GetTask(task); err != nil {
		return "", ErrUnableFindDef
	}

	orderID, descr := manager.forceDefinition(definition, odate)
	if descr != "" {
		return descr, errors.New("failed force task")
	}

	return string(orderID), nil
}

//Rerun - Orders task again
func (manager *ActiveTaskPoolManager) Rerun(id unique.TaskOrderID) (string, error) {
	defer manager.lock.RUnlock()
	manager.lock.RLock()

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), err
	}

	state := task.State()
	if state == TaskStateEndedNotOk || state == TaskStateEndedOk {
		task.SetState(TaskStateWaiting)
	} else {
		return fmt.Sprintf("rerun task:%s failed, invalid status", id), ErrInvalidStatus
	}

	return fmt.Sprintf("rerun task:%s ok", id), nil

}

//SetOk - Sets task to EndedOk status
func (manager *ActiveTaskPoolManager) SetOk(id unique.TaskOrderID) (string, error) {
	defer manager.lock.RUnlock()
	manager.lock.RLock()

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), err
	}

	state := task.State()
	if state == TaskStateEndedNotOk {
		task.SetState(TaskStateEndedOk)
	} else {
		return fmt.Sprintf("set to OK task:%s failed, invalid status", id), ErrInvalidStatus
	}

	return fmt.Sprintf("set to ok task:%s ok", id), nil
}

//Hold - Holds the taskdef. It will be not processed durning cycle
func (manager *ActiveTaskPoolManager) Hold(id unique.TaskOrderID) (string, error) {

	defer manager.lock.RUnlock()
	manager.lock.RLock()

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), err
	}

	if task.IsHeld() {
		return fmt.Sprintf("hold task:%s failed, invalid status", id), ErrInvalidStatus
	}

	task.Hold()

	return fmt.Sprintf("hold task:%s ok", id), nil
}

//Free - Frees a holded task
func (manager *ActiveTaskPoolManager) Free(id unique.TaskOrderID) (string, error) {

	defer manager.lock.RUnlock()
	manager.lock.RLock()

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), err
	}

	if !task.IsHeld() {
		return fmt.Sprintf("free task:%s failed, task is not held", id), ErrInvalidStatus
	}

	task.Free()

	return fmt.Sprintf("free task:%s ok", id), nil
}

//Confirm - Manually Confirms a task
func (manager *ActiveTaskPoolManager) Confirm(id unique.TaskOrderID) (string, error) {

	defer manager.lock.RUnlock()
	manager.lock.RLock()

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), err
	}

	result := task.SetConfirm()
	if result == false {
		return fmt.Sprintf("task with id:%s already confirmed", id), fmt.Errorf("task confirmed")
	}

	return fmt.Sprintf("confirm task:%s ok", id), nil
}

func (manager *ActiveTaskPoolManager) orderNewTasks() int {

	ordered := 0
	manager.log.Info("Ordering new tasks")

	groups := make([]string, 0)
	groups = append(groups, manager.tdm.GetGroups()...)
	taskData, _ := manager.tdm.GetTasksFromGroup(groups)

	result := manager.tdm.GetTasks(taskData...)

	for _, t := range result {

		//It is a new day procedure so skip tasks that are ordered manually
		if t.OrderType() != taskdef.OrderingManual {
			manager.orderDefinition(t, date.CurrentOdate())
			ordered++
		}
	}
	return ordered
}

//orderDefinition - Adds a new task to the Active Task Pool
//this method performs all checks
func (manager *ActiveTaskPoolManager) orderDefinition(def taskdef.TaskDefinition, odate date.Odate) (unique.TaskOrderID, string) {

	defer manager.lock.Unlock()
	manager.lock.Lock()

	manager.log.Debug("order:", def, ":", odate)

	ctx := TaskOrderContext{def: def,
		ignoreCalendar:   false,
		ignoreSubmission: false,
		odate:            odate,
		state:            &ostateCheckOtype{},
		currentOdate:     date.CurrentOdate(),
		reason:           make([]string, 0),
		log:              manager.log,
	}
	n, g, _ := def.GetInfo()

	for ctx.state.processState(&ctx) {

	}

	if ctx.isSubmited == false {
		return "", strings.Join(ctx.reason, ",")

	}
	orderID := manager.sequence.Next()

	task := newActiveTask(orderID, odate, def)
	manager.pool.addTask(orderID, task)

	manager.log.Info(fmt.Sprintf("Task %s from gorup %s ordered with id:%s odate:%s", n, g, orderID, odate))

	return orderID, strings.Join(ctx.reason, ",")
}

//forceDefinition - Forcefully adds a new task to the Active Task Pool
//this method ignores all checks
func (manager *ActiveTaskPoolManager) forceDefinition(def taskdef.TaskDefinition, odate date.Odate) (unique.TaskOrderID, string) {

	defer manager.lock.Unlock()
	manager.lock.Lock()

	manager.log.Info("force:", def, ":", odate)

	ctx := TaskOrderContext{def: def,
		ignoreCalendar:   true,
		ignoreSubmission: true,
		odate:            odate,
		state:            &ostateCheckOtype{},
		currentOdate:     date.CurrentOdate(),
		reason:           make([]string, 0),
		log:              manager.log,
	}

	n, g, _ := def.GetInfo()

	for ctx.state.processState(&ctx) {

	}

	if ctx.isSubmited == false {
		return "", strings.Join(ctx.reason, ",")

	}
	orderID := manager.sequence.Next()

	task := newActiveTask(orderID, odate, def)
	manager.pool.addTask(orderID, task)

	manager.log.Info(fmt.Sprintf("Task %s from gorup %s forced with id:%s odate:%s", n, g, orderID, odate))

	return orderID, strings.Join(ctx.reason, ",")
}

//Process - receive notification from dispatcher
func (manager *ActiveTaskPoolManager) Process(receiver events.EventReceiver, routename events.RouteName, msg events.DispatchedMessage) {

	switch routename {
	case events.RouteTaskAct:
		{
			manager.log.Debug("task action message, route:", events.RouteTaskAct, "id:", msg.MsgID())
			addmsg, istype := msg.Message().(events.RouteTaskActionMsgFormat)
			if !istype {
				er := events.ErrUnrecognizedMsgFormat
				manager.log.Error(er)
				events.ResponseToReceiver(receiver, er)
				break
			}

			result, err := manager.processAddToActivePool(addmsg)
			if err != nil {
				events.ResponseToReceiver(receiver, err)
				break
			}

			events.ResponseToReceiver(receiver, result)
		}
	case events.RouteChangeTaskState:
		{
			manager.log.Debug("task action message, route:", events.RouteChangeTaskState)
			actmsg, istype := msg.Message().(events.RouteChangeStateMsg)
			if !istype {
				er := events.ErrUnrecognizedMsgFormat
				manager.log.Error(er)
				events.ResponseToReceiver(receiver, er)
				break
			}

			result, err := manager.changeTaskState(actmsg)
			if err != nil {
				events.ResponseToReceiver(receiver, err)
				break
			}
			events.ResponseToReceiver(receiver, result)
		}
	default:
		{
			err := events.ErrInvalidRouteName
			manager.log.Debug(err)
			events.ResponseToReceiver(receiver, err)
		}
	}
}

func (manager *ActiveTaskPoolManager) changeTaskState(msg events.RouteChangeStateMsg) (events.RouteChangeStateResponseMsg, error) {

	var result string
	var err error
	test := 1

	tab := map[bool]int{true: 1, false: 0}
	//One and only one flag can be set
	test = test << tab[msg.Free] << tab[msg.Hold] << tab[msg.Rerun] << tab[msg.SetOK]

	if test != 2 {
		err = errors.New("invalid flag combination")
		return events.RouteChangeStateResponseMsg{Message: result, OrderID: msg.OrderID}, err
	}

	switch true {
	case msg.Free:
		{
			result, err = manager.Free(msg.OrderID)
		}
	case msg.Hold:
		{
			result, err = manager.Hold(msg.OrderID)
		}
	case msg.SetOK:
		{
			result, err = manager.SetOk(msg.OrderID)
		}
	case msg.Rerun:
		{
			result, err = manager.Rerun(msg.OrderID)
		}
	}
	return events.RouteChangeStateResponseMsg{Message: result, OrderID: msg.OrderID}, err

}

func (manager *ActiveTaskPoolManager) processAddToActivePool(msg events.RouteTaskActionMsgFormat) (events.RouteTaskActionResponseFormat, error) {

	var rmsg string
	var id unique.TaskOrderID
	var result events.RouteTaskActionResponseFormat
	var definition taskdef.TaskDefinition
	var err error

	if definition, err = manager.tdm.GetTask(taskdata.GroupNameData{Name: msg.Name, Group: msg.Group}); err != nil {
		return result, ErrUnableFindDef

	}

	if msg.Force {
		id, rmsg = manager.forceDefinition(definition, date.Odate(msg.Odate))

	} else {
		id, rmsg = manager.orderDefinition(definition, date.Odate(msg.Odate))
	}

	result.Data = make([]events.TaskInfoResultMsg, 1)
	result.Data[0].TaskID = id
	result.Data[0].WaitingInfo = rmsg
	result.Data[0].Name, result.Data[0].Group, _ = definition.GetInfo()

	return result, nil
}
