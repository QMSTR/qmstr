package org.qmstr.util.transformations;

public class DexClassTransformation implements TransformationFunction<String, String> {

    @Override
    public String apply(String s) throws TransformationException {
        String[] filename = s.split("\\.");
        String extension = filename[filename.length-1];
        if (!extension.equals("class")) {
            throw new TransformationException(String.format("Invalid input %s; must be class file", s));
        }
        filename[filename.length-1] = "dex";
        return String.join(".", filename);
    }
}