package company

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type mockCompanyRepo struct {
	saveFn       func(ctx context.Context, company *entity.Company) error
	findByIDFn   func(ctx context.Context, userID entity.UserID, id entity.CompanyID) (*entity.Company, error)
	listByUserFn func(ctx context.Context, userID entity.UserID) ([]*entity.Company, error)
	deleteFn     func(ctx context.Context, userID entity.UserID, id entity.CompanyID) error
}

func (m *mockCompanyRepo) Save(ctx context.Context, company *entity.Company) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, company)
	}
	return nil
}

func (m *mockCompanyRepo) FindByID(ctx context.Context, userID entity.UserID, id entity.CompanyID) (*entity.Company, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, userID, id)
	}
	return nil, repository.ErrNotFound
}

func (m *mockCompanyRepo) ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.Company, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID)
	}
	return []*entity.Company{}, nil
}

func (m *mockCompanyRepo) Delete(ctx context.Context, userID entity.UserID, id entity.CompanyID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, userID, id)
	}
	return nil
}
