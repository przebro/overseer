package journal

import (
	"context"
	"sync"
	"time"

	"github.com/przebro/overseer/common/core"
	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/config"
	"github.com/przebro/overseer/overseer/internal/events"
	"github.com/przebro/overseer/overseer/internal/unique"

	"github.com/przebro/databazaar/collection"
)

//Enumeration of task's messages that will be sent to journal when a specific event occurs
const (
	TaskHeld              = "TASK HELD, user:%s"
	TaskFreed             = "TASK FREED, user:%s"
	TaskFulfill           = "TASK PRECONDITIONS OK"
	TaskEnforce           = "TASK ENFORCED, user:%s"
	TaskRerun             = "TASK RERUN, user:%s"
	TaskSetOK             = "TASK SETOK, user:%s"
	TaskConfirmed         = "TASK CONFIRMED, user:%s"
	TaskForced            = "TASK FORCED, user:%s ODATE:%s"
	TaskOrdered           = "TASK ORDERED, user:%s ODATE:%s"
	TaskStartingRN        = "TASK STARTING RN:%d"
	TaskStartingFailedErr = "TASK STARTING FAILED worker error"
	TaskStartingFailed    = "TASK STARTING FAILED invalid worker status:%d"
	TaskStarting          = "TASK STARTING worker:%s"
	TaskComplete          = "TASK EXECUTION COMPLETE  %s"
	TaskFailed            = "TASK FAILED worker failure"
	TaskEndedNOK          = "ENDED NOT OK, RC:%d, STATUS:%d"
	TaskEndedOK           = "ENDED OK, RC:%d, STATUS:%d"
	TaskPostProc          = "TASK POST PROCESSING ends"
)

type mLogModel struct {
	TaskID  string     `json:"_id" bson:"_id"`
	Entries []LogEntry `json:"entries" bson:"entries"`
	Tstamp  time.Time  `json:"time" bson:"time"`
}

//LogEntry - represents a single task event
type LogEntry struct {
	Time        time.Time `json:"time" bson:"time"`
	ExecutionID string    `json:"eid"  bson:"eid"`
	Message     string    `json:"data" bson:"data"`
}

//TaskLogReader -reads task log
type TaskLogReader interface {
	ReadLog(id unique.TaskOrderID) []LogEntry
}

//TaskLogWriter - wrties task log
type TaskLogWriter interface {
	WriteLog(id unique.TaskOrderID, entry LogEntry)
}

//TaskJournal - writes and reads task journal data
type TaskJournal interface {
	TaskLogReader
	TaskLogWriter
	core.OverseerComponent
}

type taskLogJournal struct {
	conf     config.JournalConfiguration
	col      collection.DataCollection
	store    map[unique.TaskOrderID]mLogModel
	lock     sync.Mutex
	log      logger.AppLogger
	shutdown chan struct{}
	done     <-chan struct{}
}

//NewTaskJournal - creates a new instance of a TaskJournal
func NewTaskJournal(conf config.JournalConfiguration, dispatcher events.Dispatcher, provider *datastore.Provider, log logger.AppLogger) (TaskJournal, error) {

	var err error
	var col collection.DataCollection

	if col, err = provider.GetCollection(conf.LogCollection); err != nil {
		return nil, err
	}

	sch := make(chan struct{})

	journal := &taskLogJournal{col: col,
		store:    map[unique.TaskOrderID]mLogModel{},
		lock:     sync.Mutex{},
		log:      log,
		conf:     conf,
		shutdown: sch,
	}

	if dispatcher != nil {
		dispatcher.Subscribe(events.RoutTaskJournal, journal)
	}

	return journal, nil
}

func (journal *taskLogJournal) ReadLog(id unique.TaskOrderID) []LogEntry {

	var entries []LogEntry = []LogEntry{}
	var logs mLogModel = mLogModel{}
	var ok bool

	defer journal.lock.Unlock()
	journal.lock.Lock()

	//if there are no entries for given task, check if were stored
	if logs, ok = journal.store[id]; !ok {
		//nothing
		if err := journal.col.Get(context.Background(), string(id), &logs); err != nil {
			return entries
		}

		logs.Tstamp = time.Now()
	}

	for _, n := range logs.Entries {
		entries = append(entries, n)
	}

	return entries
}
func (journal *taskLogJournal) WriteLog(id unique.TaskOrderID, entry LogEntry) {

	logs := mLogModel{}
	var err error
	var ok bool

	defer journal.lock.Unlock()
	journal.lock.Lock()

	if logs, ok = journal.store[id]; !ok {

		if err = journal.col.Get(context.Background(), string(id), &logs); err != nil {
			logs.TaskID = string(id)
		}
	}
	logs.Tstamp = time.Now()
	logs.TaskID = string(id)
	logs.Entries = append(logs.Entries, entry)
	journal.store[id] = logs
}

func (journal *taskLogJournal) watch(interval int, shutdown <-chan struct{}) <-chan struct{} {

	inform := make(chan struct{})

	if interval == 0 {
		go func(sc <-chan struct{}, inf chan<- struct{}) {

			<-sc
			journal.sync(interval, time.Now())
			close(inf)
			return

		}(shutdown, inform)

		return inform
	}

	go func(sc <-chan struct{}, inf chan<- struct{}) {

		for {
			select {
			case t := <-time.After(time.Duration(interval) * time.Second):
				{

					journal.sync(interval, t)
				}
			case <-sc:
				{
					journal.sync(interval, time.Now())
					close(inf)
					return
				}
			}
		}

	}(shutdown, inform)

	return inform

}

func (journal *taskLogJournal) sync(interval int, t time.Time) {

	journal.lock.Lock()

	for _, n := range journal.store {
		if n.Tstamp.Add(time.Duration(interval) * time.Second).Before(t) {

			journal.col.Update(context.Background(), &n)
		}
	}

	journal.lock.Unlock()
}

//Process - receive notification from dispatcher
func (journal *taskLogJournal) Process(receiver events.EventReceiver, routename events.RouteName, msg events.DispatchedMessage) {

	switch routename {
	case events.RoutTaskJournal:
		{
			var ok bool
			var msgdata events.RouteJournalMsg
			if msgdata, ok = msg.Message().(events.RouteJournalMsg); !ok {
				journal.log.Error(events.ErrUnrecognizedMsgFormat)
				events.ResponseToReceiver(receiver, events.ErrUnrecognizedMsgFormat)
				break
			}

			journal.WriteLog(msgdata.OrderID, LogEntry{Time: msgdata.Time, ExecutionID: msgdata.ExecutionID, Message: msgdata.Msg})
			events.ResponseToReceiver(receiver, "")
		}
	default:
		{
			err := events.ErrInvalidRouteName
			journal.log.Error(err)
			events.ResponseToReceiver(receiver, err)
		}
	}
}

func (journal *taskLogJournal) Start() error {

	journal.done = journal.watch(journal.conf.SyncTime, journal.shutdown)

	return nil
}

func (journal *taskLogJournal) Shutdown() error {

	journal.shutdown <- struct{}{}
	<-journal.done
	close(journal.shutdown)

	return nil
}
