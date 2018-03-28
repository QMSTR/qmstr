package org.qmstr.util;

import groovy.ui.SystemOutputInterceptor;
import org.gradle.api.Project;
import org.gradle.api.file.FileCollection;
import org.gradle.api.tasks.SourceSetOutput;
import org.qmstr.grpc.service.Datamodel;

import javax.xml.crypto.Data;
import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.HashSet;
import java.util.Optional;
import java.util.Set;
import java.util.stream.Collectors;

public class FilenodeUtils {


    public static Datamodel.FileNode getFileNode(String path, String checksum, String type) {
        Path filepath = Paths.get(path);

        return Datamodel.FileNode.newBuilder()
                .setName(filepath.getFileName().toString())
                .setPath(filepath.toString())
                .setHash(checksum != null ? checksum : "nohash"+filepath.toString())
                .setBroken(checksum == null)
                .setNodeType(1)
                .setType(type)
                .build();

    }

    public static Datamodel.FileNode getFileNode(Path filepath, String type) {
        String checksum = Hash.getChecksum(filepath.toFile());
        String path = filepath.toString();

        return getFileNode(path, checksum, type);
    }

    public static Set<Datamodel.FileNode> processSourceFile(File sourcefile, FileCollection sourceDirs, FileCollection outDirs) {

        Datamodel.FileNode sourceNode = getFileNode(sourcefile.toPath(), "sourcecode");

        Optional<File> actualSourceDir = sourceDirs.filter(sd -> isActualSourceDir(sd, sourcefile)).getFiles().stream().findFirst();

        try {
            Path relSrcPath = actualSourceDir.orElseThrow(FileNotFoundException::new).toPath().relativize(sourcefile.toPath());
            String[] filename = relSrcPath.getFileName().toString().split("\\.");
            filename[filename.length-1] = "class";
            String className = String.join(".", filename);
            Path packageDirs = relSrcPath.getParent();
            Path classesRelPath = packageDirs.resolve(className);

            if (packageDirs != null) {
                return outDirs.filter(od -> isActualClassDir(od, classesRelPath)).getFiles().stream()
                        .map(outdir -> {
                            Path classesPath = outdir.toPath().resolve(classesRelPath);
                            Set<Path> nested = getNestedClasses(outdir.toPath().resolve(packageDirs), filename[filename.length - 2]);
                            nested.add(classesPath);
                            return nested.stream()
                                    .map(p -> getFileNode(p, "classfile"))
                                    .map(node -> node.toBuilder().addDerivedFrom(sourceNode).build())
                                    .collect(Collectors.toSet());
                        }).flatMap(sets -> sets.stream())
                        .collect(Collectors.toSet());
            }
        } catch (FileNotFoundException fnfe) {
            //TODO
        }
        return null;
    }

    private static Set<Path> getNestedClasses(Path dir, String outerclassname) {
        try {
            return Files.walk(dir)
                    .filter(p -> isNestedClass(p, outerclassname))
                    .collect(Collectors.toSet());

        } catch (IOException e) {
            e.printStackTrace();
            return new HashSet<>();
        }
    }

    private static boolean isActualSourceDir(File sourceDir, File sourceFile) {
        return sourceFile.toString().startsWith(sourceDir.toString());
    }

    private static boolean isActualClassDir(File outdir, Path classesPath) {
        return outdir.toPath().resolve(classesPath).toFile().exists();
    }

    private static boolean isNestedClass(Path classesPath, String outerClass) {
        boolean file = classesPath.toFile().isFile();
        String filename = classesPath.getFileName().toString();
        boolean clazz = filename.endsWith(".class");
        boolean starts = filename.startsWith(outerClass + "$");
        return file && clazz && starts;
    }
}
