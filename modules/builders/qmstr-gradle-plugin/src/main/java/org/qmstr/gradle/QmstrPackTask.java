package org.qmstr.gradle;

import org.gradle.api.artifacts.ConfigurationContainer;
import org.gradle.api.artifacts.PublishArtifact;
import org.gradle.api.tasks.TaskAction;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.util.FilenodeUtils;

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

        HashMap<PublishArtifact, Set<File>> arts = new HashMap<>();
        this.config
                .parallelStream()
                .filter(c -> c.isCanBeResolved())
                .forEach(c -> c.getAllArtifacts().forEach(art -> arts.put(art, c.getResolvedConfiguration().getFiles())));

        bsc.SendBuildFileNodes(arts.entrySet().parallelStream()
                .map(artEntry -> FilenodeUtils.processArtifact(artEntry.getKey().getFile(), artEntry.getValue()))
                .filter(o -> o.isPresent())
                .map(o -> o.get())
                .collect(Collectors.toSet())); 
        
    
    }

}


