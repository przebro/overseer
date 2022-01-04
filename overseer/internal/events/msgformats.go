package events

import (
	"encoding/json"
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	task "github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/przebro/overseer/overseer/internal/unique"
)

//RouteTimeOutMsgFormat - outgoing message from ticker
type RouteTimeOutMsgFormat struct {
	Year  int
	Month int
	Day   int
	Hour  int
	Min   int
	Sec   int
}

//RouteTicketCheckMsgFormat - Status of a tickets requested by a task.
type RouteTicketCheckMsgFormat struct {
	Tickets []struct {
		Name      string
		Odate     string
		Fulfilled bool
	}
}

//FlagActionData - struct for acquire flag message
type FlagActionData struct {
	Name   string
	Policy int8
}

//RouteFlagAcquireMsg - acquires flags
type RouteFlagAcquireMsg struct {
	Flags []FlagActionData
}

/*RouteFlagActionResponse - response for both RouteFlagAcquireMsg and RouteFlagReleaseMsg messages.
Contains information if all required flags were acquired; if Success is false, then the Names field will contain
the name of the flag that was not acquired
if it is a response for a RouteFlagReleaseMsg then the Names filed will contain all flags that were not unset
*/
type RouteFlagActionResponse struct {
	Success bool
	Names   []string
}

//RouteTicketInMsgFormat - Basing on this structure, tickets are added or removed from the resources manager
type RouteTicketInMsgFormat struct {
	Tickets []struct {
		Name   string
		Odate  date.Odate
		Action task.OutAction
	}
}

//TaskInfoResultMsg - Result for a request for task information
type TaskInfoResultMsg struct {
	TaskID      unique.TaskOrderID
	Odate       date.Odate
	Group       string
	Name        string
	State       int32
	RunNumber   int32
	Held        bool
	Confirmed   bool
	WaitingInfo string
}

//TaskCycleMsg - task cyclic data
type TaskCycleMsg struct {
	IsCyclic    bool
	NextRun     types.HourMinTime
	RunFrom     string
	MaxRun      int
	RunInterval int
}

//TaskDetailResultMsg - Result for a request for detailed information
type TaskDetailResultMsg struct {
	TaskInfoResultMsg
	TaskCycleMsg
	From        string
	To          string
	StartTime   string
	EndTime     string
	Worker      string
	Description string
	Tickets     []struct {
		Name      string
		Odate     date.Odate
		Fulfilled bool
	}
}

//RouteTaskActionMsgFormat - Request for a task order
type RouteTaskActionMsgFormat struct {
	Group    string
	Name     string
	TaskID   unique.TaskOrderID
	Force    bool
	Odate    date.Odate
	Username string
}

//RouteTaskActionResponseFormat - Response message for a task order or force.
type RouteTaskActionResponseFormat struct {
	Data []TaskInfoResultMsg
}

//WorkRouteCheckStatusMsg - Response with information about the status of a work
type WorkRouteCheckStatusMsg struct {
	OrderID     unique.TaskOrderID
	ExecutionID string
	WorkerName  string
}

//RouteTaskStatusResponseMsg - Response for a task status
type RouteTaskStatusResponseMsg struct {
	TaskID      string
	ExecutionID string
	Data        []string
	Ended       bool
	ReturnCode  int32
}

//RouteTaskExecutionMsg - Contains informations needed to begin a work on a remoteworker.
type RouteTaskExecutionMsg struct {
	OrderID     unique.TaskOrderID
	ExecutionID string
	Type        types.TaskType
	Variables   []task.VariableData
	Command     json.RawMessage
}

//RouteWorkResponseMsg - Contains information about the status of executing work.
type RouteWorkResponseMsg struct {
	Status      types.WorkerTaskStatus
	OrderID     unique.TaskOrderID
	ExecutionID string
	WorkerName  string
	ReturnCode  int32
	StatusCode  int32
}

//RouteChangeStateMsg - Request for setting a task into a specific state.
type RouteChangeStateMsg struct {
	Hold     bool
	Free     bool
	Rerun    bool
	SetOK    bool
	Username string
	OrderID  unique.TaskOrderID
}

//RouteChangeStateResponseMsg - Response for a change state
type RouteChangeStateResponseMsg struct {
	OrderID unique.TaskOrderID
	Message string
}

//RouteTaskCleanMsg -  Message for cleaning or termintating a task on remote worker
type RouteTaskCleanMsg struct {
	OrderID     unique.TaskOrderID
	ExecutionID string
	WorkerName  string
	Terminate   bool
}

//RouteJournalMsg -
type RouteJournalMsg struct {
	OrderID     unique.TaskOrderID
	ExecutionID string
	Time        time.Time
	Msg         string
}
