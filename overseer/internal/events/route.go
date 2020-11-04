package events

import (
	"goscheduler/common/logger"
)

//RouteName - possible routes
type RouteName string

//Route names
const (
	RouteTimeOut  RouteName = "TICKER_OUT"
	RouteTicketIn RouteName = "COND_IN"
	//	RouteTicketOut       RouteName = "COND_OUT"
	RouteTicketCheck     RouteName = "COND_CHECK"
	RouteTaskAct         RouteName = "TASK_ACT"
	RouteChangeTaskState RouteName = "TASK_STATE"
	RouteWorkLaunch      RouteName = "WORK_LAUNCH"
	RouteWorkCheck       RouteName = "WORK_CHECK"
)

//messageRoute - holds participants of route
type messgeRoute struct {
	routename    RouteName
	participants []EventParticipant
}

//MessageRoute - Route definition, performs role of a topics
// and restrict	s dispatching of an events to specific subscribers
type MessageRoute interface {
	AddParticipant(p EventParticipant)
	Remove(p EventParticipant)
	PushMessage(receiver EventReceiver, msg DispatchedMessage)
}

//AddParticipant - adds a new conversation participant
func (route *messgeRoute) AddParticipant(p EventParticipant) {

	log := logger.Get()
	log.Debug("Append new participant to:", route.routename)
	route.participants = append(route.participants, p)

}

//PushMessage - sends message to all subscribers
func (route *messgeRoute) PushMessage(receiver EventReceiver, msg DispatchedMessage) {
	log := logger.Get()

	for _, r := range route.participants {
		log.Debug("Push message route:", route.routename, msg.MsgID(), ",", msg.Created())
		r.Process(receiver, route.routename, msg)
	}
}
func (route *messgeRoute) Remove(p EventParticipant) {

	log := logger.Get()

	for i, e := range route.participants {
		if e == p {
			route.participants = append(route.participants[:i], route.participants[i+1:]...)
			log.Debug("Remove participant from:", route.routename)
		}
	}

}
