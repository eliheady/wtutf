/*
Copyright © 2025 Eli Heady
*/
package cmd

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/spf13/cobra"
	"golang.org/x/net/idna"
)

var rootCmd = &cobra.Command{
	Use:   "wtutf",
	Args:  cobra.ExactArgs(1),
	Short: "A simple utility to help me out of my ASCII-centric shell",
	Long: `This program just prints out the Unicode code points of the string you feed into it. It can also show you the punycode conversion.

$ wtutf piñata
total bytes: 8
characters: 7
    code point  (bytes)
p:      0x70    (1)
i:      0x69    (1)
n:      0x6e    (1)
̃:      0x303   (2)
a:      0x61    (1)
t:      0x74    (1)
a:      0x61    (1)
punycode:
could not punycode-convert input

$ wtutf piñata  
total bytes: 7
characters: 6
    code point  (bytes)
p:      0x70    (1)
i:      0x69    (1)
ñ:      0xf1    (2)
a:      0x61    (1)
t:      0x74    (1)
a:      0x61    (1)
punycode:
xn--piata-pta
.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(processInput(args[0]))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func processInput(ustring string) string {

	if ustring == "" {
		fmt.Println("uh, sure ... the empty string is")
	}

	var out string

	out += fmt.Sprintf("total bytes: %d\n", len(ustring))
	out += fmt.Sprintf("characters: %d\n", utf8.RuneCountInString(ustring))

	out += "    code point\t(bytes)\n"
	for _, r := range ustring {
		out += fmt.Sprintf("%c:\t%#02x\t(%d)\n", r, r, len(string(r)))
	}

	puny := idna.New(
		idna.BidiRule(),
		idna.CheckJoiners(true),
		idna.CheckHyphens(true),
		idna.ValidateForRegistration(),
		idna.ValidateLabels(true),
	)
	out += "punycode:\n"
	if pc, err := puny.ToASCII(ustring); err == nil {
		out += pc
	} else {
		out += "could not punycode-convert input"
	}
	out += "\n"
	return out
}
