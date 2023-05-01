package overseer

import (
	"time"

	"github.com/przebro/overseer/overseer/config"
)

type overseerTimer struct {
	receiver  TimeEventReciever
	receivers []TimeEventReciever
	interval  config.IntervalValue
	stop      chan struct{}
}

func newTimer(interval config.IntervalValue) *overseerTimer {
	return &overseerTimer{
		receivers: make([]TimeEventReciever, 0),
		interval:  interval,
		stop:      make(chan struct{}),
	}
}

type TimeEventReciever interface {
	ProcessTimeEvent(t time.Time)
}

func (timer *overseerTimer) Start() error {
	return timer.tickerFunc()
}

func (timer *overseerTimer) Shutdown() error {
	timer.stop <- struct{}{}
	return nil
}

func (timer *overseerTimer) tickerFunc() error {

	t := time.NewTicker(time.Duration(int(timer.interval) * int(time.Second)))
	go func() {
		for {
			select {
			case x := <-t.C:
				for _, receiver := range timer.receivers {
					receiver.ProcessTimeEvent(x)
				}
			case <-timer.stop:
				return
			}
		}
	}()

	return nil
}
func (timer *overseerTimer) AddReceiver(receiver TimeEventReciever) {
	timer.receivers = append(timer.receivers, receiver)
}
