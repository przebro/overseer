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
	"overseer/overseer/internal/pool"
	"overseer/overseer/internal/resources"
	"overseer/overseer/internal/taskdef"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		{Name: "tasks", StoreID: "servicestore"},
		{Name: "sequence", StoreID: "servicestore"},
	},
}

var rescfg = config.ResourcesConfigurartion{
	TicketSource: config.ResourceEntry{Sync: 3600, Collection: "resources"},
	FlagSource:   config.ResourceEntry{Sync: 3600, Collection: "resources"},
}

var testCollectionName = "tasks"
var taskPoolConfig config.ActivePoolConfiguration = config.ActivePoolConfiguration{
	ForceNewDayProc: true, MaxOkReturnCode: 4,
	NewDayProc: "00:30",
	SyncTime:   5,
	Collection: testCollectionName,
}

var provider *datastore.Provider
var resmanager resources.ResourceManager
var dispatcher mockDispacher = mockDispacher{}

var taskPoolT *pool.ActiveTaskPool
var activeTaskManagerT *pool.ActiveTaskPoolManager
var definitionManagerT taskdef.TaskDefinitionManager

func init() {

	var err error
	f, _ := os.Create("../../data/tests/serviceusers.json")

	data, _ := json.Marshal(securitymodel)
	f.Write(data)
	f.Close()

	f1, _ := os.Create("../../../data/tests/resources.json")
	f1.Write([]byte(`{"flags":{"_id":"flags","_rev":"","flags":[]},"tickets":{"_id":"tickets","_rev":"","tickets":[]}}`))
	f1.Close()

	log := logger.NewTestLogger()

	if provider, err = datastore.NewDataProvider(provcfg, log); err != nil {
		panic("")
	}

	if resmanager, err = resources.NewManager(&dispatcher, log, rescfg, provider); err != nil {
		fmt.Println(err)
		panic("")
	}

	f2, _ := os.Create("../../../data/tests/tasks.json")
	f2.Write([]byte("{}"))
	f2.Close()

	f3, _ := os.Create("../../../data/tests/sequence.json")
	f3.Write([]byte(`{}`))
	f3.Close()

	initTaskPool()

	path, _ := filepath.Abs("../../def")
	definitionManagerT, err = taskdef.NewManager(path, log)
	if err != nil {
		fmt.Println(err)
	}
	activeTaskManagerT, _ = pool.NewActiveTaskPoolManager(&dispatcher, definitionManagerT, taskPoolT, provider, log)

}

func initTaskPool() {
	taskPoolT, _ = pool.NewTaskPool(&dispatcher, taskPoolConfig, provider, true, logger.NewTestLogger())
}

func matchExpectedStatusFromError(err error, expected codes.Code) (bool, codes.Code) {

	var sts *status.Status
	var ok bool

	if sts, ok = status.FromError(err); !ok {
		return ok, codes.Code(9999)
	}

	if sts.Code() != expected {
		return false, sts.Code()
	}

	return true, sts.Code()
}
