package pool

import (
	"errors"
	"goscheduler/common/logger"
	"goscheduler/overseer/internal/date"
	"goscheduler/overseer/internal/events"
	"time"
)

//DailyExecutor - Executes New Day Procedure
type DailyExecutor struct {
	pool    *ActiveTaskPool
	manager *ActiveTaskPoolManager
	log     logger.AppLogger
}

//NewDailyExecutor - Creates new DailyExecutor
func NewDailyExecutor(dispatcher events.Dispatcher, manager *ActiveTaskPoolManager, pool *ActiveTaskPool) *DailyExecutor {

	daily := &DailyExecutor{pool: pool, manager: manager, log: logger.Get()}
	if dispatcher != nil {
		dispatcher.Subscribe(events.RouteTimeOut, daily)
	}

	return daily
}

//CheckDailyProcedure - Cleanups and place new tasks in the Active Pool
func (exec *DailyExecutor) CheckDailyProcedure(tm time.Time) bool {

	h, m := exec.pool.config.NewDayProc.AsTime()
	y, mth, d := date.CurrentOdate().Ymd()

	odt1 := time.Date(y, time.Month(mth), d, h, m, 0, 0, time.Local)

	return odt1.Before(tm) && date.IsBeforeCurrent(exec.pool.currentOdate, date.CurrentOdate())

}

//DailyProcedure - Cleanups and place new tasks in the Active Pool
func (exec *DailyExecutor) DailyProcedure() {

	exec.log.Info("Starting new day procedure")
	exec.pool.currentOdate = date.CurrentOdate()

	exec.pool.cleanupCompletedTasks()
	exec.manager.orderNewTasks()

}

//Process - receive notification from dispatcher
func (exec *DailyExecutor) Process(receiver events.EventReceiver, routename events.RouteName, msg events.DispatchedMessage) {

	switch routename {
	case events.RouteTimeOut:
		{
			exec.log.Debug("task action message, route:", events.RouteTimeOut, "id:", msg.MsgID())

			msgdata, istype := msg.Message().(events.RouteTimeOutMsgFormat)
			if !istype {
				er := errors.New("msg not in format")
				exec.log.Error(er)
				events.ResponseToReceiver(receiver, er)
				break
			}
			exec.log.Debug("Process time event")
			tm := time.Date(msgdata.Year, time.Month(msgdata.Month), msgdata.Day, msgdata.Hour, msgdata.Min, msgdata.Sec, 0, time.Local)
			isInTime := exec.CheckDailyProcedure(tm)
			if isInTime {
				exec.DailyProcedure()
			}
		}
	default:
		{
			err := errors.New("Invalid route name")
			exec.log.Debug(err)
			events.ResponseToReceiver(receiver, err)
		}
	}
}
