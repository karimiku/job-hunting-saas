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

type CompanyRepository struct {
	q *sqlc.Queries
}

func NewCompanyRepository(db sqlc.DBTX) *CompanyRepository {
	return &CompanyRepository{q: sqlc.New(db)}
}

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
