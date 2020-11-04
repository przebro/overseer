package events

import (
	"errors"
	"goscheduler/common/logger"
	"sync"
)

//Dispatcher - dispatch messages between subscribed objects.
type Dispatcher interface {
	PushEvent(receiver EventReceiver, route RouteName, msg DispatchedMessage) error
	Subscribe(route RouteName, participant EventParticipant)
	Unsubscribe(route RouteName, participant EventParticipant)
}
type eventDipspatcher struct {
	msgRoutes map[RouteName]MessageRoute
	log       logger.AppLogger
	lock      sync.Mutex
}

//NewDispatcher - creates new Dispatcher
func NewDispatcher() Dispatcher {

	dispatcher := eventDipspatcher{
		msgRoutes: make(map[RouteName]MessageRoute),
		log:       logger.Get(),
		lock:      sync.Mutex{}}

	return &dispatcher
}

//PushEvent - Sends an events to a specific route
func (m *eventDipspatcher) PushEvent(receiver EventReceiver, routename RouteName, msg DispatchedMessage) error {
	defer m.lock.Unlock()
	m.lock.Lock()
	route, exists := m.msgRoutes[routename]

	if !exists {
		m.log.Error("Route not defined:", routename)
		return errors.New("Route not defined")
	}

	go route.PushMessage(receiver, msg)

	return nil
}

func (m *eventDipspatcher) Subscribe(routename RouteName, participant EventParticipant) {
	defer m.lock.Unlock()
	m.lock.Lock()

	route, exists := m.msgRoutes[routename]
	if !exists {
		m.log.Debug("Creating route:", routename)
		route = &messgeRoute{participants: make([]EventParticipant, 0), routename: routename}
		m.msgRoutes[routename] = route
	}

	route.AddParticipant(participant)

}
func (m *eventDipspatcher) Unsubscribe(routename RouteName, participant EventParticipant) {
	defer m.lock.Unlock()
	m.lock.Lock()
	route, exists := m.msgRoutes[routename]
	if exists {
		route.Remove(participant)
	}

}
