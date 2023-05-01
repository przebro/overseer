package pool

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TaskPoolTestSuite struct {
	suite.Suite
}

func TestTaskPoolTestSuite(t *testing.T) {
	suite.Run(t, new(TaskPoolTestSuite))
}

func (suite *TaskPoolTestSuite) SetupSuite() {

}

/*
func TestNewTaskPool(t *testing.T) {
	if taskPoolT == nil {
		t.Error("TaskPool not initialized")
	}

	taskPoolConfig.Collection = "invalid_collection"
	_, err := NewTaskPool(mDispatcher, taskPoolConfig, provider, true, logger.NewTestLogger(), definitionManagerT, &mockWorkerManager{}, &mockResourceManager{})
	if err == nil {
		t.Error("unexpected result")
	}

}
func TestTaskAddGetDetailList(t *testing.T) {
	t.Skip()
	orderID := unique.TaskOrderID("33333")
	builder := taskdef.DummyTaskBuilder{}

	odate := date.AddDays(date.CurrentOdate(), -2)
	def, _ := builder.WithBase("test", "task", "testdescription").WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).
		WithRetention(0).Build()
	atask := &activetask.TaskInstance{TaskDefinition: def,
		orderID: orderID, orderDate: odate,
		executions: []taskExecution{{ExecutionID: "ABCD"}},
		collected:  []types.CollectedTicketModel{{Name: "ABCDEF", Odate: "20201115"}},
	}

	atask.SetState(states.TaskStateEndedOk)
	taskPoolT.addTask(orderID, atask)

	if taskPoolT.tasks.len() != 1 {
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

	if len(resultmsg.Tickets) != 1 {
		t.Error("Unexpected tickets len:", len(resultmsg.Tickets))

	}

	_, err = taskPoolT.Detail(unique.TaskOrderID("22222"))
	if err == nil {
		t.Error("Unexpected result")
	}

	lresult := taskPoolT.List("")
	if len(lresult) != 1 {
		t.Error("Unexpected len:", len(lresult))
	}

	taskPoolT.tasks.remove(orderID)

}
func TestCleanUp(t *testing.T) {
	t.Skip()
	orderID := unique.TaskOrderID("12345")
	builder := taskdef.DummyTaskBuilder{}

	odate := date.AddDays(date.CurrentOdate(), -2)
	def, _ := builder.WithBase("test", "task", "testdescription").WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).WithRetention(0).Build()
	atask := &activeTask{TaskDefinition: def, orderID: orderID, orderDate: odate, executions: []taskExecution{{}}}
	atask.SetState(states.TaskStateEndedOk)
	taskPoolT.addTask(orderID, atask)

	if taskPoolT.tasks.len() != 1 {
		t.Error("unexpected result:", taskPoolT.tasks.len())
	}

	taskPoolT.CleanupCompletedTasks()

	if taskPoolT.tasks.len() != 0 {
		t.Error("unexpected result:", taskPoolT.tasks.len())
	}

}

func TestCycleTasks(t *testing.T) {

	_, err := activeTaskManagerT.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_05"}, date.CurrentOdate(), "user")

	if err != nil {
		t.Error("unepected result")
	}

	id, err := activeTaskManagerT.Force(taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: "test"}, Name: "dummy_05"}, date.CurrentOdate(), "user")

	if err != nil {
		t.Error("unepected result")
	}

	taskPoolT.tasks.store[unique.TaskOrderID(id)].SetState(states.TaskStateEndedNotOk)

	taskPoolT.cycleTasks(time.Now())
	time.Sleep(1 * time.Second)
}

func TestStartStopQR(t *testing.T) {

	taskPoolConfig.Collection = testCollectionName

	tpool, err := NewTaskPool(mDispatcher, taskPoolConfig, provider, false, logger.NewTestLogger(), definitionManagerT, &mockWorkerManager{}, &mockResourceManager{})
	if err != nil {
		t.Error("Unexpected result")
	}

	tpool.Start()
	tpool.Resume()
	if tpool.isProcActive != true {
		t.Error("Unexpected result:", tpool.isProcActive)
	}

	tpool.Quiesce()
	if tpool.isProcActive != false {
		t.Error("Unexpected result:", tpool.isProcActive)
	}

	tpool.Resume()

	tpool.Shutdown()
	if tpool.isProcActive != false {
		t.Error("Unexpected result:", tpool.isProcActive)
	}

}
*/
