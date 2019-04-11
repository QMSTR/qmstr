# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

from . import buildservice_pb2 as buildservice__pb2
from . import datamodel_pb2 as datamodel__pb2


class BuildServiceStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Build = channel.stream_unary(
        '/service.BuildService/Build',
        request_serializer=datamodel__pb2.FileNode.SerializeToString,
        response_deserializer=buildservice__pb2.BuildResponse.FromString,
        )
    self.SendBuildError = channel.unary_unary(
        '/service.BuildService/SendBuildError',
        request_serializer=datamodel__pb2.InfoNode.SerializeToString,
        response_deserializer=buildservice__pb2.BuildResponse.FromString,
        )
    self.PushFile = channel.unary_unary(
        '/service.BuildService/PushFile',
        request_serializer=buildservice__pb2.PushFileMessage.SerializeToString,
        response_deserializer=buildservice__pb2.PushFileResponse.FromString,
        )
    self.Package = channel.stream_unary(
        '/service.BuildService/Package',
        request_serializer=datamodel__pb2.FileNode.SerializeToString,
        response_deserializer=buildservice__pb2.BuildResponse.FromString,
        )
    self.CreatePackage = channel.unary_unary(
        '/service.BuildService/CreatePackage',
        request_serializer=datamodel__pb2.PackageNode.SerializeToString,
        response_deserializer=buildservice__pb2.BuildResponse.FromString,
        )
    self.CreateProject = channel.unary_unary(
        '/service.BuildService/CreateProject',
        request_serializer=datamodel__pb2.ProjectNode.SerializeToString,
        response_deserializer=buildservice__pb2.BuildResponse.FromString,
        )
    self.GetProjectNode = channel.unary_unary(
        '/service.BuildService/GetProjectNode',
        request_serializer=datamodel__pb2.ProjectNode.SerializeToString,
        response_deserializer=datamodel__pb2.ProjectNode.FromString,
        )
    self.DeleteNode = channel.stream_unary(
        '/service.BuildService/DeleteNode',
        request_serializer=buildservice__pb2.DeleteMessage.SerializeToString,
        response_deserializer=buildservice__pb2.BuildResponse.FromString,
        )
    self.DeleteEdge = channel.unary_unary(
        '/service.BuildService/DeleteEdge',
        request_serializer=buildservice__pb2.DeleteMessage.SerializeToString,
        response_deserializer=buildservice__pb2.BuildResponse.FromString,
        )


class BuildServiceServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Build(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SendBuildError(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def PushFile(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Package(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CreatePackage(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def CreateProject(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetProjectNode(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def DeleteNode(self, request_iterator, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def DeleteEdge(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_BuildServiceServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Build': grpc.stream_unary_rpc_method_handler(
          servicer.Build,
          request_deserializer=datamodel__pb2.FileNode.FromString,
          response_serializer=buildservice__pb2.BuildResponse.SerializeToString,
      ),
      'SendBuildError': grpc.unary_unary_rpc_method_handler(
          servicer.SendBuildError,
          request_deserializer=datamodel__pb2.InfoNode.FromString,
          response_serializer=buildservice__pb2.BuildResponse.SerializeToString,
      ),
      'PushFile': grpc.unary_unary_rpc_method_handler(
          servicer.PushFile,
          request_deserializer=buildservice__pb2.PushFileMessage.FromString,
          response_serializer=buildservice__pb2.PushFileResponse.SerializeToString,
      ),
      'Package': grpc.stream_unary_rpc_method_handler(
          servicer.Package,
          request_deserializer=datamodel__pb2.FileNode.FromString,
          response_serializer=buildservice__pb2.BuildResponse.SerializeToString,
      ),
      'CreatePackage': grpc.unary_unary_rpc_method_handler(
          servicer.CreatePackage,
          request_deserializer=datamodel__pb2.PackageNode.FromString,
          response_serializer=buildservice__pb2.BuildResponse.SerializeToString,
      ),
      'CreateProject': grpc.unary_unary_rpc_method_handler(
          servicer.CreateProject,
          request_deserializer=datamodel__pb2.ProjectNode.FromString,
          response_serializer=buildservice__pb2.BuildResponse.SerializeToString,
      ),
      'GetProjectNode': grpc.unary_unary_rpc_method_handler(
          servicer.GetProjectNode,
          request_deserializer=datamodel__pb2.ProjectNode.FromString,
          response_serializer=datamodel__pb2.ProjectNode.SerializeToString,
      ),
      'DeleteNode': grpc.stream_unary_rpc_method_handler(
          servicer.DeleteNode,
          request_deserializer=buildservice__pb2.DeleteMessage.FromString,
          response_serializer=buildservice__pb2.BuildResponse.SerializeToString,
      ),
      'DeleteEdge': grpc.unary_unary_rpc_method_handler(
          servicer.DeleteEdge,
          request_deserializer=buildservice__pb2.DeleteMessage.FromString,
          response_serializer=buildservice__pb2.BuildResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'service.BuildService', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
