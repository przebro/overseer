package services

import (
	"context"
	"io"
	"net"
	"overseer/common/logger"
	"overseer/overseer/auth"
	"overseer/proto/services"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

var rescl services.ResourceServiceClient
var rsrvc *ovsResourceService

func createResourceClient(t *testing.T) services.ResourceServiceClient {

	if rescl != nil {
		return rescl
	}

	listener := bufconn.Listen(1)
	mocksrv := &mockBuffconnServer{grpcServer: grpc.NewServer(buildUnaryChain(), buildStreamChain())}

	logger.NewTestLogger()
	var err error

	if err != nil {
		t.Fatal("unable to create connection", err)
	}

	resservice := NewResourceService(resmanager)
	rsrvc = resservice.(*ovsResourceService)

	services.RegisterResourceServiceServer(mocksrv.grpcServer, resservice)

	dialer := func(ctx context.Context, s string) (net.Conn, error) {
		return listener.Dial()
	}

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithInsecure(), grpc.WithContextDialer(dialer))
	if err != nil {
		t.Fatal("unable to create connection", err)
	}

	rescl = services.NewResourceServiceClient(conn)
	go mocksrv.grpcServer.Serve(listener)

	return rescl

}

func TestAddTicket(t *testing.T) {
	client := createResourceClient(t)
	msg := &services.TicketActionMsg{}

	r, err := client.AddTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Odate = "20201120"

	r, err = client.AddTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Odate = "ABCDEF"
	msg.Name = "test"

	r, err = client.AddTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "ticket_with_very_long_name_that_exceeds_32_characters"
	msg.Odate = "20201115"

	r, err = client.AddTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "service_test_1"

	r, err = client.AddTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	msg.Name = "service_test_1"

	r, err = client.AddTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}
}

func TestDeleteTicket(t *testing.T) {
	client := createResourceClient(t)
	msg := &services.TicketActionMsg{Name: "service_test_2", Odate: ""}

	r, err := client.AddTicket(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}

	msg = &services.TicketActionMsg{Name: "service_test_3", Odate: "20201120"}

	r, err = client.AddTicket(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}

	msg.Name = "very_long_name_that_exceeds_32_characters"
	msg.Odate = ""

	r, err = client.DeleteTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "service_test_2"
	msg.Odate = "ABCDEDF"

	r, err = client.DeleteTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "service_test_2"
	msg.Odate = ""

	r, err = client.DeleteTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	msg.Name = "service_test_2"
	msg.Odate = ""

	r, err = client.DeleteTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "service_test_3"
	msg.Odate = "20201120"

	r, err = client.DeleteTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

}

func TestCheckTicket(t *testing.T) {

	client := createResourceClient(t)
	msg := &services.TicketActionMsg{Name: "service_test_4", Odate: ""}
	r, err := client.AddTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	msg = &services.TicketActionMsg{Name: "service_test_4", Odate: "20201124"}
	r, err = client.AddTicket(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}

	msg.Name = "very_long_name_that_exceeds_32_characters"
	msg.Odate = ""

	r, err = client.CheckTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "service_test_4"
	msg.Odate = "ABCDEF"

	r, err = client.CheckTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "service_test_4"
	msg.Odate = ""

	r, err = client.CheckTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	msg.Name = "service_test_4"
	msg.Odate = "20201124"

	r, err = client.CheckTicket(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}
}

func TestListTicket(t *testing.T) {

	client := createResourceClient(t)
	msg := &services.TicketActionMsg{Name: "service_test_4", Odate: ""}

	r, err := client.ListTickets(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	cnt := 0

	for {
		_, err := r.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Error(err)
			break
		}
		cnt++
	}

	if cnt == 0 {
		t.Error("unexpected result:", cnt)
	}

	msg = &services.TicketActionMsg{Name: "very_long_ticket_name_that_exceeds_32_characters", Odate: ""}

	r, err = client.ListTickets(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if _, err := r.Recv(); err == nil {
		t.Error("unexpected result")
	}

	msg = &services.TicketActionMsg{Name: "service_test_4", Odate: "123456789"}

	r, err = client.ListTickets(context.Background(), msg)

	if err != nil {
		t.Error(err)
	}

	if _, err := r.Recv(); err == nil {
		t.Error("unexpected result")
	}

}

func TestSetFlag(t *testing.T) {

	client := createResourceClient(t)
	msg := &services.FlagActionMsg{Name: "very_long_resource_name_that_exceeds_32_chracters"}

	r, err := client.SetFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "test_flag_2"
	msg.State = 2
	r, err = client.SetFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "test_flag_1"
	msg.State = 0
	r, err = client.SetFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	msg.Name = "test_flag_1"
	msg.State = 1
	//flag is already set to shared so, setting to exclusive is not allowed
	r, err = client.SetFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "test_flag_2"
	msg.State = 1
	r, err = client.SetFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}
}

func TestListFlags(t *testing.T) {

	client := createResourceClient(t)
	msg := &services.FlagActionMsg{Name: "very_long_resource_name_that_exceeds_32_chracters"}
	result, err := client.ListFlags(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if _, err := result.Recv(); err == nil {
		t.Error("unexpected result")
	}

	msg.Name = "test_flag*"
	result, err = client.ListFlags(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	_, err = result.Recv()

	if err != nil {
		t.Error("unexpected result:", err)
	}

}

func TestDestroyFlag(t *testing.T) {

	client := createResourceClient(t)
	msg := &services.FlagActionMsg{Name: "very_long_resource_name_that_exceeds_32_chracters"}
	r, err := client.DestroyFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Name = "test_flag_1"
	r, err = client.DestroyFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	msg.Name = "test_flag_1"
	r, err = client.DestroyFlag(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

}

func TestAllowedActions(t *testing.T) {

	createResourceClient(t)

	tdata := map[string]auth.UserAction{
		"AddTicket":    auth.ActionAddTicket,
		"DeleteTicket": auth.ActionRemoveTicket,
		"CheckTicket":  auth.ActionBrowse,
		"ListTickets":  auth.ActionBrowse,
		"SetFlag":      auth.ActionSetFlag,
		"DestroyFlag":  auth.ActionSetFlag,
		"ListFlags":    auth.ActionBrowse,
	}

	for k, v := range tdata {
		result := rsrvc.GetAllowedAction(k)
		if v != result {
			t.Error("unexpected result:", result, "expected:", v)
		}
	}
}
