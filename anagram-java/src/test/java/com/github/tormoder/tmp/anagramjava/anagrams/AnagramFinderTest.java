package com.github.tormoder.tmp.anagramjava.anagrams;

import junit.framework.TestCase;
import org.junit.Assert;

import java.io.FileInputStream;
import java.io.InputStreamReader;
import java.nio.charset.Charset;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;

public class AnagramFinderTest extends TestCase {

    public void testFind() throws Exception {

        ClassLoader classLoader = getClass().getClassLoader();

        // Load wanted output.
        String goldenPath = classLoader.getResource("golden.txt").getPath();
        byte[] goldenAsBytes = Files.readAllBytes(Paths.get(goldenPath));
        String goldenAsString = new String(goldenAsBytes, Charset.forName("UTF-8"));

        // Load input test data.
        String inputPath = classLoader.getResource("eventyr.txt").getPath();
        FileInputStream fileInputStream = new FileInputStream(inputPath);
        InputStreamReader inputStreamReader = new InputStreamReader(fileInputStream, Charset.forName("UTF-8"));

        // Process input.
        AnagramFinder anagramFinder = new AnagramFinder(inputStreamReader, AnagramFinder.SortMethod.COUNT);
        List<Anagrams> anagramsList = anagramFinder.find();
        StringBuilder sb = new StringBuilder();
        for (Anagrams ag: anagramsList) {
            sb.append(ag);
            sb.append('\n');
        }
        String result = sb.toString();

        // Compare.
        Assert.assertEquals(result, goldenAsString);
    }

}