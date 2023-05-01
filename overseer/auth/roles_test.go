package auth

import (
	"errors"
	"testing"

	"github.com/przebro/databazaar/result"
	"github.com/przebro/databazaar/store"
	"github.com/przebro/overseer/datastore"
	mockstore "github.com/przebro/overseer/datastore/mock"
	"github.com/przebro/overseer/overseer/config"
	"github.com/stretchr/testify/suite"
)

type AuthRoleTestSuite struct {
	suite.Suite
	roleManager *RoleManager
	store       *mockstore.MockStore
}

func TestAuthRoleTestSuite(t *testing.T) {
	suite.Run(t, new(AuthRoleTestSuite))
}

func (suite *AuthRoleTestSuite) SetupSuite() {

	suite.store = &mockstore.MockStore{}
	roleCollection := &userMockCollection{}

	roleCollection.On("Create", &dsRoleModel{
		RoleModel: RoleModel{
			Name:           "test_role",
			Description:    "description",
			Administration: true,
		},
		ID: idFormatter(rolesNamespace, "test_role"),
	}).Return(&result.BazaarResult{}, nil)

	roleCollection.On("Create", &dsRoleModel{
		RoleModel: RoleModel{
			Name:           "role_exists",
			Description:    "description",
			Administration: true,
		},
		ID: idFormatter(rolesNamespace, "role_exists"),
	}).Return(&result.BazaarResult{}, errors.New("role exists"))

	roleCollection.On("Get", idFormatter(rolesNamespace, "test_role")).Return(nil)
	roleCollection.On("Get", idFormatter(rolesNamespace, "role_not_exists")).Return(errors.New("not found"))

	suite.store.On("Collection", "auth").Return(roleCollection, nil)

	fn := func(opt store.ConnectionOptions) (store.DataStore, error) {

		return suite.store, nil
	}

	store.RegisterStoreFactory("mock", fn)

	conf := config.StoreConfiguration{ID: "userstore", ConnectionString: "mock;;data/tests?synctime=0"}
	provider, err := datastore.NewDataProvider(conf)
	suite.Nil(err)
	authconf := config.SecurityConfiguration{AllowAnonymous: true}

	suite.roleManager, err = NewRoleManager(authconf, provider)
	suite.Nil(err)
}

func (suite *AuthRoleTestSuite) TestCreateRole_Successful() {

	model := RoleModel{
		Name:           "test_role",
		Description:    "description",
		Administration: true,
	}

	err := suite.roleManager.Create(model)
	suite.Nil(err)
}
func (suite *AuthRoleTestSuite) TestCreateRole_AlreadyExists_Error() {

	model := RoleModel{
		Name:           "role_exists",
		Description:    "description",
		Administration: true,
	}

	err := suite.roleManager.Create(model)
	suite.Error(err)
}

func (suite *AuthRoleTestSuite) TestGetRole_Successful() {

	_, ok := suite.roleManager.Get("test_role")
	suite.True(ok)
}
func (suite *AuthRoleTestSuite) TestGetRole_NotExists_Error() {

	_, ok := suite.roleManager.Get("role_not_exists")
	suite.False(ok)
}
