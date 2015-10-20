package anagram_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/tormoder/tmp/anagram-go/anagram"
)

const testDataPath = "../testdata/eventyr.txt"

func TestSortCount(t *testing.T) {
	tdata := readTestData(testDataPath, t)
	result, err := anagram.Find(tdata, "count")
	if err != nil {
		t.Fatal(err)
	}
	if result != sortCountGolden {
		errorPrintStrings(result, sortCountGolden, t)
	}
}

func TestSortLexicographical(t *testing.T) {
	tdata := readTestData(testDataPath, t)
	result, err := anagram.Find(tdata, "lex")
	if err != nil {
		t.Fatal(err)
	}
	if result != sortLexGolden {
		errorPrintStrings(result, sortLexGolden, t)
	}
}

func TestSortWordSignature(t *testing.T) {
	tdata := readTestData(testDataPath, t)
	result, err := anagram.Find(tdata, "wordsig")
	if err != nil {
		t.Fatal(err)
	}
	if result != sortWordSigGolden {
		errorPrintStrings(result, sortWordSigGolden, t)
	}
}

func readTestData(path string, t *testing.T) *bytes.Buffer {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("error reading test data: %v", err)
	}
	return bytes.NewBuffer(b)
}

func errorPrintStrings(got, want string, t *testing.T) {
	t.Errorf("got != want\n")
	t.Errorf("got:\n")
	t.Error(got)
	t.Errorf("want:\n")
	t.Error(want)
}

const sortCountGolden = `dro ord rod
at ta
bar bra
bry byr
dem med
den ned
denne enden
dra rad
ende nede
engang gangen
ens sen
etter rette
glinset glinste
hellestein steinhelle
kisten skinte
kristent kristnet
krok rokk
lovt tolv
lysnet lysten
løst støl
mor rom
navn vann
niste stien
ordet torde
ristet sitter
rå år
stuen suten
søsteren søstrene
truet turte
`

const sortLexGolden = `at ta
bar bra
bry byr
dem med
den ned
denne enden
dra rad
dro ord rod
ende nede
engang gangen
ens sen
etter rette
glinset glinste
hellestein steinhelle
kisten skinte
kristent kristnet
krok rokk
lovt tolv
lysnet lysten
løst støl
mor rom
navn vann
niste stien
ordet torde
ristet sitter
rå år
stuen suten
søsteren søstrene
truet turte
`

const sortWordSigGolden = `bar bra
dra rad
engang gangen
navn vann
at ta
bry byr
ende nede
denne enden
dem med
den ned
ordet torde
dro ord rod
hellestein steinhelle
søsteren søstrene
etter rette
glinset glinste
kristent kristnet
kisten skinte
niste stien
ristet sitter
lysnet lysten
ens sen
stuen suten
truet turte
krok rokk
lovt tolv
løst støl
mor rom
rå år
`
