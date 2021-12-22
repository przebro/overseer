package aws

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"
	"overseer/ovsworker/jobs"
	"overseer/ovsworker/status"
	"path/filepath"

	"overseer/common/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	awstypes "github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type awsLambdaCaller struct {
	job           jobs.Job
	client        *lambda.Client
	payloadReader awsPayloadReader
	cancelFunc    context.CancelFunc
	function      string
	alias         string
}

const logCallLambdaHeader = "call-lambda"

func newLambdaCaller(conf aws.Config, job jobs.Job, functionName, alias string, payloadReader awsPayloadReader) awsServiceCaller {

	client := lambda.NewFromConfig(conf)
	c := &awsLambdaCaller{client: client, job: job, payloadReader: payloadReader, function: functionName, alias: alias}

	return c

}

func (j *awsLambdaCaller) Call(ctx context.Context, stat chan<- status.JobExecutionStatus) status.JobExecutionStatus {

	fpath := filepath.Join(j.job.SysoutDir, j.job.ExecutionID)

	file, err := os.Create(fpath)
	if err != nil {
		j.job.Log.Desugar().Error(logCallLambdaHeader, j.fields(zap.String("error", err.Error()), zap.String("path", fpath))...)
		return status.StatusFailed(j.job.TaskID, j.job.ExecutionID, err.Error())
	}

	payload, err := j.payloadReader.Read()

	if err != nil {
		j.job.Log.Desugar().Error(logCallLambdaHeader, j.fields(zap.String("error", err.Error()))...)
		return status.StatusFailed(j.job.TaskID, j.job.ExecutionID, err.Error())
	}

	customData := creteCustomContextData(j.job.Variables)

	go func() {

		defer file.Close()

		result, err := j.client.Invoke(ctx,
			&lambda.InvokeInput{
				FunctionName:   &j.function,
				InvocationType: awstypes.InvocationTypeRequestResponse,
				LogType:        awstypes.LogTypeTail,
				ClientContext:  &customData,
				Payload:        payload,
				Qualifier:      &j.alias,
			})

		if err != nil {
			j.job.Log.Desugar().Error(logCallLambdaHeader, j.fields(zap.String("error", err.Error()))...)
			stat <- status.StatusFailed(j.job.TaskID, j.job.ExecutionID, err.Error())
			file.Write([]byte(err.Error()))

			return
		}

		file.Write(result.Payload)

		errDescr := ""
		statCode := types.StatusCodeNormal
		if result.FunctionError != nil {
			errDescr = *result.FunctionError
			statCode = types.StatusCodeError
		}

		stat <- status.JobExecutionStatus{
			TaskID:      j.job.TaskID,
			ExecutionID: j.job.ExecutionID,
			State:       types.WorkerTaskStatusEnded,
			ReturnCode:  int(result.StatusCode),
			StatusCode:  int32(statCode),
			Reason:      errDescr,
		}
	}()

	return status.StatusExecuting(j.job.TaskID, j.job.ExecutionID)

}

func (j *awsLambdaCaller) StartJob(ctx context.Context, stat chan status.JobExecutionStatus) status.JobExecutionStatus {

	nctx, cfunc := context.WithCancel(ctx)
	j.cancelFunc = cfunc
	return j.Call(nctx, stat)
}
func (j *awsLambdaCaller) CancelJob() error {

	if j.cancelFunc == nil {
		return errors.New("failed to cancel job")
	}
	j.cancelFunc()

	return nil
}

//JobTaskID - returns ID of a task associated with this job.
func (j *awsLambdaCaller) JobTaskID() string {

	return j.job.TaskID
}

//JobExecutionID - returns ExecutionID of a task
func (j *awsLambdaCaller) JobExecutionID() string {

	return j.job.ExecutionID
}

func creteCustomContextData(data map[string]string) string {

	customData := struct{ Custom map[string]string }{Custom: map[string]string{}}
	for key, val := range data {
		customData.Custom[key] = val
	}

	b, _ := json.Marshal(&customData)

	return base64.RawStdEncoding.EncodeToString(b)
}

func (j *awsLambdaCaller) fields(fields ...zapcore.Field) []zap.Field {
	f := []zapcore.Field{
		zap.String("taskID", j.job.TaskID), zap.String("executionID", j.job.ExecutionID),
		zap.String("functionName", j.function), zap.String("alias", j.alias),
	}

	fields = append(fields, f...)

	return fields
}
