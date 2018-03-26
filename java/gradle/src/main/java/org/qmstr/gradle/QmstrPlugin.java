package org.qmstr.gradle;

import org.gradle.api.Plugin;
import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.distribution.plugins.DistributionPlugin;
import org.gradle.api.plugins.JavaPlugin;
import org.gradle.api.plugins.JavaPluginConvention;
import org.gradle.api.tasks.SourceSet;

import java.util.Collections;
import java.util.List;
import java.util.Set;

public class QmstrPlugin implements Plugin<Project> {
    public void apply(Project project) {
        project.getPluginManager().apply(JavaPlugin.class);
        project.getPluginManager().apply(DistributionPlugin.class);

        Set<Task> jartask = project.getTasksByName("jar", false);

        project.getTasks().create("qmstr", QmstrTask.class, (task) -> {
            task.setBuildServiceAddress("localhost:50051");
            task.setSourceSets(getJavaSourceSets(project));
            task.setDependsOn(jartask);
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
