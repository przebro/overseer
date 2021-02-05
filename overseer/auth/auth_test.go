package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/config"
	"testing"
)

var prvstorecfg = config.StoreProviderConfiguration{
	Store: []config.StoreConfiguration{
		config.StoreConfiguration{ID: "security",
			ConnectionString: "local;/../../data/tests?synctime=0",
		},
	},
	Collections: []config.CollectionConfiguration{
		config.CollectionConfiguration{
			Name:    "securitytest",
			StoreID: "security",
		},
	},
}

var prvconf = config.SecurityConfiguration{
	AllowAnonymous: true,
	Collection:     "securitytest",
}

var prvprovider *datastore.Provider
var prvCollectionName = "securitytest"

var testSecurity = map[string]interface{}{
	"assoc@testuser1": dsRoleAssociationModel{
		ID: "assoc@testuser1",
		RoleAssociationModel: RoleAssociationModel{
			UserID: "testuser1",
			Roles:  []string{"testrole1"},
		},
	},
	"assoc@testuser2": dsRoleAssociationModel{
		ID: "assoc@testuser2",
		RoleAssociationModel: RoleAssociationModel{
			UserID: "testuser2",
			Roles:  []string{"testrole1", "testrole2"},
		},
	},
	"role@testrole1": dsRoleModel{
		ID: "role@testrole1",
		RoleModel: RoleModel{
			Name:           "testrole1",
			Description:    "",
			Administration: true,
			AddTicket:      true,
			RemoveTicket:   true,
			Definition:     true,
		},
	},
	"role@testrole2": dsRoleModel{
		ID: "role@testrole2",
		RoleModel: RoleModel{
			Name:           "testrole2",
			Description:    "",
			Administration: false,
			AddTicket:      false,
			RemoveTicket:   false,
			Definition:     false,
			Bypass:         true,
			Confirm:        true,
			Force:          true,
			Free:           true,
		},
	},
	"user@testuser1": dsUserModel{
		ID: "user@testuser1",
		UserModel: UserModel{
			Enabled:  true,
			FullName: "Test User 1",
			Mail:     "testuser1@test.com",
			Username: "testuser1",
		},
	},
	"user@testuser2": dsUserModel{
		ID: "user@testuser2",
		UserModel: UserModel{
			Enabled:  true,
			FullName: "Test User 2",
			Mail:     "testuser2@test.com",
			Username: "testuser2",
		},
	},
}

func prvprepare(t *testing.T) {

	if prvprovider != nil {
		return
	}

	logger.NewTestLogger()

	var err error
	f, _ := os.Create("../../data/tests/securitytest.json")
	data, err := json.Marshal(testSecurity)
	if err != nil {
		t.Fatal(err)
	}
	f.Write(data)
	f.Close()

	prvprovider, err = datastore.NewDataProvider(prvstorecfg)
	fmt.Println("intiialize:", prvprovider)
	if err != nil {
		t.Fatal("unable to init store")
	}

}

func TestCreateNewAuthorizationProvider(t *testing.T) {
	prvprepare(t)

	_, err := NewAuthorizationManager(prvconf, prvprovider)
	if err != nil {
		t.Error("unexpected result:", err)
	}
}

func TestVerifyAction(t *testing.T) {

	var ok bool
	var err error

	m, err := NewAuthorizationManager(prvconf, prvprovider)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	ok, err = m.VerifyAction(context.Background(), ActionDefinition, "testuser3")

	if err == nil {
		t.Error("unexpected result:", err)
	}

	ok, err = m.VerifyAction(context.Background(), UserAction(99), "testuser1")

	if err != ErrUnableFindAction {
		t.Error("unexpected result:", err, "expected:", ErrUnableFindAction)
	}

	ok, err = m.VerifyAction(context.Background(), ActionDefinition, "testuser1")

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if ok != true {
		t.Error("unexpected result:", ok, "expected:", true)
	}

	ok, err = m.VerifyAction(context.Background(), ActionAdministration, "testuser2")

	if ok != true {
		t.Error("unexpected result:", ok, "expected:", true)
	}

	ok, err = m.VerifyAction(context.Background(), ActionBrowse, "testuser1")

	if ok != true {
		t.Error("unexpected result:", ok, "expected:", true)
	}

	ok, err = m.VerifyAction(context.Background(), ActionSetFlag, "testuser1")

	if ok != false {
		t.Error("unexpected result:", ok, "expected:", false)
	}

	ok, err = m.VerifyAction(context.Background(), ActionSetFlag, "testuser2")

	if ok != false {
		t.Error("unexpected result:", ok, "expected:", false)
	}

	ok, err = m.VerifyAction(context.Background(), ActionAddTicket, "testuser2")

	if ok != true {
		t.Error("unexpected result:", ok, "expected:", true)
	}

	ok, err = m.VerifyAction(context.Background(), ActionRemoveTicket, "testuser2")

	if ok != true {
		t.Error("unexpected result:", ok, "expected:", true)
	}

	ok, err = m.VerifyAction(context.Background(), ActionRestart, "testuser2")

	if ok != false {
		t.Error("unexpected result:", ok, "expected:", false)
	}

	ok, err = m.VerifyAction(context.Background(), ActionSetToOK, "testuser2")

	if ok != false {
		t.Error("unexpected result:", ok, "expected:", false)
	}

}

func TestHashPassword(t *testing.T) {
	pass, err := HashPassword([]byte("notsecure"))
	if err != nil {
		t.Error("unexpected error")
	}

	if pass == "" {
		t.Error("unexpected error")
	}
}
