// WTUTF: A simple UTF-8 string inspector.
package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"unicode"
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
		fmt.Print(parseFlags(cmd, args))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	var check, showRanges, strict, fromPuny, table bool
	rootCmd.PersistentFlags().BoolVarP(&check, "check", "c", false, "Check whether the string contains characters from more than one Unicode range")
	rootCmd.PersistentFlags().BoolVarP(&showRanges, "show-ranges", "r", false, "Show the Unicode script ranges included in the string")
	rootCmd.PersistentFlags().BoolVarP(&strict, "strict", "s", false, "Set strict punycode conversion rules")
	rootCmd.PersistentFlags().BoolVarP(&fromPuny, "puny", "p", false, "Convert from punycode")
	rootCmd.PersistentFlags().BoolVarP(&table, "table", "t", false, "Show table of all included unicode characters")
}

func parseFlags(cmd *cobra.Command, args []string) string {
	flags := cmd.Flags()

	showRanges, _ := flags.GetBool("show-ranges")
	strict, _ := flags.GetBool("strict")
	punyDecode, _ := flags.GetBool("puny")
	table, _ := flags.GetBool("table")

	if compare, _ := flags.GetBool("check"); compare {
		var checkResult int
		if checkMultipleRange(showRanges, args[0]) {
			checkResult = 1
		}
		os.Exit(checkResult)
	}

	return processInput(args[0], showRanges, strict, punyDecode, table)
}

// toString takes a rune returns a string with padding appropriate for the character width
func toPaddedString(r rune, defaultPad int) (out string) {
	defaultPattern := fmt.Sprintf("%% %ds", defaultPad)
	widePadLeft := fmt.Sprintf("%% %ds", defaultPad-2)
	widePattern := fmt.Sprintf("%s%%-1s", widePadLeft)
	// Unicode blocks containing 'wide' or combining characters
	fullwidth := 0xFF00 <= r && r <= 0xFFEF
	symbols := unicode.Is(unicode.Symbol, r)
	wide := unicode.Is(unicode.Ideographic, r) || unicode.Is(unicode.Extender, r)
	wide = wide || unicode.Is(unicode.Diacritic, r) || unicode.Is(unicode.Hangul, r)
	wide = wide || unicode.Is(unicode.Bopomofo, r) || unicode.Is(unicode.Han, r)
	wide = wide || unicode.Is(unicode.Hiragana, r) || unicode.Is(unicode.Katakana, r)
	wide = wide || fullwidth || symbols
	switch {
	case wide:
		out = fmt.Sprintf(widePattern, " ", politePrint(r))
	default:
		out = fmt.Sprintf(defaultPattern, politePrint(r))
	}
	return
}

type RuneCache struct {
	Printable string
	Padded    string
	Bytes     string
	Errors    []string
}

func processInput(ustring string, showRanges, strict, punyDecode, table bool) string {

	rules := []idna.Option{
		idna.BidiRule(),
		idna.CheckJoiners(true),
		idna.ValidateLabels(true),
	}
	if strict {
		rules = append(rules,
			idna.ValidateForRegistration(),
			idna.StrictDomainName(true),
		)
	}

	// todo: separate accumulation of output from
	// formatting, delegate to a printing function

	var out string
	var punyConverted bool

	if punyDecode { // ok to discard err if flag was not set
		if utfString, err := fromPuny(ustring, rules); err == nil {
			out += "      punycode:\t" + ustring + "\n"
			out += "         utf-8:\t" + utfString + "\n"
			ustring = utfString
		} else {
			out += "could not decode punycode input\n"
		}
	} else {
		if punycode, err := toPuny(ustring, rules); err == nil {
			punyConverted = true
			out += "      punycode:\t" + punycode + "\n"
		} else {
			out += "could not punycode-convert input\n"
		}
	}
	out += fmt.Sprintf("   total bytes:\t%d\n", len(ustring))
	out += fmt.Sprintf("    characters:\t%d\n", utf8.RuneCountInString(ustring))

	if showRanges {
		ranges := listRanges(ustring)
		out += "unicode ranges:\n"
		for i, count := range ranges {
			out += fmt.Sprintf("    %s: %d\n", i, count)
		}
	}

	if table {
		out += "----------------------------------\n"

		header := []string{
			fmt.Sprintf("% 17s", "code point"),
			fmt.Sprintf("% 12s", "bytes (len)"),
		}
		if !punyConverted {
			header = append(header, "conversion rules violated")
		}
		out += strings.Join(header, " | ") + "\n"

		// cache of rune validation and error enumeration
		runeCache := map[rune]RuneCache{}

		for _, r := range ustring {
			runeErrors := []string{}
			var padded, runeBytes, printable string
			if rc, ok := runeCache[r]; ok {
				printable = rc.Printable
				padded = rc.Padded
				runeBytes = rc.Bytes
				runeErrors = rc.Errors
			} else {

				printable = toPaddedString(r, 3)

				// pad the rune value with leading zeroes for every byte
				padded = fmt.Sprintf("%#0*x", (utf8.RuneLen(r) * 2), r)

				runeBytes = hex.EncodeToString([]byte(string(r)))

				if !punyConverted {
					// only check the individual runes if the whole string
					// has failed the punycode conversion.
					runeErrors = enumerateErrors(r)
				}
				runeCache[r] = RuneCache{
					Printable: printable,
					Padded:    padded,
					Bytes:     runeBytes,
					Errors:    runeErrors,
				}
			}

			bytesColumn := fmt.Sprintf("% 8s (%d)", runeBytes, utf8.RuneLen(r))
			errors := strings.Join(runeErrors, ", ")
			out += fmt.Sprintf("%s:% 13s | %s | %s\n", printable, padded, bytesColumn, errors)
		}
	}

	return out
}
