package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	var flip bool
	flag.BoolVar(&flip, "f", false, "")
	flag.BoolVar(&flip, "flip", false, "")

	var separator string
	flag.StringVar(&separator, "s", "", "")
	flag.StringVar(&separator, "separator", "", "")

	flag.Parse()

	if flag.NArg() < 2 {
		flag.Usage()
		os.Exit(1)
	}

	var err error
	var prefixFile io.Reader = os.Stdin
	if f := flag.Arg(0); f != "-" {
		prefixFile, err = os.Open(f)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	var suffixFile io.Reader = os.Stdin
	if f := flag.Arg(1); f != "-" {
		suffixFile, err = os.Open(flag.Arg(1))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
	}

	// use 'a' and 'b' because which is the prefix
	// and which is the suffix depends on if we're in
	// flip mode or not.
	fileA := prefixFile
	fileB := suffixFile

	if flip {
		fileA, fileB = fileB, fileA
	}

	// we need to read through the lines from fileB multiple
	// times. If fileB is stdin then we can't seek back to
	// the beginning of the input. To get around this, read
	// fileB into a slice of strings instead.
	bs := make([]string, 0)
	sc := bufio.NewScanner(fileB)
	for sc.Scan() {
		bs = append(bs, sc.Text())
	}

	a := bufio.NewScanner(fileA)
	for a.Scan() {

		for _, b := range bs {
			if flip {
				fmt.Printf("%s%s%s\n", b, separator, a.Text())
			} else {
				fmt.Printf("%s%s%s\n", a.Text(), separator, b)
			}
		}
	}

}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Combine the lines from two files in every combination. Use '-' to read from stdin.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  comb [OPTIONS] [PREFIXFILE|-] [SUFFIXFILE|-]\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		fmt.Fprintf(os.Stderr, "  -f, --flip             Flip mode (order by suffix)\n")
		fmt.Fprintf(os.Stderr, "  -s, --separator <str>  String to place between prefix and suffix\n")
	}
}
