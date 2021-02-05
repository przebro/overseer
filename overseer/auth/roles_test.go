package auth

import (
	"os"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/config"
	"testing"
)

var rstorecfg = config.StoreProviderConfiguration{
	Store: []config.StoreConfiguration{
		config.StoreConfiguration{ID: "userstore",
			ConnectionString: "local;/../../data/tests?synctime=0",
		},
	},
	Collections: []config.CollectionConfiguration{
		config.CollectionConfiguration{
			Name:    "authtest",
			StoreID: "userstore",
		},
	},
}

var rconf = config.SecurityConfiguration{
	AllowAnonymous: true,
	Collection:     "authtest",
}

var rprovider *datastore.Provider

func rprepare(t *testing.T) {

	if provider != nil {
		return
	}

	logger.NewTestLogger()

	var err error
	f, _ := os.Create("../../data/tests/authtest.json")
	f.Write([]byte("{}"))
	f.Close()

	provider, err = datastore.NewDataProvider(rstorecfg)
	if err != nil {
		t.Fatal("unable to init store")
	}

}

func TestNewRoleManager(t *testing.T) {

	rprepare(t)

	_, err := NewRoleManager(rconf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}
}

func TestCreateRole(t *testing.T) {
	m, err := NewRoleManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	model := RoleModel{
		Name:           "testrole",
		Description:    "description",
		Administration: true,
	}

	err = m.Create(model)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	err = m.Create(model)
	if err == nil {
		t.Error("unexpected result:", err)
	}
}

func TestGetRole(t *testing.T) {
	m, err := NewRoleManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	_, ok := m.Get("testrole")
	if ok != true {
		t.Error("unexpected result:", err)
	}

	_, ok = m.Get("testrole2")
	if ok == true {
		t.Error("unexpected result:", err)
	}
}

func TestModifyRole(t *testing.T) {
	m, err := NewRoleManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	model := RoleModel{
		Name:           "testrole2",
		Description:    "description",
		Administration: true,
	}

	err = m.Modify(model)
	if err == nil {
		t.Error("unexpected result:", nil)
	}

	err = m.Create(model)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	model.Administration = false
	model.Description = "changed description"

	err = m.Modify(model)

	if err != nil {
		t.Error("unexpected result:", nil)
	}

	nmodel, _ := m.Get("testrole2")

	if nmodel.Description != "changed description" && nmodel.Administration != false {
		t.Error("unexpected result:", nmodel.Description, nmodel.Administration, "expected:", model.Description, model.Administration)
	}

}

func TestGetAllRoles(t *testing.T) {

	m, err := NewRoleManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	r, err := m.All("")

	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	if len(r) != 2 {
		t.Error("unexpected result:", len(r), "expected 2")
	}

}
func TestDeleteRole(t *testing.T) {

	m, err := NewRoleManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	err = m.Delete("testrole2")
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	err = m.Delete("testrole2")

	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

}
