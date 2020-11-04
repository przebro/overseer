package events

import (
	"errors"
)

//EventParticipant - Participant of events dispatching
type EventParticipant interface {
	Process(receiver EventReceiver, sroutename RouteName, msg DispatchedMessage)
}

//EventReceiver - Receives information about processing of a message.
type EventReceiver interface {
	Done(data interface{})
}
type eventReceiver struct {
	done chan interface{}
}

//Done - Informs receiver about the result of a message processing. It could be an expected message or an error.
func (r *eventReceiver) Done(data interface{}) {
	r.done <- data
}

//ActiveTaskReceiver - Receiver for RouteTaskAct
type ActiveTaskReceiver interface {
	EventReceiver
	WaitForResult() (RouteTaskActionResponseFormat, error)
}

type activeTaskListReceiver struct {
	eventReceiver
}

//NewActiveTaskReceiver - Creates new ActiveTaskReceiver
func NewActiveTaskReceiver() ActiveTaskReceiver {
	l := &activeTaskListReceiver{}
	l.done = make(chan interface{})
	return l
}

//WaitForResult - Waits for a result of the processing.
func (r *activeTaskListReceiver) WaitForResult() (RouteTaskActionResponseFormat, error) {
	var result RouteTaskActionResponseFormat
	var err error = nil

	x := <-r.done
	switch val := x.(type) {
	case RouteTaskActionResponseFormat:
		{
			result = val
		}
	case error:
		{
			err = val
		}
	default:
		{
			err = errors.New("Unrecognized type")
		}
	}

	return result, err
}

//TicketCheckReciever - Receiver for a RouteTicketCheck
type TicketCheckReciever interface {
	EventReceiver
	WaitForResult() (RouteTicketCheckMsgFormat, error)
}
type ticketCheckReciever struct {
	eventReceiver
}

//NewTicketCheckReceiver - Creates a new TicketCheckReciever
func NewTicketCheckReceiver() TicketCheckReciever {
	l := &ticketCheckReciever{}
	l.done = make(chan interface{})
	return l
}

func (r *ticketCheckReciever) WaitForResult() (RouteTicketCheckMsgFormat, error) {

	var result RouteTicketCheckMsgFormat
	var err error = nil

	x := <-r.done
	switch val := x.(type) {
	case RouteTicketCheckMsgFormat:
		{
			result = val
		}
	case error:
		{
			err = val
		}
	default:
		{
			err = errors.New("Unrecognized type")
		}
	}
	return result, err
}

//WorkLaunchReceiver - Receiver for a RouteWorkLaunch
type WorkLaunchReceiver interface {
	EventReceiver
	WaitForResult() (RouteWorkResponseMsg, error)
}

type workLaunchReceiver struct {
	eventReceiver
}

//NewWorkLaunchReceiver - Creates a new WorkLaunchReceiver
func NewWorkLaunchReceiver() WorkLaunchReceiver {
	l := &workLaunchReceiver{}
	l.done = make(chan interface{})
	return l
}

func (r *workLaunchReceiver) WaitForResult() (RouteWorkResponseMsg, error) {

	var result RouteWorkResponseMsg
	var err error = nil

	x := <-r.done
	switch val := x.(type) {
	case RouteWorkResponseMsg:
		{
			result = val
		}
	case error:
		{
			err = val
		}
	default:
		{
			err = errors.New("Unrecognized type")
		}
	}
	return result, err
}

//ChangeTaskStateReceiver - Receiver for RouteChangeTaskState
type ChangeTaskStateReceiver interface {
	EventReceiver
	WaitForResult() (RouteChangeStateResponseMsg, error)
}

type changeTaskStateReceiver struct {
	eventReceiver
}

//NewChangeTaskStateReceiver - Creates a new NewChangeTaskStateReceiver
func NewChangeTaskStateReceiver() ChangeTaskStateReceiver {
	l := &changeTaskStateReceiver{}
	l.done = make(chan interface{})
	return l
}

func (r *changeTaskStateReceiver) WaitForResult() (RouteChangeStateResponseMsg, error) {

	var result RouteChangeStateResponseMsg
	var err error = nil

	x := <-r.done
	switch val := x.(type) {
	case RouteChangeStateResponseMsg:
		{
			result = val
		}
	case error:
		{
			err = val
		}
	default:
		{
			err = errors.New("Unrecognized type")
		}
	}
	return result, err
}

//ResponseToReceiver - Helper function, sends response to receiver. If receiver is nil this function does nothing
func ResponseToReceiver(r EventReceiver, data interface{}) {

	if r != nil {
		r.Done(data)
	}

}
