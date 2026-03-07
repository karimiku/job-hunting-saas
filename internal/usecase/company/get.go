package company

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type GetInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
}

type GetOutput struct {
	Company *entity.Company
}

type Get struct {
	companyRepo repository.CompanyRepository
}

func NewGet(companyRepo repository.CompanyRepository) *Get {
	return &Get{companyRepo: companyRepo}
}

func (uc *Get) Execute(ctx context.Context, input GetInput) (*GetOutput, error) {
	company, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID)
	if err != nil {
		return nil, err
	}

	return &GetOutput{Company: company}, nil
}
