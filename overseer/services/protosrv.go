package services

import (
	"fmt"
	"net"
	"overseer/common/cert"
	"overseer/common/logger"
	"overseer/overseer/config"
	"overseer/overseer/internal/events"
	"overseer/overseer/internal/resources"
	"overseer/overseer/internal/taskdef"
	"overseer/overseer/services/middleware"
	"overseer/proto/services"

	"google.golang.org/grpc"
)

//OvsGrpcServer - represents core grpc component
type OvsGrpcServer struct {
	conf       config.ServerConfiguration
	grpcServer *grpc.Server
	log        logger.AppLogger
	rm         resources.ResourceManager
	dm         taskdef.TaskDefinitionManager
	dservice   ovsDefinitionService
	dispatcher events.Dispatcher
}

//NewOvsGrpcServer - Creates new instance of a ovsGrpcServer
func NewOvsGrpcServer(disp events.Dispatcher,
	r services.ResourceServiceServer,
	d services.DefinitionServiceServer,
	t services.TaskServiceServer,
	a services.AuthenticateServiceServer,
	adm services.AdministrationServiceServer,
	stat services.StatusServiceServer,
	config config.ServerConfiguration,

) *OvsGrpcServer {

	var options []grpc.ServerOption
	var err error

	srv := &OvsGrpcServer{}
	srv.conf = config

	if options, err = buildOptions(config); err != nil {
		return nil
	}

	srv.grpcServer = grpc.NewServer(options...)
	services.RegisterResourceServiceServer(srv.grpcServer, r)
	services.RegisterDefinitionServiceServer(srv.grpcServer, d)
	services.RegisterTaskServiceServer(srv.grpcServer, t)
	services.RegisterAuthenticateServiceServer(srv.grpcServer, a)
	services.RegisterAdministrationServiceServer(srv.grpcServer, adm)
	services.RegisterStatusServiceServer(srv.grpcServer, stat)
	srv.dispatcher = disp
	srv.log = logger.Get()

	return srv
}

func (srv *OvsGrpcServer) Start() error {

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

func (srv *OvsGrpcServer) Shutdown() error {

	srv.grpcServer.GracefulStop()

	return nil
}

func buildOptions(conf config.ServerConfiguration) ([]grpc.ServerOption, error) {

	var options = []grpc.ServerOption{}

	options = append(options, buildUnaryChain())
	options = append(options, buildStreamChain())

	if conf.TLS {
		if creds, err := buildCredentials(conf.ServerCert, conf.ServerKey); err == nil {
			options = append(options, creds)
		} else {
			return nil, err
		}
	}

	return options, nil
}

func buildUnaryChain() grpc.ServerOption {

	return grpc.ChainUnaryInterceptor(middleware.GetUnaryHandlers()...)
}

func buildStreamChain() grpc.ServerOption {

	return grpc.ChainStreamInterceptor(middleware.GetStreamHandlers()...)
}

func buildCredentials(certpath, keypath string) (grpc.ServerOption, error) {

	c, err := cert.GetServerTLS(certpath, keypath)
	if err != nil {
		return nil, err
	}
	return grpc.Creds(c), nil
}
