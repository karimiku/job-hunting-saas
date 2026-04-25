package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// Task は Entry に紐づく作業や予定を表すエンティティ。
// 締切管理（deadline）と日程管理（schedule）の両方を統合的に扱う。
type Task struct {
	id        TaskID
	entryID   EntryID
	title     value.TaskTitle
	taskType  value.TaskType
	dueDate   *time.Time
	status    value.TaskStatus
	notify    bool
	memo      string
	createdAt time.Time
	updatedAt time.Time
}

// NewTask は Task を新規作成する。各値オブジェクトのバリデーションは呼び出し側で済んでいる前提。
func NewTask(entryID EntryID, title value.TaskTitle, taskType value.TaskType) *Task {
	now := time.Now()
	return &Task{
		id:        NewTaskID(),
		entryID:   entryID,
		title:     title,
		taskType:  taskType,
		dueDate:   nil,
		status:    value.TaskStatusTodo(), // 定数コンストラクタでエラー握りつぶしを回避
		notify:    false,
		memo:      "",
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructTask はDBから読み取ったデータでTaskを復元する。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructTask(
	id TaskID, entryID EntryID, title value.TaskTitle, taskType value.TaskType,
	dueDate *time.Time, status value.TaskStatus, notify bool, memo string,
	createdAt, updatedAt time.Time,
) *Task {
	return &Task{
		id:        id,
		entryID:   entryID,
		title:     title,
		taskType:  taskType,
		dueDate:   dueDate,
		status:    status,
		notify:    notify,
		memo:      memo,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID は Task の ID を返す。
func (t *Task) ID() TaskID { return t.id }

// EntryID は Task が紐づく Entry の ID を返す。
func (t *Task) EntryID() EntryID { return t.entryID }

// Title は Task のタイトルを返す。
func (t *Task) Title() value.TaskTitle { return t.title }

// TaskType は Task の種別を返す。
func (t *Task) TaskType() value.TaskType { return t.taskType }

// DueDate は Task の期日を返す。未設定の場合は nil。
func (t *Task) DueDate() *time.Time { return t.dueDate }

// Status は Task のステータスを返す。
func (t *Task) Status() value.TaskStatus { return t.status }

// Notify は Task の通知設定を返す。
func (t *Task) Notify() bool { return t.notify }

// Memo は Task のメモを返す。
func (t *Task) Memo() string { return t.memo }

// CreatedAt は Task の作成日時を返す。
func (t *Task) CreatedAt() time.Time { return t.createdAt }

// UpdatedAt は Task の更新日時を返す。
func (t *Task) UpdatedAt() time.Time { return t.updatedAt }

// UpdateTitle はタイトルを更新し、UpdatedAt を現在時刻にする。
func (t *Task) UpdateTitle(title value.TaskTitle) {
	t.title = title
	t.updatedAt = time.Now()
}

// UpdateTaskType は Task 種別を更新し、UpdatedAt を現在時刻にする。
func (t *Task) UpdateTaskType(taskType value.TaskType) {
	t.taskType = taskType
	t.updatedAt = time.Now()
}

// Complete は Task を完了状態にし、UpdatedAt を現在時刻にする。
func (t *Task) Complete() {
	t.status = value.TaskStatusDone()
	t.updatedAt = time.Now()
}

// Uncomplete は Task を未完了状態に戻し、UpdatedAt を現在時刻にする。
func (t *Task) Uncomplete() {
	t.status = value.TaskStatusTodo()
	t.updatedAt = time.Now()
}

// SetDueDate は期日を設定し、UpdatedAt を現在時刻にする。
func (t *Task) SetDueDate(dueDate time.Time) {
	t.dueDate = &dueDate
	t.updatedAt = time.Now()
}

// ClearDueDate は期日をクリアし、UpdatedAt を現在時刻にする。
func (t *Task) ClearDueDate() {
	t.dueDate = nil
	t.updatedAt = time.Now()
}

// SetNotify は通知設定を更新し、UpdatedAt を現在時刻にする。
func (t *Task) SetNotify(notify bool) {
	t.notify = notify
	t.updatedAt = time.Now()
}

// UpdateMemo は memo を更新し、UpdatedAt を現在時刻にする。
func (t *Task) UpdateMemo(memo string) {
	t.memo = memo
	t.updatedAt = time.Now()
}
