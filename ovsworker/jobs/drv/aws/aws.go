package aws

import (
	"context"
	"errors"
	"overseer/common/types"
	"overseer/ovsworker/jobs"
	"overseer/ovsworker/msgheader"
	"overseer/ovsworker/status"
	"overseer/proto/actions"

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
	Call(ctx context.Context, stat chan<- status.JobExecutionStatus)
}

//AwsJobFactory - Creates a new aws factory
func AwsJobFactory(header msgheader.TaskHeader, sysoutDir string, data []byte) (jobs.JobExecutor, error) {

	act := actions.AwsTaskAction{}
	if err := proto.Unmarshal(data, &act); err != nil {
		return nil, err
	}

	return newAwsJob(header, sysoutDir, &act)
}

//newDummyJob - factory method
func newAwsJob(header msgheader.TaskHeader, sysoutDir string, action *actions.AwsTaskAction) (jobs.JobExecutor, error) {

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
	}
	for k, v := range header.Variables {
		job.Variables[k] = v
	}

	switch act := action.GetConnection().(type) {

	case *actions.AwsTaskAction_ConnectionData:
		{
			conf, err = createConfig(context.Background(), act.ConnectionData.ProfileName, act.ConnectionData.Region)
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
				return nil, errors.New("")
			}

			caller = newLambdaCaller(conf, job, lExecData.FunctionName, lExecData.Alias, payloadReader)
		}
	case actions.AwsTaskAction_stepfunc:
		{
			sExecData := action.GetStepFunction()
			if sExecData == nil {
				return nil, errors.New("")
			}
			caller = newStepFunctionCaller(conf, job, sExecData.StateMachineARN, sExecData.ExecutionName, payloadReader)
		}
	default:
		{
			caller = nil
		}
	}

	if caller == nil {
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
