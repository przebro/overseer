package ovscli

import (
	"context"
	"fmt"
	"goscheduler/proto/services"
	"io"
	"strings"

	"google.golang.org/grpc"
)

var state []string = []string{"", "Waiting", "Starting", "Executing", "Ended OK", "Ended NOT OK", "Hold"}

type ListResult struct {
	ID     string
	Group  string
	Name   string
	Status string
	Info   []string
}

type OverseerClient struct {
	connectionString string
	conn             *grpc.ClientConn
}

func NewClient(conn string) (*OverseerClient, error) {
	var err error

	var opt []grpc.DialOption
	opt = append(opt, grpc.WithInsecure())
	cli := &OverseerClient{connectionString: conn}
	cli.conn, err = grpc.Dial(cli.connectionString, opt...)
	if err != nil {
		return nil, err
	}
	fmt.Println("connection initialized")
	return cli, nil
}

func (cli *OverseerClient) AddTicket(name, odate string) string {

	service := services.NewResourceServiceClient(cli.conn)

	result, err := service.AddTicket(context.Background(), &services.TicketActionMsg{Name: name, Odate: odate})
	if err != nil {
		return err.Error()
	}

	return result.GetMessage()
}
func (cli *OverseerClient) DelTicket(name, odate string) string {

	service := services.NewResourceServiceClient(cli.conn)
	result, err := service.DeleteTicket(context.Background(), &services.TicketActionMsg{Name: name, Odate: odate})
	if err != nil {
		return err.Error()
	}
	return result.GetMessage()
}
func (cli *OverseerClient) CheckTicket(name, odate string) string {

	service := services.NewResourceServiceClient(cli.conn)
	result, err := service.CheckTicket(context.Background(), &services.TicketActionMsg{Name: name, Odate: odate})
	if err != nil {
		return err.Error()
	}
	return result.GetMessage()
}
func (cli *OverseerClient) SetFlag(name string, value int32) string {

	service := services.NewResourceServiceClient(cli.conn)
	result, err := service.SetFlag(context.Background(), &services.FlagActionMsg{Name: name, State: value})
	if err != nil {
		return err.Error()
	}
	return result.GetMessage()
}
func (cli *OverseerClient) CheckFlag(name string, value int32) string {

	return ""
}

func (cli *OverseerClient) OrderTask(group, name, odate string, force bool) string {

	var result *services.ActionResultMsg
	var err error
	service := services.NewTaskServiceClient(cli.conn)

	if force {
		result, err = service.ForceTask(context.Background(), &services.TaskOrderMsg{TaskGroup: group, TaskName: name, Odate: odate})
	} else {
		result, err = service.OrderTask(context.Background(), &services.TaskOrderMsg{TaskGroup: group, TaskName: name, Odate: odate})
	}

	if err != nil {
		return err.Error()
	}

	return result.GetMessage()
}
func (cli *OverseerClient) ListTasks() []ListResult {

	var err error
	res := make([]ListResult, 0)
	service := services.NewTaskServiceClient(cli.conn)
	stream, err := service.ListTasks(context.Background(), &services.TaskFilterMsg{})

	if err != nil {
		return nil
	}
	for {
		result, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
		}
		lr := ListResult{
			Group:  result.GroupName,
			Name:   result.TaskName,
			ID:     result.TaskId,
			Status: state[result.TaskStatus],
			Info:   strings.Split(result.Waiting, ";"),
		}
		res = append(res, lr)
	}
	return res
}
func (cli *OverseerClient) TaskDetail(orderID string) []string {
	service := services.NewTaskServiceClient(cli.conn)
	result, err := service.TaskDetail(context.Background(), &services.TaskActionMsg{TaskID: orderID})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Group:", result.BaseData.GroupName, "Name:", result.BaseData.TaskName)
	fmt.Println("Task ID:", result.BaseData.TaskId)
	fmt.Println("Order Date:", result.BaseData.OrderDate)
	fmt.Println("Status:", state[result.BaseData.TaskStatus])
	fmt.Println("Start Time:", result.StartTime)
	fmt.Println("End Time:", result.EndTime)
	fmt.Println("Output:")
	for _, x := range result.Output {
		fmt.Println(x)
	}

	return []string{}

}

func (cli *OverseerClient) SetToOK(orderID string) string {
	service := services.NewTaskServiceClient(cli.conn)
	result, err := service.SetToOk(context.Background(), &services.TaskActionMsg{TaskID: orderID})
	if err != nil {

	}

	return result.Message
}

func (cli *OverseerClient) Rerun(orderID string) string {
	service := services.NewTaskServiceClient(cli.conn)
	result, err := service.RerunTask(context.Background(), &services.TaskActionMsg{TaskID: orderID})
	if err != nil {
		return err.Error()
	}

	return result.Message
}
