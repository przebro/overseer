package journal

import (
	"testing"
	"time"

	"github.com/przebro/databazaar/store"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/datastore/mock"
	"github.com/przebro/overseer/overseer/config"
	"github.com/stretchr/testify/suite"
)

type JournalTestSuite struct {
	suite.Suite
	journal *TaskLogJournal
	store   *mock.MockStore
}

func TestJournalTestSuite(t *testing.T) {
	suite.Run(t, new(JournalTestSuite))
}

func (s *JournalTestSuite) SetupSuite() {

	s.store = &mock.MockStore{}

	collection := &mock.MockCollection{}

	s.store.On("Collection", "journal").Return(collection, nil)

	fn := func(opt store.ConnectionOptions) (store.DataStore, error) {

		return s.store, nil
	}

	store.RegisterStoreFactory("mock", fn)

	sconfig := config.StoreConfiguration{ID: "teststore", ConnectionString: "mock;;data/tests"}

	provider, err := datastore.NewDataProvider(sconfig)
	s.Nil(err)
	j, err := NewTaskJournal(config.JournalConfiguration{SyncTime: 3600}, provider)
	s.Nil(err)
	s.journal = j
}

func (s *JournalTestSuite) TestReadLog_NoLog_Succesfull() {

	id := unique.TaskOrderID("12345")
	entries := s.journal.ReadLog(id)
	s.Equal(0, len(entries))
}

func (s *JournalTestSuite) TestReadLog_Succesful() {

	s.journal.store[unique.TaskOrderID("ABCDE")] = mLogModel{
		TaskID: "ABCDE", Entries: []LogEntry{
			{ExecutionID: "ABCDEF", Time: time.Now(), Message: "message"},
		}, Tstamp: time.Now(),
	}
	id := unique.TaskOrderID("ABCDE")
	entries := s.journal.ReadLog(id)
	s.Equal(1, len(entries))
}

func (s *JournalTestSuite) TestWriteLog_Succesful() {

	s.journal.store[unique.TaskOrderID("55555")] = mLogModel{
		TaskID: "ABCDE", Entries: []LogEntry{
			{ExecutionID: "ABCDEF", Time: time.Now(), Message: "message"},
		}, Tstamp: time.Now(),
	}
	s.journal.WriteLog(unique.TaskOrderID("55555"), LogEntry{ExecutionID: "ABCDEF", Time: time.Now(), Message: "message"})
	s.Equal(2, len(s.journal.store[unique.TaskOrderID("55555")].Entries))
}

func (s *JournalTestSuite) TestPushJournalMessage_Succesful() {

	s.journal.store[unique.TaskOrderID("66666")] = mLogModel{
		TaskID: "ABCDE", Entries: []LogEntry{
			{ExecutionID: "ABCDEF", Time: time.Now(), Message: "message"},
		}, Tstamp: time.Now(),
	}
	s.journal.PushJournalMessage(unique.TaskOrderID("66666"), "ABCDE1", time.Now(), "message")
	s.Equal(2, len(s.journal.store[unique.TaskOrderID("66666")].Entries))
}
