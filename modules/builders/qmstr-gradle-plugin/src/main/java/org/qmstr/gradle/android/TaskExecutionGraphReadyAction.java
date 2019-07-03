package org.qmstr.gradle.android;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.execution.TaskExecutionGraph;

// This action is executed when the task graph is built. After evaluation of the configuration.
// It is not used right now
public class TaskExecutionGraphReadyAction implements Action<TaskExecutionGraph> {

    private Project project;

    public TaskExecutionGraphReadyAction(Project project) {
        this.project = project;
    }

    @Override
    public void execute(TaskExecutionGraph graph) {
    }

}