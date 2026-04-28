package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/infra/postgres/sqlc"
)

// InboxClipRepository は InboxClipRepository インターフェースの PostgreSQL 実装。
type InboxClipRepository struct {
	q *sqlc.Queries
}

// NewInboxClipRepository は InboxClipRepository を新規生成する。db には pgxpool.Pool もしくは tx を渡す。
func NewInboxClipRepository(db sqlc.DBTX) *InboxClipRepository {
	return &InboxClipRepository{q: sqlc.New(db)}
}

// Create はクリップを保存する。InboxClip は不変なので update は提供しない。
func (r *InboxClipRepository) Create(ctx context.Context, clip *entity.InboxClip) error {
	if err := r.q.CreateInboxClip(ctx, sqlc.CreateInboxClipParams{
		ID:         uuid.UUID(clip.ID()),
		UserID:     uuid.UUID(clip.UserID()),
		Url:        clip.URL().String(),
		Title:      clip.Title(),
		Source:     clip.Source().String(),
		Guess:      clip.Guess(),
		CapturedAt: pgtype.Timestamptz{Time: clip.CapturedAt(), Valid: true},
	}); err != nil {
		return fmt.Errorf("postgres: CreateInboxClip: %w", err)
	}
	return nil
}

// FindByID は userID 所有のクリップを ID から取得する。存在しないか他ユーザー所有の場合は repository.ErrNotFound を返す。
func (r *InboxClipRepository) FindByID(ctx context.Context, userID entity.UserID, id entity.InboxClipID) (*entity.InboxClip, error) {
	row, err := r.q.FindInboxClipByID(ctx, sqlc.FindInboxClipByIDParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(id),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindInboxClipByID: %w", err)
	}
	return reconstructInboxClip(row)
}

// ListByUserID は userID 所有のクリップを保存日時の新しい順で返す。
func (r *InboxClipRepository) ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.InboxClip, error) {
	rows, err := r.q.ListInboxClipsByUserID(ctx, uuid.UUID(userID))
	if err != nil {
		return nil, fmt.Errorf("postgres: ListInboxClipsByUserID: %w", err)
	}
	clips := make([]*entity.InboxClip, 0, len(rows))
	for _, row := range rows {
		c, err := reconstructInboxClip(row)
		if err != nil {
			return nil, err
		}
		clips = append(clips, c)
	}
	return clips, nil
}

// Delete は userID 所有のクリップを削除する。存在しないか他ユーザー所有の場合は repository.ErrNotFound を返す。
func (r *InboxClipRepository) Delete(ctx context.Context, userID entity.UserID, id entity.InboxClipID) error {
	n, err := r.q.DeleteInboxClip(ctx, sqlc.DeleteInboxClipParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(id),
	})
	if err != nil {
		return fmt.Errorf("postgres: DeleteInboxClip: %w", err)
	}
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func reconstructInboxClip(row sqlc.InboxClip) (*entity.InboxClip, error) {
	url, err := value.NewURL(row.Url)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: inbox_clip url: %w", err)
	}
	source, err := value.NewSource(row.Source)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: inbox_clip source: %w", err)
	}
	return entity.ReconstructInboxClip(
		entity.InboxClipID(row.ID),
		entity.UserID(row.UserID),
		url,
		row.Title,
		source,
		row.Guess,
		row.CapturedAt.Time,
	), nil
}
