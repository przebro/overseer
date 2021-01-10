package auth

import (
	"context"
	"errors"
	"fmt"
	"goscheduler/datastore"

	"golang.org/x/crypto/bcrypt"

	"github.com/przebro/databazaar/collection"
)

const (
	rolesNamespace = "role"
	userNamespace  = "user"
	assocNamespace = "assoc"
)

//UserAction - type for user action
type UserAction int

//Defines possible user actions
const (
	ActionUserManagment UserAction = iota
	ActionRoleManagment
	ActionRestartAction
	ActionSetToOK
	ActionAddTicket
	ActionRemoveTicket
	ActionSetFlag
	ActionConfirm
	ActionOrder
	ActionForce
	ActionDefinition
)

func idFormatter(prefix, id string) string {
	return fmt.Sprintf("%s@%s", prefix, id)
}

//AuthenticationManager - Authenticate users
type AuthenticationManager interface {
	Authenticate(ctx context.Context, username string, password string) (bool, error)
}

//AuthorizationManager - Provides
type AuthorizationManager struct {
	col collection.DataCollection
}

func NewAuthorizationManager(collectionName string, provider *datastore.Provider) (*AuthorizationManager, error) {

	var col collection.DataCollection
	var err error

	if col, err = provider.GetCollection(collectionName); err != nil {
		return nil, err
	}

	return &AuthorizationManager{col: col}, nil
}

func (m *AuthorizationManager) VerifyAction(ctx context.Context, action UserAction, username string) (bool, error) {

	if ctx == nil {
		ctx = context.Background()
	}
	model := dsRoleAssociationModel{}
	if err := m.col.Get(ctx, idFormatter(assocNamespace, username), &model); err != nil {
		return false, errors.New("unable to get role association for give nuser")
	}

	roles := []RoleModel{}

	for x := range model.Roles {

		rmodel := RoleModel{}
		if err := m.col.Get(ctx, idFormatter(rolesNamespace, model.Roles[x]), &rmodel); err != nil {
			return false, fmt.Errorf("unable to verify action, role %s does not exists", model.Roles[x])
		}

		roles = append(roles, rmodel)

	}

	finalRole := m.getEffectiveRights(roles)

	switch action {
	case ActionUserManagment:
		return finalRole.UserManagment, nil
	case ActionRoleManagment:
		return finalRole.RoleManagment, nil
	case ActionRestartAction:
		return finalRole.Restart, nil
	case ActionSetToOK:
		return finalRole.SetToOK, nil
	case ActionAddTicket:
		return finalRole.AddTicket, nil
	case ActionRemoveTicket:
		return finalRole.RemoveTicket, nil
	case ActionSetFlag:
		return finalRole.SetFlag, nil
	case ActionConfirm:
		return finalRole.Confirm, nil
	case ActionOrder:
		return finalRole.Order, nil
	case ActionForce:
		return finalRole.Force, nil
	case ActionDefinition:
		return finalRole.Definition, nil
	}
	return false, errors.New("unable to find action")
}
func (m *AuthorizationManager) getEffectiveRights(roles []RoleModel) RoleModel {

	finalModel := RoleModel{}

	for x := range roles {

		finalModel.UserManagment = roles[x].UserManagment || finalModel.UserManagment
		finalModel.RoleManagment = roles[x].RoleManagment || finalModel.RoleManagment
		finalModel.Restart = roles[x].Restart || finalModel.Restart
		finalModel.SetToOK = roles[x].SetToOK || finalModel.SetToOK
		finalModel.AddTicket = roles[x].AddTicket || finalModel.AddTicket
		finalModel.RemoveTicket = roles[x].RemoveTicket || finalModel.RemoveTicket
		finalModel.SetFlag = roles[x].SetFlag || finalModel.SetFlag
		finalModel.Confirm = roles[x].Confirm || finalModel.Confirm
		finalModel.Order = roles[x].Order || finalModel.Order
		finalModel.Force = roles[x].Force || finalModel.Force
		finalModel.Definition = roles[x].Definition || finalModel.Definition
	}

	return finalModel
}

type userAuthenticationManager struct {
	col collection.DataCollection
}

func NewAuthenticationManager(collectionName string, provider *datastore.Provider) (AuthenticationManager, error) {

	var col collection.DataCollection
	var err error

	if col, err = provider.GetCollection(collectionName); err != nil {
		return nil, err
	}

	return &userAuthenticationManager{col: col}, nil

}
func (m *userAuthenticationManager) Authenticate(ctx context.Context, username string, password string) (bool, error) {

	if ctx == nil {
		ctx = context.Background()
	}

	dsuser := dsUserModel{}
	if err := m.col.Get(ctx, idFormatter(userNamespace, username), &dsuser); err != nil {
		return false, errors.New("user not found")
	}

	if !dsuser.Enabled {
		return false, errors.New("user account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dsuser.Password), []byte(password)); err != nil {
		return false, err
	}

	return true, nil
}
