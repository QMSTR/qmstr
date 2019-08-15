package org.qmstr.util.transformations;

import java.io.File;
import java.util.Collection;
import java.util.Set;

import org.qmstr.grpc.service.Datamodel.FileNode;

@FunctionalInterface
public interface TransformationFunction {
   Set<FileNode> apply(Collection<File> src, Collection<File> srcDirs, Collection<File> outDirs) throws TransformationException;
}