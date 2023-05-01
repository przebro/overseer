package services

import (
	"context"
	"testing"

	"github.com/przebro/overseer/overseer/auth"
	"github.com/przebro/overseer/proto/services"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockUserManager struct {
	mock.Mock
}

func (m *mockUserManager) Get(ctx context.Context, username string) (auth.UserModel, bool) {
	args := m.Called(username)
	return args.Get(0).(auth.UserModel), args.Get(1).(bool)
}
func (m *mockUserManager) Create(ctx context.Context, model auth.UserModel) error {
	args := m.Called(model)
	return args.Error(0)
}
func (m *mockUserManager) Modify(ctx context.Context, model auth.UserModel) error {
	args := m.Called(model)
	return args.Error(0)
}
func (m *mockUserManager) Delete(ctx context.Context, username string) error {
	args := m.Called(username)
	return args.Error(0)
}
func (m *mockUserManager) All(filter string) ([]auth.UserModel, error) {
	args := m.Called(filter)
	return args.Get(0).([]auth.UserModel), args.Error(1)
}

type mockRoleManager struct {
	mock.Mock
}

func (m *mockRoleManager) Get(name string) (auth.RoleModel, bool) {
	args := m.Called(name)
	return args.Get(0).(auth.RoleModel), args.Get(1).(bool)
}
func (m *mockRoleManager) Create(model auth.RoleModel) error {
	args := m.Called(model)
	return args.Error(0)
}
func (m *mockRoleManager) Modify(model auth.RoleModel) error {

	args := m.Called(model)
	return args.Error(0)
}
func (m *mockRoleManager) Delete(name string) error {
	args := m.Called(name)
	return args.Error(0)
}
func (m *mockRoleManager) All(filter string) ([]auth.RoleModel, error) {
	args := m.Called(filter)
	return args.Get(0).([]auth.RoleModel), args.Error(1)
}

type mockRoleAssociationManager struct {
	mock.Mock
}

func (m *mockRoleAssociationManager) Get(username string) (auth.RoleAssociationModel, bool) {
	args := m.Called(username)
	return args.Get(0).(auth.RoleAssociationModel), args.Get(1).(bool)
}
func (m *mockRoleAssociationManager) Create(model auth.RoleAssociationModel) error {
	args := m.Called(model)
	return args.Error(0)
}
func (m *mockRoleAssociationManager) Modify(model auth.RoleAssociationModel) error {
	args := m.Called(model)
	return args.Error(0)
}
func (m *mockRoleAssociationManager) Delete(username string) error {
	args := m.Called(username)
	return args.Error(0)
}

type adminTestSuite struct {
	suite.Suite
	service  services.AdministrationServiceServer
	srvc     *ovsAdministrationService
	umanager *mockUserManager
}

func TestAdminTestSuite(t *testing.T) {
	suite.Run(t, new(adminTestSuite))
}
func (suite *adminTestSuite) SetupSuite() {

	suite.umanager = &mockUserManager{}
	suite.service = NewAdministrationService(suite.umanager, &mockRoleManager{}, &mockRoleAssociationManager{})
	suite.srvc = suite.service.(*ovsAdministrationService)
}

func (suite *adminTestSuite) TestGetAllowedAction() {

	act := suite.srvc.GetAllowedAction("service")
	suite.Equal(auth.ActionAdministration, act)
}

func (suite *adminTestSuite) TestListUsers() {

	client := suite.service

	suite.umanager.On("All", "").Return(
		[]auth.UserModel{}, nil)

	r, err := client.ListUsers(context.Background(), &services.FilterMsg{Filter: ""})
	suite.Nil(err)
	suite.Len(r.Users, 0)
}

/*
func TestGetUser(t *testing.T) {

	client := createAdminCLient(t)

	_, err := client.GetUser(context.Background(), &services.UserMsg{Username: "user_not_exists"})

	if err == nil {
		t.Error("unexpected result")

	}
	r, err := client.GetUser(context.Background(), &services.UserMsg{Username: "testuser1"})

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.User.Username != "testuser1" {
		t.Error("unexpected result:", r.User)
	}

	if len(r.User.Roles) != 1 || r.User.Roles[0] != "testrole1" {

		t.Error("unexpected result:", r.User.Roles)
	}
}

func TestCreateUser(t *testing.T) {

	client := createAdminCLient(t)

	msg := &services.CreateUserMsg{}

	r, err := client.CreateUser(context.Background(), msg)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}
	msg.User = &services.UserAccount{}

	msg.User.Username = "very_long_name_for_test_user_that_exceeds_size_limit"
	msg.Password = "1"

	r, err = client.CreateUser(context.Background(), msg)
	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.User.Username = "testuser3"
	msg.User.Fullname = "very_long_text_in_full_name_that_exceeds_limit_of_sixtyfour_characters"

	r, err = client.CreateUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.User.Fullname = "Test User"
	msg.User.Email = "notvalidemail"

	r, err = client.CreateUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.User.Email = "field_not_valid_that_exceeds_limit_of_64_characters_email@overseer.com"

	r, err = client.CreateUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.User.Email = "testuser3@overseer.com"
	//validates password length
	r, err = client.CreateUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.User.Email = "testuser3@overseer.com"
	msg.Password = "very_long_password_that_exceeds_32_characters"
	r, err = client.CreateUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.User.Email = "testuser3@overseer.com"
	msg.Password = "notsecure"
	msg.User.Enabled = false
	msg.User.Roles = []string{"role_not_exists"}
	r, err = client.CreateUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.User.Email = "testuser3@overseer.com"
	msg.Password = "notsecure"
	msg.User.Enabled = false
	msg.User.Roles = []string{"testrole2"}
	r, err = client.CreateUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	//Canot create user that already exists
	msg.User.Username = "testuser1"

	r, err = client.CreateUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

}

func TestModifyUser(t *testing.T) {

	client := createAdminCLient(t)

	msg := &services.CreateUserMsg{
		User: &services.UserAccount{Username: "testuser5"},
	}

	r, err := client.ModifyUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.User.Username = "testuser1"
	msg.User.Fullname = "testuser1 modified"
	msg.User.Email = "testuser1overseer"
	msg.User.Enabled = true
	r, err = client.ModifyUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.User.Email = "testuser1@overseer.com"
	r, err = client.ModifyUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

}

func TestDeleteUser(t *testing.T) {

	client := createAdminCLient(t)

	msg := &services.UserMsg{Username: "user_not_exists"}

	r, err := client.DeleteUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)

	}

	msg.Username = "testuser3"

	r, err = client.DeleteUser(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)

	}
}

func TestGetRole(t *testing.T) {

	client := createAdminCLient(t)

	_, err := client.GetRole(context.Background(), &services.RoleMsg{Rolename: "role_not_exists"})

	if err == nil {
		t.Error("unexpected result")

	}
	r, err := client.GetRole(context.Background(), &services.RoleMsg{Rolename: "testrole1"})

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Role.Role.Rolename != "testrole1" {
		t.Error("unexpected result:", r.Role.Role.Rolename)
	}

}

func TestCreateRole(t *testing.T) {

	client := createAdminCLient(t)

	msg := &services.RoleDefinitionMsg{
		Role: &services.RoleMsg{},
	}
	r, err := client.CreateRole(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Role.Rolename = "very_long_rolename_that_exceeds_32_characters"
	msg.Description = "description"
	r, err = client.CreateRole(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Role = &services.RoleMsg{Rolename: "testrole3"}
	msg.Description = "very_very_long_role_description_field_that_exceeds_sixtyfour_character"
	r, err = client.CreateRole(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Role = &services.RoleMsg{Rolename: "testrole3"}
	msg.Description = "test description"
	r, err = client.CreateRole(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

	msg.Role = &services.RoleMsg{Rolename: "testrole1"}

	r, err = client.CreateRole(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

}

func TestModifyRole(t *testing.T) {

	client := createAdminCLient(t)

	msg := &services.RoleDefinitionMsg{
		Role: &services.RoleMsg{},
	}
	r, err := client.ModifyRole(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Role.Rolename = "testrole1"
	msg.Description = "role description"
	msg.Bypass = true
	r, err = client.ModifyRole(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}
}

func TestDeleteRole(t *testing.T) {
	client := createAdminCLient(t)

	msg := &services.RoleMsg{Rolename: "role_that_does_not_exists"}

	r, err := client.DeleteRole(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != false {
		t.Error("unexpected result:", r.Success, "expected:", false)
	}

	msg.Rolename = "testrole3"

	r, err = client.DeleteRole(context.Background(), msg)

	if err != nil {
		t.Error("unexpected result:", err)
	}

	if r.Success != true {
		t.Error("unexpected result:", r.Success, "expected:", true)
	}

}

func TestListRoles(t *testing.T) {

	client := createAdminCLient(t)

	r, err := client.ListRoles(context.Background(), &services.FilterMsg{Filter: ""})

	if err != nil {
		t.Error("unexpected result:", err)
	}
	if len(r.Entity) != 2 {
		t.Error("unexpected result:", len(r.Entity))
	}
	for _, x := range r.Entity {
		if x.Description == "" || x.Name == "" {
			t.Error("unexpected result")
		}
	}
}
*/
