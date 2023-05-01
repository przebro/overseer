package taskdef

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/przebro/expr"
	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/common/validator"
	"gopkg.in/yaml.v3"
)

// InTicketRelation - Restricts possible values of tickets relation.
// Possible values are: AND, OR, EXPR
type InTicketRelation string

// FlagType - Restricts possible values of flags type
type FlagType string

// SchedulingOption - how job will be ordered
type SchedulingOption string

// OutAction - Restricts possible out ticket action
// Possible values are ADD,REM
type OutAction string

// CycleFromOption -
type CycleFromOption string

var (
	//CycleFromStart - computes next time based on start time
	CycleFromStart CycleFromOption = "start"
	//CycleFromEnd - computes next time based on end time
	CycleFromEnd CycleFromOption = "end"
	//CycleFromSchedule - computes next time based on new day proc
	CycleFromSchedule CycleFromOption = "schedule"
)

// InTicketData - Holds information about the required tickets to start a task
type InTicketData struct {
	Name  string          `json:"name" yaml:"name" validate:"required,max=64,resname"`
	Odate date.OdateValue `json:"odate" yaml:"odate" validate:"odateval"`
}

// OutTicketData - Holds an action for a given ticket that is performed after the successful execution of a task.
type OutTicketData struct {
	Name   string          `json:"name" yaml:"name" validate:"required,max=64,resname"`
	Odate  date.OdateValue `json:"odate"  yaml:"odate" validate:"odateval"`
	Action OutAction       `json:"action" yaml:"action" validate:"required,oneof=ADD REM"`
}

// FlagData - Holds information about required flag resources.
// If type is set to exclusive then the state of a flag in resources manager can be only none
// If type is set to shared then the state of a flag in resources manager can be none or shared but not exclusive
type FlagData struct {
	Name string   `json:"name" yaml:"name" validate:"required,max=32,resname"`
	Type FlagType `json:"type" yaml:"type" validate:"oneof=SHR EXL"`
}

// SchedulingData - Holds informations how task should be scheduled.
type SchedulingData struct {
	OrderType  SchedulingOption  `json:"type" yaml:"type" validate:"required,oneof=manual daily weekday dayofmonth exact fromend"`
	FromTime   types.HourMinTime `json:"from" yaml:"from" validate:"omitempty,hmtime"`
	ToTime     types.HourMinTime `json:"to" yaml:"to" validate:"omitempty,hmtime"`
	Months     []time.Month      `json:"months" yaml:"months" validate:"unique"`
	Exactdates []string          `json:"exact" yaml:"exact" validate:"unique"`
	Dayvalues  []int             `json:"days"  yaml:"days"  validate:"unique"`
}

// CyclicTaskData -
type CyclicTaskData struct {
	RunFrom      CycleFromOption `json:"from" yaml:"from" validate:"omitempty,oneof=start end schedule"`
	MaxRuns      int             `json:"runs" yaml:"runs"  validate:"min=0,max=999"`
	TimeInterval int             `json:"every" yaml:"every"  validate:"min=0,max=1440"`
	IsCycle      bool            `json:"cycle" yaml:"cycle"`
}

type TaskData struct {
	Type        types.TaskType `yaml:"type" json:"type"`
	Rev         string         `yaml:"rev" json:"rev"`
	Name        string         `json:"name" yaml:"name" validate:"required,max=64,resname"`
	Group       string         `json:"group" yaml:"group" validate:"required,max=64,resname"`
	Description string         `json:"description" yaml:"description" validate:"lte=200"`
}

// Ordering options
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

// Possible flag types
const (
	FlagShared    FlagType = "SHR"
	FlagExclusive FlagType = "EXL"
)

// Possible out actions
const (
	OutActionAdd    OutAction = "ADD"
	OutActionRemove OutAction = "REM"
)

// Relation between input tickets
// Expect all: COND-1 AND COND-2 AND ...
// Expect one of them COND-1 OR COND-2 ...
// Complex evalution
const (
	InTicketAND  InTicketRelation = "AND"
	InTicketOR   InTicketRelation = "OR"
	InTicketExpr InTicketRelation = "EXPR"
)

type TaskDefinition struct {
	Definition   TaskData                      `yaml:"definition" json:"definition" validate:"required"`
	Confirm      bool                          `json:"confirm" yaml:"confirm"`
	Retain       bool                          `json:"retain" yaml:"retain" validate:"omitempty"`
	Schedule     SchedulingData                `json:"schedule" yaml:"schedule" validate:"omitempty"`
	Cyclic       CyclicTaskData                `json:"cyclic" yaml:"cyclic" validate:"omitempty"`
	Expression   string                        `json:"expr" yaml:"expr" validate:"omitempty"`
	InTickets    []InTicketData                `json:"in_tickets" yaml:"in_tickets" validate:"omitempty,dive"`
	Flags        []FlagData                    `json:"flags" yaml:"flags" validate:"omitempty,dive"`
	OutTickets   []OutTicketData               `json:"out_tickets" yaml:"out_tickets" validate:"omitempty,dive"`
	VariableList types.EnvironmentVariableList `json:"variables" yaml:"variables"  validate:"omitempty,dive"`
	Data         []byte                        `json:"spec" yaml:"spec"  validate:"omitempty"`
	Worker       string                        `json:"worker" yaml:"worker" validate:"omitempty"`
}

func (task *TaskDefinition) TimeSpan() (types.HourMinTime, types.HourMinTime) {
	return task.Schedule.FromTime, task.Schedule.ToTime
}

func (task *TaskDefinition) GetInfo() (string, string, string) {
	return task.Definition.Name, task.Definition.Group, task.Definition.Description
}
func (task *TaskDefinition) Variables() types.EnvironmentVariableList {
	return task.VariableList
}

func (task *TaskDefinition) Action() []byte {
	return task.Data
}

func (task *TaskDefinition) Payload() interface{} {
	return task.Data
}

func (task *TaskDefinition) TypeName() types.TaskType {
	return task.Definition.Type
}

// SerializeDefinition - Writes task definition to string.
func SerializeDefinition(definition *TaskDefinition) (string, error) {

	var result string

	data, err := json.Marshal(definition)
	if err != nil {
		return "", err
	}

	result = string(data)

	return result, nil
}

// ReadFromPoolDirectory - Load task from file, Wrapper function for load from string
func ReadFromPoolDirectory(path string) (*TaskDefinition, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	def, err := FromStream(data)
	if err != nil {
		return nil, err
	}

	return def, nil
}

// ReadFromFile - Load task from file, Wrapper function for load from string
func ReadFromFile(path string) (*TaskDefinition, error) {

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pth := strings.TrimSuffix(path, filepath.Ext(path))
	taskname := filepath.Base(pth)
	group := filepath.Base(filepath.Dir(pth))

	def, err := FromStream(data)
	if err != nil {
		return nil, err
	}

	n, g := def.Definition.Name, def.Definition.Group

	if n != taskname && group != g {
		return nil, errors.New("task name and group name does not match filepath")
	}
	if def.Expression != "" {
		if err := expr.Test(def.Expression); err != nil {
			return nil, fmt.Errorf("failed to parse expression:%v", err)
		}
	}

	return def, nil
}

func FromStream(data []byte) (*TaskDefinition, error) {

	if isJSON(data) {
		return fromJSON(data)
	}
	return fromYAML(data)
}

func fromJSON(data []byte) (*TaskDefinition, error) {

	def := &TaskDefinition{}

	if err := json.Unmarshal(data, def); err != nil {
		return nil, err
	}

	if err := validator.Valid.Validate(*def); err != nil {
		return nil, err
	}

	return def, nil
}

func fromYAML(data []byte) (*TaskDefinition, error) {

	def := &TaskDefinition{}

	if err := yaml.Unmarshal(data, def); err != nil {
		return nil, err
	}

	if err := validator.Valid.Validate(*def); err != nil {
		return nil, err
	}

	return def, nil
}

func isJSON(data []byte) bool {
	val := json.RawMessage{}
	return json.Unmarshal(data, &val) == nil
}
