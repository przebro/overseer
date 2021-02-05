package config

import (
	"path/filepath"
	"testing"
)

func TestConfig(t *testing.T) {

	c, err := Load("not_valid_path")
	if err == nil {
		t.Error("invalid path, expected error")
	}

	dir, _ := filepath.Abs("../../config/worker.json")
	c, err = Load(dir)
	if err != nil {
		t.Error(err)
	}

	c.GetLogConfiguration()

}
