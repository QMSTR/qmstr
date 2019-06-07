package org.qmstr.gradle;

import java.util.Collections;

import org.gradle.api.Project;
import org.gradle.api.plugins.JavaPluginConvention;
import org.gradle.api.tasks.SourceSet;

public class Utils {

    public static Iterable<SourceSet> getJavaSourceSets(Project project) {
        JavaPluginConvention plugin = project.getConvention()
                .getPlugin(JavaPluginConvention.class);
        if (plugin == null) {
            return Collections.emptyList();
        }
        return plugin.getSourceSets();
    }
}