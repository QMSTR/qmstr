# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: datamodel.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf.internal import enum_type_wrapper
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='datamodel.proto',
  package='service',
  syntax='proto3',
  serialized_options=_b('\n\026org.qmstr.grpc.service'),
  serialized_pb=_b('\n\x0f\x64\x61tamodel.proto\x12\x07service\"\xfb\x02\n\x08\x46ileNode\x12\x0b\n\x03uid\x18\x01 \x01(\t\x12\x14\n\x0c\x66ileNodeType\x18\x02 \x01(\t\x12(\n\x08\x66ileType\x18\x03 \x01(\x0e\x32\x16.service.FileNode.Type\x12\x0c\n\x04path\x18\x04 \x01(\t\x12\x0c\n\x04name\x18\x05 \x01(\t\x12\x0c\n\x04hash\x18\x06 \x01(\t\x12\x0e\n\x06\x62roken\x18\x07 \x01(\x08\x12&\n\x0b\x64\x65rivedFrom\x18\x08 \x03(\x0b\x32\x11.service.FileNode\x12)\n\x0e\x61\x64\x64itionalInfo\x18\t \x03(\x0b\x32\x11.service.InfoNode\x12/\n\x0e\x64iagnosticInfo\x18\n \x03(\x0b\x32\x17.service.DiagnosticNode\x12\'\n\x0c\x64\x65pendencies\x18\x0b \x03(\x0b\x32\x11.service.FileNode\";\n\x04Type\x12\t\n\x05UNDEF\x10\x00\x12\n\n\x06SOURCE\x10\x01\x12\x10\n\x0cINTERMEDIATE\x10\x02\x12\n\n\x06TARGET\x10\x03\"\xe6\x01\n\x08InfoNode\x12\x0b\n\x03uid\x18\x01 \x01(\t\x12\x14\n\x0cinfoNodeType\x18\x02 \x01(\t\x12\x0c\n\x04type\x18\x03 \x01(\t\x12\x17\n\x0f\x63onfidenceScore\x18\x04 \x01(\x01\x12#\n\x08\x61nalyzer\x18\x05 \x03(\x0b\x32\x11.service.Analyzer\x12-\n\tdataNodes\x18\x06 \x03(\x0b\x32\x1a.service.InfoNode.DataNode\x1a<\n\x08\x44\x61taNode\x12\x14\n\x0c\x64\x61taNodeType\x18\x01 \x01(\t\x12\x0c\n\x04type\x18\x02 \x01(\t\x12\x0c\n\x04\x64\x61ta\x18\x03 \x01(\t\"\xdc\x01\n\x0e\x44iagnosticNode\x12\x0b\n\x03uid\x18\x01 \x01(\t\x12\x1a\n\x12\x64iagnosticNodeType\x18\x02 \x01(\t\x12\x32\n\x08severity\x18\x03 \x01(\x0e\x32 .service.DiagnosticNode.Severity\x12\x0f\n\x07message\x18\x04 \x01(\t\x12#\n\x08\x61nalyzer\x18\x05 \x03(\x0b\x32\x11.service.Analyzer\"7\n\x08Severity\x12\t\n\x05UNDEF\x10\x00\x12\x08\n\x04INFO\x10\x01\x12\x0b\n\x07WARNING\x10\x02\x12\t\n\x05\x45RROR\x10\x03\"\x7f\n\x08\x41nalyzer\x12\x0b\n\x03uid\x18\x01 \x01(\t\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x18\n\x10\x61nalyzerNodeType\x18\x03 \x01(\t\x12\x12\n\ntrustLevel\x18\x04 \x01(\x03\x12*\n\x07pathSub\x18\x05 \x03(\x0b\x32\x19.service.PathSubstitution\",\n\x10PathSubstitution\x12\x0b\n\x03old\x18\x01 \x01(\t\x12\x0b\n\x03new\x18\x02 \x01(\t\"\xd6\x01\n\x0bPackageNode\x12\x0b\n\x03uid\x18\x01 \x01(\t\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x17\n\x0fpackageNodeType\x18\x04 \x01(\t\x12\"\n\x07targets\x18\x05 \x03(\x0b\x32\x11.service.FileNode\x12)\n\x0e\x61\x64\x64itionalInfo\x18\x06 \x03(\x0b\x32\x11.service.InfoNode\x12\x13\n\x0b\x62uildConfig\x18\x07 \x01(\t\x12/\n\x0e\x64iagnosticInfo\x18\x08 \x03(\x0b\x32\x17.service.DiagnosticNode\"\x94\x01\n\x0bProjectNode\x12\x0b\n\x03uid\x18\x01 \x01(\t\x12\x0c\n\x04name\x18\x02 \x01(\t\x12\x17\n\x0fprojectNodeType\x18\x03 \x01(\t\x12&\n\x08packages\x18\x04 \x03(\x0b\x32\x14.service.PackageNode\x12)\n\x0e\x61\x64\x64itionalInfo\x18\x05 \x03(\x0b\x32\x11.service.InfoNode\"<\n\x05\x45vent\x12\"\n\x05\x63lass\x18\x01 \x01(\x0e\x32\x13.service.EventClass\x12\x0f\n\x07message\x18\x02 \x01(\t\"X\n\x0eQmstrStateNode\x12\x0b\n\x03uid\x18\x01 \x01(\t\x12\x1a\n\x12qmstrStateNodeType\x18\x02 \x01(\t\x12\x1d\n\x05phase\x18\x03 \x01(\x0e\x32\x0e.service.Phase*,\n\nEventClass\x12\x07\n\x03\x41LL\x10\x00\x12\t\n\x05PHASE\x10\x01\x12\n\n\x06MODULE\x10\x02*@\n\x05Phase\x12\x08\n\x04INIT\x10\x00\x12\t\n\x05\x42UILD\x10\x01\x12\x0c\n\x08\x41NALYSIS\x10\x02\x12\n\n\x06REPORT\x10\x03\x12\x08\n\x04\x46\x41IL\x10\x04*\'\n\rExceptionType\x12\t\n\x05\x45RROR\x10\x00\x12\x0b\n\x07WARNING\x10\x01\x42\x18\n\x16org.qmstr.grpc.serviceb\x06proto3')
)

_EVENTCLASS = _descriptor.EnumDescriptor(
  name='EventClass',
  full_name='service.EventClass',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='ALL', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='PHASE', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='MODULE', index=2, number=2,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1561,
  serialized_end=1605,
)
_sym_db.RegisterEnumDescriptor(_EVENTCLASS)

EventClass = enum_type_wrapper.EnumTypeWrapper(_EVENTCLASS)
_PHASE = _descriptor.EnumDescriptor(
  name='Phase',
  full_name='service.Phase',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='INIT', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='BUILD', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='ANALYSIS', index=2, number=2,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='REPORT', index=3, number=3,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='FAIL', index=4, number=4,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1607,
  serialized_end=1671,
)
_sym_db.RegisterEnumDescriptor(_PHASE)

Phase = enum_type_wrapper.EnumTypeWrapper(_PHASE)
_EXCEPTIONTYPE = _descriptor.EnumDescriptor(
  name='ExceptionType',
  full_name='service.ExceptionType',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='ERROR', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='WARNING', index=1, number=1,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=1673,
  serialized_end=1712,
)
_sym_db.RegisterEnumDescriptor(_EXCEPTIONTYPE)

ExceptionType = enum_type_wrapper.EnumTypeWrapper(_EXCEPTIONTYPE)
ALL = 0
PHASE = 1
MODULE = 2
INIT = 0
BUILD = 1
ANALYSIS = 2
REPORT = 3
FAIL = 4
ERROR = 0
WARNING = 1


_FILENODE_TYPE = _descriptor.EnumDescriptor(
  name='Type',
  full_name='service.FileNode.Type',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='UNDEF', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='SOURCE', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='INTERMEDIATE', index=2, number=2,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='TARGET', index=3, number=3,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=349,
  serialized_end=408,
)
_sym_db.RegisterEnumDescriptor(_FILENODE_TYPE)

_DIAGNOSTICNODE_SEVERITY = _descriptor.EnumDescriptor(
  name='Severity',
  full_name='service.DiagnosticNode.Severity',
  filename=None,
  file=DESCRIPTOR,
  values=[
    _descriptor.EnumValueDescriptor(
      name='UNDEF', index=0, number=0,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='INFO', index=1, number=1,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='WARNING', index=2, number=2,
      serialized_options=None,
      type=None),
    _descriptor.EnumValueDescriptor(
      name='ERROR', index=3, number=3,
      serialized_options=None,
      type=None),
  ],
  containing_type=None,
  serialized_options=None,
  serialized_start=809,
  serialized_end=864,
)
_sym_db.RegisterEnumDescriptor(_DIAGNOSTICNODE_SEVERITY)


_FILENODE = _descriptor.Descriptor(
  name='FileNode',
  full_name='service.FileNode',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='uid', full_name='service.FileNode.uid', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='fileNodeType', full_name='service.FileNode.fileNodeType', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='fileType', full_name='service.FileNode.fileType', index=2,
      number=3, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='path', full_name='service.FileNode.path', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='name', full_name='service.FileNode.name', index=4,
      number=5, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='hash', full_name='service.FileNode.hash', index=5,
      number=6, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='broken', full_name='service.FileNode.broken', index=6,
      number=7, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='derivedFrom', full_name='service.FileNode.derivedFrom', index=7,
      number=8, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='additionalInfo', full_name='service.FileNode.additionalInfo', index=8,
      number=9, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='diagnosticInfo', full_name='service.FileNode.diagnosticInfo', index=9,
      number=10, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='dependencies', full_name='service.FileNode.dependencies', index=10,
      number=11, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
    _FILENODE_TYPE,
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=29,
  serialized_end=408,
)


_INFONODE_DATANODE = _descriptor.Descriptor(
  name='DataNode',
  full_name='service.InfoNode.DataNode',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='dataNodeType', full_name='service.InfoNode.DataNode.dataNodeType', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='service.InfoNode.DataNode.type', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='data', full_name='service.InfoNode.DataNode.data', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=581,
  serialized_end=641,
)

_INFONODE = _descriptor.Descriptor(
  name='InfoNode',
  full_name='service.InfoNode',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='uid', full_name='service.InfoNode.uid', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='infoNodeType', full_name='service.InfoNode.infoNodeType', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='type', full_name='service.InfoNode.type', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='confidenceScore', full_name='service.InfoNode.confidenceScore', index=3,
      number=4, type=1, cpp_type=5, label=1,
      has_default_value=False, default_value=float(0),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='analyzer', full_name='service.InfoNode.analyzer', index=4,
      number=5, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='dataNodes', full_name='service.InfoNode.dataNodes', index=5,
      number=6, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[_INFONODE_DATANODE, ],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=411,
  serialized_end=641,
)


_DIAGNOSTICNODE = _descriptor.Descriptor(
  name='DiagnosticNode',
  full_name='service.DiagnosticNode',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='uid', full_name='service.DiagnosticNode.uid', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='diagnosticNodeType', full_name='service.DiagnosticNode.diagnosticNodeType', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='severity', full_name='service.DiagnosticNode.severity', index=2,
      number=3, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='message', full_name='service.DiagnosticNode.message', index=3,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='analyzer', full_name='service.DiagnosticNode.analyzer', index=4,
      number=5, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
    _DIAGNOSTICNODE_SEVERITY,
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=644,
  serialized_end=864,
)


_ANALYZER = _descriptor.Descriptor(
  name='Analyzer',
  full_name='service.Analyzer',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='uid', full_name='service.Analyzer.uid', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='name', full_name='service.Analyzer.name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='analyzerNodeType', full_name='service.Analyzer.analyzerNodeType', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='trustLevel', full_name='service.Analyzer.trustLevel', index=3,
      number=4, type=3, cpp_type=2, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='pathSub', full_name='service.Analyzer.pathSub', index=4,
      number=5, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=866,
  serialized_end=993,
)


_PATHSUBSTITUTION = _descriptor.Descriptor(
  name='PathSubstitution',
  full_name='service.PathSubstitution',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='old', full_name='service.PathSubstitution.old', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='new', full_name='service.PathSubstitution.new', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=995,
  serialized_end=1039,
)


_PACKAGENODE = _descriptor.Descriptor(
  name='PackageNode',
  full_name='service.PackageNode',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='uid', full_name='service.PackageNode.uid', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='name', full_name='service.PackageNode.name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='packageNodeType', full_name='service.PackageNode.packageNodeType', index=2,
      number=4, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='targets', full_name='service.PackageNode.targets', index=3,
      number=5, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='additionalInfo', full_name='service.PackageNode.additionalInfo', index=4,
      number=6, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='buildConfig', full_name='service.PackageNode.buildConfig', index=5,
      number=7, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='diagnosticInfo', full_name='service.PackageNode.diagnosticInfo', index=6,
      number=8, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1042,
  serialized_end=1256,
)


_PROJECTNODE = _descriptor.Descriptor(
  name='ProjectNode',
  full_name='service.ProjectNode',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='uid', full_name='service.ProjectNode.uid', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='name', full_name='service.ProjectNode.name', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='projectNodeType', full_name='service.ProjectNode.projectNodeType', index=2,
      number=3, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='packages', full_name='service.ProjectNode.packages', index=3,
      number=4, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='additionalInfo', full_name='service.ProjectNode.additionalInfo', index=4,
      number=5, type=11, cpp_type=10, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1259,
  serialized_end=1407,
)


_EVENT = _descriptor.Descriptor(
  name='Event',
  full_name='service.Event',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='class', full_name='service.Event.class', index=0,
      number=1, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='message', full_name='service.Event.message', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1409,
  serialized_end=1469,
)


_QMSTRSTATENODE = _descriptor.Descriptor(
  name='QmstrStateNode',
  full_name='service.QmstrStateNode',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='uid', full_name='service.QmstrStateNode.uid', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='qmstrStateNodeType', full_name='service.QmstrStateNode.qmstrStateNodeType', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='phase', full_name='service.QmstrStateNode.phase', index=2,
      number=3, type=14, cpp_type=8, label=1,
      has_default_value=False, default_value=0,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=1471,
  serialized_end=1559,
)

_FILENODE.fields_by_name['fileType'].enum_type = _FILENODE_TYPE
_FILENODE.fields_by_name['derivedFrom'].message_type = _FILENODE
_FILENODE.fields_by_name['additionalInfo'].message_type = _INFONODE
_FILENODE.fields_by_name['diagnosticInfo'].message_type = _DIAGNOSTICNODE
_FILENODE.fields_by_name['dependencies'].message_type = _FILENODE
_FILENODE_TYPE.containing_type = _FILENODE
_INFONODE_DATANODE.containing_type = _INFONODE
_INFONODE.fields_by_name['analyzer'].message_type = _ANALYZER
_INFONODE.fields_by_name['dataNodes'].message_type = _INFONODE_DATANODE
_DIAGNOSTICNODE.fields_by_name['severity'].enum_type = _DIAGNOSTICNODE_SEVERITY
_DIAGNOSTICNODE.fields_by_name['analyzer'].message_type = _ANALYZER
_DIAGNOSTICNODE_SEVERITY.containing_type = _DIAGNOSTICNODE
_ANALYZER.fields_by_name['pathSub'].message_type = _PATHSUBSTITUTION
_PACKAGENODE.fields_by_name['targets'].message_type = _FILENODE
_PACKAGENODE.fields_by_name['additionalInfo'].message_type = _INFONODE
_PACKAGENODE.fields_by_name['diagnosticInfo'].message_type = _DIAGNOSTICNODE
_PROJECTNODE.fields_by_name['packages'].message_type = _PACKAGENODE
_PROJECTNODE.fields_by_name['additionalInfo'].message_type = _INFONODE
_EVENT.fields_by_name['class'].enum_type = _EVENTCLASS
_QMSTRSTATENODE.fields_by_name['phase'].enum_type = _PHASE
DESCRIPTOR.message_types_by_name['FileNode'] = _FILENODE
DESCRIPTOR.message_types_by_name['InfoNode'] = _INFONODE
DESCRIPTOR.message_types_by_name['DiagnosticNode'] = _DIAGNOSTICNODE
DESCRIPTOR.message_types_by_name['Analyzer'] = _ANALYZER
DESCRIPTOR.message_types_by_name['PathSubstitution'] = _PATHSUBSTITUTION
DESCRIPTOR.message_types_by_name['PackageNode'] = _PACKAGENODE
DESCRIPTOR.message_types_by_name['ProjectNode'] = _PROJECTNODE
DESCRIPTOR.message_types_by_name['Event'] = _EVENT
DESCRIPTOR.message_types_by_name['QmstrStateNode'] = _QMSTRSTATENODE
DESCRIPTOR.enum_types_by_name['EventClass'] = _EVENTCLASS
DESCRIPTOR.enum_types_by_name['Phase'] = _PHASE
DESCRIPTOR.enum_types_by_name['ExceptionType'] = _EXCEPTIONTYPE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

FileNode = _reflection.GeneratedProtocolMessageType('FileNode', (_message.Message,), dict(
  DESCRIPTOR = _FILENODE,
  __module__ = 'datamodel_pb2'
  # @@protoc_insertion_point(class_scope:service.FileNode)
  ))
_sym_db.RegisterMessage(FileNode)

InfoNode = _reflection.GeneratedProtocolMessageType('InfoNode', (_message.Message,), dict(

  DataNode = _reflection.GeneratedProtocolMessageType('DataNode', (_message.Message,), dict(
    DESCRIPTOR = _INFONODE_DATANODE,
    __module__ = 'datamodel_pb2'
    # @@protoc_insertion_point(class_scope:service.InfoNode.DataNode)
    ))
  ,
  DESCRIPTOR = _INFONODE,
  __module__ = 'datamodel_pb2'
  # @@protoc_insertion_point(class_scope:service.InfoNode)
  ))
_sym_db.RegisterMessage(InfoNode)
_sym_db.RegisterMessage(InfoNode.DataNode)

DiagnosticNode = _reflection.GeneratedProtocolMessageType('DiagnosticNode', (_message.Message,), dict(
  DESCRIPTOR = _DIAGNOSTICNODE,
  __module__ = 'datamodel_pb2'
  # @@protoc_insertion_point(class_scope:service.DiagnosticNode)
  ))
_sym_db.RegisterMessage(DiagnosticNode)

Analyzer = _reflection.GeneratedProtocolMessageType('Analyzer', (_message.Message,), dict(
  DESCRIPTOR = _ANALYZER,
  __module__ = 'datamodel_pb2'
  # @@protoc_insertion_point(class_scope:service.Analyzer)
  ))
_sym_db.RegisterMessage(Analyzer)

PathSubstitution = _reflection.GeneratedProtocolMessageType('PathSubstitution', (_message.Message,), dict(
  DESCRIPTOR = _PATHSUBSTITUTION,
  __module__ = 'datamodel_pb2'
  # @@protoc_insertion_point(class_scope:service.PathSubstitution)
  ))
_sym_db.RegisterMessage(PathSubstitution)

PackageNode = _reflection.GeneratedProtocolMessageType('PackageNode', (_message.Message,), dict(
  DESCRIPTOR = _PACKAGENODE,
  __module__ = 'datamodel_pb2'
  # @@protoc_insertion_point(class_scope:service.PackageNode)
  ))
_sym_db.RegisterMessage(PackageNode)

ProjectNode = _reflection.GeneratedProtocolMessageType('ProjectNode', (_message.Message,), dict(
  DESCRIPTOR = _PROJECTNODE,
  __module__ = 'datamodel_pb2'
  # @@protoc_insertion_point(class_scope:service.ProjectNode)
  ))
_sym_db.RegisterMessage(ProjectNode)

Event = _reflection.GeneratedProtocolMessageType('Event', (_message.Message,), dict(
  DESCRIPTOR = _EVENT,
  __module__ = 'datamodel_pb2'
  # @@protoc_insertion_point(class_scope:service.Event)
  ))
_sym_db.RegisterMessage(Event)

QmstrStateNode = _reflection.GeneratedProtocolMessageType('QmstrStateNode', (_message.Message,), dict(
  DESCRIPTOR = _QMSTRSTATENODE,
  __module__ = 'datamodel_pb2'
  # @@protoc_insertion_point(class_scope:service.QmstrStateNode)
  ))
_sym_db.RegisterMessage(QmstrStateNode)


DESCRIPTOR._options = None
# @@protoc_insertion_point(module_scope)
