package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// CompanyAliasRepository は企業別名の永続化・復元を抽象化するインターフェース。
// domain層に定義することで、usecase層がインフラ実装に依存しない(DIP)。
// CompanyAlias はイミュータブルなエンティティのため Save は提供せず、Create のみとする。
// 別名の差し替えは Delete → Create で行う。
type CompanyAliasRepository interface {
	Create(ctx context.Context, companyAlias *entity.CompanyAlias) error
	FindByID(ctx context.Context, userID entity.UserID, companyAliasID entity.CompanyAliasID) (*entity.CompanyAlias, error)
	ListByCompanyID(ctx context.Context, userID entity.UserID, companyID entity.CompanyID) ([]*entity.CompanyAlias, error)
	Delete(ctx context.Context, userID entity.UserID, companyAliasID entity.CompanyAliasID) error
}
