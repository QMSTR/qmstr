package org.qmstr.gradle;

import org.gradle.api.artifacts.ConfigurationContainer;
import org.gradle.api.artifacts.PublishArtifact;
import org.gradle.api.tasks.TaskAction;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.util.FilenodeUtils;

import java.util.HashSet;
import java.util.Set;

public class QmstrPackTask extends QmstrTask {

    private ConfigurationContainer config;
    private BuildServiceClient bsc;

    public void setProjectConfig(ConfigurationContainer configurations) {
        this.config = configurations;
    }

    @TaskAction
    void pack() {
        bsc = new BuildServiceClient(buildServiceAddress, buildServicePort);

        bsc.SendLogMessage("This is qmstr-gradle-plugin!");
        Set<PublishArtifact> arts = new HashSet<>();
        this.config.forEach(c -> c.getAllArtifacts().forEach(art -> arts.add(art)));
        arts.stream()
                .map(art -> FilenodeUtils.processArtifact(art))
                .filter(o -> o.isPresent())
                .map(o -> o.get())
                .forEach(node -> bsc.SendBuildMessage(node));
    }
}


