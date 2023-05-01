package models

import (
	"time"

	"github.com/przebro/overseer/common/types"
)

// Should use type from types
type TaskInTicket struct {
	Name      string
	Odate     string
	Fulfilled bool
}
type TaskExecution struct {
	ExecutionID string
	Worker      string
	Start       time.Time
	End         time.Time
	State       TaskState
}

type TaskCycle struct {
	IsCyclic    bool
	NextRun     types.HourMinTime
	RunFrom     string
	MaxRun      int
	RunInterval int
}
