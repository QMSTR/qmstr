package org.qmstr.gradle.android;

import static org.qmstr.gradle.android.AndroidPostTaskAction.handleElse;

import java.io.File;

import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.plugins.AppliedPlugin;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.gradle.QmstrPluginExtension;

public class AndroidPreTaskAction extends AndroidTaskAction {

    public AndroidPreTaskAction(Project project, AppliedPlugin plugin) {
        this.project = project;
        QmstrPluginExtension extension = (QmstrPluginExtension) this.project.getExtensions().findByName("qmstr");

        this.setBuildServiceAddress(extension.qmstrAddress);

        this.bsc = new BuildServiceClient(buildServiceAddress, buildServicePort);
    }

    @Override
    public void execute(Task task) {
        String apkPath = "/home/endomarkus/devel/QMSTR/qmstr-demo/demos/android-blockly/blockly-android/blocklydemo/build/outputs/apk/debug/blocklydemo-debug.apk";

        File apk = new File(apkPath);

        task.getLogger().warn("{} {}", apk.getAbsolutePath(), apk.exists() ? "exists" : "does not exist");
        task.getLogger().warn("Task {} of {} about to run", task.getName(), project.getName());
        handleElse(task);
    }
}