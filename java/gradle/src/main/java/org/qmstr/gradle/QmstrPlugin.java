package org.qmstr.gradle;

import org.gradle.api.Plugin;
import org.gradle.api.Project;
import org.gradle.api.plugins.JavaPluginConvention;
import org.gradle.api.tasks.SourceSet;

import java.util.Collections;

public class QmstrPlugin implements Plugin<Project> {
    public void apply(Project project) {
        project.getTasks().create("qmstr", QmstrTask.class, (task) -> {
            task.setBuildServiceAddress("localhost:50051");
            task.setSourceSets(getJavaSourceSets(project));
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
