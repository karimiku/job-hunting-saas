package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// User はサービス利用者を表すエンティティ。
// Google OAuth で認証されたアカウントに対応し、全データの所有者となる。
type User struct {
	id        UserID
	email     value.Email
	name      value.UserName
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(email value.Email, name value.UserName) *User {
	now := time.Now()
	return &User{
		id:        NewUserID(),
		email:     email,
		name:      name,
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructUser はDBから読み取ったデータでUserを復元する。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructUser(id UserID, email value.Email, name value.UserName, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		email:     email,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) ID() UserID           { return u.id }
func (u *User) Email() value.Email   { return u.email }
func (u *User) Name() value.UserName { return u.name }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

func (u *User) Rename(name value.UserName) {
	u.name = name
	u.updatedAt = time.Now()
}

func (u *User) ChangeEmail(email value.Email) {
	u.email = email
	u.updatedAt = time.Now()
}
