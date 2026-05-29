package companyalias

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// ListInput は CompanyAliasList ユースケースへの入力。
type ListInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
}

// ListOutput は CompanyAliasList ユースケースの出力。
type ListOutput struct {
	CompanyAliases []*entity.CompanyAlias
}

// List は企業に紐づく別名一覧を取得するUseCase。
type List struct {
	aliasRepo   repository.CompanyAliasRepository
	companyRepo repository.CompanyRepository
}

// NewList は CompanyAliasList ユースケースを生成する。
func NewList(aliasRepo repository.CompanyAliasRepository, companyRepo repository.CompanyRepository) *List {
	return &List{aliasRepo: aliasRepo, companyRepo: companyRepo}
}

// Execute はCompanyIDの存在・所有を検証し、紐づく別名一覧を検索して返す。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	if _, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID); err != nil {
		return nil, err
	}

	aliases, err := uc.aliasRepo.ListByCompanyID(ctx, input.UserID, input.CompanyID)
	if err != nil {
		return nil, err
	}

	return &ListOutput{CompanyAliases: aliases}, nil
}
