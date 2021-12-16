package taskdef

import (
	"encoding/json"
	"fmt"
	"testing"
)

var lambda string = `{"type":"lambda","connection" : {"profile" : "userprofile","region" : "us-west-1"},"payload" : {},"functionName" : "test_function","alias" : "$LATEST"}`

func TestAwsTask_Load(t *testing.T) {

	taskData := AwsTaskData{}
	if err := json.Unmarshal([]byte(lambda), &taskData); err != nil {
		t.Error(err)
	}

	fmt.Println(taskData)

	if taskData.Type == AWSActionTypeLambda {
		lambdaData := AwsLambdaTaskData{}
		if err := json.Unmarshal([]byte(lambda), &lambdaData); err != nil {
			t.Error(err)
		}

	}
}
