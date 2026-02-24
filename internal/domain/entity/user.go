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
	if strings.TrimSpace(name) == "" {
		return nil, ErrUserNameEmpty
	}

	now := time.Now()
	return &User{
		id:        NewID(),
		email:     email,
		name:      name,
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (u *User) ID() UserID          { return u.id }
func (u *User) Email() value.Email  { return u.email }
func (u *User) Name() string        { return u.name }
func (u *User) CreatedAt() time.Time { return u.createdAt }
func (u *User) UpdatedAt() time.Time { return u.updatedAt }

func (u *User) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrUserNameEmpty
	}
	u.name = name
	u.updatedAt = time.Now()
	return nil
}

func (u *User) ChangeEmail(email value.Email) {
	u.email = email
	u.updatedAt = time.Now()
}
