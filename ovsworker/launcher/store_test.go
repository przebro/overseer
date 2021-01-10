package launcher

import (
	"overseer/ovsworker/fragments"
	"testing"
)

func TestStore(t *testing.T) {

	s := NewStore()
	if s == nil {
		t.Error("failed to initialize store")
	}

	s.Add("12345", &fragments.DummyFragment{})
	if len(s.store) != 1 {
		t.Error("unexpected len")
	}

	err := s.Add("12345", &fragments.DummyFragment{})

	if err == nil {
		t.Error("Expected error")
	}

	if len(s.store) != 1 {
		t.Error("unexpected len")
	}

	_, ok := s.Get("12345")
	if ok != true {
		t.Error("unexpected result")
	}

	_, ok = s.Get("44444")
	if ok == true {
		t.Error("unexpected result")
	}

	s.Remove("12345")
	if len(s.store) != 0 {
		t.Error("unexpected len")
	}

}
