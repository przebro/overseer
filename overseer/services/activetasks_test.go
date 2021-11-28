package services

import (
	"context"
	"io"
	"net"
	"overseer/common/logger"
	"overseer/common/types/date"
	"overseer/overseer/auth"
	"overseer/overseer/internal/journal"
	"overseer/overseer/internal/unique"
	"overseer/overseer/services/handlers"
	"overseer/overseer/services/middleware"
	"overseer/proto/services"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/test/bufconn"
)

type mockJournal struct {
}

var tsl services.TaskServiceClient
var jrnl *mockJournal = &mockJournal{}
var tsrvs *ovsActiveTaskService

func (m *mockJournal) WriteLog(id unique.TaskOrderID, entry journal.LogEntry) {}
func (m *mockJournal) ReadLog(id unique.TaskOrderID) []journal.LogEntry       { return []journal.LogEntry{} }
func (m *mockJournal) Start() error                                           { return nil }
func (m *mockJournal) Shutdown() error                                        { return nil }
func (m *mockJournal) Resume() error                                          { return nil }
func (m *mockJournal) Quiesce() error                                         { return nil }

func createTaskService(t *testing.T) services.TaskServiceClient {

	if tsl != nil {
		return tsl
	}

	tcv, err := NewTokenCreatorVerifier(authcfg)
	if err != nil {
		panic("")
	}
	authhandler, err := handlers.NewServiceAuthorizeHandler(authcfg, tcv, provider, logger.NewTestLogger())

	if err != nil {
		panic("")
	}

	middleware.RegisterHandler(authhandler)

	listener := bufconn.Listen(1)
	mocksrv := &mockBuffconnServer{grpcServer: grpc.NewServer(buildUnaryChain(), buildStreamChain())}

	srvc := NewTaskService(activeTaskManagerT, taskPoolT, jrnl, logger.NewTestLogger())
	tsrvs = srvc.(*ovsActiveTaskService)

	services.RegisterTaskServiceServer(mocksrv.grpcServer, srvc)

	dialer := func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer))
	if err != nil {
		t.Fatal("unable to create connection", err)
	}

	tsl = services.NewTaskServiceClient(conn)
	go mocksrv.grpcServer.Serve(listener)

	return tsl
}

func TestOrderTask_Errors(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "username", "<anonymous>")

	_, err := service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: "abcdef", TaskName: "dummy_01"})
	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "#$SDtest", Odate: string(date.CurrentOdate()), TaskName: "dummy_01"})
	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "%^&$%dummy_01"})
	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

}
func TestOrderTask(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "username", "<anonymous>")

	r, err := service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_04"})
	if err != nil {
		t.Error("unexpected result")
	}

	if r.Success != true || !strings.Contains(r.Message, "TaskID:") {
		t.Error("unexpected result:", r.Success, ";", r.Message)
	}

	// call method directly to skip setting the name of the user in a middleware handler
	_, err = tsrvs.OrderTask(context.Background(), &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_04"})

	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Unauthenticated); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Unauthenticated)
	}

}

func TestForceTask_Errors(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	_, err := service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: "abcdef", TaskName: "dummy_01"})
	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "#$SDtest", Odate: string(date.CurrentOdate()), TaskName: "dummy_01"})
	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "%^&$%dummy_01"})
	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

}
func TestForceTask(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_04"})
	if err != nil {
		t.Error("unexpected result")
	}

	if r.Success != true || !strings.Contains(r.Message, "TaskID:") {
		t.Error("unexpected result:", r.Success, ";", r.Message)
	}

	// call method directly to omit set of the name of the user in middleware handler
	_, err = tsrvs.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_04"})

	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Unauthenticated); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Unauthenticated)
	}

}

func TestListTask(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.ListTasks(ctx, &services.TaskFilterMsg{})
	if err != nil {
		t.Error("unexpected result")
	}

	cnt := 0
	for {
		_, err := r.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error("unexpected result")
		}
		cnt++
	}

	if cnt != 2 {
		t.Error("unexpected result:", cnt, "expected:", 2)
	}

}

func TestOrderGroup_Errors(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	_, err := service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "_4@42!5terfds", Odate: string(date.CurrentOdate())})
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "TEST", Odate: "ABCDEF"})
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "", Odate: string(date.CurrentOdate())})
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "ABCDED", Odate: ""})
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Internal); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Internal)
	}
}

func TestOrderGroup(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	_, err := service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "test", Odate: ""})
	if err != nil {
		t.Error("unexpected result:", err)
	}

	// call method directly to skip setting the name of the user in a middleware handler
	_, err = tsrvs.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "test", Odate: ""})

	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Unauthenticated); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Unauthenticated)
	}
}

func TestForceGroup_Errors(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	_, err := service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "_4@42!5terfds", Odate: string(date.CurrentOdate())})
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "TEST", Odate: "ABCDEF"})
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "", Odate: string(date.CurrentOdate())})
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "ABCDED", Odate: ""})
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Internal); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Internal)
	}

}

func TestForceGroup(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "test", Odate: ""})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	// call method directly to skip setting the name of the user in a middleware handler
	_, err = tsrvs.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "test", Odate: ""})

	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Unauthenticated); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Unauthenticated)
	}

}

func TestConfirmTask(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_04"})
	if err != nil {
		t.Error("unexpected result")
	}

	msg := strings.Split(r.Message, ":")

	_, err = service.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.InvalidArgument); !ok {
		t.Error("unexpected result:", code, "expected:", codes.InvalidArgument)
	}

	_, err = service.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: "ABCDE"})
	if err == nil {
		t.Error("unexpected result:", err)
	}

	if ok, code := matchExpectedStatusFromError(err, codes.NotFound); !ok {
		t.Error("unexpected result:", code, "expected:", codes.NotFound)
	}

	r, err = service.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	// call method directly to skip setting the name of the user in a middleware handler
	_, err = tsrvs.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err == nil {
		t.Error("unexpected result")
	}

	if ok, code := matchExpectedStatusFromError(err, codes.Unauthenticated); !ok {
		t.Error("unexpected result:", code, "expected:", codes.Unauthenticated)
	}
}

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

func TestTaskGetAllowedAction(t *testing.T) {

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

		act := tsrvs.GetAllowedAction(k)
		if act != v {
			t.Error("unexpected result:", act, "expected:", v)
		}

	}

}
