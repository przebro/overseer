package events

import (
	"time"

	"github.com/przebro/overseer/overseer/internal/unique"
)

type dispatchedMessage struct {
	msgID         unique.MsgID
	correlationID unique.MsgID
	created       time.Time
	responseTo    RouteName
	message       interface{}
}

/*NewMsg - creates new DispatchedMessage without correlationID and response route
 */
func NewMsg(message interface{}) DispatchedMessage {

	id := unique.NewID()
	corrID := unique.None()
	e := dispatchedMessage{message: message, msgID: id, created: time.Now(), correlationID: corrID}
	return e
}

//NewCorrelatedMsg - creates new message in response to that with correlationID
func NewCorrelatedMsg(correlationID unique.MsgID, responseTo RouteName, message interface{}) DispatchedMessage {

	id := unique.NewID()

	e := dispatchedMessage{message: message, msgID: id, created: time.Now(), correlationID: correlationID, responseTo: responseTo}
	return e
}

//DispatchedMessage - Event data
type DispatchedMessage interface {
	Message() interface{}
	Created() time.Time
	ResponseTo() RouteName
	MsgID() unique.MsgID
	CorrelationID() unique.MsgID
}

//Message - returns message data
func (msg dispatchedMessage) Message() interface{} {

	return msg.message
}

//Created - returns time of creation of a message
func (msg dispatchedMessage) Created() time.Time {

	return msg.created
}

//ResponseTo - optional route name to response
func (msg dispatchedMessage) ResponseTo() RouteName {

	return msg.responseTo
}

//MsgID - return unique id of a message
func (msg dispatchedMessage) MsgID() unique.MsgID {

	return msg.msgID
}

//CorrelationID - msgID of an original message
func (msg dispatchedMessage) CorrelationID() unique.MsgID {

	return msg.correlationID
}
