package org.qmstr.gradle;

import org.gradle.api.artifacts.ConfigurationContainer;
import org.gradle.api.artifacts.PublishArtifact;
import org.gradle.api.tasks.TaskAction;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.util.FilenodeUtils;
import org.qmstr.util.PackagenodeUtils;

import java.io.File;
import java.util.HashMap;
import java.util.Set;
import java.util.stream.Collectors;

public class QmstrPackTask extends QmstrTask {

    private ConfigurationContainer config;
    private BuildServiceClient bsc;

    public void setProjectConfig(ConfigurationContainer configurations) {
        this.config = configurations;
    }

    @TaskAction
    void pack() {
        QmstrPluginExtension extension = (QmstrPluginExtension) getProject()
                .getExtensions().findByName("qmstr");

        this.setBuildServiceAddress(extension.qmstrAddress);

        bsc = new BuildServiceClient(buildServiceAddress, buildServicePort);


        Set<File> artifactFiles = this.config
                .parallelStream()
                .filter(c -> c.isCanBeResolved())
                .flatMap(c -> c.getAllArtifacts().stream())
                .map(art -> art.getFile())
                .collect(Collectors.toSet());


        artifactFiles.stream()
                .map(art -> PackagenodeUtils.processArtifact(art, this.getProject().getPath(), this.getProject().getVersion().toString()))
                .filter(o -> o.isPresent())
                .map(o -> o.get())
                .forEach(art -> bsc.SendPackageNode(art));
    }
}


