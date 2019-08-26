package org.qmstr.util;

import org.qmstr.util.transformations.*;
import org.qmstr.grpc.service.Datamodel;
import org.qmstr.grpc.service.Datamodel.FileNode.Type;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStream;
import java.nio.file.*;
import java.util.*;
import java.util.jar.JarEntry;
import java.util.jar.JarFile;
import java.util.stream.Collectors;

public class FilenodeUtils {

    private static final String[] SUPPORTEDFILES = new String[] { "java", "class", "jar", "dex" };

    private static final Map<Transform, TransformationFunction> srcDestMap = new HashMap<Transform, TransformationFunction>() {
        {
            put(Transform.COMPILEJAVA, new CompileJavaTransformation());
            put(Transform.DEXCLASS, new DexClassTransformation());
            put(Transform.MERGEDEX, new MergeDexTransformation());
        }
    };

    public static Datamodel.FileNode getFileNode(String path, String checksum, Datamodel.FileNode.Type type) {
        Path filepath = Paths.get(path);

        return Datamodel.FileNode.newBuilder().setName(filepath.getFileName().toString()).setPath(filepath.toString())
                .setHash(checksum != null ? checksum : "nohash" + filepath.toString()).setBroken(checksum == null)
                .setFileType(type).build();

    }

    public static Datamodel.FileNode getFileNode(Path filepath, Datamodel.FileNode.Type type) {
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

    public static boolean isSupportedFile(String filename) {
        String[] filenameArr = filename.split("\\.");
        int idx = filenameArr.length > 0 ? filenameArr.length - 1 : 0;
        return Arrays.stream(SUPPORTEDFILES).anyMatch(sf -> sf.equals(filenameArr[idx]));
    }

    @Deprecated
    public static Datamodel.FileNode.Type getTypeByFile(String filename) {
        String[] filenameArr = filename.split("\\.");
        String ext = filenameArr[filenameArr.length - 1];
        if (ext.equals("class")) {
            return Type.INTERMEDIATE;
        }
        if (ext.equals("java")) {
            return Type.SOURCE;
        }
        if (ext.equals("jar")) {
            return Type.TARGET;
        }
        return Type.UNDEF;
    }

    public static String getHash(JarFile jarfile, JarEntry jarEntry) {
        try {
            InputStream is = jarfile.getInputStream(jarEntry);
            return Hash.getChecksum(is);
        } catch (IOException e) {
            e.printStackTrace();
        }
        return null;
    }

    public static Set<Datamodel.FileNode> processSourceFiles(Transform transform, Collection<File> sourceFiles,
            Collection<File> sourceDirs, Collection<File> outDirs)
            throws TransformationException, FileNotFoundException {

        return srcDestMap.get(transform).apply(sourceFiles, sourceDirs, outDirs);
    }

    public static boolean isActualSourceDir(File sourceDir, File sourceFile) {
        return sourceFile.toString().startsWith(sourceDir.toString());
    }

    public static boolean isActualClassDir(File outdir, Path classesPath) {
        return outdir.toPath().resolve(classesPath).toFile().exists();
    }


}
