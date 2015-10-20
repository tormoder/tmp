package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/tormoder/tmp/anagram-go/anagram"
)

func main() {
	var sortMethod = flag.String("sort", "", "sort method to use: [count | lex | wordsig]")
	var parallel = flag.Bool("parallel", false, "process input in parallel")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: anagram-go [flags] [file]\n")
		flag.PrintDefaults()
	}

	flag.Parse()
	if flag.NArg() > 1 {
		flag.Usage()
		os.Exit(2)
	}

	var input io.Reader
	if flag.NArg() == 0 {
		input = os.Stdin
	} else {
		data, err := ioutil.ReadFile(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(2)
		}
		input = bytes.NewBuffer(data)
	}

	switch *sortMethod {
	case "", "count", "lex", "wordsig":
		var (
			result string
			err    error
		)
		if *parallel {
			result, err = anagram.FindParallel(input, *sortMethod)
		} else {
			result, err = anagram.Find(input, *sortMethod)
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v", err)
			os.Exit(2)
		}
		fmt.Println(result)
	default:
		fmt.Fprintf(os.Stderr, "unknown sort option: %q\n\n", *sortMethod)
		flag.Usage()
		os.Exit(2)
	}
}
