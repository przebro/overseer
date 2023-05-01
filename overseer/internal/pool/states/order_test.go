package states

import (
	"testing"
	"time"

	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/rs/zerolog/log"
)

func TestStateOrderState(t *testing.T) {

	var result bool
	var definition *taskdef.TaskDefinition
	var err error

	builder := taskdef.NewBuilder()
	definition, err = builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).Build()

	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskOrderContext{
		Log:              log.Logger,
		Odate:            date.CurrentOdate(),
		IgnoreCalendar:   false,
		IgnoreSubmission: false,
		Def:              definition,
		CurrentOdate:     date.CurrentOdate(),
	}

	state := OstateCheckOtype{}
	stchkcal := &OstateCheckCalendar{}

	state.ProcessState(&ctx)
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != stchkcal {
		t.Error("expected result: ", true, " actual:", result, " ", stchkcal)
	}

	_, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	result = state.ProcessState(&ctx)

	if result != true && ctx.State != stchkcal {
		t.Error("expected result: ", true, " actual:", result, " ", stchkcal)
	}
}

func TestStateCheckCalendar(t *testing.T) {

	var result bool

	builder := taskdef.NewBuilder()
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskOrderContext{
		Log:              log.Logger,
		Odate:            date.CurrentOdate(),
		IgnoreCalendar:   false,
		IgnoreSubmission: false,
		Def:              definition,
		CurrentOdate:     date.CurrentOdate(),
	}

	state := &OstateCheckCalendar{}
	submState := &OstateOrdered{}
	cancelState := &OstateNotSubmitted{}

	result = state.ProcessState(&ctx)

	if result != true && ctx.State != cancelState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingDaily,
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
		}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.Def = definition
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != submState {
		t.Error("expected result: ", true, " actual:", result, " ", submState)
	}

	wday := int(time.Now().Weekday())
	nday := wday + 1

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingWeek,
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Dayvalues: []int{wday},
		}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.Def = definition
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != submState {
		t.Error("expected result: ", true, " actual:", result, " ", submState)
	}

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingWeek,
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Dayvalues: []int{nday},
		}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.Def = definition
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != cancelState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	dofmonth := time.Now().Day()
	ndofmonth := time.Now().AddDate(0, 0, 1).Day()

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingDayOfMonth,
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Dayvalues: []int{dofmonth},
		}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.Def = definition
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != submState {
		t.Error("expected result: ", true, " actual:", result, " ", submState)
	}

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingDayOfMonth,
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Dayvalues: []int{ndofmonth},
		}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.Def = definition
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != cancelState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	/**/
	edate := date.CurrentOdate().FormatDate()
	odate2 := date.AddDays(date.CurrentOdate(), 2)
	nedate := odate2.FormatDate()

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType:  taskdef.OrderingExact,
			Months:     []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Exactdates: []string{edate},
		}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.Def = definition
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != submState {
		t.Error("expected result: ", true, " actual:", result, " ", submState)
	}

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType:  taskdef.OrderingExact,
			Months:     []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Exactdates: []string{nedate},
		}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.Def = definition
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != cancelState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType:  taskdef.OrderingManual,
			Months:     []time.Month{},
			Exactdates: []string{nedate},
		}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.Def = definition
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != submState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	ctx.IgnoreCalendar = true
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != submState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

}

func TestStateCheckSubmmision(t *testing.T) {

	var result bool

	builder := taskdef.NewBuilder()
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Fatal("Unable to construct task")
	}

	subodat := date.AddDays(date.CurrentOdate(), -1)

	if err != nil {
		t.Error(err)
	}

	ctx := TaskOrderContext{
		Log:              log.Logger,
		Odate:            subodat,
		IgnoreCalendar:   false,
		IgnoreSubmission: false,
		Def:              definition,
		CurrentOdate:     date.CurrentOdate(),
	}

	state := &OstateCheckCalendar{}
	ordState := &OstateOrdered{}
	cancelState := &OstateNotSubmitted{}

	result = state.ProcessState(&ctx)

	if result != true && ctx.State != cancelState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	ctx.IgnoreSubmission = true

	result = state.ProcessState(&ctx)

	if result != true && ctx.State != ordState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	ctx.IgnoreSubmission = false
	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.Def = definition
	result = state.ProcessState(&ctx)

	if result != true && ctx.State != ordState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	result = ordState.ProcessState(&ctx)
	if result != false && ctx.IsSubmited != true {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	result = cancelState.ProcessState(&ctx)
	if result != false && ctx.IsSubmited != false {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

}
