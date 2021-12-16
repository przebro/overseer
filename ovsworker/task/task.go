package task

import (
	"context"
	"overseer/common/types"
	"overseer/ovsworker/jobs"
	"overseer/ovsworker/status"
	"sync"
)

//TaskRunnerManager - executes a commissioned task
type TaskRunnerManager struct {
	store    map[string]status.JobExecutionStatus
	statChan chan status.JobExecutionStatus
	lock     sync.Mutex
}

//NewTaskExecutor - creates a new TaskRunnerManager
func NewTaskExecutor() *TaskRunnerManager {

	exec := &TaskRunnerManager{
		store:    map[string]status.JobExecutionStatus{},
		statChan: make(chan status.JobExecutionStatus),
		lock:     sync.Mutex{},
	}
	exec.updateTaskStatus()

	return exec
}

//RunTask - starts work fragment
func (exec *TaskRunnerManager) RunTask(j jobs.JobExecutor) (status.JobExecutionStatus, int) {

	tasks := 0
	status := status.JobExecutionStatus{
		TaskID:      j.JobTaskID(),
		ExecutionID: j.JobExecutionID(),
		State:       types.WorkerTaskStatusExecuting,
		ReturnCode:  0,
		StatusCode:  0,
	}

	exec.lock.Lock()
	exec.store[j.JobExecutionID()] = status
	tasks = len(exec.store)
	exec.lock.Unlock()

	go func() {
		j.StartJob(context.Background(), exec.statChan)
	}()

	return status, tasks

}

func (exec *TaskRunnerManager) updateTaskStatus() {

	go func() {
		for {
			exec.update(<-exec.statChan)
		}
	}()

}

func (exec *TaskRunnerManager) update(stat status.JobExecutionStatus) {

	defer exec.lock.Unlock()
	exec.lock.Lock()
	exec.store[stat.ExecutionID] = stat
}

//GetTaskStatus - gets fragment status
func (exec *TaskRunnerManager) GetTaskStatus(executionID string) (status.JobExecutionStatus, int, bool) {
	defer exec.lock.Unlock()
	exec.lock.Lock()

	stat, exists := exec.store[executionID]
	return stat, len(exec.store), exists
}

//TaskCount - returns the number of tasks currently processed
func (exec *TaskRunnerManager) TaskCount() int {
	defer exec.lock.Unlock()
	exec.lock.Lock()

	return len(exec.store)
}

//CleanupTask - removes a task
func (exec *TaskRunnerManager) CleanupTask(executionID string) int {

	defer exec.lock.Unlock()
	exec.lock.Lock()

	delete(exec.store, executionID)
	return len(exec.store)
}

//TerminateTask - removes a task
func (exec *TaskRunnerManager) TerminateTask(executionID string) int {

	defer exec.lock.Unlock()
	exec.lock.Lock()

	delete(exec.store, executionID)
	return len(exec.store)
}
