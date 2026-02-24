package value

import (
	"testing"
)

func TestNewTaskType(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		// 正常系
		{"deadline", "deadline", false},
		{"schedule", "schedule", false},

		// 異常系
		{"empty", "", true},
		{"unknown", "task", true},
		{"uppercase", "Deadline", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTaskType(tt.input)
			if tt.wantErr && err == nil {
				t.Errorf("NewTaskType(%q) should return error, but got nil", tt.input)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("NewTaskType(%q) should succeed, but got error: %v", tt.input, err)
			}
		})
	}
}

func TestTaskType_String(t *testing.T) {
	tt, err := NewTaskType("deadline")
	if err != nil {
		t.Fatalf("NewTaskType failed: %v", err)
	}
	if tt.String() != "deadline" {
		t.Errorf("String() = %q, want %q", tt.String(), "deadline")
	}
}

func TestTaskType_Equals(t *testing.T) {
	tests := []struct {
		name string
		a    string
		b    string
		want bool
	}{
		{"same value", "deadline", "deadline", true},
		{"different value", "deadline", "schedule", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			typeA, err := NewTaskType(tt.a)
			if err != nil {
				t.Fatalf("NewTaskType(%q) failed: %v", tt.a, err)
			}
			typeB, err := NewTaskType(tt.b)
			if err != nil {
				t.Fatalf("NewTaskType(%q) failed: %v", tt.b, err)
			}
			if got := typeA.Equals(typeB); got != tt.want {
				t.Errorf("TaskType(%q).Equals(TaskType(%q)) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestTaskType_IsSchedule(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"deadline is not schedule", "deadline", false},
		{"schedule is schedule", "schedule", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			taskType, err := NewTaskType(tt.input)
			if err != nil {
				t.Fatalf("NewTaskType(%q) failed: %v", tt.input, err)
			}
			if got := taskType.IsSchedule(); got != tt.want {
				t.Errorf("TaskType(%q).IsSchedule() = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}
