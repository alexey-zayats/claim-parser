package util

import (
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"
)

// RunIndex ...
func RunIndex(s string, edge int) int {
	step := 0
	idx := 0
	for i, w := 0, 0; i < len(s); i += w {
		step++
		_, w = utf8.DecodeRuneInString(s[i:])
		if step == edge {
			idx = i + w
			break
		}
	}
	return idx
}

// TrimNumber ...
func TrimNumber(car string) string {
	re := regexp.MustCompile(`\s+`)

	car = re.ReplaceAllString(car, "")
	car = strings.ToUpper(car)

	re1 := regexp.MustCompile(`(?i:rus?)$`)
	re2 := regexp.MustCompile(`(?i:рус?)$`)

	car = re1.ReplaceAllString(car, "")
	car = re2.ReplaceAllString(car, "")

	return car
}

// TrimSpaces ...
func TrimSpaces(str string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, str)
}
