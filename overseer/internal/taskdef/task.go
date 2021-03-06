package taskdef

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"overseer/common/types"
	"overseer/common/types/date"
	"overseer/common/validator"
	"path/filepath"
	"strings"
	"time"
)

//InTicketRelation - Restricts possible values of tickets relation.
//Possible values are: AND, OR
type InTicketRelation string

//FlagType - Restricts possible values of flags type
type FlagType string

//SchedulingOption - how job will be ordered
type SchedulingOption string

//OutAction - Restricts possible out ticket action
//Possible values are ADD,REM
type OutAction string

//InTicketData - Holds information about the required tickets to start a task
type InTicketData struct {
	Name  string          `json:"name" validate:"required,max=32,resname"`
	Odate date.OdateValue `json:"odate" validate:"odateval"`
}

//OutTicketData - Holds an action for a given ticket that is performed after the successful execution of a task.
type OutTicketData struct {
	Name   string          `json:"name" validate:"required,max=32,resname"`
	Odate  date.OdateValue `json:"odate" validate:"odateval"`
	Action OutAction       `json:"action" validate:"required,oneof=ADD REM"`
}

//FlagData - Holds information about required flag resources.
//If type is set to exclusive then the state of a flag in resources manager can be only none
//If type is set to shared then the state of a flag in resources manager can be none or shared but not exclusive
type FlagData struct {
	Name string   `json:"name" validate:"required,max=32,resname"`
	Type FlagType `json:"type" validate:"oneof=SHR EXL"`
}

//VariableData - Holds variables that will be passed to the task.
type VariableData struct {
	Name  string `json:"name" validate:"required,max=32,varname"`
	Value string `json:"value"`
}

//Expand - Expands name
func (data VariableData) Expand() string {
	return strings.Replace(data.Name, "%%", "OVS_", 1)
}

//SchedulingData - Holds informations how task should be scheduled.
type SchedulingData struct {
	OrderType    SchedulingOption  `json:"type" validate:"required,oneof=manual daily weekday dayofmonth exact fromend"`
	FromTime     types.HourMinTime `json:"from" validate:"omitempty,hmtime"`
	ToTime       types.HourMinTime `json:"to" validate:"omitempty,hmtime"`
	AllowPastSub bool              `json:"pastsub"`
	Months       []time.Month      `json:"months" validate:"unique"`
	Exactdates   []string          `json:"exact" validate:"unique"`
	Dayvalues    []int             `json:"days"  validate:"unique"`
}

type baseTaskDefinition struct {
	TaskType      types.TaskType   `json:"type" validate:"oneof=dummy os"`
	Name          string           `json:"name" validate:"required,max=32,resname"`
	Group         string           `json:"group" validate:"required,max=20,resname"`
	Description   string           `json:"description" validate:"lte=200"`
	ConfirmFlag   bool             `json:"confirm"`
	DataRetention int              `json:"retention" validate:"min=0,max=14"`
	Schedule      SchedulingData   `json:"schedule" validate:"omitempty"`
	InTickets     []InTicketData   `json:"inticket" validate:"omitempty,dive"`
	InRelation    InTicketRelation `json:"relation" validate:"required_with=InTickets,omitempty,oneof=AND OR"`
	FlagsTab      []FlagData       `json:"flags"  validate:"omitempty,dive"`
	OutTickets    []OutTicketData  `json:"outticket"  validate:"omitempty,dive"`
	TaskVariables []VariableData   `json:"variables"  validate:"omitempty,dive"`
	Data          json.RawMessage  `json:"spec"`
}

//Ordering options
const (
	//manual means that there are no specific rules for the task, however this definition will be not scheduled during the end of day procedure.
	OrderingManual SchedulingOption = "manual"
	//daily means that the task will be scheduled on each day
	OrderingDaily SchedulingOption = "daily"
	//weekday means that the task will be scheduled on a specified day of the week
	OrderingWeek SchedulingOption = "weekday"
	//dayofmonth means that the task will be scheduled on a specified exact day of a month
	OrderingDayOfMonth SchedulingOption = "dayofmonth"
	//exact means that the task will be scheduled on the exact date
	OrderingExact SchedulingOption = "exact"
	/*fromend means that the task will be scheduled on a specified day from the end of the month
	where 1 means the last day, 2 means the day before the last day, and so on, up to 14
	for instance: fromend 1 means 31 of July,30 of June and 28 of February or 29 of February if it is a leap year
				  fromend 2 means 30 of July 29 of June and 27 of February or 28 of February if it is a leap year
	*/
	OrderingFromEnd SchedulingOption = "fromend"
)

//Possible flag types
const (
	FlagShared    FlagType = "SHR"
	FlagExclusive FlagType = "EXL"
)

//Possible out actions
const (
	OutActionAdd    OutAction = "ADD"
	OutActionRemove OutAction = "REM"
)

//Relation between input tickets
//Expect all: COND-1 AND COND-2 AND ...
//Expect one of them COND-1 OR COND-2 ...
const (
	InTicketAND InTicketRelation = "AND"
	InTicketOR  InTicketRelation = "OR"
)

//TaskScheduling  - Provides information about the schedule of a task.
type TaskScheduling interface {
	OrderType() SchedulingOption
	TimeSpan() (types.HourMinTime, types.HourMinTime)
	Months() []time.Month
	AllowPast() bool
	ExactDate() []string
	Days() []int
	Calendar() SchedulingData
}

func (task *baseTaskDefinition) OrderType() SchedulingOption {
	return task.Schedule.OrderType
}
func (task *baseTaskDefinition) TimeSpan() (types.HourMinTime, types.HourMinTime) {
	return task.Schedule.FromTime, task.Schedule.ToTime
}
func (task *baseTaskDefinition) Months() []time.Month {
	return task.Schedule.Months
}
func (task *baseTaskDefinition) AllowPast() bool {
	return task.Schedule.AllowPastSub
}

func (task *baseTaskDefinition) ExactDate() []string {
	return task.Schedule.Exactdates
}

func (task *baseTaskDefinition) Calendar() SchedulingData {
	return task.Schedule
}
func (task *baseTaskDefinition) Days() []int {
	return task.Schedule.Dayvalues
}

//BaseInfo - returns base informations
type BaseInfo interface {
	//GetInfo - gets base informations about task: Name,group and description
	GetInfo() (string, string, string)
}

//GetInfo - gets base informations about task: Name,group and description
func (task *baseTaskDefinition) GetInfo() (string, string, string) {
	return task.Name, task.Group, task.Description
}

//TaskInTicket - Provides information about required tickets
type TaskInTicket interface {
	TicketsIn() []InTicketData
	Relation() InTicketRelation
}

func (task *baseTaskDefinition) TicketsIn() []InTicketData {

	return task.InTickets
}
func (task *baseTaskDefinition) Relation() InTicketRelation {

	return task.InRelation
}

//TaskOutTicket - Provides information about issue tickets after task end.
type TaskOutTicket interface {
	TicketsOut() []OutTicketData
}

func (task *baseTaskDefinition) TicketsOut() []OutTicketData {

	return task.OutTickets
}

//TaskFlag - Provides information about required flags.
type TaskFlag interface {
	Flags() []FlagData
}

func (task *baseTaskDefinition) Flags() []FlagData {

	return task.FlagsTab
}

//TaskDefinition - Definition of an active task.
type TaskDefinition interface {
	BaseInfo
	TaskScheduling
	TaskInTicket
	TaskOutTicket
	TaskFlag
	TypeName() types.TaskType
	Confirm() bool
	Retention() int
	Variables() []VariableData
	Action() interface{}
}

//TypeName - returns task's type
func (task *baseTaskDefinition) TypeName() types.TaskType {
	return task.TaskType
}

//Confirm - Is manual confirmation by operator required
func (task *baseTaskDefinition) Confirm() bool {
	return task.ConfirmFlag
}

//Retention - How many days task should be kept in active task pool after successful ending
func (task *baseTaskDefinition) Retention() int {
	return task.DataRetention
}

//For dummy task this method returns empty string, for other specific tasks
//method returns information required to execute action
func (task *baseTaskDefinition) Action() interface{} {
	return ""
}

//Variables - Gets variables from tasks definition
func (task *baseTaskDefinition) Variables() []VariableData {
	return task.TaskVariables
}

//FromDefinitionFile - Load task from file, Wrapper function for load from string
func FromDefinitionFile(path string) (TaskDefinition, error) {

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pth := strings.TrimSuffix(path, filepath.Ext(path))
	taskname := filepath.Base(pth)
	group := filepath.Base(filepath.Dir(pth))

	def, err := FromString(string(data))
	if err != nil {
		return nil, err
	}

	n, g, _ := def.GetInfo()

	if n != taskname && group != g {
		return nil, errors.New("task name and group name does not match filepath")
	}

	return def, nil
}

//FromString - loads a task definition from input string. Performs validations for enum types.
func FromString(data string) (TaskDefinition, error) {

	var result TaskDefinition
	def := baseTaskDefinition{}
	err := json.Unmarshal([]byte(data), &def)
	if err != nil {
		return nil, err
	}

	if def.TaskType == types.TypeDummy {
		result = &def
	}
	if err = validator.Valid.Validate(def); err != nil {
		return nil, err
	}

	if def.TaskType == types.TypeOs {
		data := OsTaskData{}
		json.Unmarshal([]byte(def.Data), &data)
		result = &OsTaskDefinition{baseTaskDefinition: def, Spec: data}
	}

	return result, nil
}

//SerializeDefinition - Writes task definition to string.
func SerializeDefinition(definition TaskDefinition) (string, error) {

	var result string

	data, err := json.Marshal(definition)
	if err != nil {
		return "", err
	}

	result = string(data)

	return result, nil
}
