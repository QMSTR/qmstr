package org.qmstr.gradle;

import java.io.File;
import java.util.Set;

public class QmstrPluginExtension {
    public String qmstrAddress;
    public Set<File> sourceDirs;
    public Set<File> outDirs;

    public QmstrPluginExtension() {
        this.qmstrAddress = "localhost:50051";

        String difAddress = System.getenv("QMSTR_MASTER");
        if (difAddress != null && !difAddress.trim().isEmpty()){
            this.qmstrAddress = difAddress;
        }
    }
}
