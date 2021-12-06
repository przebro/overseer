package types

//TaskType - Type of a task that can be scheduled and processed
type TaskType string

//Task types
const (
	TypeDummy       TaskType = "dummy"
	TypeOs          TaskType = "os"
	TypeAws         TaskType = "aws"
	TypeFtp         TaskType = "ftp"
	TypeFileWatcher TaskType = "filewatch"
	TypeDatabase    TaskType = "database"
)
