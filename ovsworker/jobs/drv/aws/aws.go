package aws

import (
	"context"
	"errors"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/ovsworker/jobs"
	"github.com/przebro/overseer/ovsworker/msgheader"
	"github.com/przebro/overseer/ovsworker/status"
	"github.com/przebro/overseer/proto/actions"
	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"google.golang.org/protobuf/proto"
)

func init() {
	jobs.RegisterFactory(types.TypeAws, AwsJobFactory)
}

type awsPayloadReader interface {
	Read() ([]byte, error)
}

type awsServiceCaller interface {
	jobs.JobExecutor
	Call(ctx context.Context, stat chan<- status.JobExecutionStatus) status.JobExecutionStatus
}

// AwsJobFactory - Creates a new aws factory
func AwsJobFactory(ctx context.Context, header msgheader.TaskHeader, sysoutDir string, data []byte) (jobs.JobExecutor, error) {

	log := zerolog.Ctx(ctx).With().Str("taskID", header.TaskID).Str("executionID", header.ExecutionID).Logger()

	act := actions.AwsTaskAction{}
	if err := proto.Unmarshal(data, &act); err != nil {
		log.Error().Err(err).Msg("AwsJobFactory")
		return nil, err
	}

	return newAwsJob(ctx, header, sysoutDir, &act)
}

// newDummyJob - factory method
func newAwsJob(ctx context.Context, header msgheader.TaskHeader, sysoutDir string, action *actions.AwsTaskAction) (jobs.JobExecutor, error) {

	var conf aws.Config
	var err error
	var payloadReader awsPayloadReader
	var caller awsServiceCaller

	log := zerolog.Ctx(ctx)

	job := jobs.Job{

		TaskID:      header.TaskID,
		ExecutionID: header.ExecutionID,
		Start:       make(chan status.JobExecutionStatus),
		Variables:   make(map[string]string),
		SysoutDir:   sysoutDir,
	}
	for k, v := range header.Variables {
		job.Variables[k] = v
	}

	switch act := action.GetConnection().(type) {

	case *actions.AwsTaskAction_ConnectionData:
		{
			conf, err = createConfig(context.Background(), act.ConnectionData.ProfileName, act.ConnectionData.Region)
			if err != nil {
				log.Error().Err(err).Msg("create config")
				return nil, err
			}
		}
	case *actions.AwsTaskAction_ConnectionProfileName:
		{
			err = errors.New("not implemented")
		}
	}

	if err != nil {
		return nil, err
	}

	switch src := action.GetPayloadSource().(type) {

	case *actions.AwsTaskAction_PayloadRaw:
		{
			payloadReader = newStreamReader(src.PayloadRaw)
		}
	case *actions.AwsTaskAction_PayloadFilePath:
		{
			payloadReader = newFileReader(src.PayloadFilePath)
		}
	}

	switch action.Type {
	case actions.AwsTaskAction_lambda:
		{
			lExecData := action.GetLambdaExecution()
			if lExecData == nil {
				err := errors.New("failed to get lambda execution data")
				log.Error().Err(err).Msg("create caller")
				return nil, err
			}

			caller = newLambdaCaller(conf, job, lExecData.FunctionName, lExecData.Alias, payloadReader)
			log.Info().Str("function", lExecData.FunctionName).Str("alias", lExecData.Alias).Msg("caller created")

		}
	case actions.AwsTaskAction_stepfunc:
		{
			sExecData := action.GetStepFunction()
			if sExecData == nil {
				err := errors.New("failed to get stepfunction execution data")
				log.Error().Err(err).Msg("create caller")
				return nil, err
			}
			caller = newStepFunctionCaller(conf, job, sExecData.StateMachineARN, sExecData.ExecutionName, payloadReader)
			log.Info().Str("execution_name", sExecData.ExecutionName).
				Str("machine_arn", sExecData.StateMachineARN).Msg("caller created")

		}
	default:
		{
			caller = nil
		}
	}

	if caller == nil {
		err := errors.New("failed to initialize caller, unknown type")
		log.Error().Err(err).Msg("create caller")
		return nil, err
	}

	return caller, nil
}

func createConfig(ctx context.Context, profile, region string) (aws.Config, error) {

	opts := []func(opt *config.LoadOptions) error{}

	if profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	if region != "" {
		opts = append(opts, config.WithRegion(region))
	}

	return config.LoadDefaultConfig(ctx, opts...)
}
