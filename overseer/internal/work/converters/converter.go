package converter

import (
	"encoding/json"
	"errors"
	"overseer/common/types"
	"overseer/overseer/internal/taskdef"
	"regexp"
	"strings"

	"github.com/golang/protobuf/ptypes/any"
)

var converters map[types.TaskType]TaskActionConverter = make(map[types.TaskType]TaskActionConverter)

var ErrConverterNotRegistered error = errors.New("converter not registered for given type")

type TaskActionConverter interface {
	ConvertToMsg(data json.RawMessage, variables []taskdef.VariableData) (*any.Any, error)
}

func RegisterConverter(taskType types.TaskType, c TaskActionConverter) {
	converters[taskType] = c
}

func ConvertToMsg(taskType types.TaskType, data json.RawMessage, variables []taskdef.VariableData) (*any.Any, error) {

	converter, ok := converters[taskType]
	if !ok {
		return nil, ErrConverterNotRegistered
	}

	return converter.ConvertToMsg(data, variables)
}

func ReplaceVariables(in string, variables []taskdef.VariableData) string {
	reg := regexp.MustCompile(`\%\%[A-Z0-9_]+`)

	out := in
	for _, n := range variables {
		if reg.MatchString(n.Name) {
			out = strings.Replace(out, n.Name, n.Value, -1)
		}
	}

	return out
}
