package events

import (
	"errors"
	"goscheduler/common/logger"
	"goscheduler/overseer/internal/unique"
	"testing"
	"time"
)

var log logger.AppLogger = logger.NewLogger("./logs", 2)

type mockParticipant struct {
	rout RouteName
	msg  DispatchedMessage
}

func (m *mockParticipant) Process(p EventReceiver, routename RouteName, msg DispatchedMessage) {
	m.msg = msg
	m.rout = routename

}

func TestMessage(t *testing.T) {

	msg := NewCorrelatedMsg(unique.None(), "", []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})

	if msg.CorrelationID() != unique.None() {
		t.Error("correlation id not empty")
	}

	if msg.ResponseTo() != "" {
		t.Error("responseTo not empty")
	}
}

func TestRoute(t *testing.T) {

	mock := &mockParticipant{}
	r := &messgeRoute{participants: make([]EventParticipant, 0), routename: "name"}
	now := time.Now()

	r.AddParticipant(mock)
	if len(r.participants) != 1 {
		t.Error("Add participant to route")
	}
	msg := NewMsg([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	r.PushMessage(nil, msg)
	if mock.rout != "name" {
		t.Error("invalid routename")
	}

	if len(msg.MsgID()) == 0 {
		t.Error("invalid message id")

	}

	result, isOk := mock.msg.Message().([]byte)
	if !isOk {
		t.Fatal("type assertion failed")
	}

	if mock.msg.ResponseTo() != "noroute" && len(result) != 10 {
		t.Error("invalid message from route")
	}

	post := time.Now()

	if msg.Created().Before(post) != true && now.Before(msg.Created()) {
		t.Error("invalid message creation time")
	}

}

func TestDispatcher(t *testing.T) {

	disp := NewDispatcher()

	if disp == nil {
		t.Error("unable to create dispatcher")
	}
	mock1 := &mockParticipant{}
	mock2 := &mockParticipant{}
	sdsp := eventDipspatcher{log: log, msgRoutes: make(map[RouteName]MessageRoute)}

	if len(sdsp.msgRoutes) != 0 {
		t.Error("Invalid msgroutes size")

	}
	sdsp.Subscribe("FAKE_ROUTE", mock1)
	if len(sdsp.msgRoutes) != 1 {
		t.Error("Invalid msgroutes size")

	}
	sdsp.Subscribe("FAKE_ROUTE", mock2)
	if len(sdsp.msgRoutes) != 1 {
		t.Error("Invalid msgroutes size")
	}

	fakeroute := &messgeRoute{routename: "FAKE_ROUTE2", participants: make([]EventParticipant, 0)}
	sdsp.msgRoutes["FAKE_ROUTE2"] = fakeroute

	sdsp.Subscribe("FAKE_ROUTE2", mock1)
	if len(fakeroute.participants) != 1 {
		t.Error("Invalid msgroutes size ")
	}
	sdsp.Subscribe("FAKE_ROUTE2", mock2)
	if len(fakeroute.participants) != 2 {
		t.Error("Invalid msgroutes size ")
	}

	sdsp.Unsubscribe("FAKE_ROUTE2", mock1)
	if len(fakeroute.participants) != 1 {
		t.Error("Invalid msgroutes size ")
	}

	sdsp.Unsubscribe("FAKE_ROUTE2", mock2)
	if len(fakeroute.participants) != 0 {
		t.Error("Invalid msgroutes size ")
	}

	err := sdsp.PushEvent(nil, "ROUTE_DOES_NOT_EXISTS", NewMsg("Some data"))
	if err == nil {
		t.Error("Expected result : Route not defined")
	}
	err = sdsp.PushEvent(nil, "FAKE_ROUTE", NewMsg("Some data"))
	if err != nil {
		t.Error("Unexpected result", err)
	}

	time.Sleep(100 * time.Millisecond)

	val, ok := mock1.msg.Message().(string)
	if !ok {
		t.FailNow()
	}
	if val != "Some data" {
		t.Error("Expected message != Some data")
	}

}
func TestEvent(t *testing.T) {

	cond := NewTicketCheckReceiver()
	go func() {
		time.Sleep(100 * time.Millisecond)
		cond.Done(errors.New("error"))
	}()
	_, err := cond.WaitForResult()
	if err == nil {
		t.Error("expected error message#1")
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		cond.Done("dfdsfsdfds")
	}()
	_, err = cond.WaitForResult()
	if err == nil {
		t.Error("expected error message#2")
	}
	go func() {
		time.Sleep(100 * time.Millisecond)
		cd := make([]struct {
			Name      string
			Odate     string
			Fulfilled bool
		}, 1)
		ResponseToReceiver(cond, RouteTicketCheckMsgFormat{Tickets: cd})

	}()

	result, err := cond.WaitForResult()
	if err != nil {
		t.Error("expected RouteTicketCheckMsgFormat")
	}
	if len(result.Tickets) != 1 {
		t.Error("RouteTicketCheckMsgFormat expected tickets length = 1")

	}

	task := NewActiveTaskReceiver()
	go func() {
		time.Sleep(100 * time.Millisecond)
		task.Done(errors.New("error"))
	}()
	_, err = task.WaitForResult()
	if err == nil {
		t.Error("expected error message#1")
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		task.Done("dfdsfsdfds")
	}()
	_, err = task.WaitForResult()
	if err == nil {
		t.Error("expected error message#2")
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		tinfo := make([]TaskInfoResultMsg, 1)
		ResponseToReceiver(task, RouteTaskActionResponseFormat{Data: tinfo})

	}()

	tresult, err := task.WaitForResult()
	if err != nil {
		t.Error("expected RouteTicketCheckMsgFormat")
	}
	if len(tresult.Data) != 1 {
		t.Error("RouteTicketCheckMsgFormat expected condition length = 1")

	}

}
