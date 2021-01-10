package services

import (
	"context"
	"errors"
	"overseer/common/logger"
	common "overseer/common/types"
	"overseer/ovsworker/launcher"
	"overseer/ovsworker/msgheader"
	"overseer/proto/wservices"

	"github.com/golang/protobuf/ptypes/empty"
)

type workerExecutionService struct {
	log      logger.AppLogger
	launcher *launcher.FragmentLauncher
	creator  *launcher.FragmentCreator
}

func NewWorkerExecutionService() *workerExecutionService {

	wservice := &workerExecutionService{
		log: logger.Get(),
	}
	wservice.launcher = launcher.NewFragmentLauncher()
	wservice.creator = launcher.FragmentFactory(wservice.launcher)

	return wservice

}

func (wsrvc *workerExecutionService) StartTask(ctx context.Context, msg *wservices.StartTaskMsg) (*wservices.TaskExecutionResponseMsg, error) {

	var response *wservices.TaskExecutionResponseMsg

	if msg.Type == "" {
		return nil, errors.New("message type cannot be empty")
	}
	if msg.TaskID.TaskID == "" {
		return nil, errors.New("message taskID cannot be empty")
	}

	header := msgheader.TaskHeader{Type: common.TaskType(msg.Type), TaskID: msg.TaskID.TaskID, Variables: msg.Variables}

	err := wsrvc.creator.CreateFragment(header, msg.Command.Value)
	if err != nil {
		return nil, err
	}

	ectx := context.Background()
	ch, err := wsrvc.launcher.Execute(ectx, header.TaskID)

	select {
	case result := <-ch:
		{
			response = &wservices.TaskExecutionResponseMsg{
				Started:    result.Started,
				Ended:      result.Ended,
				Output:     result.Output,
				ReturnCode: int32(result.ReturnCode),
				StatusCode: int32(result.StatusCode)}
		}
	}

	wsrvc.log.Info("TaskStatus: response for:", msg.TaskID, ";", response)
	return response, nil
}

func (wsrvc *workerExecutionService) TaskStatus(ctx context.Context, msg *wservices.TaskIdMsg) (*wservices.TaskExecutionResponseMsg, error) {

	var response *wservices.TaskExecutionResponseMsg

	result, err := wsrvc.launcher.Status(msg.TaskID)

	if err != nil {
		return nil, err
	}

	response = &wservices.TaskExecutionResponseMsg{
		Started:    result.Started,
		Ended:      result.Ended,
		Output:     result.Output,
		ReturnCode: int32(result.ReturnCode),
		StatusCode: int32(result.StatusCode)}

	wsrvc.log.Info("TaskStatus: response for:", msg.TaskID, ";", result)

	return response, nil

}
func (wsrvc *workerExecutionService) WorkerStatus(ctx context.Context, msg *empty.Empty) (*wservices.WorkerStatusResponseMsg, error) {

	n := wsrvc.launcher.Tasks()
	response := &wservices.WorkerStatusResponseMsg{Tasks: int32(n), Memused: 0, Memtotal: 0, Cpuload: 0}
	wsrvc.log.Info("WorkStatus: response:", response)
	return response, nil

}
