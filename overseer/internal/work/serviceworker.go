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
func NewWorkerMediator(conf config.WorkerConfiguration, timeout int, status chan events.RouteWorkResponseMsg) WorkerMediator {

	var err error
	worker := &workerMediator{
		config:     conf,
		timeout:    timeout,
		log:        logger.Get(),
		converter:  NewConverterChain(),
		taskStatus: status,
		wdata:      workerStatus{},
		lock:       sync.Mutex{},
	}

	if err = worker.connect(conf.WorkerHost, conf.WorkerPort, worker.timeout); err != nil {
		worker.wdata.connected = false
	} else {
		worker.wdata.connected = true
	}

	return worker
}

//WorkerMediator - WorkerMediator is responsible for communication with remote workers
type WorkerMediator interface {
	Available()
	Active() bool
	Name() string
	StartTask(msg events.RouteTaskExecutionMsg)
	RequestTaskStatusFromWorker(taskID unique.TaskOrderID)
	TerminateTask(taskID unique.TaskOrderID)
	CompleteTask(taskID unique.TaskOrderID)
}

//Name - Returns name of a worker
func (worker *workerMediator) Name() string {
	return worker.config.WorkerName
}

func (worker *workerMediator) connect(host string, port int, timeout int) error {

	defer worker.lock.Unlock()
	worker.lock.Lock()
	opt := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(time.Second * time.Duration(timeout)),
	}
	targetAddr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(targetAddr, opt...)
	if err != nil {
		worker.wdata = workerStatus{connected: false}
		return err
	}

	worker.client = wservices.NewTaskExecutionServiceClient(conn)

	return nil

}

//Avaliable - Gets a status of a worker
func (worker *workerMediator) Available() {

	go func() {

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		defer worker.lock.Unlock()
		result, err := worker.client.WorkerStatus(ctx, &empty.Empty{})
		if err != nil {
			worker.lock.Lock()
			worker.wdata = workerStatus{connected: false}
		} else {
			worker.lock.Lock()
			if !worker.wdata.connected {
				worker.wdata = workerStatus{
					connected: true,
					cpu:       int(result.Cpuload),
					memused:   int(result.Memused),
					memtotal:  int(result.Memtotal),
					tasks:     int(result.Tasks),
				}
				// if connection was broken, reconnect with worker
				worker.connect(worker.config.WorkerHost, worker.config.WorkerPort, worker.timeout)
			}
		}
	}()

}

//Active - returns whether a connection with the worker is active
func (worker *workerMediator) Active() bool {
	defer worker.lock.Unlock()
	worker.lock.Lock()
	return worker.wdata.connected
}

//StartTaskExecution - Sends a task to execution.
func (worker *workerMediator) StartTask(msg events.RouteTaskExecutionMsg) {

	status := events.RouteWorkResponseMsg{
		Output:     make([]string, 0),
		OrderID:    msg.OrderID,
		WorkerName: worker.config.WorkerName,
	}

	smsg := &wservices.StartTaskMsg{}
	smsg.TaskID = &wservices.TaskIdMsg{}
	smsg.TaskID.TaskID = string(msg.OrderID)
	smsg.Type = msg.Type

	if smsg.Command = worker.converter.Convert(msg.Command, msg.Variables); smsg.Command == nil {

		status.Status = types.WorkerTaskStatusFailed
		worker.taskStatus <- status
		return

	}

	go func() {
		s := events.RouteWorkResponseMsg{
			Output:     make([]string, 0),
			OrderID:    msg.OrderID,
			WorkerName: worker.config.WorkerName,
		}
		resp, err := worker.client.StartTask(context.Background(), smsg)
		if err != nil {
			worker.log.Error(err)
			s.Status = types.WorkerTaskStatusFailed
		} else {
			s.ReturnCode = resp.ReturnCode
			s.WorkerName = worker.config.WorkerName
			s.Output = append(s.Output, resp.Output...)
			s.Status = reverseStatusMap[resp.Status]

		}

		worker.taskStatus <- s

	}()
	status.Status = types.WorkerTaskStatusStarting
	worker.taskStatus <- status

}

//RequestTaskStatusFromWorker - sends a request for a new status of a work
func (worker *workerMediator) RequestTaskStatusFromWorker(taskID unique.TaskOrderID) {

	result := events.RouteWorkResponseMsg{Output: make([]string, 0)}

	resp, err := worker.client.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: string(taskID)})
	if err != nil {
		if s, ok := status.FromError(err); ok {
			//something really bad happen with worker and task is lost
			if s.Code() == codes.NotFound {

				result.Status = types.WorkerTaskStatusFailed
				result.WorkerName = worker.config.WorkerName
				result.OrderID = taskID
				result.Output = append(result.Output, s.Message())
				worker.taskStatus <- result
			}
		}
	} else {

		result.Status = reverseStatusMap[resp.Status]
		result.ReturnCode = resp.ReturnCode
		result.WorkerName = worker.config.WorkerName
		result.OrderID = taskID
		result.Output = append(result.Output, resp.Output...)

		worker.taskStatus <- result
	}
}

//TerminateTask - terminates a task on a remote worker
func (worker *workerMediator) TerminateTask(taskID unique.TaskOrderID) {

	go func() {
		worker.client.TerminateTask(context.Background(), &wservices.TaskIdMsg{TaskID: string(taskID)})
	}()
}

//CompleteTask - sends information that the task is complete and all resources can be released
func (worker *workerMediator) CompleteTask(taskID unique.TaskOrderID) {

	go func() {
		worker.client.CompleteTask(context.Background(), &wservices.TaskIdMsg{TaskID: string(taskID)})
	}()
}
