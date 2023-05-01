package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"path"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

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

// ServiceLoggerHandler -  provides both unary and stream handlers for middleware logging
type ServiceLoggerHandler struct {
	log *zerolog.Logger
}

// NewServiceLoggerHandler - creates a new ServiceLoggerHandler
func NewServiceLoggerHandler(log *zerolog.Logger) *ServiceLoggerHandler {

	return &ServiceLoggerHandler{log: log}
}

// Log - a unary handler
func (lp *ServiceLoggerHandler) Log(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	nctx := log.Logger.With().Str("component", "grpc-server").Logger().WithContext(ctx)
	//log := getLoggerWithFlds(lp.log, info.FullMethod)
	log := zerolog.Ctx(nctx)
	setMethodInContext(log, info.FullMethod)

	logMessage(log, unaryRcvKey, req)
	resp, err := handler(nctx, req)
	if err == nil {
		logMessage(log, unarySendKey, resp)
	}

	return resp, err
}

// StreamLog - a stream handler
func (lp *ServiceLoggerHandler) StreamLog(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

	wstream := wrapLogServerStream(ss, lp.log, info.FullMethod)

	return handler(srv, wstream)
}

// GetUnaryHandler - gets a unary handler
func (lp *ServiceLoggerHandler) GetUnaryHandler() grpc.UnaryServerInterceptor {

	return lp.Log
}

// GetStreamHandler - gets a stream handler
func (lp *ServiceLoggerHandler) GetStreamHandler() grpc.StreamServerInterceptor {
	return lp.StreamLog
}

func wrapLogServerStream(stream grpc.ServerStream, log *zerolog.Logger, method string) *wrappedLogServerStream {

	nctx := log.With().Str("component", "grpc-server").Logger().WithContext(stream.Context())

	return &wrappedLogServerStream{ServerStream: stream, method: method, WrappedContext: nctx}

}

type wrappedLogServerStream struct {
	grpc.ServerStream
	method         string
	WrappedContext context.Context
}

func (w *wrappedLogServerStream) SendMsg(m interface{}) error {

	log := zerolog.Ctx(w.WrappedContext)
	err := w.ServerStream.SendMsg(m)
	if err == nil {
		setMethodInContext(log, w.method)
		logMessage(log, streamSendKey, m)
	}

	return err
}

func (w *wrappedLogServerStream) RecvMsg(m interface{}) error {

	log := zerolog.Ctx(w.WrappedContext)
	err := w.ServerStream.RecvMsg(m)

	if err == nil {

		setMethodInContext(log, w.method)
		logMessage(log, streamRcvKey, m)
	}

	return err
}

func (w *wrappedLogServerStream) Context() context.Context {
	return w.WrappedContext
}

func setMethodInContext(log *zerolog.Logger, method string) {

	log.UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str("service", path.Dir(method)[1:]).Str("method", path.Base(method))
	})
}

func logMessage(log *zerolog.Logger, key string, gmsg interface{}) {

	if m, ok := gmsg.(proto.Message); ok {
		js := &jsonMarshaller{protomsg: m}
		bytes, _ := js.MarshalJSON()
		log.Info().RawJSON("payload", bytes).Msg("grpc-info-level")
	}
}

type jsonMarshaller struct {
	protomsg proto.Message
}

func (j *jsonMarshaller) MarshalJSON() ([]byte, error) {

	b, err := json.Marshal(j.protomsg)
	if err != nil {
		return nil, fmt.Errorf("json serializer failed: %v", err)
	}
	return b, nil
}
