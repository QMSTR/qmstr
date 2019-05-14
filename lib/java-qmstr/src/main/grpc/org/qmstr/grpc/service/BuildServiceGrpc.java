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
    comments = "Source: buildservice.proto")
public final class BuildServiceGrpc {

  private BuildServiceGrpc() {}

  public static final String SERVICE_NAME = "service.BuildService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.FileNode,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getBuildMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "Build",
      requestType = org.qmstr.grpc.service.Datamodel.FileNode.class,
      responseType = org.qmstr.grpc.service.Buildservice.BuildResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.FileNode,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getBuildMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.FileNode, org.qmstr.grpc.service.Buildservice.BuildResponse> getBuildMethod;
    if ((getBuildMethod = BuildServiceGrpc.getBuildMethod) == null) {
      synchronized (BuildServiceGrpc.class) {
        if ((getBuildMethod = BuildServiceGrpc.getBuildMethod) == null) {
          BuildServiceGrpc.getBuildMethod = getBuildMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Datamodel.FileNode, org.qmstr.grpc.service.Buildservice.BuildResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "service.BuildService", "Build"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.FileNode.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.BuildResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new BuildServiceMethodDescriptorSupplier("Build"))
                  .build();
          }
        }
     }
     return getBuildMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.InfoNode,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getSendBuildErrorMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "SendBuildError",
      requestType = org.qmstr.grpc.service.Datamodel.InfoNode.class,
      responseType = org.qmstr.grpc.service.Buildservice.BuildResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.InfoNode,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getSendBuildErrorMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.InfoNode, org.qmstr.grpc.service.Buildservice.BuildResponse> getSendBuildErrorMethod;
    if ((getSendBuildErrorMethod = BuildServiceGrpc.getSendBuildErrorMethod) == null) {
      synchronized (BuildServiceGrpc.class) {
        if ((getSendBuildErrorMethod = BuildServiceGrpc.getSendBuildErrorMethod) == null) {
          BuildServiceGrpc.getSendBuildErrorMethod = getSendBuildErrorMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Datamodel.InfoNode, org.qmstr.grpc.service.Buildservice.BuildResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.BuildService", "SendBuildError"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.InfoNode.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.BuildResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new BuildServiceMethodDescriptorSupplier("SendBuildError"))
                  .build();
          }
        }
     }
     return getSendBuildErrorMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.PushFileMessage,
      org.qmstr.grpc.service.Buildservice.PushFileResponse> getPushFileMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "PushFile",
      requestType = org.qmstr.grpc.service.Buildservice.PushFileMessage.class,
      responseType = org.qmstr.grpc.service.Buildservice.PushFileResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.PushFileMessage,
      org.qmstr.grpc.service.Buildservice.PushFileResponse> getPushFileMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.PushFileMessage, org.qmstr.grpc.service.Buildservice.PushFileResponse> getPushFileMethod;
    if ((getPushFileMethod = BuildServiceGrpc.getPushFileMethod) == null) {
      synchronized (BuildServiceGrpc.class) {
        if ((getPushFileMethod = BuildServiceGrpc.getPushFileMethod) == null) {
          BuildServiceGrpc.getPushFileMethod = getPushFileMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Buildservice.PushFileMessage, org.qmstr.grpc.service.Buildservice.PushFileResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.BuildService", "PushFile"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.PushFileMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.PushFileResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new BuildServiceMethodDescriptorSupplier("PushFile"))
                  .build();
          }
        }
     }
     return getPushFileMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getUpdatePackageNodeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "UpdatePackageNode",
      requestType = org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage.class,
      responseType = org.qmstr.grpc.service.Buildservice.BuildResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getUpdatePackageNodeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage, org.qmstr.grpc.service.Buildservice.BuildResponse> getUpdatePackageNodeMethod;
    if ((getUpdatePackageNodeMethod = BuildServiceGrpc.getUpdatePackageNodeMethod) == null) {
      synchronized (BuildServiceGrpc.class) {
        if ((getUpdatePackageNodeMethod = BuildServiceGrpc.getUpdatePackageNodeMethod) == null) {
          BuildServiceGrpc.getUpdatePackageNodeMethod = getUpdatePackageNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage, org.qmstr.grpc.service.Buildservice.BuildResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.BuildService", "UpdatePackageNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.BuildResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new BuildServiceMethodDescriptorSupplier("UpdatePackageNode"))
                  .build();
          }
        }
     }
     return getUpdatePackageNodeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.PackageNode,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getCreatePackageMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CreatePackage",
      requestType = org.qmstr.grpc.service.Datamodel.PackageNode.class,
      responseType = org.qmstr.grpc.service.Buildservice.BuildResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.PackageNode,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getCreatePackageMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.PackageNode, org.qmstr.grpc.service.Buildservice.BuildResponse> getCreatePackageMethod;
    if ((getCreatePackageMethod = BuildServiceGrpc.getCreatePackageMethod) == null) {
      synchronized (BuildServiceGrpc.class) {
        if ((getCreatePackageMethod = BuildServiceGrpc.getCreatePackageMethod) == null) {
          BuildServiceGrpc.getCreatePackageMethod = getCreatePackageMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Datamodel.PackageNode, org.qmstr.grpc.service.Buildservice.BuildResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.BuildService", "CreatePackage"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.PackageNode.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.BuildResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new BuildServiceMethodDescriptorSupplier("CreatePackage"))
                  .build();
          }
        }
     }
     return getCreatePackageMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.ProjectNode,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getCreateProjectMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CreateProject",
      requestType = org.qmstr.grpc.service.Datamodel.ProjectNode.class,
      responseType = org.qmstr.grpc.service.Buildservice.BuildResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.ProjectNode,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getCreateProjectMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.ProjectNode, org.qmstr.grpc.service.Buildservice.BuildResponse> getCreateProjectMethod;
    if ((getCreateProjectMethod = BuildServiceGrpc.getCreateProjectMethod) == null) {
      synchronized (BuildServiceGrpc.class) {
        if ((getCreateProjectMethod = BuildServiceGrpc.getCreateProjectMethod) == null) {
          BuildServiceGrpc.getCreateProjectMethod = getCreateProjectMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Datamodel.ProjectNode, org.qmstr.grpc.service.Buildservice.BuildResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.BuildService", "CreateProject"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.ProjectNode.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.BuildResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new BuildServiceMethodDescriptorSupplier("CreateProject"))
                  .build();
          }
        }
     }
     return getCreateProjectMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.PackageNode,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getUpdateProjectNodeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "UpdateProjectNode",
      requestType = org.qmstr.grpc.service.Datamodel.PackageNode.class,
      responseType = org.qmstr.grpc.service.Buildservice.BuildResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.PackageNode,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getUpdateProjectNodeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Datamodel.PackageNode, org.qmstr.grpc.service.Buildservice.BuildResponse> getUpdateProjectNodeMethod;
    if ((getUpdateProjectNodeMethod = BuildServiceGrpc.getUpdateProjectNodeMethod) == null) {
      synchronized (BuildServiceGrpc.class) {
        if ((getUpdateProjectNodeMethod = BuildServiceGrpc.getUpdateProjectNodeMethod) == null) {
          BuildServiceGrpc.getUpdateProjectNodeMethod = getUpdateProjectNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Datamodel.PackageNode, org.qmstr.grpc.service.Buildservice.BuildResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "service.BuildService", "UpdateProjectNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.PackageNode.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.BuildResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new BuildServiceMethodDescriptorSupplier("UpdateProjectNode"))
                  .build();
          }
        }
     }
     return getUpdateProjectNodeMethod;
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
    if ((getGetProjectNodeMethod = BuildServiceGrpc.getGetProjectNodeMethod) == null) {
      synchronized (BuildServiceGrpc.class) {
        if ((getGetProjectNodeMethod = BuildServiceGrpc.getGetProjectNodeMethod) == null) {
          BuildServiceGrpc.getGetProjectNodeMethod = getGetProjectNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Datamodel.ProjectNode, org.qmstr.grpc.service.Datamodel.ProjectNode>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.BuildService", "GetProjectNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.ProjectNode.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Datamodel.ProjectNode.getDefaultInstance()))
                  .setSchemaDescriptor(new BuildServiceMethodDescriptorSupplier("GetProjectNode"))
                  .build();
          }
        }
     }
     return getGetProjectNodeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.DeleteMessage,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getDeleteNodeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "DeleteNode",
      requestType = org.qmstr.grpc.service.Buildservice.DeleteMessage.class,
      responseType = org.qmstr.grpc.service.Buildservice.BuildResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.DeleteMessage,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getDeleteNodeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.DeleteMessage, org.qmstr.grpc.service.Buildservice.BuildResponse> getDeleteNodeMethod;
    if ((getDeleteNodeMethod = BuildServiceGrpc.getDeleteNodeMethod) == null) {
      synchronized (BuildServiceGrpc.class) {
        if ((getDeleteNodeMethod = BuildServiceGrpc.getDeleteNodeMethod) == null) {
          BuildServiceGrpc.getDeleteNodeMethod = getDeleteNodeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Buildservice.DeleteMessage, org.qmstr.grpc.service.Buildservice.BuildResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.CLIENT_STREAMING)
              .setFullMethodName(generateFullMethodName(
                  "service.BuildService", "DeleteNode"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.DeleteMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.BuildResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new BuildServiceMethodDescriptorSupplier("DeleteNode"))
                  .build();
          }
        }
     }
     return getDeleteNodeMethod;
  }

  private static volatile io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.DeleteMessage,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getDeleteEdgeMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "DeleteEdge",
      requestType = org.qmstr.grpc.service.Buildservice.DeleteMessage.class,
      responseType = org.qmstr.grpc.service.Buildservice.BuildResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.DeleteMessage,
      org.qmstr.grpc.service.Buildservice.BuildResponse> getDeleteEdgeMethod() {
    io.grpc.MethodDescriptor<org.qmstr.grpc.service.Buildservice.DeleteMessage, org.qmstr.grpc.service.Buildservice.BuildResponse> getDeleteEdgeMethod;
    if ((getDeleteEdgeMethod = BuildServiceGrpc.getDeleteEdgeMethod) == null) {
      synchronized (BuildServiceGrpc.class) {
        if ((getDeleteEdgeMethod = BuildServiceGrpc.getDeleteEdgeMethod) == null) {
          BuildServiceGrpc.getDeleteEdgeMethod = getDeleteEdgeMethod = 
              io.grpc.MethodDescriptor.<org.qmstr.grpc.service.Buildservice.DeleteMessage, org.qmstr.grpc.service.Buildservice.BuildResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(
                  "service.BuildService", "DeleteEdge"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.DeleteMessage.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  org.qmstr.grpc.service.Buildservice.BuildResponse.getDefaultInstance()))
                  .setSchemaDescriptor(new BuildServiceMethodDescriptorSupplier("DeleteEdge"))
                  .build();
          }
        }
     }
     return getDeleteEdgeMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static BuildServiceStub newStub(io.grpc.Channel channel) {
    return new BuildServiceStub(channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static BuildServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    return new BuildServiceBlockingStub(channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static BuildServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    return new BuildServiceFutureStub(channel);
  }

  /**
   */
  public static abstract class BuildServiceImplBase implements io.grpc.BindableService {

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.FileNode> build(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getBuildMethod(), responseObserver);
    }

    /**
     */
    public void sendBuildError(org.qmstr.grpc.service.Datamodel.InfoNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getSendBuildErrorMethod(), responseObserver);
    }

    /**
     */
    public void pushFile(org.qmstr.grpc.service.Buildservice.PushFileMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.PushFileResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getPushFileMethod(), responseObserver);
    }

    /**
     */
    public void updatePackageNode(org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getUpdatePackageNodeMethod(), responseObserver);
    }

    /**
     */
    public void createPackage(org.qmstr.grpc.service.Datamodel.PackageNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getCreatePackageMethod(), responseObserver);
    }

    /**
     */
    public void createProject(org.qmstr.grpc.service.Datamodel.ProjectNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getCreateProjectMethod(), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.PackageNode> updateProjectNode(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getUpdateProjectNodeMethod(), responseObserver);
    }

    /**
     */
    public void getProjectNode(org.qmstr.grpc.service.Datamodel.ProjectNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.ProjectNode> responseObserver) {
      asyncUnimplementedUnaryCall(getGetProjectNodeMethod(), responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.DeleteMessage> deleteNode(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      return asyncUnimplementedStreamingCall(getDeleteNodeMethod(), responseObserver);
    }

    /**
     */
    public void deleteEdge(org.qmstr.grpc.service.Buildservice.DeleteMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      asyncUnimplementedUnaryCall(getDeleteEdgeMethod(), responseObserver);
    }

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
          .addMethod(
            getBuildMethod(),
            asyncClientStreamingCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Datamodel.FileNode,
                org.qmstr.grpc.service.Buildservice.BuildResponse>(
                  this, METHODID_BUILD)))
          .addMethod(
            getSendBuildErrorMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Datamodel.InfoNode,
                org.qmstr.grpc.service.Buildservice.BuildResponse>(
                  this, METHODID_SEND_BUILD_ERROR)))
          .addMethod(
            getPushFileMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Buildservice.PushFileMessage,
                org.qmstr.grpc.service.Buildservice.PushFileResponse>(
                  this, METHODID_PUSH_FILE)))
          .addMethod(
            getUpdatePackageNodeMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage,
                org.qmstr.grpc.service.Buildservice.BuildResponse>(
                  this, METHODID_UPDATE_PACKAGE_NODE)))
          .addMethod(
            getCreatePackageMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Datamodel.PackageNode,
                org.qmstr.grpc.service.Buildservice.BuildResponse>(
                  this, METHODID_CREATE_PACKAGE)))
          .addMethod(
            getCreateProjectMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Datamodel.ProjectNode,
                org.qmstr.grpc.service.Buildservice.BuildResponse>(
                  this, METHODID_CREATE_PROJECT)))
          .addMethod(
            getUpdateProjectNodeMethod(),
            asyncClientStreamingCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Datamodel.PackageNode,
                org.qmstr.grpc.service.Buildservice.BuildResponse>(
                  this, METHODID_UPDATE_PROJECT_NODE)))
          .addMethod(
            getGetProjectNodeMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Datamodel.ProjectNode,
                org.qmstr.grpc.service.Datamodel.ProjectNode>(
                  this, METHODID_GET_PROJECT_NODE)))
          .addMethod(
            getDeleteNodeMethod(),
            asyncClientStreamingCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Buildservice.DeleteMessage,
                org.qmstr.grpc.service.Buildservice.BuildResponse>(
                  this, METHODID_DELETE_NODE)))
          .addMethod(
            getDeleteEdgeMethod(),
            asyncUnaryCall(
              new MethodHandlers<
                org.qmstr.grpc.service.Buildservice.DeleteMessage,
                org.qmstr.grpc.service.Buildservice.BuildResponse>(
                  this, METHODID_DELETE_EDGE)))
          .build();
    }
  }

  /**
   */
  public static final class BuildServiceStub extends io.grpc.stub.AbstractStub<BuildServiceStub> {
    private BuildServiceStub(io.grpc.Channel channel) {
      super(channel);
    }

    private BuildServiceStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BuildServiceStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new BuildServiceStub(channel, callOptions);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.FileNode> build(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      return asyncClientStreamingCall(
          getChannel().newCall(getBuildMethod(), getCallOptions()), responseObserver);
    }

    /**
     */
    public void sendBuildError(org.qmstr.grpc.service.Datamodel.InfoNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getSendBuildErrorMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void pushFile(org.qmstr.grpc.service.Buildservice.PushFileMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.PushFileResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getPushFileMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void updatePackageNode(org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getUpdatePackageNodeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void createPackage(org.qmstr.grpc.service.Datamodel.PackageNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getCreatePackageMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void createProject(org.qmstr.grpc.service.Datamodel.ProjectNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getCreateProjectMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.PackageNode> updateProjectNode(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      return asyncClientStreamingCall(
          getChannel().newCall(getUpdateProjectNodeMethod(), getCallOptions()), responseObserver);
    }

    /**
     */
    public void getProjectNode(org.qmstr.grpc.service.Datamodel.ProjectNode request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.ProjectNode> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getGetProjectNodeMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.DeleteMessage> deleteNode(
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      return asyncClientStreamingCall(
          getChannel().newCall(getDeleteNodeMethod(), getCallOptions()), responseObserver);
    }

    /**
     */
    public void deleteEdge(org.qmstr.grpc.service.Buildservice.DeleteMessage request,
        io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse> responseObserver) {
      asyncUnaryCall(
          getChannel().newCall(getDeleteEdgeMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   */
  public static final class BuildServiceBlockingStub extends io.grpc.stub.AbstractStub<BuildServiceBlockingStub> {
    private BuildServiceBlockingStub(io.grpc.Channel channel) {
      super(channel);
    }

    private BuildServiceBlockingStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BuildServiceBlockingStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new BuildServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public org.qmstr.grpc.service.Buildservice.BuildResponse sendBuildError(org.qmstr.grpc.service.Datamodel.InfoNode request) {
      return blockingUnaryCall(
          getChannel(), getSendBuildErrorMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Buildservice.PushFileResponse pushFile(org.qmstr.grpc.service.Buildservice.PushFileMessage request) {
      return blockingUnaryCall(
          getChannel(), getPushFileMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Buildservice.BuildResponse updatePackageNode(org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage request) {
      return blockingUnaryCall(
          getChannel(), getUpdatePackageNodeMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Buildservice.BuildResponse createPackage(org.qmstr.grpc.service.Datamodel.PackageNode request) {
      return blockingUnaryCall(
          getChannel(), getCreatePackageMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Buildservice.BuildResponse createProject(org.qmstr.grpc.service.Datamodel.ProjectNode request) {
      return blockingUnaryCall(
          getChannel(), getCreateProjectMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Datamodel.ProjectNode getProjectNode(org.qmstr.grpc.service.Datamodel.ProjectNode request) {
      return blockingUnaryCall(
          getChannel(), getGetProjectNodeMethod(), getCallOptions(), request);
    }

    /**
     */
    public org.qmstr.grpc.service.Buildservice.BuildResponse deleteEdge(org.qmstr.grpc.service.Buildservice.DeleteMessage request) {
      return blockingUnaryCall(
          getChannel(), getDeleteEdgeMethod(), getCallOptions(), request);
    }
  }

  /**
   */
  public static final class BuildServiceFutureStub extends io.grpc.stub.AbstractStub<BuildServiceFutureStub> {
    private BuildServiceFutureStub(io.grpc.Channel channel) {
      super(channel);
    }

    private BuildServiceFutureStub(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected BuildServiceFutureStub build(io.grpc.Channel channel,
        io.grpc.CallOptions callOptions) {
      return new BuildServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Buildservice.BuildResponse> sendBuildError(
        org.qmstr.grpc.service.Datamodel.InfoNode request) {
      return futureUnaryCall(
          getChannel().newCall(getSendBuildErrorMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Buildservice.PushFileResponse> pushFile(
        org.qmstr.grpc.service.Buildservice.PushFileMessage request) {
      return futureUnaryCall(
          getChannel().newCall(getPushFileMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Buildservice.BuildResponse> updatePackageNode(
        org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage request) {
      return futureUnaryCall(
          getChannel().newCall(getUpdatePackageNodeMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Buildservice.BuildResponse> createPackage(
        org.qmstr.grpc.service.Datamodel.PackageNode request) {
      return futureUnaryCall(
          getChannel().newCall(getCreatePackageMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Buildservice.BuildResponse> createProject(
        org.qmstr.grpc.service.Datamodel.ProjectNode request) {
      return futureUnaryCall(
          getChannel().newCall(getCreateProjectMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Datamodel.ProjectNode> getProjectNode(
        org.qmstr.grpc.service.Datamodel.ProjectNode request) {
      return futureUnaryCall(
          getChannel().newCall(getGetProjectNodeMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<org.qmstr.grpc.service.Buildservice.BuildResponse> deleteEdge(
        org.qmstr.grpc.service.Buildservice.DeleteMessage request) {
      return futureUnaryCall(
          getChannel().newCall(getDeleteEdgeMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_SEND_BUILD_ERROR = 0;
  private static final int METHODID_PUSH_FILE = 1;
  private static final int METHODID_UPDATE_PACKAGE_NODE = 2;
  private static final int METHODID_CREATE_PACKAGE = 3;
  private static final int METHODID_CREATE_PROJECT = 4;
  private static final int METHODID_GET_PROJECT_NODE = 5;
  private static final int METHODID_DELETE_EDGE = 6;
  private static final int METHODID_BUILD = 7;
  private static final int METHODID_UPDATE_PROJECT_NODE = 8;
  private static final int METHODID_DELETE_NODE = 9;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final BuildServiceImplBase serviceImpl;
    private final int methodId;

    MethodHandlers(BuildServiceImplBase serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_SEND_BUILD_ERROR:
          serviceImpl.sendBuildError((org.qmstr.grpc.service.Datamodel.InfoNode) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse>) responseObserver);
          break;
        case METHODID_PUSH_FILE:
          serviceImpl.pushFile((org.qmstr.grpc.service.Buildservice.PushFileMessage) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.PushFileResponse>) responseObserver);
          break;
        case METHODID_UPDATE_PACKAGE_NODE:
          serviceImpl.updatePackageNode((org.qmstr.grpc.service.Buildservice.UpdatePackageNodeMessage) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse>) responseObserver);
          break;
        case METHODID_CREATE_PACKAGE:
          serviceImpl.createPackage((org.qmstr.grpc.service.Datamodel.PackageNode) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse>) responseObserver);
          break;
        case METHODID_CREATE_PROJECT:
          serviceImpl.createProject((org.qmstr.grpc.service.Datamodel.ProjectNode) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse>) responseObserver);
          break;
        case METHODID_GET_PROJECT_NODE:
          serviceImpl.getProjectNode((org.qmstr.grpc.service.Datamodel.ProjectNode) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Datamodel.ProjectNode>) responseObserver);
          break;
        case METHODID_DELETE_EDGE:
          serviceImpl.deleteEdge((org.qmstr.grpc.service.Buildservice.DeleteMessage) request,
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse>) responseObserver);
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
        case METHODID_BUILD:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.build(
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse>) responseObserver);
        case METHODID_UPDATE_PROJECT_NODE:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.updateProjectNode(
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse>) responseObserver);
        case METHODID_DELETE_NODE:
          return (io.grpc.stub.StreamObserver<Req>) serviceImpl.deleteNode(
              (io.grpc.stub.StreamObserver<org.qmstr.grpc.service.Buildservice.BuildResponse>) responseObserver);
        default:
          throw new AssertionError();
      }
    }
  }

  private static abstract class BuildServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    BuildServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return org.qmstr.grpc.service.Buildservice.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("BuildService");
    }
  }

  private static final class BuildServiceFileDescriptorSupplier
      extends BuildServiceBaseDescriptorSupplier {
    BuildServiceFileDescriptorSupplier() {}
  }

  private static final class BuildServiceMethodDescriptorSupplier
      extends BuildServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final String methodName;

    BuildServiceMethodDescriptorSupplier(String methodName) {
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
      synchronized (BuildServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new BuildServiceFileDescriptorSupplier())
              .addMethod(getBuildMethod())
              .addMethod(getSendBuildErrorMethod())
              .addMethod(getPushFileMethod())
              .addMethod(getUpdatePackageNodeMethod())
              .addMethod(getCreatePackageMethod())
              .addMethod(getCreateProjectMethod())
              .addMethod(getUpdateProjectNodeMethod())
              .addMethod(getGetProjectNodeMethod())
              .addMethod(getDeleteNodeMethod())
              .addMethod(getDeleteEdgeMethod())
              .build();
        }
      }
    }
    return result;
  }
}
