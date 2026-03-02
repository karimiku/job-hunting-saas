package usecase

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type UpdateCompanyInput struct {
	UserID    entity.UserID
	CompanyID entity.CompanyID
	Name      string
	Memo      string
}

type UpdateCompanyOutput struct {
	Company *entity.Company
}

type UpdateCompany struct {
	companyRepo repository.CompanyRepository
}

func NewUpdateCompany(companyRepo repository.CompanyRepository) *UpdateCompany {
	return &UpdateCompany{companyRepo: companyRepo}
}

func (uc *UpdateCompany) Execute(ctx context.Context, input UpdateCompanyInput) (*UpdateCompanyOutput, error) {
	name, err := value.NewCompanyName(input.Name)
	if err != nil {
		return nil, err
	}

	company, err := uc.companyRepo.FindByID(ctx, input.UserID, input.CompanyID)
	if err != nil {
		return nil, err
	}

	company.Rename(name)
	company.UpdateMemo(input.Memo)

	if err := uc.companyRepo.Save(ctx, company); err != nil {
		return nil, err
	}

	return &UpdateCompanyOutput{Company: company}, nil
}
