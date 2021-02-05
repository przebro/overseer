package pool

import (
	"overseer/overseer/internal/unique"
	"sync"
)

//Store - Holds active tasks
type Store struct {
	store map[unique.TaskOrderID]*activeTask
	lock  sync.RWMutex
}

//NewStore - Creates a new store
func NewStore() *Store {
	return &Store{store: make(map[unique.TaskOrderID]*activeTask), lock: sync.RWMutex{}}
}

//Get - gets an active task from a store
func (s *Store) Get(taskID unique.TaskOrderID) (*activeTask, bool) {

	defer s.lock.RUnlock()
	s.lock.RLock()
	t, exists := s.store[taskID]
	return t, exists
}

//Len - returns a number of tasks in store
func (s *Store) Len() int {
	defer s.lock.RUnlock()
	s.lock.RLock()
	return len(s.store)
}

//Add - Adds a new task
func (s *Store) Add(taskID unique.TaskOrderID, task *activeTask) {

	defer s.lock.Unlock()
	s.lock.Lock()
	s.store[taskID] = task
}

//Remove - Removes a task from a store
func (s *Store) Remove(taskID unique.TaskOrderID) {

	defer s.lock.Unlock()
	s.lock.Lock()
	delete(s.store, taskID)
}

//ForEach - Performs an action for each task in store, this method should be used if task has to be modified
func (s *Store) ForEach(f func(unique.TaskOrderID, *activeTask)) {

	defer s.lock.Unlock()
	s.lock.Lock()
	for k, v := range s.store {
		f(k, v)
	}
}

//Over - Performs an action over tasks in store, this method should be used to read values from tasks
func (s *Store) Over(f func(unique.TaskOrderID, *activeTask)) {
	defer s.lock.RUnlock()
	s.lock.RLock()
	for k, v := range s.store {
		f(k, v)
	}
}
