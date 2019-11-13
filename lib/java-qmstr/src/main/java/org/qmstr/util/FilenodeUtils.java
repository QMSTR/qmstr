package org.qmstr.util;

import java.io.File;
import java.io.FileNotFoundException;
import java.io.IOException;
import java.io.InputStream;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Arrays;
import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.Optional;
import java.util.Set;
import java.util.jar.JarEntry;
import java.util.jar.JarFile;

import org.qmstr.grpc.service.Datamodel;
import org.qmstr.util.transformations.CompileJavaTransformation;
import org.qmstr.util.transformations.DexClassTransformation;
import org.qmstr.util.transformations.JarTransformation;
import org.qmstr.util.transformations.MergeDexTransformation;
import org.qmstr.util.transformations.Transform;
import org.qmstr.util.transformations.TransformationException;
import org.qmstr.util.transformations.TransformationFunction;

public class FilenodeUtils {

    private static final String[] SUPPORTEDFILES = new String[] { "java", "class", "jar", "dex" };

    private static final Map<Transform, TransformationFunction> srcDestMap = new HashMap<Transform, TransformationFunction>() {
        {
            put(Transform.COMPILEJAVA, new CompileJavaTransformation());
            put(Transform.DEXCLASS, new DexClassTransformation());
            put(Transform.MERGEDEX, new MergeDexTransformation());
            put(Transform.PACKAGECLASSESJAR, new JarTransformation("classes.jar"));
            put(Transform.PACKAGEFULLJAR, new JarTransformation("full.jar"));
        }
    };

    public static Datamodel.FileNode getFileNode(String path, String checksum) {
        Path filepath = Paths.get(path).toAbsolutePath();
        return Datamodel.FileNode.newBuilder().setName(filepath.getFileName().toString()).setPath(filepath.toString())
                .setFileData(Datamodel.FileNode.FileDataNode.newBuilder()
                        .setHash(checksum != null ? checksum : "nohash" + filepath.toString()).build())
                .build();

    }

    public static Datamodel.FileNode getFileNode(Path filepath) {
        filepath = filepath.toAbsolutePath();
        String checksum = Hash.getChecksum(filepath.toFile());
        String path = filepath.toString();

        return getFileNode(path, checksum);
    }


    public static boolean isSupportedFile(String filename) {
        String[] filenameArr = filename.split("\\.");
        int idx = filenameArr.length > 0 ? filenameArr.length - 1 : 0;
        return Arrays.stream(SUPPORTEDFILES).anyMatch(sf -> sf.equals(filenameArr[idx]));
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
