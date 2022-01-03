package pool

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/overseer/internal/events"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/przebro/overseer/overseer/internal/unique"

	"testing"
	"time"
)

func init() {
	if !isInitialized {
		setupEnv()
	}
}

func TestStateCheckTime(t *testing.T) {

	now := time.Now()
	nowPlus2 := now.Add(2 * time.Minute)
	nowPlus4 := now.Add(4 * time.Minute)
	nowMinus2 := now.Add(-2 * time.Minute)
	nowMinus4 := now.Add(-4 * time.Minute)

	strn := types.HourMinTime(fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()))
	strnp10 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowPlus2.Hour(), nowPlus2.Minute()))
	strnp20 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowPlus4.Hour(), nowPlus4.Minute()))
	strnm10 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowMinus2.Hour(), nowMinus2.Minute()))
	strnm20 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowMinus4.Hour(), nowMinus4.Minute()))

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{FromTime: "", ToTime: "", OrderType: taskdef.OrderingManual}).
		Build()
	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskExecutionContext{
		log:        logger.NewTestLogger(),
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       now,
	}
	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state := ostateCheckTime{}
	state.processState(&ctx)
	if ctx.isInTime != true {

		t.Error("expected result: ", true, " actual:", ctx.isInTime)
	}

	// now -> from "-" -> to "-"
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: strnp10, ToTime: "", OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state.processState(&ctx)
	if ctx.isInTime != false {
		t.Log(definition)
		t.Error("expected result: ", false, " actual:", ctx.isInTime, " ", strn, ",", strnp10)
	}

	// now -> from "-" -> to "-"
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: "", ToTime: strnp20, OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state.processState(&ctx)
	if ctx.isInTime != true {
		t.Log(definition)
		t.Error("expected result: ", true, " actual:", ctx.isInTime, " ", strn, ",", strnp10)
	}

	//from "-" -> now -> to "-"
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: strnm10, ToTime: strnp10, OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state.processState(&ctx)
	if ctx.isInTime != true {
		t.Log(definition)
		t.Error("expected result: ", true, " actual:", ctx.isInTime, " ", strn, ",", strnm10)
	}

	//from "-" -> to "-" -> now
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: strnm20, ToTime: strnm10, OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state.processState(&ctx)
	if ctx.isInTime != false {
		t.Log(definition)
		t.Error("expected result: ", false, " actual:", ctx.isInTime, " ", strn, ",", strnm10)
	}

	//from "-" -> to "-" -> now
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: "", ToTime: strnm10, OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state.processState(&ctx)
	if ctx.isInTime != false {
		t.Log(definition)
		t.Error("expected result: ", false, " actual:", ctx.isInTime, " ", strn, ",", strnm10)
	}

	//from "-" -> to "-" -> now
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: strn, ToTime: "", OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	//from now -> to "-"
	state.processState(&ctx)
	if ctx.isInTime != true {
		t.Log(definition)
		t.Error("expected result: ", true, " actual:", ctx.isInTime, " ", strn)
	}

	//from "-" -> to "-" -> now
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: "", ToTime: strn, OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	//from '-' -> to now
	state.processState(&ctx)
	if ctx.isInTime != false {

		t.Error("expected result: ", false, " actual:", ctx.isInTime, " ", strn)
	}

	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: strnp10, ToTime: strnp20, OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())
	ctx.isEnforced = true

	state.processState(&ctx)
	if ctx.isInTime != true {
		t.Log(definition)
		t.Error("expected result: ", true, " actual:", ctx.isInTime, " ", strn, ",", strnp10)
	}

}
func TestStateCheckCond(t *testing.T) {
	var result bool
	now := time.Now()

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).
		WithInTicekts([]taskdef.InTicketData{
			{
				Name: "TESTABC01", Odate: date.OdateValueDate,
			},
			{
				Name: "TESTABC02", Odate: date.OdateValueDate,
			},
		}, "AND").
		Build()
	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       now,
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	mDispatcher.withError = true
	state := ostateCheckConditions{}
	result = state.processState(&ctx)

	if result == true {
		t.Error("expected result: ", true, " actual:", result, " ")
	}

	mDispatcher.withError = false
	result = state.processState(&ctx)

	if result == true {
		t.Error("expected result: ", true, " actual:", result, " ")
	}

	if result == true {
		t.Error("expected result: ", true, " actual:", result, " ")
	}

	mDispatcher.Tickets["TESTABC01"] = string(date.CurrentOdate())

	result = state.processState(&ctx)

	if result == true {
		t.Error("expected result: ", true, " actual:", result, " ")
	}

	definition, err = builder.FromTemplate(definition).
		WithInTicekts([]taskdef.InTicketData{
			{
				Name: "TESTABC01", Odate: date.OdateValueDate,
			},
			{
				Name: "TESTABC02", Odate: date.OdateValueDate,
			},
		}, "OR").Build()

	if err != nil {
		t.Error(err)
	}
	ctx.isInTime = true
	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	result = state.processState(&ctx)

	if result != true {
		t.Error("expected result: ", true, " actual:", result, " ")
		t.Log(mDispatcher.Tickets)
		t.Log(ctx.task.TicketsIn())
	}

	mDispatcher.Tickets["TESTABC02"] = string(date.CurrentOdate())

	result = state.processState(&ctx)

	if result != true {
		t.Error("expected result: ", true, " actual:", result, " ")
		t.Log(mDispatcher.Tickets)
		t.Log(ctx.task.TicketsIn())
	}

	definition, err = builder.FromTemplate(definition).
		WithInTicekts([]taskdef.InTicketData{
			{
				Name: "TESTABC01", Odate: date.OdateValueDate,
			},
			{
				Name: "TESTABC02", Odate: date.OdateValueDate,
			},
		}, "AND").Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())
	result = state.processState(&ctx)

	if result != true {
		t.Error("expected result: ", true, " actual:", result, " ")
		t.Log(mDispatcher.Tickets)
		t.Log(ctx.task.TicketsIn())
	}

	definition, err = builder.FromTemplate(definition).
		WithInTicekts([]taskdef.InTicketData{
			{
				Name: "TESTABC01", Odate: date.OdateValueDate,
			},
			{
				Name: "TESTABC02", Odate: date.OdateValueDate,
			},
		}, "OR").Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())
	mDispatcher.Tickets = make(map[string]string)
	ctx.dispatcher = mDispatcher

	result = state.processState(&ctx)

	if result == true {
		t.Error("expected result: ", false, " actual:", result, " ")
		t.Log(mDispatcher.Tickets)
		t.Log(ctx.task.TicketsIn())

	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())
	ctx.isEnforced = true
	result = state.processState(&ctx)

	if result != true {
		t.Error("expected result: ", true, " actual:", result, " ")
		t.Log(mDispatcher.Tickets)
		t.Log(ctx.task.TicketsIn())
	}

	_, err = builder.FromTemplate(definition).
		WithInTicekts([]taskdef.InTicketData{}, "AND").Build()

	if err != nil {
		t.Error(err)
	}

}

func TestNewCyclicTask(t *testing.T) {

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithCyclic(taskdef.CyclicTaskData{
			IsCycle:      true,
			MaxRuns:      3,
			RunFrom:      taskdef.CycleFromStart,
			TimeInterval: 1,
		}).WithSchedule(taskdef.SchedulingData{
		OrderType: taskdef.OrderingDaily,
	}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())
	cycle := task.CycleData()

	if !cycle.IsCyclic {
		t.Error("unexpected value:", cycle.IsCyclic, "expected:", false)
	}
	if cycle.MaxRun != 3 {
		t.Error("unexpected value:", cycle.MaxRun, "expected:", 3)
	}

	if cycle.RunInterval != 1 {
		t.Error("unexpected value:", cycle.RunInterval, "expected:", 1)
	}

	if cycle.RunFrom != string(taskdef.CycleFromStart) {
		t.Error("unexpected value:", cycle.RunFrom, "expected:", taskdef.CycleFromStart)
	}

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       time.Now(),
	}

	ctx.task = task
	ctx.isEnforced = true
	state := ostatePostProcessing{}
	result := state.processState(&ctx)
	fmt.Println(result)

}

func TestCyclicState_Enforced(t *testing.T) {

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithCyclic(taskdef.CyclicTaskData{
			IsCycle:      true,
			MaxRuns:      3,
			RunFrom:      taskdef.CycleFromStart,
			TimeInterval: 1,
		}).WithSchedule(taskdef.SchedulingData{
		OrderType: taskdef.OrderingDaily,
	}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       time.Now(),
	}

	ctx.task = task
	ctx.isEnforced = true
	state := ostateCheckCyclic{}
	result := state.processState(&ctx)
	if !result {
		t.Error("unexpected result:", result, "expected:", true)
	}

}

func TestCyclicState_Normal(t *testing.T) {

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingDaily,
		}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       time.Now(),
	}

	ctx.task = task
	ctx.isEnforced = false
	state := ostateCheckCyclic{}
	result := state.processState(&ctx)
	if !result {
		t.Error("unexpected result:", result, "expected:", true)
	}
}

func TestCyclicState_BeforeNextRun(t *testing.T) {

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithCyclic(taskdef.CyclicTaskData{
			IsCycle:      true,
			MaxRuns:      3,
			RunFrom:      taskdef.CycleFromStart,
			TimeInterval: 1,
		}).WithSchedule(taskdef.SchedulingData{
		OrderType: taskdef.OrderingDaily,
	}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())
	now := time.Now().Add(2 * time.Minute)
	task.cycle.NextRun = types.FromTime(now)

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       time.Now(),
	}

	ctx.task = task
	ctx.isEnforced = false
	state := ostateCheckCyclic{}
	result := state.processState(&ctx)
	if result {
		t.Error("unexpected result:", result, "expected:", false)
	}

}

func TestStateOrderState(t *testing.T) {

	var result bool
	var definition taskdef.TaskDefinition
	var err error

	builder := taskdef.DummyTaskBuilder{}
	definition, err = builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingManual}).Build()

	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskOrderContext{
		log:              log,
		odate:            date.CurrentOdate(),
		ignoreCalendar:   false,
		ignoreSubmission: false,
		def:              definition,
		currentOdate:     date.CurrentOdate(),
	}

	state := ostateCheckOtype{}
	stchkcal := &ostateCheckCalendar{}

	state.processState(&ctx)
	result = state.processState(&ctx)

	if result != true && ctx.state != stchkcal {
		t.Error("expected result: ", true, " actual:", result, " ", stchkcal)
	}

	_, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	result = state.processState(&ctx)

	if result != true && ctx.state != stchkcal {
		t.Error("expected result: ", true, " actual:", result, " ", stchkcal)
	}
}
func TestStateConfirm(t *testing.T) {

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").WithConfirm().
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       time.Now(),
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state := &ostateConfirm{}

	result := state.processState(&ctx)

	if result == true {
		t.Error("expected result: ", false, " actual:", result)
	}

	definition, err = builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	result = state.processState(&ctx)

	if result != true {
		t.Error("expected result: ", true, " actual:", result)
	}
}

func TestStateCheckCalendar(t *testing.T) {

	var result bool

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).Build()

	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskOrderContext{
		log:              log,
		odate:            date.CurrentOdate(),
		ignoreCalendar:   false,
		ignoreSubmission: false,
		def:              definition,
		currentOdate:     date.CurrentOdate(),
	}

	state := &ostateCheckCalendar{}
	submState := &ostateCheckSubmission{}
	cancelState := &ostateNotSubmitted{}

	result = state.processState(&ctx)

	if result != true && ctx.state != cancelState {
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

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != submState {
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

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != submState {
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

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != cancelState {
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

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != submState {
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

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != cancelState {
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

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != submState {
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

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != cancelState {
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

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != submState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	ctx.ignoreCalendar = true
	result = state.processState(&ctx)

	if result != true && ctx.state != submState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

}
func TestStateCheckSubmmision(t *testing.T) {

	var result bool

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily, AllowPastSub: false}).Build()

	if err != nil {
		t.Fatal("Unable to construct task")
	}

	subodat := date.AddDays(date.CurrentOdate(), -1)

	if err != nil {
		t.Error(err)
	}

	ctx := TaskOrderContext{
		log:              log,
		odate:            subodat,
		ignoreCalendar:   false,
		ignoreSubmission: false,
		def:              definition,
		currentOdate:     date.CurrentOdate(),
	}

	state := &ostateCheckSubmission{}
	ordState := &ostateOrdered{}
	cancelState := &ostateNotSubmitted{}

	result = state.processState(&ctx)

	if result != true && ctx.state != cancelState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	ctx.ignoreSubmission = true

	result = state.processState(&ctx)

	if result != true && ctx.state != ordState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	ctx.ignoreSubmission = false
	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily, AllowPastSub: true}).Build()

	if err != nil {
		t.Error(err)
	}

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != ordState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	result = ordState.processState(&ctx)
	if result != false && ctx.isSubmited != true {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	result = cancelState.processState(&ctx)
	if result != false && ctx.isSubmited != false {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

}
func TestStatesExecEndHold(t *testing.T) {

	var result bool

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily, AllowPastSub: false}).
		WithOutTickets([]taskdef.OutTicketData{{Action: "REM", Name: "TEST", Odate: "ODATE"}}).
		Build()

	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: mDispatcher,
		time:       time.Now(),
	}

	ctx.task = newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	state := &ostateStarting{}
	runState := &ostateExecuting{}
	postState := &ostatePostProcessing{}

	result = state.processState(&ctx)

	if result != true && ctx.state != runState {
		t.Error("expected result: ", true, " actual:", result, " ", ctx.state)
	}

	mDispatcher.processNotEnded = true
	mDispatcher.withError = false
	result = runState.processState(&ctx)

	if result != false && ctx.state != runState {
		t.Error("expected result: ", true, " actual:", result, " ", ctx.state)
	}

	mDispatcher.withError = true
	result = runState.processState(&ctx)

	if result != false && ctx.state != runState {
		t.Error("expected result: ", false, " actual:", result, " ", ctx.state)
	}

	mDispatcher.withError = false
	mDispatcher.processNotEnded = false
	result = postState.processState(&ctx)
	if result != false {
		t.Error("expected result: ", false, " actual:", result)
	}
	ctx.task.SetState(TaskStateEndedOk)

	if ctx.task.State() != TaskStateEndedOk {
		t.Error("expected result: ", TaskStateEndedOk, " actual:", ctx.task.State())
	}

	if ctx.task.OrderDate() != date.CurrentOdate() {
		t.Error("expected result: ", date.CurrentOdate(), " actual:", ctx.task.OrderDate())
	}

	if ctx.task.RunNumber() != 1 {
		t.Error("expected result: ", 1, " actual:", ctx.task.RunNumber())
	}

	ctx.task.Hold()
	if ctx.task.IsHeld() != true {
		t.Error("expected result: ", true, " actual:", ctx.task.IsHeld())
	}

	ctx.task.Free()
	if ctx.task.IsHeld() != false {
		t.Error("expected result: ", false, " actual:", ctx.task.IsHeld())
	}

}

func TestGetProcState(t *testing.T) {

	stateWaiting := &ostateConfirm{}
	stateExecute := &ostateExecuting{}

	res := getProcessState(TaskStateWaiting, false)
	if res == nil && res != stateWaiting {
		t.Error("expected result: ", stateWaiting, " actual:", res)
	}

	res = getProcessState(TaskStateExecuting, false)
	if res == nil && res != stateExecute {
		t.Error("expected result: ", stateExecute, " actual:", res)
	}

	res = getProcessState(TaskStateStarting, false)
	if res != nil {
		t.Error("expected result: ", nil, " actual:", res)
	}

	res = getProcessState(TaskStateWaiting, true)
	if res != nil {
		t.Error("expected result: ", nil, " actual:", res)
	}

}

func TestStrTimeToInt(t *testing.T) {

	h, m := strTimeToInt("12:23")
	if h != 12 && m != 23 {
		t.Error("unexpected result")
	}
}

func TestCalcRealOdate(t *testing.T) {

	current := date.Odate("20201101")
	otype := []taskdef.SchedulingOption{
		taskdef.OrderingManual,
		taskdef.OrderingDaily,
		taskdef.OrderingWeek,
		taskdef.OrderingDayOfMonth,
		taskdef.OrderingExact,
		taskdef.OrderingFromEnd,
	}
	schedule := taskdef.SchedulingData{}
	for x := range otype {
		schedule.OrderType = otype[x]
		//ODATE means same date as current
		result := calcRealOdate(current, date.OdateValueDate, schedule)
		if result != current {
			t.Error("unexpected error, expecting odate:", current, "actual:", result)
		}
	}

	current = date.Odate("20201101")
	schedule.OrderType = taskdef.OrderingManual
	//Tasks ordered manually don't depend on calendar so PREV or NEXT means -1 or +1
	result := calcRealOdate(current, date.OdateValueAny, schedule)
	if result != "" {
		t.Error("unexpected result:", result, "expected empty")
	}

	current = date.Odate("20201101")
	schedule.OrderType = taskdef.OrderingManual
	//Tasks ordered manually don't depend on calendar so PREV or NEXT means -1 or +1
	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != "20201031" {
		t.Error("unexpected result:", result, "expected 20201001")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != "20201102" {
		t.Error("unexpected result:", result, "expected 20201102")
	}

	//For exact date with only a single value there NEXT and PREV should resolve to +1 and +1
	schedule.OrderType = taskdef.OrderingExact
	schedule.Exactdates = []string{"2020-11-01"}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != "20201102" {
		t.Error("unexpected result:", result, "expected 20201102")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != "20201031" {
		t.Error("unexpected result:", result, "expected 20201030")
	}

	//For exact date with current date first,  NEXT should be computed to defined value and PREV to -1
	schedule.OrderType = taskdef.OrderingExact
	schedule.Exactdates = []string{"2020-11-01", "2020-11-15"}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20201115") {
		t.Error("unexpected result:", result, "expected empty odate")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != "20201031" {
		t.Error("unexpected result:", result, "expected 20201031")
	}

	//For exact date with current date first,  PREV should be computed to defined value and NEXT to -1
	schedule.OrderType = taskdef.OrderingExact
	schedule.Exactdates = []string{"2020-10-15", "2020-11-01"}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != "20201102" {
		t.Error("unexpected result:", result, "expected 20201102")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20201015") {
		t.Error("unexpected result:", result, "expected empty odate")
	}

	//For exact date with current date in middle of values both PREV and NEXT should be Computed
	schedule.OrderType = taskdef.OrderingExact
	schedule.Exactdates = []string{"2020-10-15", "2020-11-01", "2020-11-15"}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20201115") {
		t.Error("unexpected result:", result, "expected 20201115")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20201015") {
		t.Error("unexpected result:", result, "expected 20201015")
	}

	//For a task that was forced on a non scheduled day NEXT and PERV should be computed to +1,-1
	schedule.OrderType = taskdef.OrderingExact
	schedule.Exactdates = []string{"2020-10-15", "2020-11-01"}
	current = date.Odate("20201022")

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20201023") {
		t.Error("unexpected result:", result, "expected 20201101")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20201021") {
		t.Error("unexpected result:", result, "expected 20201015")
	}

	current = date.Odate("20200923")

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20200924") {
		t.Error("unexpected result:", result, "expected 20201015")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != "20200922" {
		t.Error("unexpected result:", result, "expected 20200922")
	}

}

func TestCalcRealOdateFromEnd(t *testing.T) {

	schedule := taskdef.SchedulingData{}

	schedule.OrderType = taskdef.OrderingFromEnd
	schedule.Dayvalues = []int{1}
	schedule.Months = []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

	current := date.Odate("20201231")

	//from end in simple scenario it simply nth day of a next and previous month if they are available
	result := calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210131") {
		t.Error("unexpected result:", result, "expected empty odate")
	}

	//from end in simple scenario it simply nth day of a next and previous month if they are available
	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20201130") {
		t.Error("unexpected result:", result, "expected empty odate")
	}

	schedule.OrderType = taskdef.OrderingFromEnd
	schedule.Dayvalues = []int{1}
	schedule.Months = []time.Month{2, 4, 7}

	current = date.Odate("20200430")

	//from end in simple scenario it simply nth day of a next and previous month if they are available
	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20200731") {
		t.Error("unexpected result:", result, "expected 20200731")
	}

	//from end in simple scenario it simply nth day of a next and previous month if they are available
	//test leap year
	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20200229") {
		t.Error("unexpected result:", result, "expected 20200229")
	}

	current = date.Odate("20210430")

	//from end in simple scenario it simply nth day of a next and previous month if they are available
	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210731") {
		t.Error("unexpected result:", result, "expected 20210731")
	}

	//from end in simple scenario it simply nth day of a next and previous month if they are available
	//test leap year
	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210228") {
		t.Error("unexpected result:", result, "expected 20210228")
	}

	schedule.OrderType = taskdef.OrderingFromEnd
	schedule.Dayvalues = []int{2}
	schedule.Months = []time.Month{2, 4, 7}
	current = date.Odate("20210228")

	//from end in simple scenario it simply nth day of a next and previous month if they are available
	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210429") {
		t.Error("unexpected result:", result, "expected 20210731")
	}

	//from end in simple scenario it simply nth day of a next and previous month if they are available
	//test leap year
	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20200730") {
		t.Error("unexpected result:", result, "expected 20210228")
	}
}

func TestCalcRealOdateWeek(t *testing.T) {

	schedule := taskdef.SchedulingData{}
	schedule.OrderType = taskdef.OrderingWeek
	schedule.Dayvalues = []int{4, 5}
	schedule.Months = []time.Month{3, 7}

	current := date.Odate("20210701")

	result := calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210702") {
		t.Error("unexpected result:", result, "expected 20210702")
	}

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210326") {
		t.Error("unexpected result:", result, "expected 20210326")
	}

	schedule.Dayvalues = []int{5, 7}

	current = date.Odate("20210314")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210312") {
		t.Error("unexpected result:", result, "expected 20210312")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210319") {
		t.Error("unexpected result:", result, "expected 20210319")
	}

	schedule.Dayvalues = []int{5, 7}
	current = date.Odate("20210328")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210326") {
		t.Error("unexpected result:", result, "expected 20210326")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210702") {
		t.Error("unexpected result:", result, "expected 20210702")
	}

	schedule.Dayvalues = []int{5, 7}
	current = date.Odate("20210326")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210321") {
		t.Error("unexpected result:", result, "expected 20210321")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210328") {
		t.Error("unexpected result:", result, "expected 20210328")
	}

	schedule.Dayvalues = []int{5, 7}
	schedule.Months = []time.Month{1, 3}
	current = date.Odate("20210328")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210326") {
		t.Error("unexpected result:", result, "expected 20210326")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20220102") {
		t.Error("unexpected result:", result, "expected 20220102")
	}

	schedule.Dayvalues = []int{2, 7}
	schedule.Months = []time.Month{1, 3}
	current = date.Odate("20210324")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210323") {
		t.Error("unexpected result:", result, "expected 20210323")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210325") {
		t.Error("unexpected result:", result, "expected 20210325")
	}
}

func TestCalcRealOdateMonth(t *testing.T) {
	schedule := taskdef.SchedulingData{}
	schedule.OrderType = taskdef.OrderingDayOfMonth
	schedule.Dayvalues = []int{30}
	schedule.Months = []time.Month{1, 2, 3, 7}
	current := date.Odate("20210330")

	result := calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210130") {
		t.Error("unexpected result:", result, "expected 20210130")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210730") {
		t.Error("unexpected result:", result, "expected 20210730")
	}

	//non scheduled day
	current = date.Odate("20210328")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210327") {
		t.Error("unexpected result:", result, "expected 20210327")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210329") {
		t.Error("unexpected result:", result, "expected 20210329")
	}

	schedule.Dayvalues = []int{31}
	schedule.Months = []time.Month{1, 2, 4, 7}
	current = date.Odate("20210731")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210131") {
		t.Error("unexpected result:", result, "expected 20210131")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20220131") {
		t.Error("unexpected result:", result, "expected 20220131")
	}

	schedule.Dayvalues = []int{29, 30, 31}
	schedule.Months = []time.Month{1, 2, 4}
	current = date.Odate("20210131")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210130") {
		t.Error("unexpected result:", result, "expected 20210129")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210429") {
		t.Error("unexpected result:", result, "expected 20210429")
	}

	schedule.Dayvalues = []int{29, 30, 31}
	schedule.Months = []time.Month{1, 2, 4}
	current = date.Odate("20200131")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20200130") {
		t.Error("unexpected result:", result, "expected 20210129")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20200229") {
		t.Error("unexpected result:", result, "expected 20220229")
	}

	schedule.Dayvalues = []int{1, 17, 28, 29, 30, 31}
	schedule.Months = []time.Month{1, 2, 4}
	current = date.Odate("20210228")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210217") {
		t.Error("unexpected result:", result, "expected 20210217")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210401") {
		t.Error("unexpected result:", result, "expected 20210401")
	}

	schedule.Dayvalues = []int{2, 17, 28, 29, 30, 31}
	schedule.Months = []time.Month{1, 4, 11}
	current = date.Odate("20210102")

	result = calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20201130") {
		t.Error("unexpected result:", result, "expected 20201130")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210117") {
		t.Error("unexpected result:", result, "expected 20210217")
	}

}

func TestCalcRealOdateDaily(t *testing.T) {
	schedule := taskdef.SchedulingData{}
	schedule.OrderType = taskdef.OrderingDaily
	schedule.Dayvalues = []int{30}
	schedule.Months = []time.Month{1, 3, 7}
	current := date.Odate("20210301")

	result := calcRealOdate(current, date.OdateValuePrev, schedule)
	if result != date.Odate("20210131") {
		t.Error("unexpected result:", result, "expected 20210131")
	}

	result = calcRealOdate(current, date.OdateValueNext, schedule)
	if result != date.Odate("20210302") {
		t.Error("unexpected result:", result, "expected 20210302")
	}
}

func TestCalcRealOdateRelative(t *testing.T) {

	current := date.Odate("20201105")
	otype := []taskdef.SchedulingOption{
		taskdef.OrderingManual,
		taskdef.OrderingDaily,
		taskdef.OrderingWeek,
		taskdef.OrderingDayOfMonth,
		taskdef.OrderingExact,
		taskdef.OrderingFromEnd,
	}
	schedule := taskdef.SchedulingData{}
	for x := range otype {
		schedule.OrderType = otype[x]

		result := calcRealOdate(current, "-001", schedule)
		if result != "20201104" {
			t.Error("unexpected error, expecting odate:", "20201104", "actual:", result)
		}

		result = calcRealOdate(current, "+002", schedule)
		if result != "20201107" {
			t.Error("unexpected error, expecting odate:", "20201107", "actual:", result)
		}
	}
}

func TestGetFirstDay(t *testing.T) {

	dd := date.Odate("20210401")
	getStartOfWeek(dd, 0)

}

func Test_prepareNextCycle_From_Start(t *testing.T) {

	interval := 2
	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithCyclic(taskdef.CyclicTaskData{
			IsCycle:      true,
			MaxRuns:      3,
			RunFrom:      taskdef.CycleFromStart,
			TimeInterval: interval,
		}).WithSchedule(taskdef.SchedulingData{
		OrderType: taskdef.OrderingDaily,
	}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	tm := task.SetStartTime()
	expect := tm.Add(time.Duration(interval) * time.Minute)
	task.SetEndTime()

	task.prepareNextCycle()

	if task.CycleData().NextRun != types.FromTime(expect) {
		t.Error("unexpected value:", task.CycleData().NextRun, "expected:", types.FromTime(expect))
	}
}

func Test_prepareNextCycle_From_End(t *testing.T) {

	interval := 2
	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithCyclic(taskdef.CyclicTaskData{
			IsCycle:      true,
			MaxRuns:      3,
			RunFrom:      taskdef.CycleFromEnd,
			TimeInterval: interval,
		}).WithSchedule(taskdef.SchedulingData{
		OrderType: taskdef.OrderingDaily,
	}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	task.SetStartTime()
	tm := task.SetEndTime()
	expect := tm.Add(time.Duration(interval) * time.Minute)

	task.prepareNextCycle()

	if task.CycleData().NextRun != types.FromTime(expect) {
		t.Error("unexpected value:", task.CycleData().NextRun, "expected:", types.FromTime(expect))
	}
}

func Test_prepareNextCycle_From_Schedule(t *testing.T) {

	interval := 2
	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithCyclic(taskdef.CyclicTaskData{
			IsCycle:      true,
			MaxRuns:      3,
			RunFrom:      taskdef.CycleFromSchedule,
			TimeInterval: interval,
		}).WithSchedule(taskdef.SchedulingData{
		OrderType: taskdef.OrderingDaily,
	}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	task.SetStartTime()
	tm := task.SetEndTime()
	expect := tm.Add(time.Duration(interval) * time.Minute)

	task.prepareNextCycle()

	if task.CycleData().NextRun != types.FromTime(expect) {
		t.Error("unexpected value:", task.CycleData().NextRun, "expected:", types.FromTime(expect))
	}
}

func Test_prepareNextCycle_MaxRuns(t *testing.T) {

	interval := 2
	maxRuns := 3
	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "cyclic_test_01", "test task").
		WithCyclic(taskdef.CyclicTaskData{
			IsCycle:      true,
			MaxRuns:      maxRuns,
			RunFrom:      taskdef.CycleFromStart,
			TimeInterval: interval,
		}).WithSchedule(taskdef.SchedulingData{
		OrderType: taskdef.OrderingDaily,
	}).Build()

	if err != nil {
		t.Log(err)
	}

	task := newActiveTask(seq.Next(), date.CurrentOdate(), definition, unique.NewID())

	for i := 1; i <= maxRuns; i++ {

		task.SetStartTime()
		task.SetEndTime()
		result := task.prepareNextCycle()
		if !result && i < maxRuns {
			t.Error("unexpected value:", result, "expected:", true)
		}

		if result && i == maxRuns {
			t.Error("unexpected value:", result, "expected:", false)
		}

		task.SetExecutionID()
	}
}

func Test_buildFlagMsg_Error(t *testing.T) {

	result := buildFlagMsg(nil)
	if result != nil {
		t.Error("unexpected result:", result, "expected:", nil)
	}

	data := []taskdef.FlagData{}
	result = buildFlagMsg(data)
	if result != nil {
		t.Error("unexpected result:", result, "expected:", nil)
	}
}

func Test_buildFlagMsg(t *testing.T) {

	data := []taskdef.FlagData{
		{Name: "ABC", Type: taskdef.FlagShared},
		{Name: "ABCD", Type: taskdef.FlagExclusive},
		{Name: "ABCE", Type: taskdef.FlagShared},
	}

	result := buildFlagMsg(data)
	if result == nil {
		t.Error("unexpected result:", result, "expected:", nil)
	}
	msg, ok := result.Message().(events.RouteFlagAcquireMsg)
	if ok != true {
		t.Error("unexpected result:", ok, "expected:", true)
	}

	if len(msg.Flags) != len(data) {
		t.Error("unexpected result:", len(msg.Flags), "expected:", len(data))
	}
}

func strTimeToInt(time string) (int, int) {
	val := strings.Split(time, ":")
	h, _ := strconv.Atoi(val[0])
	m, _ := strconv.Atoi(val[1])
	return h, m

}
