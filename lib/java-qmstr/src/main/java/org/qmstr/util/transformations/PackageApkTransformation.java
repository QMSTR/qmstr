package org.qmstr.util.transformations;

import java.io.File;
import java.util.Collection;
import java.util.Set;

import org.qmstr.grpc.service.Datamodel.FileNode;

public class PackageApkTransformation implements TransformationFunction {

    @Override
	public Set<FileNode> apply(Collection<File> src, Collection<File> srcDirs, Collection<File> outDirs)
			throws TransformationException {

        return null;
    }

}