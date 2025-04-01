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
			want: "^@",
		},
		{
			name: "Control character - DEL",
			r:    0x7F,
			want: "^?",
		},
		{
			name: "Combining diacritical mark",
			r:    0x0301, // Combining acute accent
			want: "  ◌́",
		},
		{
			name: "Deprecated format character",
			r:    0x206F,
			want: " ",
		},
		{
			name: "Directional override - LRE",
			r:    0x202A,
			want: " ",
		},
		{
			name: "Implicit directional mark - LRM",
			r:    0x200E,
			want: " ",
		},
		{
			name: "Joiner character - ZWJ",
			r:    0x200D,
			want: " ",
		},
		{
			name: "Variation selector",
			r:    0xFE0F,
			want: " ",
		},
		{
			name: "Printable ASCII character",
			r:    'A',
			want: "A",
		},
		{
			name: "Wide character - CJK Unified Ideograph",
			r:    0x4E00,
			want: "一",
		},
		{
			name: "Emoji character",
			r:    0x1F600, // Grinning face
			want: "😀",
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
