package launcher

import "testing"

func TestStore(t *testing.T) {

	s := NewStore()
	if s == nil {
		t.Error("failed to initialize store")
	}

}
