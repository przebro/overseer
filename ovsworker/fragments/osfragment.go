package fragments

import (
	"errors"
	"fmt"
	"goscheduler/ovsworker/msgheader"
	"goscheduler/proto/actions"
	"io"
	"os"
	"os/exec"
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
func FactoryOS(header msgheader.TaskHeader, data []byte) (WorkFragment, error) {

	act := actions.OsTaskAction{}
	if err := proto.Unmarshal(data, &act); err != nil {
		return nil, errors.New("")
	}
	w, err := newOsFragment(header, &act)

	return w, err
}

//newOsFragment - factory method, creates a new os fragment
func newOsFragment(header msgheader.TaskHeader, action *actions.OsTaskAction) (WorkFragment, error) {

	frag := &OsWorkFragment{}
	frag.taskID = header.TaskID
	frag.start = make(chan FragmentStatus)
	frag.Status = NewStore()
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
func (frag *OsWorkFragment) StartFragment(ctx context.Context) FragmentStatus {

	var err error
	var cmdCtx context.Context
	cmdCtx, frag.cancFnc = context.WithCancel(ctx)

	cmd := exec.CommandContext(cmdCtx, frag.Command, frag.Arguments...)
	cmd.Env = append(cmd.Env, os.Environ()...)
	cmd.Env = append(cmd.Env, frag.Variables...)

	frag.stdout, err = cmd.StdoutPipe()
	if err != nil {
		frag.Status.Set(FragmentStatus{TaskID: frag.taskID,
			Started:       false,
			Ended:         true,
			ReturnCode:    999,
			PID:           0,
			Output:        []string{err.Error()},
			MarkForDelete: true,
		})
		return frag.Status.Get()
	}
	ch := make(chan struct{})
	go frag.run(cmd, ch)

	<-ch

	return frag.Status.Get()
}

//StatusFragment - gets current status of a fragment
func (frag *OsWorkFragment) StatusFragment() FragmentStatus {
	return frag.Status.Get()
}

//CancelFragment - cancels current work
func (frag *OsWorkFragment) CancelFragment() error {

	frag.cancFnc()

	return nil
}

func (frag *OsWorkFragment) run(cmd *exec.Cmd, s chan<- struct{}) {

	var err error
	err = cmd.Start()
	if err != nil {
		frag.Status.Set(FragmentStatus{TaskID: frag.taskID,
			Started:       false,
			Ended:         true,
			ReturnCode:    999,
			PID:           0,
			Output:        []string{err.Error()},
			MarkForDelete: true,
		})
		s <- struct{}{}
		return
	}

	frag.Status.Set(FragmentStatus{TaskID: frag.taskID,
		Started:    true,
		Ended:      false,
		ReturnCode: 0,
		PID:        0,
		Output:     []string{},
	})

	s <- struct{}{}
	outs := stdout(frag.stdout)
	err = cmd.Wait()

	if err != nil {
		outs = append(outs, err.Error())
	}

	frag.Status.Set(FragmentStatus{TaskID: frag.taskID,
		Started:       true,
		Ended:         true,
		ReturnCode:    cmd.ProcessState.ExitCode(),
		PID:           cmd.ProcessState.Pid(),
		Output:        outs,
		MarkForDelete: true,
	})
}

//TaskID - Returns ID of a task associated with this fragment.
func (frag *OsWorkFragment) TaskID() string {
	return frag.taskID
}

func stdout(out io.ReadCloser) []string {

	stdOutput := make([]string, 0)
	wait := sync.WaitGroup{}
	wait.Add(1)
	go func() {
		for {
			buff := make([]byte, 1024)
			bytes, err := out.Read(buff)
			if err == io.EOF {
				break
			}
			if err != nil {
				stdOutput = append(stdOutput, err.Error())
				break
			}
			stdOutput = append(stdOutput, strings.Split(string(buff[:bytes]), "\n")...)

		}

		wait.Done()
	}()

	wait.Wait()

	return stdOutput
}
