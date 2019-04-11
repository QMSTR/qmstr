# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: buildservice.proto

import sys
_b=sys.version_info[0]<3 and (lambda x:x) or (lambda x:x.encode('latin1'))
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()


from . import datamodel_pb2 as datamodel__pb2


DESCRIPTOR = _descriptor.FileDescriptor(
  name='buildservice.proto',
  package='service',
  syntax='proto3',
  serialized_options=_b('\n\026org.qmstr.grpc.service'),
  serialized_pb=_b('\n\x12\x62uildservice.proto\x12\x07service\x1a\x0f\x64\x61tamodel.proto\" \n\rBuildResponse\x12\x0f\n\x07success\x18\x01 \x01(\x08\";\n\x0fPushFileMessage\x12\x0c\n\x04name\x18\x01 \x01(\t\x12\x0c\n\x04hash\x18\x02 \x01(\t\x12\x0c\n\x04\x64\x61ta\x18\x03 \x01(\x0c\" \n\x10PushFileResponse\x12\x0c\n\x04path\x18\x01 \x01(\t\"*\n\rDeleteMessage\x12\x0b\n\x03uid\x18\x01 \x01(\t\x12\x0c\n\x04\x65\x64ge\x18\x02 \x01(\t2\xc6\x04\n\x0c\x42uildService\x12\x36\n\x05\x42uild\x12\x11.service.FileNode\x1a\x16.service.BuildResponse\"\x00(\x01\x12=\n\x0eSendBuildError\x12\x11.service.InfoNode\x1a\x16.service.BuildResponse\"\x00\x12\x41\n\x08PushFile\x12\x18.service.PushFileMessage\x1a\x19.service.PushFileResponse\"\x00\x12\x38\n\x07Package\x12\x11.service.FileNode\x1a\x16.service.BuildResponse\"\x00(\x01\x12?\n\rCreatePackage\x12\x14.service.PackageNode\x1a\x16.service.BuildResponse\"\x00\x12?\n\rCreateProject\x12\x14.service.ProjectNode\x1a\x16.service.BuildResponse\"\x00\x12>\n\x0eGetProjectNode\x12\x14.service.ProjectNode\x1a\x14.service.ProjectNode\"\x00\x12@\n\nDeleteNode\x12\x16.service.DeleteMessage\x1a\x16.service.BuildResponse\"\x00(\x01\x12>\n\nDeleteEdge\x12\x16.service.DeleteMessage\x1a\x16.service.BuildResponse\"\x00\x42\x18\n\x16org.qmstr.grpc.serviceX\x00\x62\x06proto3')
  ,
  dependencies=[datamodel__pb2.DESCRIPTOR,])




_BUILDRESPONSE = _descriptor.Descriptor(
  name='BuildResponse',
  full_name='service.BuildResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='success', full_name='service.BuildResponse.success', index=0,
      number=1, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
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
  serialized_start=48,
  serialized_end=80,
)


_PUSHFILEMESSAGE = _descriptor.Descriptor(
  name='PushFileMessage',
  full_name='service.PushFileMessage',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='name', full_name='service.PushFileMessage.name', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='hash', full_name='service.PushFileMessage.hash', index=1,
      number=2, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='data', full_name='service.PushFileMessage.data', index=2,
      number=3, type=12, cpp_type=9, label=1,
      has_default_value=False, default_value=_b(""),
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
  serialized_start=82,
  serialized_end=141,
)


_PUSHFILERESPONSE = _descriptor.Descriptor(
  name='PushFileResponse',
  full_name='service.PushFileResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='path', full_name='service.PushFileResponse.path', index=0,
      number=1, type=9, cpp_type=9, label=1,
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
  serialized_start=143,
  serialized_end=175,
)


_DELETEMESSAGE = _descriptor.Descriptor(
  name='DeleteMessage',
  full_name='service.DeleteMessage',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  fields=[
    _descriptor.FieldDescriptor(
      name='uid', full_name='service.DeleteMessage.uid', index=0,
      number=1, type=9, cpp_type=9, label=1,
      has_default_value=False, default_value=_b("").decode('utf-8'),
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR),
    _descriptor.FieldDescriptor(
      name='edge', full_name='service.DeleteMessage.edge', index=1,
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
  serialized_start=177,
  serialized_end=219,
)

DESCRIPTOR.message_types_by_name['BuildResponse'] = _BUILDRESPONSE
DESCRIPTOR.message_types_by_name['PushFileMessage'] = _PUSHFILEMESSAGE
DESCRIPTOR.message_types_by_name['PushFileResponse'] = _PUSHFILERESPONSE
DESCRIPTOR.message_types_by_name['DeleteMessage'] = _DELETEMESSAGE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

BuildResponse = _reflection.GeneratedProtocolMessageType('BuildResponse', (_message.Message,), dict(
  DESCRIPTOR = _BUILDRESPONSE,
  __module__ = 'buildservice_pb2'
  # @@protoc_insertion_point(class_scope:service.BuildResponse)
  ))
_sym_db.RegisterMessage(BuildResponse)

PushFileMessage = _reflection.GeneratedProtocolMessageType('PushFileMessage', (_message.Message,), dict(
  DESCRIPTOR = _PUSHFILEMESSAGE,
  __module__ = 'buildservice_pb2'
  # @@protoc_insertion_point(class_scope:service.PushFileMessage)
  ))
_sym_db.RegisterMessage(PushFileMessage)

PushFileResponse = _reflection.GeneratedProtocolMessageType('PushFileResponse', (_message.Message,), dict(
  DESCRIPTOR = _PUSHFILERESPONSE,
  __module__ = 'buildservice_pb2'
  # @@protoc_insertion_point(class_scope:service.PushFileResponse)
  ))
_sym_db.RegisterMessage(PushFileResponse)

DeleteMessage = _reflection.GeneratedProtocolMessageType('DeleteMessage', (_message.Message,), dict(
  DESCRIPTOR = _DELETEMESSAGE,
  __module__ = 'buildservice_pb2'
  # @@protoc_insertion_point(class_scope:service.DeleteMessage)
  ))
_sym_db.RegisterMessage(DeleteMessage)


DESCRIPTOR._options = None

_BUILDSERVICE = _descriptor.ServiceDescriptor(
  name='BuildService',
  full_name='service.BuildService',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  serialized_start=222,
  serialized_end=804,
  methods=[
  _descriptor.MethodDescriptor(
    name='Build',
    full_name='service.BuildService.Build',
    index=0,
    containing_service=None,
    input_type=datamodel__pb2._FILENODE,
    output_type=_BUILDRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='SendBuildError',
    full_name='service.BuildService.SendBuildError',
    index=1,
    containing_service=None,
    input_type=datamodel__pb2._INFONODE,
    output_type=_BUILDRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='PushFile',
    full_name='service.BuildService.PushFile',
    index=2,
    containing_service=None,
    input_type=_PUSHFILEMESSAGE,
    output_type=_PUSHFILERESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='Package',
    full_name='service.BuildService.Package',
    index=3,
    containing_service=None,
    input_type=datamodel__pb2._FILENODE,
    output_type=_BUILDRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='CreatePackage',
    full_name='service.BuildService.CreatePackage',
    index=4,
    containing_service=None,
    input_type=datamodel__pb2._PACKAGENODE,
    output_type=_BUILDRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='CreateProject',
    full_name='service.BuildService.CreateProject',
    index=5,
    containing_service=None,
    input_type=datamodel__pb2._PROJECTNODE,
    output_type=_BUILDRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='GetProjectNode',
    full_name='service.BuildService.GetProjectNode',
    index=6,
    containing_service=None,
    input_type=datamodel__pb2._PROJECTNODE,
    output_type=datamodel__pb2._PROJECTNODE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='DeleteNode',
    full_name='service.BuildService.DeleteNode',
    index=7,
    containing_service=None,
    input_type=_DELETEMESSAGE,
    output_type=_BUILDRESPONSE,
    serialized_options=None,
  ),
  _descriptor.MethodDescriptor(
    name='DeleteEdge',
    full_name='service.BuildService.DeleteEdge',
    index=8,
    containing_service=None,
    input_type=_DELETEMESSAGE,
    output_type=_BUILDRESPONSE,
    serialized_options=None,
  ),
])
_sym_db.RegisterServiceDescriptor(_BUILDSERVICE)

DESCRIPTOR.services_by_name['BuildService'] = _BUILDSERVICE

# @@protoc_insertion_point(module_scope)
