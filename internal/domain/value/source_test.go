package value

import (
	"testing"
)

func TestNewSource(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// 正常系
		{"japanese", "マイナビ", false},
		{"japanese 2", "リクナビ", false},
		{"company hp", "企業HP", false},
		{"free input", "友人紹介", false},
		{"english", "OfferBox", false},

		// 異常系
		{"empty", "", true},
		{"whitespace only", " ", true},
		{"leading space", " マイナビ", true},
		{"trailing space", "マイナビ ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSource(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("NewSource(%q) should return error, but got nil", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("NewSource(%q) should succeed, but got error: %v", tt.input, err)
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
