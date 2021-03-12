package auth

import (
	"os"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/config"
	"testing"
)

var storecfg = config.StoreProviderConfiguration{
	Store: []config.StoreConfiguration{
		{ID: "userstore",
			ConnectionString: "local;/../../data/tests?synctime=0",
		},
	},
	Collections: []config.CollectionConfiguration{
		{
			Name:    "authtest",
			StoreID: "userstore",
		},
	},
}

var conf = config.SecurityConfiguration{
	AllowAnonymous: true,
	Collection:     "authtest",
}

var notsecure = "$2a$04$EFHkGN6rDONfCE1Oa4FTcOVC4yFgsMtX4AB87cMgip4yxQpCIIixi"

var provider *datastore.Provider

func prepare(t *testing.T) {

	if provider != nil {
		return
	}

	log := logger.NewTestLogger()

	var err error
	f, _ := os.Create("../../data/tests/authtest.json")
	f.Write([]byte("{}"))
	f.Close()

	provider, err = datastore.NewDataProvider(storecfg, log)
	if err != nil {
		t.Fatal("unable to init store")
	}

}

func TestNewManager(t *testing.T) {

	prepare(t)

	_, err := NewUserManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	cfg := config.SecurityConfiguration{Collection: "invalid_name"}
	_, err = NewUserManager(cfg, provider)
	if err == nil {
		t.Error("unexpected result")
	}
}

func TestCreateUser(t *testing.T) {

	m, err := NewUserManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	model := UserModel{
		Enabled:  true,
		FullName: "test user",
		Username: "test",
		Mail:     "test@test.com",
		Password: notsecure,
	}
	err = m.Create(model)
	if err != nil {
		t.Error("create user,unexpected result:", err)
	}

	err = m.Create(model)

	if err == nil {
		t.Error("create user,unexpected result:", err)
	}
}
func TestGetUser(t *testing.T) {

	m, err := NewUserManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	_, ok := m.Get("test")
	if ok == false {
		t.Error("unexpected result:", ok, "expected:", true)
	}

	_, ok = m.Get("test2")
	if ok == true {
		t.Error("unexpected result:", ok, "expected:", false)
	}

}

func TestModifyUser(t *testing.T) {

	m, err := NewUserManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	model := UserModel{
		Enabled:  true,
		FullName: "test user",
		Username: "test2",
		Mail:     "test@test.com",
		Password: notsecure,
	}

	err = m.Modify(model)

	if err == nil {
		t.Error("unexpected result:", err)
	}

	err = m.Create(model)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	model.Mail = "changed@test.com"

	err = m.Modify(model)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	nmodel, _ := m.Get("test2")
	if nmodel.Mail != model.Mail {
		t.Error("unexpected result:", nmodel.Mail, "expected:", model.Mail)
	}

}

func TestGetAll(t *testing.T) {

	m, err := NewUserManager(conf, provider)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	r, err := m.All("")
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	if len(r) != 2 {
		t.Error("unexpected result:", len(r), "expected 2")
	}
}

func TestDeleteUser(t *testing.T) {

	m, err := NewUserManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	err = m.Delete("test2")
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	err = m.Delete("test2")

	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

}

func TestCreatePassword(t *testing.T) {

	m, err := NewUserManager(conf, provider)
	if err != nil {
		t.Error("unexpected resutlt:", err)
	}

	_, err = m.CheckChangePassword([]byte(notsecure), []byte("notsecure"), []byte("notsecure"))
	if err != nil {
		t.Error("change password, unexpected result:", err)
	}

	_, err = m.CheckChangePassword([]byte(notsecure), []byte("notsecure2"), []byte("notsecure"))
	if err == nil {
		t.Error("change password, unexpected result:", err)
	}
}
