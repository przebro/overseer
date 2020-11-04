package msgheader

import (
	common "goscheduler/common/types"
)

//TaskHeader - Common data for all tasks
type TaskHeader struct {
	TaskID    string
	Type      common.TaskType
	Variables map[string]string
}
