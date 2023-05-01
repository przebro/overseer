package aws

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/ovsworker/jobs"
	"github.com/przebro/overseer/ovsworker/status"
	"github.com/rs/zerolog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	awstypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
)

type awsStepFuncCaller struct {
	job           jobs.Job
	client        *sfn.Client
	cancelFunc    context.CancelFunc
	payloadReader awsPayloadReader
	machineARN    string
	executionName string
}

func newStepFunctionCaller(conf aws.Config, job jobs.Job, machineARN, executionName string, payloadReader awsPayloadReader) awsServiceCaller {

	client := sfn.NewFromConfig(conf)

	c := &awsStepFuncCaller{client: client, job: job, machineARN: machineARN, executionName: executionName, payloadReader: payloadReader}

	return c

}
func (c *awsStepFuncCaller) Call(ctx context.Context, stat chan<- status.JobExecutionStatus) status.JobExecutionStatus {

	log := zerolog.Ctx(ctx)
	fpath := filepath.Join(c.job.SysoutDir, c.job.ExecutionID)

	file, err := os.Create(fpath)
	if err != nil {
		log.Error().Err(err).Str("path", fpath).Msg("newStepFunctionCaller")
		return status.StatusFailed(c.job.TaskID, c.job.ExecutionID, err.Error())
	}

	payload, err := c.payloadReader.Read()
	if err != nil {
		log.Error().Err(err).Msg("newStepFunctionCaller")
		return status.StatusFailed(c.job.TaskID, c.job.ExecutionID, err.Error())
	}

	strpayload := string(payload)

	out, err := c.client.StartExecution(ctx,
		&sfn.StartExecutionInput{
			StateMachineArn: &c.machineARN,
			Name:            &c.executionName,
			Input:           &strpayload,
		},
	)

	if err != nil {
		defer file.Close()
		log.Error().Err(err).Msg("newStepFunctionCaller")
		return status.StatusFailed(c.job.TaskID, c.job.ExecutionID, err.Error())
	}

	go c.waitForExecutionEnd(ctx, file, *out.ExecutionArn, stat)

	return status.StatusExecuting(c.job.TaskID, c.job.ExecutionID)
}

func (c *awsStepFuncCaller) waitForExecutionEnd(ctx context.Context, f *os.File, executionARN string, stat chan<- status.JobExecutionStatus) {

	log := zerolog.Ctx(ctx)

	defer f.Close()
	defer c.cancelFunc()

	ticker := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-ctx.Done():
			{
				return
			}
		case <-ticker.C:
			{
				out, err := c.client.DescribeExecution(ctx, &sfn.DescribeExecutionInput{
					ExecutionArn: &executionARN,
				})

				if err != nil {
					log.Error().Err(err).Msg("newStepFunctionCaller")
					stat <- status.StatusFailed(c.job.TaskID, c.job.ExecutionID, err.Error())
				}

				if out.Status == awstypes.ExecutionStatusRunning {
					stat <- status.StatusExecuting(c.job.TaskID, c.job.ExecutionID)

				} else {

					var statusCode types.StatusCode

					switch out.Status {
					case awstypes.ExecutionStatusAborted:
						{
							statusCode = types.StatusCodeAborted
						}
					case awstypes.ExecutionStatusFailed:
						{
							statusCode = types.StatusCodeError
						}
					case awstypes.ExecutionStatusTimedOut:
						{
							statusCode = types.StatusCodeTimeout
						}
					case awstypes.ExecutionStatusSucceeded:
						{
							statusCode = types.StatusCodeNormal
						}
					}

					if out.Output != nil {
						f.Write([]byte(*out.Output))
					}

					stat <- status.StatusEnded(c.job.TaskID, c.job.ExecutionID, 0, 0, int32(statusCode))
					ticker.Stop()
					return
				}
			}

		}

	}

}

func (c *awsStepFuncCaller) StartJob(ctx context.Context, stat chan status.JobExecutionStatus) status.JobExecutionStatus {

	nctx, cfunc := context.WithCancel(ctx)
	c.cancelFunc = cfunc
	return c.Call(nctx, stat)
}
func (c *awsStepFuncCaller) CancelJob() error {

	if c.cancelFunc == nil {
		return errors.New("failed to cancel job")
	}
	c.cancelFunc()

	return nil
}

func (c *awsStepFuncCaller) JobTaskID() string {
	return c.job.TaskID
}
func (c *awsStepFuncCaller) JobExecutionID() string {
	return c.job.ExecutionID
}
