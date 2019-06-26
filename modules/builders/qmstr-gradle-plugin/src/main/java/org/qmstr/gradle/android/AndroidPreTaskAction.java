package org.qmstr.gradle.android;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.Task;

public class AndroidPreTaskAction extends AndroidTaskAction {

    public AndroidPreTaskAction(Project project) {
        this.project = project;
    }

    @Override
    public void execute(Task task) {
        if (task.getName().startsWith(dexTaskPrefix)) {
            // collect classes and sources
            task.getInputs().getSourceFiles().forEach(sf -> task.getLogger().warn("Sources for {} are {}", task.getName(), sf));
        }
    }
}