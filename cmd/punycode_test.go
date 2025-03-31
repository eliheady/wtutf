package cmd

import (
	"testing"

	"golang.org/x/net/idna"
)

func TestToPuny(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		rules   []idna.Option
		want    string
		wantErr bool
	}{
		{
			name:  "Valid ASCII input",
			input: "example",
			rules: []idna.Option{},
			want:  "example",
		},
		{
			name:  "Valid Unicode input",
			input: "exämple",
			rules: []idna.Option{},
			want:  "xn--exmple-cua",
		},
		{
			name:    "Invalid input with joiners",
			input:   string([]rune{'e', 'x', 'a', 0x0308, 'm', 'p', 'l', 'e'}),
			rules:   []idna.Option{idna.CheckJoiners(true), idna.ValidateForRegistration()},
			wantErr: true,
		},
		{
			name:    "Invalid input with hyphens",
			input:   "--example--invalid",
			rules:   []idna.Option{idna.CheckHyphens(true)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := toPuny(tt.input, tt.rules)
			if tt.wantErr && err == nil {
				t.Errorf("toPuny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("toPuny() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromPuny(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		rules   []idna.Option
		want    string
		wantErr bool
	}{
		{
			name:  "Valid Punycode input",
			input: "xn--exmple-cua",
			rules: []idna.Option{},
			want:  "exämple",
		},
		{
			name:    "Invalid Punycode input",
			input:   "xn--invalid--punycode",
			rules:   []idna.Option{idna.ValidateLabels(true), idna.ValidateForRegistration(), idna.StrictDomainName(true)},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fromPuny(tt.input, tt.rules)
			if tt.wantErr && err == nil {
				t.Errorf("fromPuny() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("fromPuny() = %v, want %v", got, tt.want)
			}
		})
	}
}
