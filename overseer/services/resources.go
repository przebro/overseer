package services

import (
	"context"
	"fmt"
	"overseer/common/logger"
	"overseer/common/validator"
	"overseer/overseer/auth"
	"overseer/overseer/internal/date"
	"overseer/overseer/internal/resources"
	"overseer/proto/services"
)

type ovsResourceService struct {
	resManager resources.ResourceManager
	log        logger.AppLogger
}

//NewResourceService - Creates new service for ResourceManager
func NewResourceService(rm resources.ResourceManager) *ovsResourceService {

	rservice := &ovsResourceService{resManager: rm, log: logger.Get()}
	return rservice
}

func (srv *ovsResourceService) AddTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	var respmsg string
	var success bool

	odate := date.Odate(msg.Odate)

	if err := validator.Valid.Validate(odate); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	name := msg.GetName()

	if err := validator.Valid.ValidateTag(name, "required,max=32"); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	result, err := srv.resManager.Add(name, odate)
	srv.log.Info(result)

	if err != nil {
		respmsg = fmt.Sprintf("%s, ticket: %s odate:%s", err, msg.GetName(), msg.GetOdate())
		success = false
	} else {
		respmsg = fmt.Sprintf("ticket: %s with odate:%s added", msg.GetName(), msg.GetOdate())
		success = true
	}

	response.Success = success
	response.Message = respmsg

	return response, nil
}
func (srv *ovsResourceService) DeleteTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	odate := date.Odate(msg.Odate)

	if err := validator.Valid.Validate(odate); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	name := msg.GetName()

	if err := validator.Valid.ValidateTag(name, "required,max=32"); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	result, _ := srv.resManager.Delete(name, odate)
	srv.log.Info(result)

	respmsg := fmt.Sprintf("ticket: %s with odate:%s does not exists", msg.GetName(), msg.GetOdate())

	if result == true {
		respmsg = fmt.Sprintf("ticket: %s with odate:%s deleted", msg.GetName(), msg.GetOdate())
	}

	response.Success = result
	response.Message = respmsg

	return response, nil
}
func (srv *ovsResourceService) CheckTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	odate := date.Odate(msg.Odate)

	if err := validator.Valid.Validate(odate); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	name := msg.GetName()

	if err := validator.Valid.ValidateTag(name, "required,max=32"); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	result := srv.resManager.Check(name, odate)

	respmsg := fmt.Sprintf("ticket: %s with odate:%s does not exists", msg.GetName(), msg.GetOdate())
	srv.log.Info(result)

	if result == true {
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

	if err := validator.Valid.ValidateTag(name, "max=8"); err != nil {
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

	if err := validator.Valid.ValidateTag(name, "required,max=32"); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil

	}

	respPolicy := func() string {
		if msg.State == int32(resources.FlagPolicyExclusive) {
			return "exclusive"
		}
		return "shared"
	}()

	ok, err := srv.resManager.Set(msg.Name, resources.FlagResourcePolicy(msg.State))
	if err != nil {
		return nil, err
	}

	response.Success = ok
	response.Message = fmt.Sprintf("%s has been set to %s", msg.Name, respPolicy)

	return response, nil
}
func (srv *ovsResourceService) ListFlags(msg *services.FlagActionMsg, lflags services.ResourceService_ListFlagsServer) error {

	err := validator.Valid.ValidateTag(msg.GetName(), "required,max=32")
	if err != nil {
		return nil
	}

	data := srv.resManager.ListFlags(msg.Name)

	for _, d := range data {
		msg := services.FlagListResultMsg{FlagName: d.Name, State: int32(d.Policy)}
		lflags.Send(&msg)

	}

	return nil
}

//GetAllowedAction - returns allowed action for given method. Implementation of handlers.AccessRestricter
func (srv *ovsResourceService) GetAllowedAction(method string) auth.UserAction {

	var action auth.UserAction

	return action
}
