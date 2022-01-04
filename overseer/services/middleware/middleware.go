package middleware

import (
	"google.golang.org/grpc"
)

var mwUnaryHandlers []UnaryHandler = []UnaryHandler{}
var mwStreamHandlers []StreamHandler = []StreamHandler{}

//UnaryHandler - common interface for grpc service handlers
type UnaryHandler interface {
	GetUnaryHandler() grpc.UnaryServerInterceptor
}

//StreamHandler - common interface for grpc stream service handlers
type StreamHandler interface {
	GetStreamHandler() grpc.StreamServerInterceptor
}

//RegisterHandler - registers a new MiddlewareUnaryHandler
func RegisterHandler(h UnaryHandler) {

	if h != nil {
		mwUnaryHandlers = append(mwUnaryHandlers, h)
	}
}

//RegisterStreamHandler - registers a new MiddlewareUnaryHandler
func RegisterStreamHandler(h StreamHandler) {

	if h != nil {
		mwStreamHandlers = append(mwStreamHandlers, h)
	}
}

//GetUnaryHandlers - gets unary handlers
func GetUnaryHandlers() []grpc.UnaryServerInterceptor {

	handlers := []grpc.UnaryServerInterceptor{}

	for x := range mwUnaryHandlers {
		handlers = append(handlers, mwUnaryHandlers[x].GetUnaryHandler())
	}

	return handlers
}

//GetStreamHandlers - gets stream handlers
func GetStreamHandlers() []grpc.StreamServerInterceptor {

	handlers := []grpc.StreamServerInterceptor{}

	for x := range mwStreamHandlers {
		handlers = append(handlers, mwStreamHandlers[x].GetStreamHandler())
	}

	return handlers
}
