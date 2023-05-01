package pool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/datastore"
	"github.com/przebro/overseer/overseer/internal/pool/activetask"
	"github.com/przebro/overseer/overseer/internal/pool/readers"
	"github.com/rs/zerolog"

	"github.com/przebro/databazaar/collection"
)

const collectionName = "tasks"

// Store - Holds active tasks
type Store struct {
	store      map[unique.TaskOrderID]*activetask.TaskInstance
	collection collection.DataCollection
	lock       sync.RWMutex
	log        zerolog.Logger
	synctime   int
}

// NewStore - Creates a new store
func NewStore(log zerolog.Logger, synctime int, provider datastore.CollectionProvider, rdr readers.ActiveDefinitionReader) (*Store, error) {

	var err error
	var col collection.DataCollection

	if col, err = provider.GetCollection(context.Background(), collectionName); err != nil {
		log.Error().Err(err).Str("collection", collectionName).Msg("unable to load collection")
		return nil, err
	}

	store := &Store{
		store:      make(map[unique.TaskOrderID]*activetask.TaskInstance),
		lock:       sync.RWMutex{},
		collection: col,
		log:        log,
		synctime:   synctime,
	}

	store.restorePool(rdr)

	return store, nil
}

// get - gets an active task from a store
func (s *Store) get(taskID unique.TaskOrderID) (*activetask.TaskInstance, bool) {

	defer s.lock.RUnlock()
	s.lock.RLock()
	t, exists := s.store[taskID]
	return t, exists
}

// len - returns a number of tasks in store
func (s *Store) len() int {
	defer s.lock.RUnlock()
	s.lock.RLock()
	return len(s.store)
}

// add - Adds a new task
func (s *Store) add(taskID unique.TaskOrderID, task *activetask.TaskInstance) {

	defer s.lock.Unlock()
	s.lock.Lock()
	s.store[taskID] = task
}

// remove - Removes a task from a store
func (s *Store) remove(taskID unique.TaskOrderID) {

	defer s.lock.Unlock()
	s.lock.Lock()
	delete(s.store, taskID)
}

// ForEach - Performs an action for each task in store, this method should be used if task has to be modified
func (s *Store) ForEach(f func(unique.TaskOrderID, *activetask.TaskInstance)) {

	defer s.lock.Unlock()
	s.lock.Lock()
	for k, v := range s.store {
		f(k, v)
	}
}

// Over - Performs an action over tasks in store, this method should be used to read values from tasks
func (s *Store) Over(f func(unique.TaskOrderID, *activetask.TaskInstance)) {
	defer s.lock.RUnlock()
	s.lock.RLock()
	for k, v := range s.store {
		f(k, v)
	}
}

func (s *Store) watch(quiesce <-chan bool, shutdown <-chan struct{}) <-chan struct{} {

	inform := make(chan struct{})
	tm := s.synctime

	go func(qch <-chan bool, sch <-chan struct{}, inf chan<- struct{}) {

		var isActive bool

		for {
			select {
			case <-time.After(time.Duration(tm) * time.Second):
				{
					if !isActive {
						continue
					}
					s.storeTasks()
				}
			case v := <-qch:
				{
					isActive = v
					s.storeTasks()
					inf <- struct{}{}
				}
			case <-sch:
				{
					s.storeTasks()
					close(inform)
					return
				}
			}
		}
	}(quiesce, shutdown, inform)
	return inform
}

func (s *Store) storeTasks() {
	tsart := time.Now()
	ilist := []interface{}{}

	s.Over(func(id unique.TaskOrderID, at *activetask.TaskInstance) { ilist = append(ilist, at.GetModel()) })
	err := s.collection.BulkUpdate(context.Background(), ilist)
	if err != nil {
		s.log.Err(err).Msg("store tasks failed")
	}

	s.log.Info().Dur("time", time.Since(tsart)).Msg("store task completed")
}

func (s *Store) restorePool(rdr readers.ActiveDefinitionReader) {
	defer s.lock.Unlock()
	s.lock.Lock()

	var err error
	var cnt int64
	tsart := time.Now()

	model := activetask.TaskPoolModel{}

	if cnt, _ = s.collection.Count(context.Background()); cnt == 0 {
		model.ID = "taskpool"
		fmt.Println("TASK POOL DOES NOT EXIST")
		s.log.Error().Err(err).Msg("TaskPool model does not exist")
		return
	}

	success := 0

	crsr, err := s.collection.All(context.Background())
	for crsr.Next(context.Background()) {
		model := activetask.ActiveTaskModel{}
		err := crsr.Decode(&model)
		if err != nil {
			s.log.Error().Err(err).Msg("error loading task")
			continue
		}

		task, err := activetask.FromModel(model, rdr)
		if err != nil {
			s.log.Error().Err(err).Msg("error loading task")
			continue
		}
		if _, exists := s.store[task.OrderID()]; exists {
			s.log.Error().Str("order_id", string(task.OrderID())).Err(err).Msg("error loading task,already exists")

			continue
		}

		s.store[task.OrderID()] = task
		success++
	}

	s.log.Info().Int("count", success).Dur("time", time.Since(tsart)).Msg("restore task completed")
}
