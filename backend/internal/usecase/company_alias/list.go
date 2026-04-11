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

// List は企業に紐づく別名一覧を取得するUseCase。
type List struct {
	aliasRepo repository.CompanyAliasRepository
}

func NewList(aliasRepo repository.CompanyAliasRepository) *List {
	return &List{aliasRepo: aliasRepo}
}

// Execute はユーザー・企業に紐づく別名一覧を検索して返す。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	aliases, err := uc.aliasRepo.ListByCompanyID(ctx, input.UserID, input.CompanyID)
	if err != nil {
		return nil, err
	}

	return &ListOutput{CompanyAliases: aliases}, nil
}
