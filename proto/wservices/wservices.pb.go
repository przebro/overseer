// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: wservices.proto

package wservices

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	any "github.com/golang/protobuf/ptypes/any"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type TaskExecutionResponseMsg_TaskStatus int32

const (
	TaskExecutionResponseMsg_RECEIVED  TaskExecutionResponseMsg_TaskStatus = 0
	TaskExecutionResponseMsg_EXECUTING TaskExecutionResponseMsg_TaskStatus = 1
	TaskExecutionResponseMsg_ENDED     TaskExecutionResponseMsg_TaskStatus = 2
	TaskExecutionResponseMsg_FAILED    TaskExecutionResponseMsg_TaskStatus = 3
	//for future use
	TaskExecutionResponseMsg_WAITING  TaskExecutionResponseMsg_TaskStatus = 4
	TaskExecutionResponseMsg_IDLE     TaskExecutionResponseMsg_TaskStatus = 5
	TaskExecutionResponseMsg_STARTING TaskExecutionResponseMsg_TaskStatus = 6
)

// Enum value maps for TaskExecutionResponseMsg_TaskStatus.
var (
	TaskExecutionResponseMsg_TaskStatus_name = map[int32]string{
		0: "RECEIVED",
		1: "EXECUTING",
		2: "ENDED",
		3: "FAILED",
		4: "WAITING",
		5: "IDLE",
		6: "STARTING",
	}
	TaskExecutionResponseMsg_TaskStatus_value = map[string]int32{
		"RECEIVED":  0,
		"EXECUTING": 1,
		"ENDED":     2,
		"FAILED":    3,
		"WAITING":   4,
		"IDLE":      5,
		"STARTING":  6,
	}
)

func (x TaskExecutionResponseMsg_TaskStatus) Enum() *TaskExecutionResponseMsg_TaskStatus {
	p := new(TaskExecutionResponseMsg_TaskStatus)
	*p = x
	return p
}

func (x TaskExecutionResponseMsg_TaskStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (TaskExecutionResponseMsg_TaskStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_wservices_proto_enumTypes[0].Descriptor()
}

func (TaskExecutionResponseMsg_TaskStatus) Type() protoreflect.EnumType {
	return &file_wservices_proto_enumTypes[0]
}

func (x TaskExecutionResponseMsg_TaskStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use TaskExecutionResponseMsg_TaskStatus.Descriptor instead.
func (TaskExecutionResponseMsg_TaskStatus) EnumDescriptor() ([]byte, []int) {
	return file_wservices_proto_rawDescGZIP(), []int{2, 0}
}

type TaskIdMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskID string `protobuf:"bytes,1,opt,name=taskID,proto3" json:"taskID,omitempty"`
}

func (x *TaskIdMsg) Reset() {
	*x = TaskIdMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wservices_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskIdMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskIdMsg) ProtoMessage() {}

func (x *TaskIdMsg) ProtoReflect() protoreflect.Message {
	mi := &file_wservices_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskIdMsg.ProtoReflect.Descriptor instead.
func (*TaskIdMsg) Descriptor() ([]byte, []int) {
	return file_wservices_proto_rawDescGZIP(), []int{0}
}

func (x *TaskIdMsg) GetTaskID() string {
	if x != nil {
		return x.TaskID
	}
	return ""
}

type StartTaskMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskID    *TaskIdMsg        `protobuf:"bytes,1,opt,name=taskID,proto3" json:"taskID,omitempty"`
	Type      string            `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	Variables map[string]string `protobuf:"bytes,4,rep,name=variables,proto3" json:"variables,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Command   *any.Any          `protobuf:"bytes,5,opt,name=Command,proto3" json:"Command,omitempty"`
}

func (x *StartTaskMsg) Reset() {
	*x = StartTaskMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wservices_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartTaskMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartTaskMsg) ProtoMessage() {}

func (x *StartTaskMsg) ProtoReflect() protoreflect.Message {
	mi := &file_wservices_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartTaskMsg.ProtoReflect.Descriptor instead.
func (*StartTaskMsg) Descriptor() ([]byte, []int) {
	return file_wservices_proto_rawDescGZIP(), []int{1}
}

func (x *StartTaskMsg) GetTaskID() *TaskIdMsg {
	if x != nil {
		return x.TaskID
	}
	return nil
}

func (x *StartTaskMsg) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *StartTaskMsg) GetVariables() map[string]string {
	if x != nil {
		return x.Variables
	}
	return nil
}

func (x *StartTaskMsg) GetCommand() *any.Any {
	if x != nil {
		return x.Command
	}
	return nil
}

type TaskExecutionResponseMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status     TaskExecutionResponseMsg_TaskStatus `protobuf:"varint,1,opt,name=status,proto3,enum=proto.TaskExecutionResponseMsg_TaskStatus" json:"status,omitempty"`
	ReturnCode int32                               `protobuf:"varint,2,opt,name=returnCode,proto3" json:"returnCode,omitempty"`
	StatusCode int32                               `protobuf:"varint,3,opt,name=statusCode,proto3" json:"statusCode,omitempty"`
	Pid        int32                               `protobuf:"varint,4,opt,name=pid,proto3" json:"pid,omitempty"`
	Output     []string                            `protobuf:"bytes,6,rep,name=output,proto3" json:"output,omitempty"`
}

func (x *TaskExecutionResponseMsg) Reset() {
	*x = TaskExecutionResponseMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wservices_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskExecutionResponseMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskExecutionResponseMsg) ProtoMessage() {}

func (x *TaskExecutionResponseMsg) ProtoReflect() protoreflect.Message {
	mi := &file_wservices_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TaskExecutionResponseMsg.ProtoReflect.Descriptor instead.
func (*TaskExecutionResponseMsg) Descriptor() ([]byte, []int) {
	return file_wservices_proto_rawDescGZIP(), []int{2}
}

func (x *TaskExecutionResponseMsg) GetStatus() TaskExecutionResponseMsg_TaskStatus {
	if x != nil {
		return x.Status
	}
	return TaskExecutionResponseMsg_RECEIVED
}

func (x *TaskExecutionResponseMsg) GetReturnCode() int32 {
	if x != nil {
		return x.ReturnCode
	}
	return 0
}

func (x *TaskExecutionResponseMsg) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

func (x *TaskExecutionResponseMsg) GetPid() int32 {
	if x != nil {
		return x.Pid
	}
	return 0
}

func (x *TaskExecutionResponseMsg) GetOutput() []string {
	if x != nil {
		return x.Output
	}
	return nil
}

type WorkerActionMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *WorkerActionMsg) Reset() {
	*x = WorkerActionMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wservices_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerActionMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerActionMsg) ProtoMessage() {}

func (x *WorkerActionMsg) ProtoReflect() protoreflect.Message {
	mi := &file_wservices_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerActionMsg.ProtoReflect.Descriptor instead.
func (*WorkerActionMsg) Descriptor() ([]byte, []int) {
	return file_wservices_proto_rawDescGZIP(), []int{3}
}

func (x *WorkerActionMsg) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *WorkerActionMsg) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type WorkerStatusResponseMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tasks    int32 `protobuf:"varint,1,opt,name=tasks,proto3" json:"tasks,omitempty"`
	Cpuload  int32 `protobuf:"varint,2,opt,name=cpuload,proto3" json:"cpuload,omitempty"`
	Memused  int32 `protobuf:"varint,3,opt,name=memused,proto3" json:"memused,omitempty"`
	Memtotal int32 `protobuf:"varint,4,opt,name=memtotal,proto3" json:"memtotal,omitempty"`
}

func (x *WorkerStatusResponseMsg) Reset() {
	*x = WorkerStatusResponseMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wservices_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerStatusResponseMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerStatusResponseMsg) ProtoMessage() {}

func (x *WorkerStatusResponseMsg) ProtoReflect() protoreflect.Message {
	mi := &file_wservices_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WorkerStatusResponseMsg.ProtoReflect.Descriptor instead.
func (*WorkerStatusResponseMsg) Descriptor() ([]byte, []int) {
	return file_wservices_proto_rawDescGZIP(), []int{4}
}

func (x *WorkerStatusResponseMsg) GetTasks() int32 {
	if x != nil {
		return x.Tasks
	}
	return 0
}

func (x *WorkerStatusResponseMsg) GetCpuload() int32 {
	if x != nil {
		return x.Cpuload
	}
	return 0
}

func (x *WorkerStatusResponseMsg) GetMemused() int32 {
	if x != nil {
		return x.Memused
	}
	return 0
}

func (x *WorkerStatusResponseMsg) GetMemtotal() int32 {
	if x != nil {
		return x.Memtotal
	}
	return 0
}

var File_wservices_proto protoreflect.FileDescriptor

var file_wservices_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x77, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x23, 0x0a, 0x09, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x16, 0x0a,
	0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74,
	0x61, 0x73, 0x6b, 0x49, 0x44, 0x22, 0xfc, 0x01, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54,
	0x61, 0x73, 0x6b, 0x4d, 0x73, 0x67, 0x12, 0x28, 0x0a, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54,
	0x61, 0x73, 0x6b, 0x49, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x44,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x40, 0x0a, 0x09, 0x76, 0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x4d, 0x73, 0x67, 0x2e, 0x56, 0x61, 0x72,
	0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x76, 0x61, 0x72,
	0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x2e, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52, 0x07, 0x43,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x1a, 0x3c, 0x0a, 0x0e, 0x56, 0x61, 0x72, 0x69, 0x61, 0x62,
	0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x3a, 0x02, 0x38, 0x01, 0x22, 0xaf, 0x02, 0x0a, 0x18, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x73,
	0x67, 0x12, 0x42, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x78,
	0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d,
	0x73, 0x67, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x74, 0x75, 0x72, 0x6e, 0x43,
	0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x72, 0x65, 0x74, 0x75, 0x72,
	0x6e, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43,
	0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x03, 0x70, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75,
	0x74, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22,
	0x65, 0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0c, 0x0a,
	0x08, 0x52, 0x45, 0x43, 0x45, 0x49, 0x56, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x45,
	0x58, 0x45, 0x43, 0x55, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x4e,
	0x44, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10,
	0x03, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x41, 0x49, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x12, 0x08,
	0x0a, 0x04, 0x49, 0x44, 0x4c, 0x45, 0x10, 0x05, 0x12, 0x0c, 0x0a, 0x08, 0x53, 0x54, 0x41, 0x52,
	0x54, 0x49, 0x4e, 0x47, 0x10, 0x06, 0x22, 0x45, 0x0a, 0x0f, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x7f, 0x0a,
	0x17, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x73, 0x67, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x18,
	0x0a, 0x07, 0x63, 0x70, 0x75, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x07, 0x63, 0x70, 0x75, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x6d, 0x75,
	0x73, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d, 0x65, 0x6d, 0x75, 0x73,
	0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x6d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6d, 0x65, 0x6d, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x32, 0xe1,
	0x02, 0x0a, 0x14, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x09, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x54, 0x61, 0x73, 0x6b, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x53, 0x74, 0x61,
	0x72, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x4d, 0x73, 0x67, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x12, 0x3b, 0x0a, 0x0d,
	0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x10, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x4d, 0x73, 0x67, 0x1a,
	0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x12, 0x3a, 0x0a, 0x0c, 0x43, 0x6f, 0x6d,
	0x70, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x4d, 0x73, 0x67, 0x1a, 0x16, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x4d, 0x73, 0x67, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b,
	0x49, 0x64, 0x4d, 0x73, 0x67, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61,
	0x73, 0x6b, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x0c, 0x57, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x1a, 0x1e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x73, 0x67,
	0x22, 0x00, 0x42, 0x11, 0x5a, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_wservices_proto_rawDescOnce sync.Once
	file_wservices_proto_rawDescData = file_wservices_proto_rawDesc
)

func file_wservices_proto_rawDescGZIP() []byte {
	file_wservices_proto_rawDescOnce.Do(func() {
		file_wservices_proto_rawDescData = protoimpl.X.CompressGZIP(file_wservices_proto_rawDescData)
	})
	return file_wservices_proto_rawDescData
}

var file_wservices_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_wservices_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_wservices_proto_goTypes = []interface{}{
	(TaskExecutionResponseMsg_TaskStatus)(0), // 0: proto.TaskExecutionResponseMsg.TaskStatus
	(*TaskIdMsg)(nil),                        // 1: proto.TaskIdMsg
	(*StartTaskMsg)(nil),                     // 2: proto.StartTaskMsg
	(*TaskExecutionResponseMsg)(nil),         // 3: proto.TaskExecutionResponseMsg
	(*WorkerActionMsg)(nil),                  // 4: proto.WorkerActionMsg
	(*WorkerStatusResponseMsg)(nil),          // 5: proto.WorkerStatusResponseMsg
	nil,                                      // 6: proto.StartTaskMsg.VariablesEntry
	(*any.Any)(nil),                          // 7: google.protobuf.Any
	(*empty.Empty)(nil),                      // 8: google.protobuf.Empty
}
var file_wservices_proto_depIdxs = []int32{
	1, // 0: proto.StartTaskMsg.taskID:type_name -> proto.TaskIdMsg
	6, // 1: proto.StartTaskMsg.variables:type_name -> proto.StartTaskMsg.VariablesEntry
	7, // 2: proto.StartTaskMsg.Command:type_name -> google.protobuf.Any
	0, // 3: proto.TaskExecutionResponseMsg.status:type_name -> proto.TaskExecutionResponseMsg.TaskStatus
	2, // 4: proto.TaskExecutionService.StartTask:input_type -> proto.StartTaskMsg
	1, // 5: proto.TaskExecutionService.TerminateTask:input_type -> proto.TaskIdMsg
	1, // 6: proto.TaskExecutionService.CompleteTask:input_type -> proto.TaskIdMsg
	1, // 7: proto.TaskExecutionService.TaskStatus:input_type -> proto.TaskIdMsg
	8, // 8: proto.TaskExecutionService.WorkerStatus:input_type -> google.protobuf.Empty
	3, // 9: proto.TaskExecutionService.StartTask:output_type -> proto.TaskExecutionResponseMsg
	4, // 10: proto.TaskExecutionService.TerminateTask:output_type -> proto.WorkerActionMsg
	4, // 11: proto.TaskExecutionService.CompleteTask:output_type -> proto.WorkerActionMsg
	3, // 12: proto.TaskExecutionService.TaskStatus:output_type -> proto.TaskExecutionResponseMsg
	5, // 13: proto.TaskExecutionService.WorkerStatus:output_type -> proto.WorkerStatusResponseMsg
	9, // [9:14] is the sub-list for method output_type
	4, // [4:9] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_wservices_proto_init() }
func file_wservices_proto_init() {
	if File_wservices_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_wservices_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskIdMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wservices_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartTaskMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wservices_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TaskExecutionResponseMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wservices_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkerActionMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_wservices_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WorkerStatusResponseMsg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_wservices_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_wservices_proto_goTypes,
		DependencyIndexes: file_wservices_proto_depIdxs,
		EnumInfos:         file_wservices_proto_enumTypes,
		MessageInfos:      file_wservices_proto_msgTypes,
	}.Build()
	File_wservices_proto = out.File
	file_wservices_proto_rawDesc = nil
	file_wservices_proto_goTypes = nil
	file_wservices_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// TaskExecutionServiceClient is the client API for TaskExecutionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TaskExecutionServiceClient interface {
	StartTask(ctx context.Context, in *StartTaskMsg, opts ...grpc.CallOption) (*TaskExecutionResponseMsg, error)
	TerminateTask(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (*WorkerActionMsg, error)
	CompleteTask(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (*WorkerActionMsg, error)
	TaskStatus(ctx context.Context, in *TaskIdMsg, opts ...grpc.CallOption) (*TaskExecutionResponseMsg, error)
	WorkerStatus(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*WorkerStatusResponseMsg, error)
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

func (c *taskExecutionServiceClient) WorkerStatus(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*WorkerStatusResponseMsg, error) {
	out := new(WorkerStatusResponseMsg)
	err := c.cc.Invoke(ctx, "/proto.TaskExecutionService/WorkerStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TaskExecutionServiceServer is the server API for TaskExecutionService service.
type TaskExecutionServiceServer interface {
	StartTask(context.Context, *StartTaskMsg) (*TaskExecutionResponseMsg, error)
	TerminateTask(context.Context, *TaskIdMsg) (*WorkerActionMsg, error)
	CompleteTask(context.Context, *TaskIdMsg) (*WorkerActionMsg, error)
	TaskStatus(context.Context, *TaskIdMsg) (*TaskExecutionResponseMsg, error)
	WorkerStatus(context.Context, *empty.Empty) (*WorkerStatusResponseMsg, error)
}

// UnimplementedTaskExecutionServiceServer can be embedded to have forward compatible implementations.
type UnimplementedTaskExecutionServiceServer struct {
}

func (*UnimplementedTaskExecutionServiceServer) StartTask(context.Context, *StartTaskMsg) (*TaskExecutionResponseMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartTask not implemented")
}
func (*UnimplementedTaskExecutionServiceServer) TerminateTask(context.Context, *TaskIdMsg) (*WorkerActionMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TerminateTask not implemented")
}
func (*UnimplementedTaskExecutionServiceServer) CompleteTask(context.Context, *TaskIdMsg) (*WorkerActionMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CompleteTask not implemented")
}
func (*UnimplementedTaskExecutionServiceServer) TaskStatus(context.Context, *TaskIdMsg) (*TaskExecutionResponseMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TaskStatus not implemented")
}
func (*UnimplementedTaskExecutionServiceServer) WorkerStatus(context.Context, *empty.Empty) (*WorkerStatusResponseMsg, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WorkerStatus not implemented")
}

func RegisterTaskExecutionServiceServer(s *grpc.Server, srv TaskExecutionServiceServer) {
	s.RegisterService(&_TaskExecutionService_serviceDesc, srv)
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

func _TaskExecutionService_WorkerStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
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
		return srv.(TaskExecutionServiceServer).WorkerStatus(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _TaskExecutionService_serviceDesc = grpc.ServiceDesc{
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
	Streams:  []grpc.StreamDesc{},
	Metadata: "wservices.proto",
}
