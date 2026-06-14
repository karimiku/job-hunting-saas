package user

import (
	"context"
	"errors"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

// DeleteAccountInput は退会ユースケースへの入力。
type DeleteAccountInput struct {
	UserID entity.UserID
}

// DeleteAccount はユーザー本人と関連データを削除する UseCase。
type DeleteAccount struct {
	accountRepo repository.AccountRepository
}

// NewDeleteAccount は DeleteAccount ユースケースを生成する。
func NewDeleteAccount(accountRepo repository.AccountRepository) *DeleteAccount {
	return &DeleteAccount{accountRepo: accountRepo}
}

// Execute はユーザー本人のアカウントを削除する。
func (uc *DeleteAccount) Execute(ctx context.Context, input DeleteAccountInput) error {
	if input.UserID.IsZero() {
		return errors.New("userID must not be empty")
	}
	return uc.accountRepo.DeleteAccount(ctx, input.UserID)
}
