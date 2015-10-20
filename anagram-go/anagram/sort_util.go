package anagram

import "sort"

type runeSlice []rune

func (p runeSlice) Len() int           { return len(p) }
func (p runeSlice) Less(i, j int) bool { return p[i] < p[j] }
func (p runeSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type lessFunc func(a1, a2 anagrams) bool

type anagramSorter struct {
	anagrams []anagrams
	less     []lessFunc
}

func (as *anagramSorter) sort(anagrams []anagrams) {
	as.anagrams = anagrams
	sort.Sort(as)
}

func orderedBy(less ...lessFunc) *anagramSorter {
	as := &anagramSorter{
		less: less,
	}
	return as
}

func (as *anagramSorter) Len() int {
	return len(as.anagrams)
}

func (as *anagramSorter) Swap(i, j int) {
	as.anagrams[i], as.anagrams[j] = as.anagrams[j], as.anagrams[i]
}

// See "SortMultiKeys" example from https://golang.org/pkg/sort/
func (as *anagramSorter) Less(i, j int) bool {
	p, q := as.anagrams[i], as.anagrams[j]
	var k int
	for k = 0; k < len(as.less)-1; k++ {
		less := as.less[k]
		switch {
		case less(p, q):
			return true
		case less(q, p):
			return false
		}
	}
	return as.less[k](p, q)
}

var nrOfAnagrams = func(a1, a2 anagrams) bool {
	return len(a1.words) > len(a2.words) // Note: > for decending.
}

var lexicoOrder = func(a1, a2 anagrams) bool {
	return a1.words[0] < a2.words[0]
}

var wordSignature = func(a1, a2 anagrams) bool {
	return a1.anagramBase < a2.anagramBase
}
