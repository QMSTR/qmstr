package org.qmstr.gradle;

import org.gradle.api.Plugin;
import org.gradle.api.Project;
import org.qmstr.gradle.android.QmstrAndroidAction;
import org.qmstr.gradle.java.QmstrJavaAction;

public class QmstrPlugin implements Plugin<Project> {

    public void apply(Project project) {
        project.getExtensions().create("qmstr", QmstrPluginExtension.class);
        project.getLogger().lifecycle("Applying QMSTR plugin to {}!", project.getName());

        // Apply plugin for all java subprojects
        project.getAllprojects().stream()
            .forEach(pro -> {
                pro.getPluginManager().apply(QmstrPlugin.class);
            });

        project.getPluginManager().withPlugin("java", new QmstrJavaAction(project));
        project.getPluginManager().withPlugin("com.android.library", new QmstrAndroidAction(project));
        project.getPluginManager().withPlugin("com.android.application", new QmstrAndroidAction(project));
    }
}
