package services

import (
	"encoding/json"
	"fmt"
	"os"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/auth"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/resources"

	"google.golang.org/grpc"
)

type mockDispacher struct {
}

func (m *mockDispacher) PushEvent(sender events.EventReceiver, route events.RouteName, msg events.DispatchedMessage) error {
	return nil
}
func (m *mockDispacher) Subscribe(route events.RouteName, participant events.EventParticipant) {

}
func (m *mockDispacher) Unsubscribe(route events.RouteName, participant events.EventParticipant) {

}

type mockBuffconnServer struct {
	grpcServer *grpc.Server
}

type testRoleAssociationModel struct {
	auth.RoleAssociationModel
	ID  string `json:"_id" bson:"_id" validate:"required,max=38"` // prefix='assoc' + @ + username
	Rev string `json:"_rev,omitempty" bson:"_rev,omitempty"`
}

type testRoleModel struct {
	auth.RoleModel
	ID  string `json:"_id" bson:"_id" validate:"required,max=37"` // prefix='role' + @ + name
	Rev string `json:"_rev,omitempty" bson:"_rev,omitempty"`
}

type testUserModel struct {
	auth.UserModel
	ID  string `json:"_id" bson:"_id" validate:"required,max=37"` // prefix='user' + @ + username
	Rev string `json:"_rev,omitempty" bson:"_rev,omitempty"`
}

var securitymodel = map[string]interface{}{
	"assoc@testuser1": testRoleAssociationModel{
		ID: "assoc@testuser1",
		RoleAssociationModel: auth.RoleAssociationModel{
			UserID: "testuser1",
			Roles:  []string{"testrole1"},
		},
	},
	"assoc@testuser2": testRoleAssociationModel{
		ID: "assoc@testuser2",
		RoleAssociationModel: auth.RoleAssociationModel{
			UserID: "testuser2",
			Roles:  []string{"testrole1", "testrole2"},
		},
	},
	"role@testrole1": testRoleModel{
		ID: "role@testrole1",
		RoleModel: auth.RoleModel{
			Name:           "testrole1",
			Description:    "Test role description",
			Administration: true,
			AddTicket:      true,
			RemoveTicket:   true,
			Definition:     true,
		},
	},
	"role@testrole2": testRoleModel{
		ID: "role@testrole2",
		RoleModel: auth.RoleModel{
			Name:           "testrole2",
			Description:    "Test role description",
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
	"user@testuser1": testUserModel{
		ID: "user@testuser1",
		UserModel: auth.UserModel{
			Enabled:  true,
			FullName: "Test User 1",
			Mail:     "testuser1@test.com",
			Username: "testuser1",
			Password: "$2a$04$2ia6jc5Ob49dwj85tNfNSeUrk6aC3AWenBaR4BZtJrX21Fn5Hp.Ui", //bcrypted: notsecure
		},
	},
	"user@testuser2": testUserModel{
		ID: "user@testuser2",
		UserModel: auth.UserModel{
			Enabled:  true,
			FullName: "Test User 2",
			Mail:     "testuser2@test.com",
			Username: "testuser2",
			Password: "$2a$04$2ia6jc5Ob49dwj85tNfNSeUrk6aC3AWenBaR4BZtJrX21Fn5Hp.Ui", //bcrypted: notsecure
		},
	},
}

var authcfg = config.SecurityConfiguration{
	Collection:     "serviceusers",
	AllowAnonymous: true,
	Timeout:        0,
	Issuer:         "testissuer",
	Secret:         "WBdumgVKBK4iTB+CR2Z2meseDrlnrg54QDSAPcFswWU=",
}

var provcfg = config.StoreProviderConfiguration{
	Store: []config.StoreConfiguration{
		{
			ID:               "servicestore",
			ConnectionString: "local;/../../data/tests?synctime=1",
		},
	},
	Collections: []config.CollectionConfiguration{

		{Name: "serviceusers", StoreID: "servicestore"},
		{Name: "resources", StoreID: "servicestore"},
	},
}

var rescfg = config.ResourcesConfigurartion{
	TicketSource: config.ResourceEntry{Sync: 1, Collection: "resources"},
	FlagSource:   config.ResourceEntry{Sync: 1, Collection: "resources"},
}

var provider *datastore.Provider
var resmanager resources.ResourceManager

var dispatcher mockDispacher = mockDispacher{}

func init() {

	var err error
	f, _ := os.Create("../../data/tests/serviceusers.json")

	data, _ := json.Marshal(securitymodel)
	f.Write(data)
	f.Close()

	log := logger.NewTestLogger()

	if provider, err = datastore.NewDataProvider(provcfg); err != nil {
		panic("")

	}

	if resmanager, err = resources.NewManager(&dispatcher, log, rescfg, provider); err != nil {
		fmt.Println(err)
		panic("")
	}

}