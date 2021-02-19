package pool

import (
	"overseer/common/types/date"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/unique"
	"overseer/overseer/taskdata"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {

	if activeTaskManagerT == nil {
		t.Error("Unexpected error")
	}
}

func TestOrder(t *testing.T) {

	_, err := activeTaskManagerT.Order(taskdata.GroupNameData{Group: "test", Name: "noexists"}, date.CurrentOdate())
	if err == nil {
		t.Error("Unexpected result")
	}

	id, err := activeTaskManagerT.Order(taskdata.GroupNameData{Group: "test", Name: "dummy_03"}, date.CurrentOdate())
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(id))

}

func TestHoldFree(t *testing.T) {

	orderid, err := activeTaskManagerT.Order(taskdata.GroupNameData{Group: "test", Name: "dummy_03"}, date.CurrentOdate())
	if err != nil {
		t.Error("Unexpected result:", err)
	}
	_, err = activeTaskManagerT.Hold(unique.TaskOrderID(orderid))
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	_, err = activeTaskManagerT.Hold(unique.TaskOrderID(orderid))

	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	_, err = activeTaskManagerT.Free(unique.TaskOrderID(orderid))

	if err != nil {
		t.Error("Unexpected result:", err)
	}

	_, err = activeTaskManagerT.Free(unique.TaskOrderID(orderid))

	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(orderid))

}
func TestConfirm(t *testing.T) {
	orderid, err := activeTaskManagerT.Force(taskdata.GroupNameData{Group: "test", Name: "dummy_04"}, date.CurrentOdate())
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	_, err = activeTaskManagerT.Confirm(unique.TaskOrderID(orderid))
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	_, err = activeTaskManagerT.Confirm(unique.TaskOrderID(orderid))
	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	_, err = activeTaskManagerT.Confirm(unique.TaskOrderID("12345"))
	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(orderid))

}

func TestSetWorkerName(t *testing.T) {

	id, err := activeTaskManagerT.Force(taskdata.GroupNameData{Group: "test", Name: "dummy_03"}, date.CurrentOdate())
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	task, _ := activeTaskManagerT.pool.tasks.get(unique.TaskOrderID(id))

	task.SetWorkerName("test_workername")
	if task.WorkerName() != "test_workername" {
		t.Error("unexpected result:", task.WorkerName(), "expected: test_workername")
	}
	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(id))
}

func TestStartEndTime(t *testing.T) {

	id, err := activeTaskManagerT.Force(taskdata.GroupNameData{Group: "test", Name: "dummy_03"}, date.CurrentOdate())
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	task, _ := activeTaskManagerT.pool.tasks.get(unique.TaskOrderID(id))

	now := time.Now()

	task.SetStartTime()
	task.SetEndTime()
	if task.StartTime().Before(now) {
		t.Error("unexpected value:", task.StartTime(), " is before:", now)
	}

	if task.EndTime().Before(now) {
		t.Error("unexpected value:", task.EndTime(), " is before:", now)
	}
	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(id))
}

func TestGetModel(t *testing.T) {

	id, err := activeTaskManagerT.Force(taskdata.GroupNameData{Group: "test", Name: "dummy_03"}, date.CurrentOdate())
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	task, _ := activeTaskManagerT.pool.tasks.get(unique.TaskOrderID(id))

	model := task.getModel()
	if model.OrderID != id {
		t.Error("unexpected result, model does not match to  the origin task, orderID is different")
	}
	if model.OrderDate != date.CurrentOdate() {
		t.Error("unexpected result, model does not match to  the origin task, ODATE is different")
	}
	if model.Confirmed != true {
		t.Error("unexpected result, model does not match to  the origin task, Confirmed is different")
	}

	if len(model.Tickets) != 1 {
		t.Error("unexpected result, model does not match to  the origin task, In Tickets are different")
	}
	if model.Tickets[0].Name != "IN-DUMMY03" {
		t.Error("unexpected result, model does not match to  the origin task, Ticket name is different")
	}

	if model.Tickets[0].Odate != string(date.CurrentOdate()) {
		t.Error("unexpected result, model does not match to  the origin task, Ticket ODATE is different")
	}

	badModel := activeTaskModel{Definition: []byte("{}")}

	if _, err = fromModel(badModel); err == nil {
		t.Error("unexpected result")
	}

	badModel.Tickets = nil

	if _, err = fromModel(badModel); err == nil {
		t.Error("unexpected result")
	}

	def, err := fromModel(model)
	if err != nil {
		t.Error("unexpected result")
	}

	if def.orderID != task.orderID {
		t.Error("unexpected result")
	}

	if def.orderDate != task.orderDate {
		t.Error("unexpected result")
	}

	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(id))
}

/*
{
    "type" : "dummy",
    "name" :"dummy_03",
    "group" : "test",
    "description" :"sample dummy task definition",
    "inticket" : [
        {"name" : "IN-DUMMY03","odate" : "ODATE" }
    ],
    "relation" :"AND",
    "flags" : [],
    "outticket" :[
        {"name" : "OK-DUMMY03","odate" : "ODATE" ,"action":"ADD"},
        {"name" : "IN-DUMMY03","odate" : "ODATE" ,"action":"REM"}

    ],
    "schedule" :{
        "type" : "manual",
        "from" : "",
        "to" : ""
    },
    "other" : {"name" :"xxx","val" :"yyyy"}

}

*/

func TestForce(t *testing.T) {

	_, err := activeTaskManagerT.Force(taskdata.GroupNameData{Group: "test", Name: "noexists"}, date.CurrentOdate())

	t.Log(err)

	if err == nil {
		t.Error("Unexpected result")
	}

	id, err := activeTaskManagerT.Force(taskdata.GroupNameData{Group: "test", Name: "dummy_03"}, date.CurrentOdate())
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(id))

}

func TestAtmProcess(t *testing.T) {

	receiver := events.NewChangeTaskStateReceiver()
	msg := events.NewMsg("")

	go func() {
		activeTaskManagerT.Process(receiver, events.RouteWorkLaunch, msg)
	}()

	_, err := receiver.WaitForResult()
	if err != events.ErrInvalidRouteName {
		t.Log("Unexpected result:", err)
	}

}

func TestChangeState(t *testing.T) {

	id, _ := activeTaskManagerT.Order(taskdata.GroupNameData{Group: "test", Name: "dummy_03"}, date.CurrentOdate())
	orderID := unique.TaskOrderID(id)

	receiver := events.NewChangeTaskStateReceiver()

	//Change state with invalid taks id
	msg := events.NewMsg(events.RouteChangeStateMsg{Hold: true, OrderID: unique.TaskOrderID("54321")})

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err := receiver.WaitForResult()
	if err == nil {
		t.Log("Unexpected result:", err)
	}

	//Test change state with invalid message
	msg = events.NewMsg("")

	go func(m events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, m)
	}(msg)

	_, err = receiver.WaitForResult()
	if err != events.ErrUnrecognizedMsgFormat {
		t.Log("Unexpected result actual:", err, "expected", events.ErrUnrecognizedMsgFormat)
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{Hold: true, OrderID: orderID})

	go func(m events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, m)
	}(msg)

	_, err = receiver.WaitForResult()
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{Free: true, OrderID: orderID})

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err != nil {
		t.Log("Unexpected result:", err)
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{Free: true, OrderID: unique.TaskOrderID("54321")})

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Log("Unexpected result:")
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{SetOK: true, OrderID: orderID})
	task, _ := activeTaskManagerT.pool.tasks.get(orderID)
	task.state = TaskStateEndedNotOk

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err != nil {
		t.Log("Unexpected result:", err)
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{SetOK: true, OrderID: orderID})

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Log("Unexpected result:", err)
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{SetOK: true, OrderID: unique.TaskOrderID("54321")})

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Log("Unexpected result:")
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{Rerun: true, OrderID: orderID})

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err != nil {
		t.Log("Unexpected result:")
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{Rerun: true, OrderID: orderID})
	task.state = TaskStateWaiting

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Log("Unexpected result:")
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{Rerun: true, OrderID: unique.TaskOrderID("54321")})
	task.state = TaskStateWaiting

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Log("Unexpected result")
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{Rerun: true, SetOK: true, OrderID: orderID})
	task.state = TaskStateWaiting

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Log("Unexpected result:", err)
	}

	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(id))

}

func TestProcesAddTask(t *testing.T) {

	receiver := events.NewActiveTaskReceiver()

	msg := events.NewMsg("")

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteTaskAct, msg)
	}(msg)

	_, err := receiver.WaitForResult()
	if err != events.ErrUnrecognizedMsgFormat {
		t.Log("Unexpected result:", err)
	}

	msg = events.NewMsg(events.RouteTaskActionMsgFormat{Group: "test", Name: "xyz"})

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteTaskAct, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Log("Unexpected result:")
	}

	msg = events.NewMsg(events.RouteTaskActionMsgFormat{Group: "test", Name: "dummy_03", Force: true, Odate: date.CurrentOdate()})

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteTaskAct, msg)
	}(msg)

	result, err := receiver.WaitForResult()
	if err != nil {
		t.Log("Unexpected result:", err)
	}
	taskID := result.Data[0].TaskID

	activeTaskManagerT.pool.tasks.remove(taskID)

	msg = events.NewMsg(events.RouteTaskActionMsgFormat{Group: "test", Name: "dummy_03", Force: false, Odate: date.CurrentOdate()})

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteTaskAct, msg)
	}(msg)

	result, err = receiver.WaitForResult()
	if err != nil {
		t.Log("Unexpected result:", err)
	}
	taskID = result.Data[0].TaskID

	activeTaskManagerT.pool.tasks.remove(taskID)
}
