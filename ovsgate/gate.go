package ovsgate

import (
	"context"
	"fmt"
	"net/http"
	"overseer/common/cert"
	"overseer/common/logger"
	"overseer/ovsgate/config"

	"overseer/proto/services"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type OverseerGateway struct {
	config *config.OverseerGatewayConfig
	log    logger.AppLogger
}

//NewInstance - creates a new instance of a OverseerGateway
func NewInstance(cfg *config.OverseerGatewayConfig, log logger.AppLogger) (*OverseerGateway, error) {

	g := &OverseerGateway{config: cfg, log: log}

	return g, nil
}

func (g *OverseerGateway) Start() error {

	var err error
	var conn *grpc.ClientConn

	dctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opt := []grpc.DialOption{grpc.WithBlock()}

	if !g.config.UseTLS {
		opt = append(opt, grpc.WithInsecure())
	} else {

		if creds, err := cert.GetClientTLS(g.config.CertPath, false); err != nil {
			return fmt.Errorf("failed to initialize connection %v", err)
		} else {
			opt = append(opt, grpc.WithTransportCredentials(creds))
		}
	}

	if conn, err = grpc.DialContext(dctx, fmt.Sprintf("%s:%d", g.config.OverseerAddress, g.config.OverseerPort), opt...); err != nil {
		g.log.Error(err)
		return err
	}

	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(grpcMetadataHandler),
		runtime.WithErrorHandler(grpcErrorHandler),
	)

	if err = initializeServices(context.Background(), gwmux, conn); err != nil {
		g.log.Error(err)
		return err
	}

	mux := http.NewServeMux()
	g.setupHandlers(mux, gwmux)

	server := &http.Server{Addr: fmt.Sprintf("%s:%d", g.config.GatewayAddress, g.config.GatewayPort), Handler: mux}
	if err = server.ListenAndServe(); err != http.ErrServerClosed {
		g.log.Error(err)
		return err
	}

	return nil
}
func (g *OverseerGateway) setupHandlers(mux *http.ServeMux, gwmux *runtime.ServeMux) {

	mux.Handle("/api/", newHttpInterceptor(gwmux, g.log))
	fs := http.FileServer(http.Dir("./www/swagger"))
	mux.Handle("/docs/", http.StripPrefix("/docs/", fs))

	ui := http.FileServer(http.Dir("./www/utils"))
	mux.Handle("/utils/", http.StripPrefix("/utils/", ui))
	mux.Handle("/css/", ui)
	mux.Handle("/js/", ui)

}
func initializeServices(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {

	for _, fn := range []func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error{
		services.RegisterAuthenticateServiceHandler,
		services.RegisterResourceServiceHandler,
		services.RegisterTaskServiceHandler,
		services.RegisterDefinitionServiceHandler,
	} {
		if err := fn(ctx, mux, conn); err != nil {
			return err
		}
	}

	return nil
}
