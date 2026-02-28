package value

import (
	"testing"
)

func TestNewUserName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// 正常系
		{"japanese name", "田中太郎", false},
		{"english name", "John Doe", false},

		// 異常系
		{"empty", "", true},
		{"whitespace only", " ", true},
		{"leading space", " 田中太郎", true},
		{"trailing space", "田中太郎 ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUserName(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("NewUserName(%q) should return error, but got nil", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("NewUserName(%q) should succeed, but got error: %v", tt.input, err)
			}
		})
	}
}

func TestUserName_String(t *testing.T) {
	n, err := NewUserName("田中太郎")
	if err != nil {
		t.Fatalf("NewUserName failed: %v", err)
	}
	if n.String() != "田中太郎" {
		t.Errorf("String() = %q, want %q", n.String(), "田中太郎")
	}
}

func TestUserName_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same value", "田中太郎", "田中太郎", true},
		{"different value", "田中太郎", "佐藤花子", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nameA, err := NewUserName(tt.a)
			if err != nil {
				t.Fatalf("NewUserName(%q) failed: %v", tt.a, err)
			}
			nameB, err := NewUserName(tt.b)
			if err != nil {
				t.Fatalf("NewUserName(%q) failed: %v", tt.b, err)
			}
			if got := nameA.Equals(nameB); got != tt.want {
				t.Errorf("UserName(%q).Equals(UserName(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
