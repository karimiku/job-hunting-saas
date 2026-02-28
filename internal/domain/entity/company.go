package entity

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrCompanyNameEmpty = errors.New("company name must not be empty")
)

// Company は応募先の企業を表すエンティティ。
// ユーザー単位で管理され、複数の Entry が紐づくマスターデータ。
type Company struct {
	id        CompanyID
	userID    UserID
	name      string
	memo      string
	createdAt time.Time
	updatedAt time.Time
}

func NewCompany(userID UserID, name string) (*Company, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, ErrCompanyNameEmpty
	}
	now := time.Now()
	return &Company{
		id:        NewCompanyID(),
		userID:    userID,
		name:      trimmed,
		memo:      "",
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructCompany はDBから読み取ったデータでCompanyを復元する。
// バリデーションをスキップする（永続化済みデータは検証済みの前提）。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructCompany(id CompanyID, userID UserID, name string, memo string, createdAt, updatedAt time.Time) *Company {
	return &Company{
		id:        id,
		userID:    userID,
		name:      name,
		memo:      memo,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (c *Company) ID() CompanyID      { return c.id }
func (c *Company) UserID() UserID      { return c.userID }
func (c *Company) Name() string        { return c.name }
func (c *Company) Memo() string        { return c.memo }
func (c *Company) CreatedAt() time.Time { return c.createdAt }
func (c *Company) UpdatedAt() time.Time { return c.updatedAt }

func (c *Company) Rename(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ErrCompanyNameEmpty
	}
	c.name = trimmed
	c.updatedAt = time.Now()
	return nil
}

func (c *Company) UpdateMemo(memo string) {
	c.memo = memo
	c.updatedAt = time.Now()
}
