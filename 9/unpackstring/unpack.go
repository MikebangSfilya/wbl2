package unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var errDigit = errors.New("ERROR: first argument is digit")

func Unpack(s string) (string, error) {
	sb := strings.Builder{}
	sb.Grow(len(s))
	var prev rune
	for _, v := range s {

		if !unicode.IsDigit(v) {
			prev = v
			sb.WriteRune(prev)
		} else {
			if prev == 0 {
				return "", errDigit
			}
			n := int(v-'0') - 1
			for range n {
				sb.WriteRune((prev))
			}

		}
	}
	return sb.String(), nil
}
