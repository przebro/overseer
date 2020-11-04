package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

//WorkerConfiguration - Holds a worker configuration
type WorkerConfiguration struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port int    `json:"port"`
}

//LogConfiguration - Holds a log configuration
type LogConfiguration struct {
	Directory string `json:"logDirectory"`
	Level     int    `json:"logLevel"`
}

//Config - configuration structure
type Config struct {
	Worker    WorkerConfiguration `json:"worker"`
	LogConfig LogConfiguration    `json:"logConfiguration"`
}

//Load - loads a configuration from file
func Load(path string) (*Config, error) {

	config := &Config{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Unable to loadConfiguration from file")
	}
	err = json.Unmarshal(data, &config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

//GetLogConfiguration - Gets the log section
func (cfg *Config) GetLogConfiguration() LogConfiguration {
	return cfg.LogConfig

}
