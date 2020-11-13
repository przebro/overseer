package launcher

import (
	"context"
	"errors"
	"goscheduler/ovsworker/fragments"
)

//FragmentLauncher - Holds fragments that are currently processed
type FragmentLauncher struct {
	store *Store
}

//NewFragmentLauncher - Creates a new fragment launcher
func NewFragmentLauncher() *FragmentLauncher {

	w := &FragmentLauncher{
		store: NewStore(),
	}
	return w
}

func (exec *FragmentLauncher) addFragment(fragment fragments.WorkFragment) error {

	f, exists := exec.store.Get(fragment.TaskID())
	if exists && f.StatusFragment().MarkForDelete {
		exec.store.Remove(fragment.TaskID())
	}

	return exec.store.Add(fragment.TaskID(), fragment)
}

//Execute - Executes a fragment
func (exec *FragmentLauncher) Execute(ctx context.Context, taskID string) (<-chan fragments.FragmentStatus, error) {

	ch := make(chan fragments.FragmentStatus)

	fragment, exists := exec.store.Get(taskID)

	if !exists {
		return nil, errors.New("unable to find fragment with given ID")
	}
	go func() { ch <- fragment.StartFragment(ctx) }()

	return ch, nil

}

//Status - Gets the status of a fragment
func (exec *FragmentLauncher) Status(taskID string) (fragments.FragmentStatus, error) {

	frag, exists := exec.store.Get(taskID)
	if !exists {
		return fragments.FragmentStatus{}, errors.New("fragment with given id does not exists")
	}

	status := frag.StatusFragment()
	if status.MarkForDelete {
		exec.store.Remove(taskID)
	}

	return status, nil
}

//Tasks - returns numeber of active fragments
func (exec *FragmentLauncher) Tasks() int {

	return len(exec.store.store)

}
