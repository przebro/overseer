package work

import (
	"context"
	"fmt"
	"overseer/common/logger"
	"overseer/common/types"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/unique"
	"overseer/proto/wservices"
	"sync"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type workerMediator struct {
	config     config.WorkerConfiguration
	connection *grpc.ClientConn
	client     wservices.TaskExecutionServiceClient
	log        logger.AppLogger
	converter  ActionConverter
	taskStatus chan events.RouteWorkResponseMsg
	timeout    int
	wdata      workerStatus
	lock       sync.Mutex
}

//NewWorkerMediator - Creates a new WorkerMediator.
func NewWorkerMediator(conf config.WorkerConfiguration, timeout int, status chan events.RouteWorkResponseMsg, log logger.AppLogger) WorkerMediator {

	worker := &workerMediator{
		config:     conf,
		timeout:    timeout,
		log:        log,
		converter:  NewConverterChain(),
		taskStatus: status,
		wdata:      workerStatus{},
		lock:       sync.Mutex{},
	}

	if client := worker.connect(conf.WorkerHost, conf.WorkerPort, worker.timeout); client == nil {
		worker.wdata.connected = false
	} else {
		worker.client = client
		worker.wdata.connected = true
	}

	return worker
}

//WorkerMediator - WorkerMediator is responsible for communication with remote workers
type WorkerMediator interface {
	Available()
	Active() workerStatus
	Name() string
	StartTask(msg events.RouteTaskExecutionMsg)
	RequestTaskStatusFromWorker(taskID unique.TaskOrderID, executionID string)
	TerminateTask(taskID unique.TaskOrderID, executionID string)
	CompleteTask(taskID unique.TaskOrderID, executionID string)
}

//Name - Returns name of a worker
func (worker *workerMediator) Name() string {
	return worker.config.WorkerName
}

func (worker *workerMediator) connect(host string, port int, timeout int) wservices.TaskExecutionServiceClient {

	opt := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second * time.Duration(timeout)),
	}
	targetAddr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(targetAddr, opt...)
	if err != nil {
		worker.lock.Lock()
		worker.wdata = workerStatus{connected: false}
		worker.lock.Unlock()
		return nil
	}

	return wservices.NewTaskExecutionServiceClient(conn)
}

//Available - Gets a status of a worker
func (worker *workerMediator) Available() {

	go func(w *workerMediator) {
		var client wservices.TaskExecutionServiceClient

		w.lock.Lock()
		client = worker.client
		w.lock.Unlock()

		if client == nil {
			if client := w.connect(w.config.WorkerHost, w.config.WorkerPort, w.timeout); client == nil {
				worker.log.Error("connect with worker failed")
			} else {
				w.lock.Lock()
				worker.client = client
				w.lock.Unlock()
			}

		} else {

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer w.lock.Unlock()
			defer cancel()
			result, err := w.client.WorkerStatus(ctx, &empty.Empty{})

			if err != nil {
				worker.log.Error("worker connection lost")
				w.lock.Lock()
				w.wdata = workerStatus{connected: false}
			} else {
				w.lock.Lock()

				w.wdata = workerStatus{
					connected:  true,
					cpu:        int(result.Cpuload),
					memused:    int(result.Memused),
					memtotal:   int(result.Memtotal),
					tasks:      int(result.Tasks),
					tasksLimit: int(result.TasksLimit),
				}
			}
		}

	}(worker)

}

//Active - returns whether a connection with the worker is active
func (worker *workerMediator) Active() workerStatus {
	defer worker.lock.Unlock()
	worker.lock.Lock()
	return worker.wdata
}

//StartTaskExecution - Sends a task to execution.
func (worker *workerMediator) StartTask(msg events.RouteTaskExecutionMsg) {

	status := events.RouteWorkResponseMsg{
		ExecutionID: msg.ExecutionID,
		OrderID:     msg.OrderID,
		WorkerName:  worker.config.WorkerName,
	}

	smsg := &wservices.StartTaskMsg{}
	smsg.TaskID = &wservices.TaskIdMsg{TaskID: string(msg.OrderID), ExecutionID: msg.ExecutionID}
	smsg.Type = msg.Type

	smsg.Variables = map[string]string{}

	for _, n := range msg.Variables {
		smsg.Variables[n.Expand()] = n.Value

	}

	if smsg.Command = worker.converter.Convert(msg.Command, msg.Variables, worker.log); smsg.Command == nil {

		status.Status = types.WorkerTaskStatusFailed
		worker.taskStatus <- status
		return
	}

	go func(worker *workerMediator) {
		s := events.RouteWorkResponseMsg{
			OrderID:     msg.OrderID,
			ExecutionID: msg.ExecutionID,
			WorkerName:  worker.config.WorkerName,
		}
		resp, err := worker.client.StartTask(context.Background(), smsg)
		if err != nil {
			worker.log.Error(err)
			s.Status = types.WorkerTaskStatusFailed
		} else {
			s.ReturnCode = resp.ReturnCode
			s.WorkerName = worker.config.WorkerName
			s.Status = reverseStatusMap[resp.Status]

			worker.setTaskInfo(int(resp.TasksLimit), int(resp.Tasks))

		}

		worker.taskStatus <- s

	}(worker)

	status.Status = types.WorkerTaskStatusStarting
	worker.taskStatus <- status

}

//RequestTaskStatusFromWorker - sends a request for a new status of a work
func (worker *workerMediator) RequestTaskStatusFromWorker(taskID unique.TaskOrderID, executionID string) {

	result := events.RouteWorkResponseMsg{OrderID: taskID, ExecutionID: executionID}

	resp, err := worker.client.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: string(taskID), ExecutionID: executionID})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			//something really bad happen with worker and task is lost
			if s.Code() == codes.NotFound {

				result.Status = types.WorkerTaskStatusFailed
				result.WorkerName = worker.config.WorkerName
				worker.taskStatus <- result
			}
		}
	} else {

		result.Status = reverseStatusMap[resp.Status]
		result.ReturnCode = resp.ReturnCode
		result.WorkerName = worker.config.WorkerName

		worker.taskStatus <- result
		worker.setTaskInfo(int(resp.TasksLimit), int(resp.Tasks))

	}
}

//TerminateTask - terminates a task on a remote worker
func (worker *workerMediator) TerminateTask(taskID unique.TaskOrderID, executionID string) {

	go func() {
		worker.client.TerminateTask(context.Background(), &wservices.TaskIdMsg{TaskID: string(taskID), ExecutionID: executionID})
	}()
}

//CompleteTask - sends information that the task is complete and all resources can be released
func (worker *workerMediator) CompleteTask(taskID unique.TaskOrderID, executionID string) {

	go func(w *workerMediator) {
		resp, err := worker.client.CompleteTask(context.Background(), &wservices.TaskIdMsg{TaskID: string(taskID), ExecutionID: executionID})
		if err != nil {
			return
		}

		w.setTaskInfo(int(resp.TasksLimit), int(resp.Tasks))

	}(worker)
}

func (worker *workerMediator) setTaskInfo(limit, tasks int) {
	defer worker.lock.Unlock()
	worker.lock.Lock()

	worker.wdata.tasks = tasks
	worker.wdata.tasksLimit = limit

}
