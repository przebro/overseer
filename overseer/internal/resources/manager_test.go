package resources

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/config"
	"github.com/przebro/overseer/overseer/internal/events"
	"github.com/przebro/overseer/overseer/internal/taskdef"
)

var mlog = logger.NewTestLogger()
var mprovider *datastore.Provider
var mdispatcher mockDispacher = mockDispacher{}
var testManager ResourceManager

var manResConfig config.ResourcesConfigurartion = config.ResourcesConfigurartion{
	TicketSource: config.ResourceEntry{Sync: 1, Collection: "mresources"},
	FlagSource:   config.ResourceEntry{Sync: 1, Collection: "mresources"},
}
var manStoreConfig config.StoreProviderConfiguration = config.StoreProviderConfiguration{
	Store: []config.StoreConfiguration{
		{ID: "teststore", ConnectionString: "local;/../../../data/tests?synctime=1"},
	},
	Collections: []config.CollectionConfiguration{
		{Name: "mresources", StoreID: "teststore"},
	},
}

func init() {

	var err error
	f, _ := os.Create("../../../data/tests/mresources.json")
	data, _ := json.Marshal(resources)

	f.Write(data)
	f.Close()

	mprovider, err = datastore.NewDataProvider(manStoreConfig, mlog)

	if err != nil {
		panic("fatal error, unable to load store")
	}

	testManager, err = NewManager(&mdispatcher, mlog, manResConfig, mprovider)
	if err != nil {
		panic("failed to initialize manager")
	}

}

func TestNewManager_Errors(t *testing.T) {

	var err error

	invalidConfig := config.ResourcesConfigurartion{
		TicketSource: config.ResourceEntry{Sync: 1, Collection: "resources"},
		FlagSource:   config.ResourceEntry{Sync: 1, Collection: "invalid"},
	}

	_, err = NewManager(&mdispatcher, mlog, invalidConfig, mprovider)

	if err == nil {
		t.Error("unexpected error")
	}
	invalidConfig.TicketSource.Collection = "invalid"

	_, err = NewManager(&mdispatcher, mlog, invalidConfig, mprovider)

	if err == nil {
		t.Error("unexpected error")
	}

}

func TestDispatchInvalid_Route(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewTicketCheckReceiver()

	go testman.Process(receiver, "ROUTE_NOT_EXISTS", events.NewMsg(""))

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Error("Invalid route name")
	}
}

func TestDispatch_Route_RouteTicketCheck_Errors(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewTicketCheckReceiver()

	go testman.Process(receiver, events.RouteTicketCheck, events.NewMsg(""))

	_, err = receiver.WaitForResult()
	if err != events.ErrUnrecognizedMsgFormat {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}
}

func TestDispatch_Route_RouteTicketCheck(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewTicketCheckReceiver()
	msg := events.RouteTicketCheckMsgFormat{
		Tickets: []struct {
			Name      string
			Odate     string
			Fulfilled bool
		}{
			{Name: "TEST_CHECK_01"},
			{Name: "TEST_CHECK_02"},
		},
	}

	go testman.Process(receiver, events.RouteTicketCheck, events.NewMsg(msg))

	result, err := receiver.WaitForResult()
	if err != nil {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}

	if len(result.Tickets) != 2 {
		t.Error("unexpected result")
	}
}

func TestDispatch_Route_RouteTicketIn_Errors(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewTicketCheckReceiver()

	go testman.Process(receiver, events.RouteTicketIn, events.NewMsg(""))

	_, err = receiver.WaitForResult()
	if err != events.ErrUnrecognizedMsgFormat {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}
}

func TestDispatch_Route_RouteTicketIn(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewTicketInReceiver()
	msg := events.RouteTicketInMsgFormat{
		Tickets: []struct {
			Name   string
			Odate  date.Odate
			Action taskdef.OutAction
		}{
			{Name: "TEST_MSG_TICKET_01", Action: taskdef.OutActionAdd},
			{Name: "TEST_MSG_TICKET_02", Action: taskdef.OutActionAdd},
		},
	}

	go testman.Process(receiver, events.RouteTicketIn, events.NewMsg(msg))

	result, err := receiver.WaitForResult()
	if err != nil {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}

	if len(result.Tickets) != 2 {
		t.Error("unexpected result")
	}
}

func TestDispatch_Route_RouteAcquireFlag(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewFlagActionReceiver()
	msg := events.RouteFlagAcquireMsg{
		Flags: []events.FlagActionData{
			{Name: "TEST_FLAG_ACQ_01", Policy: 1},
		},
	}

	go testman.Process(receiver, events.RouteFlagAcquire, events.NewMsg(msg))

	result, err := receiver.WaitForResult()
	if err != nil {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}

	if !result.Success {
		t.Error("unexpected result")
	}
}

func TestDispatch_Route_RouteFlagAcquire_Errors(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewTicketCheckReceiver()

	go testman.Process(receiver, events.RouteFlagAcquire, events.NewMsg(""))

	_, err = receiver.WaitForResult()
	if err != events.ErrUnrecognizedMsgFormat {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}
}

func TestDispatch_Route_RouteFlagAcquire_ErrorExists(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewFlagActionReceiver()
	msg := events.RouteFlagAcquireMsg{
		Flags: []events.FlagActionData{
			{Name: "TEST_FLAG_ACQ_02", Policy: 1},
		},
	}

	go testman.Process(receiver, events.RouteFlagAcquire, events.NewMsg(msg))

	result, err := receiver.WaitForResult()
	if err != nil {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}

	if !result.Success {
		t.Error("unexpected result")
	}

	go testman.Process(receiver, events.RouteFlagAcquire, events.NewMsg(msg))

	if !result.Success {
		t.Error("unexpected result")
	}

}

func TestDispatch_Route_RouteFlagRelease_Errors(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewTicketCheckReceiver()

	go testman.Process(receiver, events.RouteFlagRelase, events.NewMsg(""))

	_, err = receiver.WaitForResult()
	if err != events.ErrUnrecognizedMsgFormat {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}
}

func TestDispatch_Route_RouteReleaseFlag(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewFlagActionReceiver()
	msg := events.RouteFlagAcquireMsg{
		Flags: []events.FlagActionData{
			{Name: "TEST_FLAG_ACQ_03", Policy: 1},
		},
	}

	go testman.Process(receiver, events.RouteFlagAcquire, events.NewMsg(msg))

	result, err := receiver.WaitForResult()
	if err != nil {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}

	if !result.Success {
		t.Error("unexpected result:", false, "expected:", true)
	}

	msg = events.RouteFlagAcquireMsg{
		Flags: []events.FlagActionData{
			{Name: "TEST_FLAG_ACQ_03"},
		},
	}

	go testman.Process(receiver, events.RouteFlagRelase, events.NewMsg(msg))

	result, err = receiver.WaitForResult()
	if err != nil {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}

	if !result.Success {
		t.Error("unexpected result:", false, "expected:", true)
	}
}

func TestDispatch_Route_RouteReleaseFlag_NotExists(t *testing.T) {

	var testman *resourceManager
	var err error
	ok := false

	if testman, ok = testManager.(*resourceManager); !ok {
		t.Error("failed to get testman")
	}

	receiver := events.NewFlagActionReceiver()
	msg := events.RouteFlagAcquireMsg{
		Flags: []events.FlagActionData{
			{Name: "TEST_FLAG_ACQ_04", Policy: 1},
		},
	}

	go testman.Process(receiver, events.RouteFlagAcquire, events.NewMsg(msg))

	result, err := receiver.WaitForResult()
	if err != nil {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}

	if !result.Success {
		t.Error("unexpected result:", false, "expected:", true)
	}

	msg = events.RouteFlagAcquireMsg{
		Flags: []events.FlagActionData{
			{Name: "TEST_FLAG_ACQ_04"},
			{Name: "TEST_FLAG_ACQ_44"},
		},
	}

	go testman.Process(receiver, events.RouteFlagRelase, events.NewMsg(msg))

	result, err = receiver.WaitForResult()
	if err != nil {
		t.Error("unexpected result:", err, "expected:", events.ErrUnrecognizedMsgFormat)
	}

	if result.Success {
		t.Error("unexpected result:", false, "expected:", true)
	}
}

func TestStartShutdown(t *testing.T) {

	var err error
	manager, err = NewManager(&mdispatcher, mlog, manResConfig, mprovider)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	manager.Start()
	manager.Shutdown()

}

func TestBuildExpr(t *testing.T) {

	result := buildExpr("")
	if result != `[\w\-]*|^$` {
		t.Error("unexpected result")
	}

	result = buildExpr("*")
	if result != `^[\w\-]*$` {
		t.Error("unexpected result")
	}

	result = buildExpr("?")
	if result != `^[\w\-]{1}$` {
		t.Error("unexpected result")
	}

	result = buildExpr("AB?X")
	if result != `^AB[\w\-]{1}X$` {
		t.Error("unexpected result:", result)
	}

	result = buildExpr("AB?*")
	if result != `^AB[\w\-]{1}[\w\-]*$` {
		t.Error("unexpected result:", result)
	}
}
func TestBuildDateExpr(t *testing.T) {

	result := buildDateExpr("")
	if result != `[\d]*|^$` {
		t.Error("unexpected result")
	}

	result = buildDateExpr("*")
	if result != `^[\d]*$` {
		t.Error("unexpected result")
	}

	result = buildDateExpr("?")
	if result != `^[\d]{1}$` {
		t.Error("unexpected result")
	}

	result = buildDateExpr("2021051?")
	if result != `^2021051[\d]{1}$` {
		t.Error("unexpected result")
	}

}
