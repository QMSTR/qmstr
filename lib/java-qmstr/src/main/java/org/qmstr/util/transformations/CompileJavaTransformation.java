package org.qmstr.util.transformations;

public class CompileJavaTransformation implements TransformationFunction<String, String> {

    @Override
    public String apply(String s) throws TransformationException {
        String[] filename = s.split("\\.");
        String extension = filename[filename.length-1];
        if (!extension.equals("java")) {
            throw new TransformationException(String.format("Invalid input %s; must be java file", s));
        }
        filename[filename.length-1] = "class";
        return String.join(".", filename);
    }
}