package value

import "errors"

var (
	ErrTaskStatusEmpty   = errors.New("task status must not be empty")
	ErrTaskStatusInvalid = errors.New("task status is invalid")
)

const (
	taskStatusTodo = "todo"
	taskStatusDone = "done"
)

var validTaskStatuses = map[string]bool{
	taskStatusTodo: true,
	taskStatusDone: true,
}

type TaskStatus struct {
	value string
}

func NewTaskStatus(raw string) (TaskStatus, error) {
	if raw == "" {
		return TaskStatus{}, ErrTaskStatusEmpty
	}
	if !validTaskStatuses[raw] {
		return TaskStatus{}, ErrTaskStatusInvalid
	}
	return TaskStatus{value: raw}, nil
}

func (s TaskStatus) String() string {
	return s.value
}

func (s TaskStatus) Equals(other TaskStatus) bool {
	return s.value == other.value
}

func (s TaskStatus) IsDone() bool {
	return s.value == taskStatusDone
}
