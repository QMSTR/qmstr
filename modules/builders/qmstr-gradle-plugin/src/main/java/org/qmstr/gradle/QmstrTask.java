package org.qmstr.gradle;

import org.gradle.api.DefaultTask;

public class QmstrTask extends DefaultTask {
    protected String buildServiceAddress;
    protected int buildServicePort;

    public String getBuildServiceAddress() { return buildServiceAddress; }

    public void setBuildServiceAddress(String address) {
        String[] addressSplit = address.split(":");
        this.buildServiceAddress = addressSplit[0];
        this.buildServicePort = Integer.parseInt(addressSplit[1]);
    }
}
