package value

import (
	"errors"
	"strings"
	"testing"
)

func TestNewAIAccessTokenSecret(t *testing.T) {
	t.Parallel()

	validBody := strings.Repeat("A", minAIAccessTokenBody)
	tests := []struct {
		name    string
		raw     string
		wantErr error
	}{
		{"valid", AIAccessTokenPrefix + validBody, nil},
		{"trim", " " + AIAccessTokenPrefix + validBody + "\n", nil},
		{"empty", "", ErrAIAccessTokenEmpty},
		{"wrong prefix", "sk_" + validBody, ErrAIAccessTokenInvalid},
		{"short body", AIAccessTokenPrefix + "abc", ErrAIAccessTokenInvalid},
		{"invalid body char", AIAccessTokenPrefix + strings.Repeat("A", minAIAccessTokenBody-1) + "!", ErrAIAccessTokenInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewAIAccessTokenSecret(tt.raw)
			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want %v", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("NewAIAccessTokenSecret() failed: %v", err)
			}
			if !strings.HasPrefix(got.String(), AIAccessTokenPrefix) {
				t.Errorf("String() = %q, want prefix %q", got.String(), AIAccessTokenPrefix)
			}
			if len(got.Hash()) != 64 {
				t.Errorf("Hash() length = %d, want 64", len(got.Hash()))
			}
		})
	}
}

func TestGenerateAIAccessTokenSecret(t *testing.T) {
	t.Parallel()

	token, err := GenerateAIAccessTokenSecret()
	if err != nil {
		t.Fatalf("GenerateAIAccessTokenSecret() failed: %v", err)
	}
	if !strings.HasPrefix(token.String(), AIAccessTokenPrefix) {
		t.Errorf("generated token = %q, want prefix %q", token.String(), AIAccessTokenPrefix)
	}
	if token.Preview() == token.String() {
		t.Error("Preview() should not expose the full token")
	}
}
