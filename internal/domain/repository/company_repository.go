package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// CompanyRepository は企業情報の永続化を抽象化するインターフェース。
// Save は新規作成と更新の両方を処理する（upsert）。
// 参照系メソッドは userID によるマルチテナントスコーピングを必須とする。
type CompanyRepository interface {
	Save(ctx context.Context, company *entity.Company) error
	FindByID(ctx context.Context, userID entity.UserID, id entity.CompanyID) (*entity.Company, error)
	ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.Company, error)
	Delete(ctx context.Context, userID entity.UserID, id entity.CompanyID) error
}
