package aiaccesstoken

import (
	"context"
	"errors"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// VerifyInput は AI 連携トークン検証への入力。
type VerifyInput struct {
	RawToken string
}

// VerifyOutput は AI 連携トークン検証の出力。
type VerifyOutput struct {
	UserID  entity.UserID
	TokenID entity.AIAccessTokenID
}

// Verify は Bearer token からユーザーを解決するユースケース。
type Verify struct {
	repo repository.AIAccessTokenRepository
}

// NewVerify は Verify ユースケースを生成する。
func NewVerify(repo repository.AIAccessTokenRepository) *Verify {
	return &Verify{repo: repo}
}

// Execute は平文トークンを hash 化し、有効なトークンなら所有ユーザーを返す。
func (uc *Verify) Execute(ctx context.Context, input VerifyInput) (*VerifyOutput, error) {
	hash, err := value.NewAIAccessTokenHashFromRaw(input.RawToken)
	if err != nil {
		return nil, err
	}
	token, err := uc.repo.FindActiveByHash(ctx, hash)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, value.ErrAIAccessTokenInvalid
		}
		return nil, err
	}
	now := time.Now()
	if err := uc.repo.TouchLastUsed(ctx, token.ID(), now); err != nil {
		return nil, err
	}
	return &VerifyOutput{UserID: token.UserID(), TokenID: token.ID()}, nil
}

// VerifyBearerToken は middleware.BearerTokenVerifier を満たす adapter メソッド。
func (uc *Verify) VerifyBearerToken(ctx context.Context, rawToken string) (entity.UserID, error) {
	out, err := uc.Execute(ctx, VerifyInput{RawToken: rawToken})
	if err != nil {
		return entity.UserID{}, err
	}
	return out.UserID, nil
}
