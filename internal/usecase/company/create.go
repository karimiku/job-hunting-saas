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

type Create struct {
	companyRepo repository.CompanyRepository
}

func NewCreate(companyRepo repository.CompanyRepository) *Create {
	return &Create{companyRepo: companyRepo}
}

func (uc *Create) Execute(ctx context.Context, input CreateInput) (*CreateOutput, error) {
	// 値オブジェクトの生成でドメインバリデーションを実行する
	companyName, err := value.NewCompanyName(input.Name)
	if err != nil {
		return nil, err
	}

	// Companyの必須不変条件はコンストラクタで閉じ、任意フィールドのmemoは生成後に反映する
	company := entity.NewCompany(input.UserID, companyName)

	if input.Memo != "" {
		company.UpdateMemo(input.Memo)
	}

	if err := uc.companyRepo.Save(ctx, company); err != nil {
		return nil, err
	}

	return &CreateOutput{Company: company}, nil
}
