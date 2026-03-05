package main

import (
	"cmp"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"slices"
	"strings"
)

type Sort struct {
	dataBlocks [][]string
}

func (s *Sort) Sort(input []byte, column int) {
	column -= 1
	if column < 0 {
		column = 0
	}

	str := strings.TrimSpace(string(input))

	stringsInput := strings.Split(str, "\n")

	for _, v := range stringsInput {
		fields := strings.Fields(v)
		if len(fields) > 0 {
			s.dataBlocks = append(s.dataBlocks, fields)
		}
	}
	slog.Debug("s.data", "[]string", s.dataBlocks)

	slices.SortStableFunc(s.dataBlocks, func(a, b []string) int {
		if column >= len(a) || column >= (len(b)) {
			return 0
		}

		return cmp.Compare(a[column], b[column])
	})
	for _, v := range s.dataBlocks {
		fmt.Println(strings.Join(v, " "))
	}
}

func main() {
	setupLogger(slog.LevelInfo)

	columnPtr := flag.Int("k", 1, "Column number to sort by (1-based)")

	flag.Parse()

	data, err := os.ReadFile("file.txt")
	if err != nil {
		slog.Error("cant read file", "error", err)
		os.Exit(1)
	}
	s := Sort{}
	s.Sort(data, *columnPtr)
}

func setupLogger(level slog.Level) {
	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(handler)

	slog.SetDefault(logger)
}
