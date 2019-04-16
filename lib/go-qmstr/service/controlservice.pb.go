// Code generated by protoc-gen-go. DO NOT EDIT.
// source: controlservice.proto

package service

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type LogMessage struct {
	Msg                  []byte   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogMessage) Reset()         { *m = LogMessage{} }
func (m *LogMessage) String() string { return proto.CompactTextString(m) }
func (*LogMessage) ProtoMessage()    {}
func (*LogMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{0}
}
func (m *LogMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogMessage.Unmarshal(m, b)
}
func (m *LogMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogMessage.Marshal(b, m, deterministic)
}
func (dst *LogMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogMessage.Merge(dst, src)
}
func (m *LogMessage) XXX_Size() int {
	return xxx_messageInfo_LogMessage.Size(m)
}
func (m *LogMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_LogMessage.DiscardUnknown(m)
}

var xxx_messageInfo_LogMessage proto.InternalMessageInfo

func (m *LogMessage) GetMsg() []byte {
	if m != nil {
		return m.Msg
	}
	return nil
}

type LogResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LogResponse) Reset()         { *m = LogResponse{} }
func (m *LogResponse) String() string { return proto.CompactTextString(m) }
func (*LogResponse) ProtoMessage()    {}
func (*LogResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{1}
}
func (m *LogResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LogResponse.Unmarshal(m, b)
}
func (m *LogResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LogResponse.Marshal(b, m, deterministic)
}
func (dst *LogResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LogResponse.Merge(dst, src)
}
func (m *LogResponse) XXX_Size() int {
	return xxx_messageInfo_LogResponse.Size(m)
}
func (m *LogResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LogResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LogResponse proto.InternalMessageInfo

func (m *LogResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type QuitMessage struct {
	Kill                 bool     `protobuf:"varint,1,opt,name=kill,proto3" json:"kill,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QuitMessage) Reset()         { *m = QuitMessage{} }
func (m *QuitMessage) String() string { return proto.CompactTextString(m) }
func (*QuitMessage) ProtoMessage()    {}
func (*QuitMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{2}
}
func (m *QuitMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QuitMessage.Unmarshal(m, b)
}
func (m *QuitMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QuitMessage.Marshal(b, m, deterministic)
}
func (dst *QuitMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuitMessage.Merge(dst, src)
}
func (m *QuitMessage) XXX_Size() int {
	return xxx_messageInfo_QuitMessage.Size(m)
}
func (m *QuitMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_QuitMessage.DiscardUnknown(m)
}

var xxx_messageInfo_QuitMessage proto.InternalMessageInfo

func (m *QuitMessage) GetKill() bool {
	if m != nil {
		return m.Kill
	}
	return false
}

type QuitResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QuitResponse) Reset()         { *m = QuitResponse{} }
func (m *QuitResponse) String() string { return proto.CompactTextString(m) }
func (*QuitResponse) ProtoMessage()    {}
func (*QuitResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{3}
}
func (m *QuitResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QuitResponse.Unmarshal(m, b)
}
func (m *QuitResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QuitResponse.Marshal(b, m, deterministic)
}
func (dst *QuitResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QuitResponse.Merge(dst, src)
}
func (m *QuitResponse) XXX_Size() int {
	return xxx_messageInfo_QuitResponse.Size(m)
}
func (m *QuitResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QuitResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QuitResponse proto.InternalMessageInfo

func (m *QuitResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

type SwitchPhaseMessage struct {
	Phase                Phase    `protobuf:"varint,1,opt,name=phase,proto3,enum=service.Phase" json:"phase,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SwitchPhaseMessage) Reset()         { *m = SwitchPhaseMessage{} }
func (m *SwitchPhaseMessage) String() string { return proto.CompactTextString(m) }
func (*SwitchPhaseMessage) ProtoMessage()    {}
func (*SwitchPhaseMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{4}
}
func (m *SwitchPhaseMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchPhaseMessage.Unmarshal(m, b)
}
func (m *SwitchPhaseMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchPhaseMessage.Marshal(b, m, deterministic)
}
func (dst *SwitchPhaseMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchPhaseMessage.Merge(dst, src)
}
func (m *SwitchPhaseMessage) XXX_Size() int {
	return xxx_messageInfo_SwitchPhaseMessage.Size(m)
}
func (m *SwitchPhaseMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchPhaseMessage.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchPhaseMessage proto.InternalMessageInfo

func (m *SwitchPhaseMessage) GetPhase() Phase {
	if m != nil {
		return m.Phase
	}
	return Phase_INIT
}

type SwitchPhaseResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SwitchPhaseResponse) Reset()         { *m = SwitchPhaseResponse{} }
func (m *SwitchPhaseResponse) String() string { return proto.CompactTextString(m) }
func (*SwitchPhaseResponse) ProtoMessage()    {}
func (*SwitchPhaseResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{5}
}
func (m *SwitchPhaseResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SwitchPhaseResponse.Unmarshal(m, b)
}
func (m *SwitchPhaseResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SwitchPhaseResponse.Marshal(b, m, deterministic)
}
func (dst *SwitchPhaseResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SwitchPhaseResponse.Merge(dst, src)
}
func (m *SwitchPhaseResponse) XXX_Size() int {
	return xxx_messageInfo_SwitchPhaseResponse.Size(m)
}
func (m *SwitchPhaseResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_SwitchPhaseResponse.DiscardUnknown(m)
}

var xxx_messageInfo_SwitchPhaseResponse proto.InternalMessageInfo

func (m *SwitchPhaseResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func (m *SwitchPhaseResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type GetFileNodeMessage struct {
	FileNode             *FileNode `protobuf:"bytes,1,opt,name=fileNode,proto3" json:"fileNode,omitempty"`
	UniqueNode           bool      `protobuf:"varint,2,opt,name=uniqueNode,proto3" json:"uniqueNode,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *GetFileNodeMessage) Reset()         { *m = GetFileNodeMessage{} }
func (m *GetFileNodeMessage) String() string { return proto.CompactTextString(m) }
func (*GetFileNodeMessage) ProtoMessage()    {}
func (*GetFileNodeMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{6}
}
func (m *GetFileNodeMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFileNodeMessage.Unmarshal(m, b)
}
func (m *GetFileNodeMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFileNodeMessage.Marshal(b, m, deterministic)
}
func (dst *GetFileNodeMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFileNodeMessage.Merge(dst, src)
}
func (m *GetFileNodeMessage) XXX_Size() int {
	return xxx_messageInfo_GetFileNodeMessage.Size(m)
}
func (m *GetFileNodeMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFileNodeMessage.DiscardUnknown(m)
}

var xxx_messageInfo_GetFileNodeMessage proto.InternalMessageInfo

func (m *GetFileNodeMessage) GetFileNode() *FileNode {
	if m != nil {
		return m.FileNode
	}
	return nil
}

func (m *GetFileNodeMessage) GetUniqueNode() bool {
	if m != nil {
		return m.UniqueNode
	}
	return false
}

type StatusMessage struct {
	Phase                bool     `protobuf:"varint,1,opt,name=phase,proto3" json:"phase,omitempty"`
	Switch               bool     `protobuf:"varint,2,opt,name=switch,proto3" json:"switch,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusMessage) Reset()         { *m = StatusMessage{} }
func (m *StatusMessage) String() string { return proto.CompactTextString(m) }
func (*StatusMessage) ProtoMessage()    {}
func (*StatusMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{7}
}
func (m *StatusMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusMessage.Unmarshal(m, b)
}
func (m *StatusMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusMessage.Marshal(b, m, deterministic)
}
func (dst *StatusMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusMessage.Merge(dst, src)
}
func (m *StatusMessage) XXX_Size() int {
	return xxx_messageInfo_StatusMessage.Size(m)
}
func (m *StatusMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusMessage.DiscardUnknown(m)
}

var xxx_messageInfo_StatusMessage proto.InternalMessageInfo

func (m *StatusMessage) GetPhase() bool {
	if m != nil {
		return m.Phase
	}
	return false
}

func (m *StatusMessage) GetSwitch() bool {
	if m != nil {
		return m.Switch
	}
	return false
}

type StatusResponse struct {
	Phase                string   `protobuf:"bytes,1,opt,name=phase,proto3" json:"phase,omitempty"`
	PhaseID              Phase    `protobuf:"varint,2,opt,name=phaseID,proto3,enum=service.Phase" json:"phaseID,omitempty"`
	Switching            bool     `protobuf:"varint,3,opt,name=switching,proto3" json:"switching,omitempty"`
	Error                string   `protobuf:"bytes,4,opt,name=error,proto3" json:"error,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusResponse) Reset()         { *m = StatusResponse{} }
func (m *StatusResponse) String() string { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()    {}
func (*StatusResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{8}
}
func (m *StatusResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusResponse.Unmarshal(m, b)
}
func (m *StatusResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusResponse.Marshal(b, m, deterministic)
}
func (dst *StatusResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusResponse.Merge(dst, src)
}
func (m *StatusResponse) XXX_Size() int {
	return xxx_messageInfo_StatusResponse.Size(m)
}
func (m *StatusResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusResponse.DiscardUnknown(m)
}

var xxx_messageInfo_StatusResponse proto.InternalMessageInfo

func (m *StatusResponse) GetPhase() string {
	if m != nil {
		return m.Phase
	}
	return ""
}

func (m *StatusResponse) GetPhaseID() Phase {
	if m != nil {
		return m.PhaseID
	}
	return Phase_INIT
}

func (m *StatusResponse) GetSwitching() bool {
	if m != nil {
		return m.Switching
	}
	return false
}

func (m *StatusResponse) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

type EventMessage struct {
	Class                EventClass `protobuf:"varint,1,opt,name=class,proto3,enum=service.EventClass" json:"class,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *EventMessage) Reset()         { *m = EventMessage{} }
func (m *EventMessage) String() string { return proto.CompactTextString(m) }
func (*EventMessage) ProtoMessage()    {}
func (*EventMessage) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{9}
}
func (m *EventMessage) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventMessage.Unmarshal(m, b)
}
func (m *EventMessage) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventMessage.Marshal(b, m, deterministic)
}
func (dst *EventMessage) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventMessage.Merge(dst, src)
}
func (m *EventMessage) XXX_Size() int {
	return xxx_messageInfo_EventMessage.Size(m)
}
func (m *EventMessage) XXX_DiscardUnknown() {
	xxx_messageInfo_EventMessage.DiscardUnknown(m)
}

var xxx_messageInfo_EventMessage proto.InternalMessageInfo

func (m *EventMessage) GetClass() EventClass {
	if m != nil {
		return m.Class
	}
	return EventClass_ALL
}

type ExportRequest struct {
	Wait                 bool     `protobuf:"varint,1,opt,name=wait,proto3" json:"wait,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExportRequest) Reset()         { *m = ExportRequest{} }
func (m *ExportRequest) String() string { return proto.CompactTextString(m) }
func (*ExportRequest) ProtoMessage()    {}
func (*ExportRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{10}
}
func (m *ExportRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExportRequest.Unmarshal(m, b)
}
func (m *ExportRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExportRequest.Marshal(b, m, deterministic)
}
func (dst *ExportRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExportRequest.Merge(dst, src)
}
func (m *ExportRequest) XXX_Size() int {
	return xxx_messageInfo_ExportRequest.Size(m)
}
func (m *ExportRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ExportRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ExportRequest proto.InternalMessageInfo

func (m *ExportRequest) GetWait() bool {
	if m != nil {
		return m.Wait
	}
	return false
}

type ExportResponse struct {
	Success              bool     `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExportResponse) Reset()         { *m = ExportResponse{} }
func (m *ExportResponse) String() string { return proto.CompactTextString(m) }
func (*ExportResponse) ProtoMessage()    {}
func (*ExportResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_controlservice_788e7d435b159ea4, []int{11}
}
func (m *ExportResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExportResponse.Unmarshal(m, b)
}
func (m *ExportResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExportResponse.Marshal(b, m, deterministic)
}
func (dst *ExportResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExportResponse.Merge(dst, src)
}
func (m *ExportResponse) XXX_Size() int {
	return xxx_messageInfo_ExportResponse.Size(m)
}
func (m *ExportResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ExportResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ExportResponse proto.InternalMessageInfo

func (m *ExportResponse) GetSuccess() bool {
	if m != nil {
		return m.Success
	}
	return false
}

func init() {
	proto.RegisterType((*LogMessage)(nil), "service.LogMessage")
	proto.RegisterType((*LogResponse)(nil), "service.LogResponse")
	proto.RegisterType((*QuitMessage)(nil), "service.QuitMessage")
	proto.RegisterType((*QuitResponse)(nil), "service.QuitResponse")
	proto.RegisterType((*SwitchPhaseMessage)(nil), "service.SwitchPhaseMessage")
	proto.RegisterType((*SwitchPhaseResponse)(nil), "service.SwitchPhaseResponse")
	proto.RegisterType((*GetFileNodeMessage)(nil), "service.GetFileNodeMessage")
	proto.RegisterType((*StatusMessage)(nil), "service.StatusMessage")
	proto.RegisterType((*StatusResponse)(nil), "service.StatusResponse")
	proto.RegisterType((*EventMessage)(nil), "service.EventMessage")
	proto.RegisterType((*ExportRequest)(nil), "service.ExportRequest")
	proto.RegisterType((*ExportResponse)(nil), "service.ExportResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ControlServiceClient is the client API for ControlService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ControlServiceClient interface {
	Log(ctx context.Context, in *LogMessage, opts ...grpc.CallOption) (*LogResponse, error)
	Quit(ctx context.Context, in *QuitMessage, opts ...grpc.CallOption) (*QuitResponse, error)
	SwitchPhase(ctx context.Context, in *SwitchPhaseMessage, opts ...grpc.CallOption) (*SwitchPhaseResponse, error)
	GetPackageNode(ctx context.Context, in *PackageNode, opts ...grpc.CallOption) (*PackageNode, error)
	GetFileNode(ctx context.Context, in *GetFileNodeMessage, opts ...grpc.CallOption) (ControlService_GetFileNodeClient, error)
	GetDiagnosticNode(ctx context.Context, in *DiagnosticNode, opts ...grpc.CallOption) (ControlService_GetDiagnosticNodeClient, error)
	Status(ctx context.Context, in *StatusMessage, opts ...grpc.CallOption) (*StatusResponse, error)
	SubscribeEvents(ctx context.Context, in *EventMessage, opts ...grpc.CallOption) (ControlService_SubscribeEventsClient, error)
	ExportSnapshot(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error)
}

type controlServiceClient struct {
	cc *grpc.ClientConn
}

func NewControlServiceClient(cc *grpc.ClientConn) ControlServiceClient {
	return &controlServiceClient{cc}
}

func (c *controlServiceClient) Log(ctx context.Context, in *LogMessage, opts ...grpc.CallOption) (*LogResponse, error) {
	out := new(LogResponse)
	err := c.cc.Invoke(ctx, "/service.ControlService/Log", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlServiceClient) Quit(ctx context.Context, in *QuitMessage, opts ...grpc.CallOption) (*QuitResponse, error) {
	out := new(QuitResponse)
	err := c.cc.Invoke(ctx, "/service.ControlService/Quit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlServiceClient) SwitchPhase(ctx context.Context, in *SwitchPhaseMessage, opts ...grpc.CallOption) (*SwitchPhaseResponse, error) {
	out := new(SwitchPhaseResponse)
	err := c.cc.Invoke(ctx, "/service.ControlService/SwitchPhase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlServiceClient) GetPackageNode(ctx context.Context, in *PackageNode, opts ...grpc.CallOption) (*PackageNode, error) {
	out := new(PackageNode)
	err := c.cc.Invoke(ctx, "/service.ControlService/GetPackageNode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlServiceClient) GetFileNode(ctx context.Context, in *GetFileNodeMessage, opts ...grpc.CallOption) (ControlService_GetFileNodeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ControlService_serviceDesc.Streams[0], "/service.ControlService/GetFileNode", opts...)
	if err != nil {
		return nil, err
	}
	x := &controlServiceGetFileNodeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ControlService_GetFileNodeClient interface {
	Recv() (*FileNode, error)
	grpc.ClientStream
}

type controlServiceGetFileNodeClient struct {
	grpc.ClientStream
}

func (x *controlServiceGetFileNodeClient) Recv() (*FileNode, error) {
	m := new(FileNode)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *controlServiceClient) GetDiagnosticNode(ctx context.Context, in *DiagnosticNode, opts ...grpc.CallOption) (ControlService_GetDiagnosticNodeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ControlService_serviceDesc.Streams[1], "/service.ControlService/GetDiagnosticNode", opts...)
	if err != nil {
		return nil, err
	}
	x := &controlServiceGetDiagnosticNodeClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ControlService_GetDiagnosticNodeClient interface {
	Recv() (*DiagnosticNode, error)
	grpc.ClientStream
}

type controlServiceGetDiagnosticNodeClient struct {
	grpc.ClientStream
}

func (x *controlServiceGetDiagnosticNodeClient) Recv() (*DiagnosticNode, error) {
	m := new(DiagnosticNode)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *controlServiceClient) Status(ctx context.Context, in *StatusMessage, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/service.ControlService/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *controlServiceClient) SubscribeEvents(ctx context.Context, in *EventMessage, opts ...grpc.CallOption) (ControlService_SubscribeEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &_ControlService_serviceDesc.Streams[2], "/service.ControlService/SubscribeEvents", opts...)
	if err != nil {
		return nil, err
	}
	x := &controlServiceSubscribeEventsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ControlService_SubscribeEventsClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type controlServiceSubscribeEventsClient struct {
	grpc.ClientStream
}

func (x *controlServiceSubscribeEventsClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *controlServiceClient) ExportSnapshot(ctx context.Context, in *ExportRequest, opts ...grpc.CallOption) (*ExportResponse, error) {
	out := new(ExportResponse)
	err := c.cc.Invoke(ctx, "/service.ControlService/ExportSnapshot", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ControlServiceServer is the server API for ControlService service.
type ControlServiceServer interface {
	Log(context.Context, *LogMessage) (*LogResponse, error)
	Quit(context.Context, *QuitMessage) (*QuitResponse, error)
	SwitchPhase(context.Context, *SwitchPhaseMessage) (*SwitchPhaseResponse, error)
	GetPackageNode(context.Context, *PackageNode) (*PackageNode, error)
	GetFileNode(*GetFileNodeMessage, ControlService_GetFileNodeServer) error
	GetDiagnosticNode(*DiagnosticNode, ControlService_GetDiagnosticNodeServer) error
	Status(context.Context, *StatusMessage) (*StatusResponse, error)
	SubscribeEvents(*EventMessage, ControlService_SubscribeEventsServer) error
	ExportSnapshot(context.Context, *ExportRequest) (*ExportResponse, error)
}

func RegisterControlServiceServer(s *grpc.Server, srv ControlServiceServer) {
	s.RegisterService(&_ControlService_serviceDesc, srv)
}

func _ControlService_Log_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServiceServer).Log(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ControlService/Log",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServiceServer).Log(ctx, req.(*LogMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControlService_Quit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QuitMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServiceServer).Quit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ControlService/Quit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServiceServer).Quit(ctx, req.(*QuitMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControlService_SwitchPhase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SwitchPhaseMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServiceServer).SwitchPhase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ControlService/SwitchPhase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServiceServer).SwitchPhase(ctx, req.(*SwitchPhaseMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControlService_GetPackageNode_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PackageNode)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServiceServer).GetPackageNode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ControlService/GetPackageNode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServiceServer).GetPackageNode(ctx, req.(*PackageNode))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControlService_GetFileNode_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetFileNodeMessage)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ControlServiceServer).GetFileNode(m, &controlServiceGetFileNodeServer{stream})
}

type ControlService_GetFileNodeServer interface {
	Send(*FileNode) error
	grpc.ServerStream
}

type controlServiceGetFileNodeServer struct {
	grpc.ServerStream
}

func (x *controlServiceGetFileNodeServer) Send(m *FileNode) error {
	return x.ServerStream.SendMsg(m)
}

func _ControlService_GetDiagnosticNode_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(DiagnosticNode)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ControlServiceServer).GetDiagnosticNode(m, &controlServiceGetDiagnosticNodeServer{stream})
}

type ControlService_GetDiagnosticNodeServer interface {
	Send(*DiagnosticNode) error
	grpc.ServerStream
}

type controlServiceGetDiagnosticNodeServer struct {
	grpc.ServerStream
}

func (x *controlServiceGetDiagnosticNodeServer) Send(m *DiagnosticNode) error {
	return x.ServerStream.SendMsg(m)
}

func _ControlService_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServiceServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ControlService/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServiceServer).Status(ctx, req.(*StatusMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _ControlService_SubscribeEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EventMessage)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ControlServiceServer).SubscribeEvents(m, &controlServiceSubscribeEventsServer{stream})
}

type ControlService_SubscribeEventsServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type controlServiceSubscribeEventsServer struct {
	grpc.ServerStream
}

func (x *controlServiceSubscribeEventsServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

func _ControlService_ExportSnapshot_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ControlServiceServer).ExportSnapshot(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.ControlService/ExportSnapshot",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ControlServiceServer).ExportSnapshot(ctx, req.(*ExportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ControlService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "service.ControlService",
	HandlerType: (*ControlServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Log",
			Handler:    _ControlService_Log_Handler,
		},
		{
			MethodName: "Quit",
			Handler:    _ControlService_Quit_Handler,
		},
		{
			MethodName: "SwitchPhase",
			Handler:    _ControlService_SwitchPhase_Handler,
		},
		{
			MethodName: "GetPackageNode",
			Handler:    _ControlService_GetPackageNode_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _ControlService_Status_Handler,
		},
		{
			MethodName: "ExportSnapshot",
			Handler:    _ControlService_ExportSnapshot_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetFileNode",
			Handler:       _ControlService_GetFileNode_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "GetDiagnosticNode",
			Handler:       _ControlService_GetDiagnosticNode_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeEvents",
			Handler:       _ControlService_SubscribeEvents_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "controlservice.proto",
}

func init() {
	proto.RegisterFile("controlservice.proto", fileDescriptor_controlservice_788e7d435b159ea4)
}

var fileDescriptor_controlservice_788e7d435b159ea4 = []byte{
	// 591 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x94, 0x6d, 0x6f, 0xd3, 0x30,
	0x10, 0xc7, 0xd3, 0xf5, 0xf9, 0xda, 0x65, 0xcc, 0xeb, 0x4a, 0x15, 0xa6, 0x69, 0x18, 0x24, 0x0a,
	0x12, 0x15, 0x2a, 0xe2, 0x05, 0x8f, 0x12, 0x74, 0xa5, 0x1a, 0x2a, 0x68, 0xa4, 0x6f, 0x78, 0x9b,
	0xa6, 0x26, 0x8d, 0x96, 0xc6, 0x69, 0xec, 0x6c, 0x7c, 0x00, 0xbe, 0x22, 0xdf, 0x07, 0xc5, 0x8e,
	0x13, 0x47, 0x1d, 0x1a, 0xef, 0xec, 0xbb, 0xff, 0xfd, 0xef, 0xce, 0xf9, 0xb5, 0xd0, 0x73, 0x69,
	0xc8, 0x63, 0x1a, 0x30, 0x12, 0x5f, 0xfb, 0x2e, 0x19, 0x45, 0x31, 0xe5, 0x14, 0x35, 0xb3, 0xab,
	0x75, 0xb0, 0x72, 0xb8, 0xb3, 0xa1, 0x2b, 0x12, 0xc8, 0x0c, 0x3e, 0x05, 0x98, 0x53, 0xef, 0x2b,
	0x61, 0xcc, 0xf1, 0x08, 0xba, 0x07, 0xd5, 0x0d, 0xf3, 0x06, 0x95, 0xb3, 0xca, 0xb0, 0x6b, 0xa7,
	0x47, 0xfc, 0x04, 0x3a, 0x73, 0xea, 0xd9, 0x84, 0x45, 0x34, 0x64, 0x04, 0x0d, 0xa0, 0xc9, 0x12,
	0xd7, 0x25, 0x8c, 0x09, 0x51, 0xcb, 0x56, 0x57, 0xfc, 0x10, 0x3a, 0xdf, 0x13, 0x9f, 0x2b, 0x27,
	0x04, 0xb5, 0x2b, 0x3f, 0x08, 0x32, 0x95, 0x38, 0xe3, 0x21, 0x74, 0x53, 0xc9, 0x7f, 0x98, 0xbd,
	0x01, 0xb4, 0xb8, 0xf1, 0xb9, 0xbb, 0xbe, 0x5c, 0x3b, 0x8c, 0x28, 0xcf, 0xc7, 0x50, 0x8f, 0xd2,
	0xbb, 0x50, 0x9b, 0x63, 0x73, 0xa4, 0x96, 0x14, 0x2a, 0x5b, 0x26, 0xf1, 0x14, 0x8e, 0xb4, 0xda,
	0xbb, 0x9b, 0xa1, 0x1e, 0xd4, 0x49, 0x1c, 0xd3, 0x78, 0xb0, 0x77, 0x56, 0x19, 0xb6, 0x6d, 0x79,
	0xc1, 0x2e, 0xa0, 0x19, 0xe1, 0x9f, 0xfd, 0x80, 0x7c, 0xa3, 0xab, 0x7c, 0x84, 0xe7, 0xd0, 0xfa,
	0x99, 0x85, 0x84, 0x4d, 0x67, 0x7c, 0x98, 0x4f, 0xa1, 0xb4, 0x76, 0x2e, 0x41, 0xa7, 0x00, 0x49,
	0xe8, 0x6f, 0x13, 0x59, 0xb0, 0x27, 0xfa, 0x6a, 0x11, 0xfc, 0x1e, 0xf6, 0x17, 0xdc, 0xe1, 0x09,
	0x53, 0xfe, 0x3d, 0x7d, 0xc5, 0x56, 0xb6, 0x12, 0xea, 0x43, 0x83, 0x89, 0x95, 0x32, 0x8b, 0xec,
	0x86, 0x7f, 0x57, 0xc0, 0x94, 0xf5, 0xf9, 0x9a, 0x25, 0x83, 0xb6, 0x32, 0x18, 0x42, 0x53, 0x1c,
	0x2e, 0xce, 0x85, 0xc3, 0xee, 0xdb, 0xa9, 0x34, 0x3a, 0x81, 0xb6, 0x34, 0xf7, 0x43, 0x6f, 0x50,
	0x15, 0xdd, 0x8a, 0x40, 0xf1, 0x54, 0x35, 0xfd, 0xa9, 0x5e, 0x43, 0x77, 0x7a, 0x4d, 0xc2, 0xfc,
	0xdb, 0x3f, 0x85, 0xba, 0x1b, 0x38, 0xd9, 0x43, 0x9b, 0xe3, 0xa3, 0xbc, 0x97, 0x50, 0x4d, 0xd2,
	0x94, 0x2d, 0x15, 0xf8, 0x11, 0xec, 0x4f, 0x7f, 0x45, 0x34, 0xe6, 0x36, 0xd9, 0x26, 0x84, 0xf1,
	0x94, 0x9b, 0x1b, 0xc7, 0xe7, 0x8a, 0x9b, 0xf4, 0x8c, 0x9f, 0x81, 0xa9, 0x44, 0x77, 0x7d, 0xcc,
	0xf1, 0x9f, 0x1a, 0x98, 0x13, 0xf9, 0x13, 0x58, 0xc8, 0xae, 0x68, 0x0c, 0xd5, 0x39, 0xf5, 0x50,
	0x31, 0x46, 0x01, 0xbc, 0xd5, 0xd3, 0x83, 0xca, 0x1e, 0x1b, 0xe8, 0x15, 0xd4, 0x52, 0x54, 0x51,
	0x91, 0xd7, 0xe0, 0xb6, 0x8e, 0x4b, 0x51, 0xad, 0xec, 0x0b, 0x74, 0x34, 0xf6, 0xd0, 0x83, 0x5c,
	0xb7, 0x4b, 0xb3, 0x75, 0x72, 0x5b, 0x52, 0xf3, 0xfa, 0x00, 0xe6, 0x8c, 0xf0, 0x4b, 0xc7, 0xbd,
	0x72, 0x3c, 0x49, 0x53, 0x31, 0x8c, 0x16, 0xb5, 0x6e, 0x8d, 0x62, 0x03, 0x7d, 0x84, 0x8e, 0x06,
	0xb0, 0x36, 0xcb, 0x2e, 0xd6, 0xd6, 0x2e, 0xc4, 0xd8, 0x78, 0x51, 0x41, 0x17, 0x70, 0x38, 0x23,
	0xfc, 0xdc, 0x77, 0xbc, 0x90, 0x32, 0xee, 0xbb, 0xc2, 0xe8, 0x7e, 0xae, 0x2d, 0x27, 0xac, 0x7f,
	0x25, 0x84, 0xd5, 0x5b, 0x68, 0x48, 0x52, 0x51, 0xbf, 0xd8, 0x5b, 0x47, 0x5f, 0x2b, 0x2f, 0x23,
	0x8d, 0x0d, 0xf4, 0x0e, 0x0e, 0x16, 0xc9, 0x92, 0xb9, 0xb1, 0xbf, 0x24, 0x82, 0x21, 0x86, 0x8e,
	0xcb, 0x50, 0x29, 0x13, 0xb3, 0x1c, 0x16, 0xad, 0x27, 0x0a, 0x9f, 0x45, 0xe8, 0x44, 0x6c, 0x4d,
	0xb9, 0x36, 0x42, 0x09, 0x3e, 0x6d, 0x84, 0x32, 0x6f, 0xd8, 0xf8, 0x34, 0x80, 0x3e, 0x8d, 0xbd,
	0xd1, 0x76, 0xc3, 0x78, 0x3c, 0xf2, 0xe2, 0xc8, 0x55, 0xd2, 0x1f, 0xc6, 0xb2, 0x21, 0xfe, 0x4a,
	0x5f, 0xfe, 0x0d, 0x00, 0x00, 0xff, 0xff, 0xcb, 0x92, 0x0d, 0xf3, 0x7c, 0x05, 0x00, 0x00,
}