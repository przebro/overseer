package msgheader

import (
	common "overseer/common/types"
)

//TaskHeader - Common data for all tasks
type TaskHeader struct {
	TaskID      string
	ExecutionID string
	Type        common.TaskType
	Variables   map[string]string
}
