package services

import (
	"context"
	"fmt"
	"goscheduler/common/logger"
	"goscheduler/overseer/internal/date"
	"goscheduler/overseer/internal/resources"
	"goscheduler/proto/services"
)

type ovsResourceService struct {
	resManager resources.ResourceManager
	log        logger.AppLogger
}

//NewResourceService - Creates new service for ResourceManager
func NewResourceService(rm resources.ResourceManager) services.ResourceServiceServer {

	rservice := &ovsResourceService{resManager: rm, log: logger.Get()}
	return rservice
}

func (srv *ovsResourceService) AddTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	var respmsg string
	result, err := srv.resManager.Add(msg.GetName(), date.Odate(msg.GetOdate()))
	srv.log.Info(result)

	if err != nil {
		respmsg = fmt.Sprintf("%s, ticket: %s odate:%s", err, msg.GetName(), msg.GetOdate())
	} else {
		respmsg = fmt.Sprintf("ticket: %s with odate:%s added", msg.GetName(), msg.GetOdate())
	}

	resMsg := &services.ActionResultMsg{Success: result, Message: respmsg}

	return resMsg, nil
}
func (srv *ovsResourceService) DeleteTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	result, _ := srv.resManager.Delete(msg.GetName(), date.Odate(msg.GetOdate()))
	srv.log.Info(result)

	respmsg := fmt.Sprintf("ticket: %s with odate:%s does not exists", msg.GetName(), msg.GetOdate())

	if result == true {
		respmsg = fmt.Sprintf("ticket: %s with odate:%s deleted", msg.GetName(), msg.GetOdate())
	}

	resp := &services.ActionResultMsg{Success: result, Message: respmsg}

	return resp, nil
}
func (srv *ovsResourceService) CheckTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	result := srv.resManager.Check(msg.GetName(), date.Odate(msg.GetOdate()))

	respmsg := fmt.Sprintf("ticket: %s with odate:%s does not exists", msg.GetName(), msg.GetOdate())
	srv.log.Info(result)

	if result == true {
		respmsg = fmt.Sprintf("ticket: %s with odate:%s exists", msg.GetName(), msg.GetOdate())
	}
	resp := &services.ActionResultMsg{Success: result, Message: respmsg}

	return resp, nil
}
func (srv *ovsResourceService) ListTickets(msg *services.TicketActionMsg, lflags services.ResourceService_ListTicketsServer) error {

	data := srv.resManager.ListTickets(msg.Name, msg.Odate)

	for _, d := range data {
		msg := services.TicketListResultMsg{Name: d.Name, Odate: string(d.Odate)}
		lflags.Send(&msg)

	}

	return nil
}
func (srv *ovsResourceService) SetFlag(ctx context.Context, msg *services.FlagActionMsg) (*services.ActionResultMsg, error) {

	result := &services.ActionResultMsg{}
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

	result.Success = ok
	result.Message = fmt.Sprintf("%s has been set to %s", msg.Name, respPolicy)

	return result, nil
}
func (srv *ovsResourceService) ListFlags(msg *services.FlagActionMsg, lflags services.ResourceService_ListFlagsServer) error {

	data := srv.resManager.ListFlags(msg.Name)

	for _, d := range data {
		msg := services.FlagListResultMsg{FlagName: d.Name, State: int32(d.Policy)}
		lflags.Send(&msg)

	}

	return nil
}
