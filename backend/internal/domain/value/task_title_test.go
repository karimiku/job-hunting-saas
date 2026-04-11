package value

import (
	"errors"
	"testing"
)

func TestNewTaskTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"japanese title", "ES提出", nil},
		{"english title", "Submit Resume", nil},

		// 異常系
		{"empty", "", ErrTaskTitleEmpty},
		{"whitespace only", " ", ErrTaskTitleEmpty},
		{"leading space", " ES提出", ErrTaskTitleInvalid},
		{"trailing space", "ES提出 ", ErrTaskTitleInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTaskTitle(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTaskTitle(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestTaskTitle_String(t *testing.T) {
	title, err := NewTaskTitle("ES提出")
	if err != nil {
		t.Fatalf("NewTaskTitle failed: %v", err)
	}
	if title.String() != "ES提出" {
		t.Errorf("String() = %q, want %q", title.String(), "ES提出")
	}
}

func TestTaskTitle_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same value", "ES提出", "ES提出", true},
		{"different value", "ES提出", "一次面接", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			titleA, err := NewTaskTitle(tt.a)
			if err != nil {
				t.Fatalf("NewTaskTitle(%q) failed: %v", tt.a, err)
			}
			titleB, err := NewTaskTitle(tt.b)
			if err != nil {
				t.Fatalf("NewTaskTitle(%q) failed: %v", tt.b, err)
			}
			if got := titleA.Equals(titleB); got != tt.want {
				t.Errorf("TaskTitle(%q).Equals(TaskTitle(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
