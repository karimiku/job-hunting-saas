package value

import (
	"errors"
	"strings"
	"testing"
)

func TestNewInboxClipTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid", "株式会社サンプル 募集要項", nil},
		{"max length", strings.Repeat("あ", InboxClipTitleMaxLength), nil},
		{"empty", "", ErrInboxClipTitleEmpty},
		{"whitespace only", " ", ErrInboxClipTitleEmpty},
		{"leading space", " 株式会社サンプル", ErrInboxClipTitleInvalid},
		{"trailing space", "株式会社サンプル ", ErrInboxClipTitleInvalid},
		{"too long", strings.Repeat("あ", InboxClipTitleMaxLength+1), ErrInboxClipTitleTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			title, err := NewInboxClipTitle(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewInboxClipTitle(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if tt.wantErr == nil && title.String() != tt.input {
				t.Errorf("String() = %q, want %q", title.String(), tt.input)
			}
		})
	}
}

func TestInboxClipTitle_Equals(t *testing.T) {
	a, err := NewInboxClipTitle("株式会社サンプル")
	if err != nil {
		t.Fatalf("NewInboxClipTitle: %v", err)
	}
	b, err := NewInboxClipTitle("株式会社サンプル")
	if err != nil {
		t.Fatalf("NewInboxClipTitle: %v", err)
	}
	c, err := NewInboxClipTitle("別タイトル")
	if err != nil {
		t.Fatalf("NewInboxClipTitle: %v", err)
	}

	if !a.Equals(b) {
		t.Error("same title should be equal")
	}
	if a.Equals(c) {
		t.Error("different title should not be equal")
	}
}

func TestNewInboxClipGuess(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{"valid", "株式会社サンプル", nil},
		{"empty allowed", "", nil},
		{"max length", strings.Repeat("あ", InboxClipGuessMaxLength), nil},
		{"leading space", " 株式会社サンプル", ErrInboxClipGuessInvalid},
		{"trailing space", "株式会社サンプル ", ErrInboxClipGuessInvalid},
		{"too long", strings.Repeat("あ", InboxClipGuessMaxLength+1), ErrInboxClipGuessTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			guess, err := NewInboxClipGuess(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("NewInboxClipGuess(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			if tt.wantErr == nil && guess.String() != tt.input {
				t.Errorf("String() = %q, want %q", guess.String(), tt.input)
			}
		})
	}
}

func TestInboxClipGuess_Equals(t *testing.T) {
	a, err := NewInboxClipGuess("株式会社サンプル")
	if err != nil {
		t.Fatalf("NewInboxClipGuess: %v", err)
	}
	b, err := NewInboxClipGuess("株式会社サンプル")
	if err != nil {
		t.Fatalf("NewInboxClipGuess: %v", err)
	}
	c, err := NewInboxClipGuess("")
	if err != nil {
		t.Fatalf("NewInboxClipGuess: %v", err)
	}

	if !a.Equals(b) {
		t.Error("same guess should be equal")
	}
	if a.Equals(c) {
		t.Error("different guess should not be equal")
	}
}
