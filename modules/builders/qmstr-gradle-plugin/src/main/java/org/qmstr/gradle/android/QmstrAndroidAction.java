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
    public void execute(AppliedPlugin androidPlugin) {
        // recurse loop?
        //project.getPluginManager().apply(QmstrPlugin.class);

        //project.getGradle().getTaskGraph().whenReady(new TaskExecutionGraphReadyAction(project));
        project.getGradle().getTaskGraph().beforeTask(new AndroidPreTaskAction(project));
        project.getGradle().getTaskGraph().afterTask(new AndroidPostTaskAction(project));

    }

}