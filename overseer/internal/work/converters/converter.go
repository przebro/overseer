package converter

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"

	"github.com/przebro/overseer/common/types"

	any "google.golang.org/protobuf/types/known/anypb"
)

var converters map[types.TaskType]TaskActionConverter = make(map[types.TaskType]TaskActionConverter)

//ErrConverterNotRegistered - occurs when converter for given type is not registered
var ErrConverterNotRegistered error = errors.New("converter not registered for given type")

//TaskActionConverter - converts json raw data to any
type TaskActionConverter interface {
	ConvertToMsg(data json.RawMessage, variables types.EnvironmentVariableList) (*any.Any, error)
}

//RegisterConverter - registers converter for given task type
func RegisterConverter(taskType types.TaskType, c TaskActionConverter) {
	converters[taskType] = c
}

//ConvertToMsg - converts raw json to any
func ConvertToMsg(taskType types.TaskType, data json.RawMessage, variables types.EnvironmentVariableList) (*any.Any, error) {

	converter, ok := converters[taskType]
	if !ok {
		return nil, ErrConverterNotRegistered
	}

	return converter.ConvertToMsg(data, variables)
}

//ReplaceVariables - replace variables in input data
func ReplaceVariables(in string, variables types.EnvironmentVariableList) string {
	reg := regexp.MustCompile(`\%\%[A-Z0-9_]+`)

	out := in
	for _, n := range variables {
		if reg.MatchString(n.Name) {
			out = strings.Replace(out, n.Name, n.Value, -1)
		}
	}

	return out
}
