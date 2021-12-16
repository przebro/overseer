package dummy

import (
	"context"
	"errors"
	"overseer/common/types"
	"overseer/ovsworker/jobs"
	"overseer/ovsworker/msgheader"
	"overseer/ovsworker/status"
	"overseer/proto/actions"

	"google.golang.org/protobuf/proto"
)

func init() {
	jobs.RegisterFactory(types.TypeDummy, DummyJobFactory)
}

//DummyJobFactory - Creates a new dummy factory
func DummyJobFactory(header msgheader.TaskHeader, sysoutDir string, data []byte) (jobs.JobExecutor, error) {

	act := actions.DummyTaskAction{}
	if err := proto.Unmarshal(data, &act); err != nil {
		return nil, errors.New("")
	}
	w, err := newDummyJob(header, sysoutDir, &act)

	return w, err
}

//newDummyJob - factory method
func newDummyJob(header msgheader.TaskHeader, sysoutDir string, action *actions.DummyTaskAction) (jobs.JobExecutor, error) {

	j := &dummyJob{}
	j.TaskID = header.TaskID
	j.ExecutionID = header.ExecutionID
	j.Start = make(chan status.JobExecutionStatus)
	j.Variables = make(map[string]string)

	for k, v := range header.Variables {
		j.Variables[k] = v
	}

	return j, nil
}

//DummyJob - dummy work
type dummyJob struct {
	jobs.Job
}

//StartJob - Start a new work
func (j *dummyJob) StartJob(ctx context.Context, stat chan status.JobExecutionStatus) {

	go func() {
		status := status.JobExecutionStatus{TaskID: j.TaskID, ExecutionID: j.ExecutionID, ReturnCode: 0, PID: 0, State: types.WorkerTaskStatusEnded}
		stat <- status
	}()
}

//CancelFragment - cancels current job
func (j *dummyJob) CancelJob() error {
	return nil
}

//Running - Informs caller that job is executed.
func (j *dummyJob) Running() status.JobExecutionStatus {
	return <-j.Start

}

//TaskID - Returns ID of a task associated with this fragment.
func (j *dummyJob) JobTaskID() string {
	return j.TaskID
}

//ExecutionID - Returns ID of a current run of a task associated with this fragment.
func (j *dummyJob) JobExecutionID() string {
	return j.ExecutionID
}
