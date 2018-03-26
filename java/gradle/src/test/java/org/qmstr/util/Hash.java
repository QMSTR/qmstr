package org.qmstr.util;

import static org.junit.jupiter.api.Assertions.assertEquals;

        import org.junit.jupiter.api.Test;

import java.io.File;
import java.io.FileNotFoundException;

class HashTests{

    @Test
    void checksumTest() {
            File hashtest = new File(getClass().getResource("/hashtest").getFile());
            assertEquals("e72b4044adf37a8a925ee6b046fe2ee3146ede625f9f89636957a35e10be83aa", Hash.getChecksum(hashtest));
    }

}