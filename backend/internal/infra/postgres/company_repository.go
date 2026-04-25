// Package postgres は Repository インターフェースの PostgreSQL 実装を提供する。
// sqlc で生成された型安全クエリと pgx/v5 を使う。
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

// CompanyRepository は CompanyRepository インターフェースの PostgreSQL 実装。
type CompanyRepository struct {
	q *sqlc.Queries
}

// NewCompanyRepository は CompanyRepository を新規生成する。db には pgxpool.Pool もしくは tx を渡す。
func NewCompanyRepository(db sqlc.DBTX) *CompanyRepository {
	return &CompanyRepository{q: sqlc.New(db)}
}

// Save は Company を upsert する。同じ ID があれば更新、なければ作成。
func (r *CompanyRepository) Save(ctx context.Context, company *entity.Company) error {
	if err := r.q.UpsertCompany(ctx, sqlc.UpsertCompanyParams{
		ID:        uuid.UUID(company.ID()),
		UserID:    uuid.UUID(company.UserID()),
		Name:      company.Name().String(),
		Memo:      company.Memo(),
		CreatedAt: pgtype.Timestamptz{Time: company.CreatedAt(), Valid: true},
		UpdatedAt: pgtype.Timestamptz{Time: company.UpdatedAt(), Valid: true},
	}); err != nil {
		return fmt.Errorf("postgres: UpsertCompany: %w", err)
	}
	return nil
}

// FindByID は userID 所有の Company を ID から取得する。存在しないか他ユーザー所有の場合は repository.ErrNotFound を返す。
func (r *CompanyRepository) FindByID(ctx context.Context, userID entity.UserID, id entity.CompanyID) (*entity.Company, error) {
	row, err := r.q.FindCompanyByID(ctx, sqlc.FindCompanyByIDParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(id),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrNotFound
		}
		return nil, fmt.Errorf("postgres: FindCompanyByID: %w", err)
	}

	return reconstructCompany(row)
}

// ListByUserID は userID に紐づく Company を全件返す。
func (r *CompanyRepository) ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.Company, error) {
	rows, err := r.q.ListCompaniesByUserID(ctx, uuid.UUID(userID))
	if err != nil {
		return nil, fmt.Errorf("postgres: ListCompaniesByUserID: %w", err)
	}

	companies := make([]*entity.Company, 0, len(rows))
	for _, row := range rows {
		c, err := reconstructCompany(row)
		if err != nil {
			return nil, err
		}
		companies = append(companies, c)
	}

	return companies, nil
}

// Delete は userID 所有の Company を ID から削除する。存在しないか他ユーザー所有の場合は repository.ErrNotFound を返す。
func (r *CompanyRepository) Delete(ctx context.Context, userID entity.UserID, id entity.CompanyID) error {
	n, err := r.q.DeleteCompany(ctx, sqlc.DeleteCompanyParams{
		UserID: uuid.UUID(userID),
		ID:     uuid.UUID(id),
	})
	if err != nil {
		return fmt.Errorf("postgres: DeleteCompany: %w", err)
	}
	if n == 0 {
		return repository.ErrNotFound
	}
	return nil
}

func reconstructCompany(row sqlc.Company) (*entity.Company, error) {
	name, err := value.NewCompanyName(row.Name)
	if err != nil {
		return nil, fmt.Errorf("BUG: invalid data in DB: company name: %w", err)
	}

	return entity.ReconstructCompany(
		entity.CompanyID(row.ID),
		entity.UserID(row.UserID),
		name,
		row.Memo,
		row.CreatedAt.Time,
		row.UpdatedAt.Time,
	), nil
}
