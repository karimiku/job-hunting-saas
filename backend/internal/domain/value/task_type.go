package value

import "errors"

// ErrTaskTypeEmpty は task type が空文字のときに返されるエラー。
// ErrTaskTypeInvalid は task type が未定義の値のときに返されるエラー。
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

// TaskType はタスクの種別を表す値オブジェクト。
// deadline（締切までにやる作業）または schedule（日時確定の予定）のいずれか。
type TaskType struct {
	value string
}

// NewTaskType は raw から TaskType を生成する。空文字や未定義値は対応するエラーを返す。
func NewTaskType(raw string) (TaskType, error) {
	if raw == "" {
		return TaskType{}, ErrTaskTypeEmpty
	}
	if !validTaskTypes[raw] {
		return TaskType{}, ErrTaskTypeInvalid
	}
	return TaskType{value: raw}, nil
}

// String は task type を文字列で返す。
func (t TaskType) String() string {
	return t.value
}

// Equals は 2 つの TaskType が等しいかを判定する。
func (t TaskType) Equals(other TaskType) bool {
	return t.value == other.value
}

// IsSchedule は タスク種別が schedule (日時確定の予定) かを返す。
func (t TaskType) IsSchedule() bool {
	return t.value == taskTypeSchedule
}
