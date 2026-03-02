package usecase

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type CreateCompanyInput struct {
	UserID entity.UserID
	Name   string
	Memo   string
}

type CreateCompanyOutput struct {
	Company *entity.Company
}

// CreateCompany は新しい企業を登録するUseCase。
type CreateCompany struct {
	companyRepo repository.CompanyRepository
}

func NewCreateCompany(companyRepo repository.CompanyRepository) *CreateCompany {
	return &CreateCompany{companyRepo: companyRepo}
}

// Execute は企業名をバリデーションし、新規Companyを生成して永続化する。
func (uc *CreateCompany) Execute(ctx context.Context, input CreateCompanyInput) (*CreateCompanyOutput, error) {
	name, err := value.NewCompanyName(input.Name)
	if err != nil {
		return nil, err
	}

	company := entity.NewCompany(input.UserID, name)

	if input.Memo != "" {
		company.UpdateMemo(input.Memo)
	}

	if err := uc.companyRepo.Save(ctx, company); err != nil {
		return nil, err
	}

	return &CreateCompanyOutput{Company: company}, nil
}
