package services

import (
	"context"
	"errors"

	"fmt"
	"os"
	"overseer/common/logger"
	"overseer/common/types"
	"path/filepath"

	"overseer/ovsworker/jobs"
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
	te        *task.TaskRunnerManager
	sysoutDir string
	taskLimit int
}

//NewWorkerExecutionService - creates a new instance of a workerExecutionService
func NewWorkerExecutionService(sysoutDir string, limit int, log logger.AppLogger) (*workerExecutionService, error) {

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
		log:       log,
		sysoutDir: sysout,
		taskLimit: limit,
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

	var jobExec jobs.JobExecutor

	if jobExec, err = jobs.NewJobExecutor(header, wsrvc.sysoutDir, msg.Command.Value); err != nil {
		return nil, status.Error(codes.Aborted, err.Error())

	}

	_, num := wsrvc.te.RunTask(jobExec)

	response = &wservices.TaskExecutionResponseMsg{
		Status:     wservices.TaskExecutionResponseMsg_RECEIVED,
		TasksLimit: int32(wsrvc.taskLimit),
		Tasks:      int32(num),
	}

	wsrvc.log.Info("TaskStatus: response for:", msg.TaskID.TaskID, ",", msg.TaskID.ExecutionID, ";", response)
	return response, nil
}

func (wsrvc *workerExecutionService) TaskStatus(ctx context.Context, msg *wservices.TaskIdMsg) (*wservices.TaskExecutionResponseMsg, error) {

	var response *wservices.TaskExecutionResponseMsg
	result, num, ok := wsrvc.te.GetTaskStatus(msg.ExecutionID)

	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("task:%s does not exists", msg.TaskID))
	}

	response = &wservices.TaskExecutionResponseMsg{

		Status:     statusMap[result.State],
		ReturnCode: int32(result.ReturnCode),
		StatusCode: int32(result.StatusCode),
		Pid:        int32(result.PID),
		TasksLimit: int32(wsrvc.taskLimit),
		Tasks:      int32(num),
	}

	wsrvc.log.Info("TaskStatus:", msg.TaskID, ",", msg.ExecutionID, ";", types.RemoteTaskStatusInfo[result.State], "task processed:", num, "task limit:", wsrvc.taskLimit)

	return response, nil

}

func (wsrvc *workerExecutionService) TerminateTask(context.Context, *wservices.TaskIdMsg) (*wservices.WorkerActionMsg, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (wsrvc *workerExecutionService) CompleteTask(ctx context.Context, msg *wservices.TaskIdMsg) (*wservices.WorkerActionMsg, error) {

	tasks := wsrvc.te.CleanupTask(msg.ExecutionID)
	wsrvc.log.Info("Clean Task: ID:", msg.TaskID, "EID:", msg.ExecutionID, "tasks:", tasks, "limit:", wsrvc.taskLimit)
	resp := &wservices.WorkerActionMsg{Message: "task removed", Success: true, Tasks: int32(tasks), TasksLimit: int32(wsrvc.taskLimit)}
	return resp, nil
}

func (wsrvc *workerExecutionService) TaskOutput(*wservices.TaskIdMsg, wservices.TaskExecutionService_TaskOutputServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (wsrvc *workerExecutionService) WorkerStatus(ctx context.Context, msg *empty.Empty) (*wservices.WorkerStatusResponseMsg, error) {

	num := wsrvc.te.TaskCount()

	response := &wservices.WorkerStatusResponseMsg{Tasks: int32(num), TasksLimit: int32(wsrvc.taskLimit), Memused: 0, Memtotal: 0, Cpuload: 0}
	wsrvc.log.Info("WorkStatus: response:", "tasks:", num, "limit:", response.TasksLimit)
	return response, nil

}
