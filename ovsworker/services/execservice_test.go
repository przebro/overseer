package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"overseer/common/logger"
	"overseer/ovsworker/task"
	"overseer/proto/actions"
	"overseer/proto/wservices"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

var exservice *workerExecutionService
var lg = logger.NewTestLogger()

func init() {

	os.Mkdir("../../data/tests/sysout", os.ModePerm)
	exservice = &workerExecutionService{
		log:       lg,
		te:        task.NewTaskExecutor(),
		sysoutDir: "../../data/tests/sysout",
	}

}
func TestCreateInstance(t *testing.T) {
	inst, _ := NewWorkerExecutionService("../../data/tests/sysout", 0, lg)
	if inst == nil {
		t.Error("create instance")
	}

	_, err := NewWorkerExecutionService("../../data/tests/tasks.json", 0, lg)
	if err == nil {
		t.Error("create instance")
	}

	_, err = NewWorkerExecutionService("../../data/not_exists/sysout", 0, lg)
	if err == nil {
		t.Error("create instance")
	}

	s, _ := os.Getwd()

	inst, _ = NewWorkerExecutionService(s, 0, lg)
	if inst == nil {
		t.Error("create instance:", err)
	}

}

func TestStartTaskDummy(t *testing.T) {

	cmd, _ := anypb.New(&actions.DummyTaskAction{Data: "testdata"})
	msg := &wservices.StartTaskMsg{
		TaskID: &wservices.TaskIdMsg{TaskID: "00000", ExecutionID: "1234567"},

		Type:      "dummy",
		Variables: map[string]string{},
		Command:   cmd,
	}
	response, err := exservice.StartTask(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}

	if response.Status != wservices.TaskExecutionResponseMsg_RECEIVED {
		t.Error(response.Status)
	}

	response, err = exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00000", ExecutionID: "1234567"})
	if err != nil {
		t.Error(err)
	}

	exservice.CompleteTask(context.Background(), &wservices.TaskIdMsg{TaskID: "00000", ExecutionID: "1234567"})

	_, err = exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00000", ExecutionID: "1234567"})
	if err == nil || status.Code(err) != codes.NotFound {
		t.Error(err)
	}

}
func TestStartTaskOS(t *testing.T) {

	cmd, _ := anypb.New(&actions.OsTaskAction{Type: actions.OsTaskAction_command, CommandLine: "ls -l"})
	msg := &wservices.StartTaskMsg{
		TaskID:    &wservices.TaskIdMsg{TaskID: "00010", ExecutionID: "1234555"},
		Type:      "os",
		Variables: map[string]string{},
		Command:   cmd,
	}
	_, err := exservice.StartTask(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(1 * time.Second)
	_, err = exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00010", ExecutionID: "1234555"})
	if err != nil {
		t.Error(err)
	}

	exservice.te.CleanupTask("1234555")
}

func TestStartTask_Aws_Lambda(t *testing.T) {

	t.SkipNow()

	cmd, _ := anypb.New(&actions.AwsTaskAction{
		Type: actions.AwsTaskAction_lambda,
		Connection: &actions.AwsTaskAction_ConnectionData{
			ConnectionData: &actions.AwsConnetionData{
				Region:      "{region}",
				ProfileName: "{profile}",
			},
		},
		Execution: &actions.AwsTaskAction_LambdaExecution{
			LambdaExecution: &actions.AwsLambdaExecution{
				FunctionName: "{function_name}",
				Alias:        "{alias}",
			},
		},
		PayloadSource: &actions.AwsTaskAction_PayloadRaw{
			PayloadRaw: json.RawMessage("{\"variables\":[{\"name\":\"test\",\"value\":\"sample variable value\"}]}"),
		},
	})

	msg := &wservices.StartTaskMsg{
		TaskID: &wservices.TaskIdMsg{TaskID: "00010", ExecutionID: "1234666"},
		Type:   "aws",
		Variables: map[string]string{
			"SOME_CONTEXT_VARIABLE": "AND_THEIR_VALUE",
			"ANOTHER_VARIABLE":      "ANOTHER_VALUE",
			"OVS_SOMETHING":         "SOME_SAMPLE_VALUE",
		},
		Command: cmd,
	}

	_, err := exservice.StartTask(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(3 * time.Second)
	_, err = exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00010", ExecutionID: "1234666"})
	if err != nil {
		t.Error(err)
	}

	exservice.te.CleanupTask("1234666")
}

func TestStartTask_Aws_StepFunction(t *testing.T) {

	t.SkipNow()

	cmd, _ := anypb.New(&actions.AwsTaskAction{
		Type: actions.AwsTaskAction_stepfunc,
		Connection: &actions.AwsTaskAction_ConnectionData{
			ConnectionData: &actions.AwsConnetionData{
				Region:      "{region}",
				ProfileName: "{profile}",
			},
		},
		Execution: &actions.AwsTaskAction_StepFunction{
			StepFunction: &actions.AwsStepFunctionExecution{
				StateMachineARN: "arn:aws:states:{region}:{id}:stateMachine:{state_machine_name}",
				ExecutionName:   "{execution}",
			},
		},
		PayloadSource: &actions.AwsTaskAction_PayloadRaw{
			PayloadRaw: json.RawMessage("{\"variables\":[{\"name\":\"test\",\"value\":\"sample variable value\"}]}"),
		},
	})

	msg := &wservices.StartTaskMsg{
		TaskID: &wservices.TaskIdMsg{TaskID: "00010", ExecutionID: "1234777"},
		Type:   "aws",
		Variables: map[string]string{
			"SOME_CONTEXT_VARIABLE": "AND_THEIR_VALUE",
			"ANOTHER_VARIABLE":      "ANOTHER_VALUE",
			"OVS_SOMETHING":         "SOME_SAMPLE_VALUE",
		},
		Command: cmd,
	}

	_, err := exservice.StartTask(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}

	for {
		res, err := exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00010", ExecutionID: "1234777"})
		if err != nil {
			t.Error(err)
		}
		if res.Status == wservices.TaskExecutionResponseMsg_ENDED {
			break
		}
		fmt.Println(res)
		time.Sleep(1 * time.Second)
	}

	exservice.te.CleanupTask("1234777")
}

func TestStartTaskInvalid(t *testing.T) {

	cmd, _ := anypb.New(&actions.DummyTaskAction{Data: "testdata"})
	msg := &wservices.StartTaskMsg{
		TaskID:    &wservices.TaskIdMsg{TaskID: "11111", ExecutionID: "1234666"},
		Type:      "invalid",
		Variables: map[string]string{},
		Command:   cmd,
	}
	_, err := exservice.StartTask(context.Background(), msg)
	if err == nil || status.Code(err) != codes.Aborted {
		t.Error(err)
	}
	msg.TaskID = &wservices.TaskIdMsg{TaskID: ""}

	_, err = exservice.StartTask(context.Background(), msg)
	if err == nil || status.Code(err) != codes.Aborted {
		t.Error(err)
	}

	msg.TaskID = &wservices.TaskIdMsg{TaskID: "11111", ExecutionID: ""}

	_, err = exservice.StartTask(context.Background(), msg)
	if err == nil || status.Code(err) != codes.Aborted {
		t.Error(err)
	}

	msg.Type = ""

	_, err = exservice.StartTask(context.Background(), msg)
	if err == nil || status.Code(err) != codes.Aborted {
		t.Error(err)
	}
}
func TestWorkerStatus(t *testing.T) {

	status, err := exservice.WorkerStatus(context.Background(), &empty.Empty{})
	if err != nil {
		t.Error(err)
	}

	if status.Tasks != 0 {
		t.Error("status Tasks, invalid number expected 0,actual:", status.Tasks)
	}

	cmd, _ := anypb.New(&actions.DummyTaskAction{Data: "testdata"})
	msg := &wservices.StartTaskMsg{
		TaskID:    &wservices.TaskIdMsg{TaskID: "00020", ExecutionID: "1234777"},
		Type:      "dummy",
		Variables: map[string]string{},
		Command:   cmd,
	}

	_, err = exservice.StartTask(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}

	status, err = exservice.WorkerStatus(context.Background(), &empty.Empty{})
	if err != nil {
		t.Error(err)
	}

	if status.Tasks != 1 {
		t.Error("status Tasks, invalid number, expected 1")
	}

}

func TestTerminateTask(t *testing.T) {
	_, err := exservice.TerminateTask(context.Background(), &wservices.TaskIdMsg{TaskID: "00020", ExecutionID: "1234777"})
	if err == nil {
		t.Error("unexpected result")
	}
}

func TestTaskOutput(t *testing.T) {
	err := exservice.TaskOutput(&wservices.TaskIdMsg{TaskID: "00020", ExecutionID: "1234777"}, nil)
	if err == nil {
		t.Error("unexpected result")
	}
}

func TestSerializeMap(t *testing.T) {

	keyvalues := map[string]string{
		"sample_01":     "sample_value_01",
		"sample_02":     "sample_value_02",
		"OVS_SOMETHING": "SOME_SAMPLE_VALUE",
	}

	custom := struct {
		Custom map[string]string `json:"Custom"`
	}{}

	custom.Custom = keyvalues

	data, _ := json.Marshal(&custom)
	fmt.Println(string(data))
}
