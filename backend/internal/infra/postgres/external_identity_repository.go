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

// ExternalIdentityRepository は ExternalIdentityRepository インターフェースの PostgreSQL 実装。
type ExternalIdentityRepository struct {
	q *sqlc.Queries
}

func NewExternalIdentityRepository(db sqlc.DBTX) *ExternalIdentityRepository {
	return &ExternalIdentityRepository{q: sqlc.New(db)}
}

func (r *ExternalIdentityRepository) Save(ctx context.Context, identity *entity.ExternalIdentity) error {
	if err := r.q.InsertExternalIdentity(ctx, sqlc.InsertExternalIdentityParams{
		ID:        uuid.UUID(identity.ID()),
		UserID:    uuid.UUID(identity.UserID()),
		Provider:  sqlc.AuthProvider(identity.Provider().String()),
		Subject:   identity.Subject(),
		CreatedAt: pgtype.Timestamptz{Time: identity.CreatedAt(), Valid: true},
	}); err != nil {
		if isUniqueViolation(err) {
			return repository.ErrAlreadyExists
		}
		return fmt.Errorf("postgres: InsertExternalIdentity: %w", err)
	}
	return nil
}

func (r *ExternalIdentityRepository) FindByProviderAndSubject(ctx context.Context, provider value.AuthProvider, subject string) (*entity.ExternalIdentity, error) {
	row, err := r.q.FindExternalIdentityByProviderAndSubject(ctx, sqlc.FindExternalIdentityByProviderAndSubjectParams{
		Provider: sqlc.AuthProvider(provider.String()),
		Subject:  subject,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindExternalIdentityByProviderAndSubject: %w", err)
	}
	return reconstructExternalIdentity(row)
}

func reconstructExternalIdentity(row sqlc.ExternalIdentity) (*entity.ExternalIdentity, error) {
	provider, err := value.NewAuthProvider(string(row.Provider))
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: external identity provider: %w", err)
	}
	return entity.ReconstructExternalIdentity(
		entity.ExternalIdentityID(row.ID),
		entity.UserID(row.UserID),
		provider,
		row.Subject,
		row.CreatedAt.Time,
	), nil
}
