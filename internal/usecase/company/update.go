package company

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type UpdateInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
	Name      *string
	Memo      *string
}

type UpdateOutput struct {
	Company *entity.Company
}

// Update はPATCHセマンティクスに対応する。
// nilフィールドは現在値を維持し、非nilフィールドのみ上書きする。
type Update struct {
	companyRepo repository.CompanyRepository
}

func NewUpdate(companyRepo repository.CompanyRepository) *Update {
	return &Update{companyRepo: companyRepo}
}

func (uc *Update) Execute(ctx context.Context, input UpdateInput) (*UpdateOutput, error) {
	company, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		validatedName, err := value.NewCompanyName(*input.Name)
		if err != nil {
			return nil, err
		}
		company.Rename(validatedName)
	}

	if input.Memo != nil {
		company.UpdateMemo(*input.Memo)
	}

	if err := uc.companyRepo.Save(ctx, company); err != nil {
		return nil, err
	}

	return &UpdateOutput{Company: company}, nil
}
