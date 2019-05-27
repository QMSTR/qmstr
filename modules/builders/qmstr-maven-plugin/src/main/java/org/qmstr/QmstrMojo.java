package org.qmstr;


import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Collections;
import java.util.List;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;

import org.apache.maven.artifact.Artifact;
import org.apache.maven.artifact.handler.ArtifactHandler;
import org.apache.maven.execution.MavenSession;
import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugin.MojoExecutionException;
import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.project.MavenProject;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.grpc.service.Datamodel;
import org.qmstr.util.PackagenodeUtils;

/**
 * add built files to qmstr master server
 */
@Mojo(name = "qmstr", defaultPhase = LifecyclePhase.PACKAGE)
public class QmstrMojo extends AbstractMojo {
    private final String qmstrMasterAddress = System.getenv("QMSTR_MASTER");

    /**
     * address to connect to
     */
    @Parameter(defaultValue = "localhost", property = "qmstrAddress", required = true)
    private String qmstrAddress;

    @Parameter(defaultValue = "${session}", readonly = true)
    private MavenSession session;

    @Parameter(defaultValue = "${project}", readonly = true)
    private MavenProject project;

    /**
     * The directory for compiled classes.
     */
    @Parameter(readonly = true, required = true, defaultValue = "${project.build.outputDirectory}")
    private File outputDirectory;

    public void execute() throws MojoExecutionException {
        BuildServiceClient bsc = new BuildServiceClient(qmstrAddress);

        ArtifactHandler artifactHandler = project.getArtifact().getArtifactHandler();
        if (!artifactHandler.getLanguage().equalsIgnoreCase("java")) {
            getLog().warn("Not a java project");
            return;
        }

        Artifact art = project.getArtifact();
        if (art == null) {
            getLog().error("No artifact");
            throw new MojoExecutionException("qmstr: no artifact found");
        }

        File artFile = art.getFile();
        if (artFile == null) {
            getLog().warn(String.format("qmstr: artifact %s:%s with no file", art.getGroupId(), art.getArtifactId()));
            return;
        }
        getLog().info("Processing artifact " + art.getFile());
        Optional<Datamodel.PackageNode> pkgFileNodeOpt = PackagenodeUtils.processArtifact(art.getFile(), project.getVersion());

        Datamodel.PackageNode pkgFileNode = pkgFileNodeOpt.orElseThrow(() -> new MojoExecutionException("qmstr: no package found"));

        try {
            if (!bsc.SendPackageNode(pkgFileNode)) {
                throw new MojoExecutionException("qmstr: sending package node failed");
            }
        } finally {
            try {
                bsc.close();
            } catch (InterruptedException e) {
                throw new MojoExecutionException("qmstr: failed to close grpc channel " + e.getMessage());
            }
        }
    }
}
