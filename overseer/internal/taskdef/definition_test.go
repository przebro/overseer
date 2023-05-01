package taskdef

import (
	"fmt"
	"testing"
)

// /Users/casual/dev/goscheduler/def/samples/minimal.yaml

func TestLoadDefinition(t *testing.T) {
	path := "/Users/casual/dev/goscheduler/def/samples/sample_07.json"
	def, err := ReadFromFile(path)
	fmt.Println(def, err)
}
