package entity

import (
	"errors"
	"strings"
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

var (
	ErrUserNameEmpty = errors.New("user name must not be empty")
)

type User struct {
	id        UserID
	email     value.Email
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func NewUser(email value.Email, name string) (*User, error) {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return nil, ErrUserNameEmpty
	}

	now := time.Now()
	return &User{
		id:        NewUserID(),
		email:     email,
		name:      trimmed,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstructUser はDBから読み取ったデータでUserを復元する。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructUser(id UserID, email value.Email, name string, createdAt, updatedAt time.Time) *User {
	return &User{
		id:        id,
		email:     email,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

func (u *User) ID() UserID          { return u.id }
func (u *User) Email() value.Email  { return u.email }
func (u *User) Name() string        { return u.name }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

func (u *User) Rename(name string) error {
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return ErrUserNameEmpty
	}
	u.name = trimmed
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ChangeEmail(email value.Email) {
	u.email = email
	u.updatedAt = time.Now()
}
