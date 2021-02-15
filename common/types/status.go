package types

//WorkerTaskStatus - status of a task on a remote worker
type WorkerTaskStatus int

const (
	WorkerTaskStatusRecieved  WorkerTaskStatus = 0
	WorkerTaskStatusExecuting WorkerTaskStatus = 1
	WorkerTaskStatusEnded     WorkerTaskStatus = 2
	WorkerTaskStatusFailed    WorkerTaskStatus = 3
	WorkerTaskStatusWaiting   WorkerTaskStatus = 4
	WorkerTaskStatusIdle      WorkerTaskStatus = 5
	WorkerTaskStatusStarting  WorkerTaskStatus = 6
)

//RemoteTaskStatusInfo maps remote task status to readable form
var RemoteTaskStatusInfo = map[WorkerTaskStatus]string{
	WorkerTaskStatusRecieved:  "received",
	WorkerTaskStatusExecuting: "executing",
	WorkerTaskStatusEnded:     "ended",
	WorkerTaskStatusFailed:    "failed",
	WorkerTaskStatusWaiting:   "waiting",
	WorkerTaskStatusIdle:      "idle",
	WorkerTaskStatusStarting:  "starting",
}
