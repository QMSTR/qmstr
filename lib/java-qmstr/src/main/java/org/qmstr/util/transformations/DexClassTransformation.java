package org.qmstr.util.transformations;

import static org.qmstr.util.FilenodeUtils.isActualSourceDir;
import static org.qmstr.util.FilenodeUtils.isActualClassDir;
import java.io.File;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Collection;
import java.util.Collections;
import java.util.Set;

import org.qmstr.grpc.service.Datamodel.FileNode;
import org.qmstr.util.FilenodeUtils;

public class DexClassTransformation implements TransformationFunction {

    @Override
	public Set<FileNode> apply(Collection<File> sources, Collection<File> sourceDirs, Collection<File> outDirs)
			throws TransformationException {
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
        if (!extension.equals("class")) {
            throw new TransformationException(String.format("invalid input %s; must be class file", sourceFile.getPath()));
        }
        filename[filename.length-1] = "dex";
        String targetFileNamePath = String.join(".", filename);
        Path targetFileName = Paths.get(targetFileNamePath).getFileName();

        Path packageDirs = relSrcPath.getParent() != null ? relSrcPath.getParent() : Paths.get(".");
        Path classesRelPath = packageDirs.resolve(targetFileName);

        File destinationDir =  outDirs.stream()
            .filter(od -> isActualClassDir(od, classesRelPath))
            .findFirst()
            .orElseThrow(() -> new TransformationException(String.format("target class %s was not found", classesRelPath.toString())));

        Path classesPath = destinationDir.toPath().resolve(classesRelPath);

        FileNode sourceNode = FilenodeUtils.getFileNode(sourceFile.toPath(), FilenodeUtils.getTypeByFile(sourceFile.getName()));

        FileNode dexFileNode = FilenodeUtils.getFileNode(classesPath, FilenodeUtils.getTypeByFile(classesPath.getFileName().toString()));

        return Collections.singleton(dexFileNode.toBuilder().addDerivedFrom(sourceNode).build());
	}
}