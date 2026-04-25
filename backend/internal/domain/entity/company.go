package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// Company は応募先の企業を表すエンティティ。
// ユーザー単位で管理され、複数の Entry が紐づくマスターデータ。
type Company struct {
	id        CompanyID
	userID    UserID
	name      value.CompanyName
	memo      string
	createdAt time.Time
	updatedAt time.Time
}

// NewCompany は Company を新規作成する。各値オブジェクトのバリデーションは呼び出し側で済んでいる前提。
func NewCompany(userID UserID, name value.CompanyName) *Company {
	now := time.Now()
	return &Company{
		id:        NewCompanyID(),
		userID:    userID,
		name:      name,
		memo:      "",
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructCompany はDBから読み取ったデータでCompanyを復元する。
// バリデーションをスキップする（永続化済みデータは検証済みの前提）。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructCompany(id CompanyID, userID UserID, name value.CompanyName, memo string, createdAt, updatedAt time.Time) *Company {
	return &Company{
		id:        id,
		userID:    userID,
		name:      name,
		memo:      memo,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID は Company の ID を返す。
func (c *Company) ID() CompanyID { return c.id }

// UserID は Company を所有するユーザの ID を返す。
func (c *Company) UserID() UserID { return c.userID }

// Name は Company の企業名を返す。
func (c *Company) Name() value.CompanyName { return c.name }

// Memo は Company のメモを返す。
func (c *Company) Memo() string { return c.memo }

// CreatedAt は Company の作成日時を返す。
func (c *Company) CreatedAt() time.Time { return c.createdAt }

// UpdatedAt は Company の更新日時を返す。
func (c *Company) UpdatedAt() time.Time { return c.updatedAt }

// Rename は企業名を更新し、UpdatedAt を現在時刻にする。
func (c *Company) Rename(name value.CompanyName) {
	c.name = name
	c.updatedAt = time.Now()
}

// UpdateMemo は memo を更新し、UpdatedAt を現在時刻にする。
func (c *Company) UpdateMemo(memo string) {
	c.memo = memo
	c.updatedAt = time.Now()
}
