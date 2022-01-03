package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/config"

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
	ActionBrowse UserAction = iota
	ActionAdministration
	ActionRestart
	ActionSetToOK
	ActionAddTicket
	ActionRemoveTicket
	ActionSetFlag
	ActionConfirm
	ActionBypass
	ActionHold
	ActionFree
	ActionOrder
	ActionForce
	ActionDefinition
)

var (
	//ErrUnableFindAction - returned when an action is not found
	ErrUnableFindAction = errors.New("unable to find action")
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

//NewAuthorizationManager - creates a new instance of AuthorizationManager
func NewAuthorizationManager(conf config.SecurityConfiguration, provider *datastore.Provider) (*AuthorizationManager, error) {

	var col collection.DataCollection
	var err error

	if col, err = provider.GetCollection(conf.Collection); err != nil {
		return nil, err
	}

	return &AuthorizationManager{col: col}, nil
}

//VerifyAction - verifies if a given user is eligible to perform a specified action
func (m *AuthorizationManager) VerifyAction(ctx context.Context, action UserAction, username string) (bool, error) {

	if ctx == nil {
		ctx = context.Background()
	}
	model := dsRoleAssociationModel{}
	if err := m.col.Get(ctx, idFormatter(assocNamespace, username), &model); err != nil {
		return false, errors.New("unable to get role association for given user")
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
	case ActionAdministration:
		return finalRole.Administration, nil
	case ActionRestart:
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
	case ActionBypass:
		return finalRole.Bypass, nil
	case ActionHold:
		return finalRole.Hold, nil
	case ActionFree:
		return finalRole.Free, nil
		//a virtual action for every user that has enabled account
	case ActionBrowse:
		return true, nil
	}
	return false, ErrUnableFindAction
}

//getEffectiveRights - returns a sum of rights from all roles
func (m *AuthorizationManager) getEffectiveRights(roles []RoleModel) RoleModel {

	finalModel := RoleModel{}

	for x := range roles {

		finalModel.Administration = roles[x].Administration || finalModel.Administration
		finalModel.Restart = roles[x].Restart || finalModel.Restart
		finalModel.SetToOK = roles[x].SetToOK || finalModel.SetToOK
		finalModel.AddTicket = roles[x].AddTicket || finalModel.AddTicket
		finalModel.RemoveTicket = roles[x].RemoveTicket || finalModel.RemoveTicket
		finalModel.SetFlag = roles[x].SetFlag || finalModel.SetFlag
		finalModel.Confirm = roles[x].Confirm || finalModel.Confirm
		finalModel.Order = roles[x].Order || finalModel.Order
		finalModel.Force = roles[x].Force || finalModel.Force
		finalModel.Definition = roles[x].Definition || finalModel.Definition
		finalModel.Bypass = roles[x].Bypass || finalModel.Bypass
		finalModel.Hold = roles[x].Hold || finalModel.Hold
		finalModel.Free = roles[x].Free || finalModel.Free

	}

	return finalModel
}

type userAuthenticationManager struct {
	col collection.DataCollection
}

//NewAuthenticationManager - creates a new instance of AuthenticationManager
func NewAuthenticationManager(collectionName string, provider *datastore.Provider) (AuthenticationManager, error) {

	var col collection.DataCollection
	var err error

	if col, err = provider.GetCollection(collectionName); err != nil {
		return nil, err
	}

	return &userAuthenticationManager{col: col}, nil

}

//Authenticate - authenticates the user
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

//HashPassword - creates a new hash from given password
func HashPassword(password []byte) (string, error) {

	pass, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)

	if err != nil {
		return "", err
	}

	return string(pass), nil
}
