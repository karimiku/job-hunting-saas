package value

import "errors"

// ErrTaskStatusEmpty は task status が空文字のときに返されるエラー。
// ErrTaskStatusInvalid は task status が未定義の値のときに返されるエラー。
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

// TaskStatus はタスクの完了状態を表す値オブジェクト。
// todo（未完了）または done（完了）のいずれか。
type TaskStatus struct {
	value string
}

// NewTaskStatus は raw から TaskStatus を生成する。空文字や未定義値は対応するエラーを返す。
func NewTaskStatus(raw string) (TaskStatus, error) {
	if raw == "" {
		return TaskStatus{}, ErrTaskStatusEmpty
	}
	if !validTaskStatuses[raw] {
		return TaskStatus{}, ErrTaskStatusInvalid
	}
	return TaskStatus{value: raw}, nil
}

// String は task status を文字列で返す。
func (s TaskStatus) String() string {
	return s.value
}

// Equals は 2 つの TaskStatus が等しいかを判定する。
func (s TaskStatus) Equals(other TaskStatus) bool {
	return s.value == other.value
}

// IsDone は タスクが完了 (done) しているかを返す。
func (s TaskStatus) IsDone() bool {
	return s.value == taskStatusDone
}

// --- 定数コンストラクタ ---
// ハードコードされた既知の値に対して、エラーなしでインスタンスを返す。

// TaskStatusTodo は todo 状態の TaskStatus を返す。
func TaskStatusTodo() TaskStatus { return TaskStatus{value: taskStatusTodo} }

// TaskStatusDone は done 状態の TaskStatus を返す。
func TaskStatusDone() TaskStatus { return TaskStatus{value: taskStatusDone} }
