package resources

import (
	"encoding/json"
	"errors"
	"goscheduler/common/logger"
	"goscheduler/overseer/internal/date"
	"goscheduler/overseer/internal/events"
	"goscheduler/overseer/internal/taskdef"
	"io/ioutil"
	"strings"
	"sync"
)

type resourceManager struct {
	path       string
	dispatcher events.Dispatcher
	log        logger.AppLogger
	Resources  *resourcePool
	tlock      sync.RWMutex
	flock      sync.RWMutex
}

//TicketManager - base resources required by task to run
type TicketManager interface {
	Add(name string, odate date.Odate) (bool, error)
	Delete(name string, odate date.Odate) (bool, error)
	Check(name string, odate date.Odate) bool
	ListTickets(name string, datestr string) []TicketResource
	tsync()
}

//FlagManager - Resource that helps run tasks
type FlagManager interface {
	Set(name string, policy FlagResourcePolicy) (bool, error)
	Unset(name string) (bool, error)
	ListFlags(name string) []FlagResource
	fsync()
}

//ResourceManager - manages resources that are required by tasks
type ResourceManager interface {
	TicketManager
	FlagManager
}

//NewManager - crates new resources manager
func NewManager(dispatcher events.Dispatcher, log logger.AppLogger, directory string) (ResourceManager, error) {

	var err error
	rm := new(resourceManager)
	rm.log = log
	rm.dispatcher = dispatcher
	rm.path = directory
	rm.tlock = sync.RWMutex{}
	rm.flock = sync.RWMutex{}

	rp, err := newResourcePool(directory)
	if err != nil {
		panic(err.Error())
	}
	rm.Resources = rp

	//Subscribe for incoming messages about requests for tickets
	rm.dispatcher.Subscribe(events.RouteTicketCheck, rm)
	rm.dispatcher.Subscribe(events.RouteTicketIn, rm)

	return rm, nil
}

func (rm *resourceManager) Add(name string, odate date.Odate) (bool, error) {

	defer rm.tlock.Unlock()
	rm.tlock.Lock()

	for _, n := range rm.Resources.Tickets {
		if n.Name == name && n.Odate == odate {
			rm.log.Info("TICKET:", n.Name, n.Odate)
			return false, errors.New("ticket with given name and odate already exists")
		}
	}

	rm.Resources.Tickets = append(rm.Resources.Tickets, TicketResource{Name: name, Odate: odate})
	rm.tsync()

	return true, nil
}
func (rm *resourceManager) Delete(name string, odate date.Odate) (bool, error) {

	defer rm.tlock.Unlock()
	rm.tlock.Lock()

	for i, n := range rm.Resources.Tickets {
		if n.Name == name && n.Odate == odate {
			rm.Resources.Tickets = append(rm.Resources.Tickets[:i], rm.Resources.Tickets[i+1:]...)
			rm.tsync()
			return true, nil
		}
	}
	return false, errors.New("unable to find given condition")

}
func (rm *resourceManager) Check(name string, odate date.Odate) bool {

	defer rm.tlock.RUnlock()
	rm.tlock.RLock()

	for _, n := range rm.Resources.Tickets {
		if n.Name == name && n.Odate == odate {
			return true
		}
	}
	return false

}

//ListTickets - return a list of tickets restricted to given name and odate
func (rm *resourceManager) ListTickets(name string, datestr string) []TicketResource {

	defer rm.tlock.RUnlock()
	rm.tlock.RLock()

	result := make([]TicketResource, 0)
	for _, n := range rm.Resources.Tickets {

		if n.Name == "" {
			result = append(result, n)
			continue
		}

		if strings.HasPrefix(n.Name, name) && strings.HasPrefix(string(n.Odate), datestr) {
			result = append(result, n)
		}
	}

	return result
}

//Set - change a value of a flag
func (rm *resourceManager) Set(name string, policy FlagResourcePolicy) (bool, error) {

	defer rm.flock.Unlock()
	rm.flock.Lock()

	for i, e := range rm.Resources.Flags {
		if e.Name == name {
			if e.Policy == FlagPolicyExclusive {
				return false, errors.New("flag in use with exclusive policy")
			}
			if e.Policy == FlagPolicyShared && policy == FlagPolicyExclusive {
				return false, errors.New("unable to set shared, flag in use with exclusive policy")
			}
			rm.Resources.Flags[i].Policy = policy
			rm.Resources.Flags[i].Count++

			rm.fsync()
			return true, nil
		}
	}
	rm.Resources.Flags = append(rm.Resources.Flags, FlagResource{Name: name, Policy: policy, Count: 1})
	rm.fsync()

	return true, nil
}

//Unset - remove a flag
func (rm *resourceManager) Unset(name string) (bool, error) {

	defer rm.flock.Unlock()
	rm.flock.Lock()

	for i, e := range rm.Resources.Flags {
		if e.Name == name {
			rm.Resources.Flags[i].Count = rm.Resources.Flags[i].Count - 1
			if rm.Resources.Flags[i].Count == 0 {
				rm.Resources.Flags = append(rm.Resources.Flags[:i], rm.Resources.Flags[i+1:]...)
			}
			rm.fsync()
			return true, nil
		}
	}

	return false, errors.New("Flag with given name does not exists")
}

func (rm *resourceManager) ListFlags(name string) []FlagResource {

	defer rm.flock.RUnlock()
	rm.flock.RLock()

	result := make([]FlagResource, 0)
	for _, e := range rm.Resources.Flags {
		if strings.HasPrefix(e.Name, name) {
			result = append(result, e)
		}
	}

	return result
}

func (rm *resourceManager) tsync() {

	data, err := json.Marshal(&rm.Resources.Tickets)
	if err != nil {
		rm.log.Error("Unable to unmarshal resources", err)

	}
	err = ioutil.WriteFile(rm.Resources.tpath, data, 0644)
	if err != nil {
		rm.log.Error("Unable to sync tickets resources with file", err)
	}

}
func (rm *resourceManager) fsync() {

	data, err := json.Marshal(&rm.Resources.Flags)
	if err != nil {
		rm.log.Error("Unable to unmarshal resources", err)

	}
	err = ioutil.WriteFile(rm.Resources.fpath, data, 0644)
	if err != nil {
		rm.log.Error("Unable to sync flags resources with file", err)
	}

}

func (rm *resourceManager) Process(receiver events.EventReceiver, route events.RouteName, msg events.DispatchedMessage) {

	switch route {
	case events.RouteTicketCheck:
		{
			rm.log.Debug("receiving from route:", route, msg.MsgID(), ",", msg.Created())
			data, isOk := msg.Message().(events.RouteTicketCheckMsgFormat)
			if !isOk {
				rm.log.Error("ResourceManager: route processing error, unexpected msg format")
				if receiver != nil {
					receiver.Done(events.ErrUnrecognizedMsgFormat)
				}
				return
			}

			rm.processCheckTicketEvent(data)
			events.ResponseToReceiver(receiver, data)
		}
	case events.RouteTicketIn:
		{
			data, isOk := msg.Message().(events.RouteTicketInMsgFormat)
			if !isOk {
				rm.log.Error("ResourceManager: route processing error, unexpected msg format")
				if receiver != nil {
					receiver.Done(events.ErrUnrecognizedMsgFormat)
				}

			}
			rm.processInTicketEvent(data)
			events.ResponseToReceiver(receiver, data)
		}
	default:
		{
			err := events.ErrInvalidRouteName
			rm.log.Debug(err)
			events.ResponseToReceiver(receiver, err)
		}
	}

}
func (rm *resourceManager) processCheckTicketEvent(data events.RouteTicketCheckMsgFormat) {

	for idx, d := range data.Tickets {
		data.Tickets[idx].Fulfilled = rm.Check(d.Name, date.Odate(d.Odate))
	}
}
func (rm *resourceManager) processInTicketEvent(data events.RouteTicketInMsgFormat) {

	for _, item := range data.Tickets {
		if item.Action == taskdef.OutActionAdd {
			_, err := rm.Add(item.Name, item.Odate)
			if err != nil {
				rm.log.Error(err)
			}
		}
		if item.Action == taskdef.OutActionRemove {
			_, err := rm.Delete(item.Name, item.Odate)
			if err != nil {
				rm.log.Error(err)
			}
		}
	}
}
