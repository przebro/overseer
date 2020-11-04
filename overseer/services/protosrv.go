package services

import (
	"context"
	"fmt"
	"goscheduler/common/logger"
	"goscheduler/overseer/internal/events"
	"goscheduler/overseer/internal/resources"
	"goscheduler/overseer/internal/taskdef"
	"goscheduler/proto/services"
	"net"

	"google.golang.org/grpc"
)

type ovsGrpcServer struct {
	grpcServer *grpc.Server
	log        logger.AppLogger
	rm         resources.ResourceManager
	dm         taskdef.TaskDefinitionManager
	rservice   ovsResourceService
	dservice   ovsDefinitionService
	dispatcher events.Dispatcher
}

//OvsGrpcServer - interface for ovsGrpcServer
type OvsGrpcServer interface {
	Listen(host string, port int) error
}

//NewOvsGrpcServer - Creates new instance of a ovsGrpcServer
func NewOvsGrpcServer(disp events.Dispatcher,
	r services.ResourceServiceServer,
	d services.DefinitionServiceServer,
	t services.TaskServiceServer) OvsGrpcServer {

	srv := &ovsGrpcServer{}
	srv.grpcServer = grpc.NewServer()
	srv.rservice = ovsResourceService{}
	services.RegisterResourceServiceServer(srv.grpcServer, r)
	services.RegisterDefinitionServiceServer(srv.grpcServer, d)
	services.RegisterTaskServiceServer(srv.grpcServer, t)
	services.RegisterAuthorizeServiceServer(srv.grpcServer, srv)
	srv.dispatcher = disp
	srv.log = logger.Get()
	return srv

}

func (srv *ovsGrpcServer) Listen(host string, port int) error {

	conn := fmt.Sprintf("%s:%d", host, port)
	l, err := net.Listen("tcp", conn)
	if err != nil {
		srv.log.Error(err)
		return err
	}

	srv.log.Info("starting grpc server:", conn)

	err = srv.grpcServer.Serve(l)
	if err != nil {

	}

	return nil

}

func (srv *ovsGrpcServer) Authorize(ctx context.Context, msg *services.AuthorizeActionMsg) (*services.ActionResultMsg, error) {

	return nil, nil
}
