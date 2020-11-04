package types

//TaskType - Type of a task that can be scheduled and processed
type TaskType string

//Task types
const (
	TypeDummy    TaskType = "dummy"
	TypeOs       TaskType = "os"
	TypeFtp      TaskType = "ftp"
	TypeDatabase TaskType = "database"
)
