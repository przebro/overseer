package pool

import (
	"errors"
	"fmt"
	"goscheduler/common/logger"
	"goscheduler/common/types"
	"goscheduler/overseer/internal/date"
	"goscheduler/overseer/internal/events"
	"goscheduler/overseer/internal/taskdef"

	"goscheduler/overseer/internal/unique"
	"testing"
	"time"
)

type mockDispatcher struct {
	Tickets         map[string]string
	processNotEnded bool
	withError       bool
}

var log logger.AppLogger = logger.NewTestLogger()

func (m *mockDispatcher) PushEvent(receiver events.EventReceiver, route events.RouteName, msg events.DispatchedMessage) error {

	go func() {
		if route == events.RouteWorkLaunch {
			if m.withError {
				receiver.Done(errors.New(""))
			} else {
				dat := events.RouteWorkResponseMsg{
					Started: true,
				}
				events.ResponseToReceiver(receiver, dat)
			}

		}
		if route == events.RouteTicketIn {

		}
		if route == events.RouteWorkCheck {

			if m.withError {
				receiver.Done(errors.New(""))
			} else {

				_, iskOk := msg.Message().(events.WorkRouteCheckStatusMsg)
				if iskOk == false {
					events.ResponseToReceiver(receiver, errors.New(""))
				}
				if m.processNotEnded {
					receiver.Done(events.RouteWorkResponseMsg{Ended: false, ReturnCode: 0})
				} else {
					receiver.Done(events.RouteWorkResponseMsg{Ended: true, ReturnCode: 0})
				}

			}
		}
		if route == events.RouteTicketCheck {

			if m.withError {

				receiver.Done(errors.New(""))

			} else {
				result, iskOk := msg.Message().(events.RouteTicketCheckMsgFormat)
				if iskOk == false {
					events.ResponseToReceiver(receiver, errors.New(""))
				}

				for i, t := range result.Tickets {

					_, exists := m.Tickets[t.Name]
					if exists {
						result.Tickets[i].Fulfilled = true
					}
				}

				receiver.Done(result)

			}

		}
	}()
	return nil
}
func (m *mockDispatcher) Subscribe(route events.RouteName, participant events.EventParticipant) {

}
func (m *mockDispatcher) Unsubscribe(route events.RouteName, participant events.EventParticipant) {

}

var disp = &mockDispatcher{Tickets: make(map[string]string, 0)}

func TestStateCheckTime(t *testing.T) {

	var result bool
	now := time.Now()
	nowPlus10 := now.Add(10 * time.Minute)
	nowPlus20 := now.Add(20 * time.Minute)
	nowMinus10 := now.Add(-20 * time.Minute)
	nowMinus20 := now.Add(-20 * time.Minute)

	strn := types.HourMinTime(fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute()))
	strnp10 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowPlus10.Hour(), nowPlus10.Minute()))
	strnp20 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowPlus20.Hour(), nowPlus20.Minute()))
	strnm10 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowMinus10.Hour(), nowMinus10.Minute()))
	strnm20 := types.HourMinTime(fmt.Sprintf("%02d:%02d", nowMinus20.Hour(), nowMinus20.Minute()))

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{FromTime: "", ToTime: ""}).
		Build()
	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskExecutionContext{
		log:        logger.NewTestLogger(),
		odate:      date.CurrentOdate(),
		dispatcher: disp,
		time:       now,
	}
	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	state := ostateCheckTime{}
	result = state.processState(&ctx)
	if result != true {

		t.Error("expected result: ", true, " actual:", result)
	}

	// now -> from "-" -> to "-"
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: strnp10, ToTime: ""}).Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	result = state.processState(&ctx)
	if result == true {
		t.Log(definition)
		t.Error("expected result: ", false, " actual:", result, " ", strn, ",", strnp10)
	}

	// now -> from "-" -> to "-"
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: "", ToTime: strnp20}).Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	result = state.processState(&ctx)
	if result != true {
		t.Log(definition)
		t.Error("expected result: ", true, " actual:", result, " ", strn, ",", strnp10)
	}

	//from "-" -> now -> to "-"
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: strnm10, ToTime: strnp10}).Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	result = state.processState(&ctx)
	if result != true {
		t.Log(definition)
		t.Error("expected result: ", true, " actual:", result, " ", strn, ",", strnm10)
	}

	//from "-" -> to "-" -> now
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: strnm20, ToTime: strnm10}).Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	result = state.processState(&ctx)
	if result == true {
		t.Log(definition)
		t.Error("expected result: ", false, " actual:", result, " ", strn, ",", strnm10)
	}

	//from "-" -> to "-" -> now
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: "", ToTime: strnm10}).Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	result = state.processState(&ctx)
	if result == true {
		t.Log(definition)
		t.Error("expected result: ", false, " actual:", result, " ", strn, ",", strnm10)
	}

	//from "-" -> to "-" -> now
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: strn, ToTime: ""}).Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	//from now -> to "-"
	result = state.processState(&ctx)
	if result != true {
		t.Log(definition)
		t.Error("expected result: ", true, " actual:", result, " ", strn)
	}

	//from "-" -> to "-" -> now
	definition, err = builder.FromTemplate(definition).WithSchedule(
		taskdef.SchedulingData{FromTime: "", ToTime: strn}).Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	//from '-' -> to now
	result = state.processState(&ctx)
	if result == true {

		t.Error("expected result: ", true, " actual:", result, " ", strn)
	}

}
func TestStateCheckCond(t *testing.T) {
	var result bool
	now := time.Now()

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").
		WithInTicekts([]taskdef.InTicketData{
			taskdef.InTicketData{
				Name: "TESTABC01", Odate: date.OdateValueDate,
			},
			taskdef.InTicketData{
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
		dispatcher: disp,
		time:       now,
	}

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	disp.withError = true
	state := ostateCheckConditions{}
	result = state.processState(&ctx)

	if result == true {
		t.Error("expected result: ", true, " actual:", result, " ")
	}

	disp.withError = false
	result = state.processState(&ctx)

	if result == true {
		t.Error("expected result: ", true, " actual:", result, " ")
	}

	if result == true {
		t.Error("expected result: ", true, " actual:", result, " ")
	}

	disp.Tickets["TESTABC01"] = string(date.CurrentOdate())

	result = state.processState(&ctx)

	if result == true {
		t.Error("expected result: ", true, " actual:", result, " ")
	}

	definition, err = builder.FromTemplate(definition).
		WithInTicekts([]taskdef.InTicketData{
			taskdef.InTicketData{
				Name: "TESTABC01", Odate: date.OdateValueDate,
			},
			taskdef.InTicketData{
				Name: "TESTABC02", Odate: date.OdateValueDate,
			},
		}, "OR").Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	result = state.processState(&ctx)

	if result != true {
		t.Error("expected result: ", true, " actual:", result, " ")
		t.Log(disp.Tickets)
		t.Log(ctx.task.TicketsIn())
	}

	disp.Tickets["TESTABC02"] = string(date.CurrentOdate())

	result = state.processState(&ctx)

	if result != true {
		t.Error("expected result: ", true, " actual:", result, " ")
		t.Log(disp.Tickets)
		t.Log(ctx.task.TicketsIn())
	}

	definition, err = builder.FromTemplate(definition).
		WithInTicekts([]taskdef.InTicketData{
			taskdef.InTicketData{
				Name: "TESTABC01", Odate: date.OdateValueDate,
			},
			taskdef.InTicketData{
				Name: "TESTABC02", Odate: date.OdateValueDate,
			},
		}, "AND").Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)
	result = state.processState(&ctx)

	fmt.Println(ctx.task.TicketsIn())

	if result != true {
		t.Error("expected result: ", true, " actual:", result, " ")
		t.Log(disp.Tickets)
		t.Log(ctx.task.TicketsIn())
	}

	definition, err = builder.FromTemplate(definition).
		WithInTicekts([]taskdef.InTicketData{
			taskdef.InTicketData{
				Name: "TESTABC01", Odate: date.OdateValueDate,
			},
			taskdef.InTicketData{
				Name: "TESTABC02", Odate: date.OdateValueDate,
			},
		}, "OR").Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)
	disp.Tickets = make(map[string]string, 0)
	ctx.dispatcher = disp

	result = state.processState(&ctx)

	if result == true {
		t.Error("expected result: ", false, " actual:", result, " ")
		t.Log(disp.Tickets)
		t.Log(ctx.task.TicketsIn())

	}

}
func TestStateOrderState(t *testing.T) {

	var result bool

	builder := taskdef.DummyTaskBuilder{}
	definition, err := builder.WithBase("test", "dummy_04", "test task").
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
	result = state.processState(&ctx)

	result = state.processState(&ctx)

	if result != true && ctx.state != stchkcal {
		t.Error("expected result: ", true, " actual:", result, " ", stchkcal)
	}

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).Build()

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
		dispatcher: disp,
		time:       time.Now(),
	}

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	state := &ostateConfirm{}

	result := state.processState(&ctx)

	if result != false {
		t.Error("expected result: ", false, " actual:", result)
	}

	definition, err = builder.WithBase("test", "dummy_04", "test task").
		WithSchedule(taskdef.SchedulingData{OrderType: taskdef.OrderingDaily}).Build()

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

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
			Values:    []string{fmt.Sprintf("%d", wday)},
		}).Build()

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != submState {
		t.Error("expected result: ", true, " actual:", result, " ", submState)
	}

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingWeek,
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Values:    []string{fmt.Sprintf("%d", nday)},
		}).Build()

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
			Values:    []string{fmt.Sprintf("%02d", dofmonth)},
		}).Build()

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != submState {
		t.Error("expected result: ", true, " actual:", result, " ", submState)
	}

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingDayOfMonth,
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Values:    []string{fmt.Sprintf("%02d", ndofmonth)},
		}).Build()

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != cancelState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	/**/
	edate := date.CurrentOdate().FormatDate()
	odate2, _ := date.AddDays(date.CurrentOdate(), 2)
	nedate := odate2.FormatDate()

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingExact,
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Values:    []string{edate},
		}).Build()

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != submState {
		t.Error("expected result: ", true, " actual:", result, " ", submState)
	}

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingExact,
			Months:    []time.Month{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12},
			Values:    []string{nedate},
		}).Build()

	ctx.def = definition
	result = state.processState(&ctx)

	if result != true && ctx.state != cancelState {
		t.Error("expected result: ", true, " actual:", result, " ", cancelState)
	}

	definition, err = builder.FromTemplate(definition).
		WithSchedule(taskdef.SchedulingData{
			OrderType: taskdef.OrderingManual,
			Months:    []time.Month{},
			Values:    []string{nedate},
		}).Build()

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

	subodat, err := date.AddDays(date.CurrentOdate(), -1)
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
		WithOutTickets([]taskdef.OutTicketData{taskdef.OutTicketData{Action: "REM", Name: "TEST", Odate: "ODATE"}}).
		Build()

	if err != nil {
		t.Fatal("Unable to construct task")
	}

	ctx := TaskExecutionContext{
		log:        log,
		odate:      date.CurrentOdate(),
		dispatcher: disp,
		time:       time.Now(),
	}

	ctx.task = newActiveTask(unique.NewOrderID(), date.CurrentOdate(), definition)

	state := &ostateStarting{}
	runState := &ostateExecuting{}
	postState := &ostatePostProcessing{}
	holdState := &ostateHold{}

	result = state.processState(&ctx)

	if result != true && ctx.state != runState {
		t.Error("expected result: ", true, " actual:", result, " ", ctx.state)
	}

	disp.processNotEnded = true
	disp.withError = false
	result = runState.processState(&ctx)

	if result != false && ctx.state != runState {
		t.Error("expected result: ", true, " actual:", result, " ", ctx.state)
	}

	disp.withError = true
	result = runState.processState(&ctx)

	if result != false && ctx.state != runState {
		t.Error("expected result: ", false, " actual:", result, " ", ctx.state)
	}

	disp.withError = false
	disp.processNotEnded = false
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

	result = holdState.processState(&ctx)

	if result != false {
		t.Error("expected result: ", false, " actual:", result)
	}

	ctx.task.Hold()
	if ctx.task.State() != TaskStateHold {
		t.Error("expected result: ", TaskStateHold, " actual:", ctx.task.State())
	}

	ctx.task.Free()
	if ctx.task.State() != TaskStateWaiting {
		t.Error("expected result: ", TaskStateWaiting, " actual:", ctx.task.State())
	}

}

func TestGetProcState(t *testing.T) {

	stateWaiting := &ostateConfirm{}
	stateExecute := &ostateExecuting{}
	stateHold := &ostateHold{}

	res := getProcessState(TaskStateWaiting)
	if res == nil && res != stateWaiting {
		t.Error("expected result: ", stateWaiting, " actual:", res)
	}

	res = getProcessState(TaskStateExecuting)
	if res == nil && res != stateExecute {
		t.Error("expected result: ", stateExecute, " actual:", res)
	}

	res = getProcessState(TaskStateHold)
	if res == nil && res != stateHold {
		t.Error("expected result: ", stateHold, " actual:", res)
	}

	res = getProcessState(TaskStateStarting)
	if res != nil {
		t.Error("expected result: ", nil, " actual:", res)
	}

}
