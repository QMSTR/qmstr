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

    private static final String[] SUPPORTEDFILES = new String[] { "java", "class", "jar" };

    private static final Map<Transform, TransformationFunction<String, String>> srcDestMap = new HashMap<Transform, TransformationFunction<String, String>>() {
        {
            put(Transform.COMPILEJAVA, new CompileJavaTransformation());
            put(Transform.DEXCLASS, new DexClassTransformation());
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

    public static String getDestination(Transform transform, String source) throws TransformationException {
        TransformationFunction<String, String> transformFct = srcDestMap.getOrDefault(transform, t -> "");
        return transformFct.apply(source);
    }

    public static Set<Datamodel.FileNode> processSourceFile(Transform transform, File sourcefile,
            Collection<File> sourceDirs, Collection<File> outDirs)
            throws TransformationException, FileNotFoundException {
        
        if (sourcefile.isDirectory()) {
            return Collections.emptySet();
        }

        Datamodel.FileNode sourceNode = FilenodeUtils.getFileNode(sourcefile.toPath(),
                FilenodeUtils.getTypeByFile(sourcefile.getName()));

        Optional<File> actualSourceDir = sourceDirs.stream().filter(sd -> isActualSourceDir(sd, sourcefile))
                .findFirst();

        Path relSrcPath = actualSourceDir
                .orElseThrow(() -> new FileNotFoundException(
                        String.format("No source dir found for %s", sourcefile.getAbsolutePath())))
                .toPath().relativize(sourcefile.toPath());

        String targetFileName = getDestination(transform, relSrcPath.getFileName().toString());

        Path packageDirs = relSrcPath.getParent() != null ? relSrcPath.getParent() : Paths.get(".");
        Path classesRelPath = packageDirs.resolve(targetFileName);

        return outDirs.stream().filter(od -> FilenodeUtils.isActualClassDir(od, classesRelPath)).map(outdir -> {
            Path classesPath = outdir.toPath().resolve(classesRelPath);
            Set<Path> nested = FilenodeUtils.getNestedClasses(outdir.toPath().resolve(packageDirs), targetFileName);
            nested.add(classesPath);
            return nested.stream()
                    .map(p -> FilenodeUtils.getFileNode(p, FilenodeUtils.getTypeByFile(p.getFileName().toString())))
                    .map(node -> node.toBuilder().addDerivedFrom(sourceNode).build()).collect(Collectors.toSet());
        }).flatMap(sets -> sets.stream()).collect(Collectors.toSet());
    }

    private static boolean isActualSourceDir(File sourceDir, File sourceFile) {
        return sourceFile.toString().startsWith(sourceDir.toString());
    }

    private static boolean isActualClassDir(File outdir, Path classesPath) {
        return outdir.toPath().resolve(classesPath).toFile().exists();
    }

    private static Set<Path> getNestedClasses(Path dir, String outerclassname) {
        try {
            return Files.walk(dir).filter(p -> isNestedClass(p, outerclassname)).collect(Collectors.toSet());

        } catch (IOException e) {
            e.printStackTrace();
            return new HashSet<>();
        }
    }

    public static boolean isNestedClass(Path classesPath, String outerClass) {
        String outerClassName = outerClass.replaceAll(".class$", "");
        boolean file = classesPath.toFile().isFile();
        String filename = classesPath.getFileName().toString();
        boolean clazz = filename.endsWith(".class");
        boolean starts = filename.startsWith(outerClassName + "$");
        return file && clazz && starts;
    }

}
