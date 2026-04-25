// Package companyalias は企業別名のユースケース群を提供する。
package companyalias

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// CreateInput は CompanyAliasCreate ユースケースへの入力。
type CreateInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
	Alias     string
}

// CreateOutput は CompanyAliasCreate ユースケースの出力。
type CreateOutput struct {
	CompanyAlias *entity.CompanyAlias
}

// Create は新しい企業別名を登録するUseCase。
type Create struct {
	aliasRepo   repository.CompanyAliasRepository
	companyRepo repository.CompanyRepository
}

// NewCreate は CompanyAliasCreate ユースケースを生成する。
func NewCreate(aliasRepo repository.CompanyAliasRepository, companyRepo repository.CompanyRepository) *Create {
	return &Create{aliasRepo: aliasRepo, companyRepo: companyRepo}
}

// Execute はCompanyIDの存在・所有を検証し、Aliasをバリデーションして新規CompanyAliasを生成・永続化する。
func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	if _, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID); err != nil {
		return nil, err
	}

	alias, err := value.NewAlias(input.Alias)
	if err != nil {
		return nil, err
	}

	a := entity.NewCompanyAlias(input.UserID, input.CompanyID, alias)

	if err := uc.aliasRepo.Create(ctx, a); err != nil {
		return nil, err
	}

	return &CreateOutput{CompanyAlias: a}, nil
}
