package org.qmstr.util;

import javax.xml.bind.DatatypeConverter;
import java.io.*;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;
import java.util.Arrays;

public class Hash {

    static final int hashByteBuffer = 4096;
    static final String hashAlgo = "SHA-256";
    static final String HEXES = "0123456789abcdef";

    public static String getChecksum(File inputFile) {

        try {
            BufferedInputStream bs = new BufferedInputStream(new FileInputStream(inputFile));
            MessageDigest digest = MessageDigest.getInstance(hashAlgo);
            byte[] buffer = new byte[hashByteBuffer];
            int read;

            while ((read = bs.read(buffer)) > 0) {
                digest.update(buffer, 0, read);
            }

            // double the size is a rough estimation
            StringBuilder sb = new StringBuilder(digest.getDigestLength() * 2);

            for (byte b : digest.digest()) {
                sb.append(HEXES.charAt((b & 0xF0) >> 4))
                  .append(HEXES.charAt((b & 0x0F)));
            }

            bs.close();

            return sb.toString();

        } catch (NoSuchAlgorithmException nsae) {
            nsae.printStackTrace();
            return null;
        } catch (IOException ioe) {
            ioe.printStackTrace();
            return null;
        }
    }
}
