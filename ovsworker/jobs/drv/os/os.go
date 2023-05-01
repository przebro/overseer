package fragments

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/ovsworker/jobs"
	"github.com/przebro/overseer/ovsworker/msgheader"
	"github.com/przebro/overseer/ovsworker/status"
	"github.com/przebro/overseer/proto/actions"
	"github.com/rs/zerolog"

	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
)

func init() {
	jobs.RegisterFactory(types.TypeOs, OsJobFactory)
}

type OsStepDefinition struct {
	StepName string
	Command  string
	Args     []string
}

// osJob - Fragment that can execute OS command or script
type osJob struct {
	jobs.Job
	Command    string
	Arguments  []string
	RunAs      string
	Steps      []OsStepDefinition
	stdout     io.ReadCloser
	cancelFunc context.CancelFunc
}

// OsJobFactory - Creates new os factory
func OsJobFactory(ctx context.Context, header msgheader.TaskHeader, sysoutDir string, data []byte) (jobs.JobExecutor, error) {

	log := zerolog.Ctx(ctx)
	act := actions.OsTaskAction{}

	if err := proto.Unmarshal(data, &act); err != nil {
		log.Err(err).Msg("OsJobFactory")

		return nil, err
	}

	return newOsJob(ctx, header, sysoutDir, &act)
}

// newOsJob - factory method, creates a new os job
func newOsJob(ctx context.Context, header msgheader.TaskHeader, sysoutDir string, action *actions.OsTaskAction) (jobs.JobExecutor, error) {

	job := &osJob{}
	job.TaskID = header.TaskID
	job.ExecutionID = header.ExecutionID
	job.SysoutDir = sysoutDir
	job.Start = make(chan status.JobExecutionStatus)
	cmdarg := strings.Split(action.CommandLine, " ")
	job.Command = cmdarg[0]
	job.Arguments = cmdarg[1:]
	job.RunAs = action.Runas
	job.Steps = make([]OsStepDefinition, len(action.Steps))

	for i, step := range action.Steps {
		cmdarg := strings.Split(step.Command, " ")
		job.Steps[i] = OsStepDefinition{StepName: step.StepName, Command: cmdarg[0], Args: cmdarg[1:]}
	}

	job.Variables = make(map[string]string)
	for k, v := range header.Variables {
		job.Variables[k] = v
	}

	return job, nil
}

// StartFragment - starts a new work
func (j *osJob) StartJob(ctx context.Context, stat chan status.JobExecutionStatus) status.JobExecutionStatus {

	log := zerolog.Ctx(ctx).With().Str("task_id", j.TaskID).Str("execution_id", j.ExecutionID).Logger()
	tstat := make(chan struct{})

	go func() {

		exitCode := 0

		if j.Command != "" {
			cmd, cncl := j.prepareStep(ctx, OsStepDefinition{Command: j.Command, Args: j.Arguments, StepName: "command"})
			cmd.Env = j.prepareEnv()
			j.cancelFunc = cncl
			go j.steprun(ctx, "command", &cmd, tstat)
			_, isopen := <-tstat

			if !isopen {
				log.Err(errors.New("command failed")).Msg("command execution failed")
				stat <- status.StatusFailed(j.TaskID, j.ExecutionID, "command failed")
				return
			}

			exitCode = cmd.ProcessState.ExitCode()

		}

		for _, step := range j.Steps {

			cmd, cncl := j.prepareStep(ctx, step)
			cmd.Env = j.prepareEnv()
			j.cancelFunc = cncl
			go j.steprun(ctx, step.StepName, &cmd, tstat)
			_, isopen := <-tstat

			if !isopen {
				stat <- status.StatusFailed(j.TaskID, j.ExecutionID, "step failed")
				break
			}

			exitCode = cmd.ProcessState.ExitCode()
		}

		close(tstat)
		stat <- status.StatusEnded(j.TaskID, j.ExecutionID, exitCode, 0, 0)

	}()

	return status.StatusExecuting(j.TaskID, j.ExecutionID)
}

// CancelJob - cancels current work
func (j *osJob) CancelJob() error {

	if j.cancelFunc == nil {
		return errors.New("failed to cancel job")
	}

	j.cancelFunc()

	return nil
}

// JobTaskID - Returns ID of a task associated with this fragment.
func (j *osJob) JobTaskID() string {
	return j.TaskID
}

// JobExecutionID - Returns ID of a current run of a task associated with this fragment.
func (j *osJob) JobExecutionID() string {
	return j.ExecutionID
}

func (j *osJob) prepareEnv() []string {
	env := []string{}
	env = append(env, os.Environ()...)

	for k, v := range j.Variables {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}

	return env
}

func (j *osJob) prepareStep(ctx context.Context, step OsStepDefinition) (exec.Cmd, context.CancelFunc) {

	log := zerolog.Ctx(ctx)
	log.Info().Str("step", step.StepName).Msg("prepareStep")

	nctx, cncl := context.WithCancel(ctx)
	cmd := exec.CommandContext(nctx, step.Command, step.Args...)

	return *cmd, cncl
}

func (j *osJob) steprun(ctx context.Context, step string, cmd *exec.Cmd, stat chan<- struct{}) {

	log := zerolog.Ctx(ctx)

	sysoutPath := filepath.Join(j.SysoutDir, fmt.Sprintf("%s_%s", j.ExecutionID, step))
	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Error().Err(err).Str("step", step).Msg("failed to get stdout pipe")
	}

	log.Info().Str("step", step).Msg("step started")
	err = cmd.Start()

	if err != nil {
		log.Error().Err(err).Str("step", step).Msg("step run failed")
		close(stat)
		return
	}

	stdout(ctx, out, sysoutPath)
	cmd.Wait()
	log.Info().Str("step", step).Msg("step finished")
	stat <- struct{}{}
}

func stdout(ctx context.Context, out io.ReadCloser, path string) struct{} {

	log := zerolog.Ctx(ctx)
	wait := sync.WaitGroup{}
	file, err := os.Create(path)
	if err != nil {
		log.Error().Err(err).Msg("failed to create file")
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
