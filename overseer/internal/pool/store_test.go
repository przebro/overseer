package pool

import (
	"context"
	"testing"

	"github.com/przebro/databazaar/collection"
	"github.com/przebro/databazaar/result"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/overseer/internal/pool/activetask"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockCollectionProvider struct {
	mock.Mock
}

func (m *mockCollectionProvider) GetCollection(ctx context.Context, name string) (collection.DataCollection, error) {
	args := m.Called(name)
	return args.Get(0).(collection.DataCollection), args.Error(1)
}

func (m *mockCollectionProvider) Directory() string {
	return ""
}

type mockDefinitionReader struct {
	mock.Mock
}

type mockDataCollection struct {
	mock.Mock
}

func (m *mockDataCollection) Create(ctx context.Context, document interface{}) (*result.BazaarResult, error) {
	args := m.Called(ctx, document)
	return args.Get(0).(*result.BazaarResult), args.Error(1)
}
func (m *mockDataCollection) Get(ctx context.Context, id string, result interface{}) error {
	return nil
}
func (m *mockDataCollection) Update(ctx context.Context, doc interface{}) error {
	return nil
}
func (m *mockDataCollection) Delete(ctx context.Context, id string) error {
	return nil
}
func (m *mockDataCollection) CreateMany(ctx context.Context, docs []interface{}) ([]result.BazaarResult, error) {
	return nil, nil
}
func (m *mockDataCollection) BulkUpdate(ctx context.Context, docs []interface{}) error {
	return nil
}
func (m *mockDataCollection) All(ctx context.Context) (collection.BazaarCursor, error) {
	return nil, nil
}
func (m *mockDataCollection) Count(ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *mockDataCollection) AsQuerable() (collection.QuerableCollection, error) {

	return nil, nil
}

func (m *mockDataCollection) Type() string {
	return ""
}

func (m *mockDefinitionReader) GetActiveDefinition(refID string) (*taskdef.TaskDefinition, error) {
	args := m.Called(refID)
	return args.Get(0).(*taskdef.TaskDefinition), args.Error(1)
}

type StoreTestSuite struct {
	suite.Suite
	provider       *mockCollectionProvider
	defReader      *mockDefinitionReader
	collectionName string
}

func TestStoreTestSuite(t *testing.T) {
	suite.Run(t, new(StoreTestSuite))
}

func (suite *StoreTestSuite) SetupSuite() {

	suite.provider = &mockCollectionProvider{}
	suite.collectionName = "tasks"
	suite.defReader = &mockDefinitionReader{}

	suite.provider.On("GetCollection", suite.collectionName).Return(&mockDataCollection{}, nil)
	suite.defReader.On("GetActiveDefinition", "12345").Return(nil, nil)

}

func (suite *StoreTestSuite) TestStore() {

	store, _ := NewStore(log.Logger, 3600, suite.provider, suite.defReader)
	if store == nil {
		suite.Fail("Store not created")
	}

	orderID := unique.TaskOrderID("12345")
	task := &activetask.TaskInstance{}

	store.add(orderID, task)
	store.add(unique.TaskOrderID("33333"), &activetask.TaskInstance{})
	store.add(unique.TaskOrderID("12346"), &activetask.TaskInstance{})

	suite.Equal(3, store.len(), "Invalid store size expected 3 actual :", store.len())

	_, exists := store.get(unique.TaskOrderID("54321"))
	suite.False(exists)

	store.remove(unique.TaskOrderID("12346"))
	store.remove(unique.TaskOrderID("12346"))

	suite.Equal(2, store.len())

}
