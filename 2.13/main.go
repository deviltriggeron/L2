package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Flags is...
type Flags struct {
	fields    *string
	delimiter *string
	separated *bool
}

func parseFlags() Flags {
	f := flag.String("f", "", "specifying the numbers of fields (columns) to be output.")
	d := flag.String("d", "\t", "use a different separator (symbol). The default separator is tab ('\t').")
	s := flag.Bool("s", false, "only lines containing the separator. If the flag is specified, lines without the separator are ignored (not output).")

	flag.Parse()

	flags := Flags{
		fields:    f,
		delimiter: d,
		separated: s,
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

func parseFields(fields string) (map[int]bool, error) {
	res := make(map[int]bool)
	parts := strings.Split(fields, ",")
	for _, p := range parts {

		if strings.Contains(p, "-") {
			bounds := strings.SplitN(p, "-", 2)
			start, err1 := strconv.Atoi(strings.TrimSpace(bounds[0]))
			end, err2 := strconv.Atoi(strings.TrimSpace(bounds[1]))

			if err1 != nil || err2 != nil || start < 1 || end < 1 || start > end {
				return nil, fmt.Errorf("invalid range: %q", p)
			}

			for i := start; i <= end; i++ {
				res[i-1] = true
			}

		} else {
			val, err := strconv.Atoi(strings.TrimSpace(p))
			if err != nil {
				return nil, fmt.Errorf("invalid field: %q", p)
			}
			res[val-1] = true
		}
	}

	return res, nil
}

func parseDelimiter(delim string) error {
	if delim == "" || len(delim) != 1 {
		return fmt.Errorf("bad delimiter")
	}
	return nil
}

func cut(lines []string, flags Flags) {
	fields, err := parseFields(*flags.fields)
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse fields error:", err)
		os.Exit(1)
	}

	err = parseDelimiter(*flags.delimiter)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	for _, line := range lines {
		if *flags.separated && !strings.Contains(line, *flags.delimiter) {
			continue
		}

		if !strings.Contains(line, *flags.delimiter) {
			fmt.Println(line)
			continue
		}

		var res []string
		cols := strings.Split(line, *flags.delimiter)
		for index := range cols {
			if fields[index] {
				res = append(res, cols[index])
			}
		}

		if len(res) > 0 {
			fmt.Println(strings.Join(res, *flags.delimiter))
		}
	}
}

func main() {
	flags := parseFlags()
	var r io.Reader = os.Stdin

	if *flags.fields == "" {
		fmt.Fprintln(os.Stderr, "usage: go run main.go -f list [-s] [-w | -d delim] [file ...]")
		os.Exit(1)
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer f.Close()
	r = f

	lines, err := readLines(r)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	cut(lines, flags)
}
