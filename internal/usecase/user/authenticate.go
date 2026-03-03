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
	Created bool // 新規作成されたか
}

// Authenticate はGoogle OAuthログイン時のFind or Createを行うUseCase。
type Authenticate struct {
	userRepo repository.UserRepository
}

func NewAuthenticate(userRepo repository.UserRepository) *Authenticate {
	return &Authenticate{userRepo: userRepo}
}

// Execute はメールで既存ユーザーを検索し、いなければ新規作成して返す。
func (uc *Authenticate) Execute(ctx context.Context, input AuthenticateInput) (*AuthenticateOutput, error) {
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
		return &AuthenticateOutput{User: existingUser, Created: false}, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	newUser := entity.NewUser(email, name)
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	return &AuthenticateOutput{User: newUser, Created: true}, nil
}
