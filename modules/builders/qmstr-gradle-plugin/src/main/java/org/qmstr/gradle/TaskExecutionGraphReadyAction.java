package org.qmstr.gradle;

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
        graph.getAllTasks().stream()
            //.filter(t -> t.getName().contains("build"))
            .forEach(t -> project.getLogger().warn("found Task {}", t.getName()));
    }
}