package ovsworker

import (
	"overseer/common/core"
	"overseer/common/logger"
	"overseer/ovsworker/config"
	"overseer/ovsworker/services"
	"path/filepath"

	"github.com/pkg/errors"
)

type ovsWorker struct {
	conf          config.OverseerWorkerConfiguration
	grpcComponent core.OverseerComponent
	logger        logger.AppLogger
}

//NewWorkerService - creates a new Worker
func NewWorkerService(config config.OverseerWorkerConfiguration) (core.RunnableComponent, error) {

	var err error
	var lg logger.AppLogger
	var gs *services.OvsWorkerServer

	lg = logger.Get()

	if gs, err = createServiceServer(config); err != nil {
		return nil, err
	}

	wserver := &ovsWorker{
		logger:        lg,
		grpcComponent: gs,
		conf:          config,
	}

	return wserver, nil
}

func createServiceServer(conf config.OverseerWorkerConfiguration) (*services.OvsWorkerServer, error) {

	var sysPath string

	if !filepath.IsAbs(conf.Worker.SysoutDirectory) {
		sysPath = filepath.Join(conf.Worker.RootDirectory, conf.Worker.SysoutDirectory)
	} else {
		sysPath = conf.Worker.SysoutDirectory
	}

	srvc, err := services.NewWorkerExecutionService(sysPath)
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
