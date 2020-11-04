package fragments

import (
	"sync"
)

//FragmentStatus - Contains inforamtion about a task status.
type FragmentStatus struct {
	TaskID        string
	Started       bool
	Ended         bool
	MarkForDelete bool
	ReturnCode    int
	StatusCode    int
	PID           int

	Output []string
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
