package datastore

import (
	"context"
	"errors"
	"fmt"
	"overseer/common/logger"
	"overseer/overseer/config"

	"github.com/przebro/databazaar/store"
	_ "github.com/przebro/localstore/store"

	"github.com/przebro/databazaar/collection"
)

//Provider - serves as mediator between data storage (database,local file,...) and in memory data structure
type Provider struct {
	store       map[string]store.DataStore
	collections map[string]string
}

var (
	ErrStoreConfiguration = errors.New("Store configuration error")
)

//GetCollection - Gets collection by name
func (p *Provider) GetCollection(name string) (collection.DataCollection, error) {

	var storeName string
	var exists bool
	var st store.DataStore

	if storeName, exists = p.collections[name]; !exists {
		return nil, fmt.Errorf("%w; no mapping beetween collection name:%s and store", ErrStoreConfiguration, name)
	}
	if st, exists = p.store[storeName]; !exists {
		return nil, fmt.Errorf("%w;no mapping beetween collection name and store", ErrStoreConfiguration)
	}
	col, err := st.Collection(context.Background(), name)

	return col, err
}

//NewDataProvider - Creates a new DataProvider
func NewDataProvider(conf config.StoreProviderConfiguration) (*Provider, error) {

	var err error
	p := &Provider{}

	if p.store, err = loadStoreData(conf.Store); err != nil {
		return nil, err
	}

	if p.collections, err = loadCollectionData(conf.Collections); err != nil {
		return nil, err
	}

	return p, nil

}
func loadStoreData(conf []config.StoreConfiguration) (map[string]store.DataStore, error) {

	log := logger.Get()

	connmap := map[string]string{}
	smap := map[string]store.DataStore{}

	for _, s := range conf {
		log.Info("loading store:", s.ID, ":", s.ConnectionString)
		if s.ID == "" || s.ConnectionString == "" {
			return nil, fmt.Errorf("%w; Invalid store configuration entry", ErrStoreConfiguration)
		}
		if _, exists := connmap[s.ID]; exists {
			return nil, fmt.Errorf("%w; Duplicated store entry, id:%s", ErrStoreConfiguration, s.ID)
		}

		connmap[s.ID] = s.ConnectionString

	}

	for k, v := range connmap {

		st, err := store.NewStore(v)

		if err != nil {

			return nil, fmt.Errorf("unable to connect store:%w", err)
		}

		smap[k] = st
	}

	return smap, nil
}

func loadCollectionData(conf []config.CollectionConfiguration) (map[string]string, error) {

	cmap := map[string]string{}
	log := logger.Get()

	for _, col := range conf {

		log.Info(col.Name, ":", col.StoreID)
		if col.Name == "" || col.StoreID == "" {
			return nil, errors.New("invalid collection configuration entry")
		}

		if _, exists := cmap[col.Name]; exists {
			return nil, fmt.Errorf("duplicated collection entry, name:%s", col.Name)
		}

		cmap[col.Name] = col.StoreID
	}

	return cmap, nil

}
