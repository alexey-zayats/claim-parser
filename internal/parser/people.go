package parser

import (
	"github.com/alexey-zayats/claim-parser/internal/model"
	"regexp"
	"strings"
)

// ParsePeoples ...
func ParsePeoples(line string) ([]model.PeoplePass, bool) {
	line = strings.ReplaceAll(line, "–", "")
	line = strings.ReplaceAll(line, "—", "")
	line = strings.ReplaceAll(line, "-", "")
	line = regexp.MustCompile(`[()–,.\r\n\t;]`).ReplaceAllString(line, " ")

	line = regexp.MustCompile(`\d\.`).ReplaceAllString(line, " ")
	line = regexp.MustCompile(`^\d\s`).ReplaceAllString(line, " ")
	line = regexp.MustCompile(`\s\d\s`).ReplaceAllString(line, " ")

	lines := regexp.MustCompile(`\r?\n`).Split(line, -1)

	success := true
	var fio model.FIO

	passes := make([]model.PeoplePass, 0)
	for _, line := range lines {
		if fio, success = ParseFIO(line); success {
			passes = append(passes, model.PeoplePass{FIO: fio})
		}
	}

	return passes, success
}
