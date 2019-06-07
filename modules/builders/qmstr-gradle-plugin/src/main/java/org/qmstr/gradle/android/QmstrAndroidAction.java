package org.qmstr.gradle.android;

import java.util.Map;
import java.util.Set;
import java.util.stream.Collector;
import java.util.stream.Collectors;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.plugins.AppliedPlugin;
import org.qmstr.gradle.QmstrCompileTask;
import org.qmstr.gradle.QmstrPackTask;
import org.qmstr.gradle.QmstrPlugin;
import org.qmstr.gradle.Utils;


public class QmstrAndroidAction implements Action<AppliedPlugin> {

    Project project;

    public QmstrAndroidAction(Project project) {
        this.project = project;
    }

    @Override
    public void execute(AppliedPlugin javaPlugin) {
        // recurse loop?
        project.getPluginManager().apply(QmstrPlugin.class);

        Map<Project, Set<Task>> tasks = project.getAllTasks(true);
        Set<Task> compileTasks = tasks.values().stream()
            .flatMap(ts -> ts.stream())
            .filter(t -> t.getName().startsWith("compile"))
            .collect(Collectors.toSet());
            
        project.getLogger().debug("Found following compile tasks: %s", 
            String.join(", ",
                compileTasks.stream()
                    .map(t -> t.getName())
                    .collect(Collectors.toSet()).toArray(new String[]{})
            )
        );

        Task qmstrCompile = project.getTasks().create("qmstrbuild", QmstrCompileTask.class, (task) -> {
            task.setSourceSets(Utils.getJavaSourceSets(project));
            task.dependsOn(compileTasks.toArray());
        });
    }

}