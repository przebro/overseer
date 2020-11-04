package services

import (
	"context"
	"fmt"
	"goscheduler/common/logger"
	"goscheduler/overseer/internal/date"
	"goscheduler/overseer/internal/pool"
	"goscheduler/overseer/internal/unique"
	"goscheduler/proto/services"
)

type ovsActiveTaskService struct {
	manager  *pool.ActiveTaskPoolManager
	poolView pool.TaskViewer
	log      logger.AppLogger
}

//NewTaskService - New task service
func NewTaskService(m *pool.ActiveTaskPoolManager, p pool.TaskViewer) services.TaskServiceServer {

	tservice := &ovsActiveTaskService{manager: m, poolView: p, log: logger.Get()}

	return tservice
}

func (srv *ovsActiveTaskService) OrderTask(ctx context.Context, in *services.TaskOrderMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	result, err := srv.manager.Order(in.TaskGroup, in.TaskName, date.Odate(in.Odate))
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
	result, err := srv.manager.Force(in.TaskGroup, in.TaskName, date.Odate(in.Odate))
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
	result, err := srv.manager.Rerun(unique.TaskOrderID(in.TaskID))
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
	result, err := srv.manager.Hold(unique.TaskOrderID(in.TaskID))
	if err != nil {

	}

	response.Message = result
	response.Success = true

	return response, nil
}
func (srv *ovsActiveTaskService) FreeTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	result, err := srv.manager.Free(unique.TaskOrderID(in.TaskID))
	if err != nil {

	}

	response.Message = result
	response.Success = true

	return response, nil
}
func (srv *ovsActiveTaskService) SetToOk(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	result, err := srv.manager.SetOk(unique.TaskOrderID(in.TaskID))
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

	result, err := srv.poolView.Detail(unique.TaskOrderID(in.TaskID))
	if err != nil {
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

	return response, nil
}
