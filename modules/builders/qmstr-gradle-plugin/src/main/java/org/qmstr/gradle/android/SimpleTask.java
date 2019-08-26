package org.qmstr.gradle.android;

import java.util.regex.Pattern;

public enum SimpleTask {
    COMPILEJAVA, DEX, MERGEDEX, PACKAGEAPK, NONE;

    static Pattern dex = Pattern.compile("transformClassesWithDexBuilder\\w*");
    static Pattern mergeDex = Pattern.compile("mergeDex\\w*");
    static Pattern compileJava = Pattern.compile("compile\\w*Java\\w*");
    static Pattern packageApk = Pattern.compile("package\\w*");

    public static SimpleTask detectTask(String taskName) {
        if (compileJava.matcher(taskName).matches()) {
            return COMPILEJAVA;
        }
        if (dex.matcher(taskName).matches()) {
            return DEX;
        } 
        if (mergeDex.matcher(taskName).matches()) {
            return MERGEDEX;
        }
        if (packageApk.matcher(taskName).matches()) {
            return PACKAGEAPK;
        }
        return NONE;
    }
}