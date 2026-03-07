package companyalias

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type ListInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
}

type ListOutput struct {
	CompanyAliases []*entity.CompanyAlias
}

type List struct {
	aliasRepo repository.CompanyAliasRepository
}

func NewList(aliasRepo repository.CompanyAliasRepository) *List {
	return &List{aliasRepo: aliasRepo}
}

func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	companyAliases, err := uc.aliasRepo.ListByCompanyID(ctx, input.UserID, input.CompanyID)
	if err != nil {
		return nil, err
	}

	return &ListOutput{CompanyAliases: companyAliases}, nil
}
