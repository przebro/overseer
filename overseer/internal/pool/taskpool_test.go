package pool

import (
	"overseer/overseer/internal/date"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/taskdef"
	"overseer/overseer/internal/unique"
	"sync"
	"testing"
	"time"
)

func TestNewTaskPool(t *testing.T) {
	if taskPoolT == nil {
		t.Error("TaskPool not initialized")
	}

}
func TestTaskAddGetDetailList(t *testing.T) {

	orderID := unique.TaskOrderID("33333")
	builder := taskdef.DummyTaskBuilder{}

	odate, _ := date.AddDays(taskPoolT.currentOdate, -2)
	def, _ := builder.WithBase("test", "task", "testdescription").WithRetention(0).Build()
	atask := &activeTask{TaskDefinition: def, orderID: orderID, orderDate: odate}
	atask.state = TaskStateEndedOk

	taskPoolT.addTask(orderID, atask)

	if taskPoolT.tasks.Len() != 1 {
		t.Error("unexpected result")
	}

	_, err := taskPoolT.task(orderID)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	_, err = taskPoolT.task(unique.TaskOrderID("22222"))
	if err == nil {
		t.Error("unexpected result:")
	}

	resultmsg, err := taskPoolT.Detail(orderID)
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	if resultmsg.Name != "task" || resultmsg.Group != "test" {
		t.Error("Unexpected values, expected:", "task", "test", "actual", resultmsg.Name, resultmsg.Group)
	}

	_, err = taskPoolT.Detail(unique.TaskOrderID("22222"))
	if err == nil {
		t.Error("Unexpected result")
	}

	lresult := taskPoolT.List("")
	if len(lresult) != 1 {
		t.Error("Unexpected len")

	}

	taskPoolT.tasks.Remove(orderID)

}
func TestProcessingFlag(t *testing.T) {

	taskPoolT.PauseProcessing()
	if taskPoolT.isProcActive != false {
		t.Error("Unexpected value")
	}

	taskPoolT.ResumeProcessing()
	if taskPoolT.isProcActive != true {
		t.Error("Unexpected value")
	}
}
func TestCleanUp(t *testing.T) {

	orderID := unique.TaskOrderID("12345")
	builder := taskdef.DummyTaskBuilder{}

	odate, _ := date.AddDays(taskPoolT.currentOdate, -2)
	def, _ := builder.WithBase("test", "task", "testdescription").WithRetention(0).Build()
	atask := &activeTask{TaskDefinition: def, orderID: orderID, orderDate: odate}
	atask.state = TaskStateEndedOk
	taskPoolT.addTask(orderID, atask)

	if taskPoolT.tasks.Len() != 1 {
		t.Error("unexpected result")
	}

	taskPoolT.cleanupCompletedTasks()

	if taskPoolT.tasks.Len() != 0 {
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

	taskPoolT.Process(rcv, events.RouteTicketIn, msg)
	wg.Wait()
	if rcverr == nil {
		t.Error("Unexpected result")
	}

	x := time.Now()
	h, m, s := x.Clock()
	y, mth, d := x.Date()

	fmsg := events.NewMsg(events.RouteTaskStatusResponseMsg{})
	taskPoolT.Process(nil, events.RouteTimeOut, fmsg)

	tmsg := events.NewMsg(events.RouteTimeOutMsgFormat{Year: y, Month: int(mth), Day: d, Hour: h, Min: m, Sec: s})
	taskPoolT.Process(nil, events.RouteTimeOut, tmsg)
}
