package handlers

import (
	"context"
	"errors"
	"fmt"
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

var (
	ErrUnauthorizedAccess = "unauthorized access"
	ErrInvalidToken       = errors.New("invalid token")
	ErrUserNotAuthorized  = errors.New("user not authorized to perform this action")
)

type AccessRestricter interface {
	GetAllowedAction(method string) auth.UserAction
}

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

	if authman, err = auth.NewAuthorizationManager(security, provider); err != nil {
		return nil, err
	}

	return &ServiceAuthorizeHandler{authman: authman, verifier: tm, allowAnonymous: anonymous, log: logger.Get()}, nil

}

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
		nctx := context.WithValue(ctx, interface{}("username"), "<ANONYMOUS>")
		return handler(nctx, req)
	}

	if username, err = ap.verifyAccess(ctx, restrictedService, token, info.FullMethod); err != nil {
		return nil, status.Error(codes.Unauthenticated, ErrUnauthorizedAccess)
	}

	nctx := context.WithValue(ctx, interface{}("username"), username)

	return handler(nctx, req)
}

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
	fmt.Println(username)
	nctx := context.WithValue(ss.Context(), interface{}("username"), username)
	wrapped := WrapServerStream(ss)
	wrapped.WrappedContext = nctx

	return handler(srv, wrapped)
}

func (ap *ServiceAuthorizeHandler) GetUnaryHandler() grpc.UnaryServerInterceptor {

	return ap.Authorize
}
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

type WrappedServerStream struct {
	grpc.ServerStream
	// WrappedContext is the wrapper's own Context. You can assign it.
	WrappedContext context.Context
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context()
func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

func WrapServerStream(stream grpc.ServerStream) *WrappedServerStream {
	if existing, ok := stream.(*WrappedServerStream); ok {
		return existing
	}
	return &WrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}
