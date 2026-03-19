package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	flag "github.com/spf13/pflag"
)

var BadInput = errors.New("range start cannot be greater than range end")

// cut
type Cut struct {
	lines     []string
	fields    string
	fieldsIdx []int
	delimiter string
	separated bool
}

func (c *Cut) prepareFields() error {
	set := make(map[int]struct{})
	a := strings.Split(c.fields, ",")
	fmt.Println(a)
	res := make([]int, 0, len(a))
	for _, v := range a {

		vFields := strings.Split(v, "-")

		switch {
		case len(vFields) > 1:
			start, end := vFields[0], vFields[len(vFields)-1]
			s, err := strconv.Atoi(start)
			e, err := strconv.Atoi(end)
			if err != nil {
				return err
			}
			if s > e {
				return BadInput
			}
			for i := s; i <= e; i++ {
				set[i] = struct{}{}
			}
		case len(vFields) == 1:
			val, err := strconv.Atoi(vFields[0])
			if err != nil {
				return err
			}
			set[val] = struct{}{}
		}

	}
	for k := range set {
		res = append(res, k)
	}
	slices.Sort(res)

	c.fieldsIdx = res
	return nil
}

func (c *Cut) cut(line string) {
	columns := strings.Split(line, c.delimiter)

	if c.separated && len(columns) == 1 {
		return
	}

	var res []string
	for _, idx := range c.fieldsIdx {
		if idx > 0 && idx < len(line) {
			res = append(res, columns[idx-1])
		}
	}
	if len(res) > 0 {
		fmt.Println(strings.Join(res, c.delimiter))
	}
}

func parseFlags() *Cut {
	c := &Cut{}

	flag.StringVarP(&c.fields, "fields", "f", "", "write after found")
	flag.StringVarP(&c.delimiter, "delimiter", "d", "\t", "delimiter")
	flag.BoolVarP(&c.separated, "separated", "s", false, "separated")

	flag.Parse()

	return c
}

func (c *Cut) Read(r io.Reader, handle func(line string)) error {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		handle(scanner.Text())
	}
	return scanner.Err()
}

func main() {
	c := parseFlags()
	c.prepareFields()
	input := os.Stdin

	_ = c.Read(input, c.cut)
}
