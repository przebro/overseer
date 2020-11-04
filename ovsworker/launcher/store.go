package launcher

import (
	"fmt"
	"goscheduler/ovsworker/fragments"
	"sync"
)

//Store - Holds fragments
type Store struct {
	store map[string]fragments.WorkFragment
	lock  sync.RWMutex
}

//NewStore - Creates a new store
func NewStore() *Store {

	return &Store{store: make(map[string]fragments.WorkFragment), lock: sync.RWMutex{}}
}

//Add - Adds a fragment to a store
func (s *Store) Add(key string, fragment fragments.WorkFragment) error {

	defer s.lock.Unlock()
	s.lock.Lock()
	if _, exists := s.store[key]; exists == true {
		return fmt.Errorf("fragment with id:%s already exists", key)
	}

	s.store[key] = fragment

	return nil
}

//Remove - Removes a fragment from a store
func (s *Store) Remove(key string) {
	defer s.lock.Unlock()
	s.lock.Lock()
	delete(s.store, key)
}

//Get - Gets a fragment from a store
func (s *Store) Get(key string) (fragments.WorkFragment, bool) {
	defer s.lock.RUnlock()
	s.lock.RLock()

	v, exists := s.store[key]

	return v, exists

}
