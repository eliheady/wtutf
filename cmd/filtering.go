package cmd

import (
	"fmt"
	"io"
	"unicode"
)

// listRanges takes a string and returns a map of Unicode range
// names and the count of runes within that range
func listRanges(ustring string) map[string]int {
	rangeCounts := map[string]int{}

	for _, r := range ustring {
		rangeCounts[FindRange(r)]++
	}
	return rangeCounts
}

// checkMultipleRange takes a string and determines whether the string
// contains characters from more than one range.
// Runes from the 'Common' range are ignored.
func checkMultipleRange(w io.Writer, showRanges bool, ustring string) (multiRange bool) {
	var ranges int
	var out string
	for i, count := range listRanges(ustring) {
		if i == "Common" {
			continue
		}
		ranges++
		if showRanges {
			out += fmt.Sprintf("%s: %d\n", i, count)
		}
	}
	if ranges > 1 {
		multiRange = true
		if w != nil && out != "" {
			fmt.Fprint(w, out)
		}
	}
	return
}

func FindRange(r rune) (rangename string) {
	for i, name := range unicode.Scripts {
		if unicode.Is(name, r) {
			rangename = i
		}
	}
	return
}

// politePrint takes a rune and outputs a string safe to print in the terminal,
// and where possible with the original character. Control and formatting
// characters are filtered, some combining characters are printed.
func politePrint(r rune) string {
	switch {
	// render combining diacritics with one or two ◌ dotted circle characters
	// to keep them from being printed over the colon
	case unicode.Is(unicode.Diacritic, r):
		if r == 0x005e || r == 0x0060 {
			// ^, ` do not combine
			return " " + string(r)
		}
		if 0x035C <= r && r <= 0x0362 {
			// combining diacritical marks, 2 characters
			return "◌" + string(r) + "◌"
		}
		return " ◌" + string(r)

	// filter control characters
	// reasoning: terminal control sequences can do all sorts of damage to the
	// output.
	case unicode.IsControl(r):
		return "^?"

	// Other Alphabetic contains a mix of LTR and RTL and combining characters.
	// todo: if we implement locale detection we could use that to properly
	// handle right-to-left runes
	//case unicode.Is(unicode.Other_Alphabetic, r):
	//	return string(r) + " "

	// filter direction changing characters
	// reasoning: we are printing a single character here. If we needed to print
	// words then the directional marks should not be discarded. These checks are
	// an attempt to prevent leaving an unclosed direction change in the output.
	// ref. https://www.unicode.org/reports/tr9/#Directional_Formatting_Codes
	case unicode.Is(unicode.C, r):
		return "^?"
	// filter formatting characters
	case unicode.Is(unicode.Mn, r) || unicode.Is(unicode.Mc, r):
		return "^?"
	case unicode.Is(unicode.Bidi_Control, r):
		return "^?"
	case unicode.Is(unicode.Join_Control, r):
		return "^?"
	}
	return " " + string(r)
}
