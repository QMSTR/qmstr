package org.qmstr.gradle.android;

import java.io.File;
import java.io.FileNotFoundException;
import java.util.Collections;
import java.util.Set;
import java.util.stream.Collectors;

import org.gradle.api.Project;
import org.gradle.api.Task;
import org.gradle.api.plugins.AppliedPlugin;
import org.gradle.internal.impldep.com.fasterxml.jackson.databind.node.BooleanNode;
import org.qmstr.client.BuildServiceClient;
import org.qmstr.gradle.QmstrPluginExtension;
import org.qmstr.gradle.ResultUnavailableException;
import org.qmstr.grpc.service.Datamodel;
import org.qmstr.util.FilenodeUtils;
import org.qmstr.util.transformations.Transform;
import org.qmstr.util.transformations.TransformationException;

import static org.qmstr.gradle.android.TaskExecutionGraphReadyAction.getLibSourceDirs;
import static org.qmstr.gradle.android.TaskExecutionGraphReadyAction.getAppSourceDirs;;

public class AndroidPostTaskAction extends AndroidTaskAction {

    private final boolean lib;

    public AndroidPostTaskAction(Project project, AppliedPlugin plugin) {
        this.project = project;

        this.lib = plugin.getId().equals("com.android.library");

        qmstrExt = (QmstrPluginExtension) this.project.getExtensions().findByName("qmstr");
        this.setBuildServiceAddress(qmstrExt.qmstrAddress);
        this.bsc = new BuildServiceClient(buildServiceAddress, buildServicePort);
    }

    // This method tries to guess the root of the package hierarchy
    // The compile task will put the package hierarchy in a classes dir. So this method will try to step up until it finds a classes dir.
    // Beware that this is nothing but stupid and only for the PoC.
    private File guessSourcePath(File sourceFile) {
        if (sourceFile.getParentFile() == null || sourceFile.getParentFile().toPath().getFileName() == null) {
            return sourceFile;
        }
        if (sourceFile.getParentFile().toPath().getFileName().toString().equals("classes")) {
            return sourceFile.getParentFile();
        }
        return guessSourcePath(sourceFile.getParentFile());
    }
    
    private void handleDex(Task task) {
        // classes are dexed and can be collected
        task.getOutputs().getFiles().forEach(
                out -> task.getLogger().warn("Dex task {} output is {}", task.getName(), out.getAbsolutePath()));

        task.getInputs().getFiles().forEach(sf -> {
            Set<Datamodel.FileNode> nodes;
            try {
                // Here it becomes ugly. The output dir we get from the task is not yet the dir to look for the package file tree that holds the dex files.
                // There is another dir in the hierarchy. This might be due to multidex https://developer.android.com/studio/build/multidex
                // Anyway for now we just assume that there is one more dir called '0'
                Set<File> outdirs = task.getOutputs().getFiles().getFiles().stream()
                        .map(f -> f.toPath().resolve("0").toFile()).collect(Collectors.toSet());

                // Here it gets even uglier. The processSourceFile method assumes you have a source file and a set of input directories where your sources (inside a package hierarchy) reside.
                // This however is not the case here because we are not working with sourcesets like in a Java build.
                // Therefore we need to find the root of the package hierarchy from the filename. This is what the brain-damaged guessSourcePath method does.
                nodes = FilenodeUtils.processSourceFiles(Transform.DEXCLASS, Collections.singleton(sf),
                        Collections.singleton(guessSourcePath(sf)), outdirs);
                if (!nodes.isEmpty()) {
                    bsc.SendBuildFileNodes(nodes);
                } else {
                    bsc.SendLogMessage(String.format("No filenodes after processing %s", sf.getName()));
                }
            } catch (TransformationException e) {
                task.getLogger().warn("{} failed: {}", this.getClass().getName(), e.getMessage());
            } catch (FileNotFoundException fnfe) {
                task.getLogger().warn("{} failed: {}", this.getClass().getName(), fnfe.getMessage());
            }

        });
    }

    private void handleCompileJava(Task task) {

        Set<File> sourceDirs = lib ? getLibSourceDirs(task.getProject()) : getAppSourceDirs(task.getProject());

        task.getOutputs().getFiles().forEach(
                out -> task.getLogger().warn("Java compile task {} output is {}", task.getName(), out.getAbsolutePath()));
        
        
        Set<File> outDirs = task.getOutputs().getFiles().getFiles();

        task.getInputs().getFiles().forEach(sf -> {
            bsc.SendLogMessage("handle java compilation for " + sf.getAbsolutePath());
            Set<Datamodel.FileNode> nodes;
            try {
                nodes = FilenodeUtils.processSourceFiles(Transform.COMPILEJAVA, Collections.singleton(sf), sourceDirs, outDirs);
                if (!nodes.isEmpty()) {
                    bsc.SendBuildFileNodes(nodes);
                } else {
                    bsc.SendLogMessage(String.format("No filenodes after processing %s", sf.getName()));
                }
            } catch (TransformationException e) {
                task.getLogger().warn("{} failed: {}", this.getClass().getName(), e.getMessage());
            } catch (FileNotFoundException fnfe) {
                task.getLogger().warn("{} failed: {}", this.getClass().getName(), fnfe.getMessage());
            }

        });
    }

    private void handleMergeDex(Task task) throws ResultUnavailableException {
        task.getOutputs().getFiles().forEach(
                out -> task.getLogger().warn("Merge dex task {} output is {}", task.getName(), out.getAbsolutePath())
                );
        
        Set<File> classesDexDirs = task.getOutputs().getFiles().getFiles();

        task.getLogger().warn("Input files:");
        task.getInputs().getFiles().forEach(sf -> task.getLogger().warn(sf.getAbsolutePath()));

        Set<File> inputDirs = task.getInputs().getFiles().getFiles();
        
        Set<Datamodel.FileNode> nodes;
        try {
            nodes = FilenodeUtils.processSourceFiles(Transform.MERGEDEX, Collections.emptySet(), inputDirs, classesDexDirs);
            if (!nodes.isEmpty()) {
                bsc.SendBuildFileNodes(nodes);
            } else {
                bsc.SendLogMessage(String.format("No filenodes after processing %s", inputDirs));
            }
        } catch (TransformationException e) {
            task.getLogger().warn("{} failed: {}", this.getClass().getName(), e.getMessage());
        } catch (FileNotFoundException fnfe) {
            task.getLogger().warn("{} failed: {}", this.getClass().getName(), fnfe.getMessage());
        }
    } 

    @Override
    public void execute(Task task) {
        switch (SimpleTask.detectTask(task.getName())) {
            case COMPILEJAVA:
                handleCompileJava(task);
                break;
            case DEX:
                handleDex(task);
                break;
            case MERGEDEX:
                try {
                    handleMergeDex(task);
                } catch (ResultUnavailableException e) {
                    task.getLogger().error("Failed in merge dex");
                }
                break;
            default:
                break;
        }
    }
}