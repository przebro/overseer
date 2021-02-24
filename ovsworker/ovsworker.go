package ovsworker

import (
	"fmt"
	"net"
	"overseer/common/logger"
	"overseer/ovsworker/config"
	"overseer/ovsworker/services"
	"overseer/proto/wservices"

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

	srvc, err := services.NewWorkerExecutionService(config.Worker.SysoutDirectory)
	if err != nil {
		log.Error(err)
		return nil
	}
	wservices.RegisterTaskExecutionServiceServer(wsrvc.grpcServer, srvc)
	wsrvc.log = log

	return wsrvc
}

func (wsrvc *ovsWorkerService) Start() error {

	return wsrvc.grpcServer.Serve(wsrvc.listener)

}
