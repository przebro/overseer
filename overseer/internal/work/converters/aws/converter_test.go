package aws

import (
	"encoding/json"
	"testing"

	"github.com/przebro/overseer/common/types"
)

var input string = `{"type":"lambda","connection":{"profile":"overseer","region":"eu-west-1"},"payload":{},"functionName":"test_function_01","alias":"$LATEST"}`

func Test_Converter(t *testing.T) {

	awsConverter := &awsConverter{}
	_, err := awsConverter.ConvertToMsg(json.RawMessage(input), types.EnvironmentVariableList{})
	if err != nil {
		t.Error("unexpected result:", err)
	}
}
