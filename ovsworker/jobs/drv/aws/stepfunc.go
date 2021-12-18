package aws

import (
	"context"
	"errors"
	"os"
	"overseer/common/types"
	"overseer/ovsworker/jobs"
	"overseer/ovsworker/status"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sfn"
	awstypes "github.com/aws/aws-sdk-go-v2/service/sfn/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type awsStepFuncCaller struct {
	job           jobs.Job
	client        *sfn.Client
	cancelFunc    context.CancelFunc
	payloadReader awsPayloadReader
	machineARN    string
	executionName string
}

const logCallStepfuncHeader = "call-stepfunc"

func newStepFunctionCaller(conf aws.Config, job jobs.Job, machineARN, executionName string, payloadReader awsPayloadReader) awsServiceCaller {

	client := sfn.NewFromConfig(conf)

	c := &awsStepFuncCaller{client: client, job: job, machineARN: machineARN, executionName: executionName, payloadReader: payloadReader}

	return c

}
func (c *awsStepFuncCaller) Call(ctx context.Context, stat chan<- status.JobExecutionStatus) {

	fpath := filepath.Join(c.job.SysoutDir, c.job.ExecutionID)

	file, err := os.Create(fpath)
	if err != nil {
		c.job.Log.Desugar().Error(logCallStepfuncHeader, c.fields(zap.String("error", err.Error()), zap.String("path", fpath))...)
		stat <- status.StatusFailed(c.job.TaskID, c.job.ExecutionID, err.Error())
		return
	}

	payload, err := c.payloadReader.Read()
	if err != nil {
		c.job.Log.Desugar().Error(logCallStepfuncHeader, c.fields(zap.String("error", err.Error()))...)
		stat <- status.StatusFailed(c.job.TaskID, c.job.ExecutionID, err.Error())
		return
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
		c.job.Log.Desugar().Error(logCallStepfuncHeader, c.fields(zap.String("error", err.Error()))...)
		stat <- status.StatusFailed(c.job.TaskID, c.job.ExecutionID, err.Error())
		return
	}

	stat <- status.StatusExecuting(c.job.TaskID, c.job.ExecutionID)
	go c.waitForExecutionEnd(ctx, file, *out.ExecutionArn, stat)

}

func (c *awsStepFuncCaller) waitForExecutionEnd(ctx context.Context, f *os.File, executionARN string, stat chan<- status.JobExecutionStatus) {

	defer f.Close()
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
					c.job.Log.Desugar().Error(logCallStepfuncHeader, c.fields(zap.String("error", err.Error()))...)
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

func (c *awsStepFuncCaller) StartJob(ctx context.Context, stat chan status.JobExecutionStatus) {

	nctx, cfunc := context.WithCancel(ctx)
	c.cancelFunc = cfunc
	c.Call(nctx, stat)

}
func (c *awsStepFuncCaller) CancelJob() error {

	if c.cancelFunc == nil {
		return errors.New("cancel function is nil")
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

func (c *awsStepFuncCaller) fields(fields ...zapcore.Field) []zap.Field {
	f := []zapcore.Field{
		zap.String("taskID", c.job.TaskID), zap.String("executionID", c.job.ExecutionID),
		zap.String("stateMachineARN", c.machineARN), zap.String("executionName", c.executionName),
	}

	fields = append(fields, f...)

	return fields
}
