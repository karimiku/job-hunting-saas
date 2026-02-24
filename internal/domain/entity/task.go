package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

var (
	ErrTaskTitleEmpty = errors.New("task title must not be empty")
)

type Task struct {
	id        TaskID
	entryID   EntryID
	title     string
	taskType  value.TaskType
	dueDate   *time.Time
	status    value.TaskStatus
	notify    bool
	memo      string
	createdAt time.Time
	updatedAt time.Time
}

func NewTask(entryID EntryID, title string, taskType value.TaskType) (*Task, error) {
	if strings.TrimSpace(title) == "" {
		return nil, ErrTaskTitleEmpty
	}

	status, _ := value.NewTaskStatus("todo")

	now := time.Now()
	return &Task{
		id:        NewID(),
		entryID:   entryID,
		title:     title,
		taskType:  taskType,
		dueDate:   nil,
		status:    status,
		notify:    false,
		memo:      "",
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (t *Task) ID() TaskID              { return t.id }
func (t *Task) EntryID() EntryID        { return t.entryID }
func (t *Task) Title() string           { return t.title }
func (t *Task) TaskType() value.TaskType   { return t.taskType }
func (t *Task) DueDate() *time.Time     { return t.dueDate }
func (t *Task) Status() value.TaskStatus   { return t.status }
func (t *Task) Notify() bool            { return t.notify }
func (t *Task) Memo() string            { return t.memo }
func (t *Task) CreatedAt() time.Time    { return t.createdAt }
func (t *Task) UpdatedAt() time.Time    { return t.updatedAt }

func (t *Task) Complete() {
	t.status, _ = value.NewTaskStatus("done")
	t.updatedAt = time.Now()
}

func (t *Task) Uncomplete() {
	t.status, _ = value.NewTaskStatus("todo")
	t.updatedAt = time.Now()
}

func (t *Task) SetDueDate(dueDate time.Time) {
	t.dueDate = &dueDate
	t.updatedAt = time.Now()
}

func (t *Task) ClearDueDate() {
	t.dueDate = nil
	t.updatedAt = time.Now()
}

func (t *Task) SetNotify(notify bool) {
	t.notify = notify
	t.updatedAt = time.Now()
}

func (t *Task) UpdateMemo(memo string) {
	t.memo = memo
	t.updatedAt = time.Now()
}
