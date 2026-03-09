package value

import (
	"errors"
	"testing"
)

func TestNewPassword(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"valid 8 chars", "password", nil},
		{"valid long", "this-is-a-very-long-password-123", nil},
		{"valid with symbols", "p@ssw0rd!", nil},
		{"valid japanese 8 chars", "パスワード１２３４", nil},

		// 異常系
		{"empty", "", ErrPasswordEmpty},
		{"7 chars", "passwor", ErrPasswordTooShort},
		{"1 char", "a", ErrPasswordTooShort},
		{"7 japanese chars", "パスワード１２", ErrPasswordTooShort},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPassword(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewPassword(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestPassword_Verify(t *testing.T) {
	pw, err := NewPassword("my-secret-password")
	if err != nil {
		t.Fatalf("NewPassword failed: %v", err)
	}

	if !pw.Verify("my-secret-password") {
		t.Error("Verify should return true for correct password")
	}
	if pw.Verify("wrong-password") {
		t.Error("Verify should return false for wrong password")
	}
	if pw.Verify("") {
		t.Error("Verify should return false for empty password")
	}
}

func TestPassword_Hash_NotPlaintext(t *testing.T) {
	raw := "my-secret-password"
	pw, err := NewPassword(raw)
	if err != nil {
		t.Fatalf("NewPassword failed: %v", err)
	}

	if pw.Hash() == raw {
		t.Error("Hash() should not return plaintext password")
	}
	if pw.Hash() == "" {
		t.Error("Hash() should not be empty")
	}
}

func TestReconstructPassword_Verify(t *testing.T) {
	original, err := NewPassword("my-secret-password")
	if err != nil {
		t.Fatalf("NewPassword failed: %v", err)
	}

	reconstructed := ReconstructPassword(original.Hash())

	if !reconstructed.Verify("my-secret-password") {
		t.Error("Reconstructed password should verify correctly")
	}
	if reconstructed.Verify("wrong-password") {
		t.Error("Reconstructed password should not verify wrong password")
	}
}
