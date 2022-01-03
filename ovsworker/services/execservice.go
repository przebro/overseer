package services

import (
	"context"
	"errors"

	"fmt"
	"os"
	"path/filepath"

	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/common/types"

	"github.com/przebro/overseer/ovsworker/jobs"
	"github.com/przebro/overseer/ovsworker/msgheader"
	"github.com/przebro/overseer/ovsworker/task"
	"github.com/przebro/overseer/proto/wservices"

	"github.com/golang/protobuf/ptypes/empty"
	"go.uber.org/zap"
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

var (
	//ErrEmptyMsgType  - empty message error
	ErrEmptyMsgType error = errors.New("message type cannot be empty")
	//ErrEmptyTaskID  - empty taskID error
	ErrEmptyTaskID error = errors.New("TaskID cannot be empty")
	//ErrEmptyExecutionID  - empty executionID error
	ErrEmptyExecutionID error = errors.New("ExecutionID cannot be empty")
)

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

	wservice.te = task.NewTaskRunnerManager()

	return wservice, nil
}

func (wsrvc *workerExecutionService) StartTask(ctx context.Context, msg *wservices.StartTaskMsg) (*wservices.TaskExecutionResponseMsg, error) {

	var response *wservices.TaskExecutionResponseMsg
	var err error

	if msg.Type == "" {
		wsrvc.log.Desugar().Error("StartTask", zap.Uint32("code", uint32(codes.Aborted)),
			zap.String("descr", codes.Aborted.String()),
			zap.String("error", ErrEmptyMsgType.Error()))

		return nil, status.Error(codes.Aborted, ErrEmptyMsgType.Error())
	}
	if msg.TaskID.TaskID == "" {
		wsrvc.log.Desugar().Error("StartTask", zap.Uint32("code", uint32(codes.Aborted)),
			zap.String("descr", codes.Aborted.String()),
			zap.String("error", ErrEmptyTaskID.Error()))

		return nil, status.Error(codes.Aborted, ErrEmptyTaskID.Error())
	}

	if msg.TaskID.ExecutionID == "" {
		wsrvc.log.Desugar().Error("StartTask", zap.Uint32("code", uint32(codes.Aborted)),
			zap.String("descr", codes.Aborted.String()),
			zap.String("error", ErrEmptyTaskID.Error()))

		return nil, status.Error(codes.Aborted, ErrEmptyExecutionID.Error())
	}

	header := msgheader.TaskHeader{
		Type:        types.TaskType(msg.Type),
		TaskID:      msg.TaskID.TaskID,
		ExecutionID: msg.TaskID.ExecutionID,
		Variables:   msg.Variables,
	}

	var jobExec jobs.JobExecutor

	if jobExec, err = jobs.NewJobExecutor(header, wsrvc.sysoutDir, msg.Command.Value, wsrvc.log); err != nil {

		wsrvc.log.Desugar().Error("StartTask", zap.Uint32("code", uint32(codes.Aborted)),
			zap.String("descr", codes.Aborted.String()),
			zap.String("error", err.Error()))

		return nil, status.Error(codes.Aborted, err.Error())

	}

	wsrvc.log.Desugar().Info("StartTask",
		zap.String("descr", "Executor Created"),
		zap.String("type", string(header.Type)),
		zap.String("taskID", header.TaskID),
		zap.String("executionID", header.ExecutionID))

	status, num := wsrvc.te.RunTask(jobExec)

	wsrvc.log.Desugar().Info("StartTask", zap.Object("payload", &status))

	response = &wservices.TaskExecutionResponseMsg{
		Status:     wservices.TaskExecutionResponseMsg_RECEIVED,
		TasksLimit: int32(wsrvc.taskLimit),
		Tasks:      int32(num),
	}

	return response, nil
}

func (wsrvc *workerExecutionService) TaskStatus(ctx context.Context, msg *wservices.TaskIdMsg) (*wservices.TaskExecutionResponseMsg, error) {

	var response *wservices.TaskExecutionResponseMsg
	result, num, ok := wsrvc.te.GetTaskStatus(msg.ExecutionID)

	if !ok {
		msg := fmt.Sprintf("task:%s does not exists", msg.TaskID)

		wsrvc.log.Desugar().Error("TaskStatus", zap.Uint32("code", uint32(codes.NotFound)),
			zap.String("descr", codes.NotFound.String()),
			zap.String("error", msg))

		return nil, status.Error(codes.NotFound, msg)
	}

	wsrvc.log.Desugar().Info("TaskStatus", zap.Object("payload", &result))

	response = &wservices.TaskExecutionResponseMsg{

		Status:     statusMap[result.State],
		ReturnCode: int32(result.ReturnCode),
		StatusCode: int32(result.StatusCode),
		Pid:        int32(result.PID),
		TasksLimit: int32(wsrvc.taskLimit),
		Tasks:      int32(num),
	}

	return response, nil

}

func (wsrvc *workerExecutionService) TerminateTask(context.Context, *wservices.TaskIdMsg) (*wservices.WorkerActionMsg, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (wsrvc *workerExecutionService) CompleteTask(ctx context.Context, msg *wservices.TaskIdMsg) (*wservices.WorkerActionMsg, error) {

	tasks := wsrvc.te.CleanupTask(msg.ExecutionID)

	wsrvc.log.Desugar().Info("CompleteTask", zap.String("descr", "task removed"),
		zap.String("taskID", msg.TaskID),
		zap.String("executionID", msg.ExecutionID))

	resp := &wservices.WorkerActionMsg{Message: "task removed", Success: true, Tasks: int32(tasks), TasksLimit: int32(wsrvc.taskLimit)}
	return resp, nil
}

func (wsrvc *workerExecutionService) TaskOutput(*wservices.TaskIdMsg, wservices.TaskExecutionService_TaskOutputServer) error {
	return status.Error(codes.Unimplemented, "not implemented")
}

func (wsrvc *workerExecutionService) WorkerStatus(ctx context.Context, msg *empty.Empty) (*wservices.WorkerStatusResponseMsg, error) {

	num := wsrvc.te.TaskCount()

	response := &wservices.WorkerStatusResponseMsg{Tasks: int32(num), TasksLimit: int32(wsrvc.taskLimit), Memused: 0, Memtotal: 0, Cpuload: 0}

	return response, nil

}
