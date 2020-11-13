package pool

import (
	"fmt"
	"goscheduler/common/logger"
	"goscheduler/overseer/internal/date"
	"goscheduler/overseer/internal/events"
	"goscheduler/overseer/internal/taskdef"
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
				if date.IsInDayOfWeek(ctx.odate, ctx.def.Values()) && date.IsInMonth(ctx.odate, ctx.def.Months()) {
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
				if date.IsInDayOfMonth(ctx.odate, ctx.def.Values()) && date.IsInMonth(ctx.odate, ctx.def.Months()) {
					ctx.state = &ostateCheckSubmission{}
					ctx.log.Debug("day of month task submitted")
				}
			}
		case taskdef.OrderingExact:
			{
				if date.IsInExactDate(ctx.odate, ctx.def.Values()) {
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
func (state ostateConfirm) processState(ctx *TaskExecutionContext) bool {

	if ctx.task.Confirmed() {
		ctx.state = &ostateCheckTime{}
		return true
	}

	return false
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

	ctx.state = &ostateStarting{}

	return fulfilled
}

func (state ostateHold) processState(ctx *TaskExecutionContext) bool {

	return false
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

	if !result.Started {
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
	ctx.task.OrderID()
	msg := events.NewMsg(events.WorkRouteCheckStatusMsg{OrderID: ctx.task.OrderID(), WorkerName: ctx.task.WorkerName()})
	receiver := events.NewWorkLaunchReceiver()

	ctx.dispatcher.PushEvent(receiver, events.RouteWorkCheck, msg)

	result, err := receiver.WaitForResult()

	if err != nil {
		ctx.log.Error("State executing error:", err)
		return false
	}

	ctx.task.AddOutput(result.Output)

	if result.Ended {

		n, g, _ := ctx.task.GetInfo()
		ctx.log.Info("Task ended:", n, " group:", g, " id:", ctx.task.OrderID(), " rc:", result.ReturnCode)

		ctx.task.SetEndTime()
		ctx.state = &ostatePostProcessing{}

		if result.ReturnCode > ctx.maxRc {
			ctx.task.SetState(TaskStateEndedNotOk)
			return true
		}

		ctx.task.SetState(TaskStateEndedOk)
		return true

	}

	return false
}
func (state ostatePostProcessing) processState(ctx *TaskExecutionContext) bool {

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
		//:TODO in future this should be computed from history and from projection
		resolvedOdate := map[date.OdateValue]string{
			date.OdateValueDate: string(ctx.odate),
			date.OdateValuePrev: string(ctx.odate),
			date.OdateValueNext: string(ctx.odate),
			date.OdateValueAny:  string(date.OdateValueAny),
		}

		ticketMsg.Tickets[i].Name = t.Name
		ticketMsg.Tickets[i].Odate = date.Odate(resolvedOdate[t.Odate])
		ticketMsg.Tickets[i].Action = t.Action
		ctx.log.Info("TICKET:", t.Name, " ", date.Odate(resolvedOdate[t.Odate]), " ", t.Action)

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
	if state == TaskStateHold {
		return &ostateHold{}
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
