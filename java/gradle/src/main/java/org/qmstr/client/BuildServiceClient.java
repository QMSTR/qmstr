package org.qmstr.client;

import com.google.protobuf.ByteString;
import io.grpc.ManagedChannel;
import io.grpc.ManagedChannelBuilder;
import org.qmstr.grpc.service.*;
import org.qmstr.util.Hash;

import java.nio.file.Path;

public class BuildServiceClient {

    private final ManagedChannel channel;
    private final BuildServiceGrpc.BuildServiceBlockingStub blockingBuildStub;
    private final ControlServiceGrpc.ControlServiceBlockingStub blockingControlStub;

    public BuildServiceClient(String host, int port) {
        this(ManagedChannelBuilder.forAddress(host, port).usePlaintext(true));
    }

    public BuildServiceClient(ManagedChannelBuilder<?> channelBuilder) {
        channel = channelBuilder.build();
        blockingBuildStub = BuildServiceGrpc.newBlockingStub(channel);
        blockingControlStub = ControlServiceGrpc.newBlockingStub(channel);
    }

    public void SendBuildMessage(Datamodel.FileNode fileNode) {
        if (fileNode == null) {
            return;
        }
        Buildservice.BuildMessage bm = Buildservice.BuildMessage.newBuilder()
                .addFileNodes(fileNode)
                .build();

        this.blockingBuildStub.build(bm);
    }

    public void SendLogMessage(String message) {
        Controlservice.LogMessage logMsg = Controlservice.LogMessage.newBuilder()
                .setMsg(ByteString.copyFromUtf8(message))
                .build();
        this.blockingControlStub.log(logMsg);
    }

}