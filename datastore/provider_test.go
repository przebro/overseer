package datastore

import (
	"context"
	"errors"
	"testing"

	"github.com/przebro/databazaar/store"
	mockstore "github.com/przebro/overseer/datastore/mock"
	"github.com/przebro/overseer/overseer/config"
	"github.com/stretchr/testify/suite"
)

func init() {

}

type DatastoreTestSuite struct {
	suite.Suite
	store *mockstore.MockStore
}

func TestDatastoreTestSuite(t *testing.T) {
	suite.Run(t, new(DatastoreTestSuite))
}

func (suite *DatastoreTestSuite) SetupSuite() {

	suite.store = &mockstore.MockStore{}

	suite.store.On("Collection", "test1").Return(&mockstore.MockCollection{}, nil)
	suite.store.On("Collection", "noname").Return(nil, errors.New("no collection"))

	fn := func(opt store.ConnectionOptions) (store.DataStore, error) {

		return suite.store, nil
	}

	store.RegisterStoreFactory("mock", fn)

}

func (suite *DatastoreTestSuite) TestProvider_Successful() {

	conf := config.StoreConfiguration{ID: "test1", ConnectionString: "mock;;/../data/tests?synctime=1"}
	_, err := NewDataProvider(conf)
	suite.NoError(err)
}
func (suite *DatastoreTestSuite) TestProvider_NoDriverError() {
	conf := config.StoreConfiguration{ID: "test1", ConnectionString: "local;;/../data/tests?synctime=1"}
	_, err := NewDataProvider(conf)
	suite.NotNil(err)
}

func (suite *DatastoreTestSuite) TestGetCollection_Successful() {

	conf := config.StoreConfiguration{ID: "test1", ConnectionString: "mock;;data/tests?synctime=1"}
	p, err := NewDataProvider(conf)
	suite.Nil(err)
	_, err = p.GetCollection(context.Background(), "test1")
	suite.Nil(err)

}

func (suite *DatastoreTestSuite) TestGetCollection_NoCollection() {

	conf := config.StoreConfiguration{ID: "test1", ConnectionString: "mock;;data/tests?synctime=1"}
	p, err := NewDataProvider(conf)
	suite.Nil(err)
	_, err = p.GetCollection(context.Background(), "noname")
	suite.NotNil(err)
}

func (suite *DatastoreTestSuite) TestGetCollectionDirectory_Successful() {

	conf := config.StoreConfiguration{ID: "test1", ConnectionString: "mock;;data/tests?synctime=1"}
	p, err := NewDataProvider(conf)
	suite.Nil(err)
	suite.Equal("data/tests", p.Directory())

}
