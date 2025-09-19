package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

// helper creates a cobra.Command with the same flags parseFlags expects
func newTestCmd() *cobra.Command {
    c := &cobra.Command{Use: "test"}
    // register flags that parseFlags reads
    c.Flags().BoolP("check", "c", false, "")
    c.Flags().BoolP("show-ranges", "r", false, "")
    c.Flags().BoolP("strict", "s", false, "")
    c.Flags().BoolP("puny", "p", false, "")
    c.Flags().BoolP("table", "t", false, "")
    c.Flags().Bool("json", false, "")
    return c
}

func TestParseFlagsOutputs(t *testing.T) {
    tests := []struct {
        name       string
        setFlags   map[string]string // flag name -> value
        args       []string
        wantJSON   bool
        wantContain string
    }{
        {
            name: "json output",
            setFlags: map[string]string{"json": "true"},
            args: []string{"café"},
            wantJSON: true,
        },
        {
            name: "table output",
            setFlags: map[string]string{"table": "true"},
            args: []string{"café"},
            wantJSON: false,
            wantContain: "code point",
        },
        {
            name: "json takes precedence over table",
            setFlags: map[string]string{"json": "true", "table": "true"},
            args: []string{"café"},
            wantJSON: true,
        },
        {
            name: "puny decode shows utf-8 line",
            setFlags: map[string]string{"puny": "true"},
            args: []string{"xn--piatasafety-2db"},
            wantJSON: false,
            wantContain: "utf-8:",
        },
    }

    for _, tc := range tests {
        t.Run(tc.name, func(t *testing.T) {
            cmd := newTestCmd()
            // apply flags
            for k, v := range tc.setFlags {
                if err := cmd.Flags().Set(k, v); err != nil {
                    t.Fatalf("failed to set flag %s=%s: %v", k, v, err)
                }
            }

            out := parseFlags(cmd, tc.args)

            if tc.wantJSON {
                // should be valid JSON and map to OutputData
                var data OutputData
                if err := json.Unmarshal([]byte(out), &data); err != nil {
                    t.Fatalf("expected valid JSON output but got error: %v\noutput: %s", err, out)
                }
                if data.Input != tc.args[0] {
                    t.Fatalf("unexpected input in JSON output: got %q want %q", data.Input, tc.args[0])
                }
            }

            if tc.wantContain != "" {
                if !strings.Contains(out, tc.wantContain) {
                    t.Fatalf("output does not contain %q: %s", tc.wantContain, out)
                }
            }
        })
    }
}
