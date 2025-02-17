package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "w",
	Short: "Process a string of unknown encoding",
	Long:  `Process a string of unknown encoding`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(printMystery())
	},
}

func printMystery() string {
	const ntildeSingle = `ñ`
	//const ntildeComposed = "\x6e\x03\x03"
	const ntildeComposed = "n\u0303"
	const ntildeComposedLiteral = `ñ`

	// var composedChars []rune
	// composedChars[0] = ntildeComposed
	// composedChars[1] = '\u006E'
	// composed := string(composedChars)

	// chars := []interface{}{
	// 	ntntildeSingle,
	// 	ntntildeComposed,
	// }
	chars := []string{ntildeSingle, ntildeComposed, ntildeComposedLiteral}
	fmt.Printf("%s\n", chars[0])

	var out string

	for i, char := range chars {
		if i == 0 {
			out += "single:\n"
		} else {
			out += "composite:\n"
		}
		out += "%s format string: "
		out += fmt.Sprintf("%s\n", chars[0])
		out += "plain string: "
		out += char
		out += "\n"

		out += "quoted string: "
		out += fmt.Sprintf("%+q\n", char)

		out += "hex bytes: "
		for i := 0; i < len(char); i++ {
			out += fmt.Sprintf("% x", char[i])
		}
		out += "\n\n"
	}

	out += fmt.Sprintf("%s == %s: %v\n", ntildeSingle, ntildeComposed, bool(ntildeComposed == ntildeSingle))
	out += fmt.Sprintf("%s == %s: %v\n", ntildeComposed, ntildeComposedLiteral, bool(ntildeComposed == ntildeComposedLiteral))
	return out
}
