package value

import (
	"strings"
	"testing"
)

func TestNewESMemoCategory(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"default", "", "general", false},
		{"trim", " interview ", "interview", false},
		{"too long", strings.Repeat("a", maxESMemoCategoryLen+1), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewESMemoCategory(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewESMemoCategory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("String() = %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestNewESMemoTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"valid", "  面接で話した改善経験  ", "面接で話した改善経験", false},
		{"empty", "", "", true},
		{"whitespace", "   ", "", true},
		{"too long", strings.Repeat("a", maxESMemoTitleLen+1), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewESMemoTitle(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewESMemoTitle() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("String() = %q, want %q", got.String(), tt.want)
			}
		})
	}
}

func TestNewESMemoContent(t *testing.T) {
	got, err := NewESMemoContent("  顧客課題を分解して改善した  ")
	if err != nil {
		t.Fatalf("NewESMemoContent() failed: %v", err)
	}
	if got.String() != "顧客課題を分解して改善した" {
		t.Errorf("String() = %q", got.String())
	}

	if _, err := NewESMemoContent(" "); err == nil {
		t.Fatal("NewESMemoContent() error = nil, want error")
	}
}

func TestNewESMemoSource(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"default", "", "mcp", false},
		{"trim", " mail ", "mail", false},
		{"too long", strings.Repeat("a", maxESMemoSourceLen+1), "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewESMemoSource(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("NewESMemoSource() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && got.String() != tt.want {
				t.Errorf("String() = %q, want %q", got.String(), tt.want)
			}
		})
	}
}
