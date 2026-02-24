package value

import "errors"

var (
	ErrTaskTypeEmpty   = errors.New("task type must not be empty")
	ErrTaskTypeInvalid = errors.New("task type is invalid")
)

const (
	taskTypeDeadline = "deadline"
	taskTypeSchedule = "schedule"
)

var validTaskTypes = map[string]bool{
	taskTypeDeadline: true,
	taskTypeSchedule: true,
}

type TaskType struct {
	value string
}

func NewTaskType(raw string) (TaskType, error) {
	if raw == "" {
		return TaskType{}, ErrTaskTypeEmpty
	}
	if !validTaskTypes[raw] {
		return TaskType{}, ErrTaskTypeInvalid
	}
	return TaskType{value: raw}, nil
}

func (t TaskType) String() string {
	return t.value
}

func (t TaskType) Equals(other TaskType) bool {
	return t.value == other.value
}

func (t TaskType) IsSchedule() bool {
	return t.value == taskTypeSchedule
}
