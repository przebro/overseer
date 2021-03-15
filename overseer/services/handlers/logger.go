package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"overseer/common/logger"
	"path"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

const (
	streamRcvKey  = "grpc-rcv-stream"
	streamSendKey = "grpc-send-stream"
	unarySendKey  = "grpc-send"
	unaryRcvKey   = "grpc-rcv"
	serviceKey    = "service"
	methodKey     = "method"
	messageKey    = "message"
)

//ServiceLoggerHandler -  provides both unary and stream handlers for middleware logging
type ServiceLoggerHandler struct {
	log logger.AppLogger
}

//NewServiceLoggerHandler - creates a new ServiceLoggerHandler
func NewServiceLoggerHandler(log logger.AppLogger) *ServiceLoggerHandler {

	return &ServiceLoggerHandler{log: log}
}

//Log - a unary handler
func (lp *ServiceLoggerHandler) Log(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	log := getLoggerWithFlds(lp.log, info.FullMethod)
	logMessage(log, unaryRcvKey, req)
	resp, err := handler(ctx, req)
	if err == nil {
		logMessage(log, unarySendKey, resp)
	}

	return resp, err
}

//StreamLog - a stream handler
func (lp *ServiceLoggerHandler) StreamLog(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	wstream := wrapLogServerStream(ss, lp.log, info.FullMethod)

	return handler(srv, wstream)
}

//GetUnaryHandler - gets a unary handler
func (lp *ServiceLoggerHandler) GetUnaryHandler() grpc.UnaryServerInterceptor {

	return lp.Log
}

//GetStreamHandler - gets a stream handler
func (lp *ServiceLoggerHandler) GetStreamHandler() grpc.StreamServerInterceptor {
	return lp.StreamLog
}

func wrapLogServerStream(stream grpc.ServerStream, log logger.AppLogger, method string) *wrappedLogServerStream {

	return &wrappedLogServerStream{ServerStream: stream, log: log, method: method}

}

type wrappedLogServerStream struct {
	grpc.ServerStream
	log    logger.AppLogger
	method string
}

func (w *wrappedLogServerStream) SendMsg(m interface{}) error {

	err := w.ServerStream.SendMsg(m)
	if err == nil {
		log := getLoggerWithFlds(w.log, w.method)
		logMessage(log, streamSendKey, m)
	}

	return err
}

func (w *wrappedLogServerStream) RecvMsg(m interface{}) error {

	err := w.ServerStream.RecvMsg(m)

	if err == nil {
		log := getLoggerWithFlds(w.log, w.method)
		logMessage(log, streamRcvKey, m)
	}

	return err
}

func (w *wrappedLogServerStream) Context() context.Context {
	return w.Context()
}

func getLoggerWithFlds(log logger.AppLogger, method string) *zap.Logger {

	return log.Desugar().With(
		zap.String(serviceKey, path.Dir(method)[1:]),
		zap.String(methodKey, path.Base(method)),
	)
}

func logMessage(log *zap.Logger, key string, gmsg interface{}) {

	if m, ok := gmsg.(proto.Message); ok {

		log.Check(zapcore.InfoLevel, "grpc-info-level").Write(zap.Object(key, &jsonMarshaller{protomsg: m}))
	}
}

type jsonMarshaller struct {
	protomsg proto.Message
}

func (j *jsonMarshaller) MarshalLogObject(e zapcore.ObjectEncoder) error {

	return e.AddReflected(messageKey, j)
}

func (j *jsonMarshaller) MarshalJSON() ([]byte, error) {

	b, err := json.Marshal(j.protomsg)
	if err != nil {
		return nil, fmt.Errorf("json serializer failed: %v", err)
	}
	return b, nil
}
