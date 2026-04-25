package companyalias

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// DeleteInput は CompanyAliasDelete ユースケースへの入力。
type DeleteInput struct {
	UserID         entity.UserID
	CompanyAliasID entity.CompanyAliasID
}

// Delete は指定IDの企業別名を削除するUseCase。
type Delete struct {
	aliasRepo repository.CompanyAliasRepository
}

// NewDelete は CompanyAliasDelete ユースケースを生成する。
func NewDelete(aliasRepo repository.CompanyAliasRepository) *Delete {
	return &Delete{aliasRepo: aliasRepo}
}

// Execute はユーザーに紐づく企業別名をIDで削除する。
func (uc *Delete) Execute(ctx context.Context, input DeleteInput) error {
	return uc.aliasRepo.Delete(ctx, input.UserID, input.CompanyAliasID)
}
