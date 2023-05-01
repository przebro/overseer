package pool

/*
var taskManager *ActiveTaskPoolManager

func TestNewManager(t *testing.T) {

	if taskManager == nil {
		t.Error("Unexpected error")
	}
}

func TestOrder(t *testing.T) {

	_, err := taskManager.Order(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "noexists"}, date.CurrentOdate(), "user")
	if err == nil {
		t.Error("Unexpected result")
	}

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"
	w := mockJournalT.Collect(1, time.Now())

	id, err := taskManager.Order(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), testUser)
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

	taskManager.pool.tasks.remove(unique.TaskOrderID(id))

}

func TestHoldFree(t *testing.T) {

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"
	w := mockJournalT.Collect(1, time.Now())

	orderid, err := taskManager.Order(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), testUser)
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
	_, err = taskManager.Hold(unique.TaskOrderID(orderid), testUser)
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

	_, err = taskManager.Hold(unique.TaskOrderID(orderid), testUser)

	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	w = mockJournalT.Collect(1, time.Now())
	_, err = taskManager.Free(unique.TaskOrderID(orderid), testUser)

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

	_, err = taskManager.Free(unique.TaskOrderID(orderid), testUser)

	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	taskManager.pool.tasks.remove(unique.TaskOrderID(orderid))

}
func TestConfirm(t *testing.T) {

	orderid, err := taskManager.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_04"}, date.CurrentOdate(), "user")
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"
	w := mockJournalT.Collect(1, time.Now())

	_, err = taskManager.Confirm(unique.TaskOrderID(orderid), testUser)
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

	_, err = taskManager.Confirm(unique.TaskOrderID(orderid), "user")
	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	_, err = taskManager.Confirm(unique.TaskOrderID("12345"), "user")
	if err == nil {
		t.Error("Unexpected result, expected error")
	}

	taskManager.pool.tasks.remove(unique.TaskOrderID(orderid))

}

// func TestSetOK(t *testing.T) {

// 	testUser := "testUserT1"

// 	orderid, err := taskManager.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), testUser)
// 	if err != nil {
// 		t.Error("Unexpected result:", err)
// 	}

// 	var jrnalMsg []events.RouteJournalMsg

// 	_, err = taskManager.SetOk(unique.TaskOrderID(orderid), testUser)
// 	if err == nil {
// 		t.Error("Unexpected result")
// 	}

// 	w := mockJournalT.Collect(1, time.Now())
// 	taskManager.pool.tasks.store[unique.TaskOrderID(orderid)].executions[0].state = models.TaskStateEndedNotOk

// 	_, err = taskManager.SetOk(unique.TaskOrderID(orderid), testUser)
// 	if err != nil {
// 		t.Error("Unexpected result:", err)
// 	}

// 	if jrnalMsg = <-w; jrnalMsg == nil {
// 		t.Error("unexpected result, failed to collect data from journal")
// 		t.FailNow()
// 	}

// 	if !strings.HasPrefix(jrnalMsg[0].Msg, "TASK SETOK") || !strings.Contains(jrnalMsg[0].Msg, fmt.Sprintf("user:%s", testUser)) {
// 		t.Error("unexpected result, invalid journal entry:", jrnalMsg[0].Msg)
// 	}

// 	taskManager.pool.tasks.remove(unique.TaskOrderID(orderid))

// }

func TestSetWorkerName(t *testing.T) {

	id, err := taskManager.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), "user")
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	task, _ := taskManager.pool.tasks.get(unique.TaskOrderID(id))

	task.SetWorkerName("test_workername")
	if task.WorkerName() != "test_workername" {
		t.Error("unexpected result:", task.WorkerName(), "expected: test_workername")
	}
	taskManager.pool.tasks.remove(unique.TaskOrderID(id))
}

func TestStartEndTime(t *testing.T) {

	id, err := taskManager.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), "user")
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	task, _ := taskManager.pool.tasks.get(unique.TaskOrderID(id))

	now := time.Now()

	task.SetStartTime()
	task.SetEndTime()
	if task.StartTime().Before(now) {
		t.Error("unexpected value:", task.StartTime(), " is before:", now)
	}

	if task.EndTime().Before(now) {
		t.Error("unexpected value:", task.EndTime(), " is before:", now)
	}
	taskManager.pool.tasks.remove(unique.TaskOrderID(id))
}

func TestGetModel(t *testing.T) {

	id, err := taskManager.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), "user")
	if err != nil {
		t.Error("Unexpected result:", err)
	}

	task, _ := taskManager.pool.tasks.get(unique.TaskOrderID(id))

	model := task.GetModel()
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

	badModel := activetask.ActiveTaskModel{}

	badModel.Tickets = nil

	if _, err = activetask.FromModel(badModel, definitionManagerT); err == nil {
		t.Error("unexpected result")
	}

	def, err := activetask.FromModel(model, definitionManagerT)
	if err != nil {
		t.Error("unexpected result")
	}

	if def.OrderID() != task.OrderID() {
		t.Error("unexpected result")
	}

	if def.OrderDate() != task.OrderDate() {
		t.Error("unexpected result")
	}

	taskManager.pool.tasks.remove(unique.TaskOrderID(id))
}

func TestForce(t *testing.T) {

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"

	_, err := taskManager.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "noexists"}, date.CurrentOdate(), testUser)

	if err == nil {
		t.Error("Unexpected result")
	}

	w := mockJournalT.Collect(1, time.Now())

	id, err := taskManager.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_03"}, date.CurrentOdate(), testUser)
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

	taskManager.pool.tasks.remove(unique.TaskOrderID(id))

}

func TestChangeState(t *testing.T) {

}

func TestOrderGroup(t *testing.T) {

	_, err := taskManager.OrderGroup(taskdata.GroupData{Group: "invalid"}, date.CurrentOdate(), "user")

	if err == nil {
		t.Error("unepected result")
	}

	_, err = taskManager.OrderGroup(taskdata.GroupData{Group: "test"}, date.CurrentOdate(), "user")
	if err != nil {
		t.Error("unepected result")
	}

	for i := range taskManager.pool.tasks.store {
		taskManager.pool.tasks.remove(i)
	}
}
func TestForceGroup(t *testing.T) {

	_, err := taskManager.ForceGroup(taskdata.GroupData{Group: "invalid"}, date.CurrentOdate(), "user")

	if err == nil {
		t.Error("unepected result")
	}

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"
	w := mockJournalT.Collect(5, time.Now())

	_, err = taskManager.ForceGroup(taskdata.GroupData{Group: "test"}, date.CurrentOdate(), testUser)
	if err != nil {
		t.Error("unepected result")
	}

	if jrnalMsg = <-w; jrnalMsg == nil {
		t.Error("unexpected result, failed to collect data from journal")
		t.FailNow()
	}

	for i := range taskManager.pool.tasks.store {
		taskManager.pool.tasks.remove(i)
	}
}

func TestEnforceTask(t *testing.T) {

	id, err := taskManager.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_05"}, date.CurrentOdate(), "user")

	if err != nil {
		t.Error("unepected result")
	}

	_, err = taskManager.Enforce(unique.TaskOrderID("99999"), "user")
	if err == nil {
		t.Error("unepected result")
	}

	taskManager.pool.tasks.store[unique.TaskOrderID(id)].SetState(models.TaskStateEndedNotOk)

	_, err = taskManager.Enforce(unique.TaskOrderID(id), "user")
	if err == nil {
		t.Error("unepected result")
	}

	var jrnalMsg []events.RouteJournalMsg
	testUser := "testUserT1"
	w := mockJournalT.Collect(1, time.Now())

	taskManager.pool.tasks.store[unique.TaskOrderID(id)].SetState(models.TaskStateWaiting)

	_, err = taskManager.Enforce(unique.TaskOrderID(id), testUser)
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

	r := taskManager.OrderNewTasks()
	if r == 0 {
		t.Error("unexpected result")
	}

	for i := range taskManager.pool.tasks.store {
		taskManager.pool.tasks.remove(i)
	}
}
*/
