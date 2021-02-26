package ovsworker

import (
	"fmt"
	"net"
	"overseer/common/logger"
	"overseer/ovsworker/config"
	"overseer/ovsworker/services"
	"overseer/proto/wservices"
	"path/filepath"

	"google.golang.org/grpc"
)

type ovsWorkerService struct {
	config     *config.Config
	grpcServer *grpc.Server
	listener   net.Listener
	log        logger.AppLogger
}

type OvsWorkerService interface {
	Start() error
}

func NewWorkerService(config *config.Config) OvsWorkerService {

	log := logger.Get()
	wsrvc := &ovsWorkerService{}
	wsrvc.config = config
	wsrvc.grpcServer = grpc.NewServer()

	var err error
	conn := fmt.Sprintf("%s:%d", config.Worker.Host, config.Worker.Port)
	log.Info("Starting worker:", conn)
	wsrvc.listener, err = net.Listen("tcp", conn)
	if err != nil {
		log.Error(err)
		return nil
	}

	var defPath string

	if !filepath.IsAbs(config.Worker.SysoutDirectory) {
		defPath = filepath.Join(config.Worker.RootDirectory, config.Worker.SysoutDirectory)
	} else {
		defPath = config.Worker.SysoutDirectory
	}

	srvc, err := services.NewWorkerExecutionService(defPath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	wservices.RegisterTaskExecutionServiceServer(wsrvc.grpcServer, srvc)
	wsrvc.log = log

	return wsrvc
}

func (wsrvc *ovsWorkerService) Start() error {

	return wsrvc.grpcServer.Serve(wsrvc.listener)

}
