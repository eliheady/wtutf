// WTUTF: A simple UTF-8 string inspector.
package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"unicode"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"golang.org/x/net/idna"
)

// OutputData holds the structured output for both text and JSON formats
type OutputData struct {
	Input         string         `json:"input"`
	Punycode      string         `json:"punycode,omitempty"`
	UTF8          string         `json:"utf8,omitempty"`
	PunycodeError string         `json:"punycode_error,omitempty"`
	TotalBytes    int            `json:"total_bytes"`
	Characters    int            `json:"characters"`
	UnicodeRanges map[string]int `json:"unicode_ranges,omitempty"`
	Table         []RuneTableRow `json:"table,omitempty"`
}

type RuneTableRow struct {
	Printable string   `json:"printable"`
	CodePoint string   `json:"code_point"`
	Bytes     string   `json:"bytes"`
	Length    int      `json:"length"`
	Errors    []string `json:"errors,omitempty"`
}

type RuneCache struct {
	Printable string
	Padded    string
	Bytes     string
	Errors    []string
}

var rootCmd = &cobra.Command{
	Use:   "wtutf",
	Args:  cobra.ExactArgs(1),
	Short: "A simple utility to reduce ASCII-centrism",
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
	var check, showRanges, strict, fromPuny, table, jsonOut bool
	rootCmd.PersistentFlags().BoolVarP(&check, "check", "c", false, "Check whether the string contains characters from more than one Unicode range")
	rootCmd.PersistentFlags().BoolVarP(&showRanges, "show-ranges", "r", false, "Show the Unicode script ranges included in the string")
	rootCmd.PersistentFlags().BoolVarP(&strict, "strict", "s", false, "Set strict punycode conversion rules")
	rootCmd.PersistentFlags().BoolVarP(&fromPuny, "puny", "p", false, "Convert from punycode")
	rootCmd.PersistentFlags().BoolVarP(&table, "table", "t", false, "Show table of all included unicode characters")
	rootCmd.PersistentFlags().BoolVar(&jsonOut, "json", false, "Output results as JSON instead of plain text")
}

func parseFlags(cmd *cobra.Command, args []string) string {
	flags := cmd.Flags()

	showRanges, _ := flags.GetBool("show-ranges")
	strict, _ := flags.GetBool("strict")
	punyDecode, _ := flags.GetBool("puny")
	table, _ := flags.GetBool("table")
	jsonOut, _ := flags.GetBool("json")

	if compare, _ := flags.GetBool("check"); compare {
		var checkResult int
		if checkMultipleRange(showRanges, args[0]) {
			checkResult = 1
		}
		os.Exit(checkResult)
	}

	data := gatherOutputData(args[0], showRanges, strict, punyDecode, table)
	if jsonOut {
		b, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			return "Error encoding JSON: " + err.Error()
		}
		return string(b) + "\n"
	}
	return formatPlainText(data, showRanges, table)
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

// gatherOutputData collects all output data for a given input string
func gatherOutputData(ustring string, showRanges, strict, punyDecode, table bool) OutputData {
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

	data := OutputData{
		Input:      ustring,
		TotalBytes: len(ustring),
		Characters: utf8.RuneCountInString(ustring),
	}

	var punyConverted bool

	if punyDecode {
		if utfString, err := fromPuny(ustring, rules); err == nil {
			data.Punycode = ustring
			data.UTF8 = utfString
			ustring = utfString
		} else {
			data.PunycodeError = "could not decode punycode input"
		}
	} else {
		if punycode, err := toPuny(ustring, rules); err == nil {
			punyConverted = true
			data.Punycode = punycode
		} else {
			data.PunycodeError = "could not punycode-convert input"
		}
	}

	if showRanges {
		data.UnicodeRanges = listRanges(ustring)
	}

	if table {
		runeCache := map[rune]RuneCache{}
		for _, r := range ustring {
			var printable, padded, runeBytes string
			var runeErrors []string
			if rc, ok := runeCache[r]; ok {
				printable = rc.Printable
				padded = rc.Padded
				runeBytes = rc.Bytes
				runeErrors = rc.Errors
			} else {
				printable = toPaddedString(r, 3)
				padded = fmt.Sprintf("%#0*x", (utf8.RuneLen(r) * 2), r)
				runeBytes = hex.EncodeToString([]byte(string(r)))
				if !punyConverted {
					runeErrors = enumerateErrors(r)
				}
				runeCache[r] = RuneCache{
					Printable: printable,
					Padded:    padded,
					Bytes:     runeBytes,
					Errors:    runeErrors,
				}
			}
			row := RuneTableRow{
				Printable: printable,
				CodePoint: padded,
				Bytes:     runeBytes,
				Length:    utf8.RuneLen(r),
				Errors:    runeErrors,
			}
			data.Table = append(data.Table, row)
		}
	}

	return data
}

// formatPlainText renders OutputData as a plain text table output
func formatPlainText(data OutputData, showRanges, table bool) string {
	var b strings.Builder
	tw := tabwriter.NewWriter(&b, 0, 0, 2, ' ', 0)

	// Main output fields
	if data.Punycode != "" && data.UTF8 != "" {
		fmt.Fprintf(tw, "punycode:\t%s\n", data.Punycode)
		fmt.Fprintf(tw, "utf-8:\t%s\n", data.UTF8)
	} else if data.Punycode != "" {
		fmt.Fprintf(tw, "punycode:\t%s\n", data.Punycode)
	} else if data.PunycodeError != "" {
		fmt.Fprintf(tw, "%s\n", data.PunycodeError)
	}
	fmt.Fprintf(tw, "total bytes:\t%d\n", data.TotalBytes)
	fmt.Fprintf(tw, "characters:\t%d\n", data.Characters)

	if showRanges && data.UnicodeRanges != nil {
		fmt.Fprintf(tw, "unicode ranges:\n")
		for i, count := range data.UnicodeRanges {
			fmt.Fprintf(tw, "\t%s:	%d\n", i, count)
		}
	}

	if table && len(data.Table) > 0 {
		fmt.Fprintf(tw, "----------------------------------\n")
		header := []string{"printable", "code point", "bytes (len)"}
		hasErrors := false
		for _, row := range data.Table {
			if len(row.Errors) > 0 {
				hasErrors = true
				break
			}
		}
		if hasErrors {
			header = append(header, "conversion rules violated")
		}
		fmt.Fprintf(tw, "%s\t%s\t%s", header[0], header[1], header[2])
		if len(header) > 3 {
			fmt.Fprintf(tw, "\t%s", header[3])
		}
		fmt.Fprintf(tw, "\n")
		for _, row := range data.Table {
			bytesColumn := fmt.Sprintf("%s (%d)", row.Bytes, row.Length)
			errors := strings.Join(row.Errors, ", ")
			fmt.Fprintf(tw, "%s\t%s\t%s", row.Printable, row.CodePoint, bytesColumn)
			if hasErrors {
				fmt.Fprintf(tw, "\t%s", errors)
			}
			fmt.Fprintf(tw, "\n")
		}
	}
	tw.Flush()
	return b.String()
}
