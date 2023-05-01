package models

// TaskState - current state of a task
type TaskState int32

// Possible states of an active task
const (
	//Task waits for a time window, tickets
	TaskStateWaiting  TaskState = 1
	TaskStateStarting TaskState = 2
	//Time and tickets prerequisites were met, task was sent to worker. Task may wait for confirm or flag or being executed
	TaskStateExecuting  TaskState = 3
	TaskStateEndedOk    TaskState = 4
	TaskStateEndedNotOk TaskState = 5
	TaskStateHold       TaskState = 6
)

func (t TaskState) String() string {
	state := ""
	switch t {
	case TaskStateWaiting:
		state = "Waiting"
	case TaskStateStarting:
		state = "Starting"
	case TaskStateExecuting:
		state = "Executing"
	case TaskStateEndedOk:
		state = "Ended OK"
	case TaskStateEndedNotOk:
		state = "Ended Not OK"
	case TaskStateHold:
		state = "Held"
	default:
		state = "Unknown"
	}
	return state
}
