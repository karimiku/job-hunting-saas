package user

import (
	"context"
	"errors"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type AuthenticateInput struct {
	Email string
	Name  string
}

type AuthenticateOutput struct {
	User    *entity.User
	Created bool // 新規作成されたかどうか（初回ログイン判定に使用）
}

// Authenticate はGoogle OAuthログイン時のFind or Createを行う。
// メールアドレスで既存ユーザーを検索し、いなければ新規作成する。
type Authenticate struct {
	userRepo repository.UserRepository
}

func NewAuthenticate(userRepo repository.UserRepository) *Authenticate {
	return &Authenticate{userRepo: userRepo}
}

func (uc *Authenticate) Execute(ctx context.Context, input AuthenticateInput) (*AuthenticateOutput, error) {
	validatedEmail, err := value.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	validatedName, err := value.NewUserName(input.Name)
	if err != nil {
		return nil, err
	}

	// 既存ユーザーが見つかればそのまま返す（ログイン）
	existingUser, err := uc.userRepo.FindByEmail(ctx, validatedEmail)
	if err == nil {
		return &AuthenticateOutput{User: existingUser, Created: false}, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// ErrNotFound の場合は新規ユーザーを作成する（サインアップ）
	newUser := entity.NewUser(validatedEmail, validatedName)
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	return &AuthenticateOutput{User: newUser, Created: true}, nil
}
