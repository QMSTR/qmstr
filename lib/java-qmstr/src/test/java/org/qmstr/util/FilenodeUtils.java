package org.qmstr.util;

import static org.junit.jupiter.api.Assertions.assertEquals;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.junit.jupiter.api.Assertions.assertTrue;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

import org.junit.jupiter.api.Test;
import org.qmstr.util.FilenodeUtils.TRANSFORM;
import org.qmstr.util.transformations.TransformationException;

class FilenodeUtilsTests {

    @Test
    void compileJavaTest() throws TransformationException {
        String target = FilenodeUtils.getDestination(TRANSFORM.COMPILEJAVA, "Test.java");
        assertEquals("Test.class", target);
    }

    @Test
    void dexClassTest() throws TransformationException {
        String target = FilenodeUtils.getDestination(TRANSFORM.DEXCLASS, "Test.class");
        assertEquals("Test.dex", target);
    }

    @Test
    void dexArchiveTest() throws TransformationException {
        String target = FilenodeUtils.getDestination(TRANSFORM.DEXARCHIVE, "Test.dex");
        assertEquals("classes.dex", target);
    }

    @Test
    void isNestedTest() throws IOException {
        Path tmpDir = Files.createTempDirectory("NestedClassTest");
        Path testClass = Files.createFile(tmpDir.resolve("OuterClass$NestedClass.class"));
        boolean nested = FilenodeUtils.isNestedClass(testClass, "OuterClass.class");
        tmpDir.toFile().deleteOnExit();
        assertTrue(nested);
    }
}