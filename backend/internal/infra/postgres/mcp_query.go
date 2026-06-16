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
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres/sqlc"
	mcpuc "github.com/karimiku/job-hunting-saas/internal/usecase/mcp"
)

// MCPQuery はPostgreSQLからMCP向けコンテキストを取得するQuery実装。
type MCPQuery struct {
	q *sqlc.Queries
}

// NewMCPQuery はMCP用の参照クエリを生成する。
func NewMCPQuery(db sqlc.DBTX) *MCPQuery {
	return &MCPQuery{q: sqlc.New(db)}
}

// ListEntries はユーザーのエントリー一覧をMCP向けDTOで返す。
func (q *MCPQuery) ListEntries(ctx context.Context, userID entity.UserID) ([]mcpuc.EntryDTO, error) {
	rows, err := q.q.MCPListEntries(ctx, uuid.UUID(userID))
	if err != nil {
		return nil, fmt.Errorf("postgres: MCPListEntries: %w", err)
	}
	out := make([]mcpuc.EntryDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, entryDTOFromMCPListRow(row))
	}
	return out, nil
}

// GetEntryContext は指定エントリーと関連タスクをMCP向けDTOで返す。
func (q *MCPQuery) GetEntryContext(ctx context.Context, userID entity.UserID, entryID entity.EntryID) (*mcpuc.EntryContextDTO, error) {
	row, err := q.q.MCPGetEntryContext(ctx, sqlc.MCPGetEntryContextParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(entryID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: MCPGetEntryContext: %w", err)
	}
	tasks, err := q.q.ListTasksByEntryID(ctx, sqlc.ListTasksByEntryIDParams{
		UserID:  uuid.UUID(userID),
		EntryID: uuid.UUID(entryID),
	})
	if err != nil {
		return nil, fmt.Errorf("postgres: ListTasksByEntryID: %w", err)
	}
	taskDTOs := make([]mcpuc.TaskDTO, 0, len(tasks))
	for _, task := range tasks {
		taskDTOs = append(taskDTOs, mcpuc.TaskDTO{
			ID:        task.ID.String(),
			EntryID:   task.EntryID.String(),
			Company:   row.CompanyName,
			Title:     task.Title,
			Type:      string(task.TaskType),
			DueDate:   formatPgTime(task.DueDate),
			Status:    string(task.Status),
			Notify:    task.Notify,
			Memo:      task.Memo,
			CreatedAt: formatPgTime(task.CreatedAt),
			UpdatedAt: formatPgTime(task.UpdatedAt),
		})
	}
	return &mcpuc.EntryContextDTO{
		Entry: mcpuc.EntryDTO{
			ID:         row.ID.String(),
			CompanyID:  row.CompanyID.String(),
			Company:    row.CompanyName,
			Route:      row.Route,
			Source:     row.Source,
			SourceURL:  row.SourceUrl,
			Status:     string(row.Status),
			StageKind:  string(row.StageKind),
			StageLabel: row.StageLabel,
			Memo:       row.Memo,
			CreatedAt:  formatPgTime(row.CreatedAt),
			UpdatedAt:  formatPgTime(row.UpdatedAt),
		},
		Tasks: taskDTOs,
	}, nil
}

// ListOpenTasks は未完了タスク一覧をMCP向けDTOで返す。
func (q *MCPQuery) ListOpenTasks(ctx context.Context, userID entity.UserID) ([]mcpuc.TaskDTO, error) {
	rows, err := q.q.MCPListOpenTasks(ctx, uuid.UUID(userID))
	if err != nil {
		return nil, fmt.Errorf("postgres: MCPListOpenTasks: %w", err)
	}
	out := make([]mcpuc.TaskDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, mcpuc.TaskDTO{
			ID:        row.ID.String(),
			EntryID:   row.EntryID.String(),
			Company:   row.CompanyName,
			Title:     row.Title,
			Type:      string(row.TaskType),
			DueDate:   formatPgTime(row.DueDate),
			Status:    string(row.Status),
			Notify:    row.Notify,
			Memo:      row.Memo,
			CreatedAt: formatPgTime(row.CreatedAt),
			UpdatedAt: formatPgTime(row.UpdatedAt),
		})
	}
	return out, nil
}

// ListInboxClips はユーザーのInbox Clip一覧をMCP向けDTOで返す。
func (q *MCPQuery) ListInboxClips(ctx context.Context, userID entity.UserID) ([]mcpuc.InboxClipDTO, error) {
	rows, err := q.q.ListInboxClipsByUserID(ctx, uuid.UUID(userID))
	if err != nil {
		return nil, fmt.Errorf("postgres: ListInboxClipsByUserID: %w", err)
	}
	out := make([]mcpuc.InboxClipDTO, 0, len(rows))
	for _, row := range rows {
		out = append(out, mcpuc.InboxClipDTO{
			ID:          row.ID.String(),
			URL:         row.Url,
			Title:       row.Title,
			Source:      row.Source,
			Guess:       row.Guess,
			ContentText: row.ContentText,
			CapturedAt:  formatPgTime(row.CapturedAt),
		})
	}
	return out, nil
}

func entryDTOFromMCPListRow(row sqlc.MCPListEntriesRow) mcpuc.EntryDTO {
	return mcpuc.EntryDTO{
		ID:         row.ID.String(),
		CompanyID:  row.CompanyID.String(),
		Company:    row.CompanyName,
		Route:      row.Route,
		Source:     row.Source,
		SourceURL:  row.SourceUrl,
		Status:     string(row.Status),
		StageKind:  string(row.StageKind),
		StageLabel: row.StageLabel,
		Memo:       row.Memo,
		CreatedAt:  formatPgTime(row.CreatedAt),
		UpdatedAt:  formatPgTime(row.UpdatedAt),
	}
}

func formatPgTime(ts pgtype.Timestamptz) *string {
	if !ts.Valid {
		return nil
	}
	value := ts.Time.Format(time.RFC3339)
	return &value
}
