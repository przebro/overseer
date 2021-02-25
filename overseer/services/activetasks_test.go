package services

import (
	"context"
	"io"
	"net"
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
	"google.golang.org/grpc/test/bufconn"
)

type mockJournal struct {
}

var tsl services.TaskServiceClient
var jrnl *mockJournal = &mockJournal{}
var tsrvs *ovsActiveTaskService

func (m *mockJournal) WriteLog(id unique.TaskOrderID, entry journal.LogEntry) {}
func (m *mockJournal) ReadLog(id unique.TaskOrderID) []journal.LogEntry       { return []journal.LogEntry{} }

func createTaskService(t *testing.T) services.TaskServiceClient {

	if tsl != nil {
		return tsl
	}

	tcv, err := NewTokenCreatorVerifier(authcfg)
	if err != nil {
		panic("")
	}
	authhandler, err := handlers.NewServiceAuthorizeHandler(authcfg, tcv, provider)

	if err != nil {
		panic("")
	}

	middleware.RegisterHandler(authhandler)

	listener := bufconn.Listen(1)
	mocksrv := &mockBuffconnServer{grpcServer: grpc.NewServer(buildUnaryChain(), buildStreamChain())}

	srvc := NewTaskService(activeTaskManagerT, taskPoolT, jrnl)
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

func TestOrderTask(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: "abcdef", TaskName: "dummy_01"})
	if err != nil {
		t.Error("unexpected result")
	}

	if r.Success != false || !strings.Contains(r.Message, "'odate'") {
		t.Error("unexpected result:", r.Success, ";", r.Message)
	}

	r, err = service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "#$SDtest", Odate: string(date.CurrentOdate()), TaskName: "dummy_01"})
	if err != nil {
		t.Error("unexpected result")
	}

	if r.Success != false || !strings.Contains(r.Message, "'Group'") {
		t.Error("unexpected result:", r.Success, ";", r.Message)
	}

	r, err = service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "%^&$%dummy_01"})
	if err != nil {
		t.Error("unexpected result")
	}

	if r.Success != false || !strings.Contains(r.Message, "'Name'") {
		t.Error("unexpected result:", r.Success, ";", r.Message)
	}

	r, err = service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_04"})
	if err != nil {
		t.Error("unexpected result")
	}

	if r.Success != true || !strings.Contains(r.Message, "TaskID:") {
		t.Error("unexpected result:", r.Success, ";", r.Message)
	}
}

func TestForceTask(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: "abcdef", TaskName: "dummy_01"})
	if err != nil {
		t.Error("unexpected result")
	}

	if r.Success != false || !strings.Contains(r.Message, "'odate'") {
		t.Error("unexpected result:", r.Success, ";", r.Message)
	}

	r, err = service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "#$SDtest", Odate: string(date.CurrentOdate()), TaskName: "dummy_01"})
	if err != nil {
		t.Error("unexpected result")
	}

	if r.Success != false || !strings.Contains(r.Message, "'Group'") {
		t.Error("unexpected result:", r.Success, ";", r.Message)
	}

	r, err = service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "%^&$%dummy_01"})
	if err != nil {
		t.Error("unexpected result")
	}

	if r.Success != false || !strings.Contains(r.Message, "'Name'") {
		t.Error("unexpected result:", r.Success, ";", r.Message)
	}

	r, err = service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_04"})
	if err != nil {
		t.Error("unexpected result")
	}

	if r.Success != true || !strings.Contains(r.Message, "TaskID:") {
		t.Error("unexpected result:", r.Success, ";", r.Message)
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

func TestOrderGroup(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "_4@42!5terfds", Odate: string(date.CurrentOdate())})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "TEST", Odate: "ABCDEF"})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "", Odate: string(date.CurrentOdate())})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "ABCDED", Odate: ""})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.OrderGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "test", Odate: ""})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}
}

func TestForceGroup(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "_4@42!5terfds", Odate: string(date.CurrentOdate())})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "TEST", Odate: "ABCDEF"})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "", Odate: string(date.CurrentOdate())})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "ABCDED", Odate: ""})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.ForceGroup(ctx, &services.TaskOrderGroupMsg{TaskGroup: "test", Odate: ""})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
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

	r, err = service.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
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

	r, err = service.HoldTask(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.FreeTask(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.HoldTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	r, err = service.HoldTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	r, err = service.FreeTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	r, err = service.FreeTask(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
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

	d, err := service.TaskLog(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if d.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	d, err = service.TaskLog(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if d.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	d, err = service.TaskOutput(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if d.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	d, err = service.TaskOutput(ctx, &services.TaskActionMsg{TaskID: msg[1]})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if d.Success != true || d.Message != "Not implemented" {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}
}

func TestTaskDetail(t *testing.T) {

	service := createTaskService(t)
	ctx := context.Background()

	r, err := service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: "test", Odate: string(date.CurrentOdate()), TaskName: "dummy_05"})
	if err != nil {
		t.Error("unexpected result")
	}

	msg := strings.Split(r.Message, ":")

	d, err := service.TaskDetail(ctx, &services.TaskActionMsg{TaskID: "ABCD"})
	if err != nil {
		t.Error("unexpected result:", err)
	}
	if d.Result.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	d, err = service.TaskDetail(ctx, &services.TaskActionMsg{TaskID: msg[1]})

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
