package org.qmstr.gradle.android;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.plugins.AppliedPlugin;
import com.android.build.gradle.AppPlugin;

public class QmstrAndroidAction implements Action<AppliedPlugin> {

    Project project;

    public QmstrAndroidAction(Project project) {
        this.project = project;
    }

    @Override
    public void execute(AppliedPlugin plugin) {
        // project.getLogger().warn(plugin.getName());
        // project.getLogger().warn(plugin.getClass().getName());

        // install actions/listeners to the task graph
        project.getGradle().getTaskGraph().beforeTask(new AndroidPreTaskAction(project));
        project.getGradle().getTaskGraph().afterTask(new AndroidPostTaskAction(project));

    }

}