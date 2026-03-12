package main

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"strings"

	flag "github.com/spf13/pflag"
)

type Grep struct {
	oldLines    []string
	filterLines []string
	after       int  //-A N
	before      int  // -B N
	context     int  // -C N
	count       bool // -c
	ignoreCase  bool // -i
	invert      bool // -v
	fixed       bool // -F
	number      bool // -n
}

func (g *Grep) Read(r io.Reader) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		line := scanner.Text()
		g.oldLines = append(g.oldLines, line)
	}
	return scanner.Err()
}

func (g *Grep) filter(pattern string) func(line string) bool {

	if g.fixed {
		if g.ignoreCase {
			pattern = strings.ToLower(pattern)
		}

		return func(line string) bool {

			compareLine := line
			if g.ignoreCase {
				compareLine = strings.ToLower(line)
			}
			lineContains := strings.Contains(compareLine, pattern)

			return lineContains != g.invert
		}
	}

	regexPattern := pattern
	if g.ignoreCase {
		regexPattern = "(?i)" + regexPattern
	}
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		slog.Error("invalid regular expression", "error", err)
		os.Exit(2)
	}
	return func(line string) bool {
		match := re.MatchString(line)
		return match != g.invert
	}
}

func (g *Grep) findIdx(f func(line string) bool) []int {
	var idxSlice []int
	for i, line := range g.oldLines {
		if f(line) {
			idxSlice = append(idxSlice, i)
		}
	}
	return idxSlice

}

func parseFlags() *Grep {
	g := &Grep{}

	flag.IntVarP(&g.after, "after", "A", 0, "write after found")
	flag.IntVarP(&g.before, "before", "B", 0, "write columns before found")
	flag.IntVarP(&g.context, "context", "C", 0, "before and after found")
	flag.BoolVarP(&g.count, "count", "c", false, "count")
	flag.BoolVarP(&g.ignoreCase, "ignoreCase", "i", false, "ignore case")
	flag.BoolVarP(&g.invert, "invert", "v", false, "invert")
	flag.BoolVarP(&g.fixed, "fixed", "F", false, "fixed")
	flag.BoolVarP(&g.number, "number", "n", false, "number")

	flag.Parse()

	if g.context > 0 {
		g.after = g.context
		g.before = g.context
	}

	return g
}

func main() {
	g := parseFlags()
	args := flag.Args()
	var input io.Reader

	if len(args) == 0 {
		fmt.Println("Usage:[FLAGS] PATTERN [FILE]")
		return
	}

	pattern := args[0]

	if len(args) == 1 {
		input = os.Stdin
	} else {
		file, err := os.Open(args[1])
		if err != nil {
			slog.Error("cant open file", "error", err)
			os.Exit(1)
		}
		defer func() {
			_ = file.Close()
		}()
		input = file
	}

	if err := g.Read(input); err != nil {
		slog.Error("read error", "error", err)
		os.Exit(1)
	}

	filter := g.filter(pattern)
	foundIdx := g.findIdx(filter)

	for _, idx := range foundIdx {
		line := g.oldLines[idx]
		if g.number {
			line = fmt.Sprintf("%d: %s", idx+1, line)
		}
		g.filterLines = append(g.filterLines, line)
	}

	if g.count {
		fmt.Println(len(g.filterLines))
		return
	}

	for i := range g.filterLines {
		fmt.Println(g.filterLines[i])
	}
}
