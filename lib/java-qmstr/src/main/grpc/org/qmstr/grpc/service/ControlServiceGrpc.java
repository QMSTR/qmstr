package org.qmstr.grpc.service;

import static io.grpc.MethodDescriptor.generateFullMethodName;
import static io.grpc.stub.ClientCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ClientCalls.asyncClientStreamingCall;
import static io.grpc.stub.ClientCalls.asyncServerStreamingCall;
import static io.grpc.stub.ClientCalls.asyncUnaryCall;
import static io.grpc.stub.ClientCalls.blockingServerStreamingCall;
import static io.grpc.stub.ClientCalls.blockingUnaryCall;
import static io.grpc.stub.ClientCalls.futureUnaryCall;
import static io.grpc.stub.ServerCalls.asyncBidiStreamingCall;
import static io.grpc.stub.ServerCalls.asyncClientStreamingCall;
import static io.grpc.stub.ServerCalls.asyncServerStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnaryCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedStreamingCall;
import static io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.18.0)",
    comments = "Source: controlservice.proto")
public final class ControlServiceGrpc {

  private ControlServiceGrpc() {}

  public static final String SERVICE_NAME = "service.ControlService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.LogMessage,
      org.qmstr.grpc.service.Controlservice.LogResponse> getLogMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Log",
      requestType = org.qmstr.grpc.service.Controlservice.LogMessage.class,
      responseType = org.qmstr.grpc.service.Controlservice.LogResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.LogMessage,
      org.qmstr.grpc.service.Controlservice.LogResponse> getLogMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.LogMessage, org.qmstr.grpc.service.Controlservice.LogResponse> getLogMethod;
    if ((getLogMethod = ControlServiceGrpc.getLogMethod) == null) {
      synchronized (ControlServiceGrpc.class) {
        if ((getLogMethod = ControlServiceGrpc.getLogMethod) == null) {
          ControlServiceGrpc.getLogMethod = getLogMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Controlservice.LogMessage, org.qmstr.grpc.service.Controlservice.LogResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.ControlService", "Log"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.LogMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.LogResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new ControlServiceMethodDescriptorSupplier("Log"))
                  .build();
          }
        }
     }
     return getLogMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.QuitMessage,
      org.qmstr.grpc.service.Controlservice.QuitResponse> getQuitMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Quit",
      requestType = org.qmstr.grpc.service.Controlservice.QuitMessage.class,
      responseType = org.qmstr.grpc.service.Controlservice.QuitResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.QuitMessage,
      org.qmstr.grpc.service.Controlservice.QuitResponse> getQuitMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.QuitMessage, org.qmstr.grpc.service.Controlservice.QuitResponse> getQuitMethod;
    if ((getQuitMethod = ControlServiceGrpc.getQuitMethod) == null) {
      synchronized (ControlServiceGrpc.class) {
        if ((getQuitMethod = ControlServiceGrpc.getQuitMethod) == null) {
          ControlServiceGrpc.getQuitMethod = getQuitMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Controlservice.QuitMessage, org.qmstr.grpc.service.Controlservice.QuitResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.ControlService", "Quit"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.QuitMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.QuitResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new ControlServiceMethodDescriptorSupplier("Quit"))
                  .build();
          }
        }
     }
     return getQuitMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage,
      org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse> getSwitchPhaseMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SwitchPhase",
      requestType = org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage.class,
      responseType = org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage,
      org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse> getSwitchPhaseMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage, org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse> getSwitchPhaseMethod;
    if ((getSwitchPhaseMethod = ControlServiceGrpc.getSwitchPhaseMethod) == null) {
      synchronized (ControlServiceGrpc.class) {
        if ((getSwitchPhaseMethod = ControlServiceGrpc.getSwitchPhaseMethod) == null) {
          ControlServiceGrpc.getSwitchPhaseMethod = getSwitchPhaseMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage, org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.ControlService", "SwitchPhase"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new ControlServiceMethodDescriptorSupplier("SwitchPhase"))
                  .build();
          }
        }
     }
     return getSwitchPhaseMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.PackageRequest,
      org.qmstr.grpc.service.Datamodel.PackageNode> getGetPackageNodeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetPackageNode",
      requestType = org.qmstr.grpc.service.Controlservice.PackageRequest.class,
      responseType = org.qmstr.grpc.service.Datamodel.PackageNode.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.PackageRequest,
      org.qmstr.grpc.service.Datamodel.PackageNode> getGetPackageNodeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.PackageRequest, org.qmstr.grpc.service.Datamodel.PackageNode> getGetPackageNodeMethod;
    if ((getGetPackageNodeMethod = ControlServiceGrpc.getGetPackageNodeMethod) == null) {
      synchronized (ControlServiceGrpc.class) {
        if ((getGetPackageNodeMethod = ControlServiceGrpc.getGetPackageNodeMethod) == null) {
          ControlServiceGrpc.getGetPackageNodeMethod = getGetPackageNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Controlservice.PackageRequest, org.qmstr.grpc.service.Datamodel.PackageNode>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.ControlService", "GetPackageNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.PackageRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.PackageNode.getDefaultInstance()))
                  .setSchemaDescriptor(new ControlServiceMethodDescriptorSupplier("GetPackageNode"))
                  .build();
          }
        }
     }
     return getGetPackageNodeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.FileNode,
      org.qmstr.grpc.service.Datamodel.FileNode> getGetFileNodeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetFileNode",
      requestType = org.qmstr.grpc.service.Datamodel.FileNode.class,
      responseType = org.qmstr.grpc.service.Datamodel.FileNode.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.FileNode,
      org.qmstr.grpc.service.Datamodel.FileNode> getGetFileNodeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.FileNode, org.qmstr.grpc.service.Datamodel.FileNode> getGetFileNodeMethod;
    if ((getGetFileNodeMethod = ControlServiceGrpc.getGetFileNodeMethod) == null) {
      synchronized (ControlServiceGrpc.class) {
        if ((getGetFileNodeMethod = ControlServiceGrpc.getGetFileNodeMethod) == null) {
          ControlServiceGrpc.getGetFileNodeMethod = getGetFileNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Datamodel.FileNode, org.qmstr.grpc.service.Datamodel.FileNode>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "service.ControlService", "GetFileNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.FileNode.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.FileNode.getDefaultInstance()))
                  .setSchemaDescriptor(new ControlServiceMethodDescriptorSupplier("GetFileNode"))
                  .build();
          }
        }
     }
     return getGetFileNodeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.DiagnosticNode,
      org.qmstr.grpc.service.Datamodel.DiagnosticNode> getGetDiagnosticNodeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetDiagnosticNode",
      requestType = org.qmstr.grpc.service.Datamodel.DiagnosticNode.class,
      responseType = org.qmstr.grpc.service.Datamodel.DiagnosticNode.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.DiagnosticNode,
      org.qmstr.grpc.service.Datamodel.DiagnosticNode> getGetDiagnosticNodeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.DiagnosticNode, org.qmstr.grpc.service.Datamodel.DiagnosticNode> getGetDiagnosticNodeMethod;
    if ((getGetDiagnosticNodeMethod = ControlServiceGrpc.getGetDiagnosticNodeMethod) == null) {
      synchronized (ControlServiceGrpc.class) {
        if ((getGetDiagnosticNodeMethod = ControlServiceGrpc.getGetDiagnosticNodeMethod) == null) {
          ControlServiceGrpc.getGetDiagnosticNodeMethod = getGetDiagnosticNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Datamodel.DiagnosticNode, org.qmstr.grpc.service.Datamodel.DiagnosticNode>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "service.ControlService", "GetDiagnosticNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.DiagnosticNode.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.DiagnosticNode.getDefaultInstance()))
                  .setSchemaDescriptor(new ControlServiceMethodDescriptorSupplier("GetDiagnosticNode"))
                  .build();
          }
        }
     }
     return getGetDiagnosticNodeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.StatusMessage,
      org.qmstr.grpc.service.Controlservice.StatusResponse> getStatusMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Status",
      requestType = org.qmstr.grpc.service.Controlservice.StatusMessage.class,
      responseType = org.qmstr.grpc.service.Controlservice.StatusResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.StatusMessage,
      org.qmstr.grpc.service.Controlservice.StatusResponse> getStatusMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.StatusMessage, org.qmstr.grpc.service.Controlservice.StatusResponse> getStatusMethod;
    if ((getStatusMethod = ControlServiceGrpc.getStatusMethod) == null) {
      synchronized (ControlServiceGrpc.class) {
        if ((getStatusMethod = ControlServiceGrpc.getStatusMethod) == null) {
          ControlServiceGrpc.getStatusMethod = getStatusMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Controlservice.StatusMessage, org.qmstr.grpc.service.Controlservice.StatusResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.ControlService", "Status"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.StatusMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.StatusResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new ControlServiceMethodDescriptorSupplier("Status"))
                  .build();
          }
        }
     }
     return getStatusMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.EventMessage,
      org.qmstr.grpc.service.Datamodel.Event> getSubscribeEventsMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SubscribeEvents",
      requestType = org.qmstr.grpc.service.Controlservice.EventMessage.class,
      responseType = org.qmstr.grpc.service.Datamodel.Event.class,
      methodType = io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.EventMessage,
      org.qmstr.grpc.service.Datamodel.Event> getSubscribeEventsMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.EventMessage, org.qmstr.grpc.service.Datamodel.Event> getSubscribeEventsMethod;
    if ((getSubscribeEventsMethod = ControlServiceGrpc.getSubscribeEventsMethod) == null) {
      synchronized (ControlServiceGrpc.class) {
        if ((getSubscribeEventsMethod = ControlServiceGrpc.getSubscribeEventsMethod) == null) {
          ControlServiceGrpc.getSubscribeEventsMethod = getSubscribeEventsMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Controlservice.EventMessage, org.qmstr.grpc.service.Datamodel.Event>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.SERVER_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "service.ControlService", "SubscribeEvents"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.EventMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.Event.getDefaultInstance()))
                  .setSchemaDescriptor(new ControlServiceMethodDescriptorSupplier("SubscribeEvents"))
                  .build();
          }
        }
     }
     return getSubscribeEventsMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.ExportRequest,
      org.qmstr.grpc.service.Controlservice.ExportResponse> getExportSnapshotMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "ExportSnapshot",
      requestType = org.qmstr.grpc.service.Controlservice.ExportRequest.class,
      responseType = org.qmstr.grpc.service.Controlservice.ExportResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.ExportRequest,
      org.qmstr.grpc.service.Controlservice.ExportResponse> getExportSnapshotMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Controlservice.ExportRequest, org.qmstr.grpc.service.Controlservice.ExportResponse> getExportSnapshotMethod;
    if ((getExportSnapshotMethod = ControlServiceGrpc.getExportSnapshotMethod) == null) {
      synchronized (ControlServiceGrpc.class) {
        if ((getExportSnapshotMethod = ControlServiceGrpc.getExportSnapshotMethod) == null) {
          ControlServiceGrpc.getExportSnapshotMethod = getExportSnapshotMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Controlservice.ExportRequest, org.qmstr.grpc.service.Controlservice.ExportResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.ControlService", "ExportSnapshot"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.ExportRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Controlservice.ExportResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new ControlServiceMethodDescriptorSupplier("ExportSnapshot"))
                  .build();
          }
        }
     }
     return getExportSnapshotMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static ControlServiceStub newStub(io.grpc.Channel channel) {
    return new ControlServiceStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static ControlServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new ControlServiceBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static ControlServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new ControlServiceFutureStub(channel);
  }

  /**
   */
  public static abstract class ControlServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void log(org.qmstr.grpc.service.Controlservice.LogMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.LogResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getLogMethod(), responseObserver);
    }

    /**
     */
    public void quit(org.qmstr.grpc.service.Controlservice.QuitMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.QuitResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getQuitMethod(), responseObserver);
    }

    /**
     */
    public void switchPhase(org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getSwitchPhaseMethod(), responseObserver);
    }

    /**
     */
    public void getPackageNode(org.qmstr.grpc.service.Controlservice.PackageRequest request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.PackageNode> responseObserver) {
      asyncUnimplementedUnaryCall(getGetPackageNodeMethod(), responseObserver);
    }

    /**
     */
    public void getFileNode(org.qmstr.grpc.service.Datamodel.FileNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.FileNode> responseObserver) {
      asyncUnimplementedUnaryCall(getGetFileNodeMethod(), responseObserver);
    }

    /**
     */
    public void getDiagnosticNode(org.qmstr.grpc.service.Datamodel.DiagnosticNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.DiagnosticNode> responseObserver) {
      asyncUnimplementedUnaryCall(getGetDiagnosticNodeMethod(), responseObserver);
    }

    /**
     */
    public void status(org.qmstr.grpc.service.Controlservice.StatusMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.StatusResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getStatusMethod(), responseObserver);
    }

    /**
     */
    public void subscribeEvents(org.qmstr.grpc.service.Controlservice.EventMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.Event> responseObserver) {
      asyncUnimplementedUnaryCall(getSubscribeEventsMethod(), responseObserver);
    }

    /**
     */
    public void exportSnapshot(org.qmstr.grpc.service.Controlservice.ExportRequest request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.ExportResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getExportSnapshotMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getLogMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Controlservice.LogMessage,
                org.qmstr.grpc.service.Controlservice.LogResponse>(
                  this, METHODID_LOG)))
          .addMethod(
            getQuitMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Controlservice.QuitMessage,
                org.qmstr.grpc.service.Controlservice.QuitResponse>(
                  this, METHODID_QUIT)))
          .addMethod(
            getSwitchPhaseMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage,
                org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse>(
                  this, METHODID_SWITCH_PHASE)))
          .addMethod(
            getGetPackageNodeMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Controlservice.PackageRequest,
                org.qmstr.grpc.service.Datamodel.PackageNode>(
                  this, METHODID_GET_PACKAGE_NODE)))
          .addMethod(
            getGetFileNodeMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Datamodel.FileNode,
                org.qmstr.grpc.service.Datamodel.FileNode>(
                  this, METHODID_GET_FILE_NODE)))
          .addMethod(
            getGetDiagnosticNodeMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Datamodel.DiagnosticNode,
                org.qmstr.grpc.service.Datamodel.DiagnosticNode>(
                  this, METHODID_GET_DIAGNOSTIC_NODE)))
          .addMethod(
            getStatusMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Controlservice.StatusMessage,
                org.qmstr.grpc.service.Controlservice.StatusResponse>(
                  this, METHODID_STATUS)))
          .addMethod(
            getSubscribeEventsMethod(),
            asyncServerStreamingCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Controlservice.EventMessage,
                org.qmstr.grpc.service.Datamodel.Event>(
                  this, METHODID_SUBSCRIBE_EVENTS)))
          .addMethod(
            getExportSnapshotMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Controlservice.ExportRequest,
                org.qmstr.grpc.service.Controlservice.ExportResponse>(
                  this, METHODID_EXPORT_SNAPSHOT)))
          .build();
    }
  }

  /**
   */
  public static final class ControlServiceStub extends io.grpc.stub.AbstractStub<ControlServiceStub> {
    private ControlServiceStub(io.grpc.Channel channel) {
      super(channel);
    }

    private ControlServiceStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ControlServiceStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new ControlServiceStub(channel, callOptions);
    }

    /**
     */
    public void log(org.qmstr.grpc.service.Controlservice.LogMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.LogResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getLogMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void quit(org.qmstr.grpc.service.Controlservice.QuitMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.QuitResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getQuitMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void switchPhase(org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getSwitchPhaseMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getPackageNode(org.qmstr.grpc.service.Controlservice.PackageRequest request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.PackageNode> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getGetPackageNodeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getFileNode(org.qmstr.grpc.service.Datamodel.FileNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.FileNode> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getGetFileNodeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getDiagnosticNode(org.qmstr.grpc.service.Datamodel.DiagnosticNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.DiagnosticNode> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getGetDiagnosticNodeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void status(org.qmstr.grpc.service.Controlservice.StatusMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.StatusResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getStatusMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void subscribeEvents(org.qmstr.grpc.service.Controlservice.EventMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.Event> responseObserver) {
      asyncServerStreamingCall(
          getChannel().newCall(getSubscribeEventsMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void exportSnapshot(org.qmstr.grpc.service.Controlservice.ExportRequest request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.ExportResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getExportSnapshotMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class ControlServiceBlockingStub extends io.grpc.stub.AbstractStub<ControlServiceBlockingStub> {
    private ControlServiceBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private ControlServiceBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ControlServiceBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new ControlServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public org.qmstr.grpc.service.Controlservice.LogResponse log(org.qmstr.grpc.service.Controlservice.LogMessage request) {
      return blockingUnaryCall(
          getChannel(), getLogMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Controlservice.QuitResponse quit(org.qmstr.grpc.service.Controlservice.QuitMessage request) {
      return blockingUnaryCall(
          getChannel(), getQuitMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse switchPhase(org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage request) {
      return blockingUnaryCall(
          getChannel(), getSwitchPhaseMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Datamodel.PackageNode getPackageNode(org.qmstr.grpc.service.Controlservice.PackageRequest request) {
      return blockingUnaryCall(
          getChannel(), getGetPackageNodeMethod(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<org.qmstr.grpc.service.Datamodel.FileNode> getFileNode(
        org.qmstr.grpc.service.Datamodel.FileNode request) {
      return blockingServerStreamingCall(
          getChannel(), getGetFileNodeMethod(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<org.qmstr.grpc.service.Datamodel.DiagnosticNode> getDiagnosticNode(
        org.qmstr.grpc.service.Datamodel.DiagnosticNode request) {
      return blockingServerStreamingCall(
          getChannel(), getGetDiagnosticNodeMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Controlservice.StatusResponse status(org.qmstr.grpc.service.Controlservice.StatusMessage request) {
      return blockingUnaryCall(
          getChannel(), getStatusMethod(), getCallOptions(), request);
    }

    /**
     */
    public java.util.Iterator<org.qmstr.grpc.service.Datamodel.Event> subscribeEvents(
        org.qmstr.grpc.service.Controlservice.EventMessage request) {
      return blockingServerStreamingCall(
          getChannel(), getSubscribeEventsMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Controlservice.ExportResponse exportSnapshot(org.qmstr.grpc.service.Controlservice.ExportRequest request) {
      return blockingUnaryCall(
          getChannel(), getExportSnapshotMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class ControlServiceFutureStub extends io.grpc.stub.AbstractStub<ControlServiceFutureStub> {
    private ControlServiceFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private ControlServiceFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ControlServiceFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new ControlServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Controlservice.LogResponse> log(
        org.qmstr.grpc.service.Controlservice.LogMessage request) {
      return futureUnaryCall(
          getChannel().newCall(getLogMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Controlservice.QuitResponse> quit(
        org.qmstr.grpc.service.Controlservice.QuitMessage request) {
      return futureUnaryCall(
          getChannel().newCall(getQuitMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse> switchPhase(
        org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage request) {
      return futureUnaryCall(
          getChannel().newCall(getSwitchPhaseMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Datamodel.PackageNode> getPackageNode(
        org.qmstr.grpc.service.Controlservice.PackageRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getGetPackageNodeMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Controlservice.StatusResponse> status(
        org.qmstr.grpc.service.Controlservice.StatusMessage request) {
      return futureUnaryCall(
          getChannel().newCall(getStatusMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Controlservice.ExportResponse> exportSnapshot(
        org.qmstr.grpc.service.Controlservice.ExportRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getExportSnapshotMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_LOG = 0;
  private static final int METHODID_QUIT = 1;
  private static final int METHODID_SWITCH_PHASE = 2;
  private static final int METHODID_GET_PACKAGE_NODE = 3;
  private static final int METHODID_GET_FILE_NODE = 4;
  private static final int METHODID_GET_DIAGNOSTIC_NODE = 5;
  private static final int METHODID_STATUS = 6;
  private static final int METHODID_SUBSCRIBE_EVENTS = 7;
  private static final int METHODID_EXPORT_SNAPSHOT = 8;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final ControlServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(ControlServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_LOG:
          serviceImpl.log((org.qmstr.grpc.service.Controlservice.LogMessage) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.LogResponse>) responseObserver);
          break;
        case METHODID_QUIT:
          serviceImpl.quit((org.qmstr.grpc.service.Controlservice.QuitMessage) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.QuitResponse>) responseObserver);
          break;
        case METHODID_SWITCH_PHASE:
          serviceImpl.switchPhase((org.qmstr.grpc.service.Controlservice.SwitchPhaseMessage) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.SwitchPhaseResponse>) responseObserver);
          break;
        case METHODID_GET_PACKAGE_NODE:
          serviceImpl.getPackageNode((org.qmstr.grpc.service.Controlservice.PackageRequest) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.PackageNode>) responseObserver);
          break;
        case METHODID_GET_FILE_NODE:
          serviceImpl.getFileNode((org.qmstr.grpc.service.Datamodel.FileNode) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.FileNode>) responseObserver);
          break;
        case METHODID_GET_DIAGNOSTIC_NODE:
          serviceImpl.getDiagnosticNode((org.qmstr.grpc.service.Datamodel.DiagnosticNode) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.DiagnosticNode>) responseObserver);
          break;
        case METHODID_STATUS:
          serviceImpl.status((org.qmstr.grpc.service.Controlservice.StatusMessage) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.StatusResponse>) responseObserver);
          break;
        case METHODID_SUBSCRIBE_EVENTS:
          serviceImpl.subscribeEvents((org.qmstr.grpc.service.Controlservice.EventMessage) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.Event>) responseObserver);
          break;
        case METHODID_EXPORT_SNAPSHOT:
          serviceImpl.exportSnapshot((org.qmstr.grpc.service.Controlservice.ExportRequest) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Controlservice.ExportResponse>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class ControlServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    ControlServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return org.qmstr.grpc.service.Controlservice.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("ControlService");
    }
  }

  private static final class ControlServiceFileDescriptorSupplier
      extends ControlServiceBaseDescriptorSupplier {
    ControlServiceFileDescriptorSupplier() {}
  }

  private static final class ControlServiceMethodDescriptorSupplier
      extends ControlServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    ControlServiceMethodDescriptorSupplier(String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (ControlServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new ControlServiceFileDescriptorSupplier())
              .addMethod(getLogMethod())
              .addMethod(getQuitMethod())
              .addMethod(getSwitchPhaseMethod())
              .addMethod(getGetPackageNodeMethod())
              .addMethod(getGetFileNodeMethod())
              .addMethod(getGetDiagnosticNodeMethod())
              .addMethod(getStatusMethod())
              .addMethod(getSubscribeEventsMethod())
              .addMethod(getExportSnapshotMethod())
              .build();
        }
      }
    }
    return result;
  }
}
