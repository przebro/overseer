package work

import (
	"errors"
	"overseer/common/logger"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/unique"
)

type taskExecuteMsg struct {
	receiver events.EventReceiver
	data     events.RouteTaskExecutionMsg
}

type taskGetStatusMsg struct {
	receiver   events.EventReceiver
	orderID    unique.TaskOrderID
	workername string
}

type workerManager struct {
	askChannel    chan taskGetStatusMsg
	launchChannel chan taskExecuteMsg
	workerQueue   chan string
	log           logger.AppLogger
	workers       map[string]WorkerMediator
}

//WorkerManager - Manages actions between
type WorkerManager interface {
	Run() error
}

//NewWorkerManager - Creates a new WorkerManager
func NewWorkerManager(d events.Dispatcher, wconfig []config.WorkerConfiguration) WorkerManager {

	w := &workerManager{}
	w.log = logger.Get()
	w.askChannel = make(chan taskGetStatusMsg)
	w.launchChannel = make(chan taskExecuteMsg)
	w.workers = make(map[string]WorkerMediator)
	w.workerQueue = make(chan string)

	if d != nil {
		d.Subscribe(events.RouteWorkLaunch, w)
		d.Subscribe(events.RouteWorkCheck, w)
	}

	conv := NewConverterChain()

	for _, n := range wconfig {
		w.log.Info("Crateing service worker:", n.WorkerName, ",", n.WorkerHost, ":", n.WorkerPort)
		sworker := NewWorkerMediator(n.WorkerName, n.WorkerHost, n.WorkerPort, conv)

		if sworker != nil {
			w.workers[n.WorkerName] = sworker
		}

	}
	return w
}
func (w *workerManager) Run() error {

	go func() {
		for {

			select {
			case msg := <-w.launchChannel:
				{
					result, err := w.LaunchTask(msg.data)
					if err != nil {
						events.ResponseToReceiver(msg.receiver, err)
						break
					}
					events.ResponseToReceiver(msg.receiver, result)
				}
			case msg := <-w.askChannel:
				{
					result, err := w.CheckTaskStatus(msg.workername, msg.orderID)
					if err != nil {
						events.ResponseToReceiver(msg.receiver, err)
						break
					}
					events.ResponseToReceiver(msg.receiver, result)
				}
			}
		}

	}()
	return nil
}

func (w *workerManager) LaunchTask(msg events.RouteTaskExecutionMsg) (events.RouteWorkResponseMsg, error) {

	var result events.RouteWorkResponseMsg
	var err error = errors.New("unable to find worker")
	for _, wrkr := range w.workers {
		if wrkr.Available() {
			w.log.Info("Available worker Found,", wrkr.Name(), ", launching task:", msg.OrderID)
			result, err = wrkr.StartTaskExecution(msg)
			break
		}
	}
	return result, err
}
func (w *workerManager) CheckTaskStatus(workername string, orderID unique.TaskOrderID) (events.RouteWorkResponseMsg, error) {

	w.log.Debug("Checking worker status.")
	worker := w.workers[workername]
	result, err := worker.CheckTaskStatus(orderID)
	return result, err
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

			w.askChannel <- taskGetStatusMsg{receiver: receiver, orderID: data.OrderID, workername: data.WorkerName}
		}
	default:
		{
			events.ResponseToReceiver(receiver, "")
		}
	}

}
