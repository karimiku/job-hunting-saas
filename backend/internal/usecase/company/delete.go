package company

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type DeleteInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
}

// Delete は指定IDの企業を削除するUseCase。
type Delete struct {
	companyRepo repository.CompanyRepository
}

func NewDelete(companyRepo repository.CompanyRepository) *Delete {
	return &Delete{companyRepo: companyRepo}
}

// Execute はユーザーに紐づく企業をIDで削除する。
func (uc *Delete) Execute(ctx context.Context, input DeleteInput) error {
	return uc.companyRepo.Delete(ctx, input.UserID, input.CompanyID)
}
