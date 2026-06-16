// Package service は Entity / Value Object 単体では表現しづらいドメインルールを提供する。
package service

import (
	"context"
	"errors"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// InboxClipRegistrationService は InboxClip 登録時のドメインルールを扱う。
type InboxClipRegistrationService struct {
	repo repository.InboxClipRepository
}

// NewInboxClipRegistrationService は InboxClipRegistrationService を生成する。
func NewInboxClipRegistrationService(repo repository.InboxClipRepository) *InboxClipRegistrationService {
	return &InboxClipRegistrationService{repo: repo}
}

// Register は同一ユーザー・同一URLのクリップがあれば既存を返し、なければ新規作成する。
func (s *InboxClipRegistrationService) Register(
	ctx context.Context,
	userID entity.UserID,
	url value.URL,
	title value.InboxClipTitle,
	source value.Source,
	guess value.InboxClipGuess,
	contentText value.InboxClipContentText,
) (*entity.InboxClip, error) {
	existing, err := s.repo.FindByUserIDAndURL(ctx, userID, url)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}

	clip := entity.NewInboxClip(userID, url, title, source, guess, contentText)
	if err := s.repo.Create(ctx, clip); err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return s.repo.FindByUserIDAndURL(ctx, userID, url)
		}
		return nil, err
	}
	return clip, nil
}
