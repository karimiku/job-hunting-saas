// Package user はユーザの認証・登録ユースケースを提供する。
package user

import (
	"context"
	"errors"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type AuthenticateInput struct {
	Provider string // "google"
	Subject  string // Firebase UID
	Email    string
	Name     string
}

type AuthenticateOutput struct {
	User    *entity.User
	Created bool // 新規 User を作ったか
}

// Authenticate は Firebase 経由のログイン時に
// (provider, subject) → email の順で既存ユーザーを探し、
// なければ新規作成する UseCase。
type Authenticate struct {
	userRepo             repository.UserRepository
	externalIdentityRepo repository.ExternalIdentityRepository
}

func NewAuthenticate(userRepo repository.UserRepository, externalIdentityRepo repository.ExternalIdentityRepository) *Authenticate {
	return &Authenticate{
		userRepo:             userRepo,
		externalIdentityRepo: externalIdentityRepo,
	}
}

// Execute は Firebase 検証済み ID トークンから取れる情報で find-or-create する。
//
// フロー:
//  1. (provider, subject) で ExternalIdentity を検索 → 見つかれば対応 User を返す
//  2. email で User 検索 → 見つかれば ExternalIdentity を紐付け保存して返す
//  3. どちらもなければ User + ExternalIdentity を新規作成
func (uc *Authenticate) Execute(ctx context.Context, input AuthenticateInput) (*AuthenticateOutput, error) {
	provider, err := value.NewAuthProvider(input.Provider)
	if err != nil {
		return nil, err
	}
	if input.Subject == "" {
		return nil, errors.New("subject must not be empty")
	}
	email, err := value.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}
	name, err := value.NewUserName(input.Name)
	if err != nil {
		return nil, err
	}

	// 1. (provider, subject) で探す
	identity, err := uc.externalIdentityRepo.FindByProviderAndSubject(ctx, provider, input.Subject)
	if err == nil {
		user, err := uc.userRepo.FindByID(ctx, identity.UserID())
		if err != nil {
			return nil, err
		}
		return &AuthenticateOutput{User: user, Created: false}, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// 2. email で既存 User を探して紐付け
	existingUser, err := uc.userRepo.FindByEmail(ctx, email)
	if err == nil {
		newIdentity := entity.NewExternalIdentity(existingUser.ID(), provider, input.Subject)
		if err := uc.externalIdentityRepo.Save(ctx, newIdentity); err != nil {
			return nil, err
		}
		return &AuthenticateOutput{User: existingUser, Created: false}, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	// 3. 新規作成
	newUser := entity.NewUser(email, name)
	if err := uc.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}
	newIdentity := entity.NewExternalIdentity(newUser.ID(), provider, input.Subject)
	if err := uc.externalIdentityRepo.Save(ctx, newIdentity); err != nil {
		return nil, err
	}
	return &AuthenticateOutput{User: newUser, Created: true}, nil
}
