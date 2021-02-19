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
)

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

var storeConfig config.StoreProviderConfiguration = config.StoreProviderConfiguration{
	Store: []config.StoreConfiguration{
		{ID: "teststore", ConnectionString: "local;/../../../data/tests"},
	},
	Collections: []config.CollectionConfiguration{
		{Name: "tasks", StoreID: "teststore"},
		{Name: "sequence", StoreID: "teststore"},
	},
}

var testCollectionName = "tasks"
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
var mDispatcher = &mockDispatcher{Tickets: make(map[string]string, 0)}
var taskPoolT *ActiveTaskPool
var activeTaskManagerT *ActiveTaskPoolManager
var log logger.AppLogger = logger.NewTestLogger()

func init() {

	f, _ := os.Create("../../../data/tests/tasks.json")
	f.Write([]byte("{}"))
	f.Close()

	f2, _ := os.Create("../../../data/tests/sequence.json")
	f2.Write([]byte(`{}`))
	f2.Close()

	provider, _ = datastore.NewDataProvider(storeConfig)
	initTaskPool()
	taskPoolT.log = logger.NewTestLogger()
	path, _ := filepath.Abs("../../../def/")
	definitionManagerT, _ = taskdef.NewManager(path)
	activeTaskManagerT, _ = NewActiveTaskPoolManager(mDispatcher, definitionManagerT, taskPoolT, provider)
	activeTaskManagerT.log = log
	activeTaskManagerT.sequence = seq

}
func initTaskPool() {
	taskPoolT, _ = NewTaskPool(mDispatcher, taskPoolConfig, provider)
}
