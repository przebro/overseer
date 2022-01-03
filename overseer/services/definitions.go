package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/common/validator"
	"github.com/przebro/overseer/overseer/auth"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/przebro/overseer/overseer/taskdata"
	"github.com/przebro/overseer/proto/services"
)

type ovsDefinitionService struct {
	log        logger.AppLogger
	defManager taskdef.TaskDefinitionManager
}

//NewDefinistionService - Creates a new Definition service
func NewDefinistionService(dm taskdef.TaskDefinitionManager, log logger.AppLogger) services.DefinitionServiceServer {

	dservice := &ovsDefinitionService{defManager: dm, log: log}
	return dservice
}

//GetDefinition - Gets definition
func (srv *ovsDefinitionService) GetDefinition(ctx context.Context, msg *services.DefinitionActionMsg) (*services.DefinitionResultMsg, error) {

	var success bool
	var resultMsg string
	tdata := make([]taskdata.GroupNameData, 0)
	result := &services.DefinitionResultMsg{}

	for _, e := range msg.DefinitionMsg {

		data := taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: e.GroupName}, Name: e.TaskName}
		err := validator.Valid.Validate(data)
		if err != nil {
			result.DefinitionMsg = append(result.DefinitionMsg, &services.DefinitionDetails{Success: false, Message: err.Error()})
		} else {
			tdata = append(tdata, data)
		}
	}

	tasks := srv.defManager.GetTasks(tdata...)

	for _, t := range tasks {

		data, err := json.Marshal(t)

		if err != nil {
			success = false
			n, grp, _ := t.GetInfo()
			resultMsg = fmt.Sprintf("unable to parse definition group:%s name:%s", n, grp)
		} else {
			success = true
			resultMsg = string(data)
		}

		result.DefinitionMsg = append(result.DefinitionMsg, &services.DefinitionDetails{Success: success, Message: resultMsg})
	}

	return result, nil
}

//ListGroups - List definition groups
func (srv *ovsDefinitionService) ListGroups(ctx context.Context, filter *services.FilterMsg) (*services.DefinitionListGroupResultMsg, error) {

	result := &services.DefinitionListGroupResultMsg{GroupName: []string{}}

	if err := validator.Valid.ValidateTag(filter.Filter, "resname"); filter.Filter != "" && err != nil {
		return nil, err
	}

	srv.log.Info("Request for group")
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

//GetAllowedAction - returns allowed action for given method. Implementation of handlers.AccessRestricter
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
