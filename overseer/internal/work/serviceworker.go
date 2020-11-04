package work

import (
	"context"
	"fmt"
	"goscheduler/common/logger"
	"goscheduler/overseer/internal/events"
	"goscheduler/overseer/internal/unique"
	"goscheduler/proto/wservices"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

type workerMediator struct {
	name       string
	connection *grpc.ClientConn
	client     wservices.TaskExecutionServiceClient
	log        logger.AppLogger
	converter  ActionConverter
}

//NewWorkerMediator - Creates a new WorkerMediator.
func NewWorkerMediator(name, host string, port int, conv ActionConverter) WorkerMediator {

	var err error
	worker := &workerMediator{}
	log := logger.Get()

	opt := make([]grpc.DialOption, 0)
	opt = append(opt, grpc.WithInsecure())

	targetAddr := fmt.Sprintf("%s:%d", host, port)
	worker.connection, err = grpc.Dial(targetAddr, opt...)
	if err != nil {
		log.Error("Unable to create worker:", err)
		return nil
	}
	worker.log = log
	worker.name = name
	worker.client = wservices.NewTaskExecutionServiceClient(worker.connection)
	worker.converter = conv

	return worker

}

//WorkerMediator - WorkerMediator is responsible for communication with
type WorkerMediator interface {
	Available() bool
	Name() string
	StartTaskExecution(msg events.RouteTaskExecutionMsg) (events.RouteWorkResponseMsg, error)
	CheckTaskStatus(taskID unique.TaskOrderID) (events.RouteWorkResponseMsg, error)
}

//Name - Returns name of a worker
func (worker *workerMediator) Name() string {
	return worker.name
}

//Avaliable - Gets status of a worker
func (worker *workerMediator) Available() bool {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()
	result, err := worker.client.WorkerStatus(ctx, &empty.Empty{})
	if err != nil {
		worker.log.Error(err)
		return false
	}
	worker.log.Info("Worker avalible result:", result)

	return true
}

//StartTaskExecution - Sends a work to execution.
func (worker *workerMediator) StartTaskExecution(msg events.RouteTaskExecutionMsg) (events.RouteWorkResponseMsg, error) {

	s := events.RouteWorkResponseMsg{Output: make([]string, 0)}

	smsg := &wservices.StartTaskMsg{}
	smsg.TaskID = &wservices.TaskIdMsg{}
	smsg.TaskID.TaskID = string(msg.OrderID)
	smsg.Type = msg.Type

	if smsg.Command = worker.converter.Convert(msg.Command, msg.Variables); smsg.Command == nil {

	}

	smsg.Variables = make(map[string]string, len(msg.Variables))
	for _, x := range msg.Variables {
		name := x.Expand()
		smsg.Variables[name] = x.Value
	}
	resp, err := worker.client.StartTask(context.Background(), smsg)

	if err != nil {
		worker.log.Error(err)
		return s, err
	}

	s.ReturnCode = resp.ReturnCode
	s.Started = true
	s.Ended = resp.Ended
	s.WorkerName = worker.name
	s.Output = append(s.Output, resp.Output...)

	return s, nil

}

//CheckTaskStatus - Sends request for a status update.
func (worker *workerMediator) CheckTaskStatus(taskID unique.TaskOrderID) (events.RouteWorkResponseMsg, error) {

	s := events.RouteWorkResponseMsg{Output: make([]string, 0)}

	resp, err := worker.client.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: string(taskID)})
	worker.log.Debug("statusTaskRemote: sending response with task status,", resp)
	if err != nil {
		return s, err
	}

	s.ReturnCode = resp.ReturnCode
	s.Ended = resp.Ended
	s.Started = true
	s.WorkerName = worker.name
	s.Output = append(s.Output, resp.Output...)

	return s, nil

}
