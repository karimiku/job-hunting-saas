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

// AIAccessTokenRepository は AIAccessTokenRepository インターフェースの PostgreSQL 実装。
type AIAccessTokenRepository struct {
	q *sqlc.Queries
}

// NewAIAccessTokenRepository は AIAccessTokenRepository を新規生成する。
func NewAIAccessTokenRepository(db sqlc.DBTX) *AIAccessTokenRepository {
	return &AIAccessTokenRepository{q: sqlc.New(db)}
}

// Create はトークンを保存する。平文トークンは保存しない。
func (r *AIAccessTokenRepository) Create(ctx context.Context, token *entity.AIAccessToken) error {
	if err := r.q.CreateAIAccessToken(ctx, sqlc.CreateAIAccessTokenParams{
		ID:          uuid.UUID(token.ID()),
		UserID:      uuid.UUID(token.UserID()),
		Name:        token.Name().String(),
		TokenHash:   token.TokenHash().String(),
		TokenPrefix: token.Prefix().String(),
		CreatedAt:   pgtype.Timestamptz{Time: token.CreatedAt(), Valid: true},
		LastUsedAt:  timestamptzFromPtr(token.LastUsedAt()),
		RevokedAt:   timestamptzFromPtr(token.RevokedAt()),
	}); err != nil {
		return fmt.Errorf("postgres: CreateAIAccessToken: %w", err)
	}
	return nil
}

// ListByUserID は userID 所有のトークンを作成日時の新しい順で返す。
func (r *AIAccessTokenRepository) ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.AIAccessToken, error) {
	rows, err := r.q.ListAIAccessTokensByUserID(ctx, uuid.UUID(userID))
	if err != nil {
		return nil, fmt.Errorf("postgres: ListAIAccessTokensByUserID: %w", err)
	}
	tokens := make([]*entity.AIAccessToken, 0, len(rows))
	for _, row := range rows {
		t, err := reconstructAIAccessToken(row)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, t)
	}
	return tokens, nil
}

// FindActiveByHash は未失効トークンを hash から取得する。
func (r *AIAccessTokenRepository) FindActiveByHash(ctx context.Context, hash value.AIAccessTokenHash) (*entity.AIAccessToken, error) {
	row, err := r.q.FindActiveAIAccessTokenByHash(ctx, hash.String())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindActiveAIAccessTokenByHash: %w", err)
	}
	return reconstructAIAccessToken(row)
}

// Revoke は userID 所有のトークンを失効する。
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

// TouchLastUsed はトークンの最終利用日時を更新する。
func (r *AIAccessTokenRepository) TouchLastUsed(ctx context.Context, id entity.AIAccessTokenID, usedAt time.Time) error {
	if err := r.q.TouchAIAccessTokenLastUsed(ctx, sqlc.TouchAIAccessTokenLastUsedParams{
		ID:         uuid.UUID(id),
		LastUsedAt: pgtype.Timestamptz{Time: usedAt, Valid: true},
	}); err != nil {
		return fmt.Errorf("postgres: TouchAIAccessTokenLastUsed: %w", err)
	}
	return nil
}

func reconstructAIAccessToken(row sqlc.AiAccessToken) (*entity.AIAccessToken, error) {
	name, err := value.NewAIAccessTokenName(row.Name)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: ai_access_token name: %w", err)
	}
	hash, err := value.NewAIAccessTokenHash(row.TokenHash)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: ai_access_token hash: %w", err)
	}
	prefix, err := value.NewAIAccessTokenPrefix(row.TokenPrefix)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: ai_access_token prefix: %w", err)
	}
	return entity.ReconstructAIAccessToken(
		entity.AIAccessTokenID(row.ID),
		entity.UserID(row.UserID),
		name,
		hash,
		prefix,
		row.CreatedAt.Time,
		timePtrFromTimestamptz(row.LastUsedAt),
		timePtrFromTimestamptz(row.RevokedAt),
	), nil
}

func timestamptzFromPtr(t *time.Time) pgtype.Timestamptz {
	if t == nil {
		return pgtype.Timestamptz{}
	}
	return pgtype.Timestamptz{Time: *t, Valid: true}
}

func timePtrFromTimestamptz(t pgtype.Timestamptz) *time.Time {
	if !t.Valid {
		return nil
	}
	v := t.Time
	return &v
}
