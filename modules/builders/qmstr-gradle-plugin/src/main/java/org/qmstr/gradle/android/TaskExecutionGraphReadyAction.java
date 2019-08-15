package org.qmstr.gradle.android;

import java.io.File;
import java.util.Collections;
import java.util.Set;
import java.util.stream.Collector;
import java.util.stream.Collectors;

import com.android.build.gradle.AppExtension;
import com.android.build.gradle.LibraryExtension;
import com.android.build.gradle.api.AndroidSourceSet;

import org.gradle.api.Action;
import org.gradle.api.NamedDomainObjectContainer;
import org.gradle.api.Project;
import org.gradle.api.execution.TaskExecutionGraph;
import org.qmstr.gradle.QmstrPluginExtension;

// This action is executed when the task graph is built. After evaluation of the configuration.
// It is not used right now
public class TaskExecutionGraphReadyAction implements Action<TaskExecutionGraph> {

    private Project project;

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

    public TaskExecutionGraphReadyAction(Project project) {
        this.project = project;
    }

    @Override
    public void execute(TaskExecutionGraph graph) {
    }

}