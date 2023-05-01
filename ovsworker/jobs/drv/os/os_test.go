package fragments

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/ovsworker/msgheader"
	"github.com/przebro/overseer/ovsworker/status"
	"github.com/przebro/overseer/proto/actions"
	"github.com/rs/zerolog"
)

func TestRunSingleJob(t *testing.T) {

	ctx := context.Background()
	ctx = zerolog.New(os.Stdout).With().Timestamp().Logger().WithContext(ctx)

	ex, err := newOsJob(ctx, msgheader.TaskHeader{
		TaskID:      "2235",
		ExecutionID: "1234567",
		Type:        "os",
		Variables:   map[string]string{},
	}, "./", &actions.OsTaskAction{
		CommandLine: "ls -l /Library",
		Steps: []*actions.OsStepDefinition{
			{
				StepName: "first_step",
				Command:  "pwd",
			},
		},
	})

	if err != nil {
		t.Error("failed to create job", err)
	}

	stat := make(chan status.JobExecutionStatus)

	ex.StartJob(ctx, stat)

	for s := range stat {
		fmt.Println(s.State)
		if s.State == types.WorkerTaskStatusFailed || s.State == types.WorkerTaskStatusEnded {
			break
		}
	}
}
