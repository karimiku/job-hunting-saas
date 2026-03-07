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

type Delete struct {
	companyRepo repository.CompanyRepository
}

func NewDelete(companyRepo repository.CompanyRepository) *Delete {
	return &Delete{companyRepo: companyRepo}
}

func (uc *Delete) Execute(ctx context.Context, input DeleteInput) error {
	return uc.companyRepo.Delete(ctx, input.UserID, input.CompanyID)
}
