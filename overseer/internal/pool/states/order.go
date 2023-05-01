package states

import (
	"github.com/przebro/overseer/common/types/date"
	"github.com/przebro/overseer/overseer/internal/taskdef"
	"github.com/rs/zerolog"
)

// TaskOrderContext - use state pattern to verify if task can be added to Active Task Pool (ATP)
type TaskOrderContext struct {
	IgnoreCalendar   bool
	IgnoreSubmission bool
	Odate            date.Odate
	CurrentOdate     date.Odate
	Def              *taskdef.TaskDefinition
	State            TaskOrderState
	IsSubmited       bool
	Reason           []string
	Log              zerolog.Logger
}

type TaskOrderState interface {
	ProcessState(order *TaskOrderContext) bool
}

type OstateCheckOtype struct{}
type OstateCheckCalendar struct{}
type OstateOrdered struct{}
type OstateNotSubmitted struct{}

type taskOrderState interface {
	processState(order *TaskOrderContext) bool
}

func (state OstateCheckOtype) ProcessState(ctx *TaskOrderContext) bool {

	ctx.State = &OstateCheckCalendar{}
	return true
}

func (state OstateCheckCalendar) ProcessState(ctx *TaskOrderContext) bool {

	canceledState := &OstateNotSubmitted{}
	ctx.State = canceledState

	if ctx.IgnoreCalendar {
		ctx.Log.Debug().Msg("state check calendar processed")
		ctx.State = &OstateOrdered{}
	} else {
		switch ctx.Def.Schedule.OrderType {
		case taskdef.OrderingWeek:
			{
				ctx.Log.Debug().Msg("checking task weekly")
				if date.IsInDayOfWeek(ctx.Odate, ctx.Def.Schedule.Dayvalues) && date.IsInMonth(ctx.Odate, ctx.Def.Schedule.Months) {
					ctx.State = &OstateOrdered{}
					ctx.Log.Debug().Msg("weekly task submitted")
				}
			}
		case taskdef.OrderingDaily:
			{
				ctx.Log.Debug().Msg("checking task daily")
				if date.IsInMonth(ctx.Odate, ctx.Def.Schedule.Months) {
					ctx.State = &OstateOrdered{}
					ctx.Log.Debug().Msg("daily task submitted")
				}

			}
		case taskdef.OrderingDayOfMonth:
			{
				ctx.Log.Debug().Msg("checking task day of month")
				if date.IsInDayOfMonth(ctx.Odate, ctx.Def.Schedule.Dayvalues) && date.IsInMonth(ctx.Odate, ctx.Def.Schedule.Months) {
					ctx.State = &OstateOrdered{}
					ctx.Log.Debug().Msg("day of month task submitted")
				}
			}
		case taskdef.OrderingFromEnd:
			{
				if date.IsInFromEnd(ctx.Odate, ctx.Def.Schedule.Dayvalues) && date.IsInMonth(ctx.Odate, ctx.Def.Schedule.Months) {
					ctx.State = &OstateOrdered{}
					ctx.Log.Debug().Msg("from end of month task submitted")
				}
			}
		case taskdef.OrderingExact:
			{
				if date.IsInExactDate(ctx.Odate, ctx.Def.Schedule.Exactdates) {
					ctx.State = &OstateOrdered{}
				}

			}
		case taskdef.OrderingManual:
			{
				ctx.State = &OstateOrdered{}
			}
		}

		if ctx.State == canceledState {
			ctx.Reason = append(ctx.Reason, "Scheduling criteria does not meet")
		}
	}

	return true
}
func (state OstateOrdered) ProcessState(ctx *TaskOrderContext) bool {
	ctx.Log.Debug().Msg("state ordered")
	ctx.IsSubmited = true
	return false
}
func (state OstateNotSubmitted) ProcessState(ctx *TaskOrderContext) bool {

	ctx.Reason = append(ctx.Reason, "Task not submitted")
	ctx.IsSubmited = false
	return false
}
