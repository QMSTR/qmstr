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
    comments = "Source: analyzerservice.proto")
public final class AnalysisServiceGrpc {

  private AnalysisServiceGrpc() {}

  public static final String SERVICE_NAME = "service.AnalysisService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest,
      org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse> getGetAnalyzerConfigMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetAnalyzerConfig",
      requestType = org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest.class,
      responseType = org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest,
      org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse> getGetAnalyzerConfigMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest, org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse> getGetAnalyzerConfigMethod;
    if ((getGetAnalyzerConfigMethod = AnalysisServiceGrpc.getGetAnalyzerConfigMethod) == null) {
      synchronized (AnalysisServiceGrpc.class) {
        if ((getGetAnalyzerConfigMethod = AnalysisServiceGrpc.getGetAnalyzerConfigMethod) == null) {
          AnalysisServiceGrpc.getGetAnalyzerConfigMethod = getGetAnalyzerConfigMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest, org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.AnalysisService", "GetAnalyzerConfig"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new AnalysisServiceMethodDescriptorSupplier("GetAnalyzerConfig"))
                  .build();
          }
        }
     }
     return getGetAnalyzerConfigMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.InfoNodeMessage,
      org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendInfoNodesMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SendInfoNodes",
      requestType = org.qmstr.grpc.service.Analyzerservice.InfoNodeMessage.class,
      responseType = org.qmstr.grpc.service.Analyzerservice.SendResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.InfoNodeMessage,
      org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendInfoNodesMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.InfoNodeMessage, org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendInfoNodesMethod;
    if ((getSendInfoNodesMethod = AnalysisServiceGrpc.getSendInfoNodesMethod) == null) {
      synchronized (AnalysisServiceGrpc.class) {
        if ((getSendInfoNodesMethod = AnalysisServiceGrpc.getSendInfoNodesMethod) == null) {
          AnalysisServiceGrpc.getSendInfoNodesMethod = getSendInfoNodesMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Analyzerservice.InfoNodeMessage, org.qmstr.grpc.service.Analyzerservice.SendResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "service.AnalysisService", "SendInfoNodes"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Analyzerservice.InfoNodeMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Analyzerservice.SendResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new AnalysisServiceMethodDescriptorSupplier("SendInfoNodes"))
                  .build();
          }
        }
     }
     return getSendInfoNodesMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.FileNodeMessage,
      org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendFileNodeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SendFileNode",
      requestType = org.qmstr.grpc.service.Analyzerservice.FileNodeMessage.class,
      responseType = org.qmstr.grpc.service.Analyzerservice.SendResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.FileNodeMessage,
      org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendFileNodeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.FileNodeMessage, org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendFileNodeMethod;
    if ((getSendFileNodeMethod = AnalysisServiceGrpc.getSendFileNodeMethod) == null) {
      synchronized (AnalysisServiceGrpc.class) {
        if ((getSendFileNodeMethod = AnalysisServiceGrpc.getSendFileNodeMethod) == null) {
          AnalysisServiceGrpc.getSendFileNodeMethod = getSendFileNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Analyzerservice.FileNodeMessage, org.qmstr.grpc.service.Analyzerservice.SendResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "service.AnalysisService", "SendFileNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Analyzerservice.FileNodeMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Analyzerservice.SendResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new AnalysisServiceMethodDescriptorSupplier("SendFileNode"))
                  .build();
          }
        }
     }
     return getSendFileNodeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.PackageNodeMessage,
      org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendPackageNodeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SendPackageNode",
      requestType = org.qmstr.grpc.service.Analyzerservice.PackageNodeMessage.class,
      responseType = org.qmstr.grpc.service.Analyzerservice.SendResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.PackageNodeMessage,
      org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendPackageNodeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.PackageNodeMessage, org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendPackageNodeMethod;
    if ((getSendPackageNodeMethod = AnalysisServiceGrpc.getSendPackageNodeMethod) == null) {
      synchronized (AnalysisServiceGrpc.class) {
        if ((getSendPackageNodeMethod = AnalysisServiceGrpc.getSendPackageNodeMethod) == null) {
          AnalysisServiceGrpc.getSendPackageNodeMethod = getSendPackageNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Analyzerservice.PackageNodeMessage, org.qmstr.grpc.service.Analyzerservice.SendResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "service.AnalysisService", "SendPackageNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Analyzerservice.PackageNodeMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Analyzerservice.SendResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new AnalysisServiceMethodDescriptorSupplier("SendPackageNode"))
                  .build();
          }
        }
     }
     return getSendPackageNodeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.DiagnosticNodeMessage,
      org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendDiagnosticNodeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SendDiagnosticNode",
      requestType = org.qmstr.grpc.service.Analyzerservice.DiagnosticNodeMessage.class,
      responseType = org.qmstr.grpc.service.Analyzerservice.SendResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.DiagnosticNodeMessage,
      org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendDiagnosticNodeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Analyzerservice.DiagnosticNodeMessage, org.qmstr.grpc.service.Analyzerservice.SendResponse> getSendDiagnosticNodeMethod;
    if ((getSendDiagnosticNodeMethod = AnalysisServiceGrpc.getSendDiagnosticNodeMethod) == null) {
      synchronized (AnalysisServiceGrpc.class) {
        if ((getSendDiagnosticNodeMethod = AnalysisServiceGrpc.getSendDiagnosticNodeMethod) == null) {
          AnalysisServiceGrpc.getSendDiagnosticNodeMethod = getSendDiagnosticNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Analyzerservice.DiagnosticNodeMessage, org.qmstr.grpc.service.Analyzerservice.SendResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "service.AnalysisService", "SendDiagnosticNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Analyzerservice.DiagnosticNodeMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Analyzerservice.SendResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new AnalysisServiceMethodDescriptorSupplier("SendDiagnosticNode"))
                  .build();
          }
        }
     }
     return getSendDiagnosticNodeMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static AnalysisServiceStub newStub(io.grpc.Channel channel) {
    return new AnalysisServiceStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static AnalysisServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new AnalysisServiceBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static AnalysisServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new AnalysisServiceFutureStub(channel);
  }

  /**
   */
  public static abstract class AnalysisServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void getAnalyzerConfig(org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getGetAnalyzerConfigMethod(), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.InfoNodeMessage> sendInfoNodes(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getSendInfoNodesMethod(), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.FileNodeMessage> sendFileNode(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getSendFileNodeMethod(), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.PackageNodeMessage> sendPackageNode(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getSendPackageNodeMethod(), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.DiagnosticNodeMessage> sendDiagnosticNode(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getSendDiagnosticNodeMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetAnalyzerConfigMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest,
                org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse>(
                  this, METHODID_GET_ANALYZER_CONFIG)))
          .addMethod(
            getSendInfoNodesMethod(),
            asyncClientStreamingCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Analyzerservice.InfoNodeMessage,
                org.qmstr.grpc.service.Analyzerservice.SendResponse>(
                  this, METHODID_SEND_INFO_NODES)))
          .addMethod(
            getSendFileNodeMethod(),
            asyncClientStreamingCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Analyzerservice.FileNodeMessage,
                org.qmstr.grpc.service.Analyzerservice.SendResponse>(
                  this, METHODID_SEND_FILE_NODE)))
          .addMethod(
            getSendPackageNodeMethod(),
            asyncClientStreamingCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Analyzerservice.PackageNodeMessage,
                org.qmstr.grpc.service.Analyzerservice.SendResponse>(
                  this, METHODID_SEND_PACKAGE_NODE)))
          .addMethod(
            getSendDiagnosticNodeMethod(),
            asyncClientStreamingCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Analyzerservice.DiagnosticNodeMessage,
                org.qmstr.grpc.service.Analyzerservice.SendResponse>(
                  this, METHODID_SEND_DIAGNOSTIC_NODE)))
          .build();
    }
  }

  /**
   */
  public static final class AnalysisServiceStub extends io.grpc.stub.AbstractStub<AnalysisServiceStub> {
    private AnalysisServiceStub(io.grpc.Channel channel) {
      super(channel);
    }

    private AnalysisServiceStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AnalysisServiceStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new AnalysisServiceStub(channel, callOptions);
    }

    /**
     */
    public void getAnalyzerConfig(org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getGetAnalyzerConfigMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.InfoNodeMessage> sendInfoNodes(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse> responseObserver) {
      return asyncClientStreamingCall(
          getChannel().newCall(getSendInfoNodesMethod(), getCallOptions()), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.FileNodeMessage> sendFileNode(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse> responseObserver) {
      return asyncClientStreamingCall(
          getChannel().newCall(getSendFileNodeMethod(), getCallOptions()), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.PackageNodeMessage> sendPackageNode(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse> responseObserver) {
      return asyncClientStreamingCall(
          getChannel().newCall(getSendPackageNodeMethod(), getCallOptions()), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.DiagnosticNodeMessage> sendDiagnosticNode(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse> responseObserver) {
      return asyncClientStreamingCall(
          getChannel().newCall(getSendDiagnosticNodeMethod(), getCallOptions()), responseObserver);
    }
  }

  /**
   */
  public static final class AnalysisServiceBlockingStub extends io.grpc.stub.AbstractStub<AnalysisServiceBlockingStub> {
    private AnalysisServiceBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private AnalysisServiceBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AnalysisServiceBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new AnalysisServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse getAnalyzerConfig(org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest request) {
      return blockingUnaryCall(
          getChannel(), getGetAnalyzerConfigMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class AnalysisServiceFutureStub extends io.grpc.stub.AbstractStub<AnalysisServiceFutureStub> {
    private AnalysisServiceFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private AnalysisServiceFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AnalysisServiceFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new AnalysisServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse> getAnalyzerConfig(
        org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getGetAnalyzerConfigMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_ANALYZER_CONFIG = 0;
  private static final int METHODID_SEND_INFO_NODES = 1;
  private static final int METHODID_SEND_FILE_NODE = 2;
  private static final int METHODID_SEND_PACKAGE_NODE = 3;
  private static final int METHODID_SEND_DIAGNOSTIC_NODE = 4;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AnalysisServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(AnalysisServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_ANALYZER_CONFIG:
          serviceImpl.getAnalyzerConfig((org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigRequest) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.AnalyzerConfigResponse>) responseObserver);
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
        case METHODID_SEND_INFO_NODES:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.sendInfoNodes(
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse>) responseObserver);
        case METHODID_SEND_FILE_NODE:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.sendFileNode(
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse>) responseObserver);
        case METHODID_SEND_PACKAGE_NODE:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.sendPackageNode(
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse>) responseObserver);
        case METHODID_SEND_DIAGNOSTIC_NODE:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.sendDiagnosticNode(
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Analyzerservice.SendResponse>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class AnalysisServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    AnalysisServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return org.qmstr.grpc.service.Analyzerservice.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("AnalysisService");
    }
  }

  private static final class AnalysisServiceFileDescriptorSupplier
      extends AnalysisServiceBaseDescriptorSupplier {
    AnalysisServiceFileDescriptorSupplier() {}
  }

  private static final class AnalysisServiceMethodDescriptorSupplier
      extends AnalysisServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    AnalysisServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (AnalysisServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new AnalysisServiceFileDescriptorSupplier())
              .addMethod(getGetAnalyzerConfigMethod())
              .addMethod(getSendInfoNodesMethod())
              .addMethod(getSendFileNodeMethod())
              .addMethod(getSendPackageNodeMethod())
              .addMethod(getSendDiagnosticNodeMethod())
              .build();
        }
      }
    }
    return result;
  }
}
