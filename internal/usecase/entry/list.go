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

type List struct {
	entryRepo repository.EntryRepository
}

func NewList(entryRepo repository.EntryRepository) *List {
	return &List{entryRepo: entryRepo}
}

// Execute はフィルタ条件の各値をバリデーションしてから検索する。
// nilのフィールドはフィルタを適用しない。
func (uc *List) Execute(ctx context.Context, input ListInput) (*ListOutput, error) {
	entryFilter := repository.EntryFilter{}

	if input.Status != nil {
		validatedStatus, err := value.NewEntryStatus(*input.Status)
		if err != nil {
			return nil, err
		}
		entryFilter.Status = &validatedStatus
	}

	if input.StageKind != nil {
		validatedStageKind, err := value.NewStageKind(*input.StageKind)
		if err != nil {
			return nil, err
		}
		entryFilter.StageKind = &validatedStageKind
	}

	if input.Source != nil {
		validatedSource, err := value.NewSource(*input.Source)
		if err != nil {
			return nil, err
		}
		entryFilter.Source = &validatedSource
	}

	entries, err := uc.entryRepo.ListByUserID(ctx, input.UserID, entryFilter)
	if err != nil {
		return nil, err
	}

	return &ListOutput{Entries: entries}, nil
}
