package logger

import (
	"testing"
)

func TestLogger(t *testing.T) {
	if logger != nil {
		t.Error("logger != nil")
	}
	log := NewLogger("./logs", 2)
	if logger == nil {
		t.Error("logger == nil")
	}
	err := log.SetLevel(12)
	if err == nil {
		t.Error("Invalid level value")
	}

	err = log.SetLevel(1)
	if err != nil {
		t.Error("Invalid level value")
	}

	if log.Level() != 1 {
		t.Error("Invalid level value")
	}

}
