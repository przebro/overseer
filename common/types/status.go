package types

import (
	"github.com/przebro/overseer/common/types/unique"
)

// WorkerTaskStatus - status of a task on a remote worker
type WorkerTaskStatus int

const (
	//WorkerTaskStatusRecieved - worker receives a job, and soon it will be executed
	WorkerTaskStatusRecieved WorkerTaskStatus = 0
	//WorkerTaskStatusExecuting - task is executed
	WorkerTaskStatusExecuting WorkerTaskStatus = 1
	//WorkerTaskStatusEnded - task ended
	WorkerTaskStatusEnded WorkerTaskStatus = 2
	//WorkerTaskStatusFailed - task failed
	WorkerTaskStatusFailed WorkerTaskStatus = 3
	//WorkerTaskStatusWaiting - worker waiting - not used
	WorkerTaskStatusWaiting WorkerTaskStatus = 4
	//WorkerTaskStatusIdle - worker is idle - not used
	WorkerTaskStatusIdle WorkerTaskStatus = 5
	//WorkerTaskStatusStarting - task staring, it will be sent to the worker
	WorkerTaskStatusStarting WorkerTaskStatus = 6
	//WorkerTaskStatusWorkerBusy - worker is busy and can't accept a new task
	WorkerTaskStatusWorkerBusy WorkerTaskStatus = 7
)

// RemoteTaskStatusInfo maps remote task status to readable form
var RemoteTaskStatusInfo = map[WorkerTaskStatus]string{
	WorkerTaskStatusRecieved:   "received",
	WorkerTaskStatusExecuting:  "executing",
	WorkerTaskStatusEnded:      "ended",
	WorkerTaskStatusFailed:     "failed",
	WorkerTaskStatusWaiting:    "waiting",
	WorkerTaskStatusIdle:       "idle",
	WorkerTaskStatusStarting:   "starting",
	WorkerTaskStatusWorkerBusy: "busy",
}

// StatusCode - Subjective status of a task execution
type StatusCode int32

const (
	//StatusCodeNormal - task ended without errors
	StatusCodeNormal StatusCode = 0
	//StatusCodeWarning - task ended with rc <= 4 or with response code 2xx
	StatusCodeWarning StatusCode = 4
	//StatusCodeError - task ended with rc > 4 or response code >= 400 or response code 2xx and error message
	StatusCodeError StatusCode = 8
	//StatusCodeTimeout - execution of a task timed out
	StatusCodeTimeout StatusCode = 9
	//StatusCodeAborted - task was aborted
	StatusCodeAborted StatusCode = 10
	//StatusCodeSevereError - something really bad happened
	StatusCodeSevereError StatusCode = 999
)

type TaskPriority int32

const (
	TaskPriorityHighest TaskPriority = 0
	TaskPriorityHigh    TaskPriority = 1
	TaskPriorityNormal  TaskPriority = 2
	TaskPriorityLow     TaskPriority = 3
	TaskPriorityLowest  TaskPriority = 4
)

type TaskExecutionStatus struct {
	Status      WorkerTaskStatus
	OrderID     unique.TaskOrderID
	ExecutionID string
	WorkerName  string
	ReturnCode  int32
	StatusCode  int32
}

// WorkDescription - describes executing task
type WorkDescription interface {
	OrderID() unique.TaskOrderID
	ExecutionID() string
	WorkerName() string
}

// TaskDescription - describes task to be executed
type TaskDescription interface {
	WorkDescription
	TypeName() TaskType
	Variables() EnvironmentVariableList
	Action() []byte
	Payload() interface{}
	SetWorkerName(string)
}
