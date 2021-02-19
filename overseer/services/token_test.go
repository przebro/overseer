package services

import (
	"overseer/overseer/config"
	"testing"
)

func TestCreateTokenVerifier(t *testing.T) {
	cfg := config.SecurityConfiguration{Secret: "Invalid^%*()", Issuer: "issuer", Timeout: 0}
	if _, err := NewTokenCreatorVerifier(cfg); err == nil {
		t.Error("unexpected result")
	}

	cfg.Secret = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

	if _, err := NewTokenCreatorVerifier(cfg); err != nil {
		t.Error("unexpected result:", err)
	}

}
