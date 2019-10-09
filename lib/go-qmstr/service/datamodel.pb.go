// Code generated by protoc-gen-go. DO NOT EDIT.
// source: datamodel.proto

package service

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type EventClass int32

const (
	EventClass_ALL    EventClass = 0
	EventClass_PHASE  EventClass = 1
	EventClass_MODULE EventClass = 2
)

var EventClass_name = map[int32]string{
	0: "ALL",
	1: "PHASE",
	2: "MODULE",
}

var EventClass_value = map[string]int32{
	"ALL":    0,
	"PHASE":  1,
	"MODULE": 2,
}

func (x EventClass) String() string {
	return proto.EnumName(EventClass_name, int32(x))
}

func (EventClass) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{0}
}

type Phase int32

const (
	Phase_INIT     Phase = 0
	Phase_BUILD    Phase = 1
	Phase_ANALYSIS Phase = 2
	Phase_REPORT   Phase = 3
	Phase_FAIL     Phase = 4
)

var Phase_name = map[int32]string{
	0: "INIT",
	1: "BUILD",
	2: "ANALYSIS",
	3: "REPORT",
	4: "FAIL",
}

var Phase_value = map[string]int32{
	"INIT":     0,
	"BUILD":    1,
	"ANALYSIS": 2,
	"REPORT":   3,
	"FAIL":     4,
}

func (x Phase) String() string {
	return proto.EnumName(Phase_name, int32(x))
}

func (Phase) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{1}
}

type ExceptionType int32

const (
	ExceptionType_ERROR   ExceptionType = 0
	ExceptionType_WARNING ExceptionType = 1
)

var ExceptionType_name = map[int32]string{
	0: "ERROR",
	1: "WARNING",
}

var ExceptionType_value = map[string]int32{
	"ERROR":   0,
	"WARNING": 1,
}

func (x ExceptionType) String() string {
	return proto.EnumName(ExceptionType_name, int32(x))
}

func (ExceptionType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{2}
}

type DiagnosticNode_Severity int32

const (
	DiagnosticNode_UNDEF   DiagnosticNode_Severity = 0
	DiagnosticNode_INFO    DiagnosticNode_Severity = 1
	DiagnosticNode_WARNING DiagnosticNode_Severity = 2
	DiagnosticNode_ERROR   DiagnosticNode_Severity = 3
)

var DiagnosticNode_Severity_name = map[int32]string{
	0: "UNDEF",
	1: "INFO",
	2: "WARNING",
	3: "ERROR",
}

var DiagnosticNode_Severity_value = map[string]int32{
	"UNDEF":   0,
	"INFO":    1,
	"WARNING": 2,
	"ERROR":   3,
}

func (x DiagnosticNode_Severity) String() string {
	return proto.EnumName(DiagnosticNode_Severity_name, int32(x))
}

func (DiagnosticNode_Severity) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{2, 0}
}

type FileNode struct {
	Uid                  string                 `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	FileNodeType         string                 `protobuf:"bytes,2,opt,name=fileNodeType,proto3" json:"fileNodeType,omitempty"`
	Path                 string                 `protobuf:"bytes,4,opt,name=path,proto3" json:"path,omitempty"`
	Name                 string                 `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty"`
	FileData             *FileNode_FileDataNode `protobuf:"bytes,6,opt,name=fileData,proto3" json:"fileData,omitempty"`
	Timestamp            string                 `protobuf:"bytes,7,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	DerivedFrom          []*FileNode            `protobuf:"bytes,8,rep,name=derivedFrom,proto3" json:"derivedFrom,omitempty"`
	Dependencies         []*FileNode            `protobuf:"bytes,9,rep,name=dependencies,proto3" json:"dependencies,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *FileNode) Reset()         { *m = FileNode{} }
func (m *FileNode) String() string { return proto.CompactTextString(m) }
func (*FileNode) ProtoMessage()    {}
func (*FileNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{0}
}

func (m *FileNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileNode.Unmarshal(m, b)
}
func (m *FileNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileNode.Marshal(b, m, deterministic)
}
func (m *FileNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileNode.Merge(m, src)
}
func (m *FileNode) XXX_Size() int {
	return xxx_messageInfo_FileNode.Size(m)
}
func (m *FileNode) XXX_DiscardUnknown() {
	xxx_messageInfo_FileNode.DiscardUnknown(m)
}

var xxx_messageInfo_FileNode proto.InternalMessageInfo

func (m *FileNode) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *FileNode) GetFileNodeType() string {
	if m != nil {
		return m.FileNodeType
	}
	return ""
}

func (m *FileNode) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *FileNode) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *FileNode) GetFileData() *FileNode_FileDataNode {
	if m != nil {
		return m.FileData
	}
	return nil
}

func (m *FileNode) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func (m *FileNode) GetDerivedFrom() []*FileNode {
	if m != nil {
		return m.DerivedFrom
	}
	return nil
}

func (m *FileNode) GetDependencies() []*FileNode {
	if m != nil {
		return m.Dependencies
	}
	return nil
}

type FileNode_FileDataNode struct {
	Uid                  string            `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	FileDataNodeType     string            `protobuf:"bytes,2,opt,name=fileDataNodeType,proto3" json:"fileDataNodeType,omitempty"`
	Hash                 string            `protobuf:"bytes,3,opt,name=hash,proto3" json:"hash,omitempty"`
	AdditionalInfo       []*InfoNode       `protobuf:"bytes,4,rep,name=additionalInfo,proto3" json:"additionalInfo,omitempty"`
	DiagnosticInfo       []*DiagnosticNode `protobuf:"bytes,5,rep,name=diagnosticInfo,proto3" json:"diagnosticInfo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *FileNode_FileDataNode) Reset()         { *m = FileNode_FileDataNode{} }
func (m *FileNode_FileDataNode) String() string { return proto.CompactTextString(m) }
func (*FileNode_FileDataNode) ProtoMessage()    {}
func (*FileNode_FileDataNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{0, 0}
}

func (m *FileNode_FileDataNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileNode_FileDataNode.Unmarshal(m, b)
}
func (m *FileNode_FileDataNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileNode_FileDataNode.Marshal(b, m, deterministic)
}
func (m *FileNode_FileDataNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileNode_FileDataNode.Merge(m, src)
}
func (m *FileNode_FileDataNode) XXX_Size() int {
	return xxx_messageInfo_FileNode_FileDataNode.Size(m)
}
func (m *FileNode_FileDataNode) XXX_DiscardUnknown() {
	xxx_messageInfo_FileNode_FileDataNode.DiscardUnknown(m)
}

var xxx_messageInfo_FileNode_FileDataNode proto.InternalMessageInfo

func (m *FileNode_FileDataNode) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *FileNode_FileDataNode) GetFileDataNodeType() string {
	if m != nil {
		return m.FileDataNodeType
	}
	return ""
}

func (m *FileNode_FileDataNode) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

func (m *FileNode_FileDataNode) GetAdditionalInfo() []*InfoNode {
	if m != nil {
		return m.AdditionalInfo
	}
	return nil
}

func (m *FileNode_FileDataNode) GetDiagnosticInfo() []*DiagnosticNode {
	if m != nil {
		return m.DiagnosticInfo
	}
	return nil
}

type InfoNode struct {
	Uid                  string               `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	InfoNodeType         string               `protobuf:"bytes,2,opt,name=infoNodeType,proto3" json:"infoNodeType,omitempty"`
	Type                 string               `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	ConfidenceScore      float64              `protobuf:"fixed64,4,opt,name=confidenceScore,proto3" json:"confidenceScore,omitempty"`
	Analyzer             []*Analyzer          `protobuf:"bytes,5,rep,name=analyzer,proto3" json:"analyzer,omitempty"`
	DataNodes            []*InfoNode_DataNode `protobuf:"bytes,6,rep,name=dataNodes,proto3" json:"dataNodes,omitempty"`
	Timestamp            string               `protobuf:"bytes,7,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *InfoNode) Reset()         { *m = InfoNode{} }
func (m *InfoNode) String() string { return proto.CompactTextString(m) }
func (*InfoNode) ProtoMessage()    {}
func (*InfoNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{1}
}

func (m *InfoNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoNode.Unmarshal(m, b)
}
func (m *InfoNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoNode.Marshal(b, m, deterministic)
}
func (m *InfoNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoNode.Merge(m, src)
}
func (m *InfoNode) XXX_Size() int {
	return xxx_messageInfo_InfoNode.Size(m)
}
func (m *InfoNode) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoNode.DiscardUnknown(m)
}

var xxx_messageInfo_InfoNode proto.InternalMessageInfo

func (m *InfoNode) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *InfoNode) GetInfoNodeType() string {
	if m != nil {
		return m.InfoNodeType
	}
	return ""
}

func (m *InfoNode) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *InfoNode) GetConfidenceScore() float64 {
	if m != nil {
		return m.ConfidenceScore
	}
	return 0
}

func (m *InfoNode) GetAnalyzer() []*Analyzer {
	if m != nil {
		return m.Analyzer
	}
	return nil
}

func (m *InfoNode) GetDataNodes() []*InfoNode_DataNode {
	if m != nil {
		return m.DataNodes
	}
	return nil
}

func (m *InfoNode) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

type InfoNode_DataNode struct {
	Uid                  string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	DataNodeType         string   `protobuf:"bytes,2,opt,name=dataNodeType,proto3" json:"dataNodeType,omitempty"`
	Type                 string   `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
	Data                 string   `protobuf:"bytes,4,opt,name=data,proto3" json:"data,omitempty"`
	Timestamp            string   `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoNode_DataNode) Reset()         { *m = InfoNode_DataNode{} }
func (m *InfoNode_DataNode) String() string { return proto.CompactTextString(m) }
func (*InfoNode_DataNode) ProtoMessage()    {}
func (*InfoNode_DataNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{1, 0}
}

func (m *InfoNode_DataNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoNode_DataNode.Unmarshal(m, b)
}
func (m *InfoNode_DataNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoNode_DataNode.Marshal(b, m, deterministic)
}
func (m *InfoNode_DataNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoNode_DataNode.Merge(m, src)
}
func (m *InfoNode_DataNode) XXX_Size() int {
	return xxx_messageInfo_InfoNode_DataNode.Size(m)
}
func (m *InfoNode_DataNode) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoNode_DataNode.DiscardUnknown(m)
}

var xxx_messageInfo_InfoNode_DataNode proto.InternalMessageInfo

func (m *InfoNode_DataNode) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *InfoNode_DataNode) GetDataNodeType() string {
	if m != nil {
		return m.DataNodeType
	}
	return ""
}

func (m *InfoNode_DataNode) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *InfoNode_DataNode) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func (m *InfoNode_DataNode) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

type DiagnosticNode struct {
	Uid                  string                  `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	DiagnosticNodeType   string                  `protobuf:"bytes,2,opt,name=diagnosticNodeType,proto3" json:"diagnosticNodeType,omitempty"`
	Severity             DiagnosticNode_Severity `protobuf:"varint,3,opt,name=severity,proto3,enum=service.DiagnosticNode_Severity" json:"severity,omitempty"`
	Message              string                  `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	Analyzer             []*Analyzer             `protobuf:"bytes,5,rep,name=analyzer,proto3" json:"analyzer,omitempty"`
	Timestamp            string                  `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                `json:"-"`
	XXX_unrecognized     []byte                  `json:"-"`
	XXX_sizecache        int32                   `json:"-"`
}

func (m *DiagnosticNode) Reset()         { *m = DiagnosticNode{} }
func (m *DiagnosticNode) String() string { return proto.CompactTextString(m) }
func (*DiagnosticNode) ProtoMessage()    {}
func (*DiagnosticNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{2}
}

func (m *DiagnosticNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiagnosticNode.Unmarshal(m, b)
}
func (m *DiagnosticNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiagnosticNode.Marshal(b, m, deterministic)
}
func (m *DiagnosticNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiagnosticNode.Merge(m, src)
}
func (m *DiagnosticNode) XXX_Size() int {
	return xxx_messageInfo_DiagnosticNode.Size(m)
}
func (m *DiagnosticNode) XXX_DiscardUnknown() {
	xxx_messageInfo_DiagnosticNode.DiscardUnknown(m)
}

var xxx_messageInfo_DiagnosticNode proto.InternalMessageInfo

func (m *DiagnosticNode) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *DiagnosticNode) GetDiagnosticNodeType() string {
	if m != nil {
		return m.DiagnosticNodeType
	}
	return ""
}

func (m *DiagnosticNode) GetSeverity() DiagnosticNode_Severity {
	if m != nil {
		return m.Severity
	}
	return DiagnosticNode_UNDEF
}

func (m *DiagnosticNode) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *DiagnosticNode) GetAnalyzer() []*Analyzer {
	if m != nil {
		return m.Analyzer
	}
	return nil
}

func (m *DiagnosticNode) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

type Analyzer struct {
	Uid                  string              `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Name                 string              `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	AnalyzerNodeType     string              `protobuf:"bytes,3,opt,name=analyzerNodeType,proto3" json:"analyzerNodeType,omitempty"`
	TrustLevel           int64               `protobuf:"varint,4,opt,name=trustLevel,proto3" json:"trustLevel,omitempty"`
	PathSub              []*PathSubstitution `protobuf:"bytes,5,rep,name=pathSub,proto3" json:"pathSub,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *Analyzer) Reset()         { *m = Analyzer{} }
func (m *Analyzer) String() string { return proto.CompactTextString(m) }
func (*Analyzer) ProtoMessage()    {}
func (*Analyzer) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{3}
}

func (m *Analyzer) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Analyzer.Unmarshal(m, b)
}
func (m *Analyzer) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Analyzer.Marshal(b, m, deterministic)
}
func (m *Analyzer) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Analyzer.Merge(m, src)
}
func (m *Analyzer) XXX_Size() int {
	return xxx_messageInfo_Analyzer.Size(m)
}
func (m *Analyzer) XXX_DiscardUnknown() {
	xxx_messageInfo_Analyzer.DiscardUnknown(m)
}

var xxx_messageInfo_Analyzer proto.InternalMessageInfo

func (m *Analyzer) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *Analyzer) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Analyzer) GetAnalyzerNodeType() string {
	if m != nil {
		return m.AnalyzerNodeType
	}
	return ""
}

func (m *Analyzer) GetTrustLevel() int64 {
	if m != nil {
		return m.TrustLevel
	}
	return 0
}

func (m *Analyzer) GetPathSub() []*PathSubstitution {
	if m != nil {
		return m.PathSub
	}
	return nil
}

type PathSubstitution struct {
	Old                  string   `protobuf:"bytes,1,opt,name=old,proto3" json:"old,omitempty"`
	New                  string   `protobuf:"bytes,2,opt,name=new,proto3" json:"new,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PathSubstitution) Reset()         { *m = PathSubstitution{} }
func (m *PathSubstitution) String() string { return proto.CompactTextString(m) }
func (*PathSubstitution) ProtoMessage()    {}
func (*PathSubstitution) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{4}
}

func (m *PathSubstitution) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PathSubstitution.Unmarshal(m, b)
}
func (m *PathSubstitution) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PathSubstitution.Marshal(b, m, deterministic)
}
func (m *PathSubstitution) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PathSubstitution.Merge(m, src)
}
func (m *PathSubstitution) XXX_Size() int {
	return xxx_messageInfo_PathSubstitution.Size(m)
}
func (m *PathSubstitution) XXX_DiscardUnknown() {
	xxx_messageInfo_PathSubstitution.DiscardUnknown(m)
}

var xxx_messageInfo_PathSubstitution proto.InternalMessageInfo

func (m *PathSubstitution) GetOld() string {
	if m != nil {
		return m.Old
	}
	return ""
}

func (m *PathSubstitution) GetNew() string {
	if m != nil {
		return m.New
	}
	return ""
}

type PackageNode struct {
	Uid                  string            `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Name                 string            `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Version              string            `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	PackageNodeType      string            `protobuf:"bytes,4,opt,name=packageNodeType,proto3" json:"packageNodeType,omitempty"`
	Targets              []*FileNode       `protobuf:"bytes,5,rep,name=targets,proto3" json:"targets,omitempty"`
	AdditionalInfo       []*InfoNode       `protobuf:"bytes,6,rep,name=additionalInfo,proto3" json:"additionalInfo,omitempty"`
	BuildConfig          string            `protobuf:"bytes,7,opt,name=buildConfig,proto3" json:"buildConfig,omitempty"`
	DiagnosticInfo       []*DiagnosticNode `protobuf:"bytes,8,rep,name=diagnosticInfo,proto3" json:"diagnosticInfo,omitempty"`
	Timestamp            string            `protobuf:"bytes,9,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *PackageNode) Reset()         { *m = PackageNode{} }
func (m *PackageNode) String() string { return proto.CompactTextString(m) }
func (*PackageNode) ProtoMessage()    {}
func (*PackageNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{5}
}

func (m *PackageNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PackageNode.Unmarshal(m, b)
}
func (m *PackageNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PackageNode.Marshal(b, m, deterministic)
}
func (m *PackageNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PackageNode.Merge(m, src)
}
func (m *PackageNode) XXX_Size() int {
	return xxx_messageInfo_PackageNode.Size(m)
}
func (m *PackageNode) XXX_DiscardUnknown() {
	xxx_messageInfo_PackageNode.DiscardUnknown(m)
}

var xxx_messageInfo_PackageNode proto.InternalMessageInfo

func (m *PackageNode) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *PackageNode) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *PackageNode) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

func (m *PackageNode) GetPackageNodeType() string {
	if m != nil {
		return m.PackageNodeType
	}
	return ""
}

func (m *PackageNode) GetTargets() []*FileNode {
	if m != nil {
		return m.Targets
	}
	return nil
}

func (m *PackageNode) GetAdditionalInfo() []*InfoNode {
	if m != nil {
		return m.AdditionalInfo
	}
	return nil
}

func (m *PackageNode) GetBuildConfig() string {
	if m != nil {
		return m.BuildConfig
	}
	return ""
}

func (m *PackageNode) GetDiagnosticInfo() []*DiagnosticNode {
	if m != nil {
		return m.DiagnosticInfo
	}
	return nil
}

func (m *PackageNode) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

type ProjectNode struct {
	Uid                  string         `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	Name                 string         `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	ProjectNodeType      string         `protobuf:"bytes,3,opt,name=projectNodeType,proto3" json:"projectNodeType,omitempty"`
	Packages             []*PackageNode `protobuf:"bytes,4,rep,name=packages,proto3" json:"packages,omitempty"`
	AdditionalInfo       []*InfoNode    `protobuf:"bytes,5,rep,name=additionalInfo,proto3" json:"additionalInfo,omitempty"`
	Timestamp            string         `protobuf:"bytes,6,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *ProjectNode) Reset()         { *m = ProjectNode{} }
func (m *ProjectNode) String() string { return proto.CompactTextString(m) }
func (*ProjectNode) ProtoMessage()    {}
func (*ProjectNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{6}
}

func (m *ProjectNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ProjectNode.Unmarshal(m, b)
}
func (m *ProjectNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ProjectNode.Marshal(b, m, deterministic)
}
func (m *ProjectNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProjectNode.Merge(m, src)
}
func (m *ProjectNode) XXX_Size() int {
	return xxx_messageInfo_ProjectNode.Size(m)
}
func (m *ProjectNode) XXX_DiscardUnknown() {
	xxx_messageInfo_ProjectNode.DiscardUnknown(m)
}

var xxx_messageInfo_ProjectNode proto.InternalMessageInfo

func (m *ProjectNode) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *ProjectNode) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *ProjectNode) GetProjectNodeType() string {
	if m != nil {
		return m.ProjectNodeType
	}
	return ""
}

func (m *ProjectNode) GetPackages() []*PackageNode {
	if m != nil {
		return m.Packages
	}
	return nil
}

func (m *ProjectNode) GetAdditionalInfo() []*InfoNode {
	if m != nil {
		return m.AdditionalInfo
	}
	return nil
}

func (m *ProjectNode) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

type Event struct {
	Class                EventClass `protobuf:"varint,1,opt,name=class,proto3,enum=service.EventClass" json:"class,omitempty"`
	Message              string     `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *Event) Reset()         { *m = Event{} }
func (m *Event) String() string { return proto.CompactTextString(m) }
func (*Event) ProtoMessage()    {}
func (*Event) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{7}
}

func (m *Event) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Event.Unmarshal(m, b)
}
func (m *Event) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Event.Marshal(b, m, deterministic)
}
func (m *Event) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Event.Merge(m, src)
}
func (m *Event) XXX_Size() int {
	return xxx_messageInfo_Event.Size(m)
}
func (m *Event) XXX_DiscardUnknown() {
	xxx_messageInfo_Event.DiscardUnknown(m)
}

var xxx_messageInfo_Event proto.InternalMessageInfo

func (m *Event) GetClass() EventClass {
	if m != nil {
		return m.Class
	}
	return EventClass_ALL
}

func (m *Event) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

type QmstrStateNode struct {
	Uid                  string   `protobuf:"bytes,1,opt,name=uid,proto3" json:"uid,omitempty"`
	QmstrStateNodeType   string   `protobuf:"bytes,2,opt,name=qmstrStateNodeType,proto3" json:"qmstrStateNodeType,omitempty"`
	Phase                Phase    `protobuf:"varint,3,opt,name=phase,proto3,enum=service.Phase" json:"phase,omitempty"`
	Done                 bool     `protobuf:"varint,4,opt,name=done,proto3" json:"done,omitempty"`
	Timestamp            string   `protobuf:"bytes,5,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QmstrStateNode) Reset()         { *m = QmstrStateNode{} }
func (m *QmstrStateNode) String() string { return proto.CompactTextString(m) }
func (*QmstrStateNode) ProtoMessage()    {}
func (*QmstrStateNode) Descriptor() ([]byte, []int) {
	return fileDescriptor_f0073e28e4cd411f, []int{8}
}

func (m *QmstrStateNode) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QmstrStateNode.Unmarshal(m, b)
}
func (m *QmstrStateNode) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QmstrStateNode.Marshal(b, m, deterministic)
}
func (m *QmstrStateNode) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QmstrStateNode.Merge(m, src)
}
func (m *QmstrStateNode) XXX_Size() int {
	return xxx_messageInfo_QmstrStateNode.Size(m)
}
func (m *QmstrStateNode) XXX_DiscardUnknown() {
	xxx_messageInfo_QmstrStateNode.DiscardUnknown(m)
}

var xxx_messageInfo_QmstrStateNode proto.InternalMessageInfo

func (m *QmstrStateNode) GetUid() string {
	if m != nil {
		return m.Uid
	}
	return ""
}

func (m *QmstrStateNode) GetQmstrStateNodeType() string {
	if m != nil {
		return m.QmstrStateNodeType
	}
	return ""
}

func (m *QmstrStateNode) GetPhase() Phase {
	if m != nil {
		return m.Phase
	}
	return Phase_INIT
}

func (m *QmstrStateNode) GetDone() bool {
	if m != nil {
		return m.Done
	}
	return false
}

func (m *QmstrStateNode) GetTimestamp() string {
	if m != nil {
		return m.Timestamp
	}
	return ""
}

func init() {
	proto.RegisterEnum("service.EventClass", EventClass_name, EventClass_value)
	proto.RegisterEnum("service.Phase", Phase_name, Phase_value)
	proto.RegisterEnum("service.ExceptionType", ExceptionType_name, ExceptionType_value)
	proto.RegisterEnum("service.DiagnosticNode_Severity", DiagnosticNode_Severity_name, DiagnosticNode_Severity_value)
	proto.RegisterType((*FileNode)(nil), "service.FileNode")
	proto.RegisterType((*FileNode_FileDataNode)(nil), "service.FileNode.FileDataNode")
	proto.RegisterType((*InfoNode)(nil), "service.InfoNode")
	proto.RegisterType((*InfoNode_DataNode)(nil), "service.InfoNode.DataNode")
	proto.RegisterType((*DiagnosticNode)(nil), "service.DiagnosticNode")
	proto.RegisterType((*Analyzer)(nil), "service.Analyzer")
	proto.RegisterType((*PathSubstitution)(nil), "service.PathSubstitution")
	proto.RegisterType((*PackageNode)(nil), "service.PackageNode")
	proto.RegisterType((*ProjectNode)(nil), "service.ProjectNode")
	proto.RegisterType((*Event)(nil), "service.Event")
	proto.RegisterType((*QmstrStateNode)(nil), "service.QmstrStateNode")
}

func init() { proto.RegisterFile("datamodel.proto", fileDescriptor_f0073e28e4cd411f) }

var fileDescriptor_f0073e28e4cd411f = []byte{
	// 940 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x56, 0x5d, 0x8f, 0xf2, 0x44,
	0x14, 0xde, 0x02, 0x85, 0x72, 0x40, 0xde, 0x3a, 0x1a, 0xad, 0x1b, 0xf3, 0x86, 0x34, 0x26, 0xae,
	0xa8, 0xc4, 0xec, 0xc6, 0xcf, 0x98, 0x28, 0xef, 0x0b, 0x28, 0x09, 0xb2, 0x38, 0xec, 0xc6, 0x78,
	0x39, 0xdb, 0x99, 0x65, 0xab, 0xa5, 0xad, 0x9d, 0x01, 0x5d, 0x7f, 0x80, 0x77, 0xfa, 0x43, 0xbc,
	0xd2, 0xff, 0xe3, 0x85, 0x97, 0xfe, 0x0c, 0x33, 0xd3, 0x0f, 0xda, 0xd2, 0x25, 0xef, 0xde, 0x0d,
	0xcf, 0x79, 0xce, 0xf4, 0x9c, 0xe7, 0x9c, 0x79, 0x02, 0x3c, 0xa1, 0x44, 0x90, 0x4d, 0x40, 0x99,
	0x37, 0x0c, 0xa3, 0x40, 0x04, 0xa8, 0xc5, 0x59, 0xb4, 0x73, 0x1d, 0x66, 0xff, 0xde, 0x00, 0x63,
	0xea, 0x7a, 0x6c, 0x11, 0x50, 0x86, 0x4c, 0xa8, 0x6f, 0x5d, 0x6a, 0x69, 0x7d, 0xed, 0xac, 0x8d,
	0xe5, 0x11, 0xd9, 0xd0, 0xbd, 0x4d, 0xa2, 0x57, 0xf7, 0x21, 0xb3, 0x6a, 0x2a, 0x54, 0xc0, 0x10,
	0x82, 0x46, 0x48, 0xc4, 0x9d, 0xd5, 0x50, 0x31, 0x75, 0x96, 0x98, 0x4f, 0x36, 0xcc, 0xd2, 0x63,
	0x4c, 0x9e, 0xd1, 0x67, 0x60, 0xc8, 0xbc, 0x31, 0x11, 0xc4, 0x6a, 0xf6, 0xb5, 0xb3, 0xce, 0xf9,
	0xd3, 0x61, 0x52, 0xc6, 0x30, 0x2d, 0x41, 0x1d, 0x24, 0x43, 0xfe, 0xc0, 0x19, 0x1f, 0xbd, 0x09,
	0x6d, 0xe1, 0x6e, 0x18, 0x17, 0x64, 0x13, 0x5a, 0x2d, 0x75, 0xe9, 0x1e, 0x40, 0x17, 0xd0, 0xa1,
	0x2c, 0x72, 0x77, 0x8c, 0x4e, 0xa3, 0x60, 0x63, 0x19, 0xfd, 0xfa, 0x59, 0xe7, 0xfc, 0xe5, 0x83,
	0xcb, 0x71, 0x9e, 0x85, 0x3e, 0x84, 0x2e, 0x65, 0x21, 0xf3, 0x29, 0xf3, 0x1d, 0x97, 0x71, 0xab,
	0xfd, 0x50, 0x56, 0x81, 0x76, 0xfa, 0x8f, 0x06, 0xdd, 0x7c, 0x91, 0x15, 0xa2, 0x0d, 0xc0, 0xbc,
	0xcd, 0x31, 0x72, 0xc2, 0x1d, 0xe0, 0x52, 0xa8, 0x3b, 0xc2, 0xef, 0xac, 0x7a, 0x2c, 0x94, 0x3c,
	0xa3, 0x4f, 0xa1, 0x47, 0x28, 0x75, 0x85, 0x1b, 0xf8, 0xc4, 0x9b, 0xf9, 0xb7, 0x81, 0xd5, 0x28,
	0xd5, 0x26, 0x41, 0x55, 0x5b, 0x89, 0x88, 0xbe, 0x80, 0x1e, 0x75, 0xc9, 0xda, 0x0f, 0xb8, 0x70,
	0x1d, 0x95, 0xaa, 0xab, 0xd4, 0xd7, 0xb3, 0xd4, 0x71, 0x16, 0x8e, 0x2f, 0x28, 0xd2, 0xed, 0x3f,
	0xea, 0x60, 0xa4, 0xb7, 0x57, 0xef, 0x83, 0x9b, 0x44, 0xf3, 0xfb, 0x90, 0xc7, 0x64, 0x4b, 0x42,
	0xc6, 0x92, 0x96, 0xe4, 0x19, 0x9d, 0xc1, 0x13, 0x27, 0xf0, 0x6f, 0x5d, 0xa9, 0x22, 0x5b, 0x39,
	0x41, 0xc4, 0xd4, 0xba, 0x68, 0xb8, 0x0c, 0xa3, 0xf7, 0xc1, 0x20, 0x3e, 0xf1, 0xee, 0x7f, 0x65,
	0x51, 0x52, 0xfb, 0xbe, 0xed, 0x51, 0x12, 0xc0, 0x19, 0x05, 0x7d, 0x02, 0x6d, 0x9a, 0xe8, 0xc9,
	0xad, 0xa6, 0xe2, 0x9f, 0x1e, 0xc8, 0x34, 0xcc, 0x36, 0x6a, 0x4f, 0x3e, 0xbe, 0x52, 0xa7, 0xbf,
	0x69, 0x60, 0x1c, 0x19, 0xb1, 0x0d, 0x5d, 0x7a, 0x38, 0xde, 0x02, 0x56, 0xa9, 0x03, 0x82, 0x86,
	0xe4, 0xa4, 0x6f, 0x85, 0x1e, 0xec, 0xb6, 0x5e, 0x2a, 0xc4, 0xfe, 0xab, 0x06, 0xbd, 0xe2, 0xcc,
	0x2a, 0xca, 0x19, 0x02, 0xa2, 0x05, 0x4e, 0xae, 0xa8, 0x8a, 0x08, 0xfa, 0x1c, 0x0c, 0xce, 0x76,
	0x2c, 0x72, 0xc5, 0xbd, 0x2a, 0xaf, 0x77, 0xde, 0x7f, 0x60, 0x41, 0x86, 0xab, 0x84, 0x87, 0xb3,
	0x0c, 0x64, 0x41, 0x6b, 0xc3, 0x38, 0x27, 0x6b, 0x96, 0xf4, 0x91, 0xfe, 0x7c, 0xec, 0xf0, 0x0a,
	0x9d, 0x37, 0xcb, 0x9d, 0x7f, 0x0c, 0x46, 0xfa, 0x71, 0xd4, 0x06, 0xfd, 0x7a, 0x31, 0x9e, 0x4c,
	0xcd, 0x13, 0x64, 0x40, 0x63, 0xb6, 0x98, 0x5e, 0x9a, 0x1a, 0xea, 0x40, 0xeb, 0xbb, 0x11, 0x5e,
	0xcc, 0x16, 0x5f, 0x99, 0x35, 0xc9, 0x98, 0x60, 0x7c, 0x89, 0xcd, 0xba, 0xfd, 0xb7, 0x06, 0x46,
	0xfa, 0xb5, 0x0a, 0xb1, 0x52, 0x6f, 0xaa, 0xe5, 0xbc, 0x69, 0x00, 0x66, 0x5a, 0x55, 0x26, 0x5f,
	0x3c, 0xb7, 0x03, 0x1c, 0x3d, 0x05, 0x10, 0xd1, 0x96, 0x8b, 0x39, 0xdb, 0x31, 0x4f, 0x29, 0x50,
	0xc7, 0x39, 0x04, 0x5d, 0x40, 0x4b, 0x7a, 0xe0, 0x6a, 0x7b, 0x93, 0x68, 0xf0, 0x46, 0xa6, 0xc1,
	0x32, 0xc6, 0xb9, 0x70, 0xc5, 0x56, 0xbe, 0x5a, 0x9c, 0x32, 0xed, 0x8f, 0xc0, 0x2c, 0x07, 0x65,
	0xe9, 0x81, 0x97, 0x95, 0x1e, 0x78, 0x54, 0x22, 0x3e, 0xfb, 0x39, 0xa9, 0x5c, 0x1e, 0xed, 0x7f,
	0x6b, 0xd0, 0x59, 0x12, 0xe7, 0x47, 0xb2, 0x7e, 0xc8, 0xc2, 0xab, 0xda, 0xb5, 0xa0, 0xb5, 0x63,
	0x11, 0x77, 0x03, 0x3f, 0xe9, 0x32, 0xfd, 0x29, 0x1f, 0x6a, 0xb8, 0xbf, 0x4e, 0xe9, 0x10, 0xcf,
	0xb8, 0x0c, 0xa3, 0x77, 0xa1, 0x25, 0x48, 0xb4, 0x66, 0x82, 0x1f, 0x8c, 0x3a, 0xb3, 0xce, 0x94,
	0x51, 0x61, 0x69, 0xcd, 0x17, 0xb5, 0xb4, 0x3e, 0x74, 0x6e, 0xb6, 0xae, 0x47, 0x9f, 0x4b, 0xa3,
	0x58, 0x27, 0x2f, 0x35, 0x0f, 0x55, 0x98, 0x9e, 0xf1, 0x28, 0xd3, 0x2b, 0xee, 0x61, 0xbb, 0xbc,
	0x87, 0xff, 0x69, 0xd0, 0x59, 0x46, 0xc1, 0x0f, 0xcc, 0x11, 0x8f, 0x90, 0x58, 0x0a, 0xb9, 0x4f,
	0xca, 0x2d, 0x54, 0x19, 0x46, 0x1f, 0x80, 0x91, 0x68, 0xcb, 0x13, 0xa3, 0x7f, 0x35, 0xb7, 0x30,
	0x99, 0xe8, 0x38, 0x63, 0x55, 0xa8, 0xa9, 0xbf, 0xa8, 0x9a, 0xc7, 0x9f, 0xdc, 0x1c, 0xf4, 0xc9,
	0x8e, 0xf9, 0x02, 0xbd, 0x03, 0xba, 0xe3, 0x11, 0xce, 0x55, 0x97, 0xbd, 0xf3, 0x57, 0xb2, 0x8b,
	0x55, 0xf8, 0xb9, 0x0c, 0xe1, 0x98, 0x91, 0x77, 0x83, 0x5a, 0xc1, 0x0d, 0xec, 0x3f, 0x35, 0xe8,
	0x7d, 0xbb, 0xe1, 0x22, 0x5a, 0x09, 0x22, 0xd8, 0xc3, 0xd6, 0xf5, 0x53, 0x81, 0x93, 0xb7, 0xae,
	0xc3, 0x08, 0x7a, 0x0b, 0xf4, 0xf0, 0x8e, 0x70, 0x96, 0xf8, 0x56, 0x6f, 0x2f, 0x95, 0x44, 0x71,
	0x1c, 0x54, 0x3e, 0x1b, 0xf8, 0xf1, 0xee, 0x1a, 0x58, 0x9d, 0x8f, 0xfb, 0xec, 0xe0, 0x3d, 0x80,
	0x7d, 0x6f, 0xa8, 0x05, 0xf5, 0xd1, 0x7c, 0x6e, 0x9e, 0x48, 0x5b, 0x59, 0x7e, 0x3d, 0x5a, 0x4d,
	0x4c, 0x0d, 0x01, 0x34, 0xbf, 0xb9, 0x1c, 0x5f, 0xcf, 0x27, 0x66, 0x6d, 0xf0, 0x25, 0xe8, 0xea,
	0x7b, 0xb1, 0x1b, 0xcd, 0xae, 0x62, 0xe6, 0xb3, 0xeb, 0xd9, 0x7c, 0x6c, 0x6a, 0xa8, 0x0b, 0xc6,
	0x68, 0x31, 0x9a, 0x7f, 0xbf, 0x9a, 0xad, 0xcc, 0x9a, 0xcc, 0xc3, 0x93, 0xe5, 0x25, 0xbe, 0x32,
	0xeb, 0x92, 0x3e, 0x1d, 0xcd, 0xe6, 0x66, 0x63, 0xf0, 0x36, 0xbc, 0x34, 0xf9, 0xc5, 0x61, 0xa1,
	0x9c, 0x8d, 0x6a, 0x2c, 0x33, 0xb0, 0x93, 0xbc, 0xb1, 0x69, 0xcf, 0x2c, 0x78, 0x2d, 0x88, 0xd6,
	0x43, 0x25, 0xc5, 0x70, 0x1d, 0x85, 0x4e, 0xda, 0xf1, 0x4d, 0x53, 0xfd, 0x97, 0xbb, 0xf8, 0x3f,
	0x00, 0x00, 0xff, 0xff, 0xdb, 0x5f, 0x9a, 0x67, 0xde, 0x09, 0x00, 0x00,
}
