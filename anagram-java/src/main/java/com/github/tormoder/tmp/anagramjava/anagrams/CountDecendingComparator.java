package com.github.tormoder.tmp.anagramjava.anagrams;

import java.util.Comparator;

public class CountDecendingComparator implements Comparator<Anagrams> {

    @Override
    public int compare(Anagrams a1, Anagrams a2) {
        if (a1.count() < a2.count()) {
            return 1;
        } else if (a1.count() > a2.count()) {
           return -1;
        } else {
            return new LexicographicComparator().compare(a1, a2);
        }
    }

}
