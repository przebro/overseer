package dummy

import (
	"encoding/json"
	"overseer/common/types"
	"overseer/overseer/internal/taskdef"
	converter "overseer/overseer/internal/work/converters"
	"overseer/proto/actions"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/proto"
)

func init() {

	converter.RegisterConverter(types.TypeDummy, &dummyConverter{})
}

type dummyConverter struct {
}

func (c *dummyConverter) ConvertToMsg(data json.RawMessage, variables []taskdef.VariableData) (*any.Any, error) {

	cmd := &actions.DummyTaskAction{Data: ""}
	act, err := proto.Marshal(cmd)
	if err != nil {
		return nil, err
	}
	return &any.Any{TypeUrl: string(cmd.ProtoReflect().Descriptor().FullName()), Value: act}, nil
}
