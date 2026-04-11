package value

import (
	"errors"
	"testing"
)

func TestNewUserName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"japanese name", "田中太郎", nil},
		{"english name", "John Doe", nil},

		// 異常系
		{"empty", "", ErrUserNameEmpty},
		{"whitespace only", " ", ErrUserNameEmpty},
		{"leading space", " 田中太郎", ErrUserNameInvalid},
		{"trailing space", "田中太郎 ", ErrUserNameInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUserName(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewUserName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
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
