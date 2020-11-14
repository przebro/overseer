package services

import (
	"context"
	"goscheduler/common/logger"
	"goscheduler/ovsworker/launcher"
	"goscheduler/proto/actions"
	"goscheduler/proto/wservices"
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/anypb"
)

var l = launcher.NewFragmentLauncher()
var c = launcher.FragmentFactory(l)

var exservice *workerExecutionService = &workerExecutionService{
	log:      logger.NewTestLogger(),
	launcher: l,
	creator:  c,
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
	if response.Started != true && response.Ended != false {
		t.Error("exec response started ended, invalid values", response.Started, response.Ended)
	}

	response, err = exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00000"})
	if err != nil {
		t.Error(err)
	}
	if response.Started != true && response.Ended != true {
		t.Error("status response started ended, invalid values", response.Started, response.Ended)
	}

	_, err = exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00000"})
	if err == nil || err.Error() != "fragment with given id does not exists" {
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
	if response.Started != true && response.Ended != false {
		t.Error("exec response started ended, invalid values", response.Started, response.Ended)
	}
	time.Sleep(1 * time.Second)
	response, err = exservice.TaskStatus(context.Background(), &wservices.TaskIdMsg{TaskID: "00010"})
	if err != nil {
		t.Error(err)
	}
	if response.Started != true && response.Ended != true {
		t.Error("status response started ended, invalid values", response.Started, response.Ended)
	}

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
	if err == nil || err.Error() != "unable to construct fragment" {
		t.Error(err)
	}
	msg.TaskID = &wservices.TaskIdMsg{TaskID: ""}

	_, err = exservice.StartTask(context.Background(), msg)
	if err == nil || err.Error() != "message taskID cannot be empty" {
		t.Error(err)
	}

	msg.Type = ""

	_, err = exservice.StartTask(context.Background(), msg)
	if err == nil || err.Error() != "message type cannot be empty" {
		t.Error(err)
	}
}
func TestWorkerStatus(t *testing.T) {

	status, err := exservice.WorkerStatus(context.Background(), &empty.Empty{})
	if err != nil {
		t.Error(err)
	}

	if status.Tasks != 0 {
		t.Error("status Tasks, invalid number expected 0")
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
