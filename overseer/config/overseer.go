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
	WorkerName string `json:"name" validate:"required"`
	WorkerHost string `json:"workerHost" validate:"ipv4,required"`
	WorkerPort int    `json:"workerPort" validate:"min=1024,max=65535,required"`
}

//WorkerManagerConfiguration - setting for worker manager
type WorkerManagerConfiguration struct {
	Timeout           int                   `json:"timeout"`
	WorkerInterval    int                   `json:"interval"`
	WorkerMaxAttempts int                   `json:"attempts"`
	Workers           []WorkerConfiguration `json:"workers"`
}

//LogConfiguration - configuration for logger
type LogConfiguration struct {
	LogLevel     int    `json:"logLevel"`
	LogDirectory string `json:"logDirectory"`
}

//ActivePoolConfiguration - Active Pool Configuration section
type ActivePoolConfiguration struct {
	MaxOkReturnCode int32             `json:"maxOkReturnCode"`
	NewDayProc      types.HourMinTime `json:"newDayProc" validate:"hmtime,required"`
	ForceNewDayProc bool              `json:"forceNewDayProc"`
	Collection      string            `json:"collection"`
	SyncTime        int               `json:"syncTime" validate:"min=0,max=60"`
}

//ResourceEntry - resource configuration
type ResourceEntry struct {
	Collection string `json:"collectionName"`
	Sync       int    `json:"sync"`
}

//ResourcesConfigurartion - configuration section for tickets and flags
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
	ServiceName      string `json:"serviceName" validate:"required"`
	Host             string `json:"ovshost" validate:"ipv4,required"`
	Port             int    `json:"ovsport" validate:"min=1024,max=65535,required"`
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
	Journal             JournalConfiguration       `json:"journalConfiguration"`
}

//StoreProviderConfiguration - datastore configuration section
type StoreProviderConfiguration struct {
	Store       []StoreConfiguration      `json:"store"`
	Collections []CollectionConfiguration `json:"collections"`
}

//StoreConfiguration - store configuration entry
type StoreConfiguration struct {
	ID               string `json:"id"`
	ConnectionString string `json:"connectionString"`
}

//CollectionConfiguration - collection configuration entry
type CollectionConfiguration struct {
	StoreID string `json:"storeId"`
	Name    string `json:"name"`
}

//SecurityConfiguration - security section
type SecurityConfiguration struct {
	AllowAnonymous bool          `json:"allowAnonymous"`
	Timeout        int           `json:"timeout"`
	Issuer         string        `json:"issuer"`
	Secret         string        `json:"secret"`
	Collection     string        `json:"collectionName"`
	Providers      []interface{} `json:"providers"`
}

//JournalConfiguration - task history section
type JournalConfiguration struct {
	LogCollection string `json:"logs"`
	SyncTime      int    `json:"syncTime" validate:"min=0,max=60"`
}

//Load - Loads configuration from a file
func Load(path string) (OverseerConfiguration, error) {

	config := OverseerConfiguration{}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, errors.New("Unable to loadConfiguration from file")
	}
	err = json.Unmarshal(data, &config)

	if err != nil {
		return config, fmt.Errorf("unable to unmarshal confiuration file:%w", err)
	}

	return config, nil
}

//GetServerConfiguration - gets main configuration section
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

//GetResourceConfiguration - Gets  resource configuration section
func (cfg *OverseerConfiguration) GetResourceConfiguration() ResourcesConfigurartion {
	return cfg.Resources
}

//GetStoreProviderConfiguration - gets provider configuration section
func (cfg *OverseerConfiguration) GetStoreProviderConfiguration() StoreProviderConfiguration {
	return cfg.StoreProvider
}

//GetSecurityConfiguration - gets security configuration section
func (cfg *OverseerConfiguration) GetSecurityConfiguration() SecurityConfiguration {
	return cfg.Security
}

//GetWorkerManagerConfiguration - gets worker manager configuration section
func (cfg *OverseerConfiguration) GetWorkerManagerConfiguration() WorkerManagerConfiguration {
	return cfg.WorkerManager
}

//GetJournalConfiguration - gets journal configuration section
func (cfg *OverseerConfiguration) GetJournalConfiguration() JournalConfiguration {
	return cfg.Journal
}
