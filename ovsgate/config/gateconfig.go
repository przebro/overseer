package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"overseer/common/logger"
)

type OverseerGatewayConfig struct {
	GatewayAddress   string                     `json:"gateAddress"`
	GatewayPort      int                        `json:"gatePort"`
	OverseerAddress  string                     `json:"overseerAddress"`
	OverseerPort     int                        `json:"overseerPort"`
	UseTLS           bool                       `json:"tls"`
	CertPath         string                     `json:"cert"`
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
