package entry

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type ListInput struct {
	UserID    entity.UserID
	Status    *string
	StageKind *string
	Source    *string
}

type ListOutput struct {
	Entries []*entity.Entry
}

// List はユーザーに紐づくエントリー一覧を取得するUseCase。
type List struct {
	entryRepo repository.EntryRepository
}

func NewList(entryRepo repository.EntryRepository) *List {
	return &List{entryRepo: entryRepo}
}

// Execute はフィルタ条件をバリデーションし、エントリー一覧を検索して返す。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	filter := repository.EntryFilter{}

	if input.Status != nil {
		s, err := value.NewEntryStatus(*input.Status)
		if err != nil {
			return nil, err
		}
		filter.Status = &s
	}

	if input.StageKind != nil {
		k, err := value.NewStageKind(*input.StageKind)
		if err != nil {
			return nil, err
		}
		filter.StageKind = &k
	}

	if input.Source != nil {
		src, err := value.NewSource(*input.Source)
		if err != nil {
			return nil, err
		}
		filter.Source = &src
	}

	entries, err := uc.entryRepo.ListByUserID(ctx, input.UserID, filter)
	if err != nil {
		return nil, err
	}

	return &ListOutput{Entries: entries}, nil
}
