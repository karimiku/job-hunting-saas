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

// NewUser は User を新規作成する。
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

// ID は User の識別子を返す。
func (u *User) ID() UserID { return u.id }

// Email はメールアドレス値オブジェクトを返す。
func (u *User) Email() value.Email { return u.email }

// Name は表示名の値オブジェクトを返す。
func (u *User) Name() value.UserName { return u.name }

// CreatedAt は作成日時を返す。
func (u *User) CreatedAt() time.Time { return u.createdAt }

// UpdatedAt は最終更新日時を返す。
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

// Rename は user の表示名を更新し、UpdatedAt を現在時刻にする。
func (u *User) Rename(name value.UserName) {
	u.name = name
	u.updatedAt = time.Now()
}

// ChangeEmail は user のメールアドレスを更新し、UpdatedAt を現在時刻にする。
func (u *User) ChangeEmail(email value.Email) {
	u.email = email
	u.updatedAt = time.Now()
}
