package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/przebro/overseer/common/validator"
	"github.com/przebro/overseer/overseer/auth"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/przebro/overseer/overseer/taskdata"
	"github.com/przebro/overseer/proto/services"
	"github.com/rs/zerolog"
)

type ovsDefinitionService struct {
	defManager taskdef.TaskDefinitionManager
	services.UnimplementedDefinitionServiceServer
}

// NewDefinistionService - Creates a new Definition service
func NewDefinistionService(dm taskdef.TaskDefinitionManager) services.DefinitionServiceServer {

	dservice := &ovsDefinitionService{defManager: dm}
	return dservice
}

// GetDefinition - Gets definition
func (srv *ovsDefinitionService) GetDefinition(ctx context.Context, msg *services.DefinitionActionMsg) (*services.DefinitionResultMsg, error) {

	log := zerolog.Ctx(ctx).With().Str("service", "definition").Logger()
	var success bool
	var resultMsg string
	tdata := make([]taskdata.GroupNameData, 1)
	result := &services.DefinitionResultMsg{}

	log.Info().Msg("get definition")

	tdata[0] = taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: msg.DefinitionMsg.GroupName}, Name: msg.DefinitionMsg.TaskName}

	tasks := srv.defManager.GetTasks(tdata...)

	for _, t := range tasks {

		data, err := json.Marshal(t)

		if err != nil {
			success = false
			n, grp := t.Definition.Name, t.Definition.Group
			resultMsg = fmt.Sprintf("unable to parse definition group:%s name:%s", n, grp)
		} else {
			success = true
			resultMsg = string(data)
		}

		result.DefinitionMsg = &services.DefinitionDetails{Success: success, Message: resultMsg}
	}

	return result, nil
}

// ListGroups - List definition groups
func (srv *ovsDefinitionService) ListGroups(ctx context.Context, filter *services.FilterMsg) (*services.DefinitionListGroupResultMsg, error) {

	log := zerolog.Ctx(ctx).With().Str("service", "definition").Logger()
	result := &services.DefinitionListGroupResultMsg{GroupName: []string{}}

	if err := validator.Valid.ValidateTag(filter.Filter, "resname"); filter.Filter != "" && err != nil {
		return nil, err
	}

	log.Info().Msg("Request for group")
	groups, _ := srv.defManager.GetGroups()

	result.GroupName = append(result.GroupName, groups...)

	return result, nil
}
func (srv *ovsDefinitionService) ListDefinitionsFromGroup(ctx context.Context, groupmsg *services.GroupNameMsg) (*services.DefinitionListResultMsg, error) {

	result := &services.DefinitionListResultMsg{Definitions: []*services.DefinitionListMsg{}}

	if err := validator.Valid.ValidateTag(groupmsg.GroupName, "resname"); err != nil {
		return nil, err
	}

	tdata := taskdata.GroupData{Group: groupmsg.GroupName}
	if err := validator.Valid.Validate(tdata); err != nil {
		return nil, err
	}

	tasks, err := srv.defManager.GetTaskModelList(tdata)
	if err != nil {
		return nil, err
	}
	for _, r := range tasks {

		msg := &services.DefinitionListMsg{}
		msg.GroupName = r.Group
		msg.TaskName = r.Name
		msg.Description = r.Description
		result.Definitions = append(result.Definitions, msg)
	}

	return result, nil
}

// GetAllowedAction - returns allowed action for given method. Implementation of handlers.AccessRestricter
func (srv *ovsDefinitionService) GetAllowedAction(method string) auth.UserAction {

	var action auth.UserAction

	if strings.HasSuffix(method, "ListDefinitionsFromGroup") ||
		strings.HasSuffix(method, "ListGroups") ||
		strings.HasSuffix(method, "GetDefinition") {

		action = auth.ActionBrowse
	}

	if strings.HasSuffix(method, "LockDefinition") || strings.HasSuffix(method, "UnlockDefinition") {
		action = auth.ActionDefinition

	}

	return action
}
