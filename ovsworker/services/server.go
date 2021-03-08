package services

import (
	"fmt"
	"net"
	"overseer/common/logger"
	"overseer/ovsworker/config"
	"overseer/proto/wservices"

	"google.golang.org/grpc"
)

type OvsWorkerServer struct {
	conf       config.WorkerConfiguration
	grpcServer *grpc.Server
	log        logger.AppLogger
}

func New(config config.WorkerConfiguration, es wservices.TaskExecutionServiceServer) *OvsWorkerServer {

	wserver := &OvsWorkerServer{}
	wserver.conf = config
	wserver.grpcServer = grpc.NewServer()
	wserver.log = logger.Get()

	wservices.RegisterTaskExecutionServiceServer(wserver.grpcServer, es)

	return wserver
}

func (srv *OvsWorkerServer) Start() error {

	conn := fmt.Sprintf("%s:%d", srv.conf.Host, srv.conf.Port)
	l, err := net.Listen("tcp", conn)
	if err != nil {
		srv.log.Error(err)
		return err
	}

	srv.log.Info("starting grpc server:", conn)

	err = srv.grpcServer.Serve(l)
	if err != nil {
		return err
	}

	return nil
}

func (srv *OvsWorkerServer) Shutdown() error {

	srv.grpcServer.GracefulStop()

	return nil
}
