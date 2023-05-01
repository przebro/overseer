package auth

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/przebro/databazaar/collection"
	"github.com/przebro/databazaar/result"
	"github.com/przebro/databazaar/selector"
	"github.com/przebro/databazaar/store"
	"github.com/przebro/overseer/datastore"
	mockstore "github.com/przebro/overseer/datastore/mock"
	"github.com/przebro/overseer/overseer/config"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

var notsecure = "$2a$04$EFHkGN6rDONfCE1Oa4FTcOVC4yFgsMtX4AB87cMgip4yxQpCIIixi"

type userMockCollection struct {
	mock.Mock
}

func (m *userMockCollection) Create(ctx context.Context, document interface{}) (*result.BazaarResult, error) {
	args := m.Called(document)
	return args.Get(0).(*result.BazaarResult), args.Error(1)
}
func (m *userMockCollection) Get(ctx context.Context, id string, result interface{}) error {
	args := m.Called(id)
	return args.Error(0)
}
func (m *userMockCollection) Update(ctx context.Context, doc interface{}) error {
	return nil
}
func (m *userMockCollection) Delete(ctx context.Context, id string) error {
	return nil
}
func (m *userMockCollection) CreateMany(ctx context.Context, docs []interface{}) ([]result.BazaarResult, error) {
	return nil, nil
}
func (m *userMockCollection) BulkUpdate(ctx context.Context, docs []interface{}) error {
	return nil
}
func (m *userMockCollection) All(ctx context.Context) (collection.BazaarCursor, error) {
	return nil, nil
}
func (m *userMockCollection) Count(ctx context.Context) (int64, error) {
	return 0, nil
}
func (m *userMockCollection) AsQuerable() (collection.QuerableCollection, error) {
	return nil, nil
}
func (m *userMockCollection) Select(ctx context.Context, s selector.Expr, fld selector.Fields) (collection.BazaarCursor, error) {
	return nil, nil
}

func (m *userMockCollection) Type() string {
	return "mock"
}

type AuthUserTestSuite struct {
	suite.Suite
	userManager *UserManager
	store       *mockstore.MockStore
}

func TestAuthUserTestSuite(t *testing.T) {
	suite.Run(t, new(AuthUserTestSuite))
}

func (suite *AuthUserTestSuite) SetupSuite() {

	suite.store = &mockstore.MockStore{}
	userCollection := &userMockCollection{}

	userCollection.On("Create", &dsUserModel{
		UserModel: UserModel{
			Username: "test",
			FullName: "test user",
			Enabled:  true,
			Password: notsecure,
			Mail:     "test@test.com",
		},
		ID: idFormatter(userNamespace, "test"),
	}).Return(&result.BazaarResult{}, nil)

	userCollection.On("Create", &dsUserModel{
		UserModel: UserModel{
			Username: "user_exists",
			FullName: "user exists",
			Enabled:  true,
			Password: notsecure,
			Mail:     "user_exists@test.com",
		},
		ID: idFormatter(userNamespace, "user_exists"),
	}).Return(&result.BazaarResult{}, errors.New("user exists"))

	userCollection.On("Get", idFormatter(userNamespace, "test_user")).Return(nil)
	userCollection.On("Get", idFormatter(userNamespace, "user_not_exists")).Return(errors.New("not found"))

	suite.store.On("Collection", "auth").Return(userCollection, nil)

	fn := func(opt store.ConnectionOptions) (store.DataStore, error) {

		return suite.store, nil
	}

	store.RegisterStoreFactory("mock", fn)

	conf := config.StoreConfiguration{ID: "userstore", ConnectionString: "mock;;data/tests?synctime=0"}
	provider, err := datastore.NewDataProvider(conf)
	suite.Nil(err)
	authconf := config.SecurityConfiguration{AllowAnonymous: true}

	suite.userManager, err = NewUserManager(authconf, provider)
	suite.Nil(err)

}

func (suite *AuthUserTestSuite) TestCreateUser_Successful() {

	model := UserModel{
		Enabled:  true,
		FullName: "test user",
		Username: "test",
		Mail:     "test@test.com",
		Password: notsecure,
	}
	err := suite.userManager.Create(context.TODO(), model)
	suite.Nil(err)
}
func (suite *AuthUserTestSuite) TestCreateUser_Error() {

	model := UserModel{
		Enabled:  true,
		FullName: "user exists",
		Username: "user_exists",
		Mail:     "user_exists@test.com",
		Password: notsecure,
	}
	err := suite.userManager.Create(context.TODO(), model)
	suite.Error(err)
}

func (suite *AuthUserTestSuite) TestGetUser_Successful() {
	model, exists := suite.userManager.Get(context.TODO(), "test_user")
	suite.True(exists)
	fmt.Println(model)
}

func (suite *AuthUserTestSuite) TestGetUser_NotExists() {
	model, exists := suite.userManager.Get(context.TODO(), "user_not_exists")
	suite.False(exists)
	fmt.Println(model)
}
