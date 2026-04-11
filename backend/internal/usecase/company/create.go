package company

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type CreateInput struct {
	UserID entity.UserID
	Name   string
	Memo   string
}

type CreateOutput struct {
	Company *entity.Company
}

// Create は新しい企業を登録するUseCase。
type Create struct {
	companyRepo repository.CompanyRepository
}

func NewCreate(companyRepo repository.CompanyRepository) *Create {
	return &Create{companyRepo: companyRepo}
}

// Execute は企業名をバリデーションし、新規Companyを生成して永続化する。
func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	name, err := value.NewCompanyName(input.Name)
	if err != nil {
		return nil, err
	}

	c := entity.NewCompany(input.UserID, name)

	if input.Memo != "" {
		c.UpdateMemo(input.Memo)
	}

	if err := uc.companyRepo.Save(ctx, c); err != nil {
		return nil, err
	}

	return &CreateOutput{Company: c}, nil
}
