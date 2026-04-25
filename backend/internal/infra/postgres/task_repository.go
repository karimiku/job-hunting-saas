package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres/sqlc"
)

// TaskRepository は TaskRepository インターフェースの PostgreSQL 実装。
type TaskRepository struct {
	q *sqlc.Queries
}

// NewTaskRepository は TaskRepository を新規生成する。db には pgxpool.Pool もしくは tx を渡す。
func NewTaskRepository(db sqlc.DBTX) *TaskRepository {
	return &TaskRepository{q: sqlc.New(db)}
}

// Save は Task を upsert する。同じ ID があれば更新、なければ作成。
func (r *TaskRepository) Save(ctx context.Context, task *entity.Task) error {
	var dueDate pgtype.Timestamptz
	if task.DueDate() != nil {
		dueDate = pgtype.Timestamptz{Time: *task.DueDate(), Valid: true}
	}

	if err := r.q.UpsertTask(ctx, sqlc.UpsertTaskParams{
		ID:        uuid.UUID(task.ID()),
		EntryID:   uuid.UUID(task.EntryID()),
		Title:     task.Title().String(),
		TaskType:  sqlc.TaskType(task.TaskType().String()),
		DueDate:   dueDate,
		Status:    sqlc.TaskStatus(task.Status().String()),
		Notify:    task.Notify(),
		Memo:      task.Memo(),
		CreatedAt: pgtype.Timestamptz{Time: task.CreatedAt(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: task.UpdatedAt(), Valid: true},
	}); err != nil {
		return fmt.Errorf("postgres: UpsertTask: %w", err)
	}
	return nil
}

// FindByID は userID 所有の Task を ID から取得する。SQL で Entry の userID を JOIN 検証し、未所有なら repository.ErrNotFound を返す。
func (r *TaskRepository) FindByID(ctx context.Context, userID entity.UserID, id entity.TaskID) (*entity.Task, error) {
	row, err := r.q.FindTaskByID(ctx, sqlc.FindTaskByIDParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(id),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindTaskByID: %w", err)
	}

	return reconstructTask(row)
}

// ListByEntryID は entry に紐づく Task を全件返す。SQL で Entry の userID を JOIN 検証する。
func (r *TaskRepository) ListByEntryID(ctx context.Context, userID entity.UserID, entryID entity.EntryID) ([]*entity.Task, error) {
	rows, err := r.q.ListTasksByEntryID(ctx, sqlc.ListTasksByEntryIDParams{
		UserID:  uuid.UUID(userID),
		EntryID: uuid.UUID(entryID),
	})
	if err != nil {
		return nil, fmt.Errorf("postgres: ListTasksByEntryID: %w", err)
	}

	return reconstructTasks(rows)
}

// ListByUserIDWithDueBefore は userID 所有かつ deadline より前が期限の未完了 Task を返す。リマインダ通知用。
func (r *TaskRepository) ListByUserIDWithDueBefore(ctx context.Context, userID entity.UserID, deadline time.Time) ([]*entity.Task, error) {
	rows, err := r.q.ListTasksByUserIDWithDueBefore(ctx, sqlc.ListTasksByUserIDWithDueBeforeParams{
		UserID:  uuid.UUID(userID),
		DueDate: pgtype.Timestamptz{Time: deadline, Valid: true},
	})
	if err != nil {
		return nil, fmt.Errorf("postgres: ListTasksByUserIDWithDueBefore: %w", err)
	}

	return reconstructTasks(rows)
}

// Delete は userID 所有の Task を ID から削除する。SQL で Entry の userID を JOIN 検証し、未所有なら repository.ErrNotFound を返す。
func (r *TaskRepository) Delete(ctx context.Context, userID entity.UserID, id entity.TaskID) error {
	n, err := r.q.DeleteTask(ctx, sqlc.DeleteTaskParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(id),
	})
	if err != nil {
		return fmt.Errorf("postgres: DeleteTask: %w", err)
	}
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func reconstructTask(row sqlc.Task) (*entity.Task, error) {
	title, err := value.NewTaskTitle(row.Title)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: task title: %w", err)
	}

	taskType, err := value.NewTaskType(string(row.TaskType))
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: task type: %w", err)
	}

	status, err := value.NewTaskStatus(string(row.Status))
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: task status: %w", err)
	}

	var dueDate *time.Time
	if row.DueDate.Valid {
		t := row.DueDate.Time
		dueDate = &t
	}

	return entity.ReconstructTask(
		entity.TaskID(row.ID),
		entity.EntryID(row.EntryID),
		title,
		taskType,
		dueDate,
		status,
		row.Notify,
		row.Memo,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	), nil
}

func reconstructTasks(rows []sqlc.Task) ([]*entity.Task, error) {
	tasks := make([]*entity.Task, 0, len(rows))
	for _, row := range rows {
		t, err := reconstructTask(row)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
