package services

import (
	"context"

	"google.golang.org/grpc/metadata"
)

type MockGrpcServerStream struct {
}

func (m *MockGrpcServerStream) SetHeader(metadata.MD) error   { return nil }
func (m *MockGrpcServerStream) SendHeader(metadata.MD) error  { return nil }
func (m *MockGrpcServerStream) SetTrailer(metadata.MD)        {}
func (m *MockGrpcServerStream) Context() context.Context      { return context.Background() }
func (m *MockGrpcServerStream) SendMsg(msg interface{}) error { return nil }
func (m *MockGrpcServerStream) RecvMsg(msg interface{}) error { return nil }
