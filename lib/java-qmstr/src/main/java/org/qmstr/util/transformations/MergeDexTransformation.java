package org.qmstr.util.transformations;

import java.io.File;
import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.attribute.BasicFileAttributes;
import java.util.Collection;
import java.util.Collections;
import java.util.Set;
import java.util.function.BiPredicate;
import java.util.stream.Collectors;
import java.util.stream.Stream;

import org.qmstr.grpc.service.Datamodel.FileNode;
import org.qmstr.grpc.service.Datamodel.FileNode.Builder;
import org.qmstr.util.FilenodeUtils;

public class MergeDexTransformation implements TransformationFunction {

    @Override
	public Set<FileNode> apply(Collection<File> src, Collection<File> srcDirs, Collection<File> outDirs)
			throws TransformationException {

        Set<File> inputDexes = srcDirs.stream()
            .flatMap(srcDir -> wrapFind(
                srcDir.toPath(), 
                (path,attrs) -> attrs.isRegularFile() && path.toString().endsWith(".dex")
                )
            )
            .map(p -> p.toFile())
            .collect(Collectors.toSet());

        return Collections.singleton(outDirs.stream()
                .map(od -> od.toPath().resolve("classes.dex"))
                .map(p -> FilenodeUtils.getFileNode(p, FileNode.Type.UNDEF))
                .map(p -> {
                    Builder b = p.toBuilder();
                    inputDexes.stream()
                        .map(dex -> FilenodeUtils.getFileNode(dex.toPath(), FileNode.Type.UNDEF))
                        .forEach(dexFileNode -> b.addDerivedFrom(dexFileNode));
                    return b.build();
                })
                .findFirst()
                .orElseThrow(() -> new TransformationException(String.format("no classes.dex found in %s", outDirs))));

    }

    public static Stream<Path> wrapFind(Path start, BiPredicate<Path, BasicFileAttributes> matcher) {
        try {
            return Files.find(start, Integer.MAX_VALUE, matcher);
        } catch (IOException e) {
            throw new RuntimeException(e);
		}
    }
}