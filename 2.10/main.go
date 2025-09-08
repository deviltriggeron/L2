package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

// Flags is ...
type Flags struct {
	k *int
	n *bool
	r *bool
	u *bool
	m *bool
	b *bool
	c *bool
	h *bool
}

func parseFlags() Flags {
	k := flag.Int("k", 1, "N — sort by column #N (separator — tab by default)")
	n := flag.Bool("n", false, "sort by numeric value (strings are interpreted as numbers).")
	r := flag.Bool("r", false, "sort in reverse order.")
	u := flag.Bool("u", false, "do not output duplicate lines (only unique ones).")
	m := flag.Bool("M", false, "sort by month name (Jan, Feb, ... Dec), i.e. recognize a special date format.")
	b := flag.Bool("b", false, "ignore trailing blanks.")
	c := flag.Bool("c", false, "check if the data is sorted; if not, display a message about it.")
	h := flag.Bool("h", false, "sort by multiple result taking into account suffixes (for example, K = kilobyte, M = megabyte - human-readable sizes).")
	flag.Parse()

	return Flags{
		k: k, n: n, r: r, u: u, m: m, b: b, c: c, h: h,
	}
}

func getKey(line string, column *int) string {
	fields := strings.Fields(line)
	if *column > 0 && *column <= len(fields) {
		return fields[*column-1]
	}
	return ""
}

func readLines(r io.Reader, flags Flags) ([]string, error) {
	scanner := bufio.NewScanner(r)

	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if *flags.b {
			line = strings.TrimRight(line, " \t")
		}
		lines = append(lines, line)
	}

	return lines, scanner.Err()
}

func sortLines(flags Flags) func(i, j string) bool {
	return func(i, j string) bool {
		ki := getKey(i, flags.k)
		kj := getKey(j, flags.k)

		var less bool
		if *flags.h {
			vi, _ := humanReadSize(ki)
			vj, _ := humanReadSize(kj)
			less = vi < vj
		} else if *flags.n {
			ni, errI := strconv.ParseFloat(ki, 64)
			nj, errJ := strconv.ParseFloat(kj, 64)
			if errI == nil && errJ == nil {
				less = ni < nj
			} else {
				less = ki < kj
			}
		} else if *flags.m {
			less = parseMonth(ki) < parseMonth(kj)
		} else {
			less = ki < kj
		}

		if *flags.r {
			return !less
		}

		return less
	}
}

func deleteDuplicate(lines []string) []string {
	res := []string{}

	for i := range lines {
		if !slices.Contains(res, lines[i]) {
			res = append(res, lines[i])
		}
	}

	return res
}

func parseMonth(s string) int {
	months := map[string]int{
		"Jan": 1, "Feb": 2, "Mar": 3, "Apr": 4,
		"May": 5, "Jun": 6, "Jul": 7, "Aug": 8,
		"Sep": 9, "Oct": 10, "Nov": 11, "Dec": 12,
	}
	return months[s]
}

func checkSort(s []string, flags Flags) {
	less := sortLines(flags)
	if sort.SliceIsSorted(s, func(i, j int) bool {
		return less(s[i], s[j])
	}) {
		fmt.Println("Slice sorted")
		os.Exit(0)
	}
	fmt.Println("Slice not sorted")
	os.Exit(0)
}

func humanReadSize(s string) (int64, error) {
	multipliers := map[byte]int64{
		'K': 1024,
		'M': 1024 * 1024,
		'G': 1024 * 1024 * 1024,
	}

	if len(s) == 0 {
		return 0, fmt.Errorf("empty size string")
	}

	last := s[len(s)-1]
	if mult, ok := multipliers[last]; ok {
		val, err := strconv.ParseInt(s[:len(s)-1], 10, 64)
		if err != nil {
			return 0, err
		}
		return val * mult, nil
	}

	return strconv.ParseInt(s, 10, 64)
}

func main() {
	flags := parseFlags()
	var r io.Reader = os.Stdin

	if flag.NArg() > 0 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	}

	lines, err := readLines(r, flags)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if *flags.c {
		checkSort(lines, flags)
	}

	less := sortLines(flags)
	sort.SliceStable(lines, func(i, j int) bool {
		return less(lines[i], lines[j])
	})

	if *flags.u {
		lines = deleteDuplicate(lines)
	}

	for _, l := range lines {
		fmt.Println(l)
	}
}
