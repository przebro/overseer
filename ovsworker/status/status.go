package status

import (
	"github.com/przebro/overseer/common/types"

	"go.uber.org/zap/zapcore"
)

//JobExecutionStatus - Contains inforamtion about a task status.
type JobExecutionStatus struct {
	TaskID      string
	ExecutionID string
	State       types.WorkerTaskStatus
	ReturnCode  int
	StatusCode  int32
	PID         int
	Reason      string
}

//MarshalLogObject - MarshalLogObject
func (s *JobExecutionStatus) MarshalLogObject(e zapcore.ObjectEncoder) error {

	e.AddString("taskID", s.TaskID)
	e.AddString("executionID", s.ExecutionID)
	e.AddInt("rc", s.ReturnCode)
	e.AddInt32("stausCode", s.StatusCode)
	e.AddString("state", types.RemoteTaskStatusInfo[s.State])
	e.AddInt("stateCode", int(s.State))
	e.AddInt("pid", s.PID)
	e.AddString("reason", s.Reason)

	return nil
}

//StatusExecuting - helper method, creates status message - executing
func StatusExecuting(taskID, executionID string) JobExecutionStatus {
	return JobExecutionStatus{
		TaskID:      taskID,
		ExecutionID: executionID,
		State:       types.WorkerTaskStatusExecuting,
		ReturnCode:  0,
		PID:         0,
	}

}

//StatusEnded - helper method, creates status message - ended
func StatusEnded(taskID, executionID string, returnCode, pid int, statusCode int32) JobExecutionStatus {
	return JobExecutionStatus{
		TaskID:      taskID,
		ExecutionID: executionID,
		State:       types.WorkerTaskStatusEnded,
		ReturnCode:  returnCode,
		StatusCode:  statusCode,
		PID:         pid,
	}
}

//StatusFailed - helper method, creates status message - failed
func StatusFailed(taskID, executionID, reason string) JobExecutionStatus {
	return JobExecutionStatus{
		TaskID:      taskID,
		ExecutionID: executionID,
		State:       types.WorkerTaskStatusFailed,
		ReturnCode:  0,
		PID:         0,
		Reason:      reason,
		StatusCode:  int32(types.StatusCodeSevereError),
	}
}
