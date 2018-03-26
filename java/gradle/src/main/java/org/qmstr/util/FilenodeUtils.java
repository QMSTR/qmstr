package org.qmstr.util;

import groovy.ui.SystemOutputInterceptor;
import org.gradle.api.Project;
import org.gradle.api.file.FileCollection;
import org.gradle.api.tasks.SourceSetOutput;
import org.qmstr.grpc.service.Datamodel;

import javax.xml.crypto.Data;
import java.io.File;
import java.io.FileNotFoundException;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;

public class FilenodeUtils {

    public static Datamodel.FileNode getFileNode(Path filepath) {
        String checksum = Hash.getChecksum(filepath.toFile());


        return Datamodel.FileNode.newBuilder()
                .setName(filepath.getFileName().toString())
                .setPath(filepath.toString())
                .setHash(checksum != null ? checksum : "nohash"+filepath.toString())
                .setBroken(checksum == null)
                .build();
    }

    public static Datamodel.FileNode processSourceFile(File sourcefile, FileCollection sourceDirs, FileCollection outDirs) {

        Datamodel.FileNode sourceNode = getFileNode(sourcefile.toPath());

        Optional<File> actualSourceDir = sourceDirs.filter(sd -> isActualSourceDir(sd, sourcefile)).getFiles().stream().findFirst();

        try {
            Path relSrcPath = actualSourceDir.orElseThrow(FileNotFoundException::new).toPath().relativize(sourcefile.toPath());
            String[] filename = relSrcPath.getFileName().toString().split("\\.");
            filename[filename.length-1] = "class";
            String className = String.join(".", filename);
            Path packageDirs = relSrcPath.getParent();
            Path classesRelPath = packageDirs.resolve(className);

            if (packageDirs != null) {
                Optional<Datamodel.FileNode> root = outDirs.filter(od -> isActualClassDir(od, classesRelPath)).getFiles().stream()
                        .map(outdir -> {
                            Path classesPath = outdir.toPath().resolve(classesRelPath);
                            Datamodel.FileNode rootNode = getFileNode(classesPath);
                            return rootNode.toBuilder().addDerivedFrom(sourceNode).build();
                        }).findFirst();
                if (root.isPresent()) {
                    return root.get();
                }
            }
        } catch (FileNotFoundException fnfe) {
            //TODO
        }
        return null;
    }

    private static boolean isActualSourceDir(File sourceDir, File sourceFile) {
        return sourceFile.toString().startsWith(sourceDir.toString());
    }

    private static boolean isActualClassDir(File outdir, Path classesPath) {
        return outdir.toPath().resolve(classesPath).toFile().exists();
    }
}
