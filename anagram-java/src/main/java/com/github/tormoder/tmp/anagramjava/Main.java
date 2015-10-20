package com.github.tormoder.tmp.anagramjava;

import com.github.tormoder.tmp.anagramjava.anagrams.AnagramFinder;
import com.github.tormoder.tmp.anagramjava.anagrams.Anagrams;
import org.kohsuke.args4j.Argument;
import org.kohsuke.args4j.CmdLineException;
import org.kohsuke.args4j.CmdLineParser;
import org.kohsuke.args4j.Option;

import java.io.FileInputStream;
import java.io.IOException;
import java.io.InputStreamReader;
import java.nio.charset.Charset;
import java.util.ArrayList;
import java.util.List;

public class Main {

    @Option(name="-sort",usage="sort method to use: [count | lex | wordsig]")
    private String sortMethodString = "count";

    @Argument
    private List<String> arguments = new ArrayList<>();

    public static void main(String[] args) {
        new Main().runMain(args);
    }

    public void runMain(String[] args) {

        CmdLineParser parser = new CmdLineParser(this);

        try {
            parser.parseArgument(args);
        } catch(CmdLineException e) {
            System.err.println(e.getMessage());
            System.err.println("usage: anagram-java [flags] [file]");
            parser.printUsage(System.err);
            System.exit(2);
        }

        AnagramFinder.SortMethod sortMethod = AnagramFinder.SortMethod.COUNT;

        switch (sortMethodString) {
            case "count":
                sortMethod = AnagramFinder.SortMethod.COUNT;
                break;
            case "lex":
                sortMethod = AnagramFinder.SortMethod.LEXICOGRAPHICAL;
                break;
            case "wordsig":
                sortMethod = AnagramFinder.SortMethod.WORDSIGNATURE;
                break;
            default:
                System.err.printf("Unknown sort option: %s\n\n", sortMethodString);
                System.err.println("usage: anagram-java [flags] [file]");
                parser.printUsage(System.err);
                System.exit(2);
        }

        InputStreamReader inputStreamReader = null;
        FileInputStream fileInputStream = null;

        try {

            if (arguments.isEmpty()) {
                inputStreamReader = new InputStreamReader(System.in, Charset.forName("UTF-8"));
            } else {
                fileInputStream = new FileInputStream(arguments.get(0));
                inputStreamReader = new InputStreamReader(fileInputStream, Charset.forName("UTF-8"));
            }

            AnagramFinder anagramFinder = new AnagramFinder(inputStreamReader, sortMethod);
            List<Anagrams> anagramsList = anagramFinder.find();

            for (Anagrams ags : anagramsList) {
                System.out.println(ags);
            }

        } catch (IOException e) {
            System.err.println("Error parsing input:");
            System.err.println(e.toString());
        } finally {
            try {
                if (fileInputStream != null) {
                    fileInputStream.close();
                }
                if (inputStreamReader != null) {
                    inputStreamReader.close();
                }
            } catch (IOException e){
                e.printStackTrace();
            }
            System.exit(2);
        }
    }
}
