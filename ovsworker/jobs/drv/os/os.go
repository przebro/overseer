package fragments

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"overseer/common/types"
	"overseer/ovsworker/jobs"
	"overseer/ovsworker/msgheader"
	"overseer/ovsworker/status"
	"overseer/proto/actions"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
)

func init() {
	jobs.RegisterFactory(types.TypeOs, OsJobFactory)
}

//osJob - Fragment that can execute OS command or script
type osJob struct {
	jobs.Job
	Command    string
	Arguments  []string
	RunAs      string
	stdout     io.ReadCloser
	cancelFunc context.CancelFunc
}

//OsJobFactory - Creates new os factory
func OsJobFactory(header msgheader.TaskHeader, sysoutDir string, data []byte) (jobs.JobExecutor, error) {

	act := actions.OsTaskAction{}
	if err := proto.Unmarshal(data, &act); err != nil {
		return nil, errors.New("")
	}

	return newOsJob(header, sysoutDir, &act)
}

//newOsJob - factory method, creates a new os job
func newOsJob(header msgheader.TaskHeader, sysoutDir string, action *actions.OsTaskAction) (jobs.JobExecutor, error) {

	job := &osJob{}
	job.TaskID = header.TaskID
	job.ExecutionID = header.ExecutionID
	job.SysoutDir = sysoutDir
	job.Start = make(chan status.JobExecutionStatus)
	cmdarg := strings.Split(action.CommandLine, " ")
	job.Command = cmdarg[0]
	job.Arguments = cmdarg[1:]
	job.RunAs = action.Runas

	job.Variables = make(map[string]string)
	for k, v := range header.Variables {
		job.Variables[k] = v
	}

	return job, nil
}

//StartFragment - starts a new work
func (j *osJob) StartJob(ctx context.Context, stat chan status.JobExecutionStatus) {

	var err error
	var cmdCtx context.Context
	cmdCtx, j.cancelFunc = context.WithCancel(ctx)

	cmd := exec.CommandContext(cmdCtx, j.Command, j.Arguments...)
	cmd.Env = append(cmd.Env, os.Environ()...)

	for k, v := range j.Variables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	j.stdout, err = cmd.StdoutPipe()
	if err != nil {
		stat <- status.StatusFailed(j.TaskID, j.ExecutionID, err.Error())

	} else {
		go j.run(cmd, stat)
	}
}

//CancelJob - cancels current work
func (j *osJob) CancelJob() error {

	if j.cancelFunc == nil {
		return errors.New("cancel function is nil")
	}

	return nil
}

func (j *osJob) run(cmd *exec.Cmd, stat chan status.JobExecutionStatus) {

	err := cmd.Start()

	if err != nil {
		stat <- status.StatusFailed(j.TaskID, j.ExecutionID, err.Error())
		return
	}

	stat <- status.StatusExecuting(j.TaskID, j.ExecutionID)

	fpath := filepath.Join(j.SysoutDir, j.ExecutionID)

	stdout(j.stdout, fpath)
	cmd.Wait()

	stat <- status.StatusEnded(j.TaskID, j.ExecutionID, cmd.ProcessState.ExitCode(), cmd.ProcessState.Pid(), 0)
}

//JobTaskID - Returns ID of a task associated with this fragment.
func (j *osJob) JobTaskID() string {
	return j.TaskID
}

//JobExecutionID - Returns ID of a current run of a task associated with this fragment.
func (j *osJob) JobExecutionID() string {
	return j.ExecutionID
}

func stdout(out io.ReadCloser, path string) struct{} {

	wait := sync.WaitGroup{}
	file, _ := os.Create(path)

	wait.Add(1)
	go func(f *os.File) {
		for {
			buff := make([]byte, 1024)
			bytes, err := out.Read(buff)
			if err == io.EOF {
				break
			}
			if err != nil {
				f.WriteString(err.Error())
				break
			}

			f.Write(buff[:bytes])
		}

		wait.Done()
	}(file)

	wait.Wait()
	return struct{}{}
}
