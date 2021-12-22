package task

import "testing"

var exec *TaskRunnerManager

func TestCreateExecutor(t *testing.T) {

	exec = NewTaskRunnerManager()
	if exec == nil {
		t.Error("unexpected result, empty task executor")
	}
}
