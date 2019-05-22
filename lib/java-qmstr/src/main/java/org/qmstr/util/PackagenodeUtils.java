package org.qmstr.util;

import org.qmstr.grpc.service.Datamodel;

import java.io.File;
import java.io.IOException;
import java.nio.file.FileSystems;
import java.nio.file.PathMatcher;
import java.util.HashSet;
import java.util.Optional;
import java.util.Set;
import java.util.jar.JarFile;

public class PackagenodeUtils {
    public static Optional<Datamodel.PackageNode> processArtifact(File artifact, String version) {
        PathMatcher jarMatcher = FileSystems.getDefault().getPathMatcher("glob:*.jar");
        if (jarMatcher.matches(artifact.toPath())) {
            try {
                Set<Datamodel.FileNode> classes = new HashSet<>();
                JarFile jar = new JarFile(artifact);
                jar.stream().parallel()
                        .filter(je -> FilenodeUtils.isSupportedFile(je.getName()))
                        .forEach(je -> {
                            String hash = FilenodeUtils.getHash(jar, je);
                            classes.add(FilenodeUtils.getFileNode(je.getName(), hash, FilenodeUtils.getTypeByFile(je.getName())));
                        });
                Datamodel.PackageNode rootNode = getPackageNode(artifact.toPath().getFileName().toString(), version);
                Datamodel.PackageNode.Builder rootNodeBuilder = rootNode.toBuilder();
                classes.forEach(c -> rootNodeBuilder.addTargets(c));

                rootNode = rootNodeBuilder.build();
                return Optional.ofNullable(rootNode);

            } catch (IOException ioe) {
                //TODO
            }
        }
        return Optional.empty();
    }

    public static Datamodel.PackageNode getPackageNode(String name, String version) {

        return Datamodel.PackageNode.newBuilder()
                .setName(name)
                .setVersion(version)
                .build();
    }
}
