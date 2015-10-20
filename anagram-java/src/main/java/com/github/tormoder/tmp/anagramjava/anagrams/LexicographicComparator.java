package com.github.tormoder.tmp.anagramjava.anagrams;

import java.util.Comparator;

public class LexicographicComparator implements Comparator<Anagrams> {

    @Override
    public int compare(Anagrams a1, Anagrams a2) {
        if (a1.count() > 0 && a2.count() > 0) {
            return a1.getWords().get(0).compareTo(a2.getWords().get(0));
        }
        return 0;
    }

}
