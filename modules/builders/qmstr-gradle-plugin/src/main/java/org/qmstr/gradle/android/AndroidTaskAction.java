package org.qmstr.gradle.android;

import org.gradle.api.Action;
import org.gradle.api.Project;
import org.gradle.api.Task;
import org.qmstr.client.BuildServiceClient;

public abstract class AndroidTaskAction implements Action<Task> {

    protected Project project;
    protected static final String dexTaskPrefix = "transformClassesWithDexBuilder";
    protected BuildServiceClient bsc;
    protected String buildServiceAddress;
    protected int buildServicePort;

    public void setBuildServiceAddress(String address) {
        String[] addressSplit = address.split(":");
        this.buildServiceAddress = addressSplit[0];
        this.buildServicePort = Integer.parseInt(addressSplit[1]);
    }
}

