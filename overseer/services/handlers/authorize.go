package handlers

import (
	"context"
	"errors"
	"overseer/common/logger"
	"overseer/datastore"
	"overseer/overseer/auth"
	"overseer/overseer/config"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

const bearer = "Bearer "
const authorization = "Authorization"
const usernameKey = "username"
const anonymousKey = "<ANONYMOUS>"

var (
	ErrUnauthorizedAccess = "unauthorized access"
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotAuthorized  = errors.New("user not authorized to perform this action")
)

//AccessRestricter - restricts access to services
type AccessRestricter interface {
	GetAllowedAction(method string) auth.UserAction
}

//ServiceAuthorizeHandler - provides both unary and stream handlers for middleware authorization
type ServiceAuthorizeHandler struct {
	authman        *auth.AuthorizationManager
	verifier       *auth.TokenCreatorVerifier
	allowAnonymous bool
	log            logger.AppLogger
}

//NewServiceAuthorizeHandler - creates a new instance of a ServiceAuthorizeHandler
func NewServiceAuthorizeHandler(security config.SecurityConfiguration, tm *auth.TokenCreatorVerifier, provider *datastore.Provider, log logger.AppLogger) (*ServiceAuthorizeHandler, error) {

	var authman *auth.AuthorizationManager
	var anonymous = security.AllowAnonymous
	var err error

	if authman, err = auth.NewAuthorizationManager(security, provider); err != nil {
		return nil, err
	}

	return &ServiceAuthorizeHandler{authman: authman, verifier: tm, allowAnonymous: anonymous, log: log}, nil

}

//Authorize - handler for an unary service
func (ap *ServiceAuthorizeHandler) Authorize(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	var ok bool
	var restrictedService AccessRestricter
	var token string
	var username string
	var err error

	/*service isn't protected, in final implementation almost every service should be. A single exception is the Authenticate service,
	it is entry point for restricted access and token generator therefore it should be always accessible.
	*/
	if restrictedService, ok = info.Server.(AccessRestricter); !ok {
		return handler(ctx, req)
	}

	if token = ap.getTokenFromContext(ctx); token == "" && ap.allowAnonymous == false {
		return nil, status.Error(codes.Unauthenticated, ErrUnauthorizedAccess)
	}

	if token == "" && ap.allowAnonymous == true {
		nctx := context.WithValue(ctx, interface{}(usernameKey), anonymousKey)
		return handler(nctx, req)
	}

	if username, err = ap.verifyAccess(ctx, restrictedService, token, info.FullMethod); err != nil {
		return nil, status.Error(codes.Unauthenticated, ErrUnauthorizedAccess)
	}

	nctx := context.WithValue(ctx, interface{}(usernameKey), username)

	return handler(nctx, req)
}

//StreamAuthorize - handler for a stream service
func (ap *ServiceAuthorizeHandler) StreamAuthorize(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	var ok bool
	var token string
	var username string
	var err error
	var restrictedService AccessRestricter

	if restrictedService, ok = srv.(AccessRestricter); !ok {
		return handler(srv, ss)
	}

	if token = ap.getTokenFromContext(ss.Context()); token == "" && ap.allowAnonymous == false {
		return status.Error(codes.Unauthenticated, ErrUnauthorizedAccess)
	}

	if token == "" && ap.allowAnonymous == true {
		return handler(srv, ss)
	}

	if username, err = ap.verifyAccess(ss.Context(), restrictedService, "", info.FullMethod); err != nil {
		return err
	}

	nctx := context.WithValue(ss.Context(), interface{}(usernameKey), username)
	wrapped := wrapServerStream(ss)
	wrapped.WrappedContext = nctx

	return handler(srv, wrapped)
}

//GetUnaryHandler - gets a unary handler
func (ap *ServiceAuthorizeHandler) GetUnaryHandler() grpc.UnaryServerInterceptor {

	return ap.Authorize
}

//GetStreamHandler - gets a stream handler
func (ap *ServiceAuthorizeHandler) GetStreamHandler() grpc.StreamServerInterceptor {
	return ap.StreamAuthorize
}

func (ap *ServiceAuthorizeHandler) verifyAccess(ctx context.Context, rservice AccessRestricter, token string, method string) (string, error) {

	var result bool
	var err error

	username, err := ap.verifier.Verify(token)
	if err != nil {
		return "", ErrInvalidToken
	}

	if result, err = ap.authman.VerifyAction(ctx, rservice.GetAllowedAction(method), username); err != nil || result == false {
		return "", ErrUserNotAuthorized
	}

	return username, nil
}

func (ap *ServiceAuthorizeHandler) getTokenFromContext(ctx context.Context) string {

	var meta metadata.MD
	var ok bool

	if meta, ok = metadata.FromIncomingContext(ctx); !ok {
		return ""
	}

	authData := meta.Get(authorization)

	if len(authData) == 0 {
		ap.log.Info("authorization header not found metadata")

		return ""
	}

	return strings.TrimPrefix(authData[0], bearer)
}

type wrappedServerStream struct {
	grpc.ServerStream
	// WrappedContext is the wrapper's own Context. You can assign it.
	WrappedContext context.Context
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context()
func (w *wrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

func wrapServerStream(stream grpc.ServerStream) *wrappedServerStream {
	if existing, ok := stream.(*wrappedServerStream); ok {
		return existing
	}
	return &wrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}
