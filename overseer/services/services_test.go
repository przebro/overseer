package services

import (
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/auth"
	"github.com/przebro/overseer/overseer/config"
	"github.com/przebro/overseer/overseer/internal/resources"
	"github.com/przebro/overseer/overseer/internal/taskdef"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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

var provcfg = config.StoreConfiguration{
	ID:               "servicestore",
	ConnectionString: "local;/../../data/tests?synctime=1",
}

var rescfg = config.ResourcesConfigurartion{
	Resources: config.ResourceEntry{Sync: 3600},
}

var testCollectionName = "tasks"
var taskPoolConfig config.ActivePoolConfiguration = config.ActivePoolConfiguration{
	ForceNewDayProc: true, MaxOkReturnCode: 4,
	NewDayProc: "00:30",
	SyncTime:   5,
}

var provider *datastore.Provider
var resmanager resources.ResourceManager

var definitionManagerT taskdef.TaskDefinitionManager

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
