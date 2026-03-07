package companyalias

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type DeleteInput struct {
	UserID         entity.UserID
	CompanyAliasID entity.CompanyAliasID
}

type Delete struct {
	aliasRepo repository.CompanyAliasRepository
}

func NewDelete(aliasRepo repository.CompanyAliasRepository) *Delete {
	return &Delete{aliasRepo: aliasRepo}
}

func (uc *Delete) Execute(ctx context.Context, input DeleteInput) error {
	return uc.aliasRepo.Delete(ctx, input.UserID, input.CompanyAliasID)
}
