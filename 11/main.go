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
	input1 := []string{"Пятак", "пятка", "тяпка", "листок", "Слиток", "столик", "стол", "стул", "Кот", "Ток", "кТо", "ОТК", "ааб", "аба", "баа", "сон", "нос", "сон"}
	input2 := []string{"один", "два", "три", "ирт", "тир", "рит", "", " ", "а", "б", "в", "а"}
	input3 := []string{"ропот", "топор", "отпор", "прото", "ротор", "автор", "товар", "отвар", "рвота", "тавро"}
	input4 := []string{"яблоко", "груша", "банан", "киви", "мандарин", "слива"}
	fmt.Println(anagram(input1))
	fmt.Println(anagram(input2))
	fmt.Println(anagram(input3))
	fmt.Println(anagram(input4))
}
