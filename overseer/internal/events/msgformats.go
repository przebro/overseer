package events

import (
	"overseer/overseer/internal/date"
	task "overseer/overseer/internal/taskdef"
	"overseer/overseer/internal/unique"
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
	WaitingInfo string
}

//TaskDetailResultMsg - Result for a request for detailed information
type TaskDetailResultMsg struct {
	TaskInfoResultMsg
	Hold      bool
	Confirm   bool
	RunNumber int32
	StartTime string
	EndTime   string
	Worker    string
	Output    []string
}

//RouteTaskActionMsgFormat - Request for a task order
type RouteTaskActionMsgFormat struct {
	Group  string
	Name   string
	TaskID unique.TaskOrderID
	Force  bool
	Odate  date.Odate
}

//RouteTaskActionResponseFormat - Response message for a task order or force.
type RouteTaskActionResponseFormat struct {
	Data []TaskInfoResultMsg
}

//WorkRouteCheckStatusMsg - Response with information about the status of a work
type WorkRouteCheckStatusMsg struct {
	OrderID    unique.TaskOrderID
	WorkerName string
}

//RouteTaskStatusResponseMsg - Response for a task status
type RouteTaskStatusResponseMsg struct {
	TaskID     string
	Data       []string
	Ended      bool
	ReturnCode int32
}

//RouteTaskExecutionMsg - Contains informations needed to begin a work on a remoteworker.
type RouteTaskExecutionMsg struct {
	OrderID   unique.TaskOrderID
	Type      string
	Variables []task.VariableData
	Command   interface{}
}

//RouteWorkResponseMsg - Contains information about the status of executing work.
type RouteWorkResponseMsg struct {
	Output     []string
	WorkerName string
	Started    bool
	Ended      bool
	ReturnCode int32
}

//RouteChangeStateMsg - Request for setting a task into a specific state.
type RouteChangeStateMsg struct {
	Hold    bool
	Free    bool
	Rerun   bool
	SetOK   bool
	OrderID unique.TaskOrderID
}

//RouteChangeStateResponseMsg - Response for a change state
type RouteChangeStateResponseMsg struct {
	OrderID unique.TaskOrderID
	Message string
}
