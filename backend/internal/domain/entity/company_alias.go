// Package entity はドメインモデルのエンティティを定義する。
// 各エンティティは ID で同一性を持ち、内部状態はメソッド経由でのみ変更される。
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

// NewCompanyAlias は CompanyAlias を新規作成する。各値オブジェクトのバリデーションは呼び出し側で済んでいる前提。
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

// ID は CompanyAlias の ID を返す。
func (a *CompanyAlias) ID() CompanyAliasID { return a.id }

// UserID は CompanyAlias を所有するユーザの ID を返す。
func (a *CompanyAlias) UserID() UserID { return a.userID }

// CompanyID は別名の対象となる Company の ID を返す。
func (a *CompanyAlias) CompanyID() CompanyID { return a.companyID }

// Alias は別名を返す。
func (a *CompanyAlias) Alias() value.Alias { return a.alias }

// CreatedAt は CompanyAlias の作成日時を返す。
func (a *CompanyAlias) CreatedAt() time.Time { return a.createdAt }
