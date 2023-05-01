package mock

import (
	"context"
	"errors"
	"fmt"

	"github.com/przebro/databazaar/collection"
	"github.com/przebro/databazaar/result"
	"github.com/przebro/databazaar/selector"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateCollection(ctx context.Context, name string) (collection.DataCollection, error) {
	args := m.Called(name)
	return args.Get(0).(collection.DataCollection), args.Error(1)

}
func (m *MockStore) Collection(ctx context.Context, name string) (collection.DataCollection, error) {
	args := m.Called(name)
	if args.Get(0) == nil {
		return nil, errors.New("collection not found")
	}
	return args.Get(0).(collection.DataCollection), args.Error(1)

}
func (m *MockStore) Status(ctx context.Context) (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)

}
func (m *MockStore) Close(ctx context.Context) {
	m.Called()
}
func (m *MockStore) CollectionExists(ctx context.Context, name string) bool {
	args := m.Called()
	return args.Bool(0)
}

// MockCollection - implemenets ccollection.DataCollection interface
type MockCollection struct {
	data map[string]interface{}
}

func (m *MockCollection) Create(ctx context.Context, document interface{}) (*result.BazaarResult, error) {
	return nil, nil
}
func (m *MockCollection) Get(ctx context.Context, id string, result interface{}) error {
	return nil
}
func (m *MockCollection) Update(ctx context.Context, doc interface{}) error {
	return nil
}
func (m *MockCollection) Delete(ctx context.Context, id string) error {
	return nil
}
func (m *MockCollection) CreateMany(ctx context.Context, docs []interface{}) ([]result.BazaarResult, error) {
	return nil, nil
}
func (m *MockCollection) BulkUpdate(ctx context.Context, docs []interface{}) error {
	return nil
}
func (m *MockCollection) All(ctx context.Context) (collection.BazaarCursor, error) {
	return nil, nil
}
func (m *MockCollection) Count(ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *MockCollection) AsQuerable() (collection.QuerableCollection, error) {
	return nil, nil
}
func (m *MockCollection) Select(ctx context.Context, s selector.Expr, fld selector.Fields) (collection.BazaarCursor, error) {
	return nil, nil
}

func (m *MockCollection) Type() string {
	return "mock"
}

// MockCollectionProvider - implements datastore.CollectionProvider interface
type MockCollectionProvider struct {
	name       string
	collection collection.DataCollection
}

func NewMockCollectionProvider(name string, collection collection.DataCollection) *MockCollectionProvider {
	return &MockCollectionProvider{
		name:       name,
		collection: collection,
	}
}
func NewEmptyMockCollectionProvider(name string, collection collection.DataCollection) *MockCollectionProvider {
	return &MockCollectionProvider{
		name:       name,
		collection: &MockCollection{},
	}
}

func (m *MockCollectionProvider) GetCollection(name string) (collection.DataCollection, error) {
	if m.name != name {
		return nil, fmt.Errorf("%w;no mapping beetween collection name and store", errors.New("Store configuration error"))
	}
	return m.collection, nil
}
