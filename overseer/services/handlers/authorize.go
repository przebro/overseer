package handlers

import (
	"context"
	"errors"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/auth"
	"overseer/overseer/config"

	"google.golang.org/grpc/metadata"

	"google.golang.org/grpc"
)

var (
	ErrUnauthorizedAccess = errors.New("unauthorized access")
)

type AccessRestricter interface {
	GetAllowedAction(method string) auth.UserAction
}

// type AuthorizationHandler interface {
// 	Authorize(ctx context.Context, service string) (context.Context, error)
// }

type ServiceAuthorizeHandler struct {
	authman        *auth.AuthorizationManager
	verifier       *auth.TokenCreatorVerifier
	allowAnonymous bool
	log            logger.AppLogger
}

func NewServiceAuthorizeHandler(security config.SecurityConfiguration, tm *auth.TokenCreatorVerifier, provider *datastore.Provider) (*ServiceAuthorizeHandler, error) {

	var authman *auth.AuthorizationManager
	var anonymous = security.AllowAnonymous
	var err error

	if authman, err = auth.NewAuthorizationManager(security.Collection, provider); err != nil {
		return nil, err
	}

	return &ServiceAuthorizeHandler{authman: authman, verifier: tm, allowAnonymous: anonymous, log: logger.Get()}, nil

}

func (ap *ServiceAuthorizeHandler) Authorize(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	var ok bool
	var restrictedService AccessRestricter
	var token string

	/*service isn't protected, in final implementation almost every service should be. A single exception is the Authenticate service,
	it is entry point for restricted access and token generator therefore it should be always accessible.
	*/
	if restrictedService, ok = info.Server.(AccessRestricter); !ok {
		return handler(ctx, req)
	}

	if ap.allowAnonymous == true {
		return handler(ctx, req)
	}

	if token = ap.getTokenFromContext(ctx); token == "" {
		return nil, ErrUnauthorizedAccess
	}

	if err := ap.verifyAccess(ctx, restrictedService, token, info.FullMethod); err != nil {
		return nil, err
	}

	return handler(ctx, req)
}

func (ap *ServiceAuthorizeHandler) StreamAuthorize(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	var ok bool
	var token string
	var restrictedService AccessRestricter

	if restrictedService, ok = srv.(AccessRestricter); !ok {
		return handler(srv, ss)
	}

	if ap.allowAnonymous == true {
		return handler(srv, ss)
	}

	if token = ap.getTokenFromContext(ss.Context()); token == "" {
		return ErrUnauthorizedAccess
	}

	if err := ap.verifyAccess(ss.Context(), restrictedService, "", info.FullMethod); err != nil {
		return err
	}

	return handler(srv, ss)
}

func (ap *ServiceAuthorizeHandler) GetUnaryHandler() grpc.UnaryServerInterceptor {

	return ap.Authorize
}
func (ap *ServiceAuthorizeHandler) GetStreamHandler() grpc.StreamServerInterceptor {
	return ap.StreamAuthorize
}

func (ap *ServiceAuthorizeHandler) verifyAccess(ctx context.Context, rservice AccessRestricter, token string, method string) error {

	username, err := ap.verifier.Verify(token)
	if err != nil {
		return nil
	}

	ap.authman.VerifyAction(ctx, rservice.GetAllowedAction(method), username)

	return nil
}

func (ap *ServiceAuthorizeHandler) getTokenFromContext(ctx context.Context) string {

	var meta metadata.MD
	var ok bool

	ap.log.Info("authorization")
	if meta, ok = metadata.FromIncomingContext(ctx); !ok {
		return ""
	}

	ap.log.Info("found metadata")
	authData := meta.Get("authorization")

	if len(authData) == 0 {
		ap.log.Info("authorization header not found metadata")
		return ""
	}

	return authData[0]
}
