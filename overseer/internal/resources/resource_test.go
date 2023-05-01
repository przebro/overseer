package resources

import (
	"context"
	"sync"
	"testing"

	"github.com/przebro/databazaar/collection"
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/config"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

var manager ResourceManager

var resConfig config.ResourcesConfigurartion = config.ResourcesConfigurartion{
	Resources: config.ResourceEntry{Sync: 1},
}
var storeConfig config.StoreConfiguration = config.StoreConfiguration{ID: "teststore", ConnectionString: "local;/../../../data/tests?synctime=1"}

var provider *datastore.Provider

type mockCollectionProvider struct {
}

func (m *mockCollectionProvider) GetCollection(ctx context.Context, name string) (collection.DataCollection, error) {
	return nil, nil
}

func (m *mockCollectionProvider) Directory() string {
	return ""
}

type resourceTestSuite struct {
	suite.Suite
	manager ResourceManager
}

func TestResourcesTestSuite(t *testing.T) {
	suite.Run(t, new(resourceTestSuite))

}

func (s *resourceTestSuite) SetupSuite() {

	readWriter, _ := newResourceReadWriter("resources", &mockCollectionProvider{})

	s.manager = &ResourceManagerImpl{
		log:   log.Logger,
		rw:    readWriter,
		flock: sync.Mutex{},
	}
}
