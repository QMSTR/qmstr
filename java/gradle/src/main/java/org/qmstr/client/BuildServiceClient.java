package org.qmstr.client;

import com.google.protobuf.ByteString;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.qmstr.grpc.service.*;

public class BuildServiceClient {

    private final ManagedChannel channel;
    private final BuildServiceGrpc.BuildServiceBlockingStub blockingBuildStub;
    private final BuildServiceGrpc.BuildServiceStub asyncBuildStub;
    private final ControlServiceGrpc.ControlServiceBlockingStub blockingControlStub;

    public BuildServiceClient(String host, int port) {
        this(ManagedChannelBuilder.forAddress(host, port).usePlaintext(true));
    }

    public BuildServiceClient(ManagedChannelBuilder<?> channelBuilder) {
        channel = channelBuilder.build();
        blockingBuildStub = BuildServiceGrpc.newBlockingStub(channel);
        asyncBuildStub = BuildServiceGrpc.newStub(channel);
        blockingControlStub = ControlServiceGrpc.newBlockingStub(channel);
    }

    public void SendBuildMessage(Datamodel.FileNode fileNode) {
        Buildservice.BuildMessage bm = Buildservice.BuildMessage.newBuilder().setFileNodes(0, fileNode).build();

        this.blockingBuildStub.build(bm);
    }

    public void SendLogMessage(String message) {
        Controlservice.LogMessage logMsg = Controlservice.LogMessage.newBuilder()
                .setMsg(ByteString.copyFromUtf8(message))
                .build();
        this.blockingControlStub.log(logMsg);
    }

}