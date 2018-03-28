package org.qmstr.gradle;

import org.gradle.api.DefaultTask;
import org.gradle.api.artifacts.ConfigurationContainer;
import org.gradle.api.artifacts.PublishArtifact;
import org.gradle.api.component.Artifact;
import org.gradle.api.file.FileCollection;
import org.gradle.api.tasks.SourceSet;
import org.gradle.api.tasks.SourceSetOutput;
import org.gradle.api.tasks.TaskAction;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.grpc.service.Datamodel;
import org.qmstr.util.FilenodeUtils;
import org.qmstr.util.Hash;

import java.io.File;
import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Path;
import java.util.HashSet;
import java.util.Set;
import java.util.jar.JarEntry;
import java.util.jar.JarFile;


public class QmstrTask extends DefaultTask {
    private String buildServiceAddress;
    private int buildServicePort;
    private Iterable<SourceSet> sourceSets;
    private ConfigurationContainer config;
    private BuildServiceClient bsc;

    public String getBuildServiceAddress() { return buildServiceAddress; }

    public void setBuildServiceAddress(String address) {
        String[] addressSplit = address.split(":");
        this.buildServiceAddress = addressSplit[0];
        this.buildServicePort = Integer.parseInt(addressSplit[1]);
    }

    public void setSourceSets(Iterable<SourceSet> sources) {
        this.sourceSets = sources;
    }

    public void setProjectConfig(ConfigurationContainer configurations) {
        this.config = configurations;
    }

    @TaskAction
    void build() {
        System.out.printf("Connecting to qmstr buildservice at %s\n", this.buildServiceAddress);
        bsc = new BuildServiceClient(buildServiceAddress, buildServicePort);
        bsc.SendLogMessage("This is gradle!");

        Set<PublishArtifact> arts = new HashSet<>();
        this.config.forEach(c -> c.getAllArtifacts().forEach(art -> arts.add(art)));

        this.sourceSets.forEach(set -> {
            FileCollection sourceDirs = set.getAllJava().getSourceDirectories();
            SourceSetOutput outDirs = set.getOutput();
            bsc.SendLogMessage("Found sourceset: " + set.getName());
            set.getAllJava().forEach(js -> {
                    sourceDirs.forEach(sd -> bsc.SendLogMessage("Source dir " + sd.toString()));
                    bsc.SendLogMessage("Sending " + js);
                    bsc.SendBuildMessage(FilenodeUtils.processSourceFile(js, sourceDirs, outDirs));
            });
        });

        arts.forEach(art -> processArtifact(art));

    }

    private void processArtifact(PublishArtifact artifact) {
        if (!artifact.getExtension().equals("jar")) {
            bsc.SendLogMessage(String.format("Artifact extension %s not supported.", artifact.getExtension()));
            return;
        }

        try {
            Set<Datamodel.FileNode> classes = new HashSet<>();
            JarFile jar = new JarFile(artifact.getFile());
            jar.stream().filter(je -> je.getName().endsWith(".class"))
                    .forEach(je -> {
                        String hash = getHash(jar, je);
                        classes.add(FilenodeUtils.getFileNode(je.getName(), hash));
                        bsc.SendLogMessage(String.format("Found class %s with sha256 %s", je.getName(), hash));
                    });
            Datamodel.FileNode rootNode = FilenodeUtils.getFileNode(artifact.getFile().toPath());
            Datamodel.FileNode.Builder rootNodeBuilder = rootNode.toBuilder();
            classes.forEach(c -> rootNodeBuilder.addDerivedFrom(c));
            rootNode = rootNodeBuilder.build();
            bsc.SendBuildMessage(rootNode);

        } catch (IOException  ioe) {
            bsc.SendLogMessage("Could not read jar file " + artifact.getFile().toString());
        }

    }

    private String getHash(JarFile jarfile, JarEntry jarEntry) {
        try {
            InputStream is = jarfile.getInputStream(jarEntry);
            return Hash.getChecksum(is);
        } catch (IOException e) {
            e.printStackTrace();
        }
        return null;
    }

}