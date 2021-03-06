package pool

import (
	"overseer/common/logger"
	"overseer/common/types/date"
	"overseer/overseer/internal/events"
	"time"
)

//DailyExecutor - Executes New Day Procedure
type DailyExecutor struct {
	pool              *ActiveTaskPool
	manager           *ActiveTaskPoolManager
	log               logger.AppLogger
	lastExecutionDate date.Odate
}

//NewDailyExecutor - Creates new DailyExecutor
func NewDailyExecutor(dispatcher events.Dispatcher, manager *ActiveTaskPoolManager, pool *ActiveTaskPool, log logger.AppLogger) *DailyExecutor {

	daily := &DailyExecutor{pool: pool, manager: manager, log: log, lastExecutionDate: date.CurrentOdate()}
	if dispatcher != nil {
		dispatcher.Subscribe(events.RouteTimeOut, daily)
	}

	return daily
}

//CheckDailyProcedure - Check if it is time to start daily procedure
func (exec *DailyExecutor) CheckDailyProcedure(tm time.Time) bool {

	h, m := exec.pool.config.NewDayProc.AsTime()
	y, mth, d := date.CurrentOdate().Ymd()

	odt1 := time.Date(y, time.Month(mth), d, h, m, 0, 0, time.Local)

	return odt1.Before(tm) && date.IsBeforeCurrent(exec.lastExecutionDate, date.CurrentOdate())
}

//DailyProcedure - Cleanups and place new tasks in the Active Pool
func (exec *DailyExecutor) DailyProcedure() (int, int) {

	exec.log.Info("Starting new day procedure")
	exec.lastExecutionDate = date.CurrentOdate()

	deleted := exec.pool.cleanupCompletedTasks()
	ordered := exec.manager.orderNewTasks()

	return deleted, ordered

}

//Process - receive notification from dispatcher
func (exec *DailyExecutor) Process(receiver events.EventReceiver, routename events.RouteName, msg events.DispatchedMessage) {

	switch routename {
	case events.RouteTimeOut:
		{
			exec.log.Debug("task action message, route:", events.RouteTimeOut, "id:", msg.MsgID())

			msgdata, istype := msg.Message().(events.RouteTimeOutMsgFormat)
			if !istype {
				er := events.ErrUnrecognizedMsgFormat
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
			err := events.ErrInvalidRouteName
			exec.log.Debug(err)
			events.ResponseToReceiver(receiver, err)
		}
	}
}
