package org.qmstr.client;

import com.google.protobuf.ByteString;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.Status;
import io.grpc.stub.StreamObserver;

import org.qmstr.grpc.service.*;
import java.util.Set;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

public class BuildServiceClient {

    private final ManagedChannel channel;
    private final BuildServiceGrpc.BuildServiceStub asyncStub;
    private final ControlServiceGrpc.ControlServiceBlockingStub blockingControlStub;

    public BuildServiceClient(String host, int port) {
        this(ManagedChannelBuilder.forAddress(host, port).usePlaintext());
    }

    public BuildServiceClient(ManagedChannelBuilder<?> channelBuilder) {
        channel = channelBuilder.build();
        asyncStub = BuildServiceGrpc.newStub(channel);
        blockingControlStub = ControlServiceGrpc.newBlockingStub(channel);
    }

    public void SendBuildFileNodes(Set<Datamodel.FileNode> fileNodes) {
        if (fileNodes == null || fileNodes.isEmpty()) {
            SendLogMessage("Build with empty set");            
            return;
        }
        final CountDownLatch finishLatch = new CountDownLatch(1);
        StreamObserver<Buildservice.BuildResponse> responseObserver = new StreamObserver<Buildservice.BuildResponse>() {
            @Override
            public void onNext(Buildservice.BuildResponse response) {
                if (!response.getSuccess()){
                    SendLogMessage("Server filenode stream failed");
                }               
            }

            @Override
            public void onError(Throwable t) {
                SendLogMessage("Build Failed: " + Status.fromThrowable(t));
                finishLatch.countDown();
            }

            @Override
            public void onCompleted() {
                finishLatch.countDown();
            }
        };

        StreamObserver<Datamodel.FileNode> requestObserver = asyncStub.build(responseObserver);
        try {
            fileNodes.forEach(fileNode -> {
                requestObserver.onNext(fileNode);
                if (finishLatch.getCount() == 0) {
                    // RPC completed or errored before we finished sending.
                    // Sending further requests won't error, but they will just be thrown away.
                    return;
                }
            });
            requestObserver.onCompleted();
            if (!finishLatch.await(1, TimeUnit.MINUTES)) {
                SendLogMessage("WARNING BuildMessage can not finish within 1 minutes");
            }
        } catch (RuntimeException  e) {
            requestObserver.onError(e);
            throw e;
        } catch (InterruptedException e) {
            requestObserver.onError(e);
        }
    }

    public void SendLogMessage(String message) {
        Controlservice.LogMessage logMsg = Controlservice.LogMessage.newBuilder()
                .setMsg(ByteString.copyFromUtf8(message))
                .build();
        this.blockingControlStub.log(logMsg);
    }

}