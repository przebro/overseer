package types

//WorkerTaskStatus - status of a task on a remote worker
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
	//WorkerTaskStatusWorkerBusy -worker is busy and can't accpet a new task
	WorkerTaskStatusWorkerBusy WorkerTaskStatus = 7
)

//RemoteTaskStatusInfo maps remote task status to readable form
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

//StatusCode - Subjective status of a task execution
type StatusCode int32

const (
	//StatusCodeNormal - task ended without errors
	StatusCodeNormal StatusCode = 0
	//StatusCodeWarning - task ended with rc <= 4 or with response code 2xx
	StatusCodeWarning StatusCode = 4
	//StatusCodeError - task ended with rc > 4 or response code >= 400 or reponse code 2xx and error message
	StatusCodeError StatusCode = 8
	//StatusCodeTimeout - execution of a task timed out
	StatusCodeTimeout StatusCode = 9
	//StatusCodeAborted - task was aborted
	StatusCodeAborted StatusCode = 10
	//StatusCodeSevereError - something really bad happend
	StatusCodeSevereError StatusCode = 999
)
