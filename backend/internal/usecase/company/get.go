package company

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// GetInput は CompanyGet ユースケースへの入力。
type GetInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
}

// GetOutput は CompanyGet ユースケースの出力。
type GetOutput struct {
	Company *entity.Company
}

// Get は指定IDの企業を取得するUseCase。
type Get struct {
	companyRepo repository.CompanyRepository
}

// NewGet は CompanyGet ユースケースを生成する。
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
