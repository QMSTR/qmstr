package org.qmstr.gradle.android;

import java.io.File;
import java.io.FileNotFoundException;
import java.nio.file.Path;
import java.util.Collections;
import java.util.Set;
import java.util.stream.Collectors;

import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.file.FileCollection;
import org.gradle.api.tasks.SourceSet;
import org.gradle.api.tasks.SourceSetOutput;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.gradle.QmstrPluginExtension;
import org.qmstr.gradle.Utils;
import org.qmstr.grpc.service.Datamodel;
import org.qmstr.util.FilenodeUtils;
import org.qmstr.util.transformations.*;

public class AndroidPreTaskAction extends AndroidTaskAction {

    public AndroidPreTaskAction(Project project) {
        this.project = project;
        QmstrPluginExtension extension = (QmstrPluginExtension) this.project.getExtensions().findByName("qmstr");

        this.setBuildServiceAddress(extension.qmstrAddress);

        this.bsc = new BuildServiceClient(buildServiceAddress, buildServicePort);
    }

    @Override
    public void execute(Task task) {
        task.getLogger().warn("Task {} about to run", task.getName());
    }
}