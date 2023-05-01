package services

import (
	"context"

	"github.com/przebro/overseer/overseer/config"
	"github.com/przebro/overseer/proto/services"
	"github.com/rs/zerolog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TokenCreatorVerifier interface {
	Verify(token string) (string, error)
	Create(username string, userdata map[string]interface{}) (string, error)
}

type AuthenticationManager interface {
	Authenticate(ctx context.Context, username string, password string) (bool, error)
}

type ovsAuthenticateService struct {
	authmanager    AuthenticationManager
	tokenManager   TokenCreatorVerifier
	allowAnonymous bool
	services.UnimplementedAuthenticateServiceServer
}

// NewAuthenticateService - Creates a new Authentication service
func NewAuthenticateService(security config.SecurityConfiguration, tm TokenCreatorVerifier, authman AuthenticationManager) (services.AuthenticateServiceServer, error) {

	dservice := &ovsAuthenticateService{authmanager: authman, tokenManager: tm, allowAnonymous: security.AllowAnonymous}
	return dservice, nil
}

func (srv *ovsAuthenticateService) Authenticate(ctx context.Context, msg *services.AuthorizeActionMsg) (*services.ActionResultMsg, error) {

	res := &services.ActionResultMsg{}

	log := zerolog.Ctx(ctx).With().Str("service", "authenticate").Logger()

	if msg.Username == "" && msg.Password == "" && srv.allowAnonymous {
		res.Success = true
		res.Message = "anonymous access"
		log.Info().Str("message", res.Message).Msg("authentication successful")
		return res, nil
	}

	if msg.Username == "" || msg.Password == "" {
		log.Info().Str("message", res.Message).Msg("authentication failed")
		return res, status.Error(codes.Unauthenticated, "username and password cannot be empty")
	}

	if ok, err := srv.authmanager.Authenticate(ctx, msg.Username, msg.Password); !ok || err != nil {
		log.Info().Str("username", msg.Username).Msg("authentication failed")
		return res, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	token, err := srv.tokenManager.Create(msg.Username, map[string]interface{}{})
	if err != nil {
		log.Info().Str("message", res.Message).Msg("authentication failed")
		return res, status.Error(codes.Internal, err.Error())
	}

	res.Success = true
	res.Message = token

	return res, nil
}
