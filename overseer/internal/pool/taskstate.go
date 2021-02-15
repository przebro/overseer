package pool

import (
	"fmt"
	"overseer/common/logger"
	"overseer/common/types"
	"overseer/common/types/date"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/taskdef"
	"strconv"
	"strings"
	"time"
)

//TaskOrderContext - use state pattern to verify if task can be added to Active Task Pool (ATP)
type TaskOrderContext struct {
	ignoreCalendar   bool
	ignoreSubmission bool
	odate            date.Odate
	currentOdate     date.Odate
	def              taskdef.TaskDefinition
	state            taskOrderState
	isSubmited       bool
	reason           []string
	log              logger.AppLogger
}

//TaskExecutionContext - context for processing active task
type TaskExecutionContext struct {
	task       *activeTask
	odate      date.Odate
	time       time.Time
	state      taskExecutionState
	dispatcher events.Dispatcher
	reason     []string
	log        logger.AppLogger
	maxRc      int32
}

type ostateCheckOtype struct {
}
type ostateCheckCalendar struct {
}
type ostateCheckSubmission struct {
}
type ostateOrdered struct {
}
type ostateNotSubmitted struct {
}

//States for task in Active pool
type ostateConfirm struct {
}
type ostateCheckTime struct {
}
type ostateCheckConditions struct {
}
type ostateHold struct {
}
type ostateStarting struct {
}
type ostateExecuting struct {
}
type ostatePostProcessing struct {
}

type taskOrderState interface {
	processState(order *TaskOrderContext) bool
}

type taskExecutionState interface {
	processState(order *TaskExecutionContext) bool
}

type taskProcessingContext interface {
	State() taskOrderState
}

func (state ostateCheckOtype) processState(ctx *TaskOrderContext) bool {

	ctx.state = &ostateCheckCalendar{}
	return true
}

func (state ostateCheckCalendar) processState(ctx *TaskOrderContext) bool {

	canceledState := &ostateNotSubmitted{}
	ctx.state = canceledState

	if ctx.ignoreCalendar == true {
		ctx.log.Debug("state check calendar processed")
		ctx.state = &ostateCheckSubmission{}
	} else {
		switch ctx.def.OrderType() {
		case taskdef.OrderingWeek:
			{
				ctx.log.Debug("checking task weekly")
				if date.IsInDayOfWeek(ctx.odate, ctx.def.Days()) && date.IsInMonth(ctx.odate, ctx.def.Months()) {
					ctx.state = &ostateCheckSubmission{}
					ctx.log.Debug("weekly task submitted")
				}
			}
		case taskdef.OrderingDaily:
			{
				ctx.log.Debug("checking task daily")
				if date.IsInMonth(ctx.odate, ctx.def.Months()) {
					ctx.state = &ostateCheckSubmission{}
					ctx.log.Debug("daily task submitted")
				}

			}
		case taskdef.OrderingDayOfMonth:
			{
				ctx.log.Debug("checking task day of month")
				if date.IsInDayOfMonth(ctx.odate, ctx.def.Days()) && date.IsInMonth(ctx.odate, ctx.def.Months()) {
					ctx.state = &ostateCheckSubmission{}
					ctx.log.Debug("day of month task submitted")
				}
			}
		case taskdef.OrderingExact:
			{
				if date.IsInExactDate(ctx.odate, ctx.def.ExactDate()) {
					ctx.state = &ostateCheckSubmission{}
				}

			}
		case taskdef.OrderingManual:
			{
				ctx.state = &ostateCheckSubmission{}
			}

		}

		if ctx.state == canceledState {
			ctx.reason = append(ctx.reason, "Scheduling criteria does not meet")
		}
	}

	return true
}
func (state ostateCheckSubmission) processState(ctx *TaskOrderContext) bool {

	ctx.log.Debug()
	if ctx.ignoreSubmission == true {
		ctx.log.Debug("check submission ingored")
		ctx.state = &ostateOrdered{}
	} else {
		if ctx.def.AllowPast() == false && date.IsBeforeCurrent(ctx.odate, ctx.currentOdate) {
			ctx.log.Debug("check submission allow past")
			ctx.reason = append(ctx.reason, "Task cannot be ordered before current day")
			ctx.state = &ostateNotSubmitted{}
		} else {
			ctx.state = &ostateOrdered{}
		}
	}

	return true
}
func (state ostateOrdered) processState(ctx *TaskOrderContext) bool {
	ctx.log.Debug("state ordered")
	ctx.isSubmited = true
	return false
}
func (state ostateNotSubmitted) processState(ctx *TaskOrderContext) bool {

	ctx.reason = append(ctx.reason, "Task not submitted")
	ctx.isSubmited = false
	return false
}

func (state ostateHold) processState(ctx *TaskExecutionContext) bool {

	if ctx.task.IsHeld() {
		ctx.task.SetWaitingInfo("task is held")
		return false
	}

	ctx.state = &ostateConfirm{}
	return true
}

func (state ostateConfirm) processState(ctx *TaskExecutionContext) bool {

	if !ctx.task.Confirmed() {

		ctx.task.SetWaitingInfo("task waiting for confirmation")
		return false
	}

	ctx.state = &ostateCheckTime{}

	return true
}
func (state ostateCheckTime) processState(ctx *TaskExecutionContext) bool {

	var isInTime bool = false
	from, to := ctx.task.TimeSpan()
	ctx.task.SetWaitingInfo(fmt.Sprintf("%s - %s", from, to))
	if from == "" && to == "" {
		ctx.state = &ostateCheckConditions{}
		return true
	}
	if from != "" {
		h, m := from.AsTime()
		//if current time is after from time task is in time window
		if !ctx.time.Before(time.Date(ctx.time.Year(), ctx.time.Month(), ctx.time.Day(), h, m, 0, 0, ctx.time.Location())) {
			isInTime = true
		}
	}
	if to != "" {
		if from == "" {
			isInTime = true
		}
		h, m := to.AsTime()
		// if current time is after to time, task must wait for time window
		if !ctx.time.Before(time.Date(ctx.time.Year(), ctx.time.Month(), ctx.time.Day(), h, m, 0, 0, ctx.time.Location())) {
			isInTime = false
		}
	}
	ctx.task.SetState(TaskStateWaiting)

	if !isInTime {

		return false
	}

	ctx.state = &ostateCheckConditions{}

	return true
}
func (state ostateCheckConditions) processState(ctx *TaskExecutionContext) bool {

	ctx.task.SetState(TaskStateWaiting)

	receiver := events.NewTicketCheckReceiver()

	msgData := events.RouteTicketCheckMsgFormat{Tickets: make([]struct {
		Name, Odate string
		Fulfilled   bool
	}, 0)}

	for _, tc := range ctx.task.Tickets() {

		msgData.Tickets = append(msgData.Tickets, struct {
			Name, Odate string
			Fulfilled   bool
		}{tc.name, tc.odate, tc.fulfilled})
	}

	msg := events.NewMsg(msgData)

	ctx.dispatcher.PushEvent(receiver, events.RouteTicketCheck, msg)

	result, err := receiver.WaitForResult()
	if err != nil {
		n, g, _ := ctx.task.GetInfo()
		ctx.log.Error(err, g, " ", n)
		return false
	}

	wconds := make([]string, 0)
	for _, t := range result.Tickets {
		wconds = append(wconds, fmt.Sprintf("%s %s:%t", t.Name, t.Odate, t.Fulfilled))
		ctx.log.Debug(t.Name, "::", t.Odate, "::", t.Fulfilled)
	}
	ctx.task.SetWaitingInfo(strings.Join(wconds, ";"))

	var fulfilled bool = false
	if len(result.Tickets) == 0 {
		fulfilled = true
	} else {

		if ctx.task.Relation() == taskdef.InTicketAND {
			fulfilled = true
		}

		for _, t := range result.Tickets {
			if ctx.task.Relation() == taskdef.InTicketAND {
				fulfilled = t.Fulfilled && fulfilled
			} else {
				fulfilled = t.Fulfilled || fulfilled
			}
		}
	}
	if fulfilled {
		ctx.state = &ostateStarting{}
	}

	return fulfilled
}

func (state ostateStarting) processState(ctx *TaskExecutionContext) bool {

	ctx.task.SetState(TaskStateStarting)
	ctx.task.SetWaitingInfo("")
	ctx.task.SetRunNumber()
	n, g, _ := ctx.task.GetInfo()
	ctx.log.Info("Launching task:", n, " group:", g, " id:", ctx.task.OrderID())
	//:TODO move to separate function

	variables := make([]taskdef.VariableData, 0)
	variables = append(variables, taskdef.VariableData{Name: "%%RN", Value: fmt.Sprintf("%d", ctx.task.RunNumber())})
	variables = append(variables, taskdef.VariableData{Name: "%%ODATE", Value: fmt.Sprintf("%s", ctx.odate.ODATE())})
	variables = append(variables, ctx.task.Variables()...)

	data := events.RouteTaskExecutionMsg{
		OrderID:   ctx.task.OrderID(),
		Type:      string(ctx.task.TypeName()),
		Variables: variables,
		Command:   ctx.task.Action(),
	}

	msg := events.NewMsg(data)
	receiver := events.NewWorkLaunchReceiver()

	ctx.dispatcher.PushEvent(receiver, events.RouteWorkLaunch, msg)

	result, err := receiver.WaitForResult()

	if err != nil {
		ctx.log.Error("State executing error:", err)
		ctx.state = &ostatePostProcessing{}
		ctx.task.SetState(TaskStateEndedNotOk)
		return true
	}

	if result.Status != types.WorkerTaskStatusStarting {
		ctx.state = &ostatePostProcessing{}
		ctx.task.SetState(TaskStateEndedNotOk)
		ctx.log.Info("Task not executed:")
		return true
	}
	ctx.task.SetState(TaskStateExecuting)
	ctx.task.SetWorkerName(result.WorkerName)
	ctx.task.SetStartTime()
	ctx.state = &ostateExecuting{}

	return true
}
func (state ostateExecuting) processState(ctx *TaskExecutionContext) bool {

	ctx.log.Info("Checking task state")
	ctx.task.SetState(TaskStateExecuting)

	msg := events.NewMsg(events.WorkRouteCheckStatusMsg{OrderID: ctx.task.OrderID(), WorkerName: ctx.task.WorkerName()})
	receiver := events.NewWorkLaunchReceiver()

	ctx.dispatcher.PushEvent(receiver, events.RouteWorkCheck, msg)

	result, err := receiver.WaitForResult()
	fmt.Println(result, err)

	if err != nil {
		ctx.log.Error("State executing error:", err)
		return false
	}

	ctx.task.AddOutput(result.Output)

	if result.Status == types.WorkerTaskStatusEnded || result.Status == types.WorkerTaskStatusFailed {

		n, g, _ := ctx.task.GetInfo()
		ctx.log.Info("Task ended:", n, " group:", g, " id:", ctx.task.OrderID(), " rc:", result.ReturnCode)

		ctx.task.SetEndTime()
		ctx.state = &ostatePostProcessing{}

		if result.Status == types.WorkerTaskStatusFailed {
			ctx.task.SetState(TaskStateEndedNotOk)
		}

		if result.Status == types.WorkerTaskStatusEnded {

			if result.ReturnCode > ctx.maxRc {
				ctx.task.SetState(TaskStateEndedNotOk)
			} else {
				ctx.task.SetState(TaskStateEndedOk)
			}
		}

		return true
	}

	return false
}
func (state ostatePostProcessing) processState(ctx *TaskExecutionContext) bool {

	msg := events.NewMsg(events.RouteTaskCleanMsg{OrderID: ctx.task.OrderID(), WorkerName: ctx.task.WorkerName(), Terminate: false})
	ctx.dispatcher.PushEvent(nil, events.RouteTaskClean, msg)

	if ctx.task.State() == TaskStateEndedNotOk {
		ctx.log.Info("Task post processing ends")
		return false
	}

	n, g, _ := ctx.task.GetInfo()
	outticket := ctx.task.TicketsOut()
	ticketMsg := events.RouteTicketInMsgFormat{Tickets: make([]struct {
		Name   string
		Odate  date.Odate
		Action taskdef.OutAction
	}, len(outticket))}

	ctx.log.Info("TASK:", n, ":", g)

	for i, t := range outticket {

		realOdat := calcRealOdate(ctx.odate, t.Odate, ctx.task.TaskDefinition.Calendar())
		ticketMsg.Tickets[i].Name = t.Name
		ticketMsg.Tickets[i].Odate = realOdat
		ticketMsg.Tickets[i].Action = t.Action
		ctx.log.Info("TICKET:", t.Name, " ", realOdat, " ", t.Action)

	}

	ctx.dispatcher.PushEvent(nil, events.RouteTicketIn, events.NewMsg(ticketMsg))

	ctx.log.Info("Task post processing ends")
	return false
}

func getProcessState(state TaskState) taskExecutionState {
	if state == TaskStateWaiting {
		return &ostateConfirm{}
	}
	if state == TaskStateExecuting {
		return &ostateExecuting{}
	}
	//Any other case means that task should not be processed.
	return nil

}

func strTimeToInt(time string) (int, int) {
	val := strings.Split(time, ":")
	h, _ := strconv.Atoi(val[0])
	m, _ := strconv.Atoi(val[1])
	return h, m

}

func calcRealOdate(current date.Odate, expect date.OdateValue, schedule taskdef.SchedulingData) date.Odate {

	mths := map[time.Month]bool{}
	for i := 1; i < 13; i++ {
		mths[time.Month(i)] = false
	}

	for _, n := range schedule.Months {
		mths[n] = true
	}

	//if expected odate value is any
	if expect == date.OdateValueAny || expect == date.OdateValueNone {
		return date.OdateNone
	}

	//For all ordering types, expect = ODATE means current day
	if expect == date.OdateValueDate {
		return current
	}

	//It is explicite date relative to current date, so do just simply compute
	if expect != date.OdateValueNext && expect != date.OdateValuePrev {

		days, _ := strconv.Atoi(string(expect))
		return date.AddDays(current, days)
	}

	//Manual ordering means that no specific schedule criteria being used, NEXT and PREV means simply tomorrow and yesterday
	if schedule.OrderType == taskdef.OrderingManual {
		days := 1
		if expect == date.OdateValuePrev {
			days = -1
			return date.AddDays(current, days)
		}
	}

	var result date.Odate

	/*for daily ordering in a simple scenario means tomorrow day for NEXT and yesterday for PREV, however,
	if PREV resolves to a day in the previous month and the execution of a task is excluded in that month then PREV will resolve
	to a day in a first included month, before the current one.
	for instance, if the task can run in months:[4,5,7]
	PREV for 05-01 means 04-30, but PREV for 07-01 means 05-31
	in particular case where months:[10], PREV for 2020-10-01 means 2019-10-31.
	The same rules apply to NEXT.
	*/
	if schedule.OrderType == taskdef.OrderingDaily {

		result = calcDateDaily(current, expect, mths)
	}

	//for exact ordering, NEXT and PREV means neighbour value in the array of specified dates
	//for edge cases(single value,last value or first value) corresponding NEXT and PREV resolves to +1/-1
	//If a task is forced to run on a non scheduled day then NEXT and PREV will resolve to +1/-1
	//for instance: if execdates: [2020-05-11] and the order date is 20200501 then NEXT will resolve to 20200511 and PREV to 20200430
	if schedule.OrderType == taskdef.OrderingExact {

		result = calcDateExact(current, expect, schedule.Exactdates)
	}

	//for weekly ordering NEXT and PREV means next or previous value from specified days of a week.
	//if it is the last specified day then the first day of a next week will be used as a value for NEXT
	//respectively, if it is the first day, then the last day from a previous week will be used for PREV
	//for instance: when a task is ordered in the day of week [1 3 4 6]
	//if it is the fourth day of a week then NEXT mean 6 and PREV means 1, however,
	//if it is the first day NEXT means 3 but PREV means 6
	if schedule.OrderType == taskdef.OrderingWeek {
		result = calcDateWeek(current, expect, schedule.Dayvalues, mths)

	}

	//day of month. Task is ordered on specfic day
	// if the task is ordered on the day of the month [31] when the date is 2020-03-31 NEXT means 2020-05-31 and PREV means 2020-01-30
	//because there is no such date like 2020-02-30 and 2020-04-31

	if schedule.OrderType == taskdef.OrderingDayOfMonth {

		result = calcDateMonth(current, expect, schedule.Dayvalues, mths)
	}

	if schedule.OrderType == taskdef.OrderingFromEnd {

		result = calcDateFromEnd(current, expect, schedule.Dayvalues, mths)
	}

	return result
}

func calcDateDaily(current date.Odate, expect date.OdateValue, mths map[time.Month]bool) date.Odate {

	days := 1
	if expect == date.OdateValuePrev {
		days *= -1
	}

	planed := current

	for {
		planed = date.AddDays(planed, days)
		py, pm, day := planed.Ymd()
		if mths[time.Month(pm)] {
			planed = date.Odate(fmt.Sprintf("%d%02d%02d", py, pm, day))
			break
		}
	}

	return planed
}

func calcDateExact(current date.Odate, expect date.OdateValue, dates []string) date.Odate {

	var idx int
	var val string
	var found bool
	cdat := current.FormatDate()

	//first check if the task was forced in current day
	for idx, val = range dates {
		if val == cdat {
			found = true
			break
		}
	}
	// task was forced on a non scheduled day or it is only single value so return -1 or +1
	if !found || len(dates) == 1 {
		days := 1
		if expect == date.OdateValuePrev {
			days = -1
		}
		return date.AddDays(current, days)
	}

	//edge case for PREV(it is first the first execution) and NEXT(it is the last execution)
	if (idx == 0 && expect == date.OdateValuePrev) || (idx == len(dates)-1 && expect == date.OdateValueNext) {

		days := 1
		if expect == date.OdateValuePrev {
			days = -1
		}
		return date.AddDays(current, days)
	}

	//for any other case it is the next or previous item from the table of values
	nextval := 1
	if expect == date.OdateValuePrev {
		nextval *= -1
	}
	return date.FromDateString(dates[idx+nextval])

}

func calcDateFromEnd(current date.Odate, expect date.OdateValue, values []int, mths map[time.Month]bool) date.Odate {

	var shift = values[0]

	cm, cy := getNextMonthYear(mths, current, expect, false)

	d := getNthLastDay(cy, cm, shift)

	return date.Odate(fmt.Sprintf("%d%02d%02d", cy, cm, d))

}

func calcDateWeek(current date.Odate, expect date.OdateValue, values []int, mths map[time.Month]bool) date.Odate {

	var idx int
	var val int
	var diffWeek int
	var expectVal int
	var found bool

	cdat := current.Wday()

	for idx, val = range values {
		if val == cdat {
			found = true
			break
		}
	}

	//Task was forced on a no scheduled day
	if !found {
		days := 1
		if expect == date.OdateValuePrev {
			days = -1
		}
		return date.AddDays(current, days)
	}

	nval := 1
	if expect == date.OdateValuePrev {
		nval *= -1
	}

	refdate := current

	for {

		idx += nval

		if idx < 0 && expect == date.OdateValuePrev {
			idx = len(values) - 1
			expectVal = values[idx]
			diffWeek = -1
		} else if idx > (len(values)-1) && expect == date.OdateValueNext {
			idx = 0
			expectVal = values[idx]
			diffWeek = 1
		} else {
			expectVal = values[idx]
			diffWeek = 0
		}

		refdate = getStartOfWeek(refdate, diffWeek)
		refdate = date.AddDays(refdate, expectVal-1)
		_, cm, _ := refdate.Ymd()

		if mths[time.Month(cm)] {
			break
		}
	}

	return refdate
}

func calcDateMonth(current date.Odate, expect date.OdateValue, values []int, mths map[time.Month]bool) date.Odate {

	var expectVal int
	var diffMonth bool
	var ndate string

	var idx int
	var val int
	var found bool

	cdat := current.Day()
	for idx, val = range values {
		if val == cdat {
			found = true
			break
		}
	}

	//Task was forced on a no scheduled day
	if !found {
		days := 1
		if expect == date.OdateValuePrev {
			days = -1
		}
		return date.AddDays(current, days)
	}

	nval := 1
	if expect == date.OdateValuePrev {
		nval *= -1
	}

	refdate := current

	for {

		idx += nval

		if idx < 0 && expect == date.OdateValuePrev {
			idx = len(values) - 1
			expectVal = values[idx]
			diffMonth = true
		} else if idx > (len(values)-1) && expect == date.OdateValueNext {
			idx = 0
			expectVal = values[idx]
			diffMonth = true
		} else {
			expectVal = values[idx]
			diffMonth = false
		}

		cm, cy := getNextMonthYear(mths, refdate, expect, !diffMonth)
		lday := getNthLastDay(cy, cm, 1)

		if expectVal <= lday {
			ndate = fmt.Sprintf("%d%02d%02d", cy, cm, expectVal)
			break
		}

		refdate = date.Odate(fmt.Sprintf("%d%02d%02d", cy, cm, lday))
	}

	return date.Odate(ndate)

}

func getNextMonthYear(mths map[time.Month]bool, current date.Odate, expect date.OdateValue, incl bool) (int, int) {

	cy, cm, _ := current.Ymd()
	nval := 1
	if expect == date.OdateValuePrev {
		nval *= -1
	}

	if incl == false {
		cm += nval
	}

	if cm < 1 {
		cm = 12
		cy--
	}
	if cm > 12 {
		cm = 1
		cy++
	}

	for mths[time.Month(cm)] != true {
		cm += nval

		if cm < 1 {
			cm = 12
			cy--
		}
		if cm > 12 {
			cm = 1
			cy++
		}
	}

	return cm, cy
}

//getNthLastDay - gets the nth last day from given year and month
func getNthLastDay(year int, month int, shift int) int {

	t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	pd := t.AddDate(0, 1, 0).AddDate(0, 0, -1).Day()

	if shift > 1 {
		pd -= (shift - 1)
	}

	return pd
}

//getStartOfWeek - gets an odate of first day(monday) in  week
func getStartOfWeek(current date.Odate, shift int) date.Odate {
	y, m, d := current.Ymd()
	wday := current.Wday()

	t := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local).AddDate(0, 0, -wday+1).AddDate(0, 0, shift*7)
	return date.FromTime(t)
}
