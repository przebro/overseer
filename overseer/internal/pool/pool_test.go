package pool

import (
	"errors"
	"fmt"
	"os"
	"overseer/common/logger"
	"overseer/common/types"
	"overseer/datastore"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/taskdef"
	"overseer/overseer/internal/unique"
	"path/filepath"
	"time"
)

type mockJournal struct {
	timeout   int
	collected chan events.RouteJournalMsg
}

func (j *mockJournal) Push(msg events.RouteJournalMsg) {
	j.collected <- msg

}

func (j *mockJournal) Collect(expected int, after time.Time) <-chan []events.RouteJournalMsg {

	ch := make(chan []events.RouteJournalMsg)

	go func(timeout, expected int, done chan<- []events.RouteJournalMsg, col chan events.RouteJournalMsg) {

		collected := 0
		deadline := time.After(time.Duration(timeout) * time.Second)
		result := []events.RouteJournalMsg{}
		for {
			select {
			case <-deadline:
				{
					close(done)
					return
				}
			case d := <-col:
				{
					if d.Time.Before(after) {
						continue
					}
					result = append(result, d)
					collected++
					if collected == expected {
						done <- result
						close(done)
						return
					}

				}
			}
		}

	}(j.timeout, expected, ch, j.collected)

	return ch

}

type mockDispatcher struct {
	Tickets         map[string]string
	processNotEnded bool
	withError       bool
}

func (m *mockDispatcher) PushEvent(receiver events.EventReceiver, route events.RouteName, msg events.DispatchedMessage) error {

	go func() {
		if route == events.RouteWorkLaunch {
			if m.withError {
				receiver.Done(errors.New(""))
			} else {
				dat := events.RouteWorkResponseMsg{
					Status: types.WorkerTaskStatusExecuting,
				}
				events.ResponseToReceiver(receiver, dat)
			}

		}
		if route == events.RouteTicketIn {

		}

		if route == events.RoutTaskJournal {
			r := msg.Message().(events.RouteJournalMsg)
			mockJournalT.Push(r)
		}

		if route == events.RouteWorkCheck {

			if m.withError {
				receiver.Done(errors.New(""))
			} else {

				_, iskOk := msg.Message().(events.WorkRouteCheckStatusMsg)
				if iskOk == false {
					events.ResponseToReceiver(receiver, errors.New(""))
				}
				if m.processNotEnded {
					receiver.Done(events.RouteWorkResponseMsg{Status: types.WorkerTaskStatusExecuting, ReturnCode: 0})
				} else {
					receiver.Done(events.RouteWorkResponseMsg{Status: types.WorkerTaskStatusEnded, ReturnCode: 0})
				}

			}
		}
		if route == events.RouteTicketCheck {

			if m.withError {

				receiver.Done(errors.New(""))

			} else {
				result, iskOk := msg.Message().(events.RouteTicketCheckMsgFormat)
				if iskOk == false {
					events.ResponseToReceiver(receiver, errors.New(""))
				}

				for i, t := range result.Tickets {

					_, exists := m.Tickets[t.Name]
					if exists {
						result.Tickets[i].Fulfilled = true
					}
				}

				receiver.Done(result)

			}

		}
	}()
	return nil
}
func (m *mockDispatcher) Subscribe(route events.RouteName, participant events.EventParticipant) {

}
func (m *mockDispatcher) Unsubscribe(route events.RouteName, participant events.EventParticipant) {

}

const (
	testCollectionName = "tasks"
	testStoreTaskName  = "storetasks"
	testSequenceName   = "sequence"
)

var storeConfig config.StoreProviderConfiguration = config.StoreProviderConfiguration{
	Store: []config.StoreConfiguration{
		{ID: "teststore", ConnectionString: "local;/../../../data/tests"},
		{ID: "teststoretasks", ConnectionString: "local;/../../../data/tests?updatesync=true"},
	},
	Collections: []config.CollectionConfiguration{
		{Name: testCollectionName, StoreID: "teststore"},
		{Name: testStoreTaskName, StoreID: "teststoretasks"},
		{Name: testSequenceName, StoreID: "teststore"},
	},
}

var taskPoolConfig config.ActivePoolConfiguration = config.ActivePoolConfiguration{
	ForceNewDayProc: true, MaxOkReturnCode: 4,
	NewDayProc: "00:30",
	SyncTime:   5,
	Collection: testCollectionName,
}

type mockSequence struct {
	val int
}

func (m *mockSequence) Next() unique.TaskOrderID {

	m.val++
	return unique.TaskOrderID(fmt.Sprintf("%05d", m.val))
}

var seq = &mockSequence{val: 1}

var provider *datastore.Provider

var definitionManagerT taskdef.TaskDefinitionManager
var mDispatcher = &mockDispatcher{Tickets: make(map[string]string)}
var taskPoolT *ActiveTaskPool
var activeTaskManagerT *ActiveTaskPoolManager
var log logger.AppLogger = logger.NewTestLogger()
var mockJournalT = &mockJournal{timeout: 3, collected: make(chan events.RouteJournalMsg, 10)}

func init() {

	f, _ := os.Create(fmt.Sprintf("../../../data/tests/%s.json", testCollectionName))
	f.Write([]byte("{}"))
	f.Close()

	f1, _ := os.Create(fmt.Sprintf("../../../data/tests/%s.json", testStoreTaskName))
	f1.Write([]byte("{}"))
	f1.Close()

	f2, _ := os.Create(fmt.Sprintf("../../../data/tests/%s.json", testCollectionName))
	f2.Write([]byte(`{}`))
	f2.Close()

	provider, _ = datastore.NewDataProvider(storeConfig, log)
	initTaskPool()
	taskPoolT.log = logger.NewTestLogger()
	path, _ := filepath.Abs("../../../def/")
	definitionManagerT, _ = taskdef.NewManager(path, log)
	activeTaskManagerT, _ = NewActiveTaskPoolManager(mDispatcher, definitionManagerT, taskPoolT, provider, log)
	activeTaskManagerT.log = log
	activeTaskManagerT.sequence = seq

}
func initTaskPool() {
	taskPoolT, _ = NewTaskPool(mDispatcher, taskPoolConfig, provider, true, log)
}
