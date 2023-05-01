package taskdef

import (
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/unique"
	"github.com/przebro/overseer/common/validator"
)

// TaskBuilder - task builder.
type TaskBuilder interface {
	FromTemplate(def *TaskDefinition) TaskBuilder
	WithBase(group, name, description string) TaskBuilder
	WithSchedule(schedule SchedulingData) TaskBuilder
	WithInTicekts(in []InTicketData, expr string) TaskBuilder
	WithOutTickets(out []OutTicketData) TaskBuilder
	WithFlags(flags []FlagData) TaskBuilder
	WithConfirm() TaskBuilder
	WithCyclic(data CyclicTaskData) TaskBuilder
	WithVariables(vars types.EnvironmentVariableList) TaskBuilder
	Build() (*TaskDefinition, error)
}

// DummyTaskBuilder - dummy Task builder
type DummyTaskBuilder struct {
	def *TaskDefinition
}

func NewBuilder() TaskBuilder {

	return &DummyTaskBuilder{def: &TaskDefinition{}}
}

// FromTemplate - Helps create a new task based on an existing definition.
func (builder *DummyTaskBuilder) FromTemplate(templ *TaskDefinition) TaskBuilder {

	builder.def.Definition.Name, builder.def.Definition.Group, builder.def.Definition.Description = templ.GetInfo()
	builder.def.Definition.Type = templ.TypeName()
	builder.def.Definition.Rev = templ.Definition.Rev

	if templ.InTickets != nil {
		builder.def.InTickets = make([]InTicketData, len(templ.InTickets))
		copy(builder.def.InTickets, templ.InTickets)
	}

	builder.def.Confirm = templ.Confirm

	if templ.OutTickets != nil {
		builder.def.OutTickets = make([]OutTicketData, len(templ.OutTickets))
		copy(builder.def.OutTickets, templ.OutTickets)
	}

	builder.def.Flags = make([]FlagData, len(templ.Flags))
	copy(builder.def.Flags, templ.Flags)

	builder.def.Schedule = SchedulingData{}

	builder.def.Schedule.FromTime, builder.def.Schedule.ToTime = templ.TimeSpan()
	builder.def.Schedule.OrderType = templ.Schedule.OrderType
	builder.def.Schedule.Months = make([]time.Month, len(templ.Schedule.Months))
	copy(builder.def.Schedule.Months, templ.Schedule.Months)

	builder.def.Cyclic = templ.Cyclic

	builder.def.Schedule.Dayvalues = make([]int, len(templ.Schedule.Dayvalues))
	copy(builder.def.Schedule.Dayvalues, templ.Schedule.Dayvalues)

	builder.def.Schedule.Exactdates = make([]string, len(templ.Schedule.Exactdates))
	copy(builder.def.Schedule.Exactdates, templ.Schedule.Exactdates)

	return builder
}

// WithBase - Adds base information to the constructed task.
func (builder *DummyTaskBuilder) WithBase(group, name, description string) TaskBuilder {

	builder.def.Definition.Name = name
	builder.def.Definition.Group = group
	builder.def.Definition.Description = description
	builder.def.Definition.Rev = name + "@" + group + "@" + unique.NewID().Hex()

	return builder
}

// WithSchedule - Adds schedule information to the constructed task.
func (builder *DummyTaskBuilder) WithSchedule(schedule SchedulingData) TaskBuilder {

	builder.def.Schedule = schedule

	return builder
}

// WithInTicekts - Adds input tickets to the constructed task.
func (builder *DummyTaskBuilder) WithInTicekts(in []InTicketData, expr string) TaskBuilder {

	builder.def.InTickets = in
	builder.def.Expression = expr

	return builder
}

// WithOutTickets - Adds output tickets to the constructed task.
func (builder *DummyTaskBuilder) WithOutTickets(out []OutTicketData) TaskBuilder {

	builder.def.OutTickets = out

	return builder
}

// WithFlags - Adds flags to the constructed task.
func (builder *DummyTaskBuilder) WithFlags(flags []FlagData) TaskBuilder {

	builder.def.Flags = flags

	return builder
}

// WithConfirm - Adds confirm to the constructed task.
func (builder *DummyTaskBuilder) WithConfirm() TaskBuilder {

	builder.def.Confirm = true

	return builder
}

// WithVariables - Adds variables to the constructed task.
func (builder *DummyTaskBuilder) WithVariables(vars types.EnvironmentVariableList) TaskBuilder {

	builder.def.VariableList = vars
	return builder
}

// WithCyclic - Adds cyclic settings to the task
func (builder *DummyTaskBuilder) WithCyclic(data CyclicTaskData) TaskBuilder {

	builder.def.Cyclic = data
	return builder
}

// Build - Builds a new task definition.
func (builder *DummyTaskBuilder) Build() (*TaskDefinition, error) {

	builder.def.Definition.Type = types.TypeDummy

	//make a copy of a final product and clear builder instance of an object
	prod := builder.def
	if err := validator.Valid.Validate(*prod); err != nil {
		return nil, err
	}
	builder.def = &TaskDefinition{}
	return prod, nil
}
