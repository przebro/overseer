package services

import (
	"context"
	"fmt"

	"github.com/przebro/overseer/common/core"
	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/common/validator"
	"github.com/przebro/overseer/overseer/auth"
	"github.com/przebro/overseer/proto/services"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ovsAdministrationService struct {
	log         logger.AppLogger
	umanager    *auth.UserManager
	rmanager    *auth.RoleManager
	amanager    *auth.RoleAssociationManager
	qcomponents []core.ComponentQuiescer
}

//NewAdministrationService - returns a new instance of ovsAdministrationService
func NewAdministrationService(u *auth.UserManager, r *auth.RoleManager, a *auth.RoleAssociationManager, log logger.AppLogger, q ...core.ComponentQuiescer) services.AdministrationServiceServer {

	return &ovsAdministrationService{umanager: u, rmanager: r, amanager: a, log: log}
}

//CreateUser - Creates a new user
func (srv *ovsAdministrationService) CreateUser(ctx context.Context, msg *services.CreateUserMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	response.Success = false

	if msg.User == nil {
		response.Message = "User data required"
		return response, nil

	}

	model := auth.UserModel{
		Username: msg.User.Username,
		FullName: msg.User.Fullname,
		Mail:     msg.User.Email,
		Enabled:  msg.User.Enabled,
		Password: msg.Password,
	}

	if err := validator.Valid.Validate(model); err != nil {

		response.Message = err.Error()
		return response, nil
	}

	if err := validator.Valid.ValidateTag(msg.Password, "min=8,max=32"); err != nil {
		response.Message = err.Error()
		return response, nil

	}

	if _, ok := srv.umanager.Get(model.Username); ok {
		response.Message = "user already exists"
		return response, nil
	}

	pass, err := auth.HashPassword([]byte(model.Password))
	if err != nil {
		response.Message = err.Error()
		return response, nil
	}

	model.Password = pass

	for _, n := range msg.User.Roles {
		if _, ok := srv.rmanager.Get(n); !ok {
			response.Message = fmt.Sprintf("role %s does not exists", n)
			return response, nil
		}
	}

	if err := srv.umanager.Create(model); err != nil {
		response.Message = err.Error()
		return response, nil
	}

	amodel := auth.RoleAssociationModel{
		UserID: model.Username,
		Roles:  msg.User.Roles,
	}

	if err := srv.amanager.Create(amodel); err != nil {
		response.Message = err.Error()
		return response, nil
	}

	response.Success = true
	response.Message = "user created"

	return response, nil
}

//ModifyUser - Modifies user, this method is available for administrator only, it allows to change password
//without knowledge about the old one
func (srv *ovsAdministrationService) ModifyUser(ctx context.Context, msg *services.CreateUserMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	var umodel auth.UserModel
	var ok bool
	var pass string
	var err error

	if umodel, ok = srv.umanager.Get(msg.User.Username); !ok {
		response.Success = false
		response.Message = "user does not exists"
		return response, nil
	}

	if msg.Password == "" {
		pass = umodel.Password

	} else {
		if pass, err = auth.HashPassword([]byte(msg.Password)); err != nil {
			response.Success = false
			response.Message = err.Error()
			return response, nil
		}
	}

	model := auth.UserModel{
		Username: msg.User.Username,
		FullName: msg.User.Fullname,
		Mail:     msg.User.Email,
		Enabled:  msg.User.Enabled,
		Password: pass,
	}

	if err := validator.Valid.Validate(model); err != nil {

		response.Message = err.Error()
		response.Success = false
		return response, nil
	}

	if err = srv.umanager.Modify(model); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	amodel := auth.RoleAssociationModel{
		UserID: model.Username,
		Roles:  msg.User.Roles,
	}

	if err = srv.amanager.Modify(amodel); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	response.Success = true
	response.Message = "user modified"

	return response, nil
}
func (srv *ovsAdministrationService) DeleteUser(ctx context.Context, msg *services.UserMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	response.Success = false

	if _, ok := srv.umanager.Get(msg.Username); !ok {
		response.Message = "user does not exists"
		return response, nil
	}

	if err := srv.umanager.Delete(msg.Username); err != nil {
		response.Message = err.Error()
		return response, nil

	}

	if err := srv.amanager.Delete(msg.Username); err != nil {
		response.Message = err.Error()
		return response, nil
	}

	response.Success = true
	response.Message = "user deleted"

	return response, nil
}

//ListUsers - returns a List of users
func (srv *ovsAdministrationService) ListUsers(ctx context.Context, msg *services.FilterMsg) (*services.ListEntityResultMsg, error) {

	var umodel []auth.UserModel
	var err error
	if umodel, err = srv.umanager.All(msg.Filter); err != nil {
		return nil, err
	}

	result := &services.ListEntityResultMsg{}
	for _, m := range umodel {

		result.Entity = append(result.Entity, &services.EntityMsg{Name: m.Username, Description: m.Mail})
	}

	return result, nil
}
func (srv *ovsAdministrationService) GetUser(ctx context.Context, msg *services.UserMsg) (*services.UserResultMsg, error) {

	var model auth.UserModel
	var assoc auth.RoleAssociationModel
	var ok bool

	if err := validator.Valid.ValidateTag(msg.Username, "required,max=32"); err != nil {
		return nil, err

	}

	if model, ok = srv.umanager.Get(msg.Username); !ok {

		return nil, fmt.Errorf("user does not exists:%s", msg.Username)
	}

	assoc, _ = srv.amanager.Get(msg.Username)

	user := services.UserAccount{
		Username: model.Username,
		Fullname: model.FullName,
		Enabled:  model.Enabled,
		Email:    model.Mail,
		Roles:    assoc.Roles,
	}
	result := &services.UserResultMsg{User: &user}

	return result, nil
}

//CreateRole - creates a new roles
func (srv *ovsAdministrationService) CreateRole(ctx context.Context, msg *services.RoleDefinitionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	response.Success = false

	if msg.Role == nil || msg.Role.Rolename == "" {
		response.Message = "rolename required"
		return response, nil
	}

	model := auth.RoleModel{
		Name:           msg.Role.Rolename,
		Description:    msg.Description,
		Administration: msg.Administration,
		Restart:        msg.Restart,
		SetToOK:        msg.SetToOk,
		AddTicket:      msg.AddTicket,
		RemoveTicket:   msg.RemoveTicket,
		SetFlag:        msg.SetFlag,
		Confirm:        msg.Confirm,
		Order:          msg.Order,
		Force:          msg.Force,
		Definition:     msg.Definition,
		Bypass:         msg.Bypass,
		Hold:           msg.Hold,
		Free:           msg.Free,
	}

	if err := validator.Valid.Validate(model); err != nil {
		response.Message = err.Error()
		return response, nil
	}

	if _, ok := srv.rmanager.Get(msg.Role.Rolename); ok {
		response.Message = "role already exists"
		return response, nil
	}

	if err := srv.rmanager.Create(model); err != nil {
		response.Message = err.Error()
		return response, nil
	}

	response.Success = true
	response.Message = fmt.Sprintf("role %s created", msg.Role.Rolename)
	return response, nil
}

//ModifyRole -  an existing role
func (srv *ovsAdministrationService) ModifyRole(ctx context.Context, msg *services.RoleDefinitionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	if _, ok := srv.rmanager.Get(msg.Role.Rolename); !ok {
		response.Success = false
		response.Message = "role does not exists"
		return response, nil
	}

	model := auth.RoleModel{
		Name:           msg.Role.Rolename,
		Description:    msg.Description,
		Administration: msg.Administration,
		Restart:        msg.Restart,
		SetToOK:        msg.SetToOk,
		AddTicket:      msg.AddTicket,
		RemoveTicket:   msg.RemoveTicket,
		SetFlag:        msg.SetFlag,
		Confirm:        msg.Confirm,
		Order:          msg.Order,
		Force:          msg.Force,
		Definition:     msg.Definition,
		Bypass:         msg.Bypass,
		Hold:           msg.Hold,
		Free:           msg.Free,
	}

	if err := srv.rmanager.Modify(model); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	response.Success = true
	response.Message = "role modified"

	return response, nil

}

//DeleteRole - Removes a role
func (srv *ovsAdministrationService) DeleteRole(ctx context.Context, msg *services.RoleMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	response.Success = false

	if _, ok := srv.rmanager.Get(msg.Rolename); !ok {
		response.Message = "role does not exists"
		return response, nil
	}

	if err := srv.rmanager.Delete(msg.Rolename); err != nil {
		response.Message = err.Error()
		return response, nil
	}

	response.Success = true
	response.Message = "role deleted"

	return response, nil
}

//ListRoles - returns a list of roles
func (srv *ovsAdministrationService) ListRoles(ctx context.Context, msg *services.FilterMsg) (*services.ListEntityResultMsg, error) {

	var rmodel []auth.RoleModel
	var err error
	if rmodel, err = srv.rmanager.All(msg.Filter); err != nil {
		return nil, err
	}

	result := &services.ListEntityResultMsg{}
	for _, m := range rmodel {

		result.Entity = append(result.Entity, &services.EntityMsg{Name: m.Name, Description: m.Description})
	}

	return result, nil
}

//GetRole - returns a role
func (srv *ovsAdministrationService) GetRole(ctx context.Context, msg *services.RoleMsg) (*services.RoleResultMsg, error) {

	var model auth.RoleModel
	var ok bool

	if err := validator.Valid.ValidateTag(msg.Rolename, "required,max=32"); err != nil {
		return nil, err

	}

	if model, ok = srv.rmanager.Get(msg.Rolename); !ok {
		return nil, fmt.Errorf("role does not exists")
	}

	role := &services.RoleDefinitionMsg{
		Role:           &services.RoleMsg{Rolename: model.Name},
		Description:    model.Description,
		Administration: model.Administration,
		Restart:        model.Restart,
		SetToOk:        model.SetToOK,
		AddTicket:      model.AddTicket,
		RemoveTicket:   model.RemoveTicket,
		SetFlag:        model.SetFlag,
		Confirm:        model.Confirm,
		Order:          model.Order,
		Force:          model.Force,
		Definition:     model.Definition,
		Bypass:         model.Bypass,
		Hold:           model.Hold,
		Free:           model.Free,
	}
	result := &services.RoleResultMsg{Role: role}

	return result, nil
}

func (srv *ovsAdministrationService) Quiesce(ctx context.Context, msg *empty.Empty) (*services.ActionResultMsg, error) {

	for _, q := range srv.qcomponents {
		q.Quiesce()
	}

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (srv *ovsAdministrationService) Resume(ctx context.Context, msg *empty.Empty) (*services.ActionResultMsg, error) {

	for _, q := range srv.qcomponents {
		q.Resume()
	}

	return nil, status.Error(codes.Unimplemented, "not implemented")
}

//GetAllowedAction - returns allowed action for given method. Implementation of handlers.AccessRestricter
func (srv *ovsAdministrationService) GetAllowedAction(method string) auth.UserAction {

	return auth.ActionAdministration
}
