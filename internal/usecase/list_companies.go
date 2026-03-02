package usecase

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type ListCompaniesInput struct {
	UserID entity.UserID
}

type ListCompaniesOutput struct {
	Companies []*entity.Company
}

// ListCompanies はユーザーに紐づく企業一覧を取得するUseCase。
type ListCompanies struct {
	companyRepo repository.CompanyRepository
}

func NewListCompanies(companyRepo repository.CompanyRepository) *ListCompanies {
	return &ListCompanies{companyRepo: companyRepo}
}

// Execute はユーザーIDで企業一覧を検索して返す。
func (uc *ListCompanies) Execute(ctx context.Context, input ListCompaniesInput) (*ListCompaniesOutput, error) {
	companies, err := uc.companyRepo.ListByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &ListCompaniesOutput{Companies: companies}, nil
}
