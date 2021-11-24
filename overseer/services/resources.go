package services

import (
	"context"
	"fmt"
	"overseer/common/logger"
	"overseer/common/types/date"
	"overseer/common/validator"
	"overseer/overseer/auth"
	"overseer/overseer/internal/resources"
	"overseer/proto/services"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ovsResourceService struct {
	resManager resources.ResourceManager
	log        logger.AppLogger
}

//NewResourceService - Creates new service for ResourceManager
func NewResourceService(rm resources.ResourceManager, log logger.AppLogger) services.ResourceServiceServer {

	rservice := &ovsResourceService{resManager: rm, log: log}
	return rservice
}

func (srv *ovsResourceService) AddTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	var respmsg string

	odate := date.Odate(msg.Odate)
	name := msg.GetName()

	if err := validateTicketFields(name, odate); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := srv.resManager.Add(name, odate)
	srv.log.Info(result)

	respmsg = fmt.Sprintf("ticket: %s with odate:%s added", msg.GetName(), msg.GetOdate())
	if !result {
		respmsg = fmt.Sprintf("%s, ticket: %s odate:%s", err, msg.GetName(), msg.GetOdate())
		return response, status.Error(codes.FailedPrecondition, respmsg)
	}

	response.Success = result
	response.Message = respmsg

	return response, nil
}
func (srv *ovsResourceService) DeleteTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	odate := date.Odate(msg.Odate)
	name := msg.GetName()

	if err := validateTicketFields(name, odate); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	result, _ := srv.resManager.Delete(name, odate)
	srv.log.Info(result)

	respmsg := fmt.Sprintf("ticket: %s with odate:%s deleted", msg.GetName(), msg.GetOdate())

	if !result {
		respmsg = fmt.Sprintf("ticket: %s with odate:%s does not exists", msg.GetName(), msg.GetOdate())
		return response, status.Error(codes.FailedPrecondition, respmsg)
	}

	response.Success = result
	response.Message = respmsg

	return response, nil
}
func (srv *ovsResourceService) CheckTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	odate := date.Odate(msg.Odate)
	name := msg.GetName()

	if err := validateTicketFields(name, odate); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	result := srv.resManager.Check(name, odate)

	respmsg := fmt.Sprintf("ticket: %s with odate:%s does not exists", msg.GetName(), msg.GetOdate())
	srv.log.Info(result)

	if result {
		respmsg = fmt.Sprintf("ticket: %s with odate:%s exists", msg.GetName(), msg.GetOdate())
	}

	response.Success = result
	response.Message = respmsg

	return response, nil
}
func (srv *ovsResourceService) ListTickets(msg *services.TicketActionMsg, lflags services.ResourceService_ListTicketsServer) error {

	//Both name and odate are strings values used to filter list,
	name := msg.GetName()

	if err := validator.Valid.ValidateTag(name, "max=32"); err != nil {
		return err
	}

	odateStr := msg.GetOdate()

	if err := validator.Valid.ValidateTag(odateStr, "max=8"); err != nil {
		return err
	}

	data := srv.resManager.ListTickets(name, odateStr)

	for _, d := range data {
		msg := services.TicketListResultMsg{Name: d.Name, Odate: string(d.Odate)}
		lflags.Send(&msg)
	}

	return nil
}
func (srv *ovsResourceService) SetFlag(ctx context.Context, msg *services.FlagActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	name := msg.GetName()

	if err := validateResourceName(name); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())

	}
	if msg.State < 0 || msg.State > 1 {
		return response, status.Error(codes.InvalidArgument, "invalid flag state")

	}

	respPolicy := func() string {
		if msg.State == int32(resources.FlagPolicyExclusive) {
			return "exclusive"
		}
		return "shared"
	}()

	ok, err := srv.resManager.Set(msg.Name, resources.FlagResourcePolicy(msg.State))
	if err != nil {
		return response, status.Error(codes.FailedPrecondition, err.Error())
	}

	response.Success = ok
	response.Message = fmt.Sprintf("%s has been set to %s", msg.Name, respPolicy)

	return response, nil
}
func (srv *ovsResourceService) DestroyFlag(ctx context.Context, msg *services.FlagActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	name := msg.GetName()

	if err := validateResourceName(name); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())

	}

	ok, err := srv.resManager.DestroyFlag(msg.Name)
	if err != nil {
		return response, status.Error(codes.NotFound, err.Error())
	}

	response.Success = ok
	response.Message = fmt.Sprintf("%s has been removed", msg.Name)

	return response, nil
}

func (srv *ovsResourceService) ListFlags(msg *services.FlagActionMsg, lflags services.ResourceService_ListFlagsServer) error {

	name := msg.GetName()

	err := validateResourceName(name)
	if err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	data := srv.resManager.ListFlags(name)

	for _, d := range data {
		msg := services.FlagListResultMsg{FlagName: d.Name, State: int32(d.Policy)}
		lflags.Send(&msg)
	}

	return nil
}

func validateTicketFields(name string, odate date.Odate) error {

	if err := validator.Valid.Validate(odate); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	if err := validateResourceName(name); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}

func validateResourceName(name string) error {

	return validator.Valid.ValidateTag(name, "resvalue,required,max=32")

}

//GetAllowedAction - returns allowed action for given method. Implementation of handlers.AccessRestricter
func (srv *ovsResourceService) GetAllowedAction(method string) auth.UserAction {

	var action auth.UserAction

	if strings.HasSuffix(method, "AddTicket") {
		action = auth.ActionAddTicket
	}

	if strings.HasSuffix(method, "DeleteTicket") {
		action = auth.ActionRemoveTicket
	}

	if strings.HasSuffix(method, "CheckTicket") || strings.HasSuffix(method, "ListTickets") {
		action = auth.ActionBrowse
	}

	if strings.HasSuffix(method, "SetFlag") {
		action = auth.ActionSetFlag
	}

	if strings.HasSuffix(method, "DestroyFlag") {
		action = auth.ActionSetFlag
	}

	if strings.HasSuffix(method, "ListFlags") {
		action = auth.ActionBrowse
	}

	return action
}
