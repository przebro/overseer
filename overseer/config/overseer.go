package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"overseer/common/types"
)

//WorkerConfiguration - configuration
type WorkerConfiguration struct {
	WorkerName string `json:"name"`
	WorkerHost string `json:"workerHost"`
	WorkerPort int    `json:"workerPort"`
}

//WorkerManagerConfiguration - setting for worker manager
type WorkerManagerConfiguration struct {
	Timeout          int                   `json:"timeout"`
	WorkerInterval   int                   `json:"interval"`
	WorkerMaxAttemps int                   `json:"attemps"`
	Workers          []WorkerConfiguration `json:"workers"`
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
	Collection      string            `json:"collection"`
	SyncTime        int               `json:"syncTime"`
}
type ResourceEntry struct {
	Collection string `json:"collectionName"`
	Sync       int    `json:"sync"`
}

type ResourcesConfigurartion struct {
	TicketSource ResourceEntry `json:"tickets"`
	FlagSource   ResourceEntry `json:"flags"`
}

//IntervalValue - represents limited interval value
type IntervalValue int

//ServerConfiguration - main server parameters
type ServerConfiguration struct {
	ProcessDirectory string
	RootDirectory    string
	Host             string `json:"ovshost"`
	Port             int    `json:"ovsport"`
	TLS              bool   `json:"tls"`
	ServerCert       string `json:"cert"`
	ServerKey        string `json:"key"`
}

//OverseerConfiguration - main configuration
type OverseerConfiguration struct {
	Server              ServerConfiguration        `json:"serverConfiguration"`
	DefinitionDirectory string                     `json:"definitionDirectory"`
	Resources           ResourcesConfigurartion    `json:"ResourceConfiguration"`
	Log                 LogConfiguration           `json:"LogConfiguration"`
	PoolConfiguration   ActivePoolConfiguration    `json:"ActivePoolConfiguration"`
	TimeInterval        IntervalValue              `json:"timeInterval" validate:"min=1,max=60"`
	StoreProvider       StoreProviderConfiguration `json:"StoreProvider"`
	Security            SecurityConfiguration      `json:"security"`
	WorkerManager       WorkerManagerConfiguration `json:"WorkerConfiguration"`
}

type StoreProviderConfiguration struct {
	Store       []StoreConfiguration      `json:"store"`
	Collections []CollectionConfiguration `json:"collections"`
}

type StoreConfiguration struct {
	ID               string `json:"id"`
	ConnectionString string `json:"connectionString"`
}

type CollectionConfiguration struct {
	StoreID string `json:"storeId"`
	Name    string `json:"name"`
}

type SecurityConfiguration struct {
	AllowAnonymous bool          `json:"allowAnonymous"`
	Timeout        int           `json:"timeout"`
	Issuer         string        `json:"issuer"`
	Secret         string        `json:"secret"`
	Collection     string        `json:"collectionName"`
	Providers      []interface{} `json:"providers"`
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
		return nil, fmt.Errorf("unable to unmarshal confiuration file:%w", err)
	}

	return config, nil
}

func (cfg *OverseerConfiguration) GetServerConfiguration() ServerConfiguration {

	return cfg.Server
}

//GetLogConfiguration - Gets a log configuration section
func (cfg *OverseerConfiguration) GetLogConfiguration() LogConfiguration {
	return cfg.Log
}

//GetActivePoolConfiguration - Gets an Active Pool configuration section
func (cfg *OverseerConfiguration) GetActivePoolConfiguration() ActivePoolConfiguration {
	return cfg.PoolConfiguration
}

func (cfg *OverseerConfiguration) GetResourceConfiguration() ResourcesConfigurartion {
	return cfg.Resources
}

func (cfg *OverseerConfiguration) GetStoreProviderConfiguration() StoreProviderConfiguration {
	return cfg.StoreProvider
}

func (cfg *OverseerConfiguration) GetSecurityConfiguration() SecurityConfiguration {
	return cfg.Security
}

func (cfg *OverseerConfiguration) GetWorkerManagerConfiguration() WorkerManagerConfiguration {
	return cfg.WorkerManager
}
