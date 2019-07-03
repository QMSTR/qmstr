package org.qmstr.gradle;

public class QmstrPluginExtension {
    public String qmstrAddress;

    public QmstrPluginExtension() {
        this.qmstrAddress = "localhost:50051";

        String difAddress = System.getenv("QMSTR_MASTER");
        if (difAddress != null && !difAddress.trim().isEmpty()){
            this.qmstrAddress = difAddress;
        }
    }
}
