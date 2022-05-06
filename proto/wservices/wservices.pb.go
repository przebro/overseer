// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.20.1
// source: wservices.proto

package wservices

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TaskExecutionResponseMsg_TaskStatus int32

const (
	TaskExecutionResponseMsg_RECEIVED  TaskExecutionResponseMsg_TaskStatus = 0
	TaskExecutionResponseMsg_EXECUTING TaskExecutionResponseMsg_TaskStatus = 1
	TaskExecutionResponseMsg_ENDED     TaskExecutionResponseMsg_TaskStatus = 2
	TaskExecutionResponseMsg_FAILED    TaskExecutionResponseMsg_TaskStatus = 3
	//reserved for future use
	TaskExecutionResponseMsg_WAITING TaskExecutionResponseMsg_TaskStatus = 4
	//reserved for future use
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

	TaskID      string `protobuf:"bytes,1,opt,name=taskID,proto3" json:"taskID,omitempty"`
	ExecutionID string `protobuf:"bytes,2,opt,name=ExecutionID,proto3" json:"ExecutionID,omitempty"`
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

func (x *TaskIdMsg) GetExecutionID() string {
	if x != nil {
		return x.ExecutionID
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
	Command   *anypb.Any        `protobuf:"bytes,5,opt,name=Command,proto3" json:"Command,omitempty"`
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

func (x *StartTaskMsg) GetCommand() *anypb.Any {
	if x != nil {
		return x.Command
	}
	return nil
}

type TaskExecutionResponseMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status TaskExecutionResponseMsg_TaskStatus `protobuf:"varint,1,opt,name=status,proto3,enum=proto.TaskExecutionResponseMsg_TaskStatus" json:"status,omitempty"`
	//Code returned from executed program,
	//For dummy worker it is always set to 0
	//For os type, it contains return code of an executed command
	//For AWS Lambda it contains response code, usually 200 which means that lambda function was started successfully,
	//even if later execution was unsuccessful
	ReturnCode int32 `protobuf:"varint,2,opt,name=returnCode,proto3" json:"returnCode,omitempty"`
	//This field contains a subjective status of a task.
	StatusCode int32  `protobuf:"varint,3,opt,name=statusCode,proto3" json:"statusCode,omitempty"`
	Reason     string `protobuf:"bytes,4,opt,name=reason,proto3" json:"reason,omitempty"`
	Pid        int32  `protobuf:"varint,5,opt,name=pid,proto3" json:"pid,omitempty"`
	Tasks      int32  `protobuf:"varint,6,opt,name=tasks,proto3" json:"tasks,omitempty"`
	TasksLimit int32  `protobuf:"varint,7,opt,name=tasksLimit,proto3" json:"tasksLimit,omitempty"`
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

func (x *TaskExecutionResponseMsg) GetReason() string {
	if x != nil {
		return x.Reason
	}
	return ""
}

func (x *TaskExecutionResponseMsg) GetPid() int32 {
	if x != nil {
		return x.Pid
	}
	return 0
}

func (x *TaskExecutionResponseMsg) GetTasks() int32 {
	if x != nil {
		return x.Tasks
	}
	return 0
}

func (x *TaskExecutionResponseMsg) GetTasksLimit() int32 {
	if x != nil {
		return x.TasksLimit
	}
	return 0
}

type WorkerActionMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success    bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Message    string `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Tasks      int32  `protobuf:"varint,3,opt,name=tasks,proto3" json:"tasks,omitempty"`
	TasksLimit int32  `protobuf:"varint,4,opt,name=tasksLimit,proto3" json:"tasksLimit,omitempty"`
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

func (x *WorkerActionMsg) GetTasks() int32 {
	if x != nil {
		return x.Tasks
	}
	return 0
}

func (x *WorkerActionMsg) GetTasksLimit() int32 {
	if x != nil {
		return x.TasksLimit
	}
	return 0
}

type TaskOutputMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *TaskOutputMsg) Reset() {
	*x = TaskOutputMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wservices_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TaskOutputMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TaskOutputMsg) ProtoMessage() {}

func (x *TaskOutputMsg) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use TaskOutputMsg.ProtoReflect.Descriptor instead.
func (*TaskOutputMsg) Descriptor() ([]byte, []int) {
	return file_wservices_proto_rawDescGZIP(), []int{4}
}

func (x *TaskOutputMsg) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type WorkerStatusResponseMsg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Tasks      int32 `protobuf:"varint,1,opt,name=tasks,proto3" json:"tasks,omitempty"`
	TasksLimit int32 `protobuf:"varint,2,opt,name=tasksLimit,proto3" json:"tasksLimit,omitempty"`
	Cpuload    int32 `protobuf:"varint,3,opt,name=cpuload,proto3" json:"cpuload,omitempty"`
	Memused    int32 `protobuf:"varint,4,opt,name=memused,proto3" json:"memused,omitempty"`
	Memtotal   int32 `protobuf:"varint,5,opt,name=memtotal,proto3" json:"memtotal,omitempty"`
}

func (x *WorkerStatusResponseMsg) Reset() {
	*x = WorkerStatusResponseMsg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_wservices_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WorkerStatusResponseMsg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WorkerStatusResponseMsg) ProtoMessage() {}

func (x *WorkerStatusResponseMsg) ProtoReflect() protoreflect.Message {
	mi := &file_wservices_proto_msgTypes[5]
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
	return file_wservices_proto_rawDescGZIP(), []int{5}
}

func (x *WorkerStatusResponseMsg) GetTasks() int32 {
	if x != nil {
		return x.Tasks
	}
	return 0
}

func (x *WorkerStatusResponseMsg) GetTasksLimit() int32 {
	if x != nil {
		return x.TasksLimit
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
	0x22, 0x45, 0x0a, 0x09, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x4d, 0x73, 0x67, 0x12, 0x16, 0x0a,
	0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74,
	0x61, 0x73, 0x6b, 0x49, 0x44, 0x12, 0x20, 0x0a, 0x0b, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x45, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x44, 0x22, 0xfc, 0x01, 0x0a, 0x0c, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x54, 0x61, 0x73, 0x6b, 0x4d, 0x73, 0x67, 0x12, 0x28, 0x0a, 0x06, 0x74, 0x61, 0x73, 0x6b,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x4d, 0x73, 0x67, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b,
	0x49, 0x44, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x40, 0x0a, 0x09, 0x76, 0x61, 0x72, 0x69, 0x61, 0x62,
	0x6c, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x4d, 0x73, 0x67, 0x2e, 0x56,
	0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x09, 0x76,
	0x61, 0x72, 0x69, 0x61, 0x62, 0x6c, 0x65, 0x73, 0x12, 0x2e, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x52,
	0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x1a, 0x3c, 0x0a, 0x0e, 0x56, 0x61, 0x72, 0x69,
	0x61, 0x62, 0x6c, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65,
	0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0xe5, 0x02, 0x0a, 0x18, 0x54, 0x61, 0x73, 0x6b, 0x45,
	0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x4d, 0x73, 0x67, 0x12, 0x42, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x2a, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b,
	0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x4d, 0x73, 0x67, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x72, 0x65, 0x74, 0x75, 0x72,
	0x6e, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x72, 0x65, 0x74,
	0x75, 0x72, 0x6e, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x61, 0x73, 0x6f, 0x6e, 0x12,
	0x10, 0x0a, 0x03, 0x70, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x70, 0x69,
	0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x61, 0x73, 0x6b, 0x73,
	0x4c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x61, 0x73,
	0x6b, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x65, 0x0a, 0x0a, 0x54, 0x61, 0x73, 0x6b, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x43, 0x45, 0x49, 0x56, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x58, 0x45, 0x43, 0x55, 0x54, 0x49, 0x4e, 0x47,
	0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x45, 0x4e, 0x44, 0x45, 0x44, 0x10, 0x02, 0x12, 0x0a, 0x0a,
	0x06, 0x46, 0x41, 0x49, 0x4c, 0x45, 0x44, 0x10, 0x03, 0x12, 0x0b, 0x0a, 0x07, 0x57, 0x41, 0x49,
	0x54, 0x49, 0x4e, 0x47, 0x10, 0x04, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x44, 0x4c, 0x45, 0x10, 0x05,
	0x12, 0x0c, 0x0a, 0x08, 0x53, 0x54, 0x41, 0x52, 0x54, 0x49, 0x4e, 0x47, 0x10, 0x06, 0x22, 0x7b,
	0x0a, 0x0f, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x73,
	0x67, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x74,
	0x61, 0x73, 0x6b, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x0a, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x22, 0x23, 0x0a, 0x0d, 0x54,
	0x61, 0x73, 0x6b, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x4d, 0x73, 0x67, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61,
	0x22, 0x9f, 0x01, 0x0a, 0x17, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x73, 0x67, 0x12, 0x14, 0x0a, 0x05,
	0x74, 0x61, 0x73, 0x6b, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x61, 0x73,
	0x6b, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x4c, 0x69, 0x6d, 0x69, 0x74,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x74, 0x61, 0x73, 0x6b, 0x73, 0x4c, 0x69, 0x6d,
	0x69, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x70, 0x75, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x07, 0x63, 0x70, 0x75, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x6d, 0x75, 0x73, 0x65, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x6d,
	0x65, 0x6d, 0x75, 0x73, 0x65, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x6d, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6d, 0x65, 0x6d, 0x74, 0x6f, 0x74,
	0x61, 0x6c, 0x32, 0x9b, 0x03, 0x0a, 0x14, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x09, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x54, 0x61, 0x73, 0x6b, 0x4d, 0x73, 0x67, 0x1a, 0x1f, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x73, 0x67, 0x22, 0x00,
	0x12, 0x3b, 0x0a, 0x0d, 0x54, 0x65, 0x72, 0x6d, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x54, 0x61, 0x73,
	0x6b, 0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64,
	0x4d, 0x73, 0x67, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x6f, 0x72, 0x6b,
	0x65, 0x72, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x12, 0x3a, 0x0a,
	0x0c, 0x43, 0x6f, 0x6d, 0x70, 0x6c, 0x65, 0x74, 0x65, 0x54, 0x61, 0x73, 0x6b, 0x12, 0x10, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x4d, 0x73, 0x67, 0x1a,
	0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x12, 0x41, 0x0a, 0x0a, 0x54, 0x61, 0x73,
	0x6b, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x4d, 0x73, 0x67, 0x1a, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x73, 0x67, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x0a,
	0x54, 0x61, 0x73, 0x6b, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x10, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x4d, 0x73, 0x67, 0x1a, 0x14, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x4d,
	0x73, 0x67, 0x22, 0x00, 0x30, 0x01, 0x12, 0x48, 0x0a, 0x0c, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x57, 0x6f, 0x72, 0x6b, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x73, 0x67, 0x22, 0x00,
	0x42, 0x11, 0x5a, 0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x77, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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
var file_wservices_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_wservices_proto_goTypes = []interface{}{
	(TaskExecutionResponseMsg_TaskStatus)(0), // 0: proto.TaskExecutionResponseMsg.TaskStatus
	(*TaskIdMsg)(nil),                        // 1: proto.TaskIdMsg
	(*StartTaskMsg)(nil),                     // 2: proto.StartTaskMsg
	(*TaskExecutionResponseMsg)(nil),         // 3: proto.TaskExecutionResponseMsg
	(*WorkerActionMsg)(nil),                  // 4: proto.WorkerActionMsg
	(*TaskOutputMsg)(nil),                    // 5: proto.TaskOutputMsg
	(*WorkerStatusResponseMsg)(nil),          // 6: proto.WorkerStatusResponseMsg
	nil,                                      // 7: proto.StartTaskMsg.VariablesEntry
	(*anypb.Any)(nil),                        // 8: google.protobuf.Any
	(*emptypb.Empty)(nil),                    // 9: google.protobuf.Empty
}
var file_wservices_proto_depIdxs = []int32{
	1,  // 0: proto.StartTaskMsg.taskID:type_name -> proto.TaskIdMsg
	7,  // 1: proto.StartTaskMsg.variables:type_name -> proto.StartTaskMsg.VariablesEntry
	8,  // 2: proto.StartTaskMsg.Command:type_name -> google.protobuf.Any
	0,  // 3: proto.TaskExecutionResponseMsg.status:type_name -> proto.TaskExecutionResponseMsg.TaskStatus
	2,  // 4: proto.TaskExecutionService.StartTask:input_type -> proto.StartTaskMsg
	1,  // 5: proto.TaskExecutionService.TerminateTask:input_type -> proto.TaskIdMsg
	1,  // 6: proto.TaskExecutionService.CompleteTask:input_type -> proto.TaskIdMsg
	1,  // 7: proto.TaskExecutionService.TaskStatus:input_type -> proto.TaskIdMsg
	1,  // 8: proto.TaskExecutionService.TaskOutput:input_type -> proto.TaskIdMsg
	9,  // 9: proto.TaskExecutionService.WorkerStatus:input_type -> google.protobuf.Empty
	3,  // 10: proto.TaskExecutionService.StartTask:output_type -> proto.TaskExecutionResponseMsg
	4,  // 11: proto.TaskExecutionService.TerminateTask:output_type -> proto.WorkerActionMsg
	4,  // 12: proto.TaskExecutionService.CompleteTask:output_type -> proto.WorkerActionMsg
	3,  // 13: proto.TaskExecutionService.TaskStatus:output_type -> proto.TaskExecutionResponseMsg
	5,  // 14: proto.TaskExecutionService.TaskOutput:output_type -> proto.TaskOutputMsg
	6,  // 15: proto.TaskExecutionService.WorkerStatus:output_type -> proto.WorkerStatusResponseMsg
	10, // [10:16] is the sub-list for method output_type
	4,  // [4:10] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
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
			switch v := v.(*TaskOutputMsg); i {
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
		file_wservices_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
			NumMessages:   7,
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
