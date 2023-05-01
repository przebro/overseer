package types

import "strings"

type EnvironmentVariable struct {
	Name  string `json:"name" yaml:"name" validate:"required,max=32,varname"`
	Value string `json:"value" yaml:"value"`
}

// Expand - Expands name
func (data EnvironmentVariable) Expand() string {
	return strings.Replace(data.Name, "%%", "OVS_", 1)
}

type EnvironmentVariableList []EnvironmentVariable
