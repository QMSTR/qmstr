package org.qmstr.gradle.android;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.plugins.AppliedPlugin;

public class QmstrAndroidAction implements Action<AppliedPlugin> {

    Project project;

    public QmstrAndroidAction(Project project) {
        this.project = project;
    }

    @Override
    public void execute(AppliedPlugin plugin) {
        project.getLogger().warn("Applied plugin {} on project {}", plugin.getId(), project.getName());
        project.getGradle().getTaskGraph().afterTask(new AndroidPostTaskAction(project, plugin));
        project.getGradle().getTaskGraph().beforeTask(new AndroidPreTaskAction(project, plugin));
    }
}