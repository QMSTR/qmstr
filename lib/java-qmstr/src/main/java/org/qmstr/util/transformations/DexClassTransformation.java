package org.qmstr.util.transformations;

public class DexClassTransformation implements TransformationFunction<String, String> {

    @Override
    public String apply(String s) throws TransformationException {
        String[] filename = s.split("\\.");
        if (filename[filename.length-1] != "class") {
            throw new TransformationException("Invalid input; must be class file");
        }
        filename[filename.length-1] = "dex";
        return String.join(".", filename);
    }
}