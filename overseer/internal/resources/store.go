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
	rw     readWriter
	items  map[string]interface{}
	lock   sync.RWMutex
	stime  int
	col    collection.DataCollection
	done   <-chan struct{}
	shtdwn chan struct{}
}

func newStore(rw readWriter, stime int) (*resourceStore, error) {

	items, err := rw.Load()
	if err != nil {
		return nil, err
	}

	store := &resourceStore{items: items, lock: sync.RWMutex{}, rw: rw, shtdwn: make(chan struct{}), stime: stime}

	return store, nil
}

func (s *resourceStore) watch(tm int, shutdown <-chan struct{}) <-chan struct{} {

	inform := make(chan struct{})

	if tm == 0 {
		go func(sc <-chan struct{}, inf chan<- struct{}) {
			<-sc
			s.sync()
			close(inf)
			return

		}(shutdown, inform)

		return inform
	}

	go func(sc <-chan struct{}, inf chan<- struct{}) {

		for {
			select {
			case <-time.After(time.Duration(tm) * time.Second):
				{
					s.sync()
				}
			case <-sc:
				{
					s.sync()
					close(inf)
					return
				}
			}
		}

	}(shutdown, inform)

	return inform
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
func (s *resourceStore) Delete(key string) error {
	defer s.lock.Unlock()
	s.lock.Lock()

	if _, ok := s.items[key]; !ok {
		return errKeyNotFound
	}

	delete(s.items, key)
	return nil
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

func (s *resourceStore) start() {
	s.done = s.watch(s.stime, s.shtdwn)
}

func (s *resourceStore) shutdown() {

	s.shtdwn <- struct{}{}
	close(s.shtdwn)
	<-s.done
}

func (s *resourceStore) sync() {

	tmp := map[string]interface{}{}
	s.lock.Lock()
	for k, v := range s.items {
		tmp[k] = v
	}
	s.lock.Unlock()
	s.rw.Write(tmp)
}
