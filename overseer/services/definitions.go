package services

import (
	"context"
	"encoding/json"
	"fmt"
	"overseer/common/logger"
	"overseer/common/validator"
	"overseer/overseer/auth"
	"overseer/overseer/internal/taskdef"
	"overseer/overseer/taskdata"
	"overseer/proto/services"
	"strings"
)

type ovsDefinitionService struct {
	log        logger.AppLogger
	defManager taskdef.TaskDefinitionManager
}

//NewDefinistionService - Creates a new Definition service
func NewDefinistionService(dm taskdef.TaskDefinitionManager) services.DefinitionServiceServer {

	dservice := &ovsDefinitionService{defManager: dm, log: logger.Get()}
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

//LockDefinition - Locks a definition for edition
func (srv *ovsDefinitionService) LockDefinition(ctx context.Context, msg *services.DefinitionActionMsg) (*services.LockResultMsg, error) {

	var success bool = false
	var resultMsg string = ""

	result := &services.LockResultMsg{}
	for _, e := range msg.DefinitionMsg {

		data := taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: e.GroupName}, Name: e.TaskName}
		err := validator.Valid.Validate(data)

		if err != nil {
			result.LockResult = append(result.LockResult, &services.LockResult{Success: false, Message: err.Error()})
			continue
		}

		lockID, err := srv.defManager.Lock(data)

		if err != nil {
			resultMsg = fmt.Sprintf("Unable to acquire lock for task group:%s task:%s", e.GroupName, e.TaskName)
			success = false
			lockID = 0
		} else {

		}
		result.LockResult = append(result.LockResult, &services.LockResult{Success: success, Message: resultMsg, LockID: lockID})
	}

	return result, nil
}

//UnlockDefinition - unlocks definition
func (srv *ovsDefinitionService) UnlockDefinition(ctx context.Context, msg *services.DefinitionActionMsg) (*services.LockResultMsg, error) {

	var success bool
	var rmsg string
	result := &services.LockResultMsg{}
	for _, e := range msg.DefinitionMsg {

		err := srv.defManager.Unlock(e.LockID)

		if err != nil {
			rmsg = fmt.Sprintf("Unable to release lock")
			success = false
		} else {
			rmsg = fmt.Sprintf("Lock released")
			success = true
		}

		result.LockResult = append(result.LockResult, &services.LockResult{LockID: e.LockID, Success: success, Message: rmsg})
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
	groups := srv.defManager.GetGroups()

	for _, grp := range groups {
		result.GroupName = append(result.GroupName, grp)
	}

	return result, nil
}
func (srv *ovsDefinitionService) ListDefinitionsFromGroup(ctx context.Context, filter *services.FilterMsg) (*services.DefinitionListResultMsg, error) {

	result := &services.DefinitionListResultMsg{Definitions: []*services.DefinitionListMsg{}}

	if err := validator.Valid.ValidateTag(filter.Filter, "resname"); err != nil {
		return nil, err
	}

	tasks, err := srv.defManager.GetTaskModelList(filter.Filter)
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
