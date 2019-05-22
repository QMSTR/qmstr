package org.qmstr;


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
import org.qmstr.util.FilenodeUtils;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Collections;
import java.util.List;
import java.util.Set;
import java.util.stream.Collectors;

/**
 * add built files to qmstr master server
 */
@Mojo( name = "qmstrbuild", defaultPhase = LifecyclePhase.PROCESS_CLASSES )
public class QmstrBuildMojo extends AbstractMojo
{
    private final String qmstrMasterAddress =System.getenv("QMSTR_MASTER");

    /**
     * address to connect to
     */
    @Parameter( defaultValue = "localhost", property = "qmstrAddress", required = true )
    private String qmstrAddress;

    @Parameter( defaultValue = "${session}", readonly = true )
    private MavenSession session;
 
    @Parameter( defaultValue = "${project}", readonly = true )
    private MavenProject project;

    /**
     * The directory for compiled classes.
     *
     */
    @Parameter( readonly = true, required = true, defaultValue = "${project.build.outputDirectory}" )
    private File outputDirectory;

    public void execute() throws MojoExecutionException
    {
        BuildServiceClient bsc = new BuildServiceClient(qmstrAddress);

        getLog().info("I got triggered for " + project.getName());

        ArtifactHandler artifactHandler = project.getArtifact().getArtifactHandler();
        if (!artifactHandler.getLanguage().equalsIgnoreCase("java")) {
            getLog().warn("Not a java project");
            return;
        }

        Set<File> sourceDirs = project.getCompileSourceRoots().stream()
                .map(sds -> Paths.get(sds).toFile())
                .collect(Collectors.toSet());

        Set<File> sources = getSourceFiles(project.getCompileSourceRoots());

        getLog().info("Puttn stuff to " + outputDirectory.toString());

        sources.stream().forEach(s ->
        {
            getLog().info(s.toString());
            getLog().info(sourceDirs.stream().map(sd -> sd.toString()).collect(Collectors.joining(",")));
            Set<Datamodel.FileNode> fileNodes = FilenodeUtils.processSourceFile(s, sourceDirs, Collections.singleton(outputDirectory));
            fileNodes.stream().
                    forEach(fileNode -> getLog().info(String.format("Found %s build from %s", fileNode.getName(), fileNode.getDerivedFromList().stream().map(fileNode1 -> fileNode1.getName()).collect(Collectors.joining(",")))));

            bsc.SendBuildFileNodes(fileNodes);
        });


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
