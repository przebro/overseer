package pool

import (
	"overseer/common/types/date"
	"overseer/overseer/internal/taskdef"
	"overseer/overseer/internal/unique"
	"sync"
	"testing"
	"time"
)

func TestStore(t *testing.T) {

	store, _ := NewStore(testCollectionName, log, 3600, provider)
	if store == nil {
		t.Error("Store not created")
	}

	orderID := unique.TaskOrderID("12345")
	task := &activeTask{orderID: orderID}

	store.add(orderID, task)
	store.add(unique.TaskOrderID("33333"), &activeTask{orderID: unique.TaskOrderID("33333")})
	store.add(unique.TaskOrderID("12346"), &activeTask{orderID: unique.TaskOrderID("12346")})

	if store.len() != 3 {
		t.Error("Invalid store size expected 1 actual :", store.len())
	}

	_, exists := store.get(unique.TaskOrderID("54321"))
	if exists == true {
		t.Error("Unexpected result expected", false, "actual", exists)
	}

	ts, exists := store.get(unique.TaskOrderID("12345"))
	if exists == false {
		t.Error("Unexpected result expected", false, "actual", exists)
	}

	if ts.orderID != orderID {
		t.Error("Unexpected result expected", orderID, "actual", ts.orderID)
	}

	store.remove(unique.TaskOrderID("12346"))
	store.remove(unique.TaskOrderID("12346"))

	if store.len() != 2 {
		t.Error("Invalid store size expected 1 actual :", store.len())
	}

}

func TestStoreOver(t *testing.T) {

	store, _ := NewStore(testCollectionName, log, 3600, provider)
	if store == nil {
		t.Error("Store not created")
	}

	//sience over is non blocking for read operation, this is expected result
	expected := []string{"33333", "44444", "55555", "ABCDEF", "ABCDEF", "ABCDEF"}
	actual := []string{}
	vchan := make(chan string, 6)

	wg := sync.WaitGroup{}
	wg.Add(2)

	store.add(unique.TaskOrderID("33333"), &activeTask{orderID: unique.TaskOrderID("33333")})
	store.add(unique.TaskOrderID("44444"), &activeTask{orderID: unique.TaskOrderID("44444")})
	store.add(unique.TaskOrderID("55555"), &activeTask{orderID: unique.TaskOrderID("55555")})

	go func(s *Store) {

		s.Over(func(k unique.TaskOrderID, v *activeTask) {
			//take some time before first operation, make that second goroutine will write first
			time.Sleep(100 * time.Millisecond)

			vchan <- "ABCDEF"
		})

		wg.Done()

	}(store)

	go func(s *Store) {
		ids := []string{"33333", "44444", "55555"}
		time.Sleep(10 * time.Millisecond)
		for _, n := range ids {
			t, _ := s.get(unique.TaskOrderID(n))
			vchan <- string(t.orderID)
		}

		wg.Done()

	}(store)

	wg.Wait()
	close(vchan)

	for x := range vchan {
		actual = append(actual, x)
	}

	for x, a := range actual {
		if expected[x] != a {
			t.Error("Unexpected order, epxected", expected, "actual", actual)
		}
	}

}
func TestStoreForEach(t *testing.T) {

	store, _ := NewStore(testCollectionName, log, 3600, provider)
	if store == nil {
		t.Error("Store not created")
	}

	//sience foreach is  blocking  operation, this is expected result
	expected := []string{"ABCDEF", "ABCDEF", "ABCDEF", "33333", "44444", "55555"}
	actual := []string{}
	vchan := make(chan string, 6)

	wg := sync.WaitGroup{}
	wg.Add(2)

	store.add(unique.TaskOrderID("33333"), &activeTask{orderID: unique.TaskOrderID("33333")})
	store.add(unique.TaskOrderID("44444"), &activeTask{orderID: unique.TaskOrderID("44444")})
	store.add(unique.TaskOrderID("55555"), &activeTask{orderID: unique.TaskOrderID("55555")})

	go func(s *Store) {

		s.ForEach(func(k unique.TaskOrderID, v *activeTask) {

			time.Sleep(100 * time.Millisecond)

			vchan <- "ABCDEF"
		})

		wg.Done()

	}(store)

	go func(s *Store) {
		ids := []string{"33333", "44444", "55555"}
		time.Sleep(20 * time.Millisecond)
		for _, n := range ids {
			v, _ := s.get(unique.TaskOrderID(n))
			vchan <- string(v.orderID)
		}

		wg.Done()

	}(store)

	wg.Wait()
	close(vchan)

	for x := range vchan {
		actual = append(actual, x)
	}

	for x, a := range actual {
		if expected[x] != a {
			t.Error("Unexpected order, epxected", expected, "actual", actual)
		}
	}

}

func TestStore_StoreTask(t *testing.T) {

	store, _ := NewStore(testStoreTaskName, log, 0, provider)
	if store == nil {
		t.Error("Store not created")
	}

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{FromTime: "", ToTime: "", OrderType: taskdef.OrderingManual}).
		Build()
	if err != nil {
		t.Fatal("Unable to construct task")
	}

	active := newActiveTask(seq.Next(), date.CurrentOdate(), definition)

	store.add(active.orderID, active)

	if store.len() != 1 {
		t.Error("unexpected result:", store.len(), "expected : 1")
	}
	//:TODO rewrite test
	store.storeTasks()
	store.restoreTasks()
}
