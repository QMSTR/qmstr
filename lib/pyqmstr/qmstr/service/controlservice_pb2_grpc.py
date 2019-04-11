# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

from . import controlservice_pb2 as controlservice__pb2
from . import datamodel_pb2 as datamodel__pb2


class ControlServiceStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.Log = channel.unary_unary(
        '/service.ControlService/Log',
        request_serializer=controlservice__pb2.LogMessage.SerializeToString,
        response_deserializer=controlservice__pb2.LogResponse.FromString,
        )
    self.Quit = channel.unary_unary(
        '/service.ControlService/Quit',
        request_serializer=controlservice__pb2.QuitMessage.SerializeToString,
        response_deserializer=controlservice__pb2.QuitResponse.FromString,
        )
    self.SwitchPhase = channel.unary_unary(
        '/service.ControlService/SwitchPhase',
        request_serializer=controlservice__pb2.SwitchPhaseMessage.SerializeToString,
        response_deserializer=controlservice__pb2.SwitchPhaseResponse.FromString,
        )
    self.GetPackageNode = channel.unary_unary(
        '/service.ControlService/GetPackageNode',
        request_serializer=datamodel__pb2.PackageNode.SerializeToString,
        response_deserializer=datamodel__pb2.PackageNode.FromString,
        )
    self.GetFileNode = channel.unary_stream(
        '/service.ControlService/GetFileNode',
        request_serializer=controlservice__pb2.GetFileNodeMessage.SerializeToString,
        response_deserializer=datamodel__pb2.FileNode.FromString,
        )
    self.GetDiagnosticNode = channel.unary_stream(
        '/service.ControlService/GetDiagnosticNode',
        request_serializer=datamodel__pb2.DiagnosticNode.SerializeToString,
        response_deserializer=datamodel__pb2.DiagnosticNode.FromString,
        )
    self.Status = channel.unary_unary(
        '/service.ControlService/Status',
        request_serializer=controlservice__pb2.StatusMessage.SerializeToString,
        response_deserializer=controlservice__pb2.StatusResponse.FromString,
        )
    self.SubscribeEvents = channel.unary_stream(
        '/service.ControlService/SubscribeEvents',
        request_serializer=controlservice__pb2.EventMessage.SerializeToString,
        response_deserializer=datamodel__pb2.Event.FromString,
        )
    self.ExportSnapshot = channel.unary_unary(
        '/service.ControlService/ExportSnapshot',
        request_serializer=controlservice__pb2.ExportRequest.SerializeToString,
        response_deserializer=controlservice__pb2.ExportResponse.FromString,
        )


class ControlServiceServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def Log(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Quit(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SwitchPhase(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetPackageNode(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetFileNode(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def GetDiagnosticNode(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def Status(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def SubscribeEvents(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')

  def ExportSnapshot(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_ControlServiceServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'Log': grpc.unary_unary_rpc_method_handler(
          servicer.Log,
          request_deserializer=controlservice__pb2.LogMessage.FromString,
          response_serializer=controlservice__pb2.LogResponse.SerializeToString,
      ),
      'Quit': grpc.unary_unary_rpc_method_handler(
          servicer.Quit,
          request_deserializer=controlservice__pb2.QuitMessage.FromString,
          response_serializer=controlservice__pb2.QuitResponse.SerializeToString,
      ),
      'SwitchPhase': grpc.unary_unary_rpc_method_handler(
          servicer.SwitchPhase,
          request_deserializer=controlservice__pb2.SwitchPhaseMessage.FromString,
          response_serializer=controlservice__pb2.SwitchPhaseResponse.SerializeToString,
      ),
      'GetPackageNode': grpc.unary_unary_rpc_method_handler(
          servicer.GetPackageNode,
          request_deserializer=datamodel__pb2.PackageNode.FromString,
          response_serializer=datamodel__pb2.PackageNode.SerializeToString,
      ),
      'GetFileNode': grpc.unary_stream_rpc_method_handler(
          servicer.GetFileNode,
          request_deserializer=controlservice__pb2.GetFileNodeMessage.FromString,
          response_serializer=datamodel__pb2.FileNode.SerializeToString,
      ),
      'GetDiagnosticNode': grpc.unary_stream_rpc_method_handler(
          servicer.GetDiagnosticNode,
          request_deserializer=datamodel__pb2.DiagnosticNode.FromString,
          response_serializer=datamodel__pb2.DiagnosticNode.SerializeToString,
      ),
      'Status': grpc.unary_unary_rpc_method_handler(
          servicer.Status,
          request_deserializer=controlservice__pb2.StatusMessage.FromString,
          response_serializer=controlservice__pb2.StatusResponse.SerializeToString,
      ),
      'SubscribeEvents': grpc.unary_stream_rpc_method_handler(
          servicer.SubscribeEvents,
          request_deserializer=controlservice__pb2.EventMessage.FromString,
          response_serializer=datamodel__pb2.Event.SerializeToString,
      ),
      'ExportSnapshot': grpc.unary_unary_rpc_method_handler(
          servicer.ExportSnapshot,
          request_deserializer=controlservice__pb2.ExportRequest.FromString,
          response_serializer=controlservice__pb2.ExportResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'service.ControlService', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
