package handlers

import (
	"context"

	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
)

// LogHandler - implements middleware handlers
type LogHandler struct {
	log *zerolog.Logger
}

const (
	serviceKey   = "service"
	methodKey    = "method"
	messageKey   = "message"
	errorKey     = "error"
	workerErrKey = "worker-error"
	receiveKey   = "grpc-worker-receive"
)

// NewLogHandler - creates a new LogHandler
func NewLogHandler(log *zerolog.Logger) *LogHandler {

	return &LogHandler{log: log}
}

// Log - a unary handler
func (h *LogHandler) Log(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	nctx := log.Logger.With().Str("component", "grpc-server").
		Str("service", path.Dir(info.FullMethod)[1:]).
		Str("method", path.Base(info.FullMethod)).Logger().WithContext(ctx)

	//zlog.Info(receiveKey)

	resp, err := handler(nctx, req)

	return resp, err
}

// StreamLog - a stream handler
func (h *LogHandler) StreamLog(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	return handler(srv, ss)
}
