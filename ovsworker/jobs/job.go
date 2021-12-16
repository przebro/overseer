package jobs

import (
	"context"
	"errors"
	"overseer/common/types"
	"overseer/ovsworker/msgheader"
	"overseer/ovsworker/status"
)

var factories map[types.TaskType]JobFactorytMethod = make(map[types.TaskType]JobFactorytMethod)

//RegisterFactory - registers a new task factory
func RegisterFactory(t types.TaskType, fm JobFactorytMethod) {
	factories[t] = fm
}

//NewJobExecutor - Creates a new job that will be executed
func NewJobExecutor(header msgheader.TaskHeader, sysoutDir string, data []byte) (JobExecutor, error) {

	var job JobExecutor
	var err error

	method, exists := factories[header.Type]
	if !exists {
		return nil, errors.New("unable to create job executor")
	}

	if job, err = method(header, sysoutDir, data); err != nil {
		return nil, err
	}

	return job, nil

}

//JobFactorytMethod - Creates a new executable Job
type JobFactorytMethod func(header msgheader.TaskHeader, sysoutDir string, data []byte) (JobExecutor, error)

//Job - represents executed job
type Job struct {
	TaskID      string
	ExecutionID string
	SysoutDir   string
	Start       chan status.JobExecutionStatus
	Variables   map[string]string
}

//JobExecutor - Represents a piece of work that will be executed.
type JobExecutor interface {
	StartJob(ctx context.Context, stat chan status.JobExecutionStatus)
	CancelJob() error
	//JobTaskID - Returns ID of a task associated with this job.
	JobTaskID() string
	JobExecutionID() string
}
