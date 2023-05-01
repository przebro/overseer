package logger

import (
	"testing"
)

func TestLogConfiguration(t *testing.T) {

	conf := LoggerConfiguration{
		LogLevel:     2,
		SizeLimit:    10,
		LogDirectory: "logs",
		FilePrefix:   "prefix",
	}
	if result := conf.Prefix(); result != "prefix" {
		t.Error("unexpected result:", result, "expected:prefix")
	}

	if result := conf.Directory(); result != "logs" {
		t.Error("unexpected result:", result, "expected:logs")
	}

	if result := conf.Level(); result != 2 {
		t.Error("unexpected result:", result, "expected:2")
	}

	if result := conf.Limit(); result != 10 {
		t.Error("unexpected result:", result, "expected:10")
	}

}
