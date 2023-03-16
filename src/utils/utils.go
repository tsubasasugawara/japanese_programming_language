package utils

import (
	"strings"
	"unicode"
)

var numConv = unicode.SpecialCase{
	unicode.CaseRange{
		Lo: 0xff10,
		Hi: 0xff19,
		Delta: [unicode.MaxCase]rune{
			0,
			0x0030 - 0xff10,
			0,
		},
	},
}

func ToLower(str string) string {
	return strings.ToLowerSpecial(numConv, str)
}
