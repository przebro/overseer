package services

import (
	"context"
	"fmt"

	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/validator"
	"github.com/przebro/overseer/overseer/auth"

	"strings"

	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/overseer/internal/events"
	"github.com/przebro/overseer/overseer/internal/journal"
	"github.com/przebro/overseer/overseer/internal/pool"
	"github.com/przebro/overseer/overseer/taskdata"
	"github.com/przebro/overseer/proto/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PoolManager interface {
	OrderGroup(groupdata taskdata.GroupData, odate date.Odate, username string) ([]string, error)
	Order(task taskdata.GroupNameData, odate date.Odate, username string) (string, error)
	Force(task taskdata.GroupNameData, odate date.Odate, username string) (string, error)
	Enforce(id unique.TaskOrderID, username string) (string, error)
	Rerun(id unique.TaskOrderID, username string) (string, error)
	Hold(id unique.TaskOrderID, username string) (string, error)
	Free(id unique.TaskOrderID, username string) (string, error)
	SetOk(id unique.TaskOrderID, username string) (string, error)
	Confirm(id unique.TaskOrderID, username string) (string, error)
}

type TaskViewer interface {
	Detail(unique.TaskOrderID) (events.TaskDetailResultMsg, error)
	List(filter string) []events.TaskInfoResultMsg
}

type ovsActiveTaskService struct {
	manager  PoolManager
	poolView TaskViewer
	jrnal    journal.TaskLogReader
	services.UnimplementedTaskServiceServer
}

const errInvalidUser = "invalid user"

// NewTaskService - New task service
func NewTaskService(m PoolManager, p TaskViewer, j journal.TaskJournal) services.TaskServiceServer {

	tservice := &ovsActiveTaskService{manager: m, poolView: p, jrnal: j}

	return tservice
}

func (srv *ovsActiveTaskService) OrderGroup(ctx context.Context, in *services.TaskOrderGroupMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	odate := date.Odate(in.Odate)

	if err := validator.Valid.Validate(odate); err != nil {
		response.Success = false
		response.Message = err.Error()

		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	if odate == date.OdateNone {
		odate = date.CurrentOdate()
	}

	data := taskdata.GroupData{Group: in.TaskGroup}

	if err := validator.Valid.Validate(data); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	var username string
	var ok bool
	if username, ok = ctx.Value("username").(string); !ok {
		return response, status.Error(codes.Unauthenticated, errInvalidUser)
	}

	result, err := srv.manager.OrderGroup(data, odate, username)
	if err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, status.Error(codes.Internal, err.Error())
	}

	response.Message = fmt.Sprintf("TaskID:%s", result)
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) ForceGroup(ctx context.Context, in *services.TaskOrderGroupMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	odate := date.Odate(in.Odate)

	if err := validator.Valid.Validate(odate); err != nil {
		response.Success = false
		response.Message = err.Error()
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	if odate == date.OdateNone {
		odate = date.CurrentOdate()
	}

	data := taskdata.GroupData{Group: in.TaskGroup}

	if err := validator.Valid.Validate(data); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	var username string
	var ok bool
	if username, ok = ctx.Value("username").(string); !ok {
		return response, status.Error(codes.Unauthenticated, errInvalidUser)
	}

	result, err := srv.manager.OrderGroup(data, odate, username)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	response.Message = fmt.Sprintf("TaskID:%s", result)
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) OrderTask(ctx context.Context, in *services.TaskOrderMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	odate := date.Odate(in.Odate)

	if err := validator.Valid.Validate(odate); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	if odate == date.OdateNone {
		odate = date.CurrentOdate()
	}

	data := taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: in.TaskGroup}, Name: in.TaskName}

	if err := validator.Valid.Validate(data); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	var username string
	var ok bool

	if username, ok = ctx.Value("username").(string); !ok {
		return response, status.Error(codes.Unauthenticated, errInvalidUser)
	}

	result, err := srv.manager.Order(data, odate, username)
	if err != nil {
		fmt.Println("TASK NOT ORDERED SUCCESSFULY")
		return response, status.Error(codes.Internal, err.Error())
	}

	fmt.Println("TASK ORDERED SUCCESSFULY")

	response.Message = fmt.Sprintf("TaskID:%s", result)
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) ForceTask(ctx context.Context, in *services.TaskOrderMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	odate := date.Odate(in.Odate)
	if err := validator.Valid.Validate(odate); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	if odate == date.OdateNone {
		odate = date.CurrentOdate()
	}

	data := taskdata.GroupNameData{GroupData: taskdata.GroupData{Group: in.TaskGroup}, Name: in.TaskName}

	if err := validator.Valid.Validate(data); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	var username string
	var ok bool
	if username, ok = ctx.Value("username").(string); !ok {
		return response, status.Error(codes.Unauthenticated, errInvalidUser)
	}

	result, err := srv.manager.Force(data, odate, username)
	if err != nil {
		return response, status.Error(codes.Internal, err.Error())
	}

	response.Message = fmt.Sprintf("TaskID:%s", result)
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) RerunTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	var username string
	var ok bool
	if username, ok = ctx.Value("username").(string); !ok {
		return response, status.Error(codes.Unauthenticated, errInvalidUser)
	}

	result, err := srv.manager.Rerun(orderID, username)
	if err != nil {
		return response, setErrorResponse(result, err)
	}

	response.Message = result
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) EnforceTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	var username string
	var ok bool
	if username, ok = ctx.Value("username").(string); !ok {
		return response, status.Error(codes.Unauthenticated, errInvalidUser)
	}

	result, err := srv.manager.Enforce(orderID, username)
	if err != nil {
		return response, setErrorResponse(result, err)
	}

	response.Message = result
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) HoldTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	var username string
	var ok bool
	if username, ok = ctx.Value("username").(string); !ok {
		return response, status.Error(codes.Unauthenticated, errInvalidUser)
	}

	result, err := srv.manager.Hold(orderID, username)
	if err != nil {
		return response, setErrorResponse(result, err)
	}

	response.Success = true
	response.Message = result

	return response, nil
}
func (srv *ovsActiveTaskService) FreeTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	var username string
	var ok bool
	if username, ok = ctx.Value("username").(string); !ok {
		return response, status.Error(codes.Unauthenticated, errInvalidUser)
	}

	result, err := srv.manager.Free(orderID, username)

	if err != nil {
		return response, setErrorResponse(result, err)
	}

	response.Success = true
	response.Message = result

	return response, nil
}
func (srv *ovsActiveTaskService) SetToOk(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}

	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	var username string
	var ok bool
	if username, ok = ctx.Value("username").(string); !ok {
		return response, status.Error(codes.Unauthenticated, errInvalidUser)
	}

	result, err := srv.manager.SetOk(orderID, username)
	if err != nil {
		return response, setErrorResponse(result, err)
	}

	response.Message = result
	response.Success = true

	return response, nil
}

func (srv *ovsActiveTaskService) ConfirmTask(ctx context.Context, in *services.TaskActionMsg) (*services.ActionResultMsg, error) {

	response := &services.ActionResultMsg{}
	orderID := unique.TaskOrderID(in.TaskID)

	if err := validator.Valid.Validate(orderID); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	var username string
	var ok bool
	if username, ok = ctx.Value("username").(string); !ok {
		return response, status.Error(codes.Unauthenticated, errInvalidUser)
	}

	result, err := srv.manager.Confirm(orderID, username)
	if err != nil {
		return response, setErrorResponse(result, err)
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
			Waiting:    "",
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
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	result, err := srv.poolView.Detail(orderID)
	if err != nil {
		return response, setErrorResponse(err.Error(), err)
	}

	response.BaseData = &services.TaskListResultMsg{
		GroupName:  result.Group,
		OrderDate:  string(result.Odate),
		TaskId:     string(result.TaskID),
		TaskName:   result.Name,
		TaskStatus: result.State,
		Waiting:    "",
		Confirmed:  result.Confirmed,
		Held:       result.Held,
		RunNumber:  result.RunNumber,
	}

	response.CyclicData = &services.TaskCyclicResultMsg{
		IsCyclic: result.IsCyclic,
		NextRun:  result.NextRun.String(),
		MaxRun:   int32(result.MaxRun),
		RunFrom:  result.RunFrom,
		Interval: int32(result.RunInterval),
	}

	response.Description = result.Description
	response.EndTime = result.EndTime
	response.StartTime = result.StartTime
	response.Worker = result.Worker
	response.From = result.From
	response.To = result.To

	response.Resources = []*services.TaskResourcesMsg{}
	for _, ticket := range result.Tickets {
		response.Resources = append(response.Resources, &services.TaskResourcesMsg{
			Type:      "ticket",
			Name:      ticket.Name,
			Odate:     string(ticket.Odate),
			Satisfied: ticket.Fulfilled,
		})

	}

	response.Result = &services.ActionResultMsg{Success: true, Message: ""}

	return response, nil
}

func (srv *ovsActiveTaskService) TaskOutput(ctx context.Context, in *services.TaskActionMsg) (*services.TaskDataMsg, error) {

	response := &services.TaskDataMsg{Output: []string{}}

	orderID := unique.TaskOrderID(in.TaskID)
	if err := validator.Valid.Validate(orderID); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	return response, nil

}
func (srv *ovsActiveTaskService) TaskLog(ctx context.Context, in *services.TaskActionMsg) (*services.TaskDataMsg, error) {

	response := &services.TaskDataMsg{Output: []string{}}

	orderID := unique.TaskOrderID(in.TaskID)
	if err := validator.Valid.Validate(orderID); err != nil {
		return response, status.Error(codes.InvalidArgument, err.Error())
	}

	entries := srv.jrnal.ReadLog(orderID)
	for _, n := range entries {
		response.Output = append(response.Output, fmt.Sprintf("%s:%s", n.Time.Format("2006-01-02 15:04:05.000000"), n.Message))
	}

	return response, nil
}

// GetAllowedAction - returns allowed action for given method. Implementation of handlers.AccessRestricter
func (srv *ovsActiveTaskService) GetAllowedAction(method string) auth.UserAction {

	var action auth.UserAction

	if strings.HasSuffix(method, "ListTasks") || strings.HasSuffix(method, "TaskDetail") ||
		strings.HasSuffix(method, "TaskLog") || strings.HasSuffix(method, "TaskOutput") {
		action = auth.ActionBrowse
	}

	if strings.HasSuffix(method, "OrderTask") {
		action = auth.ActionOrder
	}

	if strings.HasSuffix(method, "ForceTask") {
		action = auth.ActionForce
	}

	if strings.HasSuffix(method, "RerunTask") || strings.HasSuffix(method, "EnforceTask") {
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

func setErrorResponse(msg string, err error) error {

	errmsg := ""
	code := codes.Code(0)

	if err == pool.ErrInvalidStatus {
		errmsg = msg
		code = codes.FailedPrecondition
	} else if err == pool.ErrUnableFindTask {
		errmsg = msg
		code = codes.NotFound
	} else {
		code = codes.Internal
		errmsg = err.Error()
	}

	return status.Error(code, errmsg)
}
