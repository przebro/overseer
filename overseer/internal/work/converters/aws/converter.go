package aws

import (
	"encoding/json"
	"errors"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	converter "github.com/przebro/overseer/overseer/internal/work/converters"

	"github.com/golang/protobuf/ptypes/any"
	"google.golang.org/protobuf/proto"
)

func init() {

	converter.RegisterConverter(types.TypeAws, &awsConverter{})
}

type awsConverter struct {
}

//ConvertToMsg - converts aws task specific data to proto message
func (c *awsConverter) ConvertToMsg(data json.RawMessage, variables []taskdef.VariableData) (*any.Any, error) {

	result := &taskdef.AwsTaskData{}

	if err := json.Unmarshal(data, result); err != nil {
		return nil, err
	}

	b := newBuilder(variables)

	if v, ok := result.IsConnection_AwsConnectionProperties(); ok {

		b.withConnectionProperties(v)

	} else if v, ok := result.IsConnection_String(); ok {

		b.withConnectionName(v)

	} else {
		return nil, errors.New("unkown connection type")
	}

	if result.Type == taskdef.AWSActionTypeLambda {

		val := taskdef.AwsLambdaTaskData{}
		if err := json.Unmarshal(data, &val); err != nil {
			return nil, err
		}
		b.withLambda(val)

	} else if result.Type == taskdef.AWSActionTypeStepFunc {

		val := taskdef.AwsStepFunctionTaskData{}
		if err := json.Unmarshal(data, &val); err != nil {
			return nil, err
		}
		b.withStepFunction(val)

	} else {
		return nil, errors.New("unkown action type")
	}

	if v, ok := result.Payload.(string); ok {

		b.withPayloadFilePath(v)

	} else if v, ok := result.Payload.(json.RawMessage); ok {

		b.withPayloadStream(v)

	} else {
		return nil, errors.New("unknown payload type")
	}

	act := b.build()

	out, err := proto.Marshal(act)
	if err != nil {
		return nil, err
	}

	return &any.Any{TypeUrl: string(act.ProtoReflect().Descriptor().FullName()), Value: out}, nil
}
