package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

//WorkerConfiguration - Holds a worker configuration
type WorkerConfiguration struct {
	ServiceName      string `json:"serviceName" validate:"required"`
	Name             string `json:"name" validate:"required"`
	Host             string `json:"host" validate:"ipv4,required"`
	Port             int    `json:"port" validate:"min=1024,max=65535,required"`
	SysoutDirectory  string `json:"sysoutDirectory" validate:"required"`
	RootDirectory    string
	ProcessDirectory string
}

//LogConfiguration - Holds a log configuration
type LogConfiguration struct {
	Directory string `json:"logDirectory"`
	Level     int    `json:"logLevel"`
}

//OverseerWorkerConfiguration - configuration
type OverseerWorkerConfiguration struct {
	Worker    WorkerConfiguration `json:"worker"`
	LogConfig LogConfiguration    `json:"logConfiguration"`
}

//Load - loads a configuration from file
func Load(path string) (OverseerWorkerConfiguration, error) {

	config := OverseerWorkerConfiguration{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, errors.New("Unable to load Configuration from file")
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

//GetLogConfiguration - Gets the log section
func (cfg OverseerWorkerConfiguration) GetLogConfiguration() LogConfiguration {
	return cfg.LogConfig

}
