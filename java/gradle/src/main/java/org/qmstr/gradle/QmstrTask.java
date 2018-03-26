package org.qmstr.gradle;

import org.gradle.api.DefaultTask;
import org.gradle.api.file.FileCollection;
import org.gradle.api.tasks.SourceSet;
import org.gradle.api.tasks.SourceSetOutput;
import org.gradle.api.tasks.TaskAction;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.grpc.service.Datamodel;
import org.qmstr.util.FilenodeUtils;

import java.io.File;
import java.nio.file.Path;


public class QmstrTask extends DefaultTask {
    private String buildServiceAddress;
    private int buildServicePort;
    private Iterable<SourceSet> sourceSets;

    public String getBuildServiceAddress() { return buildServiceAddress; }

    public void setBuildServiceAddress(String address) {
        String[] addressSplit = address.split(":");
        this.buildServiceAddress = addressSplit[0];
        this.buildServicePort = Integer.parseInt(addressSplit[1]);
    }

    public void setSourceSets(Iterable<SourceSet> sources) {
        this.sourceSets = sources;
    }

    @TaskAction
    void build() {
        System.out.printf("Connecting to qmstr buildservice at %s\n", this.buildServiceAddress);
        BuildServiceClient bsc = new BuildServiceClient(buildServiceAddress, buildServicePort);
        bsc.SendLogMessage("This is gradle!");
        this.sourceSets.forEach(set -> {
            FileCollection sourceDirs = set.getAllJava().getSourceDirectories();
            SourceSetOutput outDirs = set.getOutput();
            bsc.SendLogMessage("Found sourceset: " + set.getName());
            set.getAllJava().forEach(js -> {
                    sourceDirs.forEach(sd -> bsc.SendLogMessage("Source dir " + sd.toString()));
                    bsc.SendLogMessage("Sending " + js + " and " + sourceDirs.toString());
                    bsc.SendBuildMessage(FilenodeUtils.processSourceFile(js, sourceDirs, outDirs));
            });
        });

    }
}