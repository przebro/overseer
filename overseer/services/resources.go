package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/validator"
	"github.com/przebro/overseer/overseer/auth"
	"github.com/przebro/overseer/overseer/internal/resources"
	"github.com/przebro/overseer/proto/services"
	"github.com/rs/zerolog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResouceManager interface {
	Add(ctx context.Context, name string, odate date.Odate) (bool, error)
	Delete(ctx context.Context, name string, odate date.Odate) (bool, error)
	Set(ctx context.Context, name string, policy uint8) (bool, error)
	DestroyFlag(ctx context.Context, name string) (bool, error)
	ListResources(ctx context.Context, filter resources.ResourceFilter) []resources.ResourceModel
}

type ovsResourceService struct {
	resManager ResouceManager
	services.UnimplementedResourceServiceServer
}

// NewResourceService - Creates new service for ResourceManager
func NewResourceService(rm ResouceManager) services.ResourceServiceServer {

	rservice := &ovsResourceService{resManager: rm}
	return rservice
}

func (srv *ovsResourceService) AddTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	log := zerolog.Ctx(ctx).With().Str("service", "resources").Logger()
	response := &services.ActionResultMsg{}
	var respmsg string

	odate := date.Odate(msg.Odate)
	name := msg.GetName()

	if err := validateTicketFields(name, odate); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := srv.resManager.Add(ctx, name, odate)
	log.Info().Bool("result", result).Msg("add ticket")

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

	log := zerolog.Ctx(ctx).With().Str("service", "resources").Logger()

	response := &services.ActionResultMsg{}

	odate := date.Odate(msg.Odate)
	name := msg.GetName()

	if err := validateTicketFields(name, odate); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	result, _ := srv.resManager.Delete(ctx, name, odate)
	log.Info().Bool("result", result).Msg("delete ticket")

	respmsg := fmt.Sprintf("ticket: %s with odate:%s deleted", msg.GetName(), msg.GetOdate())

	if !result {
		respmsg = fmt.Sprintf("ticket: %s with odate:%s does not exists", msg.GetName(), msg.GetOdate())
		return response, status.Error(codes.FailedPrecondition, respmsg)
	}

	response.Success = result
	response.Message = respmsg

	return response, nil
}

// :DEPRECATED
func (srv *ovsResourceService) CheckTicket(ctx context.Context, msg *services.TicketActionMsg) (*services.ActionResultMsg, error) {

	// response := &services.ActionResultMsg{}

	// odate := date.Odate(msg.Odate)
	// name := msg.GetName()

	// if err := validateTicketFields(name, odate); err != nil {
	// 	return response, status.Error(codes.InvalidArgument, err.Error())
	// }

	// result := srv.resManager.Check(name, odate)

	// respmsg := fmt.Sprintf("ticket: %s with odate:%s does not exists", msg.GetName(), msg.GetOdate())
	// srv.log.Info(result)

	// if result {
	// 	respmsg = fmt.Sprintf("ticket: %s with odate:%s exists", msg.GetName(), msg.GetOdate())
	// }

	// response.Success = result
	// response.Message = respmsg

	return nil, nil
}

// ::DEPRECATED
func (srv *ovsResourceService) ListTickets(msg *services.TicketActionMsg, lflags services.ResourceService_ListTicketsServer) error {

	// //Both name and odate are strings values used to filter list,
	// name := msg.GetName()

	// if err := validator.Valid.ValidateTag(name, "max=32"); err != nil {
	// 	return err
	// }

	// odateStr := msg.GetOdate()

	// if err := validator.Valid.ValidateTag(odateStr, "max=8"); err != nil {
	// 	return err
	// }

	// data := srv.resManager.ListTickets(name, odateStr)
	// srv.resManager.ListResources()

	// for _, d := range data {
	// 	msg := services.TicketListResultMsg{Name: d.Name, Odate: string(d.Odate)}
	// 	lflags.Send(&msg)
	// }

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

	ok, err := srv.resManager.Set(ctx, msg.Name, uint8(msg.State))
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

	ok, err := srv.resManager.DestroyFlag(ctx, msg.Name)
	if err != nil {
		return response, status.Error(codes.NotFound, err.Error())
	}

	response.Success = ok
	response.Message = fmt.Sprintf("%s has been removed", msg.Name)

	return response, nil
}

// ::DEPRECATED
func (srv *ovsResourceService) ListFlags(msg *services.FlagActionMsg, lflags services.ResourceService_ListFlagsServer) error {

	// name := msg.GetName()

	// err := validateResourceName(name)
	// if err != nil {
	// 	return status.Error(codes.InvalidArgument, err.Error())
	// }

	// data := srv.resManager.ListFlags(name)

	// for _, d := range data {
	// 	msg := services.FlagListResultMsg{FlagName: d.Name, State: int32(d.Policy)}
	// 	lflags.Send(&msg)
	// }

	return nil
}
func (srv *ovsResourceService) ListResources(ctx context.Context, msg *services.ListResourcesMsg) (*services.ListResourcesResultMsg, error) {

	zerolog.Ctx(ctx).Info().Msg("request recieved")

	if msg == nil {
		return nil, status.Error(codes.InvalidArgument, "request body is empty")
	}

	filterProperties := msg.GetProperties()

	var ticketOptions *resources.TicketFilterOptions = nil
	var flagOptions *resources.FlagFilterOptions = nil

	response := &services.ListResourcesResultMsg{}

	restype := msg.GetResourceType()
	zerolog.Ctx(ctx).Info().Str("resource_type", restype).Str("name", msg.Name).Msg("request recieved")

	if restype == "" {
		return nil, status.Error(codes.InvalidArgument, "resource type is empty")
	}

	if restype != "*" && restype != "ticket" && restype != "flag" {
		return nil, status.Error(codes.InvalidArgument, "invalid resource type value")
	}

	if msg.Name != "" {

		err := validateResourceName(msg.Name)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	filter := resources.ResourceFilter{
		Name: msg.Name,
		Type: []resources.ResourceType{},
	}

	if restype == "*" || restype == "ticket" {

		filter.Type = append(filter.Type, resources.ResourceTypeTicket)

		ticketOptions = &resources.TicketFilterOptions{}

		odateStr := ""
		odateFrom := ""
		odateTo := ""

		if filterProperties != nil {
			if v, ok := filterProperties["odate"]; ok {
				odateStr = v
				//validate this
				ticketOptions.Odate = odateStr
			}
			if v, ok := filterProperties["from"]; ok {
				odateFrom = v
				//validate this
				ticketOptions.OrderDateFrom = odateFrom

			}
			if v, ok := filterProperties["to"]; ok {
				odateTo = v
				//validate this
				ticketOptions.OrderDateTo = odateTo
			}
		}
	}

	if restype == "*" || restype == "flag" {
		flagOptions = &resources.FlagFilterOptions{FlagPolicy: []uint8{}}
		filter.Type = append(filter.Type, resources.ResourceTypeFlag)
		if filterProperties != nil {
			if _, ok := filterProperties["shr"]; ok {
				flagOptions.FlagPolicy = append(flagOptions.FlagPolicy, resources.FlagPolicyShared)
			}
			if _, ok := filterProperties["exl"]; ok {
				flagOptions.FlagPolicy = append(flagOptions.FlagPolicy, resources.FlagPolicyShared)
			}
		}

	}

	filter.TicketFilterOptions = ticketOptions
	filter.FlagFilterOptions = flagOptions

	data := srv.resManager.ListResources(ctx, filter)

	for _, item := range data {

		resource := &services.ResourceMsg{}

		rtype := ""
		if item.Type == resources.ResourceTypeTicket {
			rtype = "ticket"
			orderDate := ""
			if item.Value != 0 {
				orderDate = string(date.FromUnix(item.Value))
			}

			props := &services.ResourceMsg_OrderDate{
				OrderDate: orderDate,
			}

			resource.Properties = props

		} else {
			rtype = "flag"
			st, cnt := getFlagStateValue(item.Value)
			props := &services.ResourceMsg_Flag{
				Flag: &services.FlagProperties{
					State: st,
					Count: cnt,
				},
			}

			resource.Properties = props
		}

		resource.Name = item.ID
		resource.ResourceType = rtype

		response.Resources = append(response.Resources, resource)
	}

	return response, nil
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

// GetAllowedAction - returns allowed action for given method. Implementation of handlers.AccessRestricter
func (srv *ovsResourceService) GetAllowedAction(method string) auth.UserAction {

	var action auth.UserAction

	if strings.HasSuffix(method, "AddTicket") {
		action = auth.ActionAddTicket
	}

	if strings.HasSuffix(method, "DeleteTicket") {
		action = auth.ActionRemoveTicket
	}

	if strings.HasSuffix(method, "SetFlag") {
		action = auth.ActionSetFlag
	}

	if strings.HasSuffix(method, "DestroyFlag") {
		action = auth.ActionSetFlag
	}

	if strings.HasSuffix(method, "ListResources") {
		action = auth.ActionBrowse
	}

	return action
}

func getFlagStateValue(p int64) (string, uint32) {
	st, cnt := uint8(p>>32), uint32(p)
	state := func(s uint8) string {
		if s == 0 {
			return "SHR"
		}
		return "EXL"
	}(st)

	return state, cnt
}
