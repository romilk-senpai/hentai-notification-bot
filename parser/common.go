package parser

import (
	"strconv"
	"strings"
	"unicode"
)

func ParseNumeric(input string) (int, error) {
	cleaned := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, input)

	output, err := strconv.Atoi(cleaned)

	if err != nil {
		return 0, err
	}

	return output, nil
}
