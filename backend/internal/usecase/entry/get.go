package entry

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
)

type GetInput struct {
	UserID  entity.UserID
	EntryID entity.EntryID
}

type GetOutput struct {
	Entry *entity.Entry
}

// Get は指定IDのエントリーを取得するUseCase。
type Get struct {
	entryRepo repository.EntryRepository
}

func NewGet(entryRepo repository.EntryRepository) *Get {
	return &Get{entryRepo: entryRepo}
}

// Execute はユーザーに紐づくエントリーをIDで検索して返す。
func (uc *Get) Execute(ctx context.Context, input GetInput) (*GetOutput, error) {
	e, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID)
	if err != nil {
		return nil, err
	}

	return &GetOutput{Entry: e}, nil
}
