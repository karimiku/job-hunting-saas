package value

import (
	"errors"
	"testing"
)

func TestNewAlias(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"japanese alias", "ソニー", nil},
		{"english alias", "Sony", nil},

		// 異常系
		{"empty", "", ErrAliasEmpty},
		{"whitespace only", " ", ErrAliasEmpty},
		{"leading space", " ソニー", ErrAliasInvalid},
		{"trailing space", "ソニー ", ErrAliasInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAlias(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewAlias(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestAlias_String(t *testing.T) {
	a, err := NewAlias("ソニー")
	if err != nil {
		t.Fatalf("NewAlias failed: %v", err)
	}
	if a.String() != "ソニー" {
		t.Errorf("String() = %q, want %q", a.String(), "ソニー")
	}
}

func TestAlias_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same value", "ソニー", "ソニー", true},
		{"different value", "ソニー", "Sony", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			aliasA, err := NewAlias(tt.a)
			if err != nil {
				t.Fatalf("NewAlias(%q) failed: %v", tt.a, err)
			}
			aliasB, err := NewAlias(tt.b)
			if err != nil {
				t.Fatalf("NewAlias(%q) failed: %v", tt.b, err)
			}
			if got := aliasA.Equals(aliasB); got != tt.want {
				t.Errorf("Alias(%q).Equals(Alias(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
