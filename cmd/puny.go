package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/net/idna"
)

func init() {
	rootCmd.AddCommand(punyCmd)
}

var punyCmd = &cobra.Command{
	Use:   "puny",
	Short: "Convert to Punycode",
	Long:  `Convert to Punycode`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(printPuny())
	},
}

func printPuny() string {
	const ntildeSingle = `piñatasafety`
	//const ntildeComposed = "\x6e\x03\x03"
	const ntildeComposed = "pin\u0303atasafety"
	const ntildeComposedLiteral = `piñatasafety`

	// var composedChars []rune
	// composedChars[0] = ntildeComposed
	// composedChars[1] = '\u006E'
	// composed := string(composedChars)

	// chars := []interface{}{
	// 	ntntildeSingle,
	// 	ntntildeComposed,
	// }
	words := []string{ntildeSingle, ntildeComposed, ntildeComposedLiteral}

	var out string

	for i, w := range words {
		if i == 0 {
			out += "single:\n"
		} else {
			out += "composite:\n"
		}
		out += "%s format string: "
		out += fmt.Sprintf("%s\n", w)
		out += "plain string: "
		out += w
		out += "\n"

		out += "\n\n"
	}

	out += fmt.Sprintf("%s == %s: %v\n", ntildeSingle, ntildeComposed, bool(ntildeComposed == ntildeSingle))
	out += fmt.Sprintf("%s == %s: %v\n", ntildeComposed, ntildeComposedLiteral, bool(ntildeComposed == ntildeComposedLiteral))

	puny := idna.New(
		idna.BidiRule(),
		idna.CheckJoiners(true),
		idna.CheckHyphens(true),
		idna.ValidateForRegistration(),
		idna.ValidateLabels(true),
	)
	out += "punycode:\n"
	if pc, err := puny.ToASCII(ntildeComposed); err == nil {
		out += pc
	} else {
		out += "could not punycode-convert input"
	}
	out += "\n"
	return out
}
