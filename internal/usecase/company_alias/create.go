package companyalias

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type CreateInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
	Alias     string
}

type CreateOutput struct {
	CompanyAlias *entity.CompanyAlias
}

type Create struct {
	aliasRepo   repository.CompanyAliasRepository
	companyRepo repository.CompanyRepository
}

func NewCreate(aliasRepo repository.CompanyAliasRepository, companyRepo repository.CompanyRepository) *Create {
	return &Create{aliasRepo: aliasRepo, companyRepo: companyRepo}
}

func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	// 指定されたCompanyが存在し、かつ操作ユーザーが所有していることを検証する
	if _, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID); err != nil {
		return nil, err
	}

	validatedAlias, err := value.NewAlias(input.Alias)
	if err != nil {
		return nil, err
	}

	companyAlias := entity.NewCompanyAlias(input.UserID, input.CompanyID, validatedAlias)

	if err := uc.aliasRepo.Create(ctx, companyAlias); err != nil {
		return nil, err
	}

	return &CreateOutput{CompanyAlias: companyAlias}, nil
}
