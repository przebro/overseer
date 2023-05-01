package services

import (
	"fmt"
	"net"

	"github.com/przebro/overseer/common/cert"
	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/ovsworker/config"
	"github.com/przebro/overseer/ovsworker/services/handlers"
	"github.com/przebro/overseer/proto/wservices"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

// OvsWorkerServer - represents worker
type OvsWorkerServer struct {
	conf       config.WorkerConfiguration
	grpcServer *grpc.Server
}

// New - creates a new instance of a OvsWorkerServer
func New(config config.WorkerConfiguration, es wservices.TaskExecutionServiceServer) *OvsWorkerServer {

	lg := log.With().Str("component", "worker").Logger()

	options, err := buildOptions(config, &lg)

	if err != nil {
		lg.Error().Err(err).Msg("build options")
		return nil
	}

	options = append(options, buildInterceptors(&lg)...)

	wserver := &OvsWorkerServer{}
	wserver.conf = config
	wserver.grpcServer = grpc.NewServer(options...)

	wservices.RegisterTaskExecutionServiceServer(wserver.grpcServer, es)

	return wserver
}

// Start - starts worker
func (srv *OvsWorkerServer) Start() error {

	lg := log.With().Str("component", "worker").Logger()
	conn := fmt.Sprintf("%s:%d", srv.conf.Host, srv.conf.Port)

	l, err := net.Listen("tcp", conn)
	if err != nil {
		lg.Error().Err(err).Msg("starting grpc server")
		return err
	}
	lg.Info().Msg("starting grpc server:" + conn)

	err = srv.grpcServer.Serve(l)
	if err != nil {
		return err
	}

	return nil
}

// Shutdown - stops worker execution
func (srv *OvsWorkerServer) Shutdown() error {

	srv.grpcServer.GracefulStop()

	return nil
}

func buildOptions(conf config.WorkerConfiguration, log *zerolog.Logger) ([]grpc.ServerOption, error) {

	var options = []grpc.ServerOption{}

	if conf.OverseerCA != "" {
		if err := cert.RegisterCA(conf.OverseerCA); err != nil {
			log.Warn().Err(err).Msg("registering ca")
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

func buildInterceptors(log *zerolog.Logger) []grpc.ServerOption {

	lhandler := handlers.NewLogHandler(log)

	opt := []grpc.ServerOption{grpc.ChainUnaryInterceptor(lhandler.Log), grpc.ChainStreamInterceptor(lhandler.StreamLog)}

	return opt
}
