package work

import (
	"overseer/common/types"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/unique"
	"overseer/proto/wservices"
)

type taskExecuteMsg struct {
	receiver events.EventReceiver
	data     events.RouteTaskExecutionMsg
}

type taskCleanMsg struct {
	receiver    events.EventReceiver
	orderID     unique.TaskOrderID
	executionID string
	workername  string
	terminate   bool
}

type taskGetStatusMsg struct {
	receiver    events.EventReceiver
	orderID     unique.TaskOrderID
	ExecutionID string
	workername  string
}

type workerStatus struct {
	connected bool
	cpu       int
	memused   int
	memtotal  int
	tasks     int
}

var reverseStatusMap = map[wservices.TaskExecutionResponseMsg_TaskStatus]types.WorkerTaskStatus{
	wservices.TaskExecutionResponseMsg_RECEIVED:  types.WorkerTaskStatusRecieved,
	wservices.TaskExecutionResponseMsg_EXECUTING: types.WorkerTaskStatusExecuting,
	wservices.TaskExecutionResponseMsg_ENDED:     types.WorkerTaskStatusEnded,
	wservices.TaskExecutionResponseMsg_FAILED:    types.WorkerTaskStatusFailed,
	wservices.TaskExecutionResponseMsg_WAITING:   types.WorkerTaskStatusWaiting,
	wservices.TaskExecutionResponseMsg_IDLE:      types.WorkerTaskStatusIdle,
	wservices.TaskExecutionResponseMsg_STARTING:  types.WorkerTaskStatusStarting,
}
