package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"overseer/common/logger"
	"overseer/common/types"
)

//OverseerGatewayConfig - holds gateway configuration
type OverseerGatewayConfig struct {
	GatewayAddress  string `json:"gateAddress"`
	GatewayPort     int    `json:"gatePort"`
	OverseerAddress string `json:"overseerAddress"`
	OverseerPort    int    `json:"overseerPort"`
	//SecurityLevel -  this applies only to internal communication through grpc, not to external service
	SecurityLevel     types.ConnectionSecurityLevel `default:"none" json:"securityLevel" validate:"oneof=none server clientandserver"`
	GatewayCert       string                        `json:"cert"`
	GatewayKey        string                        `json:"key"`
	GatewayCertPolicy types.CertPolicy              `default:"none" json:"gatewayCertPolicy" validate:"oneof=none required verify"`
	OverseerCA        string                        `json:"overseerCA"`

	LogConfiguration logger.LoggerConfiguration `json:"LogConfiguration"`
}

//Load - Loads configuration from a file
func Load(path string) (*OverseerGatewayConfig, error) {

	config := new(OverseerGatewayConfig)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("unable to loadConfiguration from file")
	}
	err = json.Unmarshal(data, &config)

	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal confiuration file:%w", err)
	}

	return config, nil
}
