package dummy

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
