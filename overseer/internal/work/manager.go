package work

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/overseer/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var errWorkerBusy = errors.New("no available workers")

//WorkerManager - Manages actions between

type WorkerWorkManager struct {
	log         zerolog.Logger
	interval    int
	maxAttempts int
	workers     map[string]*workerProxy
	shutdown    chan struct{}
}

func NewWorkerWorkManager(conf config.WorkerManagerConfiguration,
	security config.ServerSecurityConfiguration) *WorkerWorkManager {

	m := &WorkerWorkManager{
		log:         log.With().Str("component", "worker-manager").Logger(),
		interval:    conf.WorkerInterval,
		maxAttempts: conf.WorkerMaxAttempts,
		workers:     map[string]*workerProxy{},
		shutdown:    make(chan struct{}),
	}

	for _, w := range conf.Workers {
		m.log.Info().Str("worker", w.WorkerName).Str("host", w.WorkerHost).Int("port", w.WorkerPort).Msg("Creating service worker")
		sworker := newWorkerProxy(w, security, conf.Timeout, conf.WorkerInterval, &m.log)
		m.workers[w.WorkerName] = sworker
	}

	return m

}

func (m *WorkerWorkManager) Push(ctx context.Context, t types.TaskDescription, vars types.EnvironmentVariableList) (types.WorkerTaskStatus, error) {

	var worker *workerProxy
	var ok bool

	if t.WorkerName() != "" {
		if worker, ok = m.workers[t.WorkerName()]; !ok {
			return types.WorkerTaskStatusFailed, fmt.Errorf("worker %s not found", t.WorkerName())
		}

	} else {
		if workerName := m.selectWorker(); workerName != "" {
			t.SetWorkerName(workerName)
			worker = m.workers[workerName]
		} else {
			return types.WorkerTaskStatusWorkerBusy, errWorkerBusy
		}
	}

	worker.push(t, vars)

	return types.WorkerTaskStatusStarting, nil
}

func (m *WorkerWorkManager) Status(ctx context.Context, t types.WorkDescription) types.TaskExecutionStatus {

	return m.workers[t.WorkerName()].taskStatus(ctx, t)
}

func (m *WorkerWorkManager) selectWorker() string {

	type selectedWorker struct {
		name string
		workerState
	}
	states := []selectedWorker{}
	for name, w := range m.workers {
		ws := w.workerState()
		if ws.connected {
			states = append(states, selectedWorker{name: name, workerState: ws})
		}
	}

	if len(states) == 0 {
		return ""
	}

	sort.Slice(states, func(i, j int) bool {
		return states[i].lastRequest.Before(states[j].lastRequest)
	})

	return states[0].name
}
func (m *WorkerWorkManager) Start() error {

	go func() {
		for _, n := range m.workers {
			n.Run()
		}
		<-m.shutdown
		fmt.Println("shutdown recieved")
	}()
	return nil
}
func (m *WorkerWorkManager) Shutdown() error {
	m.shutdown <- struct{}{}
	return nil

}
