package usecase

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type DeleteCompanyInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
}

type DeleteCompany struct {
	companyRepo repository.CompanyRepository
}

func NewDeleteCompany(companyRepo repository.CompanyRepository) *DeleteCompany {
	return &DeleteCompany{companyRepo: companyRepo}
}

func (uc *DeleteCompany) Execute(ctx context.Context, input DeleteCompanyInput) error {
	return uc.companyRepo.Delete(ctx, input.UserID, input.CompanyID)
}
