// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.1
// source: wservices.proto

package wservices

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TaskExecutionServiceClient is the client API for TaskExecutionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TaskExecutionServiceClient interface {
	StartTask(ctx context.Context, in *StartTaskMsg, opts ...grpc.CallOption) (*TaskExecutionResponseMsg, error)
	TerminateTask(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (*WorkerActionMsg, error)
	CompleteTask(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (*WorkerActionMsg, error)
	TaskStatus(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (*TaskExecutionResponseMsg, error)
	TaskOutput(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (TaskExecutionService_TaskOutputClient, error)
	WorkerStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*WorkerStatusResponseMsg, error)
}

type taskExecutionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewTaskExecutionServiceClient(cc grpc.ClientConnInterface) TaskExecutionServiceClient {
	return &taskExecutionServiceClient{cc}
}

func (c *taskExecutionServiceClient) StartTask(ctx context.Context, in *StartTaskMsg, opts ...grpc.CallOption) (*TaskExecutionResponseMsg, error) {
	out := new(TaskExecutionResponseMsg)
	err := c.cc.Invoke(ctx, "/proto.TaskExecutionService/StartTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionServiceClient) TerminateTask(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (*WorkerActionMsg, error) {
	out := new(WorkerActionMsg)
	err := c.cc.Invoke(ctx, "/proto.TaskExecutionService/TerminateTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionServiceClient) CompleteTask(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (*WorkerActionMsg, error) {
	out := new(WorkerActionMsg)
	err := c.cc.Invoke(ctx, "/proto.TaskExecutionService/CompleteTask", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionServiceClient) TaskStatus(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (*TaskExecutionResponseMsg, error) {
	out := new(TaskExecutionResponseMsg)
	err := c.cc.Invoke(ctx, "/proto.TaskExecutionService/TaskStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *taskExecutionServiceClient) TaskOutput(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (TaskExecutionService_TaskOutputClient, error) {
	stream, err := c.cc.NewStream(ctx, &TaskExecutionService_ServiceDesc.Streams[0], "/proto.TaskExecutionService/TaskOutput", opts...)
	if err != nil {
		return nil, err
	}
	x := &taskExecutionServiceTaskOutputClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type TaskExecutionService_TaskOutputClient interface {
	Recv() (*TaskOutputMsg, error)
	grpc.ClientStream
}

type taskExecutionServiceTaskOutputClient struct {
	grpc.ClientStream
}

func (x *taskExecutionServiceTaskOutputClient) Recv() (*TaskOutputMsg, error) {
	m := new(TaskOutputMsg)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *taskExecutionServiceClient) WorkerStatus(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*WorkerStatusResponseMsg, error) {
	out := new(WorkerStatusResponseMsg)
	err := c.cc.Invoke(ctx, "/proto.TaskExecutionService/WorkerStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskExecutionServiceServer is the server API for TaskExecutionService service.
// All implementations must embed UnimplementedTaskExecutionServiceServer
// for forward compatibility
type TaskExecutionServiceServer interface {
	StartTask(context.Context, *StartTaskMsg) (*TaskExecutionResponseMsg, error)
	TerminateTask(context.Context, *TaskIdMsg) (*WorkerActionMsg, error)
	CompleteTask(context.Context, *TaskIdMsg) (*WorkerActionMsg, error)
	TaskStatus(context.Context, *TaskIdMsg) (*TaskExecutionResponseMsg, error)
	TaskOutput(*TaskIdMsg, TaskExecutionService_TaskOutputServer) error
	WorkerStatus(context.Context, *emptypb.Empty) (*WorkerStatusResponseMsg, error)
	mustEmbedUnimplementedTaskExecutionServiceServer()
}

// UnimplementedTaskExecutionServiceServer must be embedded to have forward compatible implementations.
type UnimplementedTaskExecutionServiceServer struct {
}

func (UnimplementedTaskExecutionServiceServer) StartTask(context.Context, *StartTaskMsg) (*TaskExecutionResponseMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTask not implemented")
}
func (UnimplementedTaskExecutionServiceServer) TerminateTask(context.Context, *TaskIdMsg) (*WorkerActionMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TerminateTask not implemented")
}
func (UnimplementedTaskExecutionServiceServer) CompleteTask(context.Context, *TaskIdMsg) (*WorkerActionMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteTask not implemented")
}
func (UnimplementedTaskExecutionServiceServer) TaskStatus(context.Context, *TaskIdMsg) (*TaskExecutionResponseMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskStatus not implemented")
}
func (UnimplementedTaskExecutionServiceServer) TaskOutput(*TaskIdMsg, TaskExecutionService_TaskOutputServer) error {
	return status.Errorf(codes.Unimplemented, "method TaskOutput not implemented")
}
func (UnimplementedTaskExecutionServiceServer) WorkerStatus(context.Context, *emptypb.Empty) (*WorkerStatusResponseMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WorkerStatus not implemented")
}
func (UnimplementedTaskExecutionServiceServer) mustEmbedUnimplementedTaskExecutionServiceServer() {}

// UnsafeTaskExecutionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TaskExecutionServiceServer will
// result in compilation errors.
type UnsafeTaskExecutionServiceServer interface {
	mustEmbedUnimplementedTaskExecutionServiceServer()
}

func RegisterTaskExecutionServiceServer(s grpc.ServiceRegistrar, srv TaskExecutionServiceServer) {
	s.RegisterService(&TaskExecutionService_ServiceDesc, srv)
}

func _TaskExecutionService_StartTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartTaskMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServiceServer).StartTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TaskExecutionService/StartTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServiceServer).StartTask(ctx, req.(*StartTaskMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecutionService_TerminateTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskIdMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServiceServer).TerminateTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TaskExecutionService/TerminateTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServiceServer).TerminateTask(ctx, req.(*TaskIdMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecutionService_CompleteTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskIdMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServiceServer).CompleteTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TaskExecutionService/CompleteTask",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServiceServer).CompleteTask(ctx, req.(*TaskIdMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecutionService_TaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TaskIdMsg)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServiceServer).TaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TaskExecutionService/TaskStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServiceServer).TaskStatus(ctx, req.(*TaskIdMsg))
	}
	return interceptor(ctx, in, info, handler)
}

func _TaskExecutionService_TaskOutput_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(TaskIdMsg)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(TaskExecutionServiceServer).TaskOutput(m, &taskExecutionServiceTaskOutputServer{stream})
}

type TaskExecutionService_TaskOutputServer interface {
	Send(*TaskOutputMsg) error
	grpc.ServerStream
}

type taskExecutionServiceTaskOutputServer struct {
	grpc.ServerStream
}

func (x *taskExecutionServiceTaskOutputServer) Send(m *TaskOutputMsg) error {
	return x.ServerStream.SendMsg(m)
}

func _TaskExecutionService_WorkerStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TaskExecutionServiceServer).WorkerStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.TaskExecutionService/WorkerStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TaskExecutionServiceServer).WorkerStatus(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// TaskExecutionService_ServiceDesc is the grpc.ServiceDesc for TaskExecutionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TaskExecutionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.TaskExecutionService",
	HandlerType: (*TaskExecutionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "StartTask",
			Handler:    _TaskExecutionService_StartTask_Handler,
		},
		{
			MethodName: "TerminateTask",
			Handler:    _TaskExecutionService_TerminateTask_Handler,
		},
		{
			MethodName: "CompleteTask",
			Handler:    _TaskExecutionService_CompleteTask_Handler,
		},
		{
			MethodName: "TaskStatus",
			Handler:    _TaskExecutionService_TaskStatus_Handler,
		},
		{
			MethodName: "WorkerStatus",
			Handler:    _TaskExecutionService_WorkerStatus_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "TaskOutput",
			Handler:       _TaskExecutionService_TaskOutput_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "wservices.proto",
}
