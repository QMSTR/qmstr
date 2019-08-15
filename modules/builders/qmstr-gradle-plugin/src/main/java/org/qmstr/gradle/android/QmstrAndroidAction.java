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

    public QmstrAndroidAction(Project project) {
        this.project = project;
    }

    @Override
    public void execute(AppliedPlugin plugin) {
        project.getLogger().warn("Applied plugin {} on project {}", plugin.getId(), project.getName());
        // install actions/listeners to the task graph
        //project.getGradle().getTaskGraph().whenReady(new TaskExecutionGraphReadyAction(project));
        project.getGradle().getTaskGraph().afterTask(new AndroidPostTaskAction(project, plugin));
        project.getGradle().getTaskGraph().beforeTask(new AndroidPreTaskAction(project, plugin));


    }
}