package ovsgate

import (
	"testing"

	"github.com/przebro/overseer/ovsgate/config"
)

func TestNewInstance(t *testing.T) {

	_, err := NewInstance(config.OverseerGatewayConfig{})
	if err != nil {
		t.Error("unexpected error")
	}
}
