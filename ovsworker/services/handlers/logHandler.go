package handlers

import (
	"context"
	"path"

	"github.com/przebro/overseer/common/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

//LogHandler - implements middleware handlers
type LogHandler struct {
	log logger.AppLogger
}

const (
	serviceKey   = "service"
	methodKey    = "method"
	messageKey   = "message"
	errorKey     = "error"
	workerErrKey = "worker-error"
	receiveKey   = "grpc-worker-receive"
)

//NewLogHandler - creates a new LogHandler
func NewLogHandler(log logger.AppLogger) *LogHandler {

	return &LogHandler{log: log}
}

//Log - a unary handler
func (h *LogHandler) Log(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	zlog := h.log.Desugar().With(
		zap.String(serviceKey, path.Dir(info.FullMethod)[1:]),
		zap.String(methodKey, path.Base(info.FullMethod)),
	)

	zlog.Info(receiveKey)

	resp, err := handler(ctx, req)

	return resp, err
}

//StreamLog - a stream handler
func (h *LogHandler) StreamLog(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	return handler(srv, ss)
}
