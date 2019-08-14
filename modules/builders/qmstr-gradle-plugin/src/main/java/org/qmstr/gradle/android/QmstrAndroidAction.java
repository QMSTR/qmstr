package org.qmstr.gradle.android;

import org.gradle.api.Action;
import org.gradle.api.NamedDomainObjectCollection;
import org.gradle.api.NamedDomainObjectContainer;
import org.gradle.api.Project;
import org.gradle.api.plugins.AppliedPlugin;
import org.gradle.api.tasks.SourceSet;
import org.gradle.api.tasks.SourceSetContainer;

import java.io.File;
import java.util.Set;
import java.util.stream.Collectors;

import com.android.build.api.dsl.extension.AndroidExtension;
import com.android.build.gradle.AppExtension;
import com.android.build.gradle.LibraryExtension;
import com.android.build.gradle.api.AndroidSourceSet;
import com.android.build.gradle.AppPlugin;

public class QmstrAndroidAction implements Action<AppliedPlugin> {

    Project project;
    boolean lib;

    public QmstrAndroidAction(Project project, boolean lib) {
        this.project = project;
        this.lib = lib;
    }

    private NamedDomainObjectContainer<AndroidSourceSet> getAppSourceSets() {
        AppExtension e = project.getExtensions().findByType(AppExtension.class);
        project.getLogger().warn("Found app source sets: {}", e.getSourceSets().stream().map(s -> s.getName()).collect(Collectors.joining("\n", "{", "}")));
        return e.getSourceSets();
    }

    private NamedDomainObjectContainer<AndroidSourceSet> getLibSourceSets() {
        LibraryExtension le = project.getExtensions().findByType(LibraryExtension.class);
        project.getLogger().warn("Found lib source sets: {}", le.getSourceSets().stream().map(s -> s.getName()).collect(Collectors.joining("\n", "{", "}")));
        return le.getSourceSets();
    }

    private Set<File> getSourceDirs() {
        NamedDomainObjectContainer<AndroidSourceSet> sourceSets = lib ? getLibSourceSets() : getAppSourceSets();
        return sourceSets.stream().flatMap(s -> s.getJava().getSrcDirs().stream()).collect(Collectors.toSet());
    }

    @Override
    public void execute(AppliedPlugin plugin) {


        // install actions/listeners to the task graph
        project.getGradle().getTaskGraph().afterTask(new AndroidPostTaskAction(project, getSourceDirs()));

    }
}