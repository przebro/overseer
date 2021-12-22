package fragments

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"overseer/common/logger"
	"overseer/common/types"
	"overseer/ovsworker/jobs"
	"overseer/ovsworker/msgheader"
	"overseer/ovsworker/status"
	"overseer/proto/actions"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"
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
	log        logger.AppLogger
}

//OsJobFactory - Creates new os factory
func OsJobFactory(header msgheader.TaskHeader, sysoutDir string, data []byte, log logger.AppLogger) (jobs.JobExecutor, error) {

	act := actions.OsTaskAction{}

	if err := proto.Unmarshal(data, &act); err != nil {
		log.Desugar().Error("OsJobFactory", zap.String("error", err.Error()))
		return nil, err
	}

	return newOsJob(header, sysoutDir, &act, log)
}

//newOsJob - factory method, creates a new os job
func newOsJob(header msgheader.TaskHeader, sysoutDir string, action *actions.OsTaskAction, log logger.AppLogger) (jobs.JobExecutor, error) {

	job := &osJob{}
	job.TaskID = header.TaskID
	job.ExecutionID = header.ExecutionID
	job.SysoutDir = sysoutDir
	job.Start = make(chan status.JobExecutionStatus)
	cmdarg := strings.Split(action.CommandLine, " ")
	job.Command = cmdarg[0]
	job.Arguments = cmdarg[1:]
	job.RunAs = action.Runas
	job.log = log

	job.Variables = make(map[string]string)
	for k, v := range header.Variables {
		job.Variables[k] = v
	}

	return job, nil
}

//StartFragment - starts a new work
func (j *osJob) StartJob(ctx context.Context, stat chan status.JobExecutionStatus) status.JobExecutionStatus {

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
		j.log.Desugar().Error("StartJob", zap.String("error", err.Error()))
		return status.StatusFailed(j.TaskID, j.ExecutionID, err.Error())

	}
	go j.run(cmd, stat)

	return status.StatusExecuting(j.TaskID, j.ExecutionID)

}

//CancelJob - cancels current work
func (j *osJob) CancelJob() error {

	if j.cancelFunc == nil {
		return errors.New("failed to cancel job")
	}

	j.cancelFunc()

	return nil
}

func (j *osJob) run(cmd *exec.Cmd, stat chan status.JobExecutionStatus) {

	err := cmd.Start()

	if err != nil {
		j.log.Desugar().Error("run", zap.String("error", err.Error()))

		stat <- status.StatusFailed(j.TaskID, j.ExecutionID, err.Error())
		return
	}

	fpath := filepath.Join(j.SysoutDir, j.ExecutionID)

	stdout(j.stdout, fpath, j.log)
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

func stdout(out io.ReadCloser, path string, log logger.AppLogger) struct{} {

	wait := sync.WaitGroup{}
	file, err := os.Create(path)
	if err != nil {
		log.Desugar().Error("stdout", zap.String("error", err.Error()))
		return struct{}{}
	}

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
