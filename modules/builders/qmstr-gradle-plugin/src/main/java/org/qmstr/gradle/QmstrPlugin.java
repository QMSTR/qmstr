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

import org.qmstr.gradle.java.QmstrJavaAction;

public class QmstrPlugin implements Plugin<Project> {

    public void apply(Project project) {
        project.getExtensions().create("qmstr", QmstrPluginExtension.class);

        // Apply plugin for all java subprojects
        project.getAllprojects().stream()
            .forEach(pro -> pro.getPluginManager().withPlugin("java", new QmstrJavaAction(pro)));
    }
}
