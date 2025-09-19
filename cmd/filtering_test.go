package cmd

import (
	"testing"
)

func TestPolitePrint(t *testing.T) {
	tests := []struct {
		name string
		r    rune
		want string
	}{
		{
			name: "Control character - NULL",
			r:    0x00,
			want: "^?",
		},
		{
			name: "Control character - DEL",
			r:    0x7F,
			want: "^?",
		},
		{
			name: "C1 control - 0x9A",
			r:    0x9A,
			want: "^?",
		},
		{
			name: "Combining diacritical mark",
			r:    0x0301, // Combining acute accent
			want: " ◌́",
		},
		{
			name: "Combining diacritical mark, 2 chars",
			r:    0x0361, // Double inverted breve
			want: "◌͡◌",
		},
		{
			name: "Deprecated format character",
			r:    0x206F,
			want: "^?",
		},
		{
			name: "Directional override - LRE",
			r:    0x202A,
			want: "^?",
		},
		{
			name: "Implicit directional mark - LRM",
			r:    0x200E,
			want: "^?",
		},
		{
			name: "Joiner character - ZWJ",
			r:    0x200D,
			want: "^?",
		},
		{
			name: "Variation selector",
			r:    0xFE0F,
			want: "^?",
		},
		{
			name: "Printable ASCII character",
			r:    'A',
			want: " A",
		},
		{
			name: "Wide character - CJK Unified Ideograph",
			r:    0x4E00,
			want: " 一",
		},
		{
			name: "Piñata emoji",
			r:    0x1FA85,
			want: " 🪅",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := politePrint(tt.r)
			if got != tt.want {
				t.Errorf("politePrint(%U) = %q, want %q", tt.r, got, tt.want)
			}
		})
	}
}

func TestCheckMultipleRange(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		showRanges  bool
		wantMulti   bool
		wantPrinted []string
	}{
		{"single ASCII range", "hello", false, false, nil},
		{"single ASCII with showRanges", "hello", true, false, nil},
		{"multi range Latin+Greek", "aα", false, true, nil},
		{"multi range printed", "aα", true, true, []string{"Latin", "Greek"}},
		{"common only ignored", "!@#", true, false, nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// compute expected multi-range from listRanges (ignores showRanges)
			lr := listRanges(tc.input)
			ranges := 0
			for name := range lr {
				if name == "Common" {
					continue
				}
				ranges++
			}
			expectedMulti := ranges > 1
			if expectedMulti != tc.wantMulti {
				t.Fatalf("computed expectedMulti=%v for input %q, but test expects %v", expectedMulti, tc.input, tc.wantMulti)
			}

			// if showRanges was requested, verify UnicodeRanges contains expected names
			if tc.showRanges && len(tc.wantPrinted) > 0 {
				// gatherOutputData populates UnicodeRanges when showRanges is true
				data := gatherOutputData(tc.input, tc.showRanges, false, false, false)
				for _, want := range tc.wantPrinted {
					if _, ok := data.UnicodeRanges[want]; !ok {
						t.Errorf("expected UnicodeRanges to contain %q, got %v", want, data.UnicodeRanges)
					}
				}
			}
		})
	}
}
