package fragments

import (
	common "overseer/common/types"
	"sync"
)

//FragmentStatus - Contains inforamtion about a task status.
type FragmentStatus struct {
	TaskID      string
	ExecutionID string
	State       common.WorkerTaskStatus
	ReturnCode  int
	StatusCode  int
	PID         int
}

//StatusStore - Holds status of a fragment
type StatusStore struct {
	status FragmentStatus
	lock   sync.RWMutex
}

//NewStore - Creates a new store
func NewStore() *StatusStore {

	return &StatusStore{}
}

//Set - Sets values in a store
func (store *StatusStore) Set(status FragmentStatus) {
	defer store.lock.Unlock()
	store.lock.Lock()
	store.status = status
}

//Get - Gets value from a store
func (store *StatusStore) Get() FragmentStatus {
	defer store.lock.RUnlock()
	store.lock.RLock()
	return store.status
}
