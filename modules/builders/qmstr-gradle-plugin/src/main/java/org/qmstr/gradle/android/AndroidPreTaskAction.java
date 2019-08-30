package org.qmstr.gradle.android;

import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.plugins.AppliedPlugin;

public class AndroidPreTaskAction extends AndroidTaskAction {

    public AndroidPreTaskAction(Project project, AppliedPlugin plugin) {
        super(project, plugin);
    }

    @Override
    public void execute(Task task) {
        if (debug) {
            task.getLogger().warn("Task {} of {} about to run", task.getName(), project.getName());
            logTaskInputOutput(task);
        }
    }
}