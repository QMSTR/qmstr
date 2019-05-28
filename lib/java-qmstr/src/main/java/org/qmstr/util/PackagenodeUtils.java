package org.qmstr.util;

import java.io.File;
import java.io.IOException;
import java.nio.file.FileSystems;
import java.nio.file.PathMatcher;
import java.util.HashSet;
import java.util.Optional;
import java.util.Set;
import java.util.jar.JarFile;

import org.qmstr.grpc.service.Datamodel;

public class PackagenodeUtils {
    public static Optional<Datamodel.PackageNode> processArtifact(File artifact, String packageName, String version) {
        PathMatcher jarMatcher = FileSystems.getDefault().getPathMatcher("glob:**.jar");
        if (jarMatcher.matches(artifact.toPath())) {
            try (JarFile jar = new JarFile(artifact)){
                Set<Datamodel.FileNode> classes = new HashSet<>();
                jar.stream().parallel()
                        .filter(je -> FilenodeUtils.isSupportedFile(je.getName()))
                        .forEach(je -> {
                            String hash = FilenodeUtils.getHash(jar, je);
                            classes.add(FilenodeUtils.getFileNode(je.getName(), hash, FilenodeUtils.getTypeByFile(je.getName())));
                        });
                Datamodel.PackageNode rootNode = getPackageNode(packageName, version);
                Datamodel.PackageNode.Builder rootNodeBuilder = rootNode.toBuilder();
                classes.forEach(c -> rootNodeBuilder.addTargets(c));

                rootNode = rootNodeBuilder.build();
                return Optional.ofNullable(rootNode);

            } catch (IOException ioe) {
                // Default to returning empty optional
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
