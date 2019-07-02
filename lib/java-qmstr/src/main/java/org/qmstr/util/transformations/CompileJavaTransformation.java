package org.qmstr.util.transformations;

import java.util.function.Function;

public class CompileJavaTransformation implements TransformationFunction<String, String> {

    @Override
    public String apply(String s) throws TransformationException {
        String[] filename = s.split("\\.");
        if (filename[filename.length-1] != "java") {
            throw new TransformationException("Invalid input; must be java file");
        }
        filename[filename.length-1] = "class";
        return String.join(".", filename);
    }
}