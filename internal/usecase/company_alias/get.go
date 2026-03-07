package companyalias

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type GetInput struct {
	UserID         entity.UserID
	CompanyAliasID entity.CompanyAliasID
}

type GetOutput struct {
	CompanyAlias *entity.CompanyAlias
}

type Get struct {
	aliasRepo repository.CompanyAliasRepository
}

func NewGet(aliasRepo repository.CompanyAliasRepository) *Get {
	return &Get{aliasRepo: aliasRepo}
}

func (uc *Get) Execute(ctx context.Context, input GetInput) (*GetOutput, error) {
	companyAlias, err := uc.aliasRepo.FindByID(ctx, input.UserID, input.CompanyAliasID)
	if err != nil {
		return nil, err
	}

	return &GetOutput{CompanyAlias: companyAlias}, nil
}
