package services

import (
	"context"
	"errors"

	"fmt"
	"os"
	"path/filepath"

	"github.com/przebro/overseer/common/types"
	"github.com/rs/zerolog"

	"github.com/przebro/overseer/ovsworker/jobs"
	"github.com/przebro/overseer/ovsworker/msgheader"
	"github.com/przebro/overseer/ovsworker/task"
	"github.com/przebro/overseer/proto/wservices"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
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
	te        *task.TaskRunnerManager
	sysoutDir string
	taskLimit int
	wservices.UnimplementedTaskExecutionServiceServer
}

// NewWorkerExecutionService - creates a new instance of a workerExecutionService
func NewWorkerExecutionService(sysoutDir string, limit int) (*workerExecutionService, error) {

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
		sysoutDir: sysout,
		taskLimit: limit,
	}

	wservice.te = task.NewTaskRunnerManager()

	return wservice, nil
}

func (wsrvc *workerExecutionService) StartTask(ctx context.Context, msg *wservices.StartTaskMsg) (*wservices.TaskExecutionResponseMsg, error) {

	var response *wservices.TaskExecutionResponseMsg
	var err error

	log := zerolog.Ctx(ctx).With().Str("service", "exec").Logger()

	if msg.Type == "" {
		log.Error().Str("code", codes.Aborted.String()).Err(ErrEmptyMsgType).Msg("StartTask")

		return nil, status.Error(codes.Aborted, ErrEmptyMsgType.Error())
	}
	if msg.TaskID.TaskID == "" {
		log.Error().Str("code", codes.Aborted.String()).Err(ErrEmptyTaskID).Msg("StartTask")

		return nil, status.Error(codes.Aborted, ErrEmptyTaskID.Error())
	}

	if msg.TaskID.ExecutionID == "" {
		log.Error().Str("code", codes.Aborted.String()).Err(ErrEmptyTaskID).Msg("StartTask")

		return nil, status.Error(codes.Aborted, ErrEmptyExecutionID.Error())
	}

	header := msgheader.TaskHeader{
		Type:        types.TaskType(msg.Type),
		TaskID:      msg.TaskID.TaskID,
		ExecutionID: msg.TaskID.ExecutionID,
		Variables:   msg.Variables,
	}

	var jobExec jobs.JobExecutor

	if jobExec, err = jobs.NewJobExecutor(ctx, header, wsrvc.sysoutDir, msg.Command.Value); err != nil {

		log.Error().Str("code", codes.Aborted.String()).Err(err).Msg("StartTask")

		return nil, status.Error(codes.Aborted, err.Error())

	}

	log.Info().Str("task_id", header.TaskID).Str("execution_id", header.ExecutionID).
		Str("type", string(header.Type)).Msg("StartTask")

	status, num := wsrvc.te.RunTask(jobExec)
	log.Info().Interface("status", status).Msg("StartTask")

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

	log := zerolog.Ctx(ctx).With().Str("service", "status").Logger()

	if !ok {
		msg := fmt.Sprintf("task:%s does not exists", msg.TaskID)
		log.Error().Str("code", codes.NotFound.String()).Str("error", msg).Msg("TaskStatus")

		return nil, status.Error(codes.NotFound, msg)
	}

	log.Info().Interface("status", result).Msg("StartTask")

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

func (wsrvc *workerExecutionService) TerminateTask(ctx context.Context, msg *wservices.TaskIdMsg) (*wservices.WorkerActionMsg, error) {

	log := zerolog.Ctx(ctx).With().Str("service", "termniate").Logger()
	log.Error().Msg("not implemented")

	return nil, status.Error(codes.Unimplemented, "not implemented")
}
func (wsrvc *workerExecutionService) CompleteTask(ctx context.Context, msg *wservices.TaskIdMsg) (*wservices.WorkerActionMsg, error) {

	log := zerolog.Ctx(ctx).With().Str("service", "exec").Logger()

	tasks := wsrvc.te.CleanupTask(msg.ExecutionID)
	log.Info().Int("tasks", tasks).Str("task_id", msg.TaskID).Str("execution_id", msg.ExecutionID).Msg("task complete")

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
