package proc

import (
	"time"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/common/types/date"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// DailyExecutor - Executes New Day Procedure
type DailyExecutor struct {
	pool              PoolConfig
	manager           PoolManager
	log               zerolog.Logger
	lastExecutionDate date.Odate
}

type PoolManager interface {
	OrderNewTasks() int
}

type PoolConfig interface {
	NewDayProc() types.HourMinTime
	CleanupCompletedTasks() int
}

// NewDailyExecutor - Creates new DailyExecutor
func NewDailyExecutor(manager PoolManager, pool PoolConfig) *DailyExecutor {

	daily := &DailyExecutor{pool: pool, manager: manager, log: log.With().Str("component", "procedure").Logger(), lastExecutionDate: date.CurrentOdate()}

	return daily
}

// CheckDailyProcedure - Check if it is time to start daily procedure
func (exec *DailyExecutor) CheckDailyProcedure(tm time.Time) bool {

	h, m := exec.pool.NewDayProc().AsTime()
	y, mth, d := date.CurrentOdate().Ymd()

	odt1 := time.Date(y, time.Month(mth), d, h, m, 0, 0, time.Local)

	return odt1.Before(tm) && date.IsBeforeCurrent(exec.lastExecutionDate, date.CurrentOdate())
}

// DailyProcedure - Cleanups and place new tasks in the Active Pool
func (exec *DailyExecutor) DailyProcedure() (int, int) {

	exec.log.Info().Msg("Starting new day procedure")
	exec.lastExecutionDate = date.CurrentOdate()

	deleted := exec.pool.CleanupCompletedTasks()
	ordered := exec.manager.OrderNewTasks()

	return deleted, ordered

}

func (exec *DailyExecutor) ProcessTimeEvent(t time.Time) {
	exec.log.Debug().Msg("Process time event")
	isInTime := exec.CheckDailyProcedure(t)
	if isInTime {
		exec.DailyProcedure()
	}
}
