package entry

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

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

type mockEntryRepo struct {
	saveFn       func(ctx context.Context, entry *entity.Entry) error
	findByIDFn   func(ctx context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error)
	listByUserFn func(ctx context.Context, userID entity.UserID, filter repository.EntryFilter) ([]*entity.Entry, error)
	deleteFn     func(ctx context.Context, userID entity.UserID, id entity.EntryID) error
}

func (m *mockEntryRepo) Save(ctx context.Context, entry *entity.Entry) error {
	if m.saveFn != nil {
		return m.saveFn(ctx, entry)
	}
	return nil
}

func (m *mockEntryRepo) FindByID(ctx context.Context, userID entity.UserID, id entity.EntryID) (*entity.Entry, error) {
	if m.findByIDFn != nil {
		return m.findByIDFn(ctx, userID, id)
	}
	return nil, repository.ErrNotFound
}

func (m *mockEntryRepo) ListByUserID(ctx context.Context, userID entity.UserID, filter repository.EntryFilter) ([]*entity.Entry, error) {
	if m.listByUserFn != nil {
		return m.listByUserFn(ctx, userID, filter)
	}
	return []*entity.Entry{}, nil
}

func (m *mockEntryRepo) Delete(ctx context.Context, userID entity.UserID, id entity.EntryID) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, userID, id)
	}
	return nil
}
