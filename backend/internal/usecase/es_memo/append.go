// Package es_memo は ES / 自己PR / 面接ネタ用メモのユースケースを提供する。
package es_memo

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type AppendInput struct {
	UserID   entity.UserID
	EntryID  *entity.EntryID
	Category string
	Title    string
	Content  string
	Source   string
}

type AppendOutput struct {
	Memo *entity.ESMemo
}

type Append struct {
	memoRepo  repository.ESMemoRepository
	entryRepo repository.EntryRepository
}

func NewAppend(memoRepo repository.ESMemoRepository, entryRepo repository.EntryRepository) *Append {
	return &Append{memoRepo: memoRepo, entryRepo: entryRepo}
}

func (uc *Append) Execute(ctx context.Context, input AppendInput) (*AppendOutput, error) {
	if input.EntryID != nil {
		if _, err := uc.entryRepo.FindByID(ctx, input.UserID, *input.EntryID); err != nil {
			return nil, err
		}
	}
	category, err := value.NewESMemoCategory(input.Category)
	if err != nil {
		return nil, err
	}
	title, err := value.NewESMemoTitle(input.Title)
	if err != nil {
		return nil, err
	}
	content, err := value.NewESMemoContent(input.Content)
	if err != nil {
		return nil, err
	}
	source, err := value.NewESMemoSource(input.Source)
	if err != nil {
		return nil, err
	}

	memo := entity.NewESMemo(input.UserID, input.EntryID, category, title, content, source)
	if err := uc.memoRepo.Save(ctx, memo); err != nil {
		return nil, err
	}
	return &AppendOutput{Memo: memo}, nil
}
