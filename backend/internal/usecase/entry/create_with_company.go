package entry

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// CreateWithCompanyInput は Company と Entry の同時作成ユースケースへの入力。
type CreateWithCompanyInput struct {
	UserID      entity.UserID
	CompanyName string
	Route       string
	Source      string
	SourceURL   string
	Memo        string
}

// CreateWithCompanyOutput は Company と Entry の同時作成ユースケースの出力。
type CreateWithCompanyOutput struct {
	Company *entity.Company
	Entry   *entity.Entry
}

// CreateWithCompany は新しい Company と Entry を部分作成なしで登録する UseCase。
type CreateWithCompany struct {
	repo repository.EntryWithCompanyRepository
}

// NewCreateWithCompany は CreateWithCompany ユースケースを生成する。
func NewCreateWithCompany(repo repository.EntryWithCompanyRepository) *CreateWithCompany {
	return &CreateWithCompany{repo: repo}
}

// Execute は入力を検証し、Company と Entry を同一永続化単位で作成する。
func (uc *CreateWithCompany) Execute(ctx context.Context, input CreateWithCompanyInput) (*CreateWithCompanyOutput, error) {
	companyName, err := value.NewCompanyName(input.CompanyName)
	if err != nil {
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

	company := entity.NewCompany(input.UserID, companyName)
	entry := entity.NewEntry(input.UserID, company.ID(), route, source)

	if input.SourceURL != "" {
		sourceURL, err := value.NewURL(input.SourceURL)
		if err != nil {
			return nil, err
		}
		entry.UpdateSourceURL(&sourceURL)
	}
	if input.Memo != "" {
		entry.UpdateMemo(input.Memo)
	}

	if err := uc.repo.SaveEntryWithCompany(ctx, company, entry); err != nil {
		return nil, err
	}

	return &CreateWithCompanyOutput{
		Company: company,
		Entry:   entry,
	}, nil
}
