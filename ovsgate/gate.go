package ovsgate

import (
	"context"
	"fmt"
	"net/http"

	"github.com/przebro/overseer/common/cert"
	"github.com/przebro/overseer/common/types"
	"github.com/przebro/overseer/ovsgate/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"time"

	"github.com/przebro/overseer/proto/services"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// OverseerGateway - implements grpc gateway
type OverseerGateway struct {
	config config.OverseerGatewayConfig
	log    *zerolog.Logger
}

const (
	apiHandlePatter     = "/api/"
	fileServerDirectory = "./www/swagger"
	utilsDirectory      = "./www/utils"
	docsHandlerPattern  = "/docs/"
	utilsHandlerPattern = "/utils/"
	cssHandlerPatter    = "/css/"
	jsHandlerPatter     = "/js/"
)

// NewInstance - creates a new instance of a OverseerGateway
func NewInstance(cfg config.OverseerGatewayConfig) (*OverseerGateway, error) {

	log := log.With().Str("component", "gateway").Logger()
	g := &OverseerGateway{config: cfg, log: &log}

	return g, nil
}

// Start - starts gateway
func (g *OverseerGateway) Start() error {

	var err error
	var conn *grpc.ClientConn
	log.Info().Msg("starting gateway")

	dctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	opt := []grpc.DialOption{grpc.WithBlock()}

	if g.config.SecurityLevel == types.ConnectionSecurityLevelNone {
		opt = append(opt, grpc.WithInsecure())
	} else {

		if err := cert.RegisterCA(g.config.OverseerCA); err != nil {
			log.Error().Err(err).Msg("register ca failed")
		}

		creds, err := cert.BuildClientCredentials(g.config.GatewayCert, g.config.GatewayKey, g.config.GatewayCertPolicy, g.config.SecurityLevel)
		if err != nil {
			return fmt.Errorf("failed to initialize connection %v", err)
		}

		opt = append(opt, creds)
	}

	if conn, err = grpc.DialContext(dctx, fmt.Sprintf("%s:%d", g.config.OverseerAddress, g.config.OverseerPort), opt...); err != nil {
		log.Error().Err(err).Msg("dial overseer failed")
		return err
	}

	gwmux := runtime.NewServeMux(
		runtime.WithMetadata(grpcMetadataHandler),
		runtime.WithErrorHandler(grpcErrorHandler),
	)

	if err = initializeServices(context.Background(), gwmux, conn); err != nil {
		log.Error().Err(err).Msg("initialize services failed")
		return err
	}

	mux := http.NewServeMux()
	g.setupHandlers(mux, gwmux)

	server := &http.Server{Addr: fmt.Sprintf("%s:%d", g.config.GatewayAddress, g.config.GatewayPort), Handler: mux}

	if err = server.ListenAndServe(); err != http.ErrServerClosed {
		log.Error().Err(err).Msg("listen and serve failed")
		return err
	}

	return nil
}

func (g *OverseerGateway) Shutdown() error {

	return nil
}

func (g *OverseerGateway) setupHandlers(mux *http.ServeMux, gwmux *runtime.ServeMux) {

	mux.Handle(apiHandlePatter, newHttpInterceptor(gwmux, g.log))
	fs := http.FileServer(http.Dir(fileServerDirectory))
	mux.Handle(docsHandlerPattern, http.StripPrefix(docsHandlerPattern, fs))

	ui := http.FileServer(http.Dir(utilsDirectory))
	mux.Handle(utilsHandlerPattern, http.StripPrefix(utilsHandlerPattern, ui))
	mux.Handle(cssHandlerPatter, ui)
	mux.Handle(jsHandlerPatter, ui)

}
func initializeServices(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {

	for _, fn := range []func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error{
		services.RegisterAuthenticateServiceHandler,
		services.RegisterResourceServiceHandler,
		services.RegisterTaskServiceHandler,
		services.RegisterDefinitionServiceHandler,
		services.RegisterAdministrationServiceHandler,
	} {
		if err := fn(ctx, mux, conn); err != nil {
			return err
		}
	}

	return nil
}
