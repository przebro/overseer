package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/common/types"
)

// WorkerConfiguration - configuration
type WorkerConfiguration struct {
	WorkerName string `json:"name" validate:"required"`
	WorkerHost string `json:"workerHost" validate:"ipv4,required"`
	WorkerPort int    `json:"workerPort" validate:"min=1024,max=65535,required"`
	WorkerCA   string `json:"workerCA"`
}

// WorkerManagerConfiguration - setting for worker manager
type WorkerManagerConfiguration struct {
	Timeout           int                   `json:"timeout"`
	WorkerInterval    int                   `json:"interval"`
	WorkerMaxAttempts int                   `json:"attempts"`
	Workers           []WorkerConfiguration `json:"workers"`
}

// ActivePoolConfiguration - Active Pool Configuration section
type ActivePoolConfiguration struct {
	MaxOkReturnCode int32             `json:"maxOkReturnCode"`
	NewDayProc      types.HourMinTime `json:"newDayProc" validate:"hmtime,required"`
	ForceNewDayProc bool              `json:"forceNewDayProc"`
	SyncTime        int               `json:"syncTime" validate:"min=0,max=60"`
}

// ResourceEntry - resource configuration
type ResourceEntry struct {
	Sync int `json:"sync"`
}

// ResourcesConfigurartion - configuration section for tickets and flags
type ResourcesConfigurartion struct {
	Resources ResourceEntry `json:"resources"`
}

// IntervalValue - represents limited interval value
type IntervalValue int

// ServerConfiguration - main server parameters
type ServerConfiguration struct {
	ProcessDirectory string
	RootDirectory    string
	ServiceName      string                      `json:"serviceName" validate:"required"`
	Host             string                      `json:"ovshost" validate:"ipv4,required"`
	Port             int                         `json:"ovsport" validate:"min=1024,max=65535,required"`
	Security         ServerSecurityConfiguration `json:"security"`
}

type ServerSecurityConfiguration struct {
	SecurityLevel    types.ConnectionSecurityLevel `json:"securityLevel" validate:"oneof=none server clientandserver"`
	ServerCert       string                        `json:"cert"`
	ServerKey        string                        `json:"key"`
	ClientCertPolicy types.CertPolicy              `json:"clientCertPolicy" validate:"oneof=none required verify"`
}

// OverseerConfiguration - main configuration
type OverseerConfiguration struct {
	Server              ServerConfiguration        `json:"serverConfiguration"`
	DefinitionDirectory string                     `json:"definitionDirectory"`
	Resources           ResourcesConfigurartion    `json:"ResourceConfiguration"`
	Log                 logger.LoggerConfiguration `json:"LogConfiguration"`
	PoolConfiguration   ActivePoolConfiguration    `json:"ActivePoolConfiguration"`
	TimeInterval        IntervalValue              `json:"timeInterval" validate:"min=1,max=60"`
	StoreConfiguration  StoreConfiguration         `json:"StoreConfiguration"`
	Security            SecurityConfiguration      `json:"security"`
	WorkerManager       WorkerManagerConfiguration `json:"WorkerConfiguration"`
	Journal             JournalConfiguration       `json:"journalConfiguration"`
}

// StoreConfiguration - store configuration entry
type StoreConfiguration struct {
	ID               string `json:"id"`
	ConnectionString string `json:"connectionString"`
}

// CollectionConfiguration - collection configuration entry
type CollectionConfiguration struct {
	StoreID string `json:"storeId"`
	Name    string `json:"name"`
}

// SecurityConfiguration - security section
type SecurityConfiguration struct {
	AllowAnonymous bool          `json:"allowAnonymous"`
	Timeout        int           `json:"timeout"`
	Issuer         string        `json:"issuer"`
	Secret         string        `json:"secret"`
	Providers      []interface{} `json:"providers"`
}

// JournalConfiguration - task history section
type JournalConfiguration struct {
	SyncTime int `json:"syncTime" validate:"min=0,max=60"`
}

// Load - Loads configuration from a file
func Load(path string) (OverseerConfiguration, error) {

	config := OverseerConfiguration{
		Server: ServerConfiguration{
			Security: ServerSecurityConfiguration{
				ClientCertPolicy: types.CertPolicyNone,
				SecurityLevel:    types.ConnectionSecurityLevelNone,
			},
		},
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return config, errors.New("unable to load configuration from file")
	}
	err = json.Unmarshal(data, &config)

	if err != nil {
		return config, fmt.Errorf("unable to unmarshal confiuration file:%w", err)
	}

	return config, nil
}
