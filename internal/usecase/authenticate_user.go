package usecase

import (
	"context"
	"errors"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type AuthenticateUserInput struct {
	Email string
	Name  string
}

type AuthenticateUserOutput struct {
	User    *entity.User
	Created bool // 新規作成されたか
}

// AuthenticateUser はGoogle OAuthログイン時のFind or Createを行うUseCase。
type AuthenticateUser struct {
	userRepo repository.UserRepository
}

func NewAuthenticateUser(userRepo repository.UserRepository) *AuthenticateUser {
	return &AuthenticateUser{userRepo: userRepo}
}

// Execute はメールで既存ユーザーを検索し、いなければ新規作成して返す。
func (uc *AuthenticateUser) Execute(ctx context.Context, input AuthenticateUserInput) (*AuthenticateUserOutput, error) {
	email, err := value.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	name, err := value.NewUserName(input.Name)
	if err != nil {
		return nil, err
	}

	existingUser, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil {
		return &AuthenticateUserOutput{User: existingUser, Created: false}, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	newUser := entity.NewUser(email, name)
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	return &AuthenticateUserOutput{User: newUser, Created: true}, nil
}
