package org.qmstr;


import org.apache.maven.artifact.Artifact;
import org.apache.maven.artifact.handler.ArtifactHandler;
import org.apache.maven.execution.MavenSession;
import org.apache.maven.plugin.AbstractMojo;
import org.apache.maven.plugin.MojoExecutionException;

import org.apache.maven.plugins.annotations.LifecyclePhase;
import org.apache.maven.plugins.annotations.Mojo;
import org.apache.maven.plugins.annotations.Parameter;
import org.apache.maven.project.MavenProject;
import org.qmstr.grpc.service.Datamodel;
import org.qmstr.util.FilenodeUtils;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.util.PackagenodeUtils;

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


        getLog().info("I got triggered for " + project.getName());

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

        Optional<Datamodel.PackageNode> pkgFileNodeOpt = PackagenodeUtils.processArtifact(art.getFile(), project.getVersion());

        Datamodel.PackageNode pkgFileNode = pkgFileNodeOpt.orElseThrow(() -> new MojoExecutionException("qmstr: no package found"));

        bsc.SendPackageNode(pkgFileNode);

    }

    private Set<File> getSourceFiles(List<String> sourceDirs) {
        getLog().info("collecting sources from " + sourceDirs);
        return sourceDirs.stream().map(sds -> Paths.get(sds))
                .flatMap(sd -> getSourcesFromDir(sd).stream())
                .collect(Collectors.toSet());
    }

    private Set<File> getSourcesFromDir(Path sourceDir) {
        getLog().info("descent into " + sourceDir);

        try {
            return Files.walk(sourceDir)
                    .filter(Files::isRegularFile)
                    .filter(p -> p.getFileName().toString().endsWith(".java"))
                    .map(p -> p.toFile())
                    .collect(Collectors.toSet());
        } catch (IOException ioe) {
            return Collections.emptySet();
        }
    }
}
