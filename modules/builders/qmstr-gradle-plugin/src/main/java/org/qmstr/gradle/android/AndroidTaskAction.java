package org.qmstr.gradle.android;

import java.io.File;
import java.util.Set;

import com.android.build.gradle.api.AndroidSourceSet;

import org.gradle.api.Action;
import org.gradle.api.NamedDomainObjectContainer;
import org.gradle.api.Project;
import org.gradle.api.Task;
import org.qmstr.client.BuildServiceClient;

public abstract class AndroidTaskAction implements Action<Task> {

    protected Project project;
    protected BuildServiceClient bsc;
    protected String buildServiceAddress;
    protected int buildServicePort;
    protected Set<File> sourceDirs;

    public void setBuildServiceAddress(String address) {
        String[] addressSplit = address.split(":");
        this.buildServiceAddress = addressSplit[0];
        this.buildServicePort = Integer.parseInt(addressSplit[1]);
    }
}

