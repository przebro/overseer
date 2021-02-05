package services

import (
	"context"
	"fmt"
	"overseer/common/logger"
	"overseer/common/validator"
	"overseer/overseer/auth"
	"overseer/overseer/internal/date"
	"overseer/overseer/internal/pool"
	"overseer/overseer/internal/unique"
	"overseer/overseer/taskdata"
	"overseer/proto/services"
	"strings"
)

type ovsActiveTaskService struct {
	manager  *pool.ActiveTaskPoolManager
	poolView pool.TaskViewer
	log      logger.AppLogger
}

//NewTaskService - New task service
func NewTaskService(m *pool.ActiveTaskPoolManager, p pool.TaskViewer) *ovsActiveTaskService {

	tservice := &ovsActiveTaskService{manager: m, poolView: p, log: logger.Get()}

	return tservice
}

func (srv *ovsActiveTaskService) OrderTask(ctx context.Context, in *services.TaskOrderMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	odate := date.Odate(in.Odate)

	if err := validator.Valid.Validate(odate); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	data := taskdata.GroupNameData{Group: in.TaskGroup, Name: in.TaskName}

	if err := validator.Valid.Validate(data); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil

	}

	result, err := srv.manager.Order(data, odate)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, err

	}

	response.Message = fmt.Sprintf("TaskID:%s", result)
	response.Success = true

	return response, nil
}
func (srv *ovsActiveTaskService) ForceTask(ctx context.Context, in *services.TaskOrderMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	odate := date.Odate(in.Odate)
	if err := validator.Valid.Validate(odate); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	data := taskdata.GroupNameData{Group: in.TaskGroup, Name: in.TaskName}

	if err := validator.Valid.Validate(data); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, nil

	}

	result, err := srv.manager.Force(data, odate)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, err
	}

	response.Message = fmt.Sprintf("TaskID:%s", result)
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) RerunTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {

		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	result, err := srv.manager.Rerun(orderID)
	if err != nil {

		response.Success = false

		if err == pool.ErrInvalidStatus {
			response.Message = result
		} else {
			response.Message = err.Error()
		}
		return response, nil
	}

	response.Message = result
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) HoldTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {

		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	result, err := srv.manager.Hold(orderID)
	if err != nil {

		response.Success = false

		if err == pool.ErrInvalidStatus {
			response.Message = result
		} else {
			response.Message = err.Error()
		}
		return response, nil
	}

	response.Success = true
	response.Message = result

	return response, nil
}
func (srv *ovsActiveTaskService) FreeTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {

		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	result, err := srv.manager.Free(orderID)

	if err != nil {

		response.Success = false

		if err == pool.ErrInvalidStatus {
			response.Message = result
		} else {
			response.Message = err.Error()
		}
		return response, nil
	}

	response.Success = true
	response.Message = result

	return response, nil
}
func (srv *ovsActiveTaskService) SetToOk(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {

		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	result, err := srv.manager.SetOk(orderID)
	if err != nil {

		response.Success = false

		if err == pool.ErrInvalidStatus {
			response.Message = result
		} else {
			response.Message = err.Error()
		}

		return response, nil
	}

	response.Message = result
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) ConfirmTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {

		response.Success = false
		response.Message = err.Error()
		return response, nil
	}

	result, err := srv.manager.Confirm(orderID)
	if err != nil {

		response.Success = false

		if err == pool.ErrInvalidStatus {
			response.Message = result
		} else {
			response.Message = err.Error()
		}
		return response, nil
	}

	response.Message = result
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) ListTasks(in *services.TaskFilterMsg, ltask services.TaskService_ListTasksServer) error {

	result := srv.poolView.List("")

	for _, r := range result {
		resp := &services.TaskListResultMsg{
			TaskName:   r.Name,
			GroupName:  r.Group,
			TaskId:     string(r.TaskID),
			TaskStatus: r.State,
			OrderDate:  string(r.Odate),
			Waiting:    r.WaitingInfo,
			RunNumber:  r.RunNumber,
			Confirmed:  r.Confirmed,
			Held:       r.Held,
		}
		ltask.Send(resp)
	}
	return nil
}
func (srv *ovsActiveTaskService) TaskDetail(ctx context.Context, in *services.TaskActionMsg) (*services.TaskDetailResultMsg, error) {

	response := &services.TaskDetailResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {

		response.Result = &services.ActionResultMsg{Success: false, Message: err.Error()}
		return response, nil
	}

	result, err := srv.poolView.Detail(orderID)
	if err != nil {
		response.Result = &services.ActionResultMsg{Success: false, Message: err.Error()}
		return response, nil
	}

	response.BaseData = &services.TaskListResultMsg{
		GroupName:  result.Group,
		OrderDate:  string(result.Odate),
		TaskId:     string(result.TaskID),
		TaskName:   result.Name,
		TaskStatus: result.State,
		Waiting:    result.WaitingInfo,
		Confirmed:  result.Confirmed,
		Held:       result.Held,
		RunNumber:  result.RunNumber,
	}

	response.Description = result.Description
	response.EndTime = result.EndTime
	response.StartTime = result.StartTime
	response.Output = result.Output
	response.Worker = result.Worker
	response.Output = result.Output
	response.From = result.From
	response.To = result.To

	//response.Resources = []*services.TaskResourcesMsg{}

	response.Result = &services.ActionResultMsg{Success: true, Message: ""}

	return response, nil
}

//GetAllowedAction - returns allowed action for given method. Implementation of handlers.AccessRestricter
func (srv *ovsActiveTaskService) GetAllowedAction(method string) auth.UserAction {

	var action auth.UserAction

	if strings.HasSuffix(method, "ListTasks") || strings.HasSuffix(method, "TaskDetail") {
		action = auth.ActionBrowse
	}

	if strings.HasSuffix(method, "OrderTask") {
		action = auth.ActionOrder
	}

	if strings.HasSuffix(method, "ForceTask") {
		action = auth.ActionForce
	}

	if strings.HasSuffix(method, "RerunTask") {
		action = auth.ActionRestart
	}

	if strings.HasSuffix(method, "HoldTask") {
		action = auth.ActionHold
	}

	if strings.HasSuffix(method, "FreeTask") {
		action = auth.ActionFree
	}

	if strings.HasSuffix(method, "SetToOk") {
		action = auth.ActionSetToOK
	}

	if strings.HasSuffix(method, "ConfirmTask") {
		action = auth.ActionConfirm
	}

	return action
}
