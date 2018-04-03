package org.qmstr.gradle;

import org.gradle.api.file.FileCollection;
import org.gradle.api.tasks.SourceSet;
import org.gradle.api.tasks.SourceSetOutput;
import org.gradle.api.tasks.TaskAction;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.grpc.service.Datamodel;
import org.qmstr.util.FilenodeUtils;

import java.util.Set;


public class QmstrCompileTask extends QmstrTask {
    private Iterable<SourceSet> sourceSets;
    private BuildServiceClient bsc;

    public void setSourceSets(Iterable<SourceSet> sources) {
        this.sourceSets = sources;
    }

    @TaskAction
    void build() {
        bsc = new BuildServiceClient(buildServiceAddress, buildServicePort);
        bsc.SendLogMessage("This is qmstr-gradle-plugin!");

        this.sourceSets.forEach(set -> {
            FileCollection sourceDirs = set.getAllJava().getSourceDirectories();
            SourceSetOutput outDirs = set.getOutput();
            bsc.SendLogMessage("Found sourceset: " + set.getName());
            set.getAllJava().forEach(js -> {
                    sourceDirs.forEach(sd -> bsc.SendLogMessage("Source dir " + sd.toString()));
                    bsc.SendLogMessage("Sending " + js);
                    Set<Datamodel.FileNode> nodes = FilenodeUtils.processSourceFile(js, sourceDirs, outDirs);
                    nodes.forEach(node -> bsc.SendBuildMessage(node));
            });
        });

    }




}