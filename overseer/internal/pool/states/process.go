package states

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/przebro/expr"
	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/overseer/internal/journal"
	"github.com/przebro/overseer/overseer/internal/pool/activetask"
	"github.com/przebro/overseer/overseer/internal/pool/calc"
	"github.com/przebro/overseer/overseer/internal/pool/models"
	"github.com/przebro/overseer/overseer/internal/pool/readers"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/rs/zerolog"
)

type ResourceManager interface {
	ProcessReleaseFlag(input []string) (bool, []string)
	ProcessAcquireFlag(input []types.FlagModel) (bool, []string)
	ProcessTicketAction(tickets []types.TicketActionModel) bool
	CheckTickets(in []types.CollectedTicketModel) []types.CollectedTicketModel
}

type WorkManager interface {
	Push(ctx context.Context, t types.TaskDescription, vars types.EnvironmentVariableList) (types.WorkerTaskStatus, error)
	Status(ctx context.Context, t types.WorkDescription) types.TaskExecutionStatus
}

// TaskExecutionContext - context for processing active task
type TaskExecutionContext struct {
	Task       *activetask.TaskInstance
	Odate      date.Odate
	Time       time.Time
	State      TaskExecutionState
	Log        zerolog.Logger
	IsInTime   bool
	IsEnforced bool
	MaxRc      int32
	Rmanager   ResourceManager
	Wmanager   WorkManager
	Journal    readers.JournalWriter
}

type TaskExecutionState interface {
	ProcessState(order *TaskExecutionContext) bool
}

// States for task in Active pool
// OstateConfirm - Task is ordered but waits for user's confirmation
type OstateConfirm struct{}

// OstateCheckTime -  Task is confirmed but waits for the time window
type OstateCheckTime struct{}

// OstateCheckCyclic -  Task is confirmed if it is a cyclic task check
type OstateCheckCyclic struct{}

// OstateCheckConditions - Task is already in time window but waits for conditions
type OstateCheckConditions struct{}

// OstateAcquireResources - Task has all tickets but waits for a flag or an idle worker
type OstateAcquireResources struct{}

// ostateStarting - Task has all resources and can run now
type OstateStarting struct{}

// OstateExecuting - Task is executed
type OstateExecuting struct{}

// ostatePostProcessing - Task ended, if ok, release resources and manage conditions
type OstatePostProcessing struct{}

func (state OstateConfirm) ProcessState(ctx *TaskExecutionContext) bool {

	if !ctx.Task.Confirmed() {
		ctx.Log.Debug().Str("state", "confirm").Msg("Task is not confirmed")
		return false
	}

	ctx.State = &OstateCheckCyclic{}

	return true
}
func (state OstateCheckCyclic) ProcessState(ctx *TaskExecutionContext) bool {

	if ctx.IsEnforced {
		ctx.Log.Info().Str("state", "cyclic").Msg("Task is enforced")
		ctx.State = &OstateCheckTime{}
		return true
	}

	if !ctx.Task.IsCyclic() {
		ctx.State = &OstateCheckTime{}
		return true
	}

	next := ctx.Task.NextRun()
	h, m := next.AsTime()

	if !ctx.Time.Before(time.Date(ctx.Time.Year(), ctx.Time.Month(), ctx.Time.Day(), h, m, 0, 0, ctx.Time.Location())) {
		ctx.State = &OstateCheckTime{}
		return true
	}

	return false
}
func (state OstateCheckTime) ProcessState(ctx *TaskExecutionContext) bool {

	var isInTime bool = false

	if ctx.IsEnforced {
		ctx.Log.Info().Str("state", "check_time").Msg("Task is enforced")
		ctx.IsInTime = true
		ctx.State = &OstateCheckConditions{}
		return true
	}

	from, to := ctx.Task.TimeSpan()
	if from == "" && to == "" {
		ctx.IsInTime = true
		ctx.State = &OstateCheckConditions{}
		return true
	}
	if from != "" {
		h, m := from.AsTime()
		//if current time is after from time task is in time window
		if !ctx.Time.Before(time.Date(ctx.Time.Year(), ctx.Time.Month(), ctx.Time.Day(), h, m, 0, 0, ctx.Time.Location())) {
			isInTime = true
		}
	}
	if to != "" {
		if from == "" {
			isInTime = true
		}
		h, m := to.AsTime()
		// if current time is after to time, task must wait for time window
		if !ctx.Time.Before(time.Date(ctx.Time.Year(), ctx.Time.Month(), ctx.Time.Day(), h, m, 0, 0, ctx.Time.Location())) {
			isInTime = false
		}
	}

	ctx.IsInTime = isInTime
	ctx.State = &OstateCheckConditions{}

	return true
}

func (state OstateCheckConditions) ProcessState(ctx *TaskExecutionContext) bool {

	if ctx.IsEnforced {
		ctx.Log.Info().Str("state", "check_coditions").Msg("Task is enforced")
		ctx.State = &OstateStarting{}
		return true
	}

	ctx.Task.SetState(models.TaskStateWaiting)

	var fulfilled bool = false
	var err error
	var rodate date.Odate
	payload := []types.CollectedTicketModel{}

	ex := ctx.Task.Expression

	if ex == "AND" || ex == "OR" {

		for _, tc := range ctx.Task.Tickets(ctx.Odate) {

			payload = append(payload, types.CollectedTicketModel{Name: tc.Name, Odate: date.Odate(tc.Odate)})
		}

		result := ctx.Rmanager.CheckTickets(payload)

		if len(result) == 0 {
			fulfilled = true
		} else {

			if ctx.Task.Expression == string(taskdef.InTicketAND) {
				fulfilled = true
			}

			for _, t := range result {
				if taskdef.InTicketRelation(ctx.Task.Expression) == taskdef.InTicketAND {
					fulfilled = t.Exists && fulfilled
				} else {
					fulfilled = t.Exists || fulfilled
				}
			}
		}
	} else {
		// If relation is described by expression not by simple OR / AND

		vars := map[string]interface{}{}
		var values []string

		if values, err = expr.Extract(ex); err != nil {
			n, g := ctx.Task.Definition.Name, ctx.Task.Definition.Group
			ctx.Log.Error().Str("group", g).Str("name", n).Err(err).Msg("expr")
			return false
		}

		for _, item := range values {

			tdata := strings.SplitN(item, ".", 2)

			if len(tdata) == 2 {
				rodate = calc.CalcRealOdate(ctx.Odate, date.OdateValue(tdata[1]), ctx.Task.Schedule)
			} else {
				rodate = date.Odate(date.OdateValueNone)
			}

			payload = append(payload, types.CollectedTicketModel{Name: tdata[0], Odate: date.Odate(rodate)})
		}

		for i, n := range ctx.Rmanager.CheckTickets(payload) {
			vars[values[i]] = n.Exists
		}

		if fulfilled, err = expr.Eval(ex, vars); err != nil {
			n, g := ctx.Task.Definition.Name, ctx.Task.Definition.Group
			ctx.Log.Error().Str("group", g).Str("name", n).Str("expr", ex).Err(err).Msg("invalid expression")
			return false
		}
	}

	if fulfilled && ctx.IsInTime {
		ctx.State = &OstateAcquireResources{}
		ctx.Journal.PushJournalMessage(ctx.Task.OrderID(), ctx.Task.ExecutionID(), time.Now(), journal.TaskFulfill)
	}

	return fulfilled && ctx.IsInTime
}

func (state OstateAcquireResources) ProcessState(ctx *TaskExecutionContext) bool {

	flags := ctx.Task.Flags

	if len(flags) == 0 {
		ctx.State = &OstateStarting{}
		return true
	}

	input := []types.FlagModel{}
	for _, f := range flags {
		policy := uint8(0)

		if f.Type == taskdef.FlagExclusive {
			policy = 1
		}
		input = append(input, types.FlagModel{Name: f.Name, Policy: policy})
	}

	result, _ := ctx.Rmanager.ProcessAcquireFlag(input)

	if !result {
		ctx.Task.SetState(models.TaskStateWaiting)
		return false
	}

	ctx.State = &OstateStarting{}
	return true
}

func (state OstateStarting) ProcessState(ctx *TaskExecutionContext) bool {

	ctx.Task.SetState(models.TaskStateStarting)
	variables := calc.PrepareVaribles(ctx.Task, ctx.Odate)

	result, err := ctx.Wmanager.Push(context.Background(), ctx.Task, variables)

	if err != nil {

		ctx.Journal.PushJournalMessage(ctx.Task.OrderID(), ctx.Task.ExecutionID(), time.Now(), journal.TaskStartingFailedErr)

		ctx.State = &OstatePostProcessing{}
		ctx.Task.SetState(models.TaskStateEndedNotOk)
		ctx.Log.Error().Err(err).Msg("state error")
		return true
	}

	if result == types.WorkerTaskStatusWorkerBusy {

		ctx.Task.SetState(models.TaskStateWaiting)
		flags := ctx.Task.Flags

		input := []string{}

		for _, f := range flags {
			input = append(input, f.Name)
		}

		ctx.Rmanager.ProcessReleaseFlag(input)

		return false

	}

	stime := ctx.Task.SetStartTime()

	ctx.Journal.PushJournalMessage(ctx.Task.OrderID(), ctx.Task.ExecutionID(), stime, fmt.Sprintf(journal.TaskStartingRN, ctx.Task.RunNumber()))

	if result != types.WorkerTaskStatusStarting {

		ctx.Journal.PushJournalMessage(ctx.Task.OrderID(), ctx.Task.ExecutionID(), time.Now(), fmt.Sprintf(journal.TaskStartingFailed, result))

		ctx.State = &OstatePostProcessing{}
		ctx.Task.SetState(models.TaskStateEndedNotOk)
		ctx.Log.Error().Err(err).Msg("task not executed")

		return true
	}

	ctx.Journal.PushJournalMessage(ctx.Task.OrderID(), ctx.Task.ExecutionID(), time.Now(), fmt.Sprintf(journal.TaskStarting, ctx.Task.WorkerName()))

	ctx.Task.SetState(models.TaskStateExecuting)

	ctx.State = &OstateExecuting{}

	return true
}

func (state OstateExecuting) ProcessState(ctx *TaskExecutionContext) bool {

	ctx.Task.SetState(models.TaskStateExecuting)

	result := ctx.Wmanager.Status(context.Background(), ctx.Task)
	fmt.Println(result, "::", result.Status, result.StatusCode, result.ReturnCode)

	if result.Status == types.WorkerTaskStatusEnded || result.Status == types.WorkerTaskStatusFailed {

		n, g, _ := ctx.Task.GetInfo()
		ctx.Log.Info().Str("group", g).Str("name", n).
			Int32("rc", result.ReturnCode).
			Int32("sc", result.StatusCode).
			Str("id", string(ctx.Task.OrderID())).Msg("task ended")

		tm := time.Now()

		ctx.State = &OstatePostProcessing{}

		ctx.Journal.PushJournalMessage(ctx.Task.OrderID(), ctx.Task.ExecutionID(), tm, fmt.Sprintf(journal.TaskComplete, tm.Format("2006-01-02 15:04:05.000000")))

		if result.Status == types.WorkerTaskStatusFailed {

			ctx.Journal.PushJournalMessage(ctx.Task.OrderID(), ctx.Task.ExecutionID(), tm, journal.TaskFailed)
			ctx.Task.SetState(models.TaskStateEndedNotOk)
		}

		if result.Status == types.WorkerTaskStatusEnded {

			msg := ""

			resultState := calc.ComputeTaskState(ctx.Task.Definition.Type, ctx.MaxRc, result.ReturnCode, result.StatusCode)

			if resultState == models.TaskStateEndedNotOk {
				ctx.Task.SetState(models.TaskStateEndedNotOk)
				msg = fmt.Sprintf(journal.TaskEndedNOK, result.ReturnCode, result.StatusCode)

			} else {
				ctx.Task.SetState(models.TaskStateEndedOk)
				msg = fmt.Sprintf(journal.TaskEndedOK, result.ReturnCode, result.StatusCode)
			}

			ctx.Journal.PushJournalMessage(ctx.Task.OrderID(), ctx.Task.ExecutionID(), time.Now(), msg)
		}

		return true
	}

	return false
}
func (state OstatePostProcessing) ProcessState(ctx *TaskExecutionContext) bool {

	ctx.Task.SetEndTime()

	flags := ctx.Task.Flags

	input := []string{}

	for _, f := range flags {
		input = append(input, f.Name)
	}

	ctx.Rmanager.ProcessReleaseFlag(input)

	if ctx.Task.State() == models.TaskStateEndedNotOk {

		ctx.Journal.PushJournalMessage(ctx.Task.OrderID(), ctx.Task.ExecutionID(), time.Now(), journal.TaskPostProc)
		ctx.Log.Info().Msg("Task post processing ends")

		if !ctx.Task.IsCyclic() {
			return false
		}

	} else {

		outticket := ctx.Task.OutTickets
		input := []types.TicketActionModel{}
		for _, n := range outticket {

			realOdate := calc.CalcRealOdate(ctx.Task.OrderDate(), n.Odate, ctx.Task.Schedule)
			if n.Action == taskdef.OutActionAdd {
				input = append(input, types.TicketActionModel{Name: n.Name, Action: "ADD", Odate: realOdate})

			} else {
				input = append(input, types.TicketActionModel{Name: n.Name, Action: "REM", Odate: realOdate})
			}
		}

		ctx.Rmanager.ProcessTicketAction(input)

		ctx.Journal.PushJournalMessage(ctx.Task.OrderID(), ctx.Task.ExecutionID(), time.Now(), journal.TaskPostProc)
		ctx.Log.Info().Msg("Task post processing ends")
	}

	if ctx.Task.PrepareNextCycle() {
		ctx.Task.SetExecutionID()
	}

	return false
}
