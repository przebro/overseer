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
		response.Message = err.Error()
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
	}

	result, err := srv.manager.Hold(orderID)
	response.Message = result

	if err != nil {
		response.Success = false
	} else {
		response.Success = true
	}

	return response, nil
}
func (srv *ovsActiveTaskService) FreeTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {

		response.Success = false
		response.Message = err.Error()
	}

	result, err := srv.manager.Free(orderID)
	response.Message = result

	if err != nil {
		response.Success = false
	} else {
		response.Success = true
	}

	return response, nil
}
func (srv *ovsActiveTaskService) SetToOk(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {

		response.Success = false
		response.Message = err.Error()
	}

	result, err := srv.manager.SetOk(orderID)
	if err != nil {

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
			Waiting:    r.WaitingInfo,
		}
		ltask.Send(resp)
	}
	return nil
}
func (srv *ovsActiveTaskService) TaskDetail(ctx context.Context, in *services.TaskActionMsg) (*services.TaskDetailResultMsg, error) {

	response := &services.TaskDetailResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {

		response.Success = false
		response.Message = err.Error()
	}

	result, err := srv.poolView.Detail(orderID)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, err
	}

	response.BaseData = &services.TaskListResultMsg{
		GroupName:  result.Group,
		OrderDate:  string(result.Odate),
		TaskId:     string(result.TaskID),
		TaskName:   result.Name,
		TaskStatus: result.State,
		Waiting:    result.WaitingInfo,
	}
	response.Confirm = result.Confirm
	response.EndTime = result.EndTime
	response.StartTime = result.StartTime
	response.Hold = result.Hold
	response.Output = result.Output
	response.RunNumber = result.RunNumber
	response.Worker = result.Worker
	response.Output = result.Output
	response.Success = true

	return response, nil
}

//GetAllowedAction - returns allowed action for given method. Implementation of handlers.AccessRestricter
func (srv *ovsActiveTaskService) GetAllowedAction(method string) auth.UserAction {

	var action auth.UserAction

	return action
}
