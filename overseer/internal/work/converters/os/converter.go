package os

import (
	"encoding/json"

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

//ConvertToMsg - converts os specific data to proto message
func (c *osConverter) ConvertToMsg(data json.RawMessage, variables []taskdef.VariableData) (*any.Any, error) {

	result := &taskdef.OsTaskData{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	cmdLine := converter.ReplaceVariables(result.CommandLine, variables)

	var taskType actions.OsTaskAction_OsType

	if result.ActionType == taskdef.OsActionTypeCommand {
		taskType = actions.OsTaskAction_command
	} else {
		taskType = actions.OsTaskAction_script
	}

	cmd := &actions.OsTaskAction{CommandLine: cmdLine, Runas: result.RunAs, Type: taskType}
	act, err := proto.Marshal(cmd)
	if err != nil {
		return nil, err
	}

	return &any.Any{TypeUrl: string(cmd.ProtoReflect().Descriptor().FullName()), Value: act}, nil
}
