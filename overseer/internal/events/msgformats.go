package events

import (
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/types/unique"
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

//FlagActionData - struct for acquire flag message
type FlagActionData struct {
	Name   string
	Policy uint8
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

//RouteJournalMsg -
type RouteJournalMsg struct {
	OrderID     unique.TaskOrderID
	ExecutionID string
	Time        time.Time
	Msg         string
}
