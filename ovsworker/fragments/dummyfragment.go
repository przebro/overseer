package fragments

import (
	"errors"
	"fmt"
	"overseer/common/types"
	"overseer/ovsworker/msgheader"
	"overseer/proto/actions"

	"golang.org/x/net/context"
	"google.golang.org/protobuf/proto"
)

//DummyFragment - dummy work
type DummyFragment struct {
	workFragment
}

//FactoryDummy - Creates a new dummy factory
func FactoryDummy(header msgheader.TaskHeader, data []byte) (WorkFragment, error) {

	act := actions.DummyTaskAction{}
	if err := proto.Unmarshal(data, &act); err != nil {
		return nil, errors.New("")
	}
	w, err := newDummyFragment(header, &act)

	return w, err
}

//newDummyFragment - factory method
func newDummyFragment(header msgheader.TaskHeader, action *actions.DummyTaskAction) (WorkFragment, error) {

	frag := &DummyFragment{}
	frag.taskID = header.TaskID
	frag.start = make(chan FragmentStatus)
	frag.Variables = make([]string, len(header.Variables))
	for k, v := range header.Variables {
		frag.Variables = append(frag.Variables, fmt.Sprintf("%s=%s", k, v))
	}

	return frag, nil
}

//StartFragment - Start a new work
func (frag *DummyFragment) StartFragment(ctx context.Context, stat chan FragmentStatus) {

	go func() {
		status := FragmentStatus{TaskID: frag.taskID, ReturnCode: 0, PID: 0, State: types.WorkerTaskStatusEnded}
		stat <- status
	}()
}

//CancelFragment - cancels current fragment
func (frag *DummyFragment) CancelFragment() error {
	return nil
}

//Running - Informs caller that fragment is executed.
func (frag *DummyFragment) Running() FragmentStatus {
	return <-frag.start

}

//TaskID - Returns ID of a task associated with this fragment.
func (frag *DummyFragment) TaskID() string {
	return frag.taskID
}
