package datastore

import (
	"context"
	"os"
	"strings"

	"github.com/przebro/overseer/overseer/config"

	"github.com/przebro/databazaar/store"
	//required driver
	_ "github.com/przebro/badgerstore/store"
	_ "github.com/przebro/mongostore/store"

	"github.com/przebro/databazaar/collection"
)

type CollectionProvider interface {
	GetCollection(ctx context.Context, name string) (collection.DataCollection, error)
	Directory() string
}

type Provider struct {
	store     store.DataStore
	directory string
}

func NewDataProvider(conf config.StoreConfiguration) (*Provider, error) {
	st, err := store.NewStore(conf.ConnectionString)
	if err != nil {
		if os.IsNotExist(err) {
			c, err := store.BuildOptions(conf.ConnectionString)
			if err != nil {
				return nil, err
			}
			os.Mkdir(c.Path, 0755)
		} else {
			return nil, err
		}
	}

	pth := strings.Split(conf.ConnectionString, ";")[2]
	if strings.Contains(pth, "?") {
		pth = strings.Split(pth, "?")[0]
	}

	return &Provider{store: st, directory: pth}, nil
}

func (ds *Provider) GetCollection(ctx context.Context, name string) (collection.DataCollection, error) {
	return ds.store.Collection(ctx, name)
}
func (ds *Provider) Directory() string {
	return ds.directory
}

func (ds *Provider) Exists(ctx context.Context, collection string) bool {
	return ds.store.CollectionExists(context.Background(), collection)
}
