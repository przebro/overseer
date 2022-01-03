package task

import (
	"context"
	"sync"

	"github.com/przebro/overseer/ovsworker/jobs"
	"github.com/przebro/overseer/ovsworker/status"
)

//TaskRunnerManager - executes a commissioned task
type TaskRunnerManager struct {
	store    map[string]status.JobExecutionStatus
	jobs     map[string]jobs.JobExecutor
	statChan chan status.JobExecutionStatus
	lock     *sync.Mutex
}

//NewTaskRunnerManager - creates a new TaskRunnerManager
func NewTaskRunnerManager() *TaskRunnerManager {

	exec := &TaskRunnerManager{
		store:    map[string]status.JobExecutionStatus{},
		statChan: make(chan status.JobExecutionStatus),
		jobs:     map[string]jobs.JobExecutor{},
		lock:     &sync.Mutex{},
	}
	exec.updateTaskStatus()

	return exec
}

//RunTask - starts work fragment
func (exec *TaskRunnerManager) RunTask(j jobs.JobExecutor) (status.JobExecutionStatus, int) {

	tasks := 0
	exec.lock.Lock()
	defer exec.lock.Unlock()

	status := j.StartJob(context.Background(), exec.statChan)

	exec.store[j.JobExecutionID()] = status
	exec.jobs[j.JobExecutionID()] = j
	tasks = len(exec.store)

	return status, tasks
}

func (exec *TaskRunnerManager) updateTaskStatus() {

	go func() {
		for {
			value, ok := <-exec.statChan
			if ok {
				exec.update(value)
			} else {
				return
			}

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
	delete(exec.jobs, executionID)
	delete(exec.store, executionID)
	return len(exec.store)
}

//TerminateTask - removes a task
func (exec *TaskRunnerManager) TerminateTask(executionID string) int {

	defer exec.lock.Unlock()
	exec.lock.Lock()
	exec.jobs[executionID].CancelJob()

	return len(exec.store)
}
