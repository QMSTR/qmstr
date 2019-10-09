package org.qmstr.util.transformations;

import java.io.File;
import java.util.Collection;
import java.util.Collections;
import java.util.Set;
import java.util.stream.Collectors;

import org.qmstr.grpc.service.Datamodel.FileNode;
import org.qmstr.grpc.service.Datamodel.FileNode.Builder;
import org.qmstr.util.FilenodeUtils;

public class JarTransformation implements TransformationFunction {

    private final String resultFilename;

    public JarTransformation(String resultFilename) {
        this.resultFilename = resultFilename;
    }

    @Override
	public Set<FileNode> apply(Collection<File> src, Collection<File> srcDirs, Collection<File> out)
			throws TransformationException {

        Set<File> inputClasses = srcDirs.stream()
            .filter(f -> (f.toPath().toString().endsWith(".class") || f.toPath().toString().endsWith(".jar")))
            .collect(Collectors.toSet());

        if (inputClasses.isEmpty()) {
            throw new TransformationException("No input for jar transformation.");
        }

        return Collections.singleton(out.stream()
                .filter(f -> f.toPath().getFileName().toString().endsWith(this.resultFilename))
                .map(f -> FilenodeUtils.getFileNode(f.toPath()))
                .map(p -> {
                    Builder b = p.toBuilder();
                    inputClasses.stream()
                        .map(clazz -> FilenodeUtils.getFileNode(clazz.toPath()))
                        .forEach(classFileNode -> b.addDerivedFrom(classFileNode));
                    return b.build();
                })
                .findFirst()
                .orElseThrow(() -> new TransformationException(String.format("no %s found in %s", this.resultFilename, out))));
    }
}