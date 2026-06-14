package inmemory

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// AccountRepository はメモリ上でユーザーアカウント削除を提供する。
type AccountRepository struct {
	userRepo *UserRepository
}

// NewAccountRepository は AccountRepository を新規生成する。
func NewAccountRepository(userRepo *UserRepository) *AccountRepository {
	return &AccountRepository{userRepo: userRepo}
}

// DeleteAccount は User を削除する。
func (r *AccountRepository) DeleteAccount(ctx context.Context, userID entity.UserID) error {
	return r.userRepo.Delete(ctx, userID)
}
