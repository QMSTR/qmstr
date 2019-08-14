package org.qmstr.gradle.android;

import java.util.regex.Pattern;

public enum SimpleTask {
    COMPILEJAVA, DEX, NONE;

    static Pattern dex = Pattern.compile("transformClassesWithDexBuilder\\w*");
    static Pattern compileJava = Pattern.compile("compile\\w*Java\\w*");

    public static SimpleTask detectTask(String taskName) {
        if (compileJava.matcher(taskName).matches()) {
            return COMPILEJAVA;
        }
        if (dex.matcher(taskName).matches()) {
            return DEX;
        } 
        return NONE;
    }
}