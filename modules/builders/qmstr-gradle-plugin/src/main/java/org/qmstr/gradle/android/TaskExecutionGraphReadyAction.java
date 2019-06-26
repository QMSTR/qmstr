package org.qmstr.gradle.android;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.execution.TaskExecutionGraph;

public class TaskExecutionGraphReadyAction implements Action<TaskExecutionGraph> {

    private Project project;

    public TaskExecutionGraphReadyAction(Project project) {
        this.project = project;
    }

    @Override
    public void execute(TaskExecutionGraph graph) {
    }

}