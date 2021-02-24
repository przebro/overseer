package fragments

import (
	"golang.org/x/net/context"
)

type workFragment struct {
	taskID      string
	executionID string
	sysoutDir   string
	start       chan FragmentStatus
	Variables   []string
}

//WorkFragment - Represents a piece of work that will be executed.
type WorkFragment interface {
	StartFragment(ctx context.Context, stat chan FragmentStatus)
	CancelFragment() error
	//TaskID - Returns ID of a task associated with this fragment.
	TaskID() string
	ExecutionID() string
}
