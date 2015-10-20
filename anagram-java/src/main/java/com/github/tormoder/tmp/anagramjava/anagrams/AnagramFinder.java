package com.github.tormoder.tmp.anagramjava.anagrams;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStreamReader;
import java.util.*;

public class AnagramFinder {

    public enum SortMethod {
        COUNT(new CountDecendingComparator()), LEXICOGRAPHICAL(new LexicographicComparator()), WORDSIGNATURE(new WordSignatureComparator());

        private final Comparator<Anagrams> comparator;

        SortMethod(Comparator<Anagrams> comparator) {
            this.comparator = comparator;
        }
    }

    private final InputStreamReader inputStreamReader;
    private final SortMethod sortMethod;
    private final Map<String, Anagrams> anagramsMap;

    public AnagramFinder(InputStreamReader inputStreamReader, SortMethod sortMethod) {
        this.inputStreamReader = inputStreamReader;
        this.sortMethod = sortMethod;
        this.anagramsMap = new HashMap<>();
    }

    public List<Anagrams> find() throws IOException {

        String inputLine;
        String wordSignature;
        char[] wordChars;
        Anagrams anagrams;

        try (BufferedReader br = new BufferedReader(this.inputStreamReader)) {
             while ((inputLine = br.readLine()) != null) {
                wordChars = inputLine.toCharArray();
                Arrays.sort(wordChars);
                wordSignature = new String(wordChars);
                anagrams = anagramsMap.get(wordSignature);
                if (anagrams == null) {
                    anagrams = new Anagrams(wordSignature);
                    anagramsMap.put(wordSignature, anagrams);
                }
                anagrams.addWord(inputLine);
            }

            List<Anagrams> anagramsList = new ArrayList<>();
            for (Anagrams ags: anagramsMap.values()) {
                if (ags.count() > 1) {
                    anagramsList.add(ags);
                }
            }

            Collections.sort(anagramsList, this.sortMethod.comparator);

            return anagramsList;
        }
    }

}
