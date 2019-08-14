package org.qmstr.gradle;

import org.gradle.api.Plugin;
import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.distribution.plugins.DistributionPlugin;
import org.gradle.api.plugins.JavaPlugin;
import org.gradle.api.plugins.JavaPluginConvention;
import org.gradle.api.tasks.SourceSet;

import org.qmstr.gradle.android.QmstrAndroidAction;

import java.util.Collections;
import java.util.Set;

import org.qmstr.gradle.java.QmstrJavaAction;

public class QmstrPlugin implements Plugin<Project> {

    public void apply(Project project) {
        project.getExtensions().create("qmstr", QmstrPluginExtension.class);
        project.getLogger().lifecycle("Applying QMSTR plugin to {}!", project.getName());

        // Apply plugin for all java subprojects
        project.getAllprojects().stream()
            .forEach(pro -> {
                pro.getPluginManager().withPlugin("java", new QmstrJavaAction(pro));
                pro.getPluginManager().withPlugin("com.android.library", new QmstrAndroidAction(pro, true));
                pro.getPluginManager().withPlugin("com.android.application", new QmstrAndroidAction(pro, false));
            });
    }
}
