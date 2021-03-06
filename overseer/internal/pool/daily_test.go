package pool

import (
	"fmt"
	"overseer/common/logger"
	"overseer/common/types"
	"overseer/common/types/date"
	"overseer/datastore"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/taskdef"
	"sync"
	"testing"
	"time"
)

var pManager *ActiveTaskPoolManager = &ActiveTaskPoolManager{}
var daily *DailyExecutor

func init() {

	if taskPoolT == nil {
		provider, _ = datastore.NewDataProvider(storeConfig, log)
		initTaskPool()
	}
	daily = NewDailyExecutor(mDispatcher, pManager, taskPoolT, log)
}

func TestDailyExecutor(t *testing.T) {

	daily = NewDailyExecutor(mDispatcher, pManager, taskPoolT, log)
	if daily == nil {
		t.Error("Daile executor not initialized")
	}

}

func TestProcessDaily(t *testing.T) {
	var rcverr error
	daily.log = logger.NewTestLogger()
	rcv := events.NewTicketCheckReceiver()

	msg := events.NewMsg("test data")
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		_, rcverr = rcv.WaitForResult()
		wg.Done()
	}()

	daily.Process(rcv, events.RouteTicketIn, msg)
	wg.Wait()
	if rcverr == nil {
		t.Error("Unexpected result")
	}

	x := time.Now()
	h, m, s := x.Clock()
	y, mth, d := x.Date()

	fmsg := events.NewMsg(events.RouteTaskStatusResponseMsg{})
	daily.Process(nil, events.RouteTimeOut, fmsg)

	tmsg := events.NewMsg(events.RouteTimeOutMsgFormat{Year: y, Month: int(mth), Day: d, Hour: h, Min: m, Sec: s})
	daily.Process(nil, events.RouteTimeOut, tmsg)
}

func TestCheckDailyProcedure(t *testing.T) {

	daily.log = logger.NewTestLogger()

	tm := time.Now()
	h, m, _ := tm.Clock()

	result := daily.CheckDailyProcedure(tm)
	if result != false {
		t.Error("Unexpected value:", result)
	}

	taskPoolT.config.NewDayProc = types.HourMinTime(fmt.Sprintf("%2d:%2d", h, m-2))
	daily.lastExecutionDate = date.AddDays(daily.lastExecutionDate, -1)

	result = daily.CheckDailyProcedure(tm)
	if result != true {
		t.Error("Unexpected value:", result)
	}

	taskPoolT.config.NewDayProc = types.HourMinTime(fmt.Sprintf("%2d:%2d", h, m+2))
	daily.lastExecutionDate = date.CurrentOdate()

	result = daily.CheckDailyProcedure(tm)
	if result == true {
		t.Error("Unexpected value:", result)
	}
}

func TestDailyProc(t *testing.T) {

	daily.log = logger.NewTestLogger()
	taskPoolT.log = daily.log
	daily.manager.tdm, _ = taskdef.NewManager("../../def/test", log)
	pManager.log = daily.log
	del, ord := daily.DailyProcedure()
	if del != 0 || ord != 0 {
		t.Error("Unexpected values")
	}

}
