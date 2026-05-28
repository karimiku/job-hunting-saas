// Package inboxclip は Chrome 拡張等で保存される求人ページクリップのユースケース群を提供する。
package inboxclip

import (
	"context"
	"errors"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
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
	repo repository.InboxClipRepository
}

// NewCreate は Create ユースケースを生成する。
func NewCreate(repo repository.InboxClipRepository) *Create {
	return &Create{repo: repo}
}

// Execute は値オブジェクトを生成してから InboxClip を保存する。
//
// 重複URLの最低限対応 (#98): 同一ユーザーが同じ URL のクリップを既に保存している場合は、
// 新規作成せず既存クリップをそのまま返す（冪等）。これにより拡張からの多重保存で
// Inbox がノイズで溢れるのを防ぐ。完全な重複統合は P1 のスコープ外。
func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	url, err := value.NewURL(input.URL)
	if err != nil {
		return nil, err
	}
	source, err := value.NewSource(input.Source)
	if err != nil {
		return nil, err
	}

	// 既存の同一 URL クリップがあれば、それを返して新規作成をスキップする。
	if existing, err := uc.repo.FindByUserIDAndURL(ctx, input.UserID, url); err == nil {
		return &CreateOutput{Clip: existing}, nil
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	clip := entity.NewInboxClip(input.UserID, url, input.Title, source, input.Guess)
	if err := uc.repo.Create(ctx, clip); err != nil {
		return nil, err
	}
	return &CreateOutput{Clip: clip}, nil
}
