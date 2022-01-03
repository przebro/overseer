package ovscli

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/przebro/overseer/common/cert"
	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/proto/services"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var state []string = []string{"", "Waiting", "Starting", "Executing", "Ended OK", "Ended NOT OK", "Hold"}

//ListResult - contains base information about the task
type ListResult struct {
	ID     string
	Group  string
	Name   string
	Status string
	Info   []string
}

//TicketValue - represents ticket value
type TicketValue struct {
	Name, Odate string
}

//OverseerClient - holds connection to ovs server
type OverseerClient struct {
	conn  *grpc.ClientConn
	token string
}

//CreateClient - creates a new instance of OverseerClient
func CreateClient() *OverseerClient {

	return &OverseerClient{}
}

//Close - closes current connection
func (cli *OverseerClient) Close() {
	if cli.conn != nil {
		cli.conn.Close()
	}
}

//Connect - setup  connection to server
func (cli *OverseerClient) Connect(addr string, serverCA, clientCertPath, clientKeyPath string, level types.ConnectionSecurityLevel, policy types.CertPolicy) string {

	var opt []grpc.DialOption
	var err error

	dctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opt = append(opt, grpc.WithBlock())
	if level == types.ConnectionSecurityLevelNone {
		opt = append(opt, grpc.WithInsecure())
	} else {
		creds, err := cert.BuildClientCredentials(serverCA, clientCertPath, clientKeyPath, policy, level)
		if err != nil {
			return fmt.Sprintf("failed to initialize connection:%v\n", err)
		}
		opt = append(opt, creds)
	}

	if cli.conn != nil {
		cli.token = ""
	}

	cli.conn, err = grpc.DialContext(dctx, addr, opt...)
	if err != nil {
		return fmt.Sprintf("failed to initialize connection:%v\n", err)
	}

	return fmt.Sprintf("connected to:%s", addr)

}

//Authenticate - authenticates user against server
func (cli *OverseerClient) Authenticate(username, passwd string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}
	service := services.NewAuthenticateServiceClient(cli.conn)

	r, err := service.Authenticate(context.Background(), &services.AuthorizeActionMsg{Username: username, Password: passwd})
	if err != nil {
		return "", err
	}

	if r.Success {
		cli.token = r.Message
		return "client authenticated", nil
	}

	return "", fmt.Errorf(r.Message)
}

//AddTicket - adds a ticket with given name and odate
func (cli *OverseerClient) AddTicket(name, odate string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	service := services.NewResourceServiceClient(cli.conn)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	result, err := service.AddTicket(ctx, &services.TicketActionMsg{Name: name, Odate: odate})
	if err != nil {
		return "", err
	}

	return result.Message, nil
}

//DelTicket - removes a ticket
func (cli *OverseerClient) DelTicket(name, odate string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", cli.token)
	service := services.NewResourceServiceClient(cli.conn)

	result, err := service.DeleteTicket(ctx, &services.TicketActionMsg{Name: name, Odate: odate})
	if err != nil {
		return "", err
	}
	return result.Message, nil
}

//CheckTicket - checks if given ticket exists
func (cli *OverseerClient) CheckTicket(name, odate string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	service := services.NewResourceServiceClient(cli.conn)
	result, err := service.CheckTicket(ctx, &services.TicketActionMsg{Name: name, Odate: odate})
	if err != nil {
		return "", err
	}
	return result.Message, nil
}

//ListTickets - returns a list of tickets
func (cli *OverseerClient) ListTickets(name, odate string) ([]TicketValue, error) {

	var result = []TicketValue{}
	if cli.conn == nil {
		return []TicketValue{}, fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	service := services.NewResourceServiceClient(cli.conn)
	r, err := service.ListTickets(ctx, &services.TicketActionMsg{Name: name, Odate: odate})
	if err != nil {
		return []TicketValue{}, err
	}

	for {
		t, err := r.Recv()
		if err != nil && err != io.EOF {
			return []TicketValue{}, err
		}
		if err == io.EOF {
			break
		}
		result = append(result, TicketValue{Name: t.Name, Odate: t.Odate})

	}

	return result, nil
}

//SetFlag - sets a new flag or change current flag
func (cli *OverseerClient) SetFlag(name string, value int32) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	service := services.NewResourceServiceClient(cli.conn)
	result, err := service.SetFlag(ctx, &services.FlagActionMsg{Name: name, State: value})
	if err != nil {
		return "", err
	}
	return result.Message, nil
}

//DestroyFlag - destroys flag resource
func (cli *OverseerClient) DestroyFlag(name string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	service := services.NewResourceServiceClient(cli.conn)
	result, err := service.DestroyFlag(ctx, &services.FlagActionMsg{Name: name})
	if err != nil {
		return "", err
	}

	return result.Message, nil
}

//OrderTask - orders a task definition to active task pool
func (cli *OverseerClient) OrderTask(group, name, odate string, force bool) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	var result *services.ActionResultMsg
	var msg string
	var err error
	service := services.NewTaskServiceClient(cli.conn)

	if force {
		result, err = service.ForceTask(ctx, &services.TaskOrderMsg{TaskGroup: group, TaskName: name, Odate: odate})
	} else {
		result, err = service.OrderTask(ctx, &services.TaskOrderMsg{TaskGroup: group, TaskName: name, Odate: odate})
	}

	if err == nil {
		msg = result.Message
	}

	return msg, err
}

//ListTasks - returns a list of tasks in active task pool
func (cli *OverseerClient) ListTasks() ([]ListResult, error) {

	if cli.conn == nil {
		return []ListResult{}, fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	var err error
	var result = []ListResult{}

	service := services.NewTaskServiceClient(cli.conn)
	stream, err := service.ListTasks(ctx, &services.TaskFilterMsg{})

	if err != nil {
		return result, err
	}
	for {
		r, err := stream.Recv()
		if err != nil && err != io.EOF {
			break
		}
		if err == io.EOF {
			break
		}
		lr := ListResult{
			Group:  r.GroupName,
			Name:   r.TaskName,
			ID:     r.TaskId,
			Status: state[r.TaskStatus],
			Info:   strings.Split(r.Waiting, ";"),
		}
		result = append(result, lr)
	}
	return result, nil
}

//TaskDetail - returns a detailed information about task
func (cli *OverseerClient) TaskDetail(orderID string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	service := services.NewTaskServiceClient(cli.conn)
	result, err := service.TaskDetail(ctx, &services.TaskActionMsg{TaskID: orderID})

	if err != nil {
		return "", err
	}

	details := []string{
		fmt.Sprintf("Group:%s", result.BaseData.GroupName),
		fmt.Sprintf("Name:%s", result.BaseData.TaskName),
		fmt.Sprintf("TaskID:%s", result.BaseData.TaskId),
		fmt.Sprintf("Order Date:%s", result.BaseData.OrderDate),
		fmt.Sprintf("Run Number:%d", result.BaseData.RunNumber),
		fmt.Sprintf("Confirmed:%v", result.BaseData.Confirmed),
		fmt.Sprintf("Held:%v", result.BaseData.Held),
		fmt.Sprintf("Status:%s", state[result.BaseData.TaskStatus]),
		fmt.Sprintf("Start Time:%s", result.StartTime),
		fmt.Sprintf("End Time:%s", result.EndTime),
	}

	return strings.Join(details, "\n"), nil

}

//SetToOK - sets a task to status ended ok
func (cli *OverseerClient) SetToOK(orderID string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	service := services.NewTaskServiceClient(cli.conn)
	result, err := service.SetToOk(ctx, &services.TaskActionMsg{TaskID: orderID})
	if err != nil {
		return "", err
	}

	return result.Message, nil
}

//Rerun - reruns an ended task
func (cli *OverseerClient) Rerun(orderID string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	service := services.NewTaskServiceClient(cli.conn)
	result, err := service.RerunTask(ctx, &services.TaskActionMsg{TaskID: orderID})
	if err != nil {
		return "", err
	}

	return result.Message, nil
}

//Confirm - confirms manually a task
func (cli *OverseerClient) Confirm(orderID string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	service := services.NewTaskServiceClient(cli.conn)
	result, err := service.ConfirmTask(ctx, &services.TaskActionMsg{TaskID: orderID})
	if err != nil {
		return "", err
	}

	return result.Message, nil
}

//Hold - holds a task
func (cli *OverseerClient) Hold(orderID string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	service := services.NewTaskServiceClient(cli.conn)
	result, err := service.HoldTask(ctx, &services.TaskActionMsg{TaskID: orderID})
	if err != nil {
		return "", err
	}

	return result.Message, nil
}

//Free - frees a holded task
func (cli *OverseerClient) Free(orderID string) (string, error) {

	if cli.conn == nil {
		return "", fmt.Errorf("client not connected,connect first")
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", cli.token)

	service := services.NewTaskServiceClient(cli.conn)
	result, err := service.HoldTask(ctx, &services.TaskActionMsg{TaskID: orderID})
	if err != nil {
		return "", err
	}

	return result.Message, nil
}
