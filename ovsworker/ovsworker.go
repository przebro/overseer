package ovsworker

import (
	"fmt"
	"path/filepath"

	"github.com/przebro/overseer/common/core"
	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/ovsworker/config"
	"github.com/przebro/overseer/ovsworker/services"
	"github.com/przebro/overseer/proto/wservices"

	"github.com/pkg/errors"
)

type ovsWorker struct {
	conf          config.OverseerWorkerConfiguration
	grpcComponent core.OverseerComponent
	logger        logger.AppLogger
	wservices.UnimplementedTaskExecutionServiceServer
}

//NewWorkerService - creates a new Worker
func NewWorkerService(config config.OverseerWorkerConfiguration, log logger.AppLogger) (core.RunnableComponent, error) {

	var err error
	var gs *services.OvsWorkerServer

	if gs, err = createServiceServer(config, log); err != nil {
		return nil, err
	}

	wserver := &ovsWorker{
		logger:        log,
		grpcComponent: gs,
		conf:          config,
	}

	return wserver, nil
}

func createServiceServer(conf config.OverseerWorkerConfiguration, log logger.AppLogger) (*services.OvsWorkerServer, error) {

	if !filepath.IsAbs(conf.Worker.SysoutDirectory) {
		return nil, fmt.Errorf("filepath is not absolute:%s", conf.Worker.SysoutDirectory)
	}

	srvc, err := services.NewWorkerExecutionService(conf.Worker.SysoutDirectory, conf.Worker.TaskLimit, log)
	if err != nil {
		return nil, err
	}

	grpcsrv := services.New(conf.Worker, srvc, log)
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
