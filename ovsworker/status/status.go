package status

import "overseer/common/types"

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

func StatusExecuting(taskID, executionID string) JobExecutionStatus {
	return JobExecutionStatus{
		TaskID:      taskID,
		ExecutionID: executionID,
		State:       types.WorkerTaskStatusExecuting,
		ReturnCode:  0,
		PID:         0,
	}

}
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
