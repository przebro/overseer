package core

import (
	"os"
	"os/signal"
)

//ComponentStarter - starts a component
type ComponentStarter interface {
	Start() error
}

//ComponentShutdowner - shutdowns a component
type ComponentShutdowner interface {
	Shutdown() error
}

//OverseerComponent - starts and shutdowns a component
type OverseerComponent interface {
	ComponentStarter
	ComponentShutdowner
}

//ComponentQuiescer - pauses and resumes the activity of a component
type ComponentQuiescer interface {
	OverseerComponent
	Quiesce() error
	Resume() error
}

//RunnableComponent - represents a core component
type RunnableComponent interface {
	OverseerComponent
	ServiceName() string
}

//OverseerRunner - wraps runnable component into a runnable unit
type OverseerRunner interface {
	Run() error
}

type runFunc func(RunnableComponent, chan<- struct{}) error

//stdRunFunc - starts
func stdRunFunc(c RunnableComponent, done chan<- struct{}) error {

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	go handleSignal(ch, c, done)

	return c.Start()
}

var runnerFunc runFunc = stdRunFunc

type runner struct {
	fn       runFunc
	runnable RunnableComponent
	wait     chan struct{}
}

//NewRunner - Create a new OverseerRunner that can be started as console program or as a windows service.
func NewRunner(r RunnableComponent) OverseerRunner {

	return &runner{fn: runnerFunc, runnable: r, wait: make(chan struct{})}
}

//Run - runs runnable unit
func (r *runner) Run() error {

	r.fn(r.runnable, r.wait)
	<-r.wait
	return nil
}

//handleSignal - handles kill or interrupt signal and safely turns off all components
func handleSignal(signal <-chan os.Signal, c RunnableComponent, done chan<- struct{}) error {

	<-signal
	err := c.Shutdown()
	done <- struct{}{}
	return err
}
