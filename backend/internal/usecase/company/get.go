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

// Get は指定IDの企業を取得するUseCase。
type Get struct {
	companyRepo repository.CompanyRepository
}

func NewGet(companyRepo repository.CompanyRepository) *Get {
	return &Get{companyRepo: companyRepo}
}

// Execute はユーザーに紐づく企業をIDで検索して返す。
func (uc *Get) Execute(ctx context.Context, input GetInput) (*GetOutput, error) {
	c, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID)
	if err != nil {
		return nil, err
	}

	return &GetOutput{Company: c}, nil
}
