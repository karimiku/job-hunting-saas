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

// TaskStatus はタスクの完了状態を表す値オブジェクト。
// todo（未完了）または done（完了）のいずれか。
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

// --- 定数コンストラクタ ---
// ハードコードされた既知の値に対して、エラーなしでインスタンスを返す。

func TaskStatusTodo() TaskStatus { return TaskStatus{value: taskStatusTodo} }
func TaskStatusDone() TaskStatus { return TaskStatus{value: taskStatusDone} }
