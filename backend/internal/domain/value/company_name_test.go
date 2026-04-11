package value

import (
	"errors"
	"testing"
)

func TestNewCompanyName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"japanese name", "トヨタ自動車", nil},
		{"english name", "Google", nil},

		// 異常系
		{"empty", "", ErrCompanyNameEmpty},
		{"whitespace only", " ", ErrCompanyNameEmpty},
		{"leading space", " トヨタ自動車", ErrCompanyNameInvalid},
		{"trailing space", "トヨタ自動車 ", ErrCompanyNameInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewCompanyName(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewCompanyName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestCompanyName_String(t *testing.T) {
	n, err := NewCompanyName("トヨタ自動車")
	if err != nil {
		t.Fatalf("NewCompanyName failed: %v", err)
	}
	if n.String() != "トヨタ自動車" {
		t.Errorf("String() = %q, want %q", n.String(), "トヨタ自動車")
	}
}

func TestCompanyName_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same value", "トヨタ自動車", "トヨタ自動車", true},
		{"different value", "トヨタ自動車", "ソニー", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nameA, err := NewCompanyName(tt.a)
			if err != nil {
				t.Fatalf("NewCompanyName(%q) failed: %v", tt.a, err)
			}
			nameB, err := NewCompanyName(tt.b)
			if err != nil {
				t.Fatalf("NewCompanyName(%q) failed: %v", tt.b, err)
			}
			if got := nameA.Equals(nameB); got != tt.want {
				t.Errorf("CompanyName(%q).Equals(CompanyName(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
