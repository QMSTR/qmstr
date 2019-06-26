package org.qmstr.gradle.android;

import org.gradle.api.Project;
import org.gradle.api.Task;

public class AndroidPostTaskAction extends AndroidTaskAction {

    public AndroidPostTaskAction(Project project) {
        this.project = project;
    }

    @Override
    public void execute(Task task) {
        if (task.getName().startsWith(dexTaskPrefix)) {
        task.getOutputs().getFiles().forEach(out -> task.getLogger().warn("Task {} output is {}", task.getName(), out.getAbsolutePath()));
        }
    }

}