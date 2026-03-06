package main

import (
	"bufio"
	"cmp"
	"fmt"
	"io"
	"log/slog"
	"os"
	"slices"
	"strconv"
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

func (s *Sorter) parseFloat(v string) float64 {
	if v == "" {
		return 0
	}
	val, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0
	}
	return val
}

func (s *Sorter) Read(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		s.Lines = append(s.Lines, line)
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}
	return nil
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

	s.Lines = append(s.Lines[:0], result...)
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

		if s.Num {
			aFloat := s.parseFloat(aVal)
			bFloat := s.parseFloat(bVal)

			if s.Reverse {
				return cmp.Compare(bFloat, aFloat)
			}

			return cmp.Compare(aFloat, bFloat)

		}

		if s.Reverse {
			return cmp.Compare(bVal, aVal)
		}

		return cmp.Compare(aVal, bVal)
	})
}

func (s *Sorter) Sort() {
	s.Column -= 1
	if s.Column < 0 {
		s.Column = 0
	}

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
	numPtr := flag.BoolP("numeric", "n", false, "Compare according to string numerical value")

	flag.Parse()

	var input io.Reader
	args := flag.Args()

	if len(args) == 0 {
		input = os.Stdin
	} else {
		file, err := os.Open(args[0])
		if err != nil {
			slog.Error("cant open file", "error", err)
			os.Exit(1)
		}
		defer func() {
			_ = file.Close()
		}()
		input = file
	}

	s := Sorter{
		Column:  *columnPtr,
		Reverse: *reversePtr,
		Uniq:    *uniquePtr,
		Num:     *numPtr,
	}

	if err := s.Read(input); err != nil {
		slog.Error("read error", "error", err)
		os.Exit(1)
	}

	s.Sort()
}

func setupLogger(level slog.Level) {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)
}
