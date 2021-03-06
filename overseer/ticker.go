package overseer

import (
	"overseer/common/logger"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"time"
)

type overseerTimer struct{ log logger.AppLogger }

type ovsTimer interface {
	tickerFunc(dispatcher events.Dispatcher, interval config.IntervalValue) error
}

func (timer *overseerTimer) tickerFunc(dispatcher events.Dispatcher, interval config.IntervalValue) error {

	t := time.NewTicker(time.Duration(int(interval) * int(time.Second)))
	go func() {
		for {
			x := <-t.C
			h, m, s := x.Clock()
			y, mth, d := x.Date()
			msgdata := events.RouteTimeOutMsgFormat{Year: y, Month: int(mth), Day: d, Hour: h, Min: m, Sec: s}
			msg := events.NewMsg(msgdata)
			err := dispatcher.PushEvent(nil, events.RouteTimeOut, msg)

			if err != nil {
				timer.log.Info("Unable to Push events:", err)
				continue
			}
		}
	}()

	return nil
}
