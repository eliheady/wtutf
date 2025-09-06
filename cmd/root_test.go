package cmd

import (
	"strings"
	"testing"
)

func TestTableOutput(t *testing.T) {
	input := "café"
	data := gatherOutputData(input, false, false, false, true)
	outStr := formatPlainText(data, false, true)

	if !strings.Contains(outStr, "code point") || !strings.Contains(outStr, "bytes (len)") {
		t.Errorf("table header missing in output: %s", outStr)
	}
	if !strings.Contains(outStr, "c") || !strings.Contains(outStr, "a") || !strings.Contains(outStr, "f") || !strings.Contains(outStr, "é") {
		t.Errorf("expected runes not found in table output: %s", outStr)
	}
}

func TestJSONOutput(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		showRanges bool
		wantRunes  []string
	}{
		{
			name:       "basic table",
			input:      "café",
			showRanges: false,
			wantRunes:  []string{"c", "a", "f", "é"},
		},
		{
			name:       "table with unicode ranges",
			input:      "café",
			showRanges: true,
			wantRunes:  []string{"c", "a", "f", "é"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := gatherOutputData(tc.input, tc.showRanges, false, false, true)
			if data.Input != tc.input {
				t.Errorf("expected input %q, got %q", tc.input, data.Input)
			}
			if len(data.Table) == 0 {
				t.Errorf("expected non-empty Table, got empty")
			}
			for _, wantRune := range tc.wantRunes {
				found := false
				for _, row := range data.Table {
					if strings.TrimSpace(row.Printable) == wantRune {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected rune %q not found in Table: %+v", wantRune, data.Table)
				}
			}
			if tc.showRanges && len(data.UnicodeRanges) == 0 {
				t.Errorf("expected UnicodeRanges to be populated when showRanges is true")
			}
			if !tc.showRanges && len(data.UnicodeRanges) != 0 {
				t.Errorf("expected UnicodeRanges to be empty when showRanges is false")
			}
		})
	}
}
