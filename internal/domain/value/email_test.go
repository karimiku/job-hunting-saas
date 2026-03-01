package value

import (
	"errors"
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"basic email", "user@example.com", nil},
		{"subdomain", "test.name@company.co.jp", nil},
		{"minimal", "a@b.co", nil},
		{"underscore", "user_name@example.com", nil},
		{"hyphen", "user-name@example.com", nil},
		{"plus tag", "user+tag@example.com", nil},
		{"uppercase input", "USER@EXAMPLE.COM", nil},

		// 異常系
		{"empty", "", ErrEmailEmpty},
		{"no at sign", "invalid", ErrEmailInvalid},
		{"no local part", "@example.com", ErrEmailInvalid},
		{"no domain", "user@", ErrEmailInvalid},
		{"domain starts with dot", "user@.com", ErrEmailInvalid},
		{"no tld", "user@example", ErrEmailInvalid},
		{"double at sign", "user@@example.com", ErrEmailInvalid},
		{"leading dot local", ".user@example.com", ErrEmailInvalid},
		{"trailing dot local", "user.@example.com", ErrEmailInvalid},
		{"consecutive dots local", "us..er@example.com", ErrEmailInvalid},
		{"consecutive dots domain", "user@example..com", ErrEmailInvalid},
		{"leading space", " user@example.com", ErrEmailInvalid},
		{"trailing space", "user@example.com ", ErrEmailInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmail(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewEmail(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestEmail_String(t *testing.T) {
	email, err := NewEmail("user@example.com")
	if err != nil {
		t.Fatalf("NewEmail failed: %v", err)
	}
	if email.String() != "user@example.com" {
		t.Errorf("String() = %q, want %q", email.String(), "user@example.com")
	}
}

func TestEmail_String_lowercased(t *testing.T) {
	email, err := NewEmail("USER@EXAMPLE.COM")
	if err != nil {
		t.Fatalf("NewEmail failed: %v", err)
	}
	if email.String() != "user@example.com" {
		t.Errorf("String() = %q, want %q", email.String(), "user@example.com")
	}
}

func TestEmail_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same value", "a@b.com", "a@b.com", true},
		{"different value", "a@b.com", "x@y.com", false},
		{"reflexive", "user@example.com", "user@example.com", true},
		{"case normalized", "USER@EXAMPLE.COM", "user@example.com", true},
		{"different local", "user+1@example.com", "user+2@example.com", false},
		{"different domain", "user@example.com", "user@sub.example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			emailA, err := NewEmail(tt.a)
			if err != nil {
				t.Fatalf("NewEmail(%q) failed: %v", tt.a, err)
			}
			emailB, err := NewEmail(tt.b)
			if err != nil {
				t.Fatalf("NewEmail(%q) failed: %v", tt.b, err)
			}
			if got := emailA.Equals(emailB); got != tt.want {
				t.Errorf("Email(%q).Equals(Email(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
