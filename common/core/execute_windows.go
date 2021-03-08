// +build windows

package core

import (
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

var evlog *eventlog.Log

func init() {

	runnerFunc = runFuncWin
}

type windowsService struct {
	component RunnableComponent
	done      chan<- struct{}
}

func (ws *windowsService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {

	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	go ws.component.Start()
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				{
					changes <- c.CurrentStatus
				}
			case svc.Stop, svc.Shutdown:
				{
					ws.component.Shutdown()
					ws.done <- struct{}{}
					return
				}
			}
		}
	}
}

func runFuncWin(comp RunnableComponent, done chan<- struct{}) error {

	var isService bool
	var err error
	isService, err = svc.IsWindowsService()
	if err != nil {
		return err
	}
	if isService {

		evlog, err = eventlog.Open(comp.ServiceName())
		if err != nil {
			return err
		}
		evlog.Warning(1, "starting as service...")
		err = svc.Run(comp.ServiceName(), &windowsService{component: comp, done: done})

	} else {
		fmt.Println("Not a service, normal startup...")
		err = stdRunFunc(comp, done)
	}

	return err
}
