package com.github.tormoder.tmp.anagramjava.anagrams;

import java.util.Comparator;

public class WordSignatureComparator implements Comparator<Anagrams> {

    @Override
    public int compare(Anagrams a1, Anagrams a2) {
        return a1.getWordSignature().compareTo(a2.getWordSignature());
    }

}
