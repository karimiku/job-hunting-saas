// Package inboxclip は Chrome 拡張等で保存される求人ページクリップのユースケース群を提供する。
package inboxclip

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	domainservice "github.com/karimiku/job-hunting-saas/internal/domain/service"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// CreateInput は InboxClip 作成への入力。
type CreateInput struct {
	UserID entity.UserID
	URL    string
	Title  string
	Source string
	Guess  string
}

// CreateOutput は InboxClip 作成の出力。
type CreateOutput struct {
	Clip *entity.InboxClip
}

// Create は新しいクリップを保存するユースケース。
type Create struct {
	registration *domainservice.InboxClipRegistrationService
}

// NewCreate は Create ユースケースを生成する。
func NewCreate(repo repository.InboxClipRepository) *Create {
	return &Create{registration: domainservice.NewInboxClipRegistrationService(repo)}
}

// Execute は値オブジェクトを生成してから InboxClip を保存する。
func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	url, err := value.NewURL(input.URL)
	if err != nil {
		return nil, err
	}
	title, err := value.NewInboxClipTitle(input.Title)
	if err != nil {
		return nil, err
	}
	source, err := value.NewSource(input.Source)
	if err != nil {
		return nil, err
	}
	guess, err := value.NewInboxClipGuess(input.Guess)
	if err != nil {
		return nil, err
	}

	clip, err := uc.registration.Register(ctx, input.UserID, url, title, source, guess)
	if err != nil {
		return nil, err
	}
	return &CreateOutput{Clip: clip}, nil
}
