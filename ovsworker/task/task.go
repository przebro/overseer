package task

import (
	"context"
	"overseer/common/types"
	"overseer/ovsworker/fragments"
	"sync"
)

//TaskExecutor - executes a commissioned task
type TaskExecutor struct {
	store    map[string]fragments.FragmentStatus
	statChan chan fragments.FragmentStatus
	lock     sync.Mutex
}

//NewTaskExecutor - creates a new TaskExecutor
func NewTaskExecutor() *TaskExecutor {

	exec := &TaskExecutor{
		store:    map[string]fragments.FragmentStatus{},
		statChan: make(chan fragments.FragmentStatus),
		lock:     sync.Mutex{},
	}
	exec.updateTaskStatus()

	return exec
}

//ExecuteTask - starts work fragment
func (exec *TaskExecutor) ExecuteTask(fragment fragments.WorkFragment) fragments.FragmentStatus {

	status := fragments.FragmentStatus{
		TaskID:      fragment.TaskID(),
		ExecutionID: fragment.ExecutionID(),
		State:       types.WorkerTaskStatusExecuting,
		ReturnCode:  0,
		StatusCode:  0,
	}

	exec.lock.Lock()
	exec.store[fragment.ExecutionID()] = status
	exec.lock.Unlock()

	go func() {
		fragment.StartFragment(context.Background(), exec.statChan)
	}()

	return status

}

func (exec *TaskExecutor) updateTaskStatus() {

	go func() {
		for {
			exec.update(<-exec.statChan)
		}
	}()

}

func (exec *TaskExecutor) update(status fragments.FragmentStatus) {

	defer exec.lock.Unlock()
	exec.lock.Lock()
	exec.store[status.ExecutionID] = status
}

//GetTaskStatus - gets fragment status
func (exec *TaskExecutor) GetTaskStatus(executionID string) (fragments.FragmentStatus, bool) {
	defer exec.lock.Unlock()
	exec.lock.Lock()

	stat, exists := exec.store[executionID]
	return stat, exists
}

//TaskCount - returns the number of tasks currently processed
func (exec *TaskExecutor) TaskCount() int {
	defer exec.lock.Unlock()
	exec.lock.Lock()

	return len(exec.store)
}

//CleanupTask - removes a task
func (exec *TaskExecutor) CleanupTask(executionID string) {

	defer exec.lock.Unlock()
	exec.lock.Lock()

	delete(exec.store, executionID)
}

//TerminateTask - removes a task
func (exec *TaskExecutor) TerminateTask(executionID string) {

	defer exec.lock.Unlock()
	exec.lock.Lock()

	delete(exec.store, executionID)
}