package org.qmstr.util.transformations;

import static org.qmstr.util.FilenodeUtils.isActualSourceDir;
import static org.qmstr.util.FilenodeUtils.isActualClassDir;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Collection;
import java.util.Collections;
import java.util.Set;
import java.util.stream.Collectors;

import org.qmstr.grpc.service.Datamodel.FileNode;
import org.qmstr.util.FilenodeUtils;

public class CompileJavaTransformation implements TransformationFunction {

    @Override
    public Set<FileNode> apply(Collection<File> sources, Collection<File> sourceDirs, Collection<File> outDirs) throws TransformationException {
        if (sources.size() != 1) {
            throw new TransformationException(String.format("invalid number of source files %d; must be 1", sources.size()));
        }

        File sourceFile = sources.stream().findFirst().orElseThrow(() -> new TransformationException("failed to get source file"));

        File actualSourceDir = sourceDirs.stream()
                .filter(sd -> isActualSourceDir(sd, sourceFile))
                .findFirst()
                .orElseThrow(() -> new TransformationException(
                        String.format("No source dir found for %s", sourceFile.getAbsolutePath())));

        Path relSrcPath = actualSourceDir
                .toPath().relativize(sourceFile.toPath());

        String[] filename = relSrcPath.toString().split("\\.");
        String extension = filename[filename.length-1];
        if (!extension.equals("java")) {
            throw new TransformationException(String.format("invalid input %s; must be java file", sourceFile.getPath()));
        }
        filename[filename.length-1] = "class";
        String targetFileNamePath = String.join(".", filename);
        Path targetFileName = Paths.get(targetFileNamePath).getFileName();

        Path packageDirs = relSrcPath.getParent() != null ? relSrcPath.getParent() : Paths.get(".");
        Path classesRelPath = packageDirs.resolve(targetFileName);

        File destinationDir =  outDirs.stream()
            .filter(od -> isActualClassDir(od, classesRelPath))
            .findFirst()
            .orElseThrow(() -> new TransformationException(String.format("target class %s was not found", classesRelPath.toString())));

        Path classesPath = destinationDir.toPath().resolve(classesRelPath);
        Set<File> nested = getNestedClasses(destinationDir.toPath().resolve(packageDirs), targetFileName);
        nested.add(classesPath.toFile());

        FileNode sourceNode = FilenodeUtils.getFileNode(sourceFile.toPath(), FilenodeUtils.getTypeByFile(sourceFile.getName()));

        return nested.stream()
            .map(df -> FilenodeUtils.getFileNode(df.toPath(), FilenodeUtils.getTypeByFile(df.toPath().getFileName().toString())))
            .map(dfnode -> dfnode.toBuilder().addDerivedFrom(sourceNode).build())
            .collect(Collectors.toSet());
    }

    private static Set<File> getNestedClasses(Path dir, Path outerclass) {
        try {
            return Files.walk(dir)
                .filter(p -> isNestedClass(p, outerclass.toString()))
                .map(p -> p.toFile())
                .collect(Collectors.toSet());

        } catch (IOException e) {
            e.printStackTrace();
        }
        return Collections.emptySet();
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