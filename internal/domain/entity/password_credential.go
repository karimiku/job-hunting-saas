package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// PasswordCredential はメールアドレス+パスワードによる認証情報を表すエンティティ。
// User とは 1:1 の関係で、パスワード認証を選択したユーザーのみ保持する。
type PasswordCredential struct {
	id        PasswordCredentialID
	userID    UserID
	password  value.Password
	createdAt time.Time
	updatedAt time.Time
}

func NewPasswordCredential(userID UserID, password value.Password) *PasswordCredential {
	now := time.Now()
	return &PasswordCredential{
		id:        NewPasswordCredentialID(),
		userID:    userID,
		password:  password,
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructPasswordCredential はDBから読み取ったデータでPasswordCredentialを復元する。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructPasswordCredential(id PasswordCredentialID, userID UserID, password value.Password, createdAt, updatedAt time.Time) *PasswordCredential {
	return &PasswordCredential{
		id:        id,
		userID:    userID,
		password:  password,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (c *PasswordCredential) ID() PasswordCredentialID { return c.id }
func (c *PasswordCredential) UserID() UserID           { return c.userID }
func (c *PasswordCredential) Password() value.Password { return c.password }
func (c *PasswordCredential) CreatedAt() time.Time     { return c.createdAt }
func (c *PasswordCredential) UpdatedAt() time.Time     { return c.updatedAt }

// ChangePassword はパスワードを変更する。
func (c *PasswordCredential) ChangePassword(password value.Password) {
	c.password = password
	c.updatedAt = time.Now()
}
