package taskdef

//OsActionType - type of a possible os action
type OsActionType string

//Possible values for OsActionType
const (
	OsActionTypeCommand OsActionType = "command"
	OsActionTypeScript  OsActionType = "script"
)

//OsTaskData - specific data for OS task
type OsTaskData struct {
	ActionType  OsActionType `json:"type"`
	CommandLine string       `json:"command"`
	RunAs       string       `json:"runas"`
}

//OsTaskDefinition - definition of a OS task
type OsTaskDefinition struct {
	baseTaskDefinition
	Spec OsTaskData
}

//Action - Returns action defined in a task.
func (os *OsTaskDefinition) Action() interface{} {
	return os.Spec
}
