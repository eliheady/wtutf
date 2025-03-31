package cmd

import "fmt"

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
func checkMultipleRange(ustring string) (multiRange bool) {
	var ranges int
	var out string
	for i, count := range listRanges(ustring) {
		if i == "Common" {
			continue
		}
		ranges++
		out += fmt.Sprintf("%s: %d\n", i, count)
	}
	if ranges > 1 {
		multiRange = true
		fmt.Print(out)
	}
	return
}

// politePrint takes a rune and outputs a human readable conversion
func politePrint(r rune) string {
	// todo: replace these filters with ranges defined in unicode/table.go

	// there's a box of mysteries to explore if you print random crap
	// into a terminal. Let's try to be nice and avoid trashing the
	// user's environment. It was their input, but they could be
	// drunk. Or a cat.
	politeCharmap := []string{
		"^@", "^A", "^B", "^C", "^D", "^E", "^F", "^G", "^H", "^I", "^J", "^K", "^L",
		"^M", "^N", "^O", "^P", "^Q", "^R", "^S", "^T", "^U", "^V", "^W", "^X", "^Y",
		"^Z", "^[", "^\\", "^]", "^^", "^_",
	}

	switch {
	case 0x035C <= r && r <= 0x0362: // combining diacritical marks, 2 characters
		return " ◌" + string(r) + "◌"
	case 0x0300 <= r && r <= 0x036F: // combining diacritical marks
		return "  ◌" + string(r)
	// filter control characters
	// reasoning: terminal control sequences can do all sorts of damage to the
	// output.  we will remove them and put in the caret notation for C0 and 'C1'
	// for C1 unicode control characters
	case int(r) < len(politeCharmap):
		return politeCharmap[int(r)] // C0 controls
	case int(r) == 127:
		return "^?" // DEL
	case 0x206A <= r && r <= 0x206F: // deprecated format characters
		return " "
	case 0x80 <= r && r <= 0x9F:
		return "C1" // Unicode C1 controls
	// filter direction changing characters
	// reasoning: we are printing a single character here. If we needed to print
	// words then the directional marks should not be discarded. These checks are
	// an attempt to prevent leaving an unclosed direction change in the output.
	// ref. https://www.unicode.org/reports/tr9/#Directional_Formatting_Codes
	case 0x202A <= r && r <= 0x202E: // directional overrides (LRE, RLE, PDF, LRO, RLO)
		return " "
	case r == 0x061C, r == 0x200E, r == 0x200F: // implicit directional marks
		return " "
	case 0x2066 <= r && r <= 0x2069: // directional isolates (LRI, RLI, FSI, PDI)
		return " "
	case r == 0x200d || r == 0x2060: // joiners
		return " "
	case 0xFE00 <= r && r <= 0xFE0F: // variation selectors
		return " "
	}
	return string(r)
}
