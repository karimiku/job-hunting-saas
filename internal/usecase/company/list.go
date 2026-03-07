package company

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type ListInput struct {
	UserID entity.UserID
}

type ListOutput struct {
	Companies []*entity.Company
}

type List struct {
	companyRepo repository.CompanyRepository
}

func NewList(companyRepo repository.CompanyRepository) *List {
	return &List{companyRepo: companyRepo}
}

func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	companies, err := uc.companyRepo.ListByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &ListOutput{Companies: companies}, nil
}
