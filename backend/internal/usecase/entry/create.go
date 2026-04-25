// Package entry は応募エントリに対するユースケース群を提供する。
package entry

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// CreateInput は EntryCreate ユースケースへの入力。
type CreateInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
	Route     string
	Source    string
	Memo      string
}

// CreateOutput は EntryCreate ユースケースの出力。
type CreateOutput struct {
	Entry *entity.Entry
}

// Create は新しいエントリーを登録するUseCase。
type Create struct {
	entryRepo   repository.EntryRepository
	companyRepo repository.CompanyRepository
}

// NewCreate は EntryCreate ユースケースを生成する。
func NewCreate(entryRepo repository.EntryRepository, companyRepo repository.CompanyRepository) *Create {
	return &Create{entryRepo: entryRepo, companyRepo: companyRepo}
}

// Execute はCompanyIDの存在・所有を検証し、Route/Sourceをバリデーションして新規Entryを生成・永続化する。
func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	if _, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID); err != nil {
		return nil, err
	}

	route, err := value.NewRoute(input.Route)
	if err != nil {
		return nil, err
	}

	source, err := value.NewSource(input.Source)
	if err != nil {
		return nil, err
	}

	e := entity.NewEntry(input.UserID, input.CompanyID, route, source)

	if input.Memo != "" {
		e.UpdateMemo(input.Memo)
	}

	if err := uc.entryRepo.Save(ctx, e); err != nil {
		return nil, err
	}

	return &CreateOutput{Entry: e}, nil
}
