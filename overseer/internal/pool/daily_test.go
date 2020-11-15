package pool

import (
	"goscheduler/common/logger"
	"goscheduler/overseer/config"
	"goscheduler/overseer/internal/events"
	"sync"
	"testing"
	"time"
)

type dDispatcher struct {
	Tickets         map[string]string
	processNotEnded bool
	withError       bool
}

func (m *dDispatcher) PushEvent(receiver events.EventReceiver, route events.RouteName, msg events.DispatchedMessage) error {
	return nil
}

func (m *dDispatcher) Subscribe(route events.RouteName, participant events.EventParticipant) {

}
func (m *dDispatcher) Unsubscribe(route events.RouteName, participant events.EventParticipant) {

}

var ddsp events.Dispatcher = &dDispatcher{}
var dpool *ActiveTaskPool = NewTaskPool(dsp, config.ActivePoolConfiguration{ForceNewDayProc: true, MaxOkReturnCode: 4, NewDayProc: "!0:30"})

var pManager *ActiveTaskPoolManager = &ActiveTaskPoolManager{}
var daily *DailyExecutor = NewDailyExecutor(ddsp, pManager, dpool)

func TestDailyExecutor(t *testing.T) {

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
