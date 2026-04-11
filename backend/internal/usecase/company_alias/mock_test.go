package companyalias

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type mockAliasRepo struct {
	createFn        func(ctx context.Context, alias *entity.CompanyAlias) error
	findByIDFn      func(ctx context.Context, userID entity.UserID, id entity.CompanyAliasID) (*entity.CompanyAlias, error)
	listByCompanyFn func(ctx context.Context, userID entity.UserID, companyID entity.CompanyID) ([]*entity.CompanyAlias, error)
	deleteFn        func(ctx context.Context, userID entity.UserID, id entity.CompanyAliasID) error
}

func (m *mockAliasRepo) Create(ctx context.Context, alias *entity.CompanyAlias) error {
	if m.createFn != nil {
		return m.createFn(ctx, alias)
	}
	return nil
}

func (m *mockAliasRepo) FindByID(ctx context.Context, userID entity.UserID, id entity.CompanyAliasID) (*entity.CompanyAlias, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, userID, id)
	}
	return nil, repository.ErrNotFound
}

func (m *mockAliasRepo) ListByCompanyID(ctx context.Context, userID entity.UserID, companyID entity.CompanyID) ([]*entity.CompanyAlias, error) {
	if m.listByCompanyFn != nil {
		return m.listByCompanyFn(ctx, userID, companyID)
	}
	return []*entity.CompanyAlias{}, nil
}

func (m *mockAliasRepo) Delete(ctx context.Context, userID entity.UserID, id entity.CompanyAliasID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, userID, id)
	}
	return nil
}

type mockCompanyRepo struct {
	findByIDFn func(ctx context.Context, userID entity.UserID, id entity.CompanyID) (*entity.Company, error)
}

func (m *mockCompanyRepo) Save(ctx context.Context, company *entity.Company) error {
	return nil
}

func (m *mockCompanyRepo) FindByID(ctx context.Context, userID entity.UserID, id entity.CompanyID) (*entity.Company, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, userID, id)
	}
	return nil, repository.ErrNotFound
}

func (m *mockCompanyRepo) ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.Company, error) {
	return []*entity.Company{}, nil
}

func (m *mockCompanyRepo) Delete(ctx context.Context, userID entity.UserID, id entity.CompanyID) error {
	return nil
}
