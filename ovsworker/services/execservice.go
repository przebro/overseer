package services

import (
	"context"
	"errors"

	"fmt"
	"os"
	"overseer/common/logger"
	"overseer/common/types"
	"path/filepath"

	"overseer/ovsworker/fragments"
	"overseer/ovsworker/msgheader"
	"overseer/ovsworker/task"
	"overseer/proto/wservices"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var statusMap = map[types.WorkerTaskStatus]wservices.TaskExecutionResponseMsg_TaskStatus{
	types.WorkerTaskStatusRecieved:  wservices.TaskExecutionResponseMsg_RECEIVED,
	types.WorkerTaskStatusExecuting: wservices.TaskExecutionResponseMsg_EXECUTING,
	types.WorkerTaskStatusEnded:     wservices.TaskExecutionResponseMsg_ENDED,
	types.WorkerTaskStatusFailed:    wservices.TaskExecutionResponseMsg_FAILED,
	types.WorkerTaskStatusWaiting:   wservices.TaskExecutionResponseMsg_WAITING,
	types.WorkerTaskStatusIdle:      wservices.TaskExecutionResponseMsg_IDLE,
}

type workerExecutionService struct {
	log       logger.AppLogger
	te        *task.TaskExecutor
	sysoutDir string
}

func NewWorkerExecutionService(sysoutDir string) (*workerExecutionService, error) {

	var sysout string
	var err error
	var nfo os.FileInfo

	if !filepath.IsAbs(sysoutDir) {
		if sysout, err = filepath.Abs(sysoutDir); err != nil {
			return nil, err
		}
	} else {
		sysout = sysoutDir
	}

	if nfo, err = os.Stat(sysout); err != nil {
		return nil, err
	}

	if !nfo.IsDir() {
		return nil, errors.New("sysout path points to a file not to the directory")
	}

	wservice := &workerExecutionService{
		log:       logger.Get(),
		sysoutDir: sysout,
	}

	wservice.te = task.NewTaskExecutor()

	return wservice, nil

}

func (wsrvc *workerExecutionService) StartTask(ctx context.Context, msg *wservices.StartTaskMsg) (*wservices.TaskExecutionResponseMsg, error) {

	var response *wservices.TaskExecutionResponseMsg
	var err error

	if msg.Type == "" {
		return nil, status.Error(codes.Aborted, "message type cannot be empty")
	}
	if msg.TaskID.TaskID == "" {
		return nil, status.Error(codes.Aborted, "message taskID cannot be empty")
	}

	if msg.TaskID.ExecutionID == "" {
		return nil, status.Error(codes.Aborted, "message ExecutionID cannot be empty")
	}

	header := msgheader.TaskHeader{
		Type:        types.TaskType(msg.Type),
		TaskID:      msg.TaskID.TaskID,
		ExecutionID: msg.TaskID.ExecutionID,
		Variables:   msg.Variables,
	}

	var frag fragments.WorkFragment

	if frag, err = fragments.CreateWorkFragment(header, wsrvc.sysoutDir, msg.Command.Value); err != nil {
		return nil, status.Error(codes.Aborted, err.Error())

	}

	wsrvc.te.ExecuteTask(frag)

	response = &wservices.TaskExecutionResponseMsg{
		Status: wservices.TaskExecutionResponseMsg_RECEIVED,
	}

	wsrvc.log.Info("TaskStatus: response for:", msg.TaskID.TaskID, ",", msg.TaskID.ExecutionID, ";", response)
	return response, nil
}

func (wsrvc *workerExecutionService) TaskStatus(ctx context.Context, msg *wservices.TaskIdMsg) (*wservices.TaskExecutionResponseMsg, error) {

	var response *wservices.TaskExecutionResponseMsg
	result, ok := wsrvc.te.GetTaskStatus(msg.ExecutionID)

	if ok != true {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("task:%s does not exists", msg.TaskID))
	}

	response = &wservices.TaskExecutionResponseMsg{

		Status:     statusMap[result.State],
		ReturnCode: int32(result.ReturnCode),
		StatusCode: int32(result.StatusCode),
		Pid:        int32(result.PID),
	}

	wsrvc.log.Info("TaskStatus: response for:", msg.TaskID, ",", msg.ExecutionID, ";", types.RemoteTaskStatusInfo[result.State])

	return response, nil

}

func (wsrvc *workerExecutionService) TerminateTask(context.Context, *wservices.TaskIdMsg) (*wservices.WorkerActionMsg, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (wsrvc *workerExecutionService) CompleteTask(ctx context.Context, msg *wservices.TaskIdMsg) (*wservices.WorkerActionMsg, error) {

	resp := &wservices.WorkerActionMsg{Message: "task removed", Success: true}
	wsrvc.te.CleanupTask(msg.ExecutionID)
	return resp, nil
}

func (wsrvc *workerExecutionService) TaskOutput(*wservices.TaskIdMsg, wservices.TaskExecutionService_TaskOutputServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (wsrvc *workerExecutionService) WorkerStatus(ctx context.Context, msg *empty.Empty) (*wservices.WorkerStatusResponseMsg, error) {

	num := wsrvc.te.TaskCount()

	response := &wservices.WorkerStatusResponseMsg{Tasks: int32(num), Memused: 0, Memtotal: 0, Cpuload: 0}
	wsrvc.log.Info("WorkStatus: response:", response)
	return response, nil

}
