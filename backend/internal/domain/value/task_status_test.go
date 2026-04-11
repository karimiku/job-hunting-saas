package value

import (
	"errors"
	"testing"
)

func TestNewTaskStatus(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		// 正常系
		{"todo", "todo", nil},
		{"done", "done", nil},

		// 異常系
		{"empty", "", ErrTaskStatusEmpty},
		{"unknown", "in_progress", ErrTaskStatusInvalid},
		{"uppercase", "Todo", ErrTaskStatusInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTaskStatus(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("NewTaskStatus(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestTaskStatus_String(t *testing.T) {
	s, err := NewTaskStatus("todo")
	if err != nil {
		t.Fatalf("NewTaskStatus failed: %v", err)
	}
	if s.String() != "todo" {
		t.Errorf("String() = %q, want %q", s.String(), "todo")
	}
}

func TestTaskStatus_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same value", "todo", "todo", true},
		{"different value", "todo", "done", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			statusA, err := NewTaskStatus(tt.a)
			if err != nil {
				t.Fatalf("NewTaskStatus(%q) failed: %v", tt.a, err)
			}
			statusB, err := NewTaskStatus(tt.b)
			if err != nil {
				t.Fatalf("NewTaskStatus(%q) failed: %v", tt.b, err)
			}
			if got := statusA.Equals(statusB); got != tt.want {
				t.Errorf("TaskStatus(%q).Equals(TaskStatus(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestTaskStatus_IsDone(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"todo is not done", "todo", false},
		{"done is done", "done", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := NewTaskStatus(tt.input)
			if err != nil {
				t.Fatalf("NewTaskStatus(%q) failed: %v", tt.input, err)
			}
			if got := status.IsDone(); got != tt.want {
				t.Errorf("TaskStatus(%q).IsDone() = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
