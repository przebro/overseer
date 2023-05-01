package services

import (
	"context"
	"fmt"
	"io"
	"testing"

	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/overseer/auth"
	"github.com/przebro/overseer/overseer/internal/events"
	"github.com/przebro/overseer/overseer/internal/journal"
	"github.com/przebro/overseer/overseer/internal/pool"
	"github.com/przebro/overseer/overseer/taskdata"
	"github.com/przebro/overseer/proto/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"google.golang.org/grpc/codes"
)

type mockPoolManager struct {
	mock.Mock
}

func (m *mockPoolManager) OrderGroup(groupdata taskdata.GroupData, odate date.Odate, username string) ([]string, error) {
	args := m.Called(groupdata, odate, username)
	return args.Get(0).([]string), args.Error(1)
}
func (m *mockPoolManager) Order(task taskdata.GroupNameData, odate date.Odate, username string) (string, error) {
	args := m.Called(task, odate, username)
	return args.Get(0).(string), args.Error(1)
}
func (m *mockPoolManager) Force(task taskdata.GroupNameData, odate date.Odate, username string) (string, error) {
	args := m.Called(task, odate, username)
	return args.Get(0).(string), args.Error(1)
}
func (m *mockPoolManager) Enforce(id unique.TaskOrderID, username string) (string, error) {
	args := m.Called(id, username)
	return args.Get(0).(string), args.Error(1)
}
func (m *mockPoolManager) Rerun(id unique.TaskOrderID, username string) (string, error) {
	args := m.Called(id, username)
	return args.Get(0).(string), args.Error(1)
}
func (m *mockPoolManager) Hold(id unique.TaskOrderID, username string) (string, error) {
	args := m.Called(id, username)
	return args.Get(0).(string), args.Error(1)
}
func (m *mockPoolManager) Free(id unique.TaskOrderID, username string) (string, error) {
	args := m.Called(id, username)
	return args.Get(0).(string), args.Error(1)
}
func (m *mockPoolManager) SetOk(id unique.TaskOrderID, username string) (string, error) {
	args := m.Called(id, username)
	return args.Get(0).(string), args.Error(1)
}
func (m *mockPoolManager) Confirm(id unique.TaskOrderID, username string) (string, error) {
	args := m.Called(id, username)
	return args.Get(0).(string), args.Error(1)
}

type mockTaskViewer struct {
	mock.Mock
}

func (m *mockTaskViewer) Detail(orderID unique.TaskOrderID) (events.TaskDetailResultMsg, error) {
	args := m.Called(orderID)
	return args.Get(0).(events.TaskDetailResultMsg), args.Error(1)

}
func (m *mockTaskViewer) List(filter string) []events.TaskInfoResultMsg {
	args := m.Called(filter)
	return args.Get(0).([]events.TaskInfoResultMsg)
}

type mockJournal struct{}

var jrnl *mockJournal = &mockJournal{}

func (m *mockJournal) WriteLog(id unique.TaskOrderID, entry journal.LogEntry) {}
func (m *mockJournal) ReadLog(id unique.TaskOrderID) []journal.LogEntry       { return []journal.LogEntry{} }
func (m *mockJournal) Start() error                                           { return nil }
func (m *mockJournal) Shutdown() error                                        { return nil }
func (m *mockJournal) Resume() error                                          { return nil }
func (m *mockJournal) Quiesce() error                                         { return nil }

type mockListTaskStream struct {
	MockGrpcServerStream
	channel chan *services.TaskListResultMsg
}

func (m *mockListTaskStream) Send(input *services.TaskListResultMsg) error {
	m.channel <- input
	return nil
}
func (m *mockListTaskStream) Recv() (*services.TaskListResultMsg, error) {

	if len(m.channel) == 0 {
		return nil, io.EOF
	}
	result := <-m.channel
	return result, nil
}

type activeTaskTestSuite struct {
	suite.Suite
	serviceServer services.TaskServiceServer
	ovsService    *ovsActiveTaskService
	poolManager   *mockPoolManager
	poolViewer    *mockTaskViewer
}

func TestTaskSuite(t *testing.T) {
	suite.Run(t, new(activeTaskTestSuite))
}
func (suite *activeTaskTestSuite) SetupSuite() {

	suite.poolViewer = &mockTaskViewer{}
	suite.poolManager = &mockPoolManager{}
	suite.serviceServer = NewTaskService(suite.poolManager, suite.poolViewer, jrnl)
	suite.ovsService = suite.serviceServer.(*ovsActiveTaskService)
}

func (suite *activeTaskTestSuite) TestOrderTask_Errors() {

	service := suite.serviceServer
	ctx := context.Background()
	ctx = context.WithValue(ctx, "username", "<anonymous>")

	_, err := service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: "abcdef", TaskName: "dummy_01"})
	suite.NotNil(err, "unexpected result")
	_, code := matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "#$SDtest", Odate: string(date.CurrentOdate()), TaskName: "dummy_01"})
	suite.NotNil(err, "unexpected result")
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "%^&$%dummy_01"})
	suite.NotNil(err, "unexpected result")
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)
}

func (suite *activeTaskTestSuite) TestOrderTask_Success() {

	input := &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_04"}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "username", "<anonymous>")

	suite.poolManager.On("Order", taskdata.GroupNameData{
		Name:      input.TaskName,
		GroupData: taskdata.GroupData{Group: input.TaskGroup},
	}, date.Odate(input.Odate), "<anonymous>").Return("12345", nil)

	service := suite.serviceServer

	r, err := service.OrderTask(ctx, input)
	suite.Nil(err, "unexpected result", err)
	suite.True(r.Success)
	suite.Contains(r.Message, "TaskID:")

}

func (suite *activeTaskTestSuite) TestForceTask_Errors() {

	service := suite.serviceServer
	ctx := context.Background()

	_, err := service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: "abcdef", TaskName: "dummy_01"})
	suite.NotNil(err)
	_, code := matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "#$SDtest", Odate: string(date.CurrentOdate()), TaskName: "dummy_01"})
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "%^&$%dummy_01"})
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

}

func (suite *activeTaskTestSuite) TestForceTask() {

	service := suite.serviceServer
	ctx := context.Background()
	ctx = context.WithValue(ctx, "username", "<anonymous>")
	input := &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_04"}

	suite.poolManager.On("Force", taskdata.GroupNameData{
		Name:      input.TaskName,
		GroupData: taskdata.GroupData{Group: input.TaskGroup},
	}, date.Odate(input.Odate), "<anonymous>").Return("12345", nil)

	r, err := service.ForceTask(ctx, input)
	suite.Nil(err)
	suite.True(r.Success)
	suite.Contains(r.Message, "TaskID:")
}

func (suite *activeTaskTestSuite) TestListTask() {

	service := suite.serviceServer

	suite.poolViewer.On("List", "").Return(
		[]events.TaskInfoResultMsg{
			{
				TaskID:      "12345",
				Odate:       date.CurrentOdate(),
				Group:       "Test",
				Name:        "Test_Task_01",
				State:       0,
				RunNumber:   1,
				Held:        false,
				Confirmed:   true,
				WaitingInfo: "",
			},
			{
				TaskID:      "55555",
				Odate:       date.CurrentOdate(),
				Group:       "Test",
				Name:        "Test_Task_02",
				State:       0,
				RunNumber:   1,
				Held:        false,
				Confirmed:   true,
				WaitingInfo: "",
			},
		},
	)
	out := &mockListTaskStream{channel: make(chan *services.TaskListResultMsg, 10)}
	err := service.ListTasks(&services.TaskFilterMsg{}, out)
	suite.Nil(err)
	collected := []*services.TaskListResultMsg{}

	for {
		r, err := out.Recv()
		if err == io.EOF {
			break
		}
		collected = append(collected, r)
	}

	suite.Len(collected, 2)
}

func (suite *activeTaskTestSuite) TestOrderGroup_Errors() {

	service := suite.serviceServer
	ctx := context.Background()

	_, err := service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "_4@42!5terfds", Odate: string(date.CurrentOdate())})
	suite.NotNil(err)
	_, code := matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "TEST", Odate: "ABCDEF"})
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "", Odate: string(date.CurrentOdate())})
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "ABCDED", Odate: ""})
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.Unauthenticated)
	suite.Equal(codes.Unauthenticated, code)

	suite.poolManager.On("OrderGroup",
		taskdata.GroupData{
			Group: "ABCDED"}, date.CurrentOdate(), "<anonymous>").Return([]string{}, pool.ErrUnableFindGroup)

	ctx = context.WithValue(ctx, "username", "<anonymous>")
	_, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "ABCDED", Odate: ""})
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.Internal)
	suite.Equal(codes.Internal, code)
}

func (suite *activeTaskTestSuite) TestOrderGroup() {

	service := suite.serviceServer
	ctx := context.WithValue(context.Background(), "username", "<anonymous>")

	suite.poolManager.On("OrderGroup",
		taskdata.GroupData{
			Group: "test"}, date.CurrentOdate(), "<anonymous>").
		Return([]string{"TASK_01", "Task_02"}, nil)

	_, err := service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "test", Odate: ""})
	suite.Nil(err)

}

func (suite *activeTaskTestSuite) TestForceGroup_Errors() {

	service := suite.serviceServer
	ctx := context.Background()

	_, err := service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "_4@42!5terfds", Odate: string(date.CurrentOdate())})
	suite.NotNil(err)
	_, code := matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "TEST", Odate: "ABCDEF"})
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "", Odate: string(date.CurrentOdate())})
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "ABCDED", Odate: ""})
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.Unauthenticated)
	suite.Equal(codes.Unauthenticated, code)

	suite.poolManager.On("OrderGroup",
		taskdata.GroupData{
			Group: "ABCDED"}, date.CurrentOdate(), "<anonymous>").Return([]string{}, pool.ErrUnableFindGroup)

	ctx = context.WithValue(ctx, "username", "<anonymous>")

	_, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "ABCDED", Odate: ""})
	suite.NotNil(err)

}

func (suite *activeTaskTestSuite) TestForceGroup() {

	service := suite.serviceServer
	ctx := context.WithValue(context.Background(), "username", "<anonymous>")

	suite.poolManager.On("OrderGroup",
		taskdata.GroupData{
			Group: "test_force"}, date.CurrentOdate(), "<anonymous>").Return([]string{"", ""}, nil)

	r, err := service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "test_force", Odate: ""})

	suite.Nil(err)
	suite.True(r.Success)

}
func (suite *activeTaskTestSuite) TestConfirmTask_Errors() {
	service := suite.serviceServer
	ctx := context.Background()

	_, err := service.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
	suite.NotNil(err)
	_, code := matchExpectedStatusFromError(err, codes.InvalidArgument)
	suite.Equal(codes.InvalidArgument, code)

	_, err = service.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: "ABCDE"})
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.Unauthenticated)
	suite.Equal(codes.Unauthenticated, code)

	suite.poolManager.On("Confirm", unique.TaskOrderID("12345"), "<anonymous>").
		Return("", pool.ErrUnableFindTask)

	ctx = context.WithValue(ctx, "username", "<anonymous>")
	_, err = service.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: "12345"})
	suite.NotNil(err)

}
func (suite *activeTaskTestSuite) TestConfirmTask_Success() {

	service := suite.serviceServer
	ctx := context.WithValue(context.Background(), "username", "<anonymous>")
	suite.poolManager.On("Confirm", unique.TaskOrderID("123ab"), "<anonymous>").
		Return("", nil)

	_, err := service.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: "123ab"})
	suite.Nil(err)

}

/*
	func TestHoldFree_Errors_NotFound(t *testing.T) {
		service := createTaskService(t)
		ctx := context.Background()

		_, err := service.HoldTask(ctx, &services.TaskActionMsg{TaskID: "ABCDE"})
		if err == nil {
			t.Error("unexpected result:", err)
		}

		if ok, code := matchExpectedStatusFromError(err, codes.NotFound); !ok {
			t.Error("unexpected result:", code, "expected:", codes.NotFound)
		}

		_, err = service.FreeTask(ctx, &services.TaskActionMsg{TaskID: "ABCDE"})
		if err == nil {
			t.Error("unexpected result:", err)
		}

		if ok, code := matchExpectedStatusFromError(err, codes.NotFound); !ok {
			t.Error("unexpected result:", code, "expected:", codes.NotFound)
		}

}

func TestHoldFree_Errors_InvalidArgument(t *testing.T) {

		service := createTaskService(t)
		ctx := context.Background()

		_, err := service.HoldTask(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
		if err == nil {
			t.Error("unexpected result:", err)
		}

		if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
			t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
		}

		_, err = service.FreeTask(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
		if err == nil {
			t.Error("unexpected result:", err)
		}

		if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
			t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
		}
	}

func TestHoldFree(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_05"})
	if err != nil {
		t.Error("unexpected result")
	}

	msg := strings.Split(r.Message, ":")

	r, err = service.HoldTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	_, err = service.HoldTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.FailedPrecondition); !ok {
		t.Error("unexpected result:", code, "expected:", codes.FailedPrecondition)
	}

	r, err = service.FreeTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	_, err = service.FreeTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.FailedPrecondition); !ok {
		t.Error("unexpected result:", code, "expected:", codes.FailedPrecondition)
	}

	// call method directly to skip setting the name of the user in a middleware handler
	_, err = tsrvs.FreeTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Unauthenticated); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Unauthenticated)
	}

}

	func TestTaskLog_Errors(t *testing.T) {
		service := createTaskService(t)
		ctx := context.Background()

		_, err := service.TaskLog(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
		if err == nil {
			t.Error("unexpected result:", err)
		}

		if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
			t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
		}
	}

	func TestTaskOutput_Errors(t *testing.T) {
		service := createTaskService(t)
		ctx := context.Background()

		_, err := service.TaskOutput(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
		if err == nil {
			t.Error("unexpected result:", err)
		}

		if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
			t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
		}
	}

func TestTaskLogTaskOutput(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_05"})
	if err != nil {
		t.Error("unexpected result")
	}

	msg := strings.Split(r.Message, ":")

	_, err = service.TaskLog(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}

	_, err = service.TaskOutput(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}

}

func TestTaskDetail_Errors(t *testing.T) {

		service := createTaskService(t)
		ctx := context.Background()

		_, err := service.TaskDetail(ctx, &services.TaskActionMsg{TaskID: "ABCD1"})

		if err == nil {
			t.Error("unexpected result:", err)
		}

		if ok, code := matchExpectedStatusFromError(err, codes.NotFound); !ok {
			t.Error("unexpected result:", code, "expected:", codes.NotFound)
		}
	}

func TestTaskDetail(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()
	msg := &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_05"}

	r, err := service.ForceTask(ctx, msg)
	if err != nil {
		t.Error("unexpected result")
	}

	rmsg := strings.Split(r.Message, ":")

	d, err := service.TaskDetail(ctx, &services.TaskActionMsg{TaskID: rmsg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if d.Result.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

}
*/
func (suite *activeTaskTestSuite) TestTaskGetAllowedAction() {

	tdata := map[string]auth.UserAction{
		"ListTasks":   auth.ActionBrowse,
		"TaskDetail":  auth.ActionBrowse,
		"TaskLog":     auth.ActionBrowse,
		"TaskOutput":  auth.ActionBrowse,
		"OrderTask":   auth.ActionOrder,
		"ForceTask":   auth.ActionForce,
		"RerunTask":   auth.ActionRestart,
		"EnforceTask": auth.ActionRestart,
		"HoldTask":    auth.ActionHold,
		"FreeTask":    auth.ActionFree,
		"SetToOk":     auth.ActionSetToOK,
		"ConfirmTask": auth.ActionConfirm,
	}

	for k, v := range tdata {

		act := suite.ovsService.GetAllowedAction(k)
		if act != v {
			suite.Fail(fmt.Sprintf("unexpected result:%v expected:%v", act, v))
		}

	}

}
