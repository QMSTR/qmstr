package org.qmstr.util;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.io.File;

import org.junit.jupiter.api.Test;

class HashTests{

    @Test
    void checksumTest() {
            File hashtest = new File(getClass().getResource("/hashtest").getFile());
            assertEquals("45e51db6f37b0b8af21c7822ed4b470e6565f931", Hash.getChecksum(hashtest));
    }

}