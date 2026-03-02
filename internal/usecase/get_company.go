package usecase

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type GetCompanyInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
}

type GetCompanyOutput struct {
	Company *entity.Company
}

// GetCompany は指定IDの企業を取得するUseCase。
type GetCompany struct {
	companyRepo repository.CompanyRepository
}

func NewGetCompany(companyRepo repository.CompanyRepository) *GetCompany {
	return &GetCompany{companyRepo: companyRepo}
}

// Execute はユーザーに紐づく企業をIDで検索して返す。
func (uc *GetCompany) Execute(ctx context.Context, input GetCompanyInput) (*GetCompanyOutput, error) {
	company, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID)
	if err != nil {
		return nil, err
	}

	return &GetCompanyOutput{Company: company}, nil
}
