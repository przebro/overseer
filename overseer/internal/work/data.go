package work

import (
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/proto/wservices"
)

// WorkDescription - describes executing task
type WorkDescription interface {
	OrderID() unique.TaskOrderID
	ExecutionID() string
	WorkerName() string
}

// TaskDescription - describes task to be executed
type TaskDescription interface {
	WorkDescription
	TypeName() types.TaskType
	Variables() types.EnvironmentVariableList
	Action() []byte
	Payload() interface{}
	SetWorkerName(string)
}

type workerState struct {
	connected   bool
	cpu         int
	memused     int
	memtotal    int
	tasks       int
	tasksLimit  int
	lastRequest time.Time
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
