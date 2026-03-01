package value

import (
	"errors"
	"testing"
)

func TestNewEntryStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"in progress", "in_progress", nil},
		{"offered", "offered", nil},
		{"accepted", "accepted", nil},
		{"rejected", "rejected", nil},
		{"withdrawn", "withdrawn", nil},

		// 異常系
		{"empty", "", ErrEntryStatusEmpty},
		{"unknown value", "pending", ErrEntryStatusInvalid},
		{"uppercase", "In_Progress", ErrEntryStatusInvalid},
		{"all caps", "OFFERED", ErrEntryStatusInvalid},
		{"with space", "in progress", ErrEntryStatusInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewEntryStatus(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewEntryStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestEntryStatus_String(t *testing.T) {
	status, err := NewEntryStatus("in_progress")
	if err != nil {
		t.Fatalf("NewEntryStatus failed: %v", err)
	}
	if status.String() != "in_progress" {
		t.Errorf("String() = %q, want %q", status.String(), "in_progress")
	}
}

func TestEntryStatus_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same value", "offered", "offered", true},
		{"different value", "offered", "rejected", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusA, err := NewEntryStatus(tt.a)
			if err != nil {
				t.Fatalf("NewEntryStatus(%q) failed: %v", tt.a, err)
			}
			statusB, err := NewEntryStatus(tt.b)
			if err != nil {
				t.Fatalf("NewEntryStatus(%q) failed: %v", tt.b, err)
			}
			if got := statusA.Equals(statusB); got != tt.want {
				t.Errorf("EntryStatus(%q).Equals(EntryStatus(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestEntryStatus_IsOpen(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"in_progress is open", "in_progress", true},
		{"offered is open", "offered", true},
		{"accepted is not open", "accepted", false},
		{"rejected is not open", "rejected", false},
		{"withdrawn is not open", "withdrawn", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := NewEntryStatus(tt.input)
			if err != nil {
				t.Fatalf("NewEntryStatus(%q) failed: %v", tt.input, err)
			}
			if got := status.IsOpen(); got != tt.want {
				t.Errorf("EntryStatus(%q).IsOpen() = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
