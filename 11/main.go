package main

import (
	"fmt"
	"slices"
	"strings"
)

func anagram(s []string) map[string][]string {
	m := make(map[string][]string, len(s))

	for _, v := range s {
		vL := strings.ToLower(v)
		r := []rune(vL)
		slices.Sort(r)
		m[string(r)] = append(m[string(r)], vL)
	}
	res := make(map[string][]string, len(m))
	for _, v := range m {
		if len(v) > 1 {
			first := v[0]
			slices.Sort(v)
			v = slices.Compact(v)
			res[first] = v
		}
	}

	return res
}

func main() {
	input := []string{"Пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол"}
	fmt.Println(anagram(input))
}
