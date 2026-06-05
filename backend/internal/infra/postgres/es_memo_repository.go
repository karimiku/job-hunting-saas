package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres/sqlc"
)

// ESMemoRepository はPostgreSQLを使ったESメモRepository実装。
type ESMemoRepository struct {
	q *sqlc.Queries
}

// NewESMemoRepository はESメモRepositoryを生成する。
func NewESMemoRepository(db sqlc.DBTX) *ESMemoRepository {
	return &ESMemoRepository{q: sqlc.New(db)}
}

// Save はESメモを永続化する。
func (r *ESMemoRepository) Save(ctx context.Context, memo *entity.ESMemo) error {
	var entryID pgtype.UUID
	if memo.EntryID() != nil {
		entryID = pgtype.UUID{Bytes: uuid.UUID(*memo.EntryID()), Valid: true}
	}
	if _, err := r.q.CreateESMemo(ctx, sqlc.CreateESMemoParams{
		ID:        uuid.UUID(memo.ID()),
		UserID:    uuid.UUID(memo.UserID()),
		EntryID:   entryID,
		Category:  memo.Category().String(),
		Title:     memo.Title().String(),
		Content:   memo.Content().String(),
		Source:    memo.Source().String(),
		CreatedAt: pgtype.Timestamptz{Time: memo.CreatedAt(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: memo.UpdatedAt(), Valid: true},
	}); err != nil {
		return fmt.Errorf("postgres: CreateESMemo: %w", err)
	}
	return nil
}

// ListByUserID はユーザーのESメモを新しい順に取得する。
func (r *ESMemoRepository) ListByUserID(ctx context.Context, userID entity.UserID, limit int32) ([]*entity.ESMemo, error) {
	rows, err := r.q.ListESMemosByUserID(ctx, sqlc.ListESMemosByUserIDParams{
		UserID: uuid.UUID(userID),
		Limit:  limit,
	})
	if err != nil {
		return nil, fmt.Errorf("postgres: ListESMemosByUserID: %w", err)
	}
	memos := make([]*entity.ESMemo, 0, len(rows))
	for _, row := range rows {
		memo, err := reconstructESMemo(row)
		if err != nil {
			return nil, err
		}
		memos = append(memos, memo)
	}
	return memos, nil
}

func reconstructESMemo(row sqlc.EsMemo) (*entity.ESMemo, error) {
	category, err := value.NewESMemoCategory(row.Category)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: es memo category: %w", err)
	}
	title, err := value.NewESMemoTitle(row.Title)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: es memo title: %w", err)
	}
	content, err := value.NewESMemoContent(row.Content)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: es memo content: %w", err)
	}
	source, err := value.NewESMemoSource(row.Source)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: es memo source: %w", err)
	}
	var entryID *entity.EntryID
	if row.EntryID.Valid {
		id := entity.EntryID(row.EntryID.Bytes)
		entryID = &id
	}
	return entity.ReconstructESMemo(
		entity.ESMemoID(row.ID),
		entity.UserID(row.UserID),
		entryID,
		category,
		title,
		content,
		source,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	), nil
}
