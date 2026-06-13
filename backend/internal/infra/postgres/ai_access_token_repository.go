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
)

// AIAccessTokenRepository は AIAccessTokenRepository インターフェースの PostgreSQL 実装。
type AIAccessTokenRepository struct {
	q *sqlc.Queries
}

// NewAIAccessTokenRepository は AIAccessTokenRepository を新規生成する。
func NewAIAccessTokenRepository(db sqlc.DBTX) *AIAccessTokenRepository {
	return &AIAccessTokenRepository{q: sqlc.New(db)}
}

// Save はAIアクセストークンを保存する。
func (r *AIAccessTokenRepository) Save(ctx context.Context, token *entity.AIAccessToken) error {
	if err := r.q.CreateAIAccessToken(ctx, sqlc.CreateAIAccessTokenParams{
		ID:           uuid.UUID(token.ID()),
		UserID:       uuid.UUID(token.UserID()),
		Name:         token.Name(),
		TokenHash:    token.TokenHash(),
		TokenPreview: token.TokenPreview(),
		LastUsedAt:   nullableTime(token.LastUsedAt()),
		RevokedAt:    nullableTime(token.RevokedAt()),
		CreatedAt:    pgtype.Timestamptz{Time: token.CreatedAt(), Valid: true},
		UpdatedAt:    pgtype.Timestamptz{Time: token.UpdatedAt(), Valid: true},
	}); err != nil {
		if isUniqueViolation(err) {
			return repository.ErrAlreadyExists
		}
		return fmt.Errorf("postgres: CreateAIAccessToken: %w", err)
	}
	return nil
}

// FindActiveByHash は有効なAIアクセストークンをハッシュから取得する。
func (r *AIAccessTokenRepository) FindActiveByHash(ctx context.Context, tokenHash string) (*entity.AIAccessToken, error) {
	row, err := r.q.FindActiveAIAccessTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindActiveAIAccessTokenByHash: %w", err)
	}
	return reconstructAIAccessToken(row), nil
}

// ListByUserID は指定ユーザーのAIアクセストークンを新しい順に返す。
func (r *AIAccessTokenRepository) ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.AIAccessToken, error) {
	rows, err := r.q.ListAIAccessTokensByUserID(ctx, uuid.UUID(userID))
	if err != nil {
		return nil, fmt.Errorf("postgres: ListAIAccessTokensByUserID: %w", err)
	}
	out := make([]*entity.AIAccessToken, 0, len(rows))
	for _, row := range rows {
		out = append(out, reconstructAIAccessToken(row))
	}
	return out, nil
}

// TouchLastUsed はAIアクセストークンの最終利用日時を更新する。
func (r *AIAccessTokenRepository) TouchLastUsed(ctx context.Context, id entity.AIAccessTokenID, usedAt time.Time) error {
	n, err := r.q.TouchAIAccessTokenLastUsed(ctx, sqlc.TouchAIAccessTokenLastUsedParams{
		ID:         uuid.UUID(id),
		LastUsedAt: pgtype.Timestamptz{Time: usedAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("postgres: TouchAIAccessTokenLastUsed: %w", err)
	}
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

// Revoke はAIアクセストークンを失効する。
func (r *AIAccessTokenRepository) Revoke(ctx context.Context, userID entity.UserID, id entity.AIAccessTokenID, revokedAt time.Time) error {
	n, err := r.q.RevokeAIAccessToken(ctx, sqlc.RevokeAIAccessTokenParams{
		UserID:    uuid.UUID(userID),
		ID:        uuid.UUID(id),
		RevokedAt: pgtype.Timestamptz{Time: revokedAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("postgres: RevokeAIAccessToken: %w", err)
	}
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func reconstructAIAccessToken(row sqlc.AiAccessToken) *entity.AIAccessToken {
	return entity.ReconstructAIAccessToken(
		entity.AIAccessTokenID(row.ID),
		entity.UserID(row.UserID),
		row.Name,
		row.TokenHash,
		row.TokenPreview,
		timePtr(row.LastUsedAt),
		timePtr(row.RevokedAt),
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	)
}

func nullableTime(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func timePtr(ts pgtype.Timestamptz) *time.Time {
	if !ts.Valid {
		return nil
	}
	return &ts.Time
}
