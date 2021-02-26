package journal

import (
	"os"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/unique"
	"testing"
	"time"
)

type mockDisp struct {
}

var storeConfig config.StoreProviderConfiguration = config.StoreProviderConfiguration{
	Store: []config.StoreConfiguration{
		{ID: "teststore", ConnectionString: "local;/../../../data/tests"},
	},
	Collections: []config.CollectionConfiguration{
		{Name: "journal", StoreID: "teststore"},
	},
}

var conf = config.JournalConfiguration{LogCollection: "journal", SyncTime: 0}
var provider *datastore.Provider
var jrnal TaskJournal

func init() {
	logger.NewTestLogger()
	f2, _ := os.Create("../../../data/tests/journal.json")
	f2.Write([]byte(`{}`))
	f2.Close()

	provider, _ = datastore.NewDataProvider(storeConfig)
	jrnal, _ = NewTaskJournal(conf, nil, provider)
}

func TestNewTaskJournal(t *testing.T) {
	var err error
	cfg := config.JournalConfiguration{LogCollection: "_invalid_collection", SyncTime: 0}

	_, err = NewTaskJournal(cfg, nil, provider)
	if err == nil {
		t.Error("unexpected result:", err)
	}
}

func TestReadWriteLog(t *testing.T) {

	id := unique.TaskOrderID("12345")
	tJrnal := jrnal.(*taskLogJournal)
	tJrnal.store[id] = mLogModel{
		TaskID: "12345", Entries: []LogEntry{}, Tstamp: time.Now(),
	}

	jrnal.WriteLog(id, LogEntry{ExecutionID: "ABCDEF", Time: time.Now(), Message: "message"})

	if len(tJrnal.store[id].Entries) != 1 {
		t.Error("unexpected result:")
	}

	entries := jrnal.ReadLog(id)
	if len(entries) != 1 {
		t.Error("unexpected result:")
	}

	id2 := unique.TaskOrderID("55555")
	jrnal.WriteLog(id2, LogEntry{ExecutionID: "ABCDEF", Time: time.Now(), Message: "message"})
	if len(tJrnal.store) != 2 {
		t.Error("unexpected result:")
	}

	id3 := unique.TaskOrderID("66666")
	entries = jrnal.ReadLog(id3)
	if len(tJrnal.store) != 2 {
		t.Error("unexpected result:", len(tJrnal.store))
	}

	if len(entries) != 0 {
		t.Error("unexpected result:", len(entries))
	}

}

func TestProcess(t *testing.T) {

	id := unique.TaskOrderID("12345")
	tJrnal := jrnal.(*taskLogJournal)

	tJrnal.Process(nil, events.RoutTaskJournal,
		events.NewMsg(events.RouteJournalMsg{
			OrderID:     id,
			ExecutionID: "ABCDEF",
			Time:        time.Now(),
			Msg:         "message",
		},
		),
	)

	entries := jrnal.ReadLog(id)
	if len(entries) != 2 {
		t.Error("unexpected result:", len(entries))
	}

}
