package taskdef

import (
	"overseer/common/types"
	"overseer/common/validator"
	"time"
)

//TaskBuilder - task builder.
type TaskBuilder interface {
	FromTemplate(def TaskDefinition) TaskBuilder
	WithBase(group, name, description string) TaskBuilder
	WithSchedule(schedule SchedulingData) TaskBuilder
	WithInTicekts(in []InTicketData, relation InTicketRelation) TaskBuilder
	WithOutTickets(out []OutTicketData) TaskBuilder
	WithFlags(flags []FlagData) TaskBuilder
	WithConfirm() TaskBuilder
	WithVariables(vars []VariableData) TaskBuilder
	WithRetention(days int) TaskBuilder
	Build() (TaskDefinition, error)
}

//DummyTaskBuilder - dummy Task builder
type DummyTaskBuilder struct {
	def baseTaskDefinition
}

//FromTemplate - Helps create a new task based on an existing definition.
func (builder *DummyTaskBuilder) FromTemplate(templ TaskDefinition) TaskBuilder {

	builder.def.Name, builder.def.Group, builder.def.Description = templ.GetInfo()
	builder.def.TaskType = templ.TypeName()

	if templ.TicketsIn() != nil {
		builder.def.InTickets = make([]InTicketData, len(templ.TicketsIn()))
		copy(builder.def.InTickets, templ.TicketsIn())
	}

	builder.def.InRelation = templ.Relation()
	builder.def.DataRetention = templ.Retention()
	builder.def.ConfirmFlag = templ.Confirm()

	if templ.TicketsOut() != nil {
		builder.def.OutTickets = make([]OutTicketData, len(templ.TicketsOut()))
		copy(builder.def.OutTickets, templ.TicketsOut())
	}

	builder.def.FlagsTab = make([]FlagData, len(templ.Flags()))
	copy(builder.def.FlagsTab, templ.Flags())

	builder.def.Schedule = SchedulingData{}

	builder.def.Schedule.FromTime, builder.def.Schedule.ToTime = templ.TimeSpan()
	builder.def.Schedule.AllowPastSub = templ.AllowPast()
	builder.def.Schedule.OrderType = templ.OrderType()
	builder.def.Schedule.Months = make([]time.Month, len(templ.Months()))
	copy(builder.def.Schedule.Months, templ.Months())

	builder.def.Schedule.Dayvalues = make([]int, len(templ.Days()))
	copy(builder.def.Schedule.Dayvalues, templ.Days())

	builder.def.Schedule.Exactdates = make([]string, len(templ.ExactDate()))
	copy(builder.def.Schedule.Exactdates, templ.ExactDate())

	return builder
}

//WithBase - Adds base information to the constructed task.
func (builder *DummyTaskBuilder) WithBase(group, name, description string) TaskBuilder {

	builder.def.Name = name
	builder.def.Group = group
	builder.def.Description = description

	return builder
}

//WithSchedule - Adds schedule information to the constructed task.
func (builder *DummyTaskBuilder) WithSchedule(schedule SchedulingData) TaskBuilder {

	builder.def.Schedule = schedule

	return builder
}

//WithInTicekts - Adds input tickets to the constructed task.
func (builder *DummyTaskBuilder) WithInTicekts(in []InTicketData, relation InTicketRelation) TaskBuilder {

	builder.def.InRelation = relation
	builder.def.InTickets = in

	return builder
}

//WithOutTickets - Adds output tickets to the constructed task.
func (builder *DummyTaskBuilder) WithOutTickets(out []OutTicketData) TaskBuilder {

	builder.def.OutTickets = out

	return builder
}

//WithFlags - Adds flags to the constructed task.
func (builder *DummyTaskBuilder) WithFlags(flags []FlagData) TaskBuilder {

	builder.def.FlagsTab = flags

	return builder
}

//WithConfirm - Adds confirm to the constructed task.
func (builder *DummyTaskBuilder) WithConfirm() TaskBuilder {

	builder.def.ConfirmFlag = true

	return builder
}

//WithVariables - Adds variables to the constructed task.
func (builder *DummyTaskBuilder) WithVariables(vars []VariableData) TaskBuilder {

	builder.def.TaskVariables = vars
	return builder
}

//WithRetention - Adds retention to the constructed task.
func (builder *DummyTaskBuilder) WithRetention(days int) TaskBuilder {

	builder.def.DataRetention = days

	return builder
}

//Build - Builds a new task definition.
func (builder *DummyTaskBuilder) Build() (TaskDefinition, error) {

	builder.def.TaskType = types.TypeDummy

	//make a copy of a final product and clear builder instance of an object
	prod := builder.def
	if err := validator.Valid.Validate(prod); err != nil {
		return nil, err
	}
	builder.def = baseTaskDefinition{}
	return &prod, nil
}
