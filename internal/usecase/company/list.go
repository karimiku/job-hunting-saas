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

// List はユーザーに紐づく企業一覧を取得するUseCase。
type List struct {
	companyRepo repository.CompanyRepository
}

func NewList(companyRepo repository.CompanyRepository) *List {
	return &List{companyRepo: companyRepo}
}

// Execute はユーザーIDで企業一覧を検索して返す。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	companies, err := uc.companyRepo.ListByUserID(ctx, input.UserID)
	if err != nil {
		return nil, err
	}

	return &ListOutput{Companies: companies}, nil
}
