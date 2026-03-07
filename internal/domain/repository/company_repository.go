package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// CompanyRepository は企業エンティティの永続化・復元を抽象化するインターフェース。
// domain層に定義することで、usecase層がインフラ実装に依存しない(DIP)。
// Save は新規作成と更新の両方を処理する（upsert）。存在判定の責務はusecase側にある。
// Save 以外のメソッドは、他ユーザーのデータへのアクセスを防ぐため userID でスコープする。
type CompanyRepository interface {
	Save(ctx context.Context, company *entity.Company) error
	FindByID(ctx context.Context, userID entity.UserID, companyID entity.CompanyID) (*entity.Company, error)
	ListByUserID(ctx context.Context, userID entity.UserID) ([]*entity.Company, error)
	Delete(ctx context.Context, userID entity.UserID, companyID entity.CompanyID) error
}
