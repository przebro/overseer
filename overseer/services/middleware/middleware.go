package middleware

import (
	"fmt"

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

func GetUnaryHandlers() []grpc.UnaryServerInterceptor {

	fmt.Println("getting handlers")
	handlers := []grpc.UnaryServerInterceptor{}

	for x := range mwUnaryHandlers {
		handlers = append(handlers, mwUnaryHandlers[x].GetUnaryHandler())
	}

	return handlers
}

func GetStreamHandlers() []grpc.StreamServerInterceptor {

	handlers := []grpc.StreamServerInterceptor{}

	return handlers
}
