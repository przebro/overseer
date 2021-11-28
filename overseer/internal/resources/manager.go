package resources

import (
	"errors"
	"overseer/common/core"
	"overseer/common/logger"
	"overseer/common/types/date"
	"overseer/datastore"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/taskdef"
	"regexp"
	"sort"
	"sync"
)

type resourceManager struct {
	dispatcher events.Dispatcher
	log        logger.AppLogger
	tstore     *resourceStore
	fstore     *resourceStore
	flock      sync.Mutex
}

//TicketManager - base resources required by task to run
type TicketManager interface {
	Add(name string, odate date.Odate) (bool, error)
	Delete(name string, odate date.Odate) (bool, error)
	Check(name string, odate date.Odate) bool
	ListTickets(name string, datestr string) []TicketResource
}

//FlagManager - Resource that helps run tasks
type FlagManager interface {
	Set(name string, policy FlagResourcePolicy) (bool, error)
	Unset(name string) (bool, error)
	DestroyFlag(name string) (bool, error)
	ListFlags(name string) []FlagResource
}

//ResourceManager - manages resources that are required by tasks
type ResourceManager interface {
	TicketManager
	FlagManager
	core.OverseerComponent
}

//NewManager - crates new resources manager
func NewManager(dispatcher events.Dispatcher, log logger.AppLogger, rconfig config.ResourcesConfigurartion, provider *datastore.Provider) (ResourceManager, error) {

	var err error

	var tstore *resourceStore
	var fstore *resourceStore

	trw, err := newTicketReadWriter(rconfig.TicketSource.Collection, "tickets", provider)
	if err != nil {
		return nil, err
	}

	tstore, err = newStore(trw, rconfig.TicketSource.Sync)
	if err != nil {
		return nil, err
	}

	frw, err := newFlagReadWriter(rconfig.FlagSource.Collection, "flags", provider)
	if err != nil {
		return nil, err
	}

	fstore, err = newStore(frw, rconfig.FlagSource.Sync)
	if err != nil {
		return nil, err
	}

	rm := &resourceManager{
		log:        log,
		dispatcher: dispatcher,
		tstore:     tstore,
		fstore:     fstore,
		flock:      sync.Mutex{},
	}

	//Subscribe for incoming messages about requests for tickets
	rm.dispatcher.Subscribe(events.RouteTicketCheck, rm)
	rm.dispatcher.Subscribe(events.RouteTicketIn, rm)
	rm.dispatcher.Subscribe(events.RouteFlagAcquire, rm)
	rm.dispatcher.Subscribe(events.RouteFlagRelase, rm)

	return rm, nil
}

func (rm *resourceManager) Add(name string, odate date.Odate) (bool, error) {

	err := rm.tstore.Insert(name+string(odate), TicketResource{Name: name, Odate: odate})
	if err != nil {
		return false, errors.New("ticket with given name and odate already exists")
	}
	rm.log.Info("TICKET:", name, odate)

	return true, nil
}
func (rm *resourceManager) Delete(name string, odate date.Odate) (bool, error) {

	key := name + string(odate)
	if err := rm.tstore.Delete(key); err != nil {
		return false, err
	}

	return true, nil
}
func (rm *resourceManager) Check(name string, odate date.Odate) bool {

	key := name + string(odate)
	_, ok := rm.tstore.Get(key)
	return ok
}

//ListTickets - return a list of tickets restricted to given name and odate
func (rm *resourceManager) ListTickets(name string, datestr string) []TicketResource {

	tickets := rm.tstore.All()
	var matchName bool
	var matchDate bool
	var err error

	result := make([]TicketResource, 0)
	nexpr := buildExpr(name)
	dexpr := buildDateExpr(datestr)

	for _, n := range tickets {

		if matchName, err = regexp.Match(nexpr, []byte(n.(TicketResource).Name)); err != nil {
			return []TicketResource{}
		}

		if matchDate, err = regexp.Match(dexpr, []byte(n.(TicketResource).Odate)); err != nil {
			return []TicketResource{}
		}

		if matchName && matchDate {
			result = append(result, n.(TicketResource))
		}

	}

	sort.Sort(ticketSorter{result})

	return result
}

//Set - change a value of a flag
func (rm *resourceManager) Set(name string, policy FlagResourcePolicy) (bool, error) {

	defer rm.flock.Unlock()
	rm.flock.Lock()

	var v interface{}
	var ok bool

	if v, ok = rm.fstore.Get(name); !ok {
		rm.log.Debug("CREATE FLAG:", name, "COUNT:", 1, "TYPE:", policy)
		rm.fstore.Insert(name, FlagResource{Name: name, Policy: policy, Count: 1})
		return true, nil
	}

	flag := v.(FlagResource)
	if flag.Policy == FlagPolicyExclusive {
		rm.log.Debug("ACQ ERR FLAG:", name, "TYPE:", policy, "flag in use with exclusive policy")
		return false, errors.New("flag in use with exclusive policy")
	}

	if flag.Policy == FlagPolicyShared && policy == FlagPolicyExclusive && flag.Count != 0 {
		rm.log.Debug("ACQ ERR FLAG SHR:", name, "TYPE:", flag.Policy, "flag in use with shared,trying exclusive")
		return false, errors.New("unable to set shared, flag in use with exclusive policy")
	}

	flag.Count++
	flag.Policy = policy
	rm.log.Debug("ACQ FLAG SUCCESS:", flag.Name, "TYPE:", flag.Policy, "COUNT:", flag.Count)
	rm.fstore.Update(flag.Name, flag)

	return true, nil
}

//Unset - remove a flag
func (rm *resourceManager) Unset(name string) (bool, error) {

	defer rm.flock.Unlock()
	rm.flock.Lock()

	var v interface{}
	var ok bool

	if v, ok = rm.fstore.Get(name); !ok {
		return false, errors.New("flag with given name does not exists")
	}

	flag := v.(FlagResource)
	flag.Count--
	rm.log.Debug("UNSET FLAG:", flag.Name, "COUNT:", flag.Count)
	if flag.Count == 0 {
		rm.fstore.Delete(name)
		rm.log.Debug("FLAG REMOVED")
	} else {
		rm.fstore.Update(flag.Name, flag)
	}

	return true, nil

}

//DestroyFlag - forcefully removes a flag
func (rm *resourceManager) DestroyFlag(name string) (bool, error) {

	defer rm.flock.Unlock()
	rm.flock.Lock()

	var ok bool

	if _, ok = rm.fstore.Get(name); !ok {
		return false, errors.New("flag with given name does not exists")
	}

	rm.fstore.Delete(name)

	return true, nil

}

func (rm *resourceManager) ListFlags(name string) []FlagResource {

	var matchName bool
	var err error
	flags := rm.fstore.All()
	result := make([]FlagResource, 0)

	nexpr := buildExpr(name)

	for _, n := range flags {

		if matchName, err = regexp.Match(nexpr, []byte(n.(FlagResource).Name)); err != nil {
			return []FlagResource{}
		}

		if matchName {
			result = append(result, n.(FlagResource))
		}
	}

	return result
}

//Start - starts the task pool
func (rm *resourceManager) Start() error {

	rm.tstore.start()
	rm.fstore.start()
	return nil
}

//Shutdown - shutdowns task pool
func (rm *resourceManager) Shutdown() error {

	rm.tstore.shutdown()
	rm.fstore.shutdown()

	return nil
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
	case events.RouteFlagAcquire:
		{
			data, ok := msg.Message().(events.RouteFlagAcquireMsg)
			if !ok {
				rm.log.Error("ResourceManager: route processing error, unexpected msg format")
				if receiver != nil {
					receiver.Done(events.ErrUnrecognizedMsgFormat)
				}
			}

			result := rm.processAcquireFlag(data)
			events.ResponseToReceiver(receiver, result)

		}
	case events.RouteFlagRelase:
		{
			data, ok := msg.Message().(events.RouteFlagAcquireMsg)
			if !ok {
				rm.log.Error("ResourceManager: route processing error, unexpected msg format")
				if receiver != nil {
					receiver.Done(events.ErrUnrecognizedMsgFormat)
				}
			}

			result := rm.processReleaseFlag(data)
			events.ResponseToReceiver(receiver, result)
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

func (rm *resourceManager) processAcquireFlag(data events.RouteFlagAcquireMsg) events.RouteFlagActionResponse {

	var requiredFlagNames []string = []string{}
	ok := true

	aflags := []struct {
		Name   string
		Policy FlagResourcePolicy
	}{}

	for _, f := range data.Flags {
		policy := FlagResourcePolicy(f.Policy)
		// if flag is not acquired, rollback changes that were made
		if ok, _ = rm.Set(f.Name, policy); !ok {
			requiredFlagNames = append(requiredFlagNames, f.Name)
			break
		}
		aflags = append(aflags, struct {
			Name   string
			Policy FlagResourcePolicy
		}{Name: f.Name, Policy: policy})
	}

	if !ok {

		for _, f := range aflags {
			rm.Unset(f.Name)
		}
	}

	return events.RouteFlagActionResponse{Success: ok, Names: requiredFlagNames}
}

func (rm *resourceManager) processReleaseFlag(data events.RouteFlagAcquireMsg) events.RouteFlagActionResponse {

	var flagNames []string = []string{}
	var success bool = true

	for _, f := range data.Flags {
		if ok, _ := rm.Unset(f.Name); !ok {
			flagNames = append(flagNames, f.Name)
			success = success && false
		} else {
			success = success && true
		}
	}

	return events.RouteFlagActionResponse{Success: success, Names: flagNames}
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

func buildExpr(value string) string {

	expr := ""

	if value == "" {
		return `[\w\-]*|^$`
	}

	expr += "^"

	for _, c := range value {

		if c == '*' {
			expr += `[\w\-]*`
			continue
		}
		if c == '?' {
			expr += `[\w\-]{1}`
			continue
		}

		expr += string(c)
	}

	expr += "$"

	return expr
}

func buildDateExpr(value string) string {

	expr := ""

	if value == "" {
		return `[\d]*|^$`
	}

	expr += "^"

	for _, c := range value {

		if c == '*' {
			expr += `[\d]*`
			continue
		}
		if c == '?' {
			expr += `[\d]{1}`
			continue
		}

		expr += string(c)
	}

	expr += "$"

	return expr
}
