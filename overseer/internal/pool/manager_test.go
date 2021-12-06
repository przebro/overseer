package pool

import (
	"fmt"
	"overseer/common/types/date"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/unique"
	"overseer/overseer/taskdata"
	"strings"
	"testing"
	"time"
)

func init() {
	if !isInitialized {
		setupEnv()
	}
}

func TestNewManager(t *testing.T) {

	if activeTaskManagerT == nil {
		t.Error("Unexpected error")
	}
}

func TestOrder(t *testing.T) {

	_, err := activeTaskManagerT.Order(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "noexists"}, date.CurrentOdate(), "user")
	if err == nil {
		t.Error("Unexpected result")
	}

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"
	w := mockJournalT.Collect(1, time.Now())

	id, err := activeTaskManagerT.Order(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), testUser)
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	if jrnalMsg = <-w; jrnalMsg == nil {
		t.Error("unexpected result, failed to collect data from journal")
		t.FailNow()
	}

	if !strings.HasPrefix(jrnalMsg[0].Msg, "TASK ORDERED") || !strings.Contains(jrnalMsg[0].Msg, fmt.Sprintf("user:%s", testUser)) {
		t.Error("unexpected result, invalid journal entry:", jrnalMsg[0].Msg)
	}

	pastDate := date.AddDays(date.CurrentOdate(), -1)
	_, err = activeTaskManagerT.Order(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_05"}, pastDate, "user")
	if err == nil {
		t.Error("Unexpected result:", err)
	}

	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(id))

}

func TestHoldFree(t *testing.T) {

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"
	w := mockJournalT.Collect(1, time.Now())

	orderid, err := activeTaskManagerT.Order(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), testUser)
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	if jrnalMsg = <-w; jrnalMsg == nil {
		t.Error("unexpected result, failed to collect data from journal")
		t.FailNow()
	}

	if !strings.HasPrefix(jrnalMsg[0].Msg, "TASK ORDERED") || !strings.Contains(jrnalMsg[0].Msg, fmt.Sprintf("user:%s", testUser)) {
		t.Error("unexpected result, invalid journal entry:", jrnalMsg[0].Msg)
	}

	w = mockJournalT.Collect(1, time.Now())
	_, err = activeTaskManagerT.Hold(unique.TaskOrderID(orderid), testUser)
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	if jrnalMsg = <-w; jrnalMsg == nil {
		t.Error("unexpected result, failed to collect data from journal")
		t.FailNow()
	}

	if !strings.HasPrefix(jrnalMsg[0].Msg, "TASK HELD") || !strings.Contains(jrnalMsg[0].Msg, fmt.Sprintf("user:%s", testUser)) {
		t.Error("unexpected result, invalid journal entry:", jrnalMsg[0].Msg)
	}

	_, err = activeTaskManagerT.Hold(unique.TaskOrderID(orderid), testUser)

	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	w = mockJournalT.Collect(1, time.Now())
	_, err = activeTaskManagerT.Free(unique.TaskOrderID(orderid), testUser)

	if err != nil {
		t.Error("Unexpected result:", err)
	}

	if jrnalMsg = <-w; jrnalMsg == nil {
		t.Error("unexpected result, failed to collect data from journal")
		t.FailNow()
	}

	if !strings.HasPrefix(jrnalMsg[0].Msg, "TASK FREED") || !strings.Contains(jrnalMsg[0].Msg, fmt.Sprintf("user:%s", testUser)) {
		t.Error("unexpected result, invalid journal entry:", jrnalMsg[0].Msg)
	}

	_, err = activeTaskManagerT.Free(unique.TaskOrderID(orderid), testUser)

	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(orderid))

}
func TestConfirm(t *testing.T) {

	orderid, err := activeTaskManagerT.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_04"}, date.CurrentOdate(), "user")
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"
	w := mockJournalT.Collect(1, time.Now())

	_, err = activeTaskManagerT.Confirm(unique.TaskOrderID(orderid), testUser)
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	if jrnalMsg = <-w; jrnalMsg == nil {
		t.Error("unexpected result, failed to collect data from journal")
		t.FailNow()
	}

	if !strings.HasPrefix(jrnalMsg[0].Msg, "TASK CONFIRMED") || !strings.Contains(jrnalMsg[0].Msg, fmt.Sprintf("user:%s", testUser)) {
		t.Error("unexpected result, invalid journal entry:", jrnalMsg[0].Msg)
	}

	_, err = activeTaskManagerT.Confirm(unique.TaskOrderID(orderid), "user")
	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	_, err = activeTaskManagerT.Confirm(unique.TaskOrderID("12345"), "user")
	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(orderid))

}

func TestSetOK(t *testing.T) {

	testUser := "testUserT1"

	orderid, err := activeTaskManagerT.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), testUser)
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	var jrnalMsg []events.RouteJournalMsg

	_, err = activeTaskManagerT.SetOk(unique.TaskOrderID(orderid), testUser)
	if err == nil {
		t.Error("Unexpected result")
	}

	w := mockJournalT.Collect(1, time.Now())
	activeTaskManagerT.pool.tasks.store[unique.TaskOrderID(orderid)].executions[0].state = TaskStateEndedNotOk

	_, err = activeTaskManagerT.SetOk(unique.TaskOrderID(orderid), testUser)
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	if jrnalMsg = <-w; jrnalMsg == nil {
		t.Error("unexpected result, failed to collect data from journal")
		t.FailNow()
	}

	if !strings.HasPrefix(jrnalMsg[0].Msg, "TASK SETOK") || !strings.Contains(jrnalMsg[0].Msg, fmt.Sprintf("user:%s", testUser)) {
		t.Error("unexpected result, invalid journal entry:", jrnalMsg[0].Msg)
	}

	activeTaskManagerT.pool.tasks.remove(unique.TaskOrderID(orderid))

}

func TestSetWorkerName(t *testing.T) {

	id, err := activeTaskManagerT.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), "user")
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

	id, err := activeTaskManagerT.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), "user")
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

	id, err := activeTaskManagerT.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), "user")
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

	badModel := activeTaskModel{}

	badModel.Tickets = nil

	if _, err = fromModel(badModel, definitionManagerT); err == nil {
		t.Error("unexpected result")
	}

	def, err := fromModel(model, definitionManagerT)
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

func TestForce(t *testing.T) {

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"

	_, err := activeTaskManagerT.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "noexists"}, date.CurrentOdate(), testUser)

	if err == nil {
		t.Error("Unexpected result")
	}

	w := mockJournalT.Collect(1, time.Now())

	id, err := activeTaskManagerT.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), testUser)
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	if jrnalMsg = <-w; jrnalMsg == nil {
		t.Error("unexpected result, failed to collect data from journal")
		t.FailNow()
	}

	if !strings.HasPrefix(jrnalMsg[0].Msg, "TASK FORCED") || !strings.Contains(jrnalMsg[0].Msg, fmt.Sprintf("user:%s", testUser)) {
		t.Error("unexpected result, invalid journal entry:", jrnalMsg[0].Msg)
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

	id, _ := activeTaskManagerT.Order(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), "user")
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
	task.SetState(TaskStateEndedNotOk)

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
	task.SetState(TaskStateWaiting)

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Log("Unexpected result:")
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{Rerun: true, OrderID: unique.TaskOrderID("54321")})
	task.SetState(TaskStateWaiting)

	go func(msg events.DispatchedMessage) {
		activeTaskManagerT.Process(receiver, events.RouteChangeTaskState, msg)
	}(msg)

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Log("Unexpected result")
	}

	msg = events.NewMsg(events.RouteChangeStateMsg{Rerun: true, SetOK: true, OrderID: orderID})
	task.SetState(TaskStateWaiting)

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

func TestOrderGroup(t *testing.T) {

	_, err := activeTaskManagerT.OrderGroup(taskdata.GroupData{Group: "invalid"}, date.CurrentOdate(), "user")

	if err == nil {
		t.Error("unepected result")
	}

	_, err = activeTaskManagerT.OrderGroup(taskdata.GroupData{Group: "test"}, date.CurrentOdate(), "user")
	if err != nil {
		t.Error("unepected result")
	}

	for i := range activeTaskManagerT.pool.tasks.store {
		activeTaskManagerT.pool.tasks.remove(i)
	}
}
func TestForceGroup(t *testing.T) {

	_, err := activeTaskManagerT.ForceGroup(taskdata.GroupData{Group: "invalid"}, date.CurrentOdate(), "user")

	if err == nil {
		t.Error("unepected result")
	}

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"
	w := mockJournalT.Collect(5, time.Now())

	_, err = activeTaskManagerT.ForceGroup(taskdata.GroupData{Group: "test"}, date.CurrentOdate(), testUser)
	if err != nil {
		t.Error("unepected result")
	}

	if jrnalMsg = <-w; jrnalMsg == nil {
		t.Error("unexpected result, failed to collect data from journal")
		t.FailNow()
	}

	for i := range activeTaskManagerT.pool.tasks.store {
		activeTaskManagerT.pool.tasks.remove(i)
	}
}

func TestEnforceTask(t *testing.T) {

	id, err := activeTaskManagerT.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_05"}, date.CurrentOdate(), "user")

	if err != nil {
		t.Error("unepected result")
	}

	_, err = activeTaskManagerT.Enforce(unique.TaskOrderID("99999"), "user")
	if err == nil {
		t.Error("unepected result")
	}

	activeTaskManagerT.pool.tasks.store[unique.TaskOrderID(id)].SetState(TaskStateEndedNotOk)

	_, err = activeTaskManagerT.Enforce(unique.TaskOrderID(id), "user")
	if err == nil {
		t.Error("unepected result")
	}

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"
	w := mockJournalT.Collect(1, time.Now())

	activeTaskManagerT.pool.tasks.store[unique.TaskOrderID(id)].SetState(TaskStateWaiting)

	_, err = activeTaskManagerT.Enforce(unique.TaskOrderID(id), testUser)
	if err != nil {
		t.Error("unepected result")
	}

	if jrnalMsg = <-w; jrnalMsg == nil {
		t.Error("unexpected result, failed to collect data from journal")
		t.FailNow()
	}

	if !strings.HasPrefix(jrnalMsg[0].Msg, "TASK ENFORCED") || !strings.Contains(jrnalMsg[0].Msg, fmt.Sprintf("user:%s", testUser)) {
		t.Error("unexpected result, invalid journal entry:", jrnalMsg[0].Msg)
	}

}

func TestOrderNewTasks(t *testing.T) {

	r := activeTaskManagerT.orderNewTasks()
	if r == 0 {
		t.Error("unexpected result")
	}

	for i := range activeTaskManagerT.pool.tasks.store {
		activeTaskManagerT.pool.tasks.remove(i)
	}
}
