package services

import (
	"context"
	"fmt"
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

var exservice *workerExecutionService = &workerExecutionService{
	log: logger.NewTestLogger(),
	te:  task.NewTaskExecutor(),
}

func TestCreateInstance(t *testing.T) {
	inst := NewWorkerExecutionService()
	if inst == nil {
		t.Error("create instance")
	}
}

func TestStartTaskDummy(t *testing.T) {

	cmd, _ := anypb.New(&actions.DummyTaskAction{Data: "testdata"})
	msg := &wservices.StartTaskMsg{
		TaskID:    &wservices.TaskIdMsg{TaskID: "00000"},
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

	response, err = exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00000"})
	if err != nil {
		t.Error(err)
	}

	exservice.CompleteTask(context.Background(), &wservices.TaskIdMsg{TaskID: "00000"})

	_, err = exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00000"})
	if err == nil || status.Code(err) != codes.NotFound {
		t.Error(err)
	}

}
func TestStartTaskOS(t *testing.T) {

	cmd, _ := anypb.New(&actions.OsTaskAction{Type: "command", CommandLine: "ls -l"})
	msg := &wservices.StartTaskMsg{
		TaskID:    &wservices.TaskIdMsg{TaskID: "00010"},
		Type:      "os",
		Variables: map[string]string{},
		Command:   cmd,
	}
	response, err := exservice.StartTask(context.Background(), msg)
	if err != nil {
		t.Error(err)
	}

	time.Sleep(1 * time.Second)
	response, err = exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00010"})
	if err != nil {
		t.Error(err)
	}

	fmt.Println(response)

	exservice.te.CleanupTask("00010")
}

func TestStartTaskInvalid(t *testing.T) {

	cmd, _ := anypb.New(&actions.DummyTaskAction{Data: "testdata"})
	msg := &wservices.StartTaskMsg{
		TaskID:    &wservices.TaskIdMsg{TaskID: "11111"},
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
		TaskID:    &wservices.TaskIdMsg{TaskID: "00020"},
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
