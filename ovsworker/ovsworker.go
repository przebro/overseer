package ovsworker

import (
	"fmt"
	"goscheduler/common/logger"
	"goscheduler/ovsworker/config"
	"goscheduler/ovsworker/launcher"
	"goscheduler/ovsworker/services"
	"goscheduler/proto/wservices"
	"net"

	"google.golang.org/grpc"
)

type ovsWorkerService struct {
	config     *config.Config
	grpcServer *grpc.Server
	listener   net.Listener
	workExec   *launcher.FragmentLauncher
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

	srvc := services.NewWorkerExecutionService()
	wservices.RegisterTaskExecutionServiceServer(wsrvc.grpcServer, srvc)
	wsrvc.workExec = launcher.NewFragmentLauncher()
	wsrvc.log = log

	return wsrvc
}

func (wsrvc *ovsWorkerService) Start() error {

	return wsrvc.grpcServer.Serve(wsrvc.listener)

}
