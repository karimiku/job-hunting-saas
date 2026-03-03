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
	Name      string
	Memo      string
}

type UpdateOutput struct {
	Company *entity.Company
}

// Update は既存企業の名前・メモを更新するUseCase。
type Update struct {
	companyRepo repository.CompanyRepository
}

func NewUpdate(companyRepo repository.CompanyRepository) *Update {
	return &Update{companyRepo: companyRepo}
}

// Execute は企業名をバリデーションし、既存Companyを取得して名前・メモを更新する。
func (uc *Update) Execute(ctx context.Context, input UpdateInput) (*UpdateOutput, error) {
	name, err := value.NewCompanyName(input.Name)
	if err != nil {
		return nil, err
	}

	c, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID)
	if err != nil {
		return nil, err
	}

	c.Rename(name)
	c.UpdateMemo(input.Memo)

	if err := uc.companyRepo.Save(ctx, c); err != nil {
		return nil, err
	}

	return &UpdateOutput{Company: c}, nil
}
