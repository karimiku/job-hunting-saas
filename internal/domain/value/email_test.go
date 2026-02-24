package value

import (
	"testing"
)

func TestNewEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// 正常系
		{"basic email", "user@example.com", false},
		{"subdomain", "test.name@company.co.jp", false},
		{"minimal", "a@b.co", false},
		{"underscore", "user_name@example.com", false},
		{"hyphen", "user-name@example.com", false},
		{"plus tag", "user+tag@example.com", false},
		{"uppercase input", "USER@EXAMPLE.COM", false},

		// 異常系
		{"empty", "", true},
		{"no at sign", "invalid", true},
		{"no local part", "@example.com", true},
		{"no domain", "user@", true},
		{"domain starts with dot", "user@.com", true},
		{"no tld", "user@example", true},
		{"double at sign", "user@@example.com", true},
		{"leading dot local", ".user@example.com", true},
		{"trailing dot local", "user.@example.com", true},
		{"consecutive dots local", "us..er@example.com", true},
		{"consecutive dots domain", "user@example..com", true},
		{"leading space", " user@example.com", true},
		{"trailing space", "user@example.com ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEmail(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("NewEmail(%q) should return error, but got nil", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("NewEmail(%q) should succeed, but got error: %v", tt.input, err)
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
