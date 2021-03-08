package work

import (
	"overseer/common/logger"
	"overseer/common/types"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/unique"
	"sync"
	"time"
)

type workerManager struct {
	askChannel    chan taskGetStatusMsg
	resultChannel chan events.RouteWorkResponseMsg
	launchChannel chan taskExecuteMsg
	cleanChannel  chan taskCleanMsg
	workStatus    chan struct{}

	log     logger.AppLogger
	workers map[string]WorkerMediator
	status  map[string]events.RouteWorkResponseMsg
	lock    sync.Mutex
}

//WorkerManager - Manages actions between
type WorkerManager interface {
	Run() error
}

//NewWorkerManager - Creates a new WorkerManager
func NewWorkerManager(d events.Dispatcher, conf config.WorkerManagerConfiguration) WorkerManager {

	w := &workerManager{}
	w.log = logger.Get()
	w.askChannel = make(chan taskGetStatusMsg)
	w.resultChannel = make(chan events.RouteWorkResponseMsg)
	w.launchChannel = make(chan taskExecuteMsg)
	w.cleanChannel = make(chan taskCleanMsg)
	w.workers = make(map[string]WorkerMediator)
	w.workStatus = make(chan struct{})
	w.lock = sync.Mutex{}
	w.status = map[string]events.RouteWorkResponseMsg{}

	if d != nil {
		d.Subscribe(events.RouteWorkLaunch, w)
		d.Subscribe(events.RouteWorkCheck, w)
		d.Subscribe(events.RouteTimeOut, w)
		d.Subscribe(events.RouteTaskClean, w)
	}

	for _, n := range conf.Workers {
		w.log.Info("Creating service worker:", n.WorkerName, ",", n.WorkerHost, ":", n.WorkerPort)
		sworker := NewWorkerMediator(n, conf.Timeout, w.resultChannel)
		w.workers[n.WorkerName] = sworker
	}

	go func() {
		w.updateWorkers(conf.WorkerInterval)
	}()

	return w
}
func (w *workerManager) Run() error {

	go func() {
		for {
			select {
			case msg := <-w.launchChannel:
				{
					//Task request to process a work on remote worker
					result := w.startTask(msg.data)
					events.ResponseToReceiver(msg.receiver, result)
				}
			case msg := <-w.askChannel:
				{
					//task asking for status or sending an information that remote task should be cleaned
					result := w.getTaskStatus(msg.workername, msg.ExecutionID, msg.orderID)
					events.ResponseToReceiver(msg.receiver, result)
				}
			case msg := <-w.resultChannel:
				{
					//An information from worker has been returned, update status
					w.updateTaskStatus(msg)
				}
			case msg := <-w.cleanChannel:
				{
					w.cleanupTask(msg.workername, msg.orderID, msg.executionID, msg.terminate)
				}
			case <-w.workStatus:
				{
					//time out,request workers for actual task statuses
					w.requestTaskStatus()
				}
			}
		}

	}()
	return nil
}

func (w *workerManager) startTask(msg events.RouteTaskExecutionMsg) events.RouteWorkResponseMsg {

	var response events.RouteWorkResponseMsg
	var wname string

	if wname = w.getWorker(); wname == "" {
		response = events.RouteWorkResponseMsg{
			Status:      types.WorkerTaskStatusFailed,
			OrderID:     msg.OrderID,
			ExecutionID: msg.ExecutionID,
			WorkerName:  "",
		}
		return response
	}

	defer w.lock.Unlock()
	w.lock.Lock()

	response = events.RouteWorkResponseMsg{
		Status:      types.WorkerTaskStatusStarting,
		OrderID:     msg.OrderID,
		ExecutionID: msg.ExecutionID,
		WorkerName:  wname,
	}

	w.status[msg.ExecutionID] = response

	go func() { w.workers[wname].StartTask(msg) }()

	return response
}

func (w *workerManager) getWorker() string {
	for name, wrkr := range w.workers {
		if wrkr.Active() {
			return name
		}
	}
	return ""
}

func (w *workerManager) updateTaskStatus(msg events.RouteWorkResponseMsg) {

	defer w.lock.Unlock()
	w.lock.Lock()

	w.status[msg.ExecutionID] = msg
}

func (w *workerManager) getTaskStatus(workername string, ExecutionID string, orderID unique.TaskOrderID) events.RouteWorkResponseMsg {
	defer w.lock.Unlock()
	w.lock.Lock()

	stat := w.status[ExecutionID]
	return stat
}

func (w *workerManager) requestTaskStatus() {
	defer w.lock.Unlock()
	w.lock.Lock()

	for execid, msg := range w.status {
		go func(orderID unique.TaskOrderID, executionID string, workername string) {
			w.workers[workername].RequestTaskStatusFromWorker(orderID, executionID)
		}(msg.OrderID, execid, msg.WorkerName)

	}
}

func (w *workerManager) cleanupTask(worker string, orderID unique.TaskOrderID, executionID string, terminate bool) {

	defer w.lock.Unlock()
	w.lock.Lock()

	delete(w.status, executionID)

	if worker == "" {
		return
	}

	if terminate {
		w.workers[worker].TerminateTask(orderID, executionID)
	} else {
		w.workers[worker].CompleteTask(orderID, executionID)
	}
}

func (w *workerManager) updateWorkers(interval int) {
	for {
		time.Sleep(time.Duration(interval) * time.Second)
		for _, worker := range w.workers {
			worker.Available()
		}
	}
}

func (w *workerManager) Process(receiver events.EventReceiver, routename events.RouteName, msg events.DispatchedMessage) {

	switch routename {
	case events.RouteWorkLaunch:
		{
			data, isOk := msg.Message().(events.RouteTaskExecutionMsg)
			if isOk == false {
				err := events.ErrUnrecognizedMsgFormat
				w.log.Debug("worker,", events.RouteWorkLaunch, ",", err)
				events.ResponseToReceiver(receiver, err)
				break
			}
			w.launchChannel <- taskExecuteMsg{receiver: receiver, data: data}
		}
	case events.RouteWorkCheck:
		{
			data, isOk := msg.Message().(events.WorkRouteCheckStatusMsg)
			if isOk == false {
				err := events.ErrUnrecognizedMsgFormat
				w.log.Debug("worker,", events.RouteWorkCheck, ",", err)
				events.ResponseToReceiver(receiver, err)
				break
			}

			w.askChannel <- taskGetStatusMsg{receiver: receiver, orderID: data.OrderID, ExecutionID: data.ExecutionID, workername: data.WorkerName}
		}
	case events.RouteTaskClean:
		{
			data, isOk := msg.Message().(events.RouteTaskCleanMsg)
			if isOk == false {
			}
			w.cleanChannel <- taskCleanMsg{terminate: data.Terminate, orderID: data.OrderID, executionID: data.ExecutionID, workername: data.WorkerName}
		}
	case events.RouteTimeOut:
		{
			w.workStatus <- struct{}{}
		}
	default:
		{
			events.ResponseToReceiver(receiver, "")
		}
	}
}
