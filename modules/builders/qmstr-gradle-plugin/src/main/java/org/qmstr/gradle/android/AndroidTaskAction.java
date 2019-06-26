package org.qmstr.gradle.android;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.Task;

public abstract class AndroidTaskAction implements Action<Task> {

    protected Project project;
    protected static final String dexTaskPrefix = "transformClassesWithDexBuilder";
}
