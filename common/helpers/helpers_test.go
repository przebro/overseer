package helpers

import (
	"path/filepath"
	"testing"
)

func TestGetDirectories(t *testing.T) {

	r, d, e := GetDirectories("../../../")
	if e != nil {
		t.Error(e)
	}
	if r != filepath.Dir(d) {
		t.Error("base dir not in root")
	}
}
