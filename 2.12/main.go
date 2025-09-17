package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
)

// Flags is ...
type Flags struct {
	after   *int
	before  *int
	context *int
	count   *bool
	ignore  *bool
	invert  *bool
	fixed   *bool
	lineNum *bool
}

func parseFlags() Flags {
	A := flag.Int("A", 0, "")
	B := flag.Int("B", 0, "")
	C := flag.Int("C", 0, "")
	c := flag.Bool("c", false, "")
	i := flag.Bool("i", false, "")
	v := flag.Bool("v", false, "")
	F := flag.Bool("F", false, "")
	n := flag.Bool("n", false, "")

	flag.Parse()

	if *C > 0 {
		*A, *B = *C, *C
	}

	flags := Flags{
		after: A, before: B, context: C, count: c,
		ignore: i, invert: v, fixed: F, lineNum: n,
	}

	return flags
}

func readLines(r io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(r)

	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	return lines, scanner.Err()
}

func compilePattern(flags Flags, pattern string) (*regexp.Regexp, error) {
	if *flags.fixed {
		if *flags.ignore {
			pattern = "(?i)" + regexp.QuoteMeta(pattern)
		} else {
			pattern = regexp.QuoteMeta(pattern)
		}
	} else {
		if *flags.ignore {
			pattern = "(?i)" + pattern
		}
	}
	return regexp.Compile(pattern)
}

func grep(flags Flags, lines []string, pattern string) error {
	re, err := compilePattern(flags, pattern)
	if err != nil {
		return err
	}

	var matches []int
	for i, line := range lines {
		found := re.MatchString(line)
		if *flags.invert {
			found = !found
		}
		if found {
			matches = append(matches, i)
		}
	}

	if *flags.count {
		fmt.Println(len(matches))
		return nil
	}

	toPrint := make(map[int]bool)
	for _, idx := range matches {
		start := idx - *flags.before
		if start < 0 {
			start = 0
		}
		end := idx + *flags.after
		if end >= len(lines) {
			end = len(lines) - 1
		}
		for i := start; i <= end; i++ {
			toPrint[i] = true
		}
	}

	for i := range lines {
		if toPrint[i] {
			if *flags.lineNum {
				fmt.Printf("%d:%s\n", i+1, lines[i])
			} else {
				fmt.Println(lines[i])
			}
		}
	}

	return nil
}

func main() {
	f := parseFlags()
	var r io.Reader = os.Stdin

	args := flag.Args()
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: grep [OPTIONS] PATTERN [FILE]")
		os.Exit(1)
	}
	pattern := args[0]

	if len(args) > 1 {
		f, err := os.Open(flag.Arg(1))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}

	lines, err := readLines(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	err = grep(f, lines, pattern)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
