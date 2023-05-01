package os

import (
	"encoding/json"
	"errors"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	converter "github.com/przebro/overseer/overseer/internal/work/converters"
	"github.com/przebro/overseer/proto/actions"

	"google.golang.org/protobuf/proto"
	any "google.golang.org/protobuf/types/known/anypb"
)

func init() {

	converter.RegisterConverter(types.TypeOs, &osConverter{})
}

type osConverter struct {
}

// ConvertToMsg - converts os specific data to proto message
func (c *osConverter) ConvertToMsg(data []byte, variables types.EnvironmentVariableList) (*any.Any, error) {

	result := &taskdef.OsTaskData{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	cmdLine := converter.ReplaceVariables(result.CommandLine, variables)

	steps := make([]*actions.OsStepDefinition, len(result.Steps))

	for _, step := range result.Steps {
		steps = append(steps, &actions.OsStepDefinition{StepName: step.Name, Command: step.Command})
	}

	cmd := &actions.OsTaskAction{CommandLine: cmdLine, Runas: result.RunAs, Steps: steps}
	act, err := proto.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	return &any.Any{TypeUrl: string(cmd.ProtoReflect().Descriptor().FullName()), Value: act}, nil
}

func CreateMessage(object interface{}, variables types.EnvironmentVariableList) (*any.Any, error) {

	data, ok := object.(*taskdef.OsTaskData)
	if !ok {
		return nil, errors.New("invalid data type")
	}

	cmdLine := converter.ReplaceVariables(data.CommandLine, variables)

	steps := make([]*actions.OsStepDefinition, len(data.Steps))

	for _, step := range data.Steps {
		steps = append(
			steps, &actions.OsStepDefinition{
				StepName: step.Name,
				Command:  converter.ReplaceVariables(step.Command, variables),
			})
	}

	cmd := &actions.OsTaskAction{
		CommandLine: cmdLine,
		Runas:       data.RunAs,
		Steps:       steps,
	}

	act, err := proto.Marshal(cmd)

	if err != nil {
		return nil, err
	}

	return &any.Any{TypeUrl: string(cmd.ProtoReflect().Descriptor().FullName()), Value: act}, nil
}
