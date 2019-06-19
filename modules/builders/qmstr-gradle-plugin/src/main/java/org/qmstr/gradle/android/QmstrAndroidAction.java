package org.qmstr.gradle.android;

import java.util.Map;
import java.util.Set;
import java.util.stream.Collector;
import java.util.stream.Collectors;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.plugins.AppliedPlugin;
import org.gradle.api.tasks.TaskContainer;
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
    public void execute(AppliedPlugin androidPlugin) {
        // recurse loop?
        //project.getPluginManager().apply(QmstrPlugin.class);


        //Map<Project, Set<Task>> tasks = project.getAllTasks(true);
        TaskContainer tasks = project.getTasks();
//        Set<Task> compileTasks = tasks.stream()
//            .flatMap(ts -> ts.stream())
//            .filter(t -> t.getName().startsWith("compile"))
//          .collect(Collectors.toSet());
            
        project.getLogger().lifecycle("Found following compile tasks: {}", 
            String.join(", ",
                tasks.stream()
                    .map(t -> t.getName())
                    .collect(Collectors.toSet()).toArray(new String[]{})
            )
        );

        //Task qmstrCompile = project.getTasks().create("qmstrbuild", QmstrCompileTask.class, (task) -> {
            //task.setSourceSets(Utils.getJavaSourceSets(project));
            //task.dependsOn(compileTasks.toArray());
        //});
    }

}