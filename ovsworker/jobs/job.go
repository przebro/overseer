package jobs

import (
	"context"
	"errors"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/ovsworker/msgheader"
	"github.com/przebro/overseer/ovsworker/status"
	"github.com/rs/zerolog/log"
)

var factories map[types.TaskType]JobFactorytMethod = make(map[types.TaskType]JobFactorytMethod)

// RegisterFactory - registers a new task factory
func RegisterFactory(t types.TaskType, fm JobFactorytMethod) {
	factories[t] = fm
}

// NewJobExecutor - Creates a new job that will be executed
func NewJobExecutor(ctx context.Context, header msgheader.TaskHeader, sysoutDir string, data []byte) (JobExecutor, error) {

	var job JobExecutor
	var err error

	lg := log.Ctx(ctx).With().Str("service", "exec").Logger()

	method, exists := factories[header.Type]
	if !exists {
		lg.Error().Err(err).Msg("method does not exist")
		return nil, errors.New("unable to create job executor")
	}

	if job, err = method(ctx, header, sysoutDir, data); err != nil {
		lg.Error().Err(err).Msg("unable to create Job")
		return nil, err
	}

	return job, nil

}

// JobFactorytMethod - Creates a new executable Job
type JobFactorytMethod func(ctx context.Context, header msgheader.TaskHeader, sysoutDir string, data []byte) (JobExecutor, error)

// Job - represents executed job
type Job struct {
	TaskID      string
	ExecutionID string
	SysoutDir   string
	Start       chan status.JobExecutionStatus
	Variables   map[string]string
}

// JobExecutor - Represents a piece of work that will be executed.
type JobExecutor interface {
	StartJob(ctx context.Context, stat chan status.JobExecutionStatus) status.JobExecutionStatus
	CancelJob() error
	//JobTaskID - Returns ID of a task associated with this job.
	JobTaskID() string
	JobExecutionID() string
}
