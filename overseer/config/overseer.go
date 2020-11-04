package config

import (
	"encoding/json"
	"errors"
	"goscheduler/common/types"
	"io/ioutil"
)

//WorkerConfiguration - configuration
type WorkerConfiguration struct {
	WorkerName string `json:"name"`
	WorkerHost string `json:"workerHost"`
	WorkerPort int    `json:"workerPort"`
}

//LogConfiguration - configuration for logger
type LogConfiguration struct {
	LogLevel     int    `json:"logLevel"`
	LogDirectory string `json:"logDirectory"`
}

//ActivePoolConfiguration - Active Pool Configuration section
type ActivePoolConfiguration struct {
	MaxOkReturnCode int32             `json:"maxOkReturnCode"`
	NewDayProc      types.HourMinTime `json:"newDayProc"`
	ForceNewDayProc bool              `json:"forceNewDayProc"`
}

//IntervalValue - represents limited interval value
type IntervalValue int

//OverseerConfiguration - main configuration
type OverseerConfiguration struct {
	ProcessDirectory    string
	RootDirectory       string
	Host                string                  `json:"ovshost"`
	Port                int                     `json:"ovsport"`
	DefinitionDirectory string                  `json:"definitionDirectory"`
	ResourceFilePath    string                  `json:"resourceDirectory"`
	Log                 LogConfiguration        `json:"LogConfiguration"`
	PoolConfiguration   ActivePoolConfiguration `json:"ActivePoolConfiguration"`
	TimeInterval        IntervalValue           `json:"timeInterval" validate:"min=1,max=60"`
	WorkerTimeout       int
	WorkerMaxAttemps    int
	Workers             []WorkerConfiguration `json:"workers"`
}

//Load - Loads configuration from a file
func Load(path string) (*OverseerConfiguration, error) {

	config := new(OverseerConfiguration)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.New("Unable to loadConfiguration from file")
	}
	err = json.Unmarshal(data, &config)

	if err != nil {
		return nil, errors.New("unable to unmarshal confiuration file")
	}

	return config, nil
}

//GetLogConfiguration - Gets a log configuration section
func (cfg *OverseerConfiguration) GetLogConfiguration() LogConfiguration {
	return cfg.Log
}

//GetActivePoolConfiguration - Gets an Active Pool configuration section
func (cfg *OverseerConfiguration) GetActivePoolConfiguration() ActivePoolConfiguration {
	return cfg.PoolConfiguration
}
