package cmd

import (
	"unicode"
)

func FindRange(r rune) (rangename string) {
	for i, name := range unicode.Scripts {
		if unicode.Is(name, r) {
			rangename = i
		}
	}
	return
}
