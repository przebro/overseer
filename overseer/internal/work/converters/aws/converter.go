package aws

import (
	"encoding/json"
	"errors"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	converter "github.com/przebro/overseer/overseer/internal/work/converters"

	"google.golang.org/protobuf/proto"
	any "google.golang.org/protobuf/types/known/anypb"
)

func init() {

	converter.RegisterConverter(types.TypeAws, &awsConverter{})
}

type awsConverter struct {
}

// ConvertToMsg - converts aws task specific data to proto message
func (c *awsConverter) ConvertToMsg(data []byte, variables types.EnvironmentVariableList) (*any.Any, error) {

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
		return nil, errors.New("unknown connection type")
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
		return nil, errors.New("unknown action type")
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

func CreateMessage(object interface{}, variables types.EnvironmentVariableList) (*any.Any, error) {

	switch v := object.(type) {
	case *taskdef.AwsLambdaTaskData:
		return createLambdaMessage(v, variables)
	case *taskdef.AwsStepFunctionTaskData:
		return createStepFunctionMessage(v, variables)
	}

	return nil, errors.New("unknown type")
}

func createLambdaMessage(object *taskdef.AwsLambdaTaskData, variables types.EnvironmentVariableList) (*any.Any, error) {

	b := newBuilder(variables)

	if v, ok := object.IsConnection_AwsConnectionProperties(); ok {
		b.withConnectionProperties(v)
	} else if v, ok := object.IsConnection_String(); ok {
		b.withConnectionName(v)
	} else {
		return nil, errors.New("unknown connection type")
	}

	b.withLambda(*object)

	act := b.build()

	out, err := proto.Marshal(act)
	if err != nil {
		return nil, err
	}

	return &any.Any{TypeUrl: string(act.ProtoReflect().Descriptor().FullName()), Value: out}, nil
}

func createStepFunctionMessage(object *taskdef.AwsStepFunctionTaskData, variables types.EnvironmentVariableList) (*any.Any, error) {

	b := newBuilder(variables)

	if v, ok := object.IsConnection_AwsConnectionProperties(); ok {
		b.withConnectionProperties(v)

	} else if v, ok := object.IsConnection_String(); ok {
		b.withConnectionName(v)
	} else {
		return nil, errors.New("unknown connection type")
	}

	b.withStepFunction(*object)

	act := b.build()

	out, err := proto.Marshal(act)
	if err != nil {
		return nil, err
	}

	return &any.Any{TypeUrl: string(act.ProtoReflect().Descriptor().FullName()), Value: out}, nil
}
