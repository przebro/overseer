package resources

import (
	"errors"
	"sync"
	"time"

	"github.com/przebro/databazaar/collection"
)

var (
	errKeyNotFound = errors.New("key not found")
	errKeyExists   = errors.New("key aleready exists")
)

//resourceStore - inmemmory structure with sync and backup
type resourceStore struct {
	rw    readWriter
	items map[string]interface{}
	lock  sync.RWMutex
	stime int
	col   collection.DataCollection
}

func newStore(rw readWriter, stime int) (*resourceStore, error) {

	items, err := rw.Load()
	if err != nil {
		return nil, err
	}

	store := &resourceStore{items: items, lock: sync.RWMutex{}, rw: rw}

	if stime > 0 {
		go watch(store, stime)
	}

	return store, nil
}

func watch(s *resourceStore, tm int) {
	t := time.NewTicker(time.Duration(tm) * time.Second)
	for {
		select {
		case <-t.C:
			{
				s.Sync()
			}
		}
	}
}

func (s *resourceStore) Insert(key string, item interface{}) error {

	defer s.lock.Unlock()
	s.lock.Lock()
	if _, ok := s.items[key]; !ok {
		s.items[key] = item
		return nil
	}

	return errKeyExists
}
func (s *resourceStore) Get(key string) (interface{}, bool) {

	defer s.lock.RUnlock()
	s.lock.RLock()

	item, ok := s.items[key]
	return item, ok
}
func (s *resourceStore) Update(key string, item interface{}) error {

	defer s.lock.Unlock()
	s.lock.Lock()
	if _, ok := s.items[key]; ok {
		s.items[key] = item
		return nil
	}

	return errKeyNotFound
}
func (s *resourceStore) Delete(key string) {
	defer s.lock.Unlock()
	s.lock.Lock()

	delete(s.items, key)
}

func (s *resourceStore) All() []interface{} {

	defer s.lock.RUnlock()
	s.lock.RLock()
	col := make([]interface{}, len(s.items), len(s.items))
	i := 0
	for _, v := range s.items {
		col[i] = v
		i++
	}

	return col

}

func (s *resourceStore) Sync() {

	tmp := map[string]interface{}{}
	s.lock.Lock()
	for k, v := range s.items {
		tmp[k] = v
	}
	s.lock.Unlock()

	s.rw.Write(tmp)

}
