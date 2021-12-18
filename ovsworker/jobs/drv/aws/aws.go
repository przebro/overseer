package aws

import (
	"context"
	"errors"
	"overseer/common/logger"
	"overseer/common/types"
	"overseer/ovsworker/jobs"
	"overseer/ovsworker/msgheader"
	"overseer/ovsworker/status"
	"overseer/proto/actions"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"go.uber.org/zap"
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
	Call(ctx context.Context, stat chan<- status.JobExecutionStatus)
}

//AwsJobFactory - Creates a new aws factory
func AwsJobFactory(header msgheader.TaskHeader, sysoutDir string, data []byte, log logger.AppLogger) (jobs.JobExecutor, error) {

	act := actions.AwsTaskAction{}
	if err := proto.Unmarshal(data, &act); err != nil {
		log.Desugar().Error("AwsJobFactory", zap.String("error", err.Error()))
		return nil, err
	}

	return newAwsJob(header, sysoutDir, &act, log)
}

//newDummyJob - factory method
func newAwsJob(header msgheader.TaskHeader, sysoutDir string, action *actions.AwsTaskAction, log logger.AppLogger) (jobs.JobExecutor, error) {

	var conf aws.Config
	var err error
	var payloadReader awsPayloadReader
	var caller awsServiceCaller

	job := jobs.Job{

		TaskID:      header.TaskID,
		ExecutionID: header.ExecutionID,
		Start:       make(chan status.JobExecutionStatus),
		Variables:   make(map[string]string),
		SysoutDir:   sysoutDir,
		Log:         log,
	}
	for k, v := range header.Variables {
		job.Variables[k] = v
	}

	switch act := action.GetConnection().(type) {

	case *actions.AwsTaskAction_ConnectionData:
		{
			conf, err = createConfig(context.Background(), act.ConnectionData.ProfileName, act.ConnectionData.Region)
			if err != nil {
				log.Desugar().Error("newAwsJob", zap.String("error", err.Error()))
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
				log.Desugar().Error("create caller", zap.String("error", "failed to get lambda execution data"))
				return nil, errors.New("failed to get lambda execution data")
			}

			caller = newLambdaCaller(conf, job, lExecData.FunctionName, lExecData.Alias, payloadReader)
			log.Desugar().Info("create caller", zap.String("created", "lambda"),
				zap.String("function", lExecData.FunctionName),
				zap.String("alias", lExecData.Alias))
		}
	case actions.AwsTaskAction_stepfunc:
		{
			sExecData := action.GetStepFunction()
			if sExecData == nil {
				log.Desugar().Error("create caller", zap.String("error", "failed to get stepfunction execution data"))
				return nil, errors.New("failed to get stepfunction execution data")
			}
			caller = newStepFunctionCaller(conf, job, sExecData.StateMachineARN, sExecData.ExecutionName, payloadReader)
			log.Desugar().Info("create caller", zap.String("created", "stepfunc"),
				zap.String("machineARN", sExecData.StateMachineARN),
				zap.String("executionName", sExecData.ExecutionName))

		}
	default:
		{
			caller = nil
		}
	}

	if caller == nil {
		log.Desugar().Error("create caller", zap.String("error", "failed to initialize caller, unknown type"))
		return nil, errors.New("failed to initialize caller, unknown type")
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
