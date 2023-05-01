package taskdef

import (
	"encoding/json"
)

// OsActionType - type of a possible os action
type OsActionType string

// Possible values for OsActionType
const (
	OsActionTypeCommand OsActionType = "command"
	OsActionTypeScript  OsActionType = "script"
)

type StepDefinition struct {
	Name    string `json:"name"`
	Command string `json:"exec"`
}

// OsTaskData - specific data for OS task
type OsTaskData struct {
	CommandLine string           `json:"command" yaml:"command"`
	RunAs       string           `json:"runas" yaml:"runas"`
	Steps       []StepDefinition `json:"steps" yaml:"steps"`
}

// OsTaskDefinition - definition of a OS task
type OsTaskDefinition struct {
	TaskDefinition
	Spec OsTaskData
}

type AWSActionType string

const (
	AWSActionTypeLambda   AWSActionType = "lambda"
	AWSActionTypeStepFunc AWSActionType = "stepfunc"
	AWSActionTypeBatch    AWSActionType = "batch"
)

type AwsConnectionConfig interface{}

type AwsConnectionProperties struct {
	Profile string `json:"profile,omitempty"`
	Region  string `json:"region,omitempty"`
}

type AwsTaskData struct {
	Type       AWSActionType `json:"type"`
	Connection interface{}   `json:"connection"`
	Payload    interface{}   `json:"payload"`
}

// IsConnection_String - checks if this holds Connection as string
func (p *AwsTaskData) IsConnection_String() (string, bool) {
	if v, ok := p.Connection.(string); ok {
		return v, true
	}
	return "", false
}

// IsConnection_AwsConnectionProperties - checks if this holds Connection as an instance of AwsConnectionProperties
func (p *AwsTaskData) IsConnection_AwsConnectionProperties() (AwsConnectionProperties, bool) {

	if v, ok := p.Connection.(AwsConnectionProperties); ok {
		return v, true
	}
	return AwsConnectionProperties{}, false
}

type awsTaskData AwsTaskData

// UnmarshalJSON - parses []byte and stores value in current AwsTaskData instance
func (p *AwsTaskData) UnmarshalJSON(b []byte) error {

	var taskData awsTaskData
	if err := json.Unmarshal(b, &taskData); err != nil {
		return err
	}

	switch value := taskData.Connection.(type) {
	case string:
		{
			p.Connection = value
		}
	case map[string]interface{}:
		{
			out, _ := json.Marshal(value)

			connProp := AwsConnectionProperties{}

			if err := json.Unmarshal(out, &connProp); err != nil {
				return UnknownTypeError{}
			}

			p.Connection = connProp

		}
	default:
		return UnknownTypeError{}
	}

	switch value := taskData.Payload.(type) {
	case string:
		{
			p.Payload = value
		}
	case map[string]interface{}:
		{
			out, _ := json.Marshal(value)
			if ok := json.Valid(out); !ok {
				return UnknownTypeError{}
			}
			p.Payload = json.RawMessage(out)
		}
	}

	p.Type = taskData.Type

	return nil
}

type AwsLambdaTaskData struct {
	AwsTaskData
	FunctionName  string `json:"functionName"`
	FunctionAlias string `json:"alias"`
}

func (p *AwsLambdaTaskData) UnmarshalJSON(b []byte) error {

	inline := struct {
		FunctionName  string `json:"functionName"`
		FunctionAlias string `json:"alias"`
	}{}

	if err := json.Unmarshal(b, &inline); err != nil {
		return err
	}
	p.FunctionName = inline.FunctionName
	p.FunctionAlias = inline.FunctionAlias

	return nil
}

type AwsStepFunctionTaskData struct {
	AwsTaskData
	StateMachine  string `json:"stateMachineARN"`
	ExecutionName string `json:"executionName"`
}

func (p *AwsStepFunctionTaskData) UnmarshalJSON(b []byte) error {

	inline := struct {
		StateMachine  string `json:"stateMachineARN"`
		ExecutionName string `json:"executionName"`
	}{}

	if err := json.Unmarshal(b, &inline); err != nil {
		return err
	}

	p.StateMachine = inline.StateMachine
	p.ExecutionName = inline.ExecutionName

	return nil
}

type UnknownTypeError struct {
}

func (e UnknownTypeError) Error() string {
	return "failed to unmarshal struct, unknown type"
}
