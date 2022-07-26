package aws

import (
	"encoding/json"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	converter "github.com/przebro/overseer/overseer/internal/work/converters"
	"github.com/przebro/overseer/proto/actions"
)

type awsTaskDatabuilder interface {
	build() *actions.AwsTaskAction
}

type taskDataBuilder struct {
	action *actions.AwsTaskAction
	v      types.EnvironmentVariableList
}

func newBuilder(v types.EnvironmentVariableList) *taskDataBuilder {
	return &taskDataBuilder{action: &actions.AwsTaskAction{}, v: v}
}

func (b *taskDataBuilder) build() *actions.AwsTaskAction {
	return b.action
}
func (b *taskDataBuilder) withLambda(lambda taskdef.AwsLambdaTaskData) *taskDataBuilder {

	b.action.Type = actions.AwsTaskAction_lambda
	b.action.Execution = &actions.AwsTaskAction_LambdaExecution{
		LambdaExecution: &actions.AwsLambdaExecution{
			FunctionName: lambda.FunctionName,
			Alias:        lambda.FunctionAlias,
		},
	}

	return b
}

func (b *taskDataBuilder) withStepFunction(stepfunc taskdef.AwsStepFunctionTaskData) *taskDataBuilder {

	b.action.Type = actions.AwsTaskAction_stepfunc

	execName := converter.ReplaceVariables(stepfunc.ExecutionName, b.v)

	b.action.Execution = &actions.AwsTaskAction_StepFunction{
		StepFunction: &actions.AwsStepFunctionExecution{
			StateMachineARN: stepfunc.StateMachine,
			ExecutionName:   execName,
		},
	}

	return b
}

func (b *taskDataBuilder) withConnectionProperties(props taskdef.AwsConnectionProperties) *taskDataBuilder {
	b.action.Connection = &actions.AwsTaskAction_ConnectionData{
		ConnectionData: &actions.AwsConnetionData{
			ProfileName: props.Profile,
			Region:      props.Region,
		},
	}
	return b
}

func (b *taskDataBuilder) withConnectionName(connectionName string) *taskDataBuilder {
	b.action.Connection = &actions.AwsTaskAction_ConnectionProfileName{
		ConnectionProfileName: connectionName,
	}

	return b
}

func (b *taskDataBuilder) withPayloadStream(stream json.RawMessage) *taskDataBuilder {

	b.action.PayloadSource = &actions.AwsTaskAction_PayloadRaw{
		PayloadRaw: stream,
	}

	return b
}

func (b *taskDataBuilder) withPayloadFilePath(filepath string) *taskDataBuilder {

	b.action.PayloadSource = &actions.AwsTaskAction_PayloadFilePath{
		PayloadFilePath: filepath,
	}

	return b
}
