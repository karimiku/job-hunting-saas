// Package aiaccesstoken は AI / MCP 連携用アクセストークンのユースケース群を提供する。
package aiaccesstoken

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// ErrUserIDRequired は認証済み userID が入力されていないときに返される。
var ErrUserIDRequired = errors.New("userID is required")

const defaultAIIntegrationLabel = "AI連携トークン"

// CreateInput は AI 連携トークン作成への入力。
type CreateInput struct {
	UserID entity.UserID
	Name   string
}

// CreateOutput は AI 連携トークン作成の出力。
type CreateOutput struct {
	Token    *entity.AIAccessToken
	RawToken string
}

// Create は新しい AI 連携トークンを発行するユースケース。
type Create struct {
	repo repository.AIAccessTokenRepository
}

// NewCreate は Create ユースケースを生成する。
func NewCreate(repo repository.AIAccessTokenRepository) *Create {
	return &Create{repo: repo}
}

// Execute は平文トークンを一度だけ返し、DBには hash と prefix のみ保存する。
func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	if input.UserID.IsZero() {
		return nil, ErrUserIDRequired
	}
	tokenName := strings.TrimSpace(input.Name)
	if tokenName == "" {
		tokenName = defaultAIIntegrationLabel
	}
	name, err := value.NewAIAccessTokenName(tokenName)
	if err != nil {
		return nil, err
	}
	rawToken, err := value.GenerateAIAccessTokenRaw()
	if err != nil {
		return nil, fmt.Errorf("generate AI access token: %w", err)
	}
	hash, err := value.NewAIAccessTokenHashFromRaw(rawToken)
	if err != nil {
		return nil, err
	}
	prefix, err := value.NewAIAccessTokenPrefixFromRaw(rawToken)
	if err != nil {
		return nil, err
	}

	token := entity.NewAIAccessToken(input.UserID, name, hash, prefix)
	if err := uc.repo.Create(ctx, token); err != nil {
		return nil, err
	}
	return &CreateOutput{Token: token, RawToken: rawToken}, nil
}
