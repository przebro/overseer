package config

import "testing"

func TestLoad_Errors(t *testing.T) {

	path := ""
	_, err := Load(path)
	if err == nil {
		t.Error("unexpected result")
	}

	path = "gateconfig.go"

	_, err = Load(path)
	if err == nil {
		t.Error("unexpected result")
	}
}

func TestLoad(t *testing.T) {

	path := "../../config/gateway.json"
	_, err := Load(path)
	if err != nil {
		t.Error("unexpected result:", err)
	}

}
