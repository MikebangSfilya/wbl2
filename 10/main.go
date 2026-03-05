package main

import (
	"cmp"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"

	flag "github.com/spf13/pflag"
)

type Sorter struct {
	Lines   []string
	set     map[string]struct{}
	Column  int
	Reverse bool
	Num     bool
	Uniq    bool
}

func (s *Sorter) prepare(input []byte) {
	str := strings.TrimSpace(string(input))
	stringsInput := strings.Split(str, "\n")
	s.Lines = stringsInput
	slog.Debug("s.data", "[]string", stringsInput)
}

func (s *Sorter) unique() {
	s.set = make(map[string]struct{}, len(s.Lines))
	var result []string
	for _, v := range s.Lines {
		if _, ok := s.set[v]; !ok {
			s.set[v] = struct{}{}
			result = append(result, v)
		}
	}

	s.Lines = s.Lines[:0]
	for _, v := range result {
		s.Lines = append(s.Lines, v)
	}
}

func (s *Sorter) sort() {
	slices.SortStableFunc(s.Lines, func(a, b string) int {
		aField := strings.Fields(a)
		bField := strings.Fields(b)

		var aVal, bVal string
		if s.Column < len(aField) {
			aVal = aField[s.Column]
		}
		if s.Column < len(bField) {
			bVal = bField[s.Column]
		}

		if s.Reverse {
			return cmp.Compare(bVal, aVal)
		}

		return cmp.Compare(aVal, bVal)
	})
}

func (s *Sorter) Sort(input []byte) {
	s.Column -= 1
	if s.Column < 0 {
		s.Column = 0
	}

	s.prepare(input)
	if s.Uniq {
		s.unique()
	}
	s.sort()

	for _, v := range s.Lines {
		fmt.Println(v)
	}

}

func main() {
	setupLogger(slog.LevelDebug)

	columnPtr := flag.IntP("column", "k", 1, "Column number to sort by")
	reversePtr := flag.BoolP("reverse", "r", false, "Reverse sort")
	uniquePtr := flag.BoolP("unique", "u", false, "Unique lines only")

	flag.Parse()

	data, err := os.ReadFile("file.txt")
	if err != nil {
		slog.Error("cant read file", "error", err)
		os.Exit(1)
	}
	s := Sorter{
		Column:  *columnPtr,
		Reverse: *reversePtr,
		Uniq:    *uniquePtr,
	}

	s.Sort(data)
}

func setupLogger(level slog.Level) {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)
}
