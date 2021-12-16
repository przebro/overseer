package os

import (
	"encoding/json"
	"overseer/common/types"
	"overseer/overseer/internal/taskdef"
	converter "overseer/overseer/internal/work/converters"
	"overseer/proto/actions"
	"regexp"
	"strings"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/proto"
)

func init() {

	converter.RegisterConverter(types.TypeDummy, &osConverter{})
}

type osConverter struct {
}

//ConvertToMsg - converts os specific data to proto message
func (c *osConverter) ConvertToMsg(data json.RawMessage, variables []taskdef.VariableData) (*any.Any, error) {

	reg := regexp.MustCompile(`\%\%[A-Z0-9_]+`)
	result := &taskdef.OsTaskData{}
	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}
	cmdLine := result.CommandLine
	for _, n := range variables {
		if reg.MatchString(n.Name) {
			cmdLine = strings.Replace(cmdLine, n.Name, n.Value, -1)
		}
	}

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
