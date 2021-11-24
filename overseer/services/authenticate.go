package services

import (
	"context"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/auth"
	"overseer/overseer/config"
	"overseer/proto/services"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ovsAuthenticateService struct {
	authmanager    auth.AuthenticationManager
	tokenManager   *auth.TokenCreatorVerifier
	allowAnonymous bool
	log            logger.AppLogger
}

//NewAuthenticateService - Creates a new Authentication service
func NewAuthenticateService(security config.SecurityConfiguration, tm *auth.TokenCreatorVerifier, provider *datastore.Provider, log logger.AppLogger) (services.AuthenticateServiceServer, error) {

	authman, err := auth.NewAuthenticationManager(security.Collection, provider)
	if err != nil {
		return nil, err
	}

	dservice := &ovsAuthenticateService{log: log, authmanager: authman, tokenManager: tm, allowAnonymous: security.AllowAnonymous}
	return dservice, nil
}

func (srv *ovsAuthenticateService) Authenticate(ctx context.Context, msg *services.AuthorizeActionMsg) (*services.ActionResultMsg, error) {

	res := &services.ActionResultMsg{}

	if msg.Username == "" && msg.Password == "" && srv.allowAnonymous {
		res.Success = true
		res.Message = "anonymous access"
		srv.log.Info("authentication ok:", res.Message)
		return res, nil
	}

	if msg.Username == "" || msg.Password == "" {
		srv.log.Info("authentication failed:", res.Message)
		return res, status.Error(codes.Unauthenticated, "username and password cannot be empty")
	}

	if ok, err := srv.authmanager.Authenticate(ctx, msg.Username, msg.Password); !ok || err != nil {
		srv.log.Info("authentication failed:", msg.Username)
		return res, status.Error(codes.Unauthenticated, "invalid username or password")
	}

	token, err := srv.tokenManager.Create(msg.Username, map[string]interface{}{})
	if err != nil {
		srv.log.Info("authentication failed:", res.Message)
		return res, status.Error(codes.Internal, err.Error())
	}

	res.Success = true
	res.Message = token

	return res, nil
}
