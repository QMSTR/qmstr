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
    value = "by gRPC proto compiler (version 1.20.0)",
    comments = "Source: reportservice.proto")
public final class ReportServiceGrpc {

  private ReportServiceGrpc() {}

  public static final String SERVICE_NAME = "service.ReportService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Reportservice.ReporterConfigRequest,
      org.qmstr.grpc.service.Reportservice.ReporterConfigResponse> getGetReporterConfigMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetReporterConfig",
      requestType = org.qmstr.grpc.service.Reportservice.ReporterConfigRequest.class,
      responseType = org.qmstr.grpc.service.Reportservice.ReporterConfigResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Reportservice.ReporterConfigRequest,
      org.qmstr.grpc.service.Reportservice.ReporterConfigResponse> getGetReporterConfigMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Reportservice.ReporterConfigRequest, org.qmstr.grpc.service.Reportservice.ReporterConfigResponse> getGetReporterConfigMethod;
    if ((getGetReporterConfigMethod = ReportServiceGrpc.getGetReporterConfigMethod) == null) {
      synchronized (ReportServiceGrpc.class) {
        if ((getGetReporterConfigMethod = ReportServiceGrpc.getGetReporterConfigMethod) == null) {
          ReportServiceGrpc.getGetReporterConfigMethod = getGetReporterConfigMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Reportservice.ReporterConfigRequest, org.qmstr.grpc.service.Reportservice.ReporterConfigResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.ReportService", "GetReporterConfig"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Reportservice.ReporterConfigRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Reportservice.ReporterConfigResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new ReportServiceMethodDescriptorSupplier("GetReporterConfig"))
                  .build();
          }
        }
     }
     return getGetReporterConfigMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Reportservice.InfoDataRequest,
      org.qmstr.grpc.service.Reportservice.InfoDataResponse> getGetInfoDataMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetInfoData",
      requestType = org.qmstr.grpc.service.Reportservice.InfoDataRequest.class,
      responseType = org.qmstr.grpc.service.Reportservice.InfoDataResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Reportservice.InfoDataRequest,
      org.qmstr.grpc.service.Reportservice.InfoDataResponse> getGetInfoDataMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Reportservice.InfoDataRequest, org.qmstr.grpc.service.Reportservice.InfoDataResponse> getGetInfoDataMethod;
    if ((getGetInfoDataMethod = ReportServiceGrpc.getGetInfoDataMethod) == null) {
      synchronized (ReportServiceGrpc.class) {
        if ((getGetInfoDataMethod = ReportServiceGrpc.getGetInfoDataMethod) == null) {
          ReportServiceGrpc.getGetInfoDataMethod = getGetInfoDataMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Reportservice.InfoDataRequest, org.qmstr.grpc.service.Reportservice.InfoDataResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.ReportService", "GetInfoData"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Reportservice.InfoDataRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Reportservice.InfoDataResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new ReportServiceMethodDescriptorSupplier("GetInfoData"))
                  .build();
          }
        }
     }
     return getGetInfoDataMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.ProjectNode,
      org.qmstr.grpc.service.Datamodel.ProjectNode> getGetProjectNodeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetProjectNode",
      requestType = org.qmstr.grpc.service.Datamodel.ProjectNode.class,
      responseType = org.qmstr.grpc.service.Datamodel.ProjectNode.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.ProjectNode,
      org.qmstr.grpc.service.Datamodel.ProjectNode> getGetProjectNodeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.ProjectNode, org.qmstr.grpc.service.Datamodel.ProjectNode> getGetProjectNodeMethod;
    if ((getGetProjectNodeMethod = ReportServiceGrpc.getGetProjectNodeMethod) == null) {
      synchronized (ReportServiceGrpc.class) {
        if ((getGetProjectNodeMethod = ReportServiceGrpc.getGetProjectNodeMethod) == null) {
          ReportServiceGrpc.getGetProjectNodeMethod = getGetProjectNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Datamodel.ProjectNode, org.qmstr.grpc.service.Datamodel.ProjectNode>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.ReportService", "GetProjectNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.ProjectNode.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.ProjectNode.getDefaultInstance()))
                  .setSchemaDescriptor(new ReportServiceMethodDescriptorSupplier("GetProjectNode"))
                  .build();
          }
        }
     }
     return getGetProjectNodeMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static ReportServiceStub newStub(io.grpc.Channel channel) {
    return new ReportServiceStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static ReportServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new ReportServiceBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static ReportServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new ReportServiceFutureStub(channel);
  }

  /**
   */
  public static abstract class ReportServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public void getReporterConfig(org.qmstr.grpc.service.Reportservice.ReporterConfigRequest request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Reportservice.ReporterConfigResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getGetReporterConfigMethod(), responseObserver);
    }

    /**
     */
    public void getInfoData(org.qmstr.grpc.service.Reportservice.InfoDataRequest request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Reportservice.InfoDataResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getGetInfoDataMethod(), responseObserver);
    }

    /**
     */
    public void getProjectNode(org.qmstr.grpc.service.Datamodel.ProjectNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.ProjectNode> responseObserver) {
      asyncUnimplementedUnaryCall(getGetProjectNodeMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getGetReporterConfigMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Reportservice.ReporterConfigRequest,
                org.qmstr.grpc.service.Reportservice.ReporterConfigResponse>(
                  this, METHODID_GET_REPORTER_CONFIG)))
          .addMethod(
            getGetInfoDataMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Reportservice.InfoDataRequest,
                org.qmstr.grpc.service.Reportservice.InfoDataResponse>(
                  this, METHODID_GET_INFO_DATA)))
          .addMethod(
            getGetProjectNodeMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Datamodel.ProjectNode,
                org.qmstr.grpc.service.Datamodel.ProjectNode>(
                  this, METHODID_GET_PROJECT_NODE)))
          .build();
    }
  }

  /**
   */
  public static final class ReportServiceStub extends io.grpc.stub.AbstractStub<ReportServiceStub> {
    private ReportServiceStub(io.grpc.Channel channel) {
      super(channel);
    }

    private ReportServiceStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ReportServiceStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new ReportServiceStub(channel, callOptions);
    }

    /**
     */
    public void getReporterConfig(org.qmstr.grpc.service.Reportservice.ReporterConfigRequest request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Reportservice.ReporterConfigResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getGetReporterConfigMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getInfoData(org.qmstr.grpc.service.Reportservice.InfoDataRequest request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Reportservice.InfoDataResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getGetInfoDataMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getProjectNode(org.qmstr.grpc.service.Datamodel.ProjectNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.ProjectNode> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getGetProjectNodeMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class ReportServiceBlockingStub extends io.grpc.stub.AbstractStub<ReportServiceBlockingStub> {
    private ReportServiceBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private ReportServiceBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ReportServiceBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new ReportServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public org.qmstr.grpc.service.Reportservice.ReporterConfigResponse getReporterConfig(org.qmstr.grpc.service.Reportservice.ReporterConfigRequest request) {
      return blockingUnaryCall(
          getChannel(), getGetReporterConfigMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Reportservice.InfoDataResponse getInfoData(org.qmstr.grpc.service.Reportservice.InfoDataRequest request) {
      return blockingUnaryCall(
          getChannel(), getGetInfoDataMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Datamodel.ProjectNode getProjectNode(org.qmstr.grpc.service.Datamodel.ProjectNode request) {
      return blockingUnaryCall(
          getChannel(), getGetProjectNodeMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class ReportServiceFutureStub extends io.grpc.stub.AbstractStub<ReportServiceFutureStub> {
    private ReportServiceFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private ReportServiceFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected ReportServiceFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new ReportServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Reportservice.ReporterConfigResponse> getReporterConfig(
        org.qmstr.grpc.service.Reportservice.ReporterConfigRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getGetReporterConfigMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Reportservice.InfoDataResponse> getInfoData(
        org.qmstr.grpc.service.Reportservice.InfoDataRequest request) {
      return futureUnaryCall(
          getChannel().newCall(getGetInfoDataMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Datamodel.ProjectNode> getProjectNode(
        org.qmstr.grpc.service.Datamodel.ProjectNode request) {
      return futureUnaryCall(
          getChannel().newCall(getGetProjectNodeMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_REPORTER_CONFIG = 0;
  private static final int METHODID_GET_INFO_DATA = 1;
  private static final int METHODID_GET_PROJECT_NODE = 2;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final ReportServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(ReportServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_GET_REPORTER_CONFIG:
          serviceImpl.getReporterConfig((org.qmstr.grpc.service.Reportservice.ReporterConfigRequest) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Reportservice.ReporterConfigResponse>) responseObserver);
          break;
        case METHODID_GET_INFO_DATA:
          serviceImpl.getInfoData((org.qmstr.grpc.service.Reportservice.InfoDataRequest) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Reportservice.InfoDataResponse>) responseObserver);
          break;
        case METHODID_GET_PROJECT_NODE:
          serviceImpl.getProjectNode((org.qmstr.grpc.service.Datamodel.ProjectNode) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.ProjectNode>) responseObserver);
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

  private static abstract class ReportServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    ReportServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return org.qmstr.grpc.service.Reportservice.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("ReportService");
    }
  }

  private static final class ReportServiceFileDescriptorSupplier
      extends ReportServiceBaseDescriptorSupplier {
    ReportServiceFileDescriptorSupplier() {}
  }

  private static final class ReportServiceMethodDescriptorSupplier
      extends ReportServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    ReportServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (ReportServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new ReportServiceFileDescriptorSupplier())
              .addMethod(getGetReporterConfigMethod())
              .addMethod(getGetInfoDataMethod())
              .addMethod(getGetProjectNodeMethod())
              .build();
        }
      }
    }
    return result;
  }
}
