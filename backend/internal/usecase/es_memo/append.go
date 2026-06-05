// Package esmemo は ES / 自己PR / 面接ネタ用メモのユースケースを提供する。
package esmemo

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// AppendInput はESメモ追記ユースケースへの入力。
type AppendInput struct {
	UserID   entity.UserID
	EntryID  *entity.EntryID
	Category string
	Title    string
	Content  string
	Source   string
}

// AppendOutput はESメモ追記ユースケースの出力。
type AppendOutput struct {
	Memo *entity.ESMemo
}

// Append はESメモを追記するUseCase。
type Append struct {
	memoRepo  repository.ESMemoRepository
	entryRepo repository.EntryRepository
}

// NewAppend はESメモ追記ユースケースを生成する。
func NewAppend(memoRepo repository.ESMemoRepository, entryRepo repository.EntryRepository) *Append {
	return &Append{memoRepo: memoRepo, entryRepo: entryRepo}
}

// Execute はEntryIDの存在・所有を検証し、ESメモを生成・永続化する。
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
