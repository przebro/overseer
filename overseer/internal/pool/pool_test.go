package pool

import (
	"errors"
	"overseer/common/logger"
	"overseer/common/types"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/taskdef"
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

var definitionManagerT taskdef.TaskDefinitionManager
var mDispatcher = &mockDispatcher{Tickets: make(map[string]string, 0)}
var taskPoolT *ActiveTaskPool = NewTaskPool(mDispatcher, config.ActivePoolConfiguration{ForceNewDayProc: true, MaxOkReturnCode: 4, NewDayProc: "00:30"})
var activeTaskManagerT *ActiveTaskPoolManager
var log logger.AppLogger = logger.NewTestLogger()

func init() {

	taskPoolT.log = log
	path, _ := filepath.Abs("../../../def/")
	definitionManagerT, _ = taskdef.NewManager(path)
	activeTaskManagerT = NewActiveTaskPoolManager(mDispatcher, definitionManagerT, taskPoolT)
	activeTaskManagerT.log = log

}
