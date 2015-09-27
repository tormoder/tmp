package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: fuzztotest [gofuzz workdir]\n")
	fmt.Fprintf(os.Stderr, "\noptions:\n")
	flag.PrintDefaults()
	os.Exit(2)
}

var (
	tnsuffix = flag.String("tnsuffix", "FuzzCrashers", "test method name suffix")
	pkg      = flag.String("pkg", "main", "package name")
)

func isQuoted(fi os.FileInfo) bool {
	return strings.HasSuffix(fi.Name(), ".quoted")
}

type gen struct {
	bytes.Buffer
}

func (g *gen) n() {
	g.WriteByte('\n')
}

func (g *gen) p(ss ...string) {
	for _, s := range ss {
		g.WriteString(s)
	}
	g.WriteByte('\n')
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
	}

	bpath := filepath.Join(flag.Arg(0), "crashers")
	files, err := ioutil.ReadDir(bpath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading directory: %v\n", err)
		os.Exit(2)
	}

	var qfiles []string
	for _, file := range files {
		if isQuoted(file) {
			qfiles = append(qfiles, file.Name())
		}
	}

	if len(qfiles) == 0 {
		fmt.Fprintf(os.Stderr, "found no '.quoted' files\n", err)
		os.Exit(2)
	}

	g := new(gen)
	g.writeHeader()
	g.p("var inputs = [...]string{")
	for i, file := range qfiles {
		filep := filepath.Join(bpath, file)
		input, err := ioutil.ReadFile(filep)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file %q: %v\n", file, err)
			os.Exit(2)
		}
		g.p("// ", fmt.Sprintf("%d", i), ".")
		input = bytes.TrimSuffix(input, []byte{'\n'})
		g.Write(input)
		g.p(",")
		g.n()
	}
	g.p("}")
	g.n()
	g.writeTest()

	fmted, err := format.Source(g.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "error formatting code: %v\n\n", err)
		fmt.Fprintf(os.Stderr, "%s", g.Bytes())
		os.Exit(2)
	}

	fmt.Fprintf(os.Stdout, "%s", string(fmted))
}

func (g *gen) writeHeader() {
	g.p("package ", *pkg)
	g.n()
	g.p("import (")
	g.p("\"bytes\"")
	g.p("\"testing\"")
	g.p(")")
	g.n()
}

func (g *gen) writeTest() {
	g.p("func Test", *tnsuffix, "(t *testing.T) {")
	g.p("for i, input := range inputs {")
	g.p("t.Logf(\"Crasher input: %d\", i)")
	g.p("Decode(bytes.NewReader([]byte(input)))")
	g.p("}")
	g.p("}")
}
