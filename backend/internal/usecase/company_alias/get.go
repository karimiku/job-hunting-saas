package companyalias

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// GetInput は CompanyAliasGet ユースケースへの入力。
type GetInput struct {
	UserID         entity.UserID
	CompanyAliasID entity.CompanyAliasID
}

// GetOutput は CompanyAliasGet ユースケースの出力。
type GetOutput struct {
	CompanyAlias *entity.CompanyAlias
}

// Get は指定IDの企業別名を取得するUseCase。
type Get struct {
	aliasRepo repository.CompanyAliasRepository
}

// NewGet は CompanyAliasGet ユースケースを生成する。
func NewGet(aliasRepo repository.CompanyAliasRepository) *Get {
	return &Get{aliasRepo: aliasRepo}
}

// Execute はユーザーに紐づく企業別名をIDで検索して返す。
func (uc *Get) Execute(ctx context.Context, input GetInput) (*GetOutput, error) {
	a, err := uc.aliasRepo.FindByID(ctx, input.UserID, input.CompanyAliasID)
	if err != nil {
		return nil, err
	}

	return &GetOutput{CompanyAlias: a}, nil
}
