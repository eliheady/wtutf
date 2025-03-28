/*
WTUTF: A simple UTF-8 string inspector.

# Copyright © 2025 Eli Heady

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package cmd

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"golang.org/x/net/idna"
)

var rootCmd = &cobra.Command{
	Use:   "wtutf",
	Args:  cobra.ExactArgs(1),
	Short: "A simple utility to help me out of my ASCII-centric shell",
	Long:  `This program just prints out the Unicode code points of the string you feed into it. It can also show you the punycode conversion of your string, or failure reasons if conversion isn't possible.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(processInput(cmd, args[0]))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var strict bool
	rootCmd.PersistentFlags().BoolVarP(&strict, "strict", "s", false, "Set strict punycode conversion rules")
}

// toPuny takes a string and a slice of []idna.Option rules and calls
// idna.ToASCII, returning the punycode string and error
func toPuny(s string, rules []idna.Option) (string, error) {
	punyRules := idna.New(
		rules...,
	)
	return punyRules.ToASCII(s)
}

// canPunyConvert takes a rune and a slice of []idna.Option rules and attempts
// the ToASCII conversion, then returns a bool indicating success
func canPunyConvert(s string, rules []idna.Option) bool {
	_, err := toPuny(s, rules)
	return err == nil
}

// enumerateErrors takes a rune, checks several punycode conversion rules and
// reports the failures as a single string
func enumerateErrors(r rune) string {
	rules := map[string][]idna.Option{
		"CheckBidi (RFC 5893)":                       {idna.BidiRule()},
		"CheckJoiners (RFC 5892)":                    {idna.CheckJoiners(true)},
		"CheckHyphens (UTS 46)":                      {idna.CheckHyphens(true)},
		"ValidateForRegistration (RFC 5891)":         {idna.ValidateForRegistration()},
		"ValidateLabels (RFC 5891)":                  {idna.ValidateLabels(true)},
		"UseSTD3ASCIIRules (RFC 1034, 5891, UTS 46)": {idna.StrictDomainName(true), idna.ValidateLabels(true)},
	}

	var allErrors []string

	for i, ruleset := range rules {
		if !canPunyConvert(string(r), ruleset) {
			allErrors = append(allErrors, i)
		}
	}

	return strings.Join(allErrors, ", ")

}

// politePrint takes a rune and outputs a human readable conversion
func politePrint(r rune) string {
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
	// todo: these filters are probably available in some form in the utf8 package
	case int(r) < len(politeCharmap):
		return politeCharmap[r] // C0 controls
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
	}
	return string(r)
}

func processInput(cmd *cobra.Command, ustring string) string {

	flags := cmd.Flags()
	var out string

	out += fmt.Sprintf("total bytes:\t%d\n", len(ustring))
	out += fmt.Sprintf("characters:\t%d\n", utf8.RuneCountInString(ustring))

	rules := []idna.Option{
		idna.BidiRule(),
		idna.CheckJoiners(true),
		idna.ValidateLabels(true),
	}
	if strict, _ := flags.GetBool("strict"); strict { // ok to discard err if flag was not set
		rules = append(rules,
			idna.ValidateForRegistration(),
			idna.StrictDomainName(true),
		)
	}

	var converted bool
	out += "punycode:\t"
	if punycode, err := toPuny(ustring, rules); err == nil {
		converted = true
		out += punycode
	} else {
		out += "could not punycode-convert input"
	}
	out += "\n"

	header := "      code point | bytes "
	if !converted {
		header += "| conversion rules violated"
	}
	out += header + "\n"

	for _, r := range ustring {
		// todo: implement caching of the rune validation and error enumeration
		moreErrors := ""
		var padded string
		printable := fmt.Sprintf("% 3s", politePrint(r))
		switch {
		case r > 0xFFFF:
			padded = fmt.Sprintf("%#06x", r)
		case r > 0xFF:
			padded = fmt.Sprintf("%#04x", r)
		default:
			padded = fmt.Sprintf("%#02x", r)
		}
		if !converted {
			moreErrors = fmt.Sprintf(" | %s", enumerateErrors(r))
		}
		out += fmt.Sprintf("%s:\t% 8s |  (%d) %s\n", printable, padded, utf8.RuneLen(r), moreErrors)
	}

	return out
}
