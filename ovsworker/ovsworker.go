package ovsworker

import (
	"fmt"
	"path/filepath"

	"github.com/przebro/overseer/common/core"
	"github.com/przebro/overseer/ovsworker/config"
	"github.com/przebro/overseer/ovsworker/services"
	"github.com/przebro/overseer/proto/wservices"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"
)

type ovsWorker struct {
	conf          config.OverseerWorkerConfiguration
	grpcComponent core.OverseerComponent
	wservices.UnimplementedTaskExecutionServiceServer
}

// NewWorkerService - creates a new Worker
func NewWorkerService(config config.OverseerWorkerConfiguration) (core.RunnableComponent, error) {

	var err error
	var gs *services.OvsWorkerServer

	log := log.With().Str("component", "worker").Logger()

	if gs, err = createServiceServer(config, &log); err != nil {
		return nil, err
	}

	wserver := &ovsWorker{
		grpcComponent: gs,
		conf:          config,
	}

	return wserver, nil
}

func createServiceServer(conf config.OverseerWorkerConfiguration, log *zerolog.Logger) (*services.OvsWorkerServer, error) {

	if !filepath.IsAbs(conf.Worker.SysoutDirectory) {
		return nil, fmt.Errorf("filepath is not absolute:%s", conf.Worker.SysoutDirectory)
	}

	srvc, err := services.NewWorkerExecutionService(conf.Worker.SysoutDirectory, conf.Worker.TaskLimit)
	if err != nil {
		return nil, err
	}

	grpcsrv := services.New(conf.Worker, srvc)
	if grpcsrv == nil {
		return nil, errors.New("unable to initialize grpc server")
	}

	return grpcsrv, nil
}

func (wsrvc *ovsWorker) Start() error {

	return wsrvc.grpcComponent.Start()
}

func (wsrvc *ovsWorker) Shutdown() error {

	return wsrvc.grpcComponent.Shutdown()
}

func (wsrvc *ovsWorker) ServiceName() string {

	return wsrvc.conf.Worker.ServiceName
}
