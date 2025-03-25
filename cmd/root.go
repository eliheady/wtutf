/*
Copyright © 2025 Eli Heady
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

// toPuny takes a rune and an idna.Option rule and calls idna.ToASCII,
// returning the punycode string and error
func toPuny(s string, rules []idna.Option) (string, error) {
	punyRules := idna.New(
		rules...,
	)
	return punyRules.ToASCII(s)
}

// canPunyConvert takes a rune and an idna.Option rule and attempts the ToASCII
// conversion, then reports success
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

	header := "    code point | bytes "
	if !converted {
		header += " | conversion errors"
	}
	out += header + "\n"

	for _, r := range ustring {
		safeToPrint := r
		moreErrors := ""
		var padded string
		switch {
		case r < 0x20:
			safeToPrint = ' '
			padded = fmt.Sprintf("%#02x", r)
		case r > 0xff:
			padded = fmt.Sprintf("%#04x", r)
		case r > 0xffff:
			padded = fmt.Sprintf("%#06x", r)
		default:
			padded = fmt.Sprintf("%#02x", r)
		}
		if !converted {
			moreErrors = fmt.Sprintf("\t| %s", enumerateErrors(r))
		}
		out += fmt.Sprintf("%c:\t% 6s | (%d)%s\n", safeToPrint, padded, len(string(r)), moreErrors)
	}

	return out
}
