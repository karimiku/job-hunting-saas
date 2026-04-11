package value

import (
	"errors"
	"testing"
)

func TestNewSource(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"japanese", "マイナビ", nil},
		{"japanese 2", "リクナビ", nil},
		{"company hp", "企業HP", nil},
		{"free input", "友人紹介", nil},
		{"english", "OfferBox", nil},

		// 異常系
		{"empty", "", ErrSourceEmpty},
		{"whitespace only", " ", ErrSourceEmpty},
		{"leading space", " マイナビ", ErrSourceInvalid},
		{"trailing space", "マイナビ ", ErrSourceInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSource(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewSource(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestSource_String(t *testing.T) {
	s, err := NewSource("マイナビ")
	if err != nil {
		t.Fatalf("NewSource failed: %v", err)
	}
	if s.String() != "マイナビ" {
		t.Errorf("String() = %q, want %q", s.String(), "マイナビ")
	}
}

func TestSource_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same value", "マイナビ", "マイナビ", true},
		{"different value", "マイナビ", "リクナビ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srcA, err := NewSource(tt.a)
			if err != nil {
				t.Fatalf("NewSource(%q) failed: %v", tt.a, err)
			}
			srcB, err := NewSource(tt.b)
			if err != nil {
				t.Fatalf("NewSource(%q) failed: %v", tt.b, err)
			}
			if got := srcA.Equals(srcB); got != tt.want {
				t.Errorf("Source(%q).Equals(Source(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
