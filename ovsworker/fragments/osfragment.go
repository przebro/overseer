package fragments

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"overseer/common/types"
	"overseer/ovsworker/msgheader"
	"overseer/proto/actions"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
)

//OsWorkFragment - Fragment that can execute OS command or script
type OsWorkFragment struct {
	workFragment
	Command   string
	Arguments []string
	RunAs     string
	stdout    io.ReadCloser
	cancFnc   context.CancelFunc
}

//FactoryOS - Creates new os factory
func FactoryOS(header msgheader.TaskHeader, sysoutDir string, data []byte) (WorkFragment, error) {

	act := actions.OsTaskAction{}
	if err := proto.Unmarshal(data, &act); err != nil {
		return nil, errors.New("")
	}
	w, err := newOsFragment(header, sysoutDir, &act)

	return w, err
}

//newOsFragment - factory method, creates a new os fragment
func newOsFragment(header msgheader.TaskHeader, sysoutDir string, action *actions.OsTaskAction) (WorkFragment, error) {

	frag := &OsWorkFragment{}
	frag.taskID = header.TaskID
	frag.executionID = header.ExecutionID
	frag.sysoutDir = sysoutDir
	frag.start = make(chan FragmentStatus)
	cmdarg := strings.Split(action.CommandLine, " ")
	frag.Command = cmdarg[0]
	frag.Arguments = cmdarg[1:]
	frag.RunAs = action.Runas

	frag.Variables = make([]string, len(header.Variables))
	for k, v := range header.Variables {
		frag.Variables = append(frag.Variables, fmt.Sprintf("%s=%s", k, v))
	}

	return frag, nil
}

//StartFragment - starts a new work
func (frag *OsWorkFragment) StartFragment(ctx context.Context, stat chan FragmentStatus) {

	var err error
	var cmdCtx context.Context
	cmdCtx, frag.cancFnc = context.WithCancel(ctx)

	cmd := exec.CommandContext(cmdCtx, frag.Command, frag.Arguments...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, frag.Variables...)

	frag.stdout, err = cmd.StdoutPipe()
	if err != nil {
		stat <- FragmentStatus{
			TaskID:      frag.taskID,
			ExecutionID: frag.executionID,
			State:       types.WorkerTaskStatusFailed,
			ReturnCode:  999,
			PID:         0,
		}

	} else {
		go frag.run(cmd, stat)
	}

}

//CancelFragment - cancels current work
func (frag *OsWorkFragment) CancelFragment() error {

	frag.cancFnc()

	return nil
}

func (frag *OsWorkFragment) run(cmd *exec.Cmd, stat chan FragmentStatus) {

	var err error
	err = cmd.Start()

	if err != nil {
		stat <- FragmentStatus{
			TaskID:      frag.taskID,
			ExecutionID: frag.executionID,
			State:       types.WorkerTaskStatusFailed,
			ReturnCode:  999,
			PID:         0,
		}
		return

	}

	stat <- FragmentStatus{
		TaskID:      frag.taskID,
		ExecutionID: frag.executionID,
		State:       types.WorkerTaskStatusExecuting,
		ReturnCode:  0,
		PID:         0,
	}

	fpath := filepath.Join(frag.sysoutDir, frag.executionID)

	stdout(frag.stdout, fpath)
	err = cmd.Wait()

	stat <- FragmentStatus{
		TaskID:      frag.taskID,
		ExecutionID: frag.executionID,
		State:       types.WorkerTaskStatusEnded,
		ReturnCode:  cmd.ProcessState.ExitCode(),
		PID:         cmd.ProcessState.Pid(),
	}
}

//TaskID - Returns ID of a task associated with this fragment.
func (frag *OsWorkFragment) TaskID() string {
	return frag.taskID
}

//ExecutionID - Returns ID of a current run of a task associated with this fragment.
func (frag *OsWorkFragment) ExecutionID() string {
	return frag.executionID
}

func stdout(out io.ReadCloser, path string) struct{} {

	wait := sync.WaitGroup{}
	fmt.Println(path)
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
