package repository

import (
	"context"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
)

// EntryWithCompanyRepository は Company と Entry の同時作成を永続化する。
//
// SaveEntryWithCompany は両方の保存を同一トランザクションで扱い、片方だけが残る
// orphan 状態を作らない。
type EntryWithCompanyRepository interface {
	SaveEntryWithCompany(ctx context.Context, company *entity.Company, entry *entity.Entry) error
}
