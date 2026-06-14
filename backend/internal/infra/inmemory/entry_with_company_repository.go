package inmemory

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// EntryWithCompanyRepository はメモリ上で Company と Entry を同時保存する。
type EntryWithCompanyRepository struct {
	companyRepo *CompanyRepository
	entryRepo   *EntryRepository
}

// NewEntryWithCompanyRepository は EntryWithCompanyRepository を新規生成する。
func NewEntryWithCompanyRepository(companyRepo *CompanyRepository, entryRepo *EntryRepository) *EntryWithCompanyRepository {
	return &EntryWithCompanyRepository{
		companyRepo: companyRepo,
		entryRepo:   entryRepo,
	}
}

// SaveEntryWithCompany は Company と Entry を保存する。
func (r *EntryWithCompanyRepository) SaveEntryWithCompany(ctx context.Context, company *entity.Company, entry *entity.Entry) error {
	if err := r.companyRepo.Save(ctx, company); err != nil {
		return err
	}
	if err := r.entryRepo.Save(ctx, entry); err != nil {
		return err
	}
	return nil
}
