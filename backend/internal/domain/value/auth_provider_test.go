package value

import (
	"errors"
	"testing"
)

func TestNewAuthProvider(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"google", "google", nil},
		{"supabase", "supabase", nil},
		{"empty", "", ErrAuthProviderEmpty},
		{"invalid", "facebook", ErrAuthProviderInvalid},
		{"uppercase", "Google", ErrAuthProviderInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewAuthProvider(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewAuthProvider(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestAuthProvider_Equals(t *testing.T) {
	a := AuthProviderGoogle()
	b := AuthProviderGoogle()
	c := AuthProviderSupabase()

	if !a.Equals(b) {
		t.Error("same providers should be equal")
	}
	if a.Equals(c) {
		t.Error("different providers should not be equal")
	}
}

func TestAuthProvider_String(t *testing.T) {
	p := AuthProviderGoogle()
	if p.String() != "google" {
		t.Errorf("String() = %q, want %q", p.String(), "google")
	}
	if got := AuthProviderSupabase().String(); got != "supabase" {
		t.Errorf("String() = %q, want %q", got, "supabase")
	}
}
