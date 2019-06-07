package org.qmstr.gradle.java;

import java.util.Set;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.plugins.AppliedPlugin;
import org.qmstr.gradle.QmstrCompileTask;
import org.qmstr.gradle.QmstrPackTask;
import org.qmstr.gradle.QmstrPlugin;
import org.qmstr.gradle.Utils;


public class QmstrJavaAction implements Action<AppliedPlugin> {

    Project project;

    public QmstrJavaAction(Project project) {
        this.project = project;
    }

    @Override
    public void execute(AppliedPlugin javaPlugin) {
        // recurse loop?
        project.getPluginManager().apply(QmstrPlugin.class);

        Set<Task> classes = project.getTasksByName("classes", false);
        Task qmstrCompile = project.getTasks().create("qmstrbuild", QmstrCompileTask.class, (task) -> {
            task.setSourceSets(Utils.getJavaSourceSets(project));
            task.dependsOn(classes);
        });

        Set<Task> jarTask = project.getTasksByName("jar", false);
        project.getTasks().create("qmstr", QmstrPackTask.class, (task) -> {
            task.dependsOn(qmstrCompile, jarTask);
            task.setProjectConfig(project.getConfigurations());
        });
    }

}