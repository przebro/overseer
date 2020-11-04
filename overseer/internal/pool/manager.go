package pool

import (
	"errors"
	"fmt"
	"goscheduler/common/logger"
	"goscheduler/overseer/internal/date"
	"goscheduler/overseer/internal/events"
	"goscheduler/overseer/internal/taskdef"
	"goscheduler/overseer/internal/unique"
	"strings"
	"sync"
)

//ActiveTaskPoolManager - Manages tasks in Active task pool
type ActiveTaskPoolManager struct {
	lock sync.RWMutex
	log  logger.AppLogger
	tdm  taskdef.TaskDefinitionManager
	pool *ActiveTaskPool
}

//NewActiveTaskPoolManager - Creates a new Managager
func NewActiveTaskPoolManager(dispatcher events.Dispatcher, tdm taskdef.TaskDefinitionManager, pool *ActiveTaskPool) *ActiveTaskPoolManager {
	manager := &ActiveTaskPoolManager{}
	manager.lock = sync.RWMutex{}
	manager.tdm = tdm
	manager.pool = pool
	manager.log = logger.Get()

	if dispatcher != nil {
		dispatcher.Subscribe(events.RouteTaskAct, manager)
		dispatcher.Subscribe(events.RouteChangeTaskState, manager)
	}

	return manager
}

//Order - Orders a new task, this method checks if all precoditions are met before it adds a new task
func (manager *ActiveTaskPoolManager) Order(group, name string, odate date.Odate) (string, error) {

	def, err := manager.tdm.GetTasks(taskdef.TaskData{Group: group, Name: name})
	if err != nil {
		return "", err
	}

	orderID, descr := manager.orderDefinition(def[0], odate)
	if descr != "" {
		return descr, errors.New("failed force task")
	}

	return string(orderID), nil
}

//Force - forces a new task, this method does not check for precondtions
func (manager *ActiveTaskPoolManager) Force(group, name string, odate date.Odate) (string, error) {

	def, err := manager.tdm.GetTasks(taskdef.TaskData{Group: group, Name: name})
	if err != nil {
		return "", err
	}

	orderID, descr := manager.forceDefinition(def[0], odate)
	if descr != "" {
		return descr, errors.New("failed force task")
	}

	return string(orderID), nil
}

//Rerun - Orders task again
func (manager *ActiveTaskPoolManager) Rerun(id unique.TaskOrderID) (string, error) {
	defer manager.lock.RUnlock()
	manager.lock.RLock()

	fmt.Println("TASK ID :", id)

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), err
	}

	state := task.State()
	if state == TaskStateEndedNotOk || state == TaskStateEndedOk {
		task.SetState(TaskStateWaiting)
	} else {
		return fmt.Sprintf("rerun task:%s failed, invalid status", id), nil
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
		return fmt.Sprintf("rerun task:%s failed, invalid status", id), nil
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
	task.Hold()

	return "", nil
}

//Free - Frees a holded task
func (manager *ActiveTaskPoolManager) Free(id unique.TaskOrderID) (string, error) {

	defer manager.lock.RUnlock()
	manager.lock.RLock()

	task, err := manager.pool.task(id)
	if err != nil {
		return fmt.Sprintf("task with id:%s does not exists", id), err
	}
	task.Free()

	return "", nil
}

func (manager *ActiveTaskPoolManager) orderNewTasks() {

	manager.log.Info("Ordering new tasks")

	groups := make([]string, 0)
	groups = append(groups, manager.tdm.GetGroups()...)
	taskData, _ := manager.tdm.GetTasksFromGroup(groups)

	tasks, err := manager.tdm.GetTasks(taskData...)
	if err != nil {
		manager.log.Error(err)
	}

	for _, t := range tasks {
		//It is a new day procedure so skip tasks that are ordered manually
		if t.OrderType() != taskdef.OrderingManual {
			manager.orderDefinition(t, manager.pool.currentOdate)
		}

	}
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
		currentOdate:     manager.pool.currentOdate,
		reason:           make([]string, 0),
		log:              manager.log,
	}
	n, g, _ := def.GetInfo()

	for ctx.state.processState(&ctx) {

	}

	if ctx.isSubmited == false {
		return "", strings.Join(ctx.reason, ",")

	}
	orderID := unique.NewOrderID()

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
		currentOdate:     manager.pool.currentOdate,
		reason:           make([]string, 0),
		log:              manager.log,
	}

	n, g, _ := def.GetInfo()

	for ctx.state.processState(&ctx) {

	}

	if ctx.isSubmited == false {
		return "", strings.Join(ctx.reason, ",")

	}
	orderID := unique.NewOrderID()

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
				er := errors.New("msg not in format")
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
				er := errors.New("msg not in format")
				manager.log.Error(er)
				events.ResponseToReceiver(receiver, er)
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
			err := errors.New("Invalid route name")
			manager.log.Debug(err)
			events.ResponseToReceiver(receiver, err)
		}
	}
}

func (manager *ActiveTaskPoolManager) changeTaskState(msg events.RouteChangeStateMsg) (events.RouteChangeStateResponseMsg, error) {

	var result string
	var err error
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

	def, err := manager.tdm.GetTasks(taskdef.TaskData{Name: msg.Name, Group: msg.Group})

	if err != nil {
		manager.log.Error(err)
		return result, err
	}

	if msg.Force {
		id, rmsg = manager.forceDefinition(def[0], date.Odate(msg.Odate))

	} else {
		id, rmsg = manager.orderDefinition(def[0], date.Odate(msg.Odate))
	}

	result.Data = make([]events.TaskInfoResultMsg, 1)
	result.Data[0].TaskID = id
	result.Data[0].WaitingInfo = rmsg
	if err != nil {
		result.Data[0].WaitingInfo = strings.Join([]string{result.Data[0].WaitingInfo, err.Error()}, ",")
	}
	result.Data[0].Name, result.Data[0].Group, _ = def[0].GetInfo()

	return result, nil
}
