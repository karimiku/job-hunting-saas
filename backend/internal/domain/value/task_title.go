package value

import (
	"errors"
	"strings"
)

// ErrTaskTitleEmpty は task title が空文字のときに返されるエラー。
// ErrTaskTitleInvalid は task title の形式が不正なときに返されるエラー。
var (
	ErrTaskTitleEmpty   = errors.New("task title must not be empty")
	ErrTaskTitleInvalid = errors.New("task title format is invalid")
)

// TaskTitle はタスクの件名を表す値オブジェクト。
type TaskTitle struct {
	value string
}

// NewTaskTitle は raw から TaskTitle を生成する。空文字や不正値は対応するエラーを返す。
func NewTaskTitle(raw string) (TaskTitle, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return TaskTitle{}, ErrTaskTitleEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return TaskTitle{}, ErrTaskTitleInvalid
	}
	return TaskTitle{value: raw}, nil
}

// String は task title を文字列で返す。
func (t TaskTitle) String() string {
	return t.value
}

// Equals は 2 つの TaskTitle が等しいかを判定する。
func (t TaskTitle) Equals(other TaskTitle) bool {
	return t.value == other.value
}
