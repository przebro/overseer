package services

import (
	"fmt"
	"net"

	"github.com/przebro/overseer/common/cert"
	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/overseer/config"
	"github.com/przebro/overseer/overseer/services/middleware"
	"github.com/przebro/overseer/proto/services"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
)

// OvsGrpcServer - represents core grpc component
type OvsGrpcServer struct {
	conf       config.ServerConfiguration
	grpcServer *grpc.Server
	log        *zerolog.Logger
}

// NewOvsGrpcServer - Creates new instance of a ovsGrpcServer
func NewOvsGrpcServer(
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
	lg := log.With().Str("component", "grpc-server").Logger()
	srv.log = &lg

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

	return srv
}

// Start - starts a service
func (srv *OvsGrpcServer) Start() error {

	conn := fmt.Sprintf("%s:%d", srv.conf.Host, srv.conf.Port)
	l, err := net.Listen("tcp", conn)
	if err != nil {
		srv.log.Error().Err(err).Msg("listener failed")
		return err
	}

	srv.log.Info().Str("connection", conn).Msg("starting grpc server")

	err = srv.grpcServer.Serve(l)
	if err != nil {
		return err
	}

	return nil
}

// Shutdown - shutdowns a service
func (srv *OvsGrpcServer) Shutdown() error {

	srv.grpcServer.GracefulStop()

	return nil
}

func buildOptions(conf config.ServerConfiguration) ([]grpc.ServerOption, error) {

	var options = []grpc.ServerOption{}

	options = append(options, buildUnaryChain())
	options = append(options, buildStreamChain())

	if conf.Security.SecurityLevel != types.ConnectionSecurityLevelNone {

		if creds, err := buildCredentials(conf.Security.ServerCert, conf.Security.ServerKey, conf.Security.ClientCertPolicy, conf.Security.SecurityLevel); err == nil {
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

func buildCredentials(certpath, keypath string, clientPolicy types.CertPolicy, level types.ConnectionSecurityLevel) (grpc.ServerOption, error) {

	var err error

	creds, err := cert.BuildServerCredentials(certpath, keypath, clientPolicy, level)

	if err != nil {
		return nil, err
	}
	return creds, nil
}
