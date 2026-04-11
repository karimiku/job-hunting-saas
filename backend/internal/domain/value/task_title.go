package value

import (
	"errors"
	"strings"
)

var (
	ErrTaskTitleEmpty   = errors.New("task title must not be empty")
	ErrTaskTitleInvalid = errors.New("task title format is invalid")
)

// TaskTitle はタスクの件名を表す値オブジェクト。
type TaskTitle struct {
	value string
}

func NewTaskTitle(raw string) (TaskTitle, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return TaskTitle{}, ErrTaskTitleEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return TaskTitle{}, ErrTaskTitleInvalid
	}
	return TaskTitle{value: raw}, nil
}

func (t TaskTitle) String() string {
	return t.value
}

func (t TaskTitle) Equals(other TaskTitle) bool {
	return t.value == other.value
}
