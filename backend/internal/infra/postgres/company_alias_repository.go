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

// CompanyAliasRepository は CompanyAliasRepository インターフェースの PostgreSQL 実装。
type CompanyAliasRepository struct {
	q *sqlc.Queries
}

// NewCompanyAliasRepository は CompanyAliasRepository を新規生成する。db には pgxpool.Pool もしくは tx を渡す。
func NewCompanyAliasRepository(db sqlc.DBTX) *CompanyAliasRepository {
	return &CompanyAliasRepository{q: sqlc.New(db)}
}

// Create は CompanyAlias を新規登録する。CompanyAlias は不変エンティティのため upsert はしない。
func (r *CompanyAliasRepository) Create(ctx context.Context, alias *entity.CompanyAlias) error {
	if err := r.q.CreateCompanyAlias(ctx, sqlc.CreateCompanyAliasParams{
		ID:        uuid.UUID(alias.ID()),
		UserID:    uuid.UUID(alias.UserID()),
		CompanyID: uuid.UUID(alias.CompanyID()),
		Alias:     alias.Alias().String(),
		CreatedAt: pgtype.Timestamptz{Time: alias.CreatedAt(), Valid: true},
	}); err != nil {
		if isUniqueViolation(err) {
			return repository.ErrAlreadyExists
		}
		return fmt.Errorf("postgres: CreateCompanyAlias: %w", err)
	}
	return nil
}

// FindByID は userID 所有の CompanyAlias を ID から取得する。存在しないか他ユーザー所有の場合は repository.ErrNotFound を返す。
func (r *CompanyAliasRepository) FindByID(ctx context.Context, userID entity.UserID, id entity.CompanyAliasID) (*entity.CompanyAlias, error) {
	row, err := r.q.FindCompanyAliasByID(ctx, sqlc.FindCompanyAliasByIDParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(id),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindCompanyAliasByID: %w", err)
	}

	return reconstructCompanyAlias(row)
}

// ListByCompanyID は userID 所有かつ companyID に紐づく CompanyAlias を全件返す。
func (r *CompanyAliasRepository) ListByCompanyID(ctx context.Context, userID entity.UserID, companyID entity.CompanyID) ([]*entity.CompanyAlias, error) {
	rows, err := r.q.ListCompanyAliasesByCompanyID(ctx, sqlc.ListCompanyAliasesByCompanyIDParams{
		UserID:    uuid.UUID(userID),
		CompanyID: uuid.UUID(companyID),
	})
	if err != nil {
		return nil, fmt.Errorf("postgres: ListCompanyAliasesByCompanyID: %w", err)
	}

	aliases := make([]*entity.CompanyAlias, 0, len(rows))
	for _, row := range rows {
		a, err := reconstructCompanyAlias(row)
		if err != nil {
			return nil, err
		}
		aliases = append(aliases, a)
	}

	return aliases, nil
}

// Delete は userID 所有の CompanyAlias を ID から削除する。存在しないか他ユーザー所有の場合は repository.ErrNotFound を返す。
func (r *CompanyAliasRepository) Delete(ctx context.Context, userID entity.UserID, id entity.CompanyAliasID) error {
	n, err := r.q.DeleteCompanyAlias(ctx, sqlc.DeleteCompanyAliasParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(id),
	})
	if err != nil {
		return fmt.Errorf("postgres: DeleteCompanyAlias: %w", err)
	}
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func reconstructCompanyAlias(row sqlc.CompanyAlias) (*entity.CompanyAlias, error) {
	alias, err := value.NewAlias(row.Alias)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: company alias: %w", err)
	}

	return entity.ReconstructCompanyAlias(
		entity.CompanyAliasID(row.ID),
		entity.UserID(row.UserID),
		entity.CompanyID(row.CompanyID),
		alias,
		row.CreatedAt.Time,
	), nil
}
