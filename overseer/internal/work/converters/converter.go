package converter

import (
	"encoding/json"
	"errors"
	"overseer/common/types"
	"overseer/overseer/internal/taskdef"

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
