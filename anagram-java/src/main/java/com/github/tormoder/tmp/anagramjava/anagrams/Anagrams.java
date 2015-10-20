package com.github.tormoder.tmp.anagramjava.anagrams;

import org.apache.commons.lang3.StringUtils;

import java.util.ArrayList;
import java.util.List;

public class Anagrams {
        final String wordSignature;
        final List<String> words;

        Anagrams(String wordSignature) {
            this.wordSignature = wordSignature;
            this.words = new ArrayList<>();
        }

        List<String> getWords() {
            return words;
        }

        String getWordSignature() {
            return wordSignature;
        }

        void addWord(String word) {
            this.words.add(word);
        }

        int count() {
            if (this.getWords() == null) {
                return 0;
            }

            return this.getWords().size();
        }

        @Override
        public String toString() {
            return StringUtils.join(this.words, " ");
        }
    }
