package work

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/przebro/overseer/common/cert"
	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/overseer/config"
	converter "github.com/przebro/overseer/overseer/internal/work/converters"
	"github.com/przebro/overseer/proto/wservices"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

type workerProxy struct {
	log       *zerolog.Logger
	client    wservices.TaskExecutionServiceClient
	status    map[string]types.TaskExecutionStatus
	config    config.WorkerConfiguration
	security  config.ServerSecurityConfiguration
	lock      sync.Mutex
	state     workerState
	timeout   int
	interval  int
	shutdown  chan struct{}
	cleanchnl chan string
}

func newWorkerProxy(conf config.WorkerConfiguration,
	security config.ServerSecurityConfiguration,
	timeout int,
	interval int,
	log *zerolog.Logger) *workerProxy {

	worker := &workerProxy{
		config:    conf,
		log:       log,
		status:    map[string]types.TaskExecutionStatus{},
		security:  security,
		lock:      sync.Mutex{},
		state:     workerState{},
		timeout:   timeout,
		interval:  interval,
		shutdown:  make(chan struct{}),
		cleanchnl: make(chan string),
	}

	if client := worker.connect(conf, security.SecurityLevel, security.ClientCertPolicy, timeout); client == nil {
		worker.state.connected = false
	} else {
		worker.client = client
		worker.state.connected = true
	}

	return worker
}

func (w *workerProxy) connect(conf config.WorkerConfiguration, level types.ConnectionSecurityLevel, policy types.CertPolicy,
	timeout int) wservices.TaskExecutionServiceClient {

	opt := []grpc.DialOption{
		grpc.WithBlock(),
	}

	if level == types.ConnectionSecurityLevelNone {
		opt = append(opt, grpc.WithInsecure())
	} else {

		if err := cert.RegisterCA(conf.WorkerCA); err != nil {
			w.log.Warn().Err(err).Msg("connect")
		}

		result, err := cert.BuildClientCredentials(w.security.ServerCert, w.security.ServerKey, policy, level)
		if err != nil {
			w.log.Warn().Err(err).Msg("connect")
			return nil
		}

		opt = append(opt, result)
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))

	targetAddr := fmt.Sprintf("%s:%d", conf.WorkerHost, conf.WorkerPort)
	conn, err := grpc.DialContext(ctx, targetAddr, opt...)
	if err != nil {
		w.log.Error().Err(err).Msg("dial")
		w.lock.Lock()
		w.state.connected = false
		w.lock.Unlock()
		return nil
	}

	return wservices.NewTaskExecutionServiceClient(conn)
}

func (w *workerProxy) push(t TaskDescription, variables types.EnvironmentVariableList) {

	defer w.lock.Unlock()

	status := types.TaskExecutionStatus{
		Status:      types.WorkerTaskStatusStarting,
		OrderID:     t.OrderID(),
		WorkerName:  w.config.WorkerName,
		ExecutionID: t.ExecutionID(),
	}
	w.lock.Lock()
	w.status[t.ExecutionID()] = status

	go func() {

		status := types.TaskExecutionStatus{
			OrderID:     t.OrderID(),
			WorkerName:  w.config.WorkerName,
			ExecutionID: t.ExecutionID(),
		}

		smsg := &wservices.StartTaskMsg{}
		smsg.TaskID = &wservices.TaskIdMsg{TaskID: string(t.OrderID()), ExecutionID: t.ExecutionID()}
		smsg.Type = string(t.TypeName())

		smsg.Variables = map[string]string{}
		for _, n := range variables {
			smsg.Variables[n.Expand()] = n.Value

		}
		var err error

		smsg.Command, err = converter.ConvertToMsg(t.TypeName(), t.Action(), variables)
		if err != nil {
			w.log.Error().Err(err).Msg("start_task")
		}

		resp, err := w.client.StartTask(context.Background(), smsg)
		if err != nil {
			w.log.Error().Err(err).Msg("start_task")
			status.Status = types.WorkerTaskStatusFailed
		} else {
			status.ReturnCode = resp.ReturnCode
			status.StatusCode = resp.StatusCode
			status.Status = reverseStatusMap[resp.Status]
		}

		w.lock.Lock()
		defer w.lock.Unlock()
		w.status[status.ExecutionID] = status

	}()
}

func (w *workerProxy) updateStatus() {
	go func(worker *workerProxy) {
		var client wservices.TaskExecutionServiceClient

		w.lock.Lock()
		client = worker.client
		w.lock.Unlock()

		if client == nil {
			if client := w.connect(w.config, w.security.SecurityLevel, w.security.ClientCertPolicy, worker.timeout); client == nil {
				w.log.Error().Msg("worker not available")
				return
			}

			w.lock.Lock()
			worker.client = client
			w.lock.Unlock()
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer w.lock.Unlock()
		defer cancel()
		result, err := w.client.WorkerStatus(ctx, &empty.Empty{})

		if err != nil {
			w.log.Error().Str("worker", w.config.WorkerName).Err(err).Msg("worker not available")
			w.lock.Lock()
			worker.state.connected = false
		} else {
			w.lock.Lock()
			fmt.Println("status updated")
			worker.state = workerState{
				connected:   true,
				cpu:         int(result.Cpuload),
				memused:     int(result.Memused),
				memtotal:    int(result.Memtotal),
				tasks:       int(result.Tasks),
				tasksLimit:  int(result.TasksLimit),
				lastRequest: time.Now(),
			}
		}

	}(w)

}

func (w *workerProxy) workerState() workerState {

	defer w.lock.Unlock()
	w.lock.Lock()
	return w.state

}

func (w *workerProxy) taskStatus(ctx context.Context, t WorkDescription) types.TaskExecutionStatus {

	defer w.lock.Unlock()
	w.lock.Lock()

	taskStatus := w.status[t.ExecutionID()]

	go func() {
		result := types.TaskExecutionStatus{OrderID: t.OrderID(), ExecutionID: t.ExecutionID(), WorkerName: t.WorkerName()}
		resp, err := w.client.TaskStatus(ctx, &wservices.TaskIdMsg{TaskID: string(t.OrderID()), ExecutionID: t.ExecutionID()})
		if err != nil {
			if s, ok := status.FromError(err); ok {
				//something really bad happen with worker and task is lost
				if s.Code() == codes.NotFound {
					result.Status = types.WorkerTaskStatusFailed
				}
			}
			w.log.Error().Err(err).Msg("task_status")
			return
		}

		result.Status = reverseStatusMap[resp.Status]
		result.ReturnCode = resp.ReturnCode
		result.StatusCode = resp.StatusCode

		defer w.lock.Unlock()
		w.lock.Lock()

		taskStatus, ok := w.status[t.ExecutionID()]

		if ok {
			if (result.Status == types.WorkerTaskStatusEnded || result.Status == types.WorkerTaskStatusFailed) &&
				result.Status == taskStatus.Status {
				w.cleanchnl <- t.ExecutionID()
				return
			}
		}

		w.status[t.ExecutionID()] = result

	}()

	return taskStatus
}

func (w *workerProxy) cleanup(executionID string) {
	defer w.lock.Unlock()
	w.lock.Lock()
	taskStatus, ok := w.status[executionID]
	if !ok {
		return
	}

	w.client.CompleteTask(context.Background(), &wservices.TaskIdMsg{TaskID: string(taskStatus.OrderID), ExecutionID: executionID})
}

func (w *workerProxy) Run() {

	go func() {

		ticker := time.NewTicker(time.Duration(w.interval) * time.Second)
		for {
			select {
			case <-ticker.C:
				{
					w.log.Info().Msg("request for update status")
					w.updateStatus()
				}
			case id := <-w.cleanchnl:
				w.cleanup(id)
			case <-w.shutdown:
				ticker.Stop()
				//what if shutdown called and there are task to clean ???
				return

			}
		}
	}()
}

func (w *workerProxy) Shutdown() {
	w.shutdown <- struct{}{}
}
