package org.qmstr.util;

import org.gradle.api.artifacts.PublishArtifact;
import org.gradle.api.file.FileCollection;
import org.qmstr.grpc.service.Datamodel;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.HashSet;
import java.util.Optional;
import java.util.Set;
import java.util.jar.JarEntry;
import java.util.jar.JarFile;
import java.util.stream.Collectors;

public class FilenodeUtils {

    private static final String[] SUPPORTEDFILES = new String[]{"java", "class", "jar"};

    public static Datamodel.FileNode getFileNode(String path, String checksum, String type) {
        Path filepath = Paths.get(path);

        return Datamodel.FileNode.newBuilder()
                .setName(filepath.getFileName().toString())
                .setPath(filepath.toString())
                .setHash(checksum != null ? checksum : "nohash"+filepath.toString())
                .setBroken(checksum == null)
                .setType(type)
                .build();

    }

    public static Datamodel.FileNode getFileNode(Path filepath, String type) {
        String checksum = Hash.getChecksum(filepath.toFile());
        String path = filepath.toString();

        return getFileNode(path, checksum, type);
    }

    public static Optional<Datamodel.FileNode> getFileNode(Path filepath) {
        if (isSupportedFile(filepath.toString())) {
            return Optional.of(getFileNode(filepath, getTypeByFile(filepath.toString())));
        }
        return Optional.empty();
    }


    public static Set<Datamodel.FileNode> processSourceFile(File sourcefile, FileCollection sourceDirs, FileCollection outDirs) {

        Datamodel.FileNode sourceNode = getFileNode(sourcefile.toPath(), getTypeByFile(sourcefile.getName()));

        Optional<File> actualSourceDir = sourceDirs
                .filter(sd -> isActualSourceDir(sd, sourcefile))
                .getFiles()
                .stream()
                .findFirst();

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
                                    .map(p -> getFileNode(p, getTypeByFile(p.getFileName().toString())))
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

    public static Optional<Datamodel.FileNode> processArtifact(PublishArtifact artifact, Set<File> dependencySet) {
        if (artifact.getExtension().equals("jar")) {
            try {
                Set<Datamodel.FileNode> classes = new HashSet<>();
                JarFile jar = new JarFile(artifact.getFile());
                jar.stream().parallel()
                        .filter(je -> isSupportedFile(je.getName()))
                        .forEach(je -> {
                            String hash = getHash(jar, je);
                            classes.add(FilenodeUtils.getFileNode(je.getName(), hash, getTypeByFile(je.getName())));
                        });
                Datamodel.FileNode rootNode = FilenodeUtils.getFileNode(artifact.getFile().toPath(), getTypeByFile(artifact.getFile().getName()));
                Datamodel.FileNode.Builder rootNodeBuilder = rootNode.toBuilder();
                classes.forEach(c -> rootNodeBuilder.addDerivedFrom(c));

                dependencySet.parallelStream()
                        .map(f -> FilenodeUtils.getFileNode(f.toPath()))
                        .filter(o -> o.isPresent())
                        .map(o -> o.get())
                        .forEach(depNode -> rootNodeBuilder.addDerivedFrom(depNode));

                rootNode = rootNodeBuilder.build();
                return Optional.ofNullable(rootNode);

            } catch (IOException ioe) {
                //TODO
            }
        }
        return Optional.empty();
    }

    private static boolean isSupportedFile(String filename) {
        String[] filenameArr = filename.split("\\.");
        int idx = filenameArr.length > 0 ? filenameArr.length-1 : 0;
        return Arrays.stream(SUPPORTEDFILES).anyMatch(sf -> sf.equals(filenameArr[idx]));
    }

    private static String getTypeByFile(String filename) {
        String[] filenameArr = filename.split("\\.");
        String ext = filenameArr[filenameArr.length-1];
        if (ext.equals("class")) {
            return "classfile";
        }
        if (ext.equals("java")) {
            return "sourcecode";
        }
        if (ext.equals("jar")) {
            return "jarfile";
        }
        return "";
    }

    private static String getHash(JarFile jarfile, JarEntry jarEntry) {
        try {
            InputStream is = jarfile.getInputStream(jarEntry);
            return Hash.getChecksum(is);
        } catch (IOException e) {
            e.printStackTrace();
        }
        return null;
    }
}
