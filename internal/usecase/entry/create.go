package entry

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type CreateInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
	Route     string
	Source    string
	Memo      string
}

type CreateOutput struct {
	Entry *entity.Entry
}

type Create struct {
	entryRepo   repository.EntryRepository
	companyRepo repository.CompanyRepository
}

func NewCreate(entryRepo repository.EntryRepository, companyRepo repository.CompanyRepository) *Create {
	return &Create{entryRepo: entryRepo, companyRepo: companyRepo}
}

func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	// 指定されたCompanyが存在し、かつ操作ユーザーが所有していることを検証する
	if _, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID); err != nil {
		return nil, err
	}

	validatedRoute, err := value.NewRoute(input.Route)
	if err != nil {
		return nil, err
	}

	validatedSource, err := value.NewSource(input.Source)
	if err != nil {
		return nil, err
	}

	entry := entity.NewEntry(input.UserID, input.CompanyID, validatedRoute, validatedSource)

	if input.Memo != "" {
		entry.UpdateMemo(input.Memo)
	}

	if err := uc.entryRepo.Save(ctx, entry); err != nil {
		return nil, err
	}

	return &CreateOutput{Entry: entry}, nil
}
