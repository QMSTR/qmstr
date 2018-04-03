package org.qmstr.gradle;

import org.gradle.api.Plugin;
import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.distribution.plugins.DistributionPlugin;
import org.gradle.api.plugins.JavaPlugin;
import org.gradle.api.plugins.JavaPluginConvention;
import org.gradle.api.tasks.SourceSet;

import java.util.Collections;
import java.util.Set;

public class QmstrPlugin implements Plugin<Project> {

    public void apply(Project project) {
        project.getPluginManager().apply(JavaPlugin.class);
        project.getPluginManager().apply(DistributionPlugin.class);

        // Apply plugin for all java subprojects
        project.getAllprojects().stream()
                .filter(pro -> pro.getPlugins().hasPlugin(JavaPlugin.class))
                .forEach(pro -> pro.getPluginManager().apply(QmstrPlugin.class));


        Set<Task> jarTask = project.getTasksByName("jar", false);
        Set<Task> classes = project.getTasksByName("classes", false);

        Task qmstrCompile = project.getTasks().create("qmstrbuild", QmstrCompileTask.class, (task) -> {
            task.setBuildServiceAddress("localhost:50051");
            task.setSourceSets(getJavaSourceSets(project));
            task.dependsOn(classes);
        });

        project.getTasks().create("qmstr", QmstrPackTask.class, (task) -> {
            task.dependsOn(qmstrCompile, jarTask);
            task.setBuildServiceAddress("localhost:50051");
            task.setProjectConfig(project.getConfigurations());
        });
    }

    private static Iterable<SourceSet> getJavaSourceSets(Project project) {
        JavaPluginConvention plugin = project.getConvention()
                .getPlugin(JavaPluginConvention.class);
        if (plugin == null) {
            return Collections.emptyList();
        }
        return plugin.getSourceSets();
    }
}
