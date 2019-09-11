package org.qmstr.gradle.android;

import java.io.File;
import java.util.Collections;
import java.util.HashSet;
import java.util.Set;
import java.util.stream.Collectors;

import com.android.build.gradle.AppExtension;
import com.android.build.gradle.LibraryExtension;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.plugins.AppliedPlugin;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.gradle.QmstrPluginExtension;

public abstract class AndroidTaskAction implements Action<Task> {

    protected Project project;
    protected BuildServiceClient bsc;
    protected String buildServiceAddress;
    protected int buildServicePort;
    protected QmstrPluginExtension qmstrExt;
    protected final boolean debug;


    public AndroidTaskAction(Project project, AppliedPlugin plugin) {
        this.project = project;
        QmstrPluginExtension extension = (QmstrPluginExtension) this.project.getExtensions().findByName("qmstr");

        this.debug = extension.debug;
        this.setBuildServiceAddress(extension.qmstrAddress);

        this.bsc = new BuildServiceClient(buildServiceAddress, buildServicePort);
    }


    public void setBuildServiceAddress(String address) {
        String[] addressSplit = address.split(":");
        this.buildServiceAddress = addressSplit[0];
        this.buildServicePort = Integer.parseInt(addressSplit[1]);
    }

    public static Set<File> getSourceDirs(Project project) {
        Set<File> sources = new HashSet<>();
        sources.addAll(getAppSourceDirs(project));
        sources.addAll(getLibSourceDirs(project));
        return sources;
    }

    public static Set<File> getAppSourceDirs(Project project) {
        AppExtension e = project.getExtensions().findByType(AppExtension.class);
        if (e == null) {
            return Collections.emptySet();
        }
        project.getLogger().warn("Found app source sets: {}", e.getSourceSets().stream().map(s -> s.getName()).collect(Collectors.joining("\n", "{", "}")));
        return e.getSourceSets().stream().flatMap(s -> s.getJava().getSrcDirs().stream()).collect(Collectors.toSet());
    }

    public static Set<File> getLibSourceDirs(Project project) {
        LibraryExtension le = project.getExtensions().findByType(LibraryExtension.class);
        if (le == null) {
            return Collections.emptySet();
        }
        project.getLogger().warn("Found lib source sets: {}", le.getSourceSets().stream().map(s -> s.getName()).collect(Collectors.joining("\n", "{", "}")));
        return le.getSourceSets().stream().flatMap(s -> s.getJava().getSrcDirs().stream()).collect(Collectors.toSet());
    }

    public static void logTaskInputOutput(Task task) {
        task.getLogger().warn("====================>");
        task.getLogger().warn("project {} handle {} task\nInput:", task.getProject().getName(), task.getName());
        task.getInputs().getFiles().forEach(
            in -> task.getLogger().warn(in.getAbsolutePath())
        );

        task.getLogger().warn("Output:");
        task.getOutputs().getFiles().forEach(
                out -> task.getLogger().warn(out.getAbsolutePath()));
        task.getLogger().warn("<====================");
    }
}

