package services

import (
	"context"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/auth"
	"overseer/overseer/config"
	"overseer/proto/services"
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
		return res, nil
	}

	if msg.Username == "" || msg.Password == "" {
		res.Success = false
		res.Message = "username or password cannot be empty"
		return res, nil
	}

	if ok, err := srv.authmanager.Authenticate(ctx, msg.Username, msg.Password); ok != true || err != nil {
		res.Success = false
		res.Message = "invalid username or password"
		srv.log.Info("authentication failed:", msg.Username)
		return res, nil
	}

	token, err := srv.tokenManager.Create(msg.Username, map[string]interface{}{})
	if err != nil {
		res.Success = false
		res.Message = err.Error()
		return res, nil
	}

	res.Success = true
	res.Message = token

	return res, nil
}
