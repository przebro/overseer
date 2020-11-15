package pool

import (
	"goscheduler/common/logger"
	"goscheduler/overseer/config"
	"goscheduler/overseer/internal/date"
	"goscheduler/overseer/internal/events"
	"goscheduler/overseer/internal/taskdef"
	"goscheduler/overseer/internal/unique"
	"sync"
	"testing"
	"time"
)

type mDispatcher struct {
	Tickets         map[string]string
	processNotEnded bool
	withError       bool
}

func (m *mDispatcher) PushEvent(receiver events.EventReceiver, route events.RouteName, msg events.DispatchedMessage) error {
	return nil
}

func (m *mDispatcher) Subscribe(route events.RouteName, participant events.EventParticipant) {

}
func (m *mDispatcher) Unsubscribe(route events.RouteName, participant events.EventParticipant) {

}

var dsp events.Dispatcher = &mDispatcher{}
var tpool *ActiveTaskPool = NewTaskPool(dsp, config.ActivePoolConfiguration{ForceNewDayProc: true, MaxOkReturnCode: 4, NewDayProc: "!0:30"})

func TestNewTaskPool(t *testing.T) {
	if tpool == nil {
		t.Error("TaskPool not initialized")
	}

}
func TestTaskAddGetDetailList(t *testing.T) {
	tpool.log = logger.NewTestLogger()

	orderID := unique.TaskOrderID("33333")
	builder := taskdef.DummyTaskBuilder{}

	odate, _ := date.AddDays(tpool.currentOdate, -2)
	def, _ := builder.WithBase("test", "task", "testdescription").WithRetention(0).Build()
	atask := &activeTask{TaskDefinition: def, orderID: orderID, orderDate: odate}
	atask.state = TaskStateEndedOk

	tpool.addTask(orderID, atask)

	if tpool.tasks.Len() != 1 {
		t.Error("unexpected result")
	}

	_, err := tpool.task(orderID)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	_, err = tpool.task(unique.TaskOrderID("22222"))
	if err == nil {
		t.Error("unexpected result:")
	}

	resultmsg, err := tpool.Detail(orderID)
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	if resultmsg.Name != "task" || resultmsg.Group != "test" {
		t.Error("Unexpected values, expected:", "task", "test", "actual", resultmsg.Name, resultmsg.Group)
	}

	_, err = tpool.Detail(unique.TaskOrderID("22222"))
	if err == nil {
		t.Error("Unexpected result")
	}

	lresult := tpool.List("")
	if len(lresult) != 1 {
		t.Error("Unexpected len")

	}

	tpool.tasks.Remove(orderID)

}
func TestProcessingFlag(t *testing.T) {

	tpool.PauseProcessing()
	if tpool.isProcActive != false {
		t.Error("Unexpected value")
	}

	tpool.ResumeProcessing()
	if tpool.isProcActive != true {
		t.Error("Unexpected value")
	}
}
func TestCleanUp(t *testing.T) {

	orderID := unique.TaskOrderID("12345")
	builder := taskdef.DummyTaskBuilder{}

	odate, _ := date.AddDays(tpool.currentOdate, -2)
	def, _ := builder.WithBase("test", "task", "testdescription").WithRetention(0).Build()
	atask := &activeTask{TaskDefinition: def, orderID: orderID, orderDate: odate}
	atask.state = TaskStateEndedOk
	tpool.addTask(orderID, atask)

	if tpool.tasks.Len() != 1 {
		t.Error("unexpected result")
	}

	tpool.cleanupCompletedTasks()

	if tpool.tasks.Len() != 0 {
		t.Error("unexpected result")
	}

}

func TestProcess(t *testing.T) {

	var rcverr error

	rcv := events.NewTicketCheckReceiver()

	msg := events.NewMsg("test data")
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_, rcverr = rcv.WaitForResult()
		wg.Done()
	}()

	tpool.Process(rcv, events.RouteTicketIn, msg)
	wg.Wait()
	if rcverr == nil {
		t.Error("Unexpected result")
	}

	x := time.Now()
	h, m, s := x.Clock()
	y, mth, d := x.Date()

	fmsg := events.NewMsg(events.RouteTaskStatusResponseMsg{})
	tpool.Process(nil, events.RouteTimeOut, fmsg)

	tmsg := events.NewMsg(events.RouteTimeOutMsgFormat{Year: y, Month: int(mth), Day: d, Hour: h, Min: m, Sec: s})
	tpool.Process(nil, events.RouteTimeOut, tmsg)
}
