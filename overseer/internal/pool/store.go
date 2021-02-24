package pool

import (
	"context"
	"fmt"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/internal/unique"
	"sync"
	"time"

	"github.com/przebro/databazaar/collection"
)

//Store - Holds active tasks
type Store struct {
	store      map[unique.TaskOrderID]*activeTask
	collection collection.DataCollection
	lock       sync.RWMutex
	syncTime   int
	log        logger.AppLogger
}

//NewStore - Creates a new store
func NewStore(collectionName string, log logger.AppLogger, provider *datastore.Provider, syncTime int) (*Store, error) {

	var err error
	var col collection.DataCollection

	if col, err = provider.GetCollection(collectionName); err != nil {
		log.Error("unable to load collection:", collectionName)
		return nil, err
	}

	store := &Store{
		store:      make(map[unique.TaskOrderID]*activeTask),
		lock:       sync.RWMutex{},
		collection: col,
		syncTime:   syncTime,
		log:        log,
	}

	store.restoreTasks()

	if store.syncTime > 0 {
		store.sync()
	}

	return store, nil
}

//get - gets an active task from a store
func (s *Store) get(taskID unique.TaskOrderID) (*activeTask, bool) {

	defer s.lock.RUnlock()
	s.lock.RLock()
	t, exists := s.store[taskID]
	return t, exists
}

//len - returns a number of tasks in store
func (s *Store) len() int {
	defer s.lock.RUnlock()
	s.lock.RLock()
	return len(s.store)
}

//add - Adds a new task
func (s *Store) add(taskID unique.TaskOrderID, task *activeTask) {

	defer s.lock.Unlock()
	s.lock.Lock()
	s.store[taskID] = task
}

//remove - Removes a task from a store
func (s *Store) remove(taskID unique.TaskOrderID) {

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

func (s *Store) sync() {
	go func() {
		for {
			select {
			case <-time.After(time.Duration(s.syncTime) * time.Second):
				{
					s.storeTasks()
				}
			}
		}
	}()
}

func (s *Store) storeTasks() {
	tsart := time.Now()
	ilist := []interface{}{}

	s.Over(func(id unique.TaskOrderID, at *activeTask) { ilist = append(ilist, at.getModel()) })
	err := s.collection.BulkUpdate(context.Background(), ilist)
	if err != nil {
		fmt.Println(err)
	}

	s.log.Info("store task complete:", time.Since(tsart))
}

func (s *Store) restoreTasks() {
	defer s.lock.Unlock()
	s.lock.Lock()

	var err error
	var cnt int64
	tsart := time.Now()

	model := taskPoolModel{}

	if cnt, _ = s.collection.Count(context.Background()); cnt == 0 {
		model.ID = "taskpool"
		s.log.Error("TaskPool model does not exist:", err)
		return
	}

	success := 0

	crsr, err := s.collection.All(context.Background())
	for crsr.Next(context.Background()) {
		model := activeTaskModel{}
		err := crsr.Decode(&model)
		if err != nil {
			s.log.Error("error loading task:", err)
			continue
		}

		task, err := fromModel(model)
		if err != nil {
			s.log.Error("error loading task:", err)
			continue
		}
		if _, exists := s.store[task.orderID]; exists {
			s.log.Error("error loading task, task with id", task.orderID, "already exists")
			continue
		}

		s.store[task.orderID] = task
		success++
	}

	s.log.Info("restore task complete, task restored:", success, ",time:", time.Since(tsart))
}
