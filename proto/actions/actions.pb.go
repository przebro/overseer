// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.20.1
// source: actions.proto

package actions

import (
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

type AwsTaskAction_AwsType int32

const (
	AwsTaskAction_lambda   AwsTaskAction_AwsType = 0
	AwsTaskAction_stepfunc AwsTaskAction_AwsType = 1
	AwsTaskAction_batch    AwsTaskAction_AwsType = 2
)

// Enum value maps for AwsTaskAction_AwsType.
var (
	AwsTaskAction_AwsType_name = map[int32]string{
		0: "lambda",
		1: "stepfunc",
		2: "batch",
	}
	AwsTaskAction_AwsType_value = map[string]int32{
		"lambda":   0,
		"stepfunc": 1,
		"batch":    2,
	}
)

func (x AwsTaskAction_AwsType) Enum() *AwsTaskAction_AwsType {
	p := new(AwsTaskAction_AwsType)
	*p = x
	return p
}

func (x AwsTaskAction_AwsType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (AwsTaskAction_AwsType) Descriptor() protoreflect.EnumDescriptor {
	return file_actions_proto_enumTypes[0].Descriptor()
}

func (AwsTaskAction_AwsType) Type() protoreflect.EnumType {
	return &file_actions_proto_enumTypes[0]
}

func (x AwsTaskAction_AwsType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use AwsTaskAction_AwsType.Descriptor instead.
func (AwsTaskAction_AwsType) EnumDescriptor() ([]byte, []int) {
	return file_actions_proto_rawDescGZIP(), []int{6, 0}
}

type DummyTaskAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data string `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *DummyTaskAction) Reset() {
	*x = DummyTaskAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_actions_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DummyTaskAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DummyTaskAction) ProtoMessage() {}

func (x *DummyTaskAction) ProtoReflect() protoreflect.Message {
	mi := &file_actions_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DummyTaskAction.ProtoReflect.Descriptor instead.
func (*DummyTaskAction) Descriptor() ([]byte, []int) {
	return file_actions_proto_rawDescGZIP(), []int{0}
}

func (x *DummyTaskAction) GetData() string {
	if x != nil {
		return x.Data
	}
	return ""
}

type OsStepDefinition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StepName string `protobuf:"bytes,1,opt,name=stepName,proto3" json:"stepName,omitempty"`
	Command  string `protobuf:"bytes,2,opt,name=command,proto3" json:"command,omitempty"`
}

func (x *OsStepDefinition) Reset() {
	*x = OsStepDefinition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_actions_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OsStepDefinition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OsStepDefinition) ProtoMessage() {}

func (x *OsStepDefinition) ProtoReflect() protoreflect.Message {
	mi := &file_actions_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OsStepDefinition.ProtoReflect.Descriptor instead.
func (*OsStepDefinition) Descriptor() ([]byte, []int) {
	return file_actions_proto_rawDescGZIP(), []int{1}
}

func (x *OsStepDefinition) GetStepName() string {
	if x != nil {
		return x.StepName
	}
	return ""
}

func (x *OsStepDefinition) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

type OsTaskAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CommandLine string              `protobuf:"bytes,2,opt,name=commandLine,proto3" json:"commandLine,omitempty"`
	Runas       string              `protobuf:"bytes,3,opt,name=runas,proto3" json:"runas,omitempty"`
	Steps       []*OsStepDefinition `protobuf:"bytes,4,rep,name=steps,proto3" json:"steps,omitempty"`
}

func (x *OsTaskAction) Reset() {
	*x = OsTaskAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_actions_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OsTaskAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OsTaskAction) ProtoMessage() {}

func (x *OsTaskAction) ProtoReflect() protoreflect.Message {
	mi := &file_actions_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OsTaskAction.ProtoReflect.Descriptor instead.
func (*OsTaskAction) Descriptor() ([]byte, []int) {
	return file_actions_proto_rawDescGZIP(), []int{2}
}

func (x *OsTaskAction) GetCommandLine() string {
	if x != nil {
		return x.CommandLine
	}
	return ""
}

func (x *OsTaskAction) GetRunas() string {
	if x != nil {
		return x.Runas
	}
	return ""
}

func (x *OsTaskAction) GetSteps() []*OsStepDefinition {
	if x != nil {
		return x.Steps
	}
	return nil
}

type AwsLambdaExecution struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FunctionName string `protobuf:"bytes,1,opt,name=functionName,proto3" json:"functionName,omitempty"`
	Alias        string `protobuf:"bytes,2,opt,name=alias,proto3" json:"alias,omitempty"`
}

func (x *AwsLambdaExecution) Reset() {
	*x = AwsLambdaExecution{}
	if protoimpl.UnsafeEnabled {
		mi := &file_actions_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AwsLambdaExecution) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AwsLambdaExecution) ProtoMessage() {}

func (x *AwsLambdaExecution) ProtoReflect() protoreflect.Message {
	mi := &file_actions_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AwsLambdaExecution.ProtoReflect.Descriptor instead.
func (*AwsLambdaExecution) Descriptor() ([]byte, []int) {
	return file_actions_proto_rawDescGZIP(), []int{3}
}

func (x *AwsLambdaExecution) GetFunctionName() string {
	if x != nil {
		return x.FunctionName
	}
	return ""
}

func (x *AwsLambdaExecution) GetAlias() string {
	if x != nil {
		return x.Alias
	}
	return ""
}

type AwsStepFunctionExecution struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StateMachineARN string `protobuf:"bytes,1,opt,name=stateMachineARN,proto3" json:"stateMachineARN,omitempty"`
	ExecutionName   string `protobuf:"bytes,2,opt,name=executionName,proto3" json:"executionName,omitempty"`
}

func (x *AwsStepFunctionExecution) Reset() {
	*x = AwsStepFunctionExecution{}
	if protoimpl.UnsafeEnabled {
		mi := &file_actions_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AwsStepFunctionExecution) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AwsStepFunctionExecution) ProtoMessage() {}

func (x *AwsStepFunctionExecution) ProtoReflect() protoreflect.Message {
	mi := &file_actions_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AwsStepFunctionExecution.ProtoReflect.Descriptor instead.
func (*AwsStepFunctionExecution) Descriptor() ([]byte, []int) {
	return file_actions_proto_rawDescGZIP(), []int{4}
}

func (x *AwsStepFunctionExecution) GetStateMachineARN() string {
	if x != nil {
		return x.StateMachineARN
	}
	return ""
}

func (x *AwsStepFunctionExecution) GetExecutionName() string {
	if x != nil {
		return x.ExecutionName
	}
	return ""
}

type AwsConnetionData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ProfileName string `protobuf:"bytes,3,opt,name=profileName,proto3" json:"profileName,omitempty"`
	Region      string `protobuf:"bytes,4,opt,name=region,proto3" json:"region,omitempty"`
}

func (x *AwsConnetionData) Reset() {
	*x = AwsConnetionData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_actions_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AwsConnetionData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AwsConnetionData) ProtoMessage() {}

func (x *AwsConnetionData) ProtoReflect() protoreflect.Message {
	mi := &file_actions_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AwsConnetionData.ProtoReflect.Descriptor instead.
func (*AwsConnetionData) Descriptor() ([]byte, []int) {
	return file_actions_proto_rawDescGZIP(), []int{5}
}

func (x *AwsConnetionData) GetProfileName() string {
	if x != nil {
		return x.ProfileName
	}
	return ""
}

func (x *AwsConnetionData) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

type AwsTaskAction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type AwsTaskAction_AwsType `protobuf:"varint,1,opt,name=type,proto3,enum=proto.AwsTaskAction_AwsType" json:"type,omitempty"`
	// Types that are assignable to Connection:
	//
	//	*AwsTaskAction_ConnectionData
	//	*AwsTaskAction_ConnectionProfileName
	Connection isAwsTaskAction_Connection `protobuf_oneof:"connection"`
	// Types that are assignable to Execution:
	//
	//	*AwsTaskAction_LambdaExecution
	//	*AwsTaskAction_StepFunction
	Execution isAwsTaskAction_Execution `protobuf_oneof:"execution"`
	// Types that are assignable to PayloadSource:
	//
	//	*AwsTaskAction_PayloadRaw
	//	*AwsTaskAction_PayloadFilePath
	PayloadSource isAwsTaskAction_PayloadSource `protobuf_oneof:"payloadSource"`
}

func (x *AwsTaskAction) Reset() {
	*x = AwsTaskAction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_actions_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AwsTaskAction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AwsTaskAction) ProtoMessage() {}

func (x *AwsTaskAction) ProtoReflect() protoreflect.Message {
	mi := &file_actions_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AwsTaskAction.ProtoReflect.Descriptor instead.
func (*AwsTaskAction) Descriptor() ([]byte, []int) {
	return file_actions_proto_rawDescGZIP(), []int{6}
}

func (x *AwsTaskAction) GetType() AwsTaskAction_AwsType {
	if x != nil {
		return x.Type
	}
	return AwsTaskAction_lambda
}

func (m *AwsTaskAction) GetConnection() isAwsTaskAction_Connection {
	if m != nil {
		return m.Connection
	}
	return nil
}

func (x *AwsTaskAction) GetConnectionData() *AwsConnetionData {
	if x, ok := x.GetConnection().(*AwsTaskAction_ConnectionData); ok {
		return x.ConnectionData
	}
	return nil
}

func (x *AwsTaskAction) GetConnectionProfileName() string {
	if x, ok := x.GetConnection().(*AwsTaskAction_ConnectionProfileName); ok {
		return x.ConnectionProfileName
	}
	return ""
}

func (m *AwsTaskAction) GetExecution() isAwsTaskAction_Execution {
	if m != nil {
		return m.Execution
	}
	return nil
}

func (x *AwsTaskAction) GetLambdaExecution() *AwsLambdaExecution {
	if x, ok := x.GetExecution().(*AwsTaskAction_LambdaExecution); ok {
		return x.LambdaExecution
	}
	return nil
}

func (x *AwsTaskAction) GetStepFunction() *AwsStepFunctionExecution {
	if x, ok := x.GetExecution().(*AwsTaskAction_StepFunction); ok {
		return x.StepFunction
	}
	return nil
}

func (m *AwsTaskAction) GetPayloadSource() isAwsTaskAction_PayloadSource {
	if m != nil {
		return m.PayloadSource
	}
	return nil
}

func (x *AwsTaskAction) GetPayloadRaw() []byte {
	if x, ok := x.GetPayloadSource().(*AwsTaskAction_PayloadRaw); ok {
		return x.PayloadRaw
	}
	return nil
}

func (x *AwsTaskAction) GetPayloadFilePath() string {
	if x, ok := x.GetPayloadSource().(*AwsTaskAction_PayloadFilePath); ok {
		return x.PayloadFilePath
	}
	return ""
}

type isAwsTaskAction_Connection interface {
	isAwsTaskAction_Connection()
}

type AwsTaskAction_ConnectionData struct {
	ConnectionData *AwsConnetionData `protobuf:"bytes,2,opt,name=connectionData,proto3,oneof"`
}

type AwsTaskAction_ConnectionProfileName struct {
	// reserved for future use. The name of a profile stored in external data store
	ConnectionProfileName string `protobuf:"bytes,3,opt,name=connectionProfileName,proto3,oneof"`
}

func (*AwsTaskAction_ConnectionData) isAwsTaskAction_Connection() {}

func (*AwsTaskAction_ConnectionProfileName) isAwsTaskAction_Connection() {}

type isAwsTaskAction_Execution interface {
	isAwsTaskAction_Execution()
}

type AwsTaskAction_LambdaExecution struct {
	LambdaExecution *AwsLambdaExecution `protobuf:"bytes,4,opt,name=lambdaExecution,proto3,oneof"`
}

type AwsTaskAction_StepFunction struct {
	StepFunction *AwsStepFunctionExecution `protobuf:"bytes,5,opt,name=stepFunction,proto3,oneof"`
}

func (*AwsTaskAction_LambdaExecution) isAwsTaskAction_Execution() {}

func (*AwsTaskAction_StepFunction) isAwsTaskAction_Execution() {}

type isAwsTaskAction_PayloadSource interface {
	isAwsTaskAction_PayloadSource()
}

type AwsTaskAction_PayloadRaw struct {
	PayloadRaw []byte `protobuf:"bytes,6,opt,name=payloadRaw,proto3,oneof"`
}

type AwsTaskAction_PayloadFilePath struct {
	PayloadFilePath string `protobuf:"bytes,7,opt,name=payloadFilePath,proto3,oneof"`
}

func (*AwsTaskAction_PayloadRaw) isAwsTaskAction_PayloadSource() {}

func (*AwsTaskAction_PayloadFilePath) isAwsTaskAction_PayloadSource() {}

var File_actions_proto protoreflect.FileDescriptor

var file_actions_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x25, 0x0a, 0x0f, 0x44, 0x75, 0x6d, 0x6d, 0x79, 0x54,
	0x61, 0x73, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74,
	0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0x48, 0x0a,
	0x10, 0x4f, 0x73, 0x53, 0x74, 0x65, 0x70, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x74, 0x65, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x74, 0x65, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x22, 0x75, 0x0a, 0x0c, 0x4f, 0x73, 0x54, 0x61, 0x73,
	0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x4c, 0x69, 0x6e, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x4c, 0x69, 0x6e, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x72, 0x75, 0x6e,
	0x61, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x72, 0x75, 0x6e, 0x61, 0x73, 0x12,
	0x2d, 0x0a, 0x05, 0x73, 0x74, 0x65, 0x70, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4f, 0x73, 0x53, 0x74, 0x65, 0x70, 0x44, 0x65, 0x66,
	0x69, 0x6e, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x73, 0x74, 0x65, 0x70, 0x73, 0x22, 0x4e,
	0x0a, 0x12, 0x41, 0x77, 0x73, 0x4c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x45, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x66, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x66, 0x75, 0x6e, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x69, 0x61,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x61, 0x6c, 0x69, 0x61, 0x73, 0x22, 0x6a,
	0x0a, 0x18, 0x41, 0x77, 0x73, 0x53, 0x74, 0x65, 0x70, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x0f, 0x73, 0x74,
	0x61, 0x74, 0x65, 0x4d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x41, 0x52, 0x4e, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x4d, 0x61, 0x63, 0x68, 0x69, 0x6e,
	0x65, 0x41, 0x52, 0x4e, 0x12, 0x24, 0x0a, 0x0d, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f,
	0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x65, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x4c, 0x0a, 0x10, 0x41, 0x77,
	0x73, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x12, 0x20,
	0x0a, 0x0b, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0xf4, 0x03, 0x0a, 0x0d, 0x41, 0x77, 0x73,
	0x54, 0x61, 0x73, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x30, 0x0a, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x41, 0x77, 0x73, 0x54, 0x61, 0x73, 0x6b, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x41,
	0x77, 0x73, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x41, 0x0a, 0x0e,
	0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x77, 0x73,
	0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x48, 0x00, 0x52,
	0x0e, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x36, 0x0a, 0x15, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f,
	0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x15, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x45, 0x0a, 0x0f, 0x6c, 0x61, 0x6d, 0x62, 0x64,
	0x61, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x77, 0x73, 0x4c, 0x61, 0x6d, 0x62,
	0x64, 0x61, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x01, 0x52, 0x0f, 0x6c,
	0x61, 0x6d, 0x62, 0x64, 0x61, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x45,
	0x0a, 0x0c, 0x73, 0x74, 0x65, 0x70, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x77, 0x73,
	0x53, 0x74, 0x65, 0x70, 0x46, 0x75, 0x6e, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x45, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x69, 0x6f, 0x6e, 0x48, 0x01, 0x52, 0x0c, 0x73, 0x74, 0x65, 0x70, 0x46, 0x75, 0x6e,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x20, 0x0a, 0x0a, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64,
	0x52, 0x61, 0x77, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x02, 0x52, 0x0a, 0x70, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x52, 0x61, 0x77, 0x12, 0x2a, 0x0a, 0x0f, 0x70, 0x61, 0x79, 0x6c, 0x6f,
	0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x74, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x48, 0x02, 0x52, 0x0f, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x46, 0x69, 0x6c, 0x65, 0x50,
	0x61, 0x74, 0x68, 0x22, 0x2e, 0x0a, 0x07, 0x41, 0x77, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0a,
	0x0a, 0x06, 0x6c, 0x61, 0x6d, 0x62, 0x64, 0x61, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x73, 0x74,
	0x65, 0x70, 0x66, 0x75, 0x6e, 0x63, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x62, 0x61, 0x74, 0x63,
	0x68, 0x10, 0x02, 0x42, 0x0c, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x42, 0x0b, 0x0a, 0x09, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x42, 0x0f,
	0x0a, 0x0d, 0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x42,
	0x0f, 0x5a, 0x0d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_actions_proto_rawDescOnce sync.Once
	file_actions_proto_rawDescData = file_actions_proto_rawDesc
)

func file_actions_proto_rawDescGZIP() []byte {
	file_actions_proto_rawDescOnce.Do(func() {
		file_actions_proto_rawDescData = protoimpl.X.CompressGZIP(file_actions_proto_rawDescData)
	})
	return file_actions_proto_rawDescData
}

var file_actions_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_actions_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_actions_proto_goTypes = []interface{}{
	(AwsTaskAction_AwsType)(0),       // 0: proto.AwsTaskAction.AwsType
	(*DummyTaskAction)(nil),          // 1: proto.DummyTaskAction
	(*OsStepDefinition)(nil),         // 2: proto.OsStepDefinition
	(*OsTaskAction)(nil),             // 3: proto.OsTaskAction
	(*AwsLambdaExecution)(nil),       // 4: proto.AwsLambdaExecution
	(*AwsStepFunctionExecution)(nil), // 5: proto.AwsStepFunctionExecution
	(*AwsConnetionData)(nil),         // 6: proto.AwsConnetionData
	(*AwsTaskAction)(nil),            // 7: proto.AwsTaskAction
}
var file_actions_proto_depIdxs = []int32{
	2, // 0: proto.OsTaskAction.steps:type_name -> proto.OsStepDefinition
	0, // 1: proto.AwsTaskAction.type:type_name -> proto.AwsTaskAction.AwsType
	6, // 2: proto.AwsTaskAction.connectionData:type_name -> proto.AwsConnetionData
	4, // 3: proto.AwsTaskAction.lambdaExecution:type_name -> proto.AwsLambdaExecution
	5, // 4: proto.AwsTaskAction.stepFunction:type_name -> proto.AwsStepFunctionExecution
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_actions_proto_init() }
func file_actions_proto_init() {
	if File_actions_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_actions_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DummyTaskAction); i {
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
		file_actions_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OsStepDefinition); i {
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
		file_actions_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OsTaskAction); i {
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
		file_actions_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AwsLambdaExecution); i {
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
		file_actions_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AwsStepFunctionExecution); i {
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
		file_actions_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AwsConnetionData); i {
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
		file_actions_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AwsTaskAction); i {
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
	file_actions_proto_msgTypes[6].OneofWrappers = []interface{}{
		(*AwsTaskAction_ConnectionData)(nil),
		(*AwsTaskAction_ConnectionProfileName)(nil),
		(*AwsTaskAction_LambdaExecution)(nil),
		(*AwsTaskAction_StepFunction)(nil),
		(*AwsTaskAction_PayloadRaw)(nil),
		(*AwsTaskAction_PayloadFilePath)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_actions_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_actions_proto_goTypes,
		DependencyIndexes: file_actions_proto_depIdxs,
		EnumInfos:         file_actions_proto_enumTypes,
		MessageInfos:      file_actions_proto_msgTypes,
	}.Build()
	File_actions_proto = out.File
	file_actions_proto_rawDesc = nil
	file_actions_proto_goTypes = nil
	file_actions_proto_depIdxs = nil
}
