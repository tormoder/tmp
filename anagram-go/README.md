# anagram-go

Installation:

```
go get github.com/tormoder/tmp/anagram-go
```

Usage:

```
usage: anagram-go [flags] [file]
  -parallel
        process input in parallel
  -sort string
        sort method to use: [count | lex | wordsig]
```

Read input from file:

```
anagram-go input.txt
```

Read from standard input:

```
anagram-go < input.txt
```
