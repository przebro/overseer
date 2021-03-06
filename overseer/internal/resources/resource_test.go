package resources

import (
	"encoding/json"
	"os"
	"overseer/common/logger"
	"overseer/common/types/date"
	"overseer/datastore"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/taskdef"
	"testing"
)

var manager ResourceManager

type mockDispacher struct {
}

func (m *mockDispacher) PushEvent(sender events.EventReceiver, route events.RouteName, msg events.DispatchedMessage) error {
	return nil
}
func (m *mockDispacher) Subscribe(route events.RouteName, participant events.EventParticipant) {

}
func (m *mockDispacher) Unsubscribe(route events.RouteName, participant events.EventParticipant) {

}

var dispatcher mockDispacher = mockDispacher{}

var resConfig config.ResourcesConfigurartion = config.ResourcesConfigurartion{
	TicketSource: config.ResourceEntry{Sync: 1, Collection: "resources"},
	FlagSource:   config.ResourceEntry{Sync: 1, Collection: "resources"},
}
var storeConfig config.StoreProviderConfiguration = config.StoreProviderConfiguration{
	Store: []config.StoreConfiguration{
		{ID: "teststore", ConnectionString: "local;/../../../data/tests?synctime=1"},
	},
	Collections: []config.CollectionConfiguration{
		{Name: "resources", StoreID: "teststore"},
	},
}

var provider *datastore.Provider

var resources = map[string]interface{}{
	"tickets": TicketsResourceModel{
		ID: "tickets",
		Tickets: []TicketResource{
			{Name: "TEST_TICKET", Odate: ""},
		},
	},
	"flags": FlagsResourceModel{
		ID: "flags",
		Flags: []FlagResource{
			{Name: "TEST_FLAG", Count: 1, Policy: FlagPolicyShared},
		},
	},
}

func init() {

	var err error
	log := logger.NewTestLogger()
	f, _ := os.Create("../../../data/tests/resources.json")
	data, _ := json.Marshal(resources)

	f.Write(data)
	f.Close()

	provider, err = datastore.NewDataProvider(storeConfig, log)

	if err != nil {
		panic("fatal error, unable to load store")
	}

}

func TestNewManager(t *testing.T) {

	var err error
	tlog := logger.NewTestLogger()
	manager, err = NewManager(&dispatcher, tlog, resConfig, provider)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	invalidConfig := config.ResourcesConfigurartion{
		TicketSource: config.ResourceEntry{Sync: 1, Collection: "resources"},
		FlagSource:   config.ResourceEntry{Sync: 1, Collection: "invalid"},
	}

	_, err = NewManager(&dispatcher, tlog, invalidConfig, provider)

	if err == nil {
		t.Error("unexpected error")
	}
	invalidConfig.TicketSource.Collection = "invalid"

	_, err = NewManager(&dispatcher, tlog, invalidConfig, provider)

	if err == nil {
		t.Error("unexpected error")
	}

}

func TestStartShutdown(t *testing.T) {

	tlog := logger.NewTestLogger()
	m, err := NewManager(&dispatcher, tlog, resConfig, provider)
	if err != nil {
		t.Error("unexpected error:", err)
	}
	m.Start()
	m.Shutdown()

}

func TestAddCondition(t *testing.T) {

	_, err := manager.Add("RCOND_A_01", "")

	if err != nil {
		t.Error(err)
	}

	_, err = manager.Add("RCOND_A_01", "0909")

	if err != nil {
		t.Error(err)
	}

	_, err = manager.Add("RCOND_A_01", "0909")

	if err == nil {
		t.Error(err)
	}
	manager.Delete("RCOND_A_01", "0909")
	manager.Delete("RCOND_A_01", "")

}
func TestDeleteCondition(t *testing.T) {

	//Check for condition tha already doesn't exists
	_, err := manager.Delete("COND_D_01", "")
	if err == nil {
		t.Error(err)
	}

	_, err = manager.Add("COND_D_01", "")
	if err != nil {
		t.Fatal(err)
	}

	_, err = manager.Add("COND_D_01", "0909")
	if err != nil {
		t.Fatal(err)
	}

	//Remove condition
	_, err = manager.Delete("COND_D_01", "0909")
	if err != nil {
		t.Error(err)
	}

	// Remove condition that does not exists
	_, err = manager.Delete("COND_D_01", "0909")
	if err == nil {
		t.Error(err)
	}

	//Remove condition
	_, err = manager.Delete("COND_D_01", "")
	if err != nil {
		t.Error(err)
	}
}
func TestCheckCondition(t *testing.T) {

	res, err := manager.Add("COND_C_01", "")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)

	res = manager.Check("COND_C_01", "")

	if res == false {
		t.Error(err)
	}

	res = manager.Check("COND_C_01", "0909")

	if res == true {
		t.Error(err)
	}

	manager.Delete("COND_C_01", "")
}
func TestListTickets(t *testing.T) {

	manager.Add("TESTL_01", "20201013")
	manager.Add("TESTL_01", "")
	manager.Add("TESTL_01", "20201009")
	manager.Add("TESTL_02", "20201009")

	res := manager.ListTickets("TESTL_*", "")
	if len(res) != 4 {
		t.Error("invalid number of conditions", len(res))
	}

	res = manager.ListTickets("TESTL_01", "")
	if len(res) != 3 {
		t.Error("invalid number of conditions", len(res))
	}

	res = manager.ListTickets("TESTL_01", "20*")
	if len(res) != 2 {
		t.Error("invalid number of conditions:", len(res))
	}

	res = manager.ListTickets("TESTL_01", "202010*")
	if len(res) != 2 {
		t.Error("invalid number of conditions", len(res))
	}

	manager.Delete("TESTL_01", "20200909")
	manager.Delete("TESTL_01", "")
	manager.Delete("TESTL_01", "20201009")
	manager.Delete("TESTL_02", "20201009")
}

func TestFlags(t *testing.T) {

	testman := &resourceManager{}
	testman.dispatcher = &dispatcher
	testman.log = logger.NewTestLogger()
	trw, err := newTicketReadWriter("resources", "tickets", provider)
	if err != nil {
		t.Fatal(err)
	}
	testman.tstore, err = newStore(trw, 3600)

	if err != nil {
		t.Fatal(err)
	}

	frw, err := newFlagReadWriter("resources", "flags", provider)
	if err != nil {
		t.Fatal(err)
	}

	testman.fstore, err = newStore(frw, 3600)

	if err != nil {
		t.Fatal(err)
	}

	//Sets a flag to exclusive
	result, err := testman.Set("FLAG_01", FlagPolicyExclusive)
	if err != nil {
		t.Error("Set flag failed#1:", err)
	}
	if result != true {
		t.Error("Set flag failed#2:", err)
	}
	//Setting a flag to exclusive if already exists one is invalid.
	result, err = testman.Set("FLAG_01", FlagPolicyExclusive)
	if err == nil {
		t.Error("Set flag failed#3:", err)
	}
	if result != false {
		t.Error("Set flag failed#4:", err)
	}
	//Setting a flag to shared if there is a flag with exclusive is invalid
	result, _ = testman.Set("FLAG_01", FlagPolicyShared)
	if result == true {
		t.Error("Set flag failed#5")
	}

	list := testman.ListFlags("FLAG*")
	if len(list) != 1 {
		t.Error("List flag failed#1, expected len:1, result", len(list))
	}

	result, err = testman.Unset("FLAG_01")

	if err != nil {
		t.Error(err)
	}

	if result != true {
		t.Error("Unset flag failed#1")
	}

	list = testman.ListFlags("FLAG*")
	if len(list) != 0 {
		t.Error("List flag failed#2")
	}

	_, err = testman.Unset("FLAG_01")
	if err == nil {
		t.Error("Unset flag failed#2", err)
	}

	//Sets flag to shared
	result, err = testman.Set("FLAG_02", FlagPolicyShared)

	if err != nil {
		t.Error(err)
	}

	if result != true {
		t.Error("Set flag failed#6")
	}

	//It is possible to set a flag to shared if exists already
	result, err = testman.Set("FLAG_02", FlagPolicyShared)

	if err != nil {
		t.Error(err)
	}

	if result != true {
		t.Error("Set flag failed#7")
	}

	//If there is a shared flag, exclusive is not possible
	result, err = testman.Set("FLAG_02", FlagPolicyExclusive)

	if err == nil {
		t.Error("unexpected result")
	}

	if result == true {
		t.Error("Set flag failed#8")
	}

	list = testman.ListFlags("FLAG")
	lflag, _ := testman.fstore.Get("FLAG_02")
	lresource := lflag.(FlagResource)

	if len(list) != 1 && lresource.Count != 2 {
		t.Error("List flag failed#3")
	}

	_, err = testman.Unset("FLAG_02")
	if err != nil {
		t.Error("Unset flag failed#3", err)
	}

	list = testman.ListFlags("FLAG*")

	lflag, _ = testman.fstore.Get("FLAG_02")
	lresource = lflag.(FlagResource)

	if len(list) != 1 && lresource.Count != 1 {
		t.Error("List flag failed#4")
	}

	_, err = testman.Unset("FLAG_02")
	if err != nil {
		t.Error("Unset flag failed#4", err)
	}

	result, err = testman.Set("FLAG_03", FlagPolicyShared)

	if err != nil {
		t.Error(err)
	}

	if result != true {
		t.Error("Set flag failed#8")
	}

	result, err = testman.DestroyFlag("FLAG_03")

	if err != nil {
		t.Error(err)
	}

	if result != true {
		t.Error("destroy flag failed#1")
	}

	result, err = testman.DestroyFlag("FLAG_03")

	if err == nil {
		t.Error("unexpected result")
	}

	if result != false {
		t.Error("Set flag failed#2")
	}

}
func TestDispatch(t *testing.T) {

	var e error
	testman := &resourceManager{}
	testman.dispatcher = &dispatcher

	testman.log = logger.NewTestLogger()

	trw, err := newTicketReadWriter("resources", "tickets", provider)
	if err != nil {
		t.Fatal(e)
	}

	testman.tstore, e = newStore(trw, 3600)

	if e != nil {
		t.Fatal(e)
	}

	frw, err := newFlagReadWriter("resources", "flags", provider)
	if err != nil {
		t.Fatal(e)
	}

	testman.fstore, e = newStore(frw, 3600)

	if e != nil {
		t.Fatal(e)
	}

	receiver := events.NewTicketCheckReceiver()

	go testman.Process(receiver, "ROUTE_NOT_EXISTS", events.NewMsg(""))

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Error("Invalid route name")
	}

	msg := events.RouteTicketCheckMsgFormat{}
	msg.Tickets = []struct {
		Name      string
		Odate     string
		Fulfilled bool
	}{
		{Name: "CCCCC", Odate: "20200202", Fulfilled: false},
	}

	go testman.Process(receiver, events.RouteTicketCheck, events.NewMsg(msg))
	result, err := receiver.WaitForResult()
	if err != nil {
		t.Error("expected response")
	}
	if len(result.Tickets) != 1 {
		t.Error("unexpected data")
	}

	go testman.Process(receiver, events.RouteTicketCheck, events.NewMsg("msg"))
	result, err = receiver.WaitForResult()
	if err == nil {
		t.Error("expected response")
	}
	if err != events.ErrUnrecognizedMsgFormat {
		t.Error("unexpected data")
	}

	ticketMsg := events.RouteTicketInMsgFormat{Tickets: make([]struct {
		Name   string
		Odate  date.Odate
		Action taskdef.OutAction
	}, 2)}

	ticketMsg.Tickets[0].Action = taskdef.OutActionRemove
	ticketMsg.Tickets[0].Name = "ADDTEST01"
	ticketMsg.Tickets[0].Odate = "ODAT"

	ticketMsg.Tickets[1].Action = taskdef.OutActionAdd
	ticketMsg.Tickets[1].Name = "ADDTEST02"
	ticketMsg.Tickets[1].Odate = "ODAT"

	go testman.Process(nil, events.RouteTicketIn, events.NewMsg(ticketMsg))

	go testman.Process(receiver, events.RouteTicketIn, events.NewMsg(""))

	_, err = receiver.WaitForResult()
	if err == nil {
		t.Error("unexpected data expected error actual nil")
	}

}

func TestCreateReadWriter(t *testing.T) {

	_, err := newTicketReadWriter("invalid", "tickets", provider)
	if err == nil {
		t.Error("unexpected result")
	}

	_, err = newTicketReadWriter("resources", "tickets", provider)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	_, err = newFlagReadWriter("invalid", "tickets", provider)
	if err == nil {
		t.Error("unexpected result")
	}

	_, err = newFlagReadWriter("resources", "flags", provider)

	if err != nil {
		t.Error("unexpected result:", err)
	}
}
