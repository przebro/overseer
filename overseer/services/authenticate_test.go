package services

import (
	"context"
	"testing"

	"github.com/przebro/overseer/overseer/config"
	"github.com/przebro/overseer/proto/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"google.golang.org/grpc/codes"
)

type mockTokenCreator struct {
	mock.Mock
}

func (m *mockTokenCreator) Verify(token string) (string, error) {
	args := m.Called(token)
	return args.Get(0).(string), args.Error(1)
}
func (m *mockTokenCreator) Create(username string, userdata map[string]interface{}) (string, error) {
	args := m.Called(username, userdata)
	return args.Get(0).(string), args.Error(1)
}

type mockAuthenticationManager struct {
	mock.Mock
}

func (m *mockAuthenticationManager) Authenticate(ctx context.Context, username string, password string) (bool, error) {
	args := m.Called(ctx, username, password)
	return args.Bool(0), args.Error(1)
}

type authenticateTestSuite struct {
	suite.Suite
	service services.AuthenticateServiceServer
	creator *mockTokenCreator
	manager *mockAuthenticationManager
	asrvc   *ovsAuthenticateService
}

func TestAuthenticsteTestSuite(t *testing.T) {
	suite.Run(t, new(authenticateTestSuite))
}

func (suite *authenticateTestSuite) SetupSuite() {

	var authcfg = config.SecurityConfiguration{
		AllowAnonymous: true,
		Timeout:        0,
		Issuer:         "testissuer",
		Secret:         "WBdumgVKBK4iTB+CR2Z2meseDrlnrg54QDSAPcFswWU=",
	}

	suite.creator = &mockTokenCreator{}
	suite.manager = &mockAuthenticationManager{}
	suite.service, _ = NewAuthenticateService(authcfg, suite.creator, suite.manager)
	suite.asrvc = suite.service.(*ovsAuthenticateService)
}

func (suite *authenticateTestSuite) TestAuthenticate_Anonymous_User_Success() {

	service := suite.service
	suite.asrvc.allowAnonymous = true
	msg := &services.AuthorizeActionMsg{Username: "", Password: ""}
	r, err := service.Authenticate(context.Background(), msg)

	suite.Nil(err)
	suite.Equal(r.Message, "anonymous access")
}
func (suite *authenticateTestSuite) TestAuthenticate_Anonymous_User_Fail() {

	service := suite.service

	msg := &services.AuthorizeActionMsg{Username: "", Password: ""}

	suite.asrvc.allowAnonymous = false

	_, err := service.Authenticate(context.Background(), msg)
	suite.NotNil(err)

	_, code := matchExpectedStatusFromError(err, codes.Unauthenticated)
	suite.Equal(code, codes.Unauthenticated)

}

func (suite *authenticateTestSuite) TestAuthenticate_User_Fail() {

	service := suite.service
	suite.asrvc.allowAnonymous = false
	ctx := context.Background()

	msg := &services.AuthorizeActionMsg{Username: "testuser1", Password: ""}

	_, err := service.Authenticate(ctx, msg)
	suite.NotNil(err)
	_, code := matchExpectedStatusFromError(err, codes.Unauthenticated)
	suite.Equal(codes.Unauthenticated, code)

	msg.Username = "testuser1"
	msg.Password = "invalid_password"

	suite.manager.On("Authenticate", ctx, msg.Username, msg.Password).Return(
		false, nil,
	)

	_, err = service.Authenticate(ctx, msg)
	suite.NotNil(err)
	_, code = matchExpectedStatusFromError(err, codes.Unauthenticated)
	suite.Equal(codes.Unauthenticated, code)
}

func (suite *authenticateTestSuite) TestAuthenticate_User_Success() {

	service := suite.service
	suite.asrvc.allowAnonymous = false
	ctx := context.Background()

	msg := &services.AuthorizeActionMsg{Username: "testuser1", Password: "notsecure"}

	suite.manager.On("Authenticate", ctx, msg.Username, msg.Password).Return(
		true, nil,
	)
	suite.creator.On("Create", msg.Username, map[string]interface{}{}).Return("a_token_created", nil)

	r, err := service.Authenticate(ctx, msg)
	suite.Nil(err)
	suite.Equal(r.Success, true)
	suite.Equal(r.Message, "a_token_created")

}
