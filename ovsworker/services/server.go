package services

import (
	"fmt"
	"net"

	"github.com/przebro/overseer/common/cert"
	"github.com/przebro/overseer/common/logger"
	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/ovsworker/config"
	"github.com/przebro/overseer/ovsworker/services/handlers"
	"github.com/przebro/overseer/proto/wservices"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

//OvsWorkerServer - represents worker
type OvsWorkerServer struct {
	conf       config.WorkerConfiguration
	grpcServer *grpc.Server
	log        logger.AppLogger
}

//New - creates a new instance of a OvsWorkerServer
func New(config config.WorkerConfiguration, es wservices.TaskExecutionServiceServer, log logger.AppLogger) *OvsWorkerServer {

	options, err := buildOptions(config, log)

	if err != nil {
		zlog := log.Desugar()
		zlog.Error("build options", zap.String("error", err.Error()))
		return nil
	}

	options = append(options, buildInterceptors(log)...)

	wserver := &OvsWorkerServer{}
	wserver.conf = config
	wserver.grpcServer = grpc.NewServer(options...)
	wserver.log = log

	wservices.RegisterTaskExecutionServiceServer(wserver.grpcServer, es)

	return wserver
}

//Start - starts worker
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

//Shutdown - stops worker execution
func (srv *OvsWorkerServer) Shutdown() error {

	srv.grpcServer.GracefulStop()

	return nil
}

func buildOptions(conf config.WorkerConfiguration, log logger.AppLogger) ([]grpc.ServerOption, error) {

	var options = []grpc.ServerOption{}

	if conf.OverseerCA != "" {
		if err := cert.RegisterCA(conf.OverseerCA); err != nil {
			zlog := log.Desugar()
			zlog.Warn("build options", zap.String("error", err.Error()))
		}

	}

	if conf.SecurityLevel != types.ConnectionSecurityLevelNone {
		if creds, err := cert.BuildServerCredentials(conf.WorkerCert, conf.WorkerKey, conf.WorkerCertPolicy, conf.SecurityLevel); err == nil {
			options = append(options, creds)
		} else {
			return nil, err
		}
	}

	return options, nil
}

func buildInterceptors(log logger.AppLogger) []grpc.ServerOption {

	lhandler := handlers.NewLogHandler(log)

	opt := []grpc.ServerOption{grpc.ChainUnaryInterceptor(lhandler.Log), grpc.ChainStreamInterceptor(lhandler.StreamLog)}

	return opt
}
