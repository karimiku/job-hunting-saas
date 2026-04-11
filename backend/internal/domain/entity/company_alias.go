package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// CompanyAlias は企業名の表記揺れを吸収するユーザー単位の別名辞書。
// イミュータブル（作成後に変更しない）。変更が必要な場合は削除→再作成する。
type CompanyAlias struct {
	id        CompanyAliasID
	userID    UserID
	companyID CompanyID
	alias     value.Alias
	createdAt time.Time
}

func NewCompanyAlias(userID UserID, companyID CompanyID, alias value.Alias) *CompanyAlias {
	return &CompanyAlias{
		id:        NewCompanyAliasID(),
		userID:    userID,
		companyID: companyID,
		alias:     alias,
		createdAt: time.Now(),
	}
}

// ReconstructCompanyAlias はDBから読み取ったデータでCompanyAliasを復元する。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructCompanyAlias(id CompanyAliasID, userID UserID, companyID CompanyID, alias value.Alias, createdAt time.Time) *CompanyAlias {
	return &CompanyAlias{
		id:        id,
		userID:    userID,
		companyID: companyID,
		alias:     alias,
		createdAt: createdAt,
	}
}

func (a *CompanyAlias) ID() CompanyAliasID  { return a.id }
func (a *CompanyAlias) UserID() UserID      { return a.userID }
func (a *CompanyAlias) CompanyID() CompanyID { return a.companyID }
func (a *CompanyAlias) Alias() value.Alias  { return a.alias }
func (a *CompanyAlias) CreatedAt() time.Time { return a.createdAt }
