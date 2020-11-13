package services

import (
	"context"
	"encoding/json"
	"fmt"
	"goscheduler/common/logger"
	"goscheduler/common/validator"
	"goscheduler/overseer/internal/taskdef"
	"goscheduler/overseer/taskdata"
	"goscheduler/proto/services"
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

	var success bool = false
	var resultMsg string = ""
	tdata := make([]taskdata.GroupNameData, 0)
	result := &services.DefinitionResultMsg{}

	for _, e := range msg.DefinitionMsg {

		data := taskdata.GroupNameData{Group: e.GroupName, Name: e.TaskName}
		err := validator.Valid.Validate(data)
		if err != nil {
			result.DefinitionMsg = append(result.DefinitionMsg, &services.DefinitionResult{Success: false, Message: err.Error()})
		} else {
			tdata = append(tdata, data)
		}
	}

	tasks := srv.defManager.GetTasks(tdata...)

	for _, t := range tasks {

		if t.Result == false {
			success = false
			resultMsg = t.Msg.Error()

		} else {
			data, err := json.Marshal(t.Definition)

			if err != nil {
				success = false
				n, grp, _ := t.Definition.GetInfo()
				resultMsg = fmt.Sprintf("unable to parse definition group:%s name:%s", n, grp)
			} else {
				success = true
				resultMsg = string(data)
			}

		}

		result.DefinitionMsg = append(result.DefinitionMsg, &services.DefinitionResult{Success: success, Message: resultMsg})
	}

	return result, nil
}

//LockDefinition - Locks a definition for edition
func (srv *ovsDefinitionService) LockDefinition(ctx context.Context, msg *services.DefinitionActionMsg) (*services.LockResultMsg, error) {

	var success bool = false
	var resultMsg string = ""

	result := &services.LockResultMsg{}
	for _, e := range msg.DefinitionMsg {

		data := taskdata.GroupNameData{Group: e.GroupName, Name: e.TaskName}
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

	var success bool = false
	var rmsg string = ""
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
func (srv *ovsDefinitionService) ListGroups(msg *services.GroupActionMsg, ldef services.DefinitionService_ListGroupsServer) error {

	srv.log.Info("Request for group")
	result := srv.defManager.GetGroups()

	for _, e := range result {

		msg := &services.DefinitionListGroupResultMsg{GroupName: e}
		ldef.Send(msg)
	}

	return nil
}
func (srv *ovsDefinitionService) ListDefinitionsFromGroup(msg *services.DefinitionActionMsg, ldef services.DefinitionService_ListDefinitionsFromGroupServer) error {

	groups := make([]string, 0)
	for _, e := range msg.DefinitionMsg {

		groups = append(groups, e.GroupName)
	}
	result, err := srv.defManager.GetTasksFromGroup(groups)
	if err != nil {
		return err
	}
	for _, r := range result {

		msg := &services.DefinitionListResultMsg{}
		msg.GroupName = r.Group
		msg.TaskName = r.Name

		ldef.Send(msg)
	}

	return nil
}
