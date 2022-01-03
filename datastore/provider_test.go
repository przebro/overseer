package datastore

import (
	"errors"
	"os"
	"testing"

	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/overseer/config"
)

var resourceString = `{"flags":{"_id":"flags","_rev":"","flags":null},"tickets":{"_id":"tickets","_rev":"","tickets":[]}}`

func init() {
	f, _ := os.Create("../data/tests/resources.json")
	f.Write([]byte(resourceString))
	f.Close()
}

func TestProvider(t *testing.T) {

	log := logger.NewTestLogger()

	conf := config.StoreProviderConfiguration{
		Store: []config.StoreConfiguration{
			{ID: "", ConnectionString: ""},
		},
		Collections: []config.CollectionConfiguration{
			{Name: "", StoreID: ""},
		},
	}
	_, err := NewDataProvider(conf, log)

	if err == nil || !errors.Is(err, ErrStoreConfiguration) {
		t.Errorf("unexpected result")
	}

	conf.Store[0] = config.StoreConfiguration{ID: "test1", ConnectionString: "local;/../data/tests?synctime=1"}
	conf.Store = append(conf.Store, config.StoreConfiguration{ID: "test1", ConnectionString: "local;/../data/tests?synctime=1"})

	_, err = NewDataProvider(conf, log)

	if err == nil || !errors.Is(err, ErrStoreConfiguration) {
		t.Error("unexpected result:", err)
	}

	conf.Store[1] = config.StoreConfiguration{ID: "test2", ConnectionString: "local;/data/tests?synctime=1"}

	_, err = NewDataProvider(conf, log)

	if err == nil {
		t.Error("unexpected result", err)
	}

	conf.Store = []config.StoreConfiguration{{ID: "test1", ConnectionString: "local;/../data/tests?synctime=1"}}

	_, err = NewDataProvider(conf, log)

	if err == nil {
		t.Error("unexpected result", err)
	}

}

func TestProviderCollections(t *testing.T) {

	log := logger.NewTestLogger()

	conf := config.StoreProviderConfiguration{
		Store: []config.StoreConfiguration{
			{ID: "test1", ConnectionString: "local;/../data/tests?synctime=1"},
		},
		Collections: []config.CollectionConfiguration{
			{Name: "", StoreID: ""},
		},
	}
	_, err := NewDataProvider(conf, log)

	if err == nil {
		t.Error("unexpected result", err)
	}

	conf.Collections[0] = config.CollectionConfiguration{Name: "resources", StoreID: "test1"}

	_, err = NewDataProvider(conf, log)

	if err != nil {
		t.Error("unexpected result", err)
	}

	conf.Collections = append(conf.Collections, config.CollectionConfiguration{Name: "resources", StoreID: "test1"})

	_, err = NewDataProvider(conf, log)

	if err == nil {
		t.Error("unexpected result", err)
	}

}

func TestGetCollection(t *testing.T) {

	log := logger.NewTestLogger()

	conf := config.StoreProviderConfiguration{
		Store: []config.StoreConfiguration{
			{ID: "test1", ConnectionString: "local;/../data/tests?synctime=1"},
		},
		Collections: []config.CollectionConfiguration{
			{Name: "resources", StoreID: "test1"},
		},
	}
	p, err := NewDataProvider(conf, log)

	if err != nil {
		t.Error("unexpected result", err)
	}

	_, err = p.GetCollection("noname")
	if err == nil {
		t.Error("unexpected error")
	}

	col, err := p.GetCollection("resources")
	if err != nil {
		t.Error("unexpected error:", err)
	}
	if col == nil {
		t.Error("unexpected error nil collection")
	}

	conf.Collections = []config.CollectionConfiguration{{Name: "resources2", StoreID: "not_exists"}}

	prov, err := NewDataProvider(conf, log)

	if err != nil {
		t.Error("unexpected result", err)
	}

	_, err = prov.GetCollection("resources2")

	if err == nil {
		t.Error("unexpected result")
	}
}
