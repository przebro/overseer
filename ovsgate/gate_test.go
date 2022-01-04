package ovsgate

import (
	"testing"

	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/ovsgate/config"
)

func TestNewInstance(t *testing.T) {

	log := logger.NewTestLogger()
	_, err := NewInstance(config.OverseerGatewayConfig{}, log)
	if err != nil {
		t.Error("unexpected error")
	}
}
