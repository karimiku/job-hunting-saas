package entity

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrCompanyNameEmpty = errors.New("company name must not be empty")
)

type Company struct {
	id        CompanyID
	userID    UserID
	name      string
	memo      string
	createdAt time.Time
	updatedAt time.Time
}

func NewCompany(userID UserID, name string) (*Company, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrCompanyNameEmpty
	}
	now := time.Now()
	return &Company{
		id:        NewID(),
		userID:    userID,
		name:      name,
		memo:      "",
		createdAt: now,
		updatedAt: now,
	}, nil
}

func (c *Company) ID() CompanyID {
	return c.id
}

func (c *Company) UserID() UserID {
	return c.userID
}

func (c *Company) Name() string {
	return c.name
}

func (c *Company) Memo() string {
	return c.memo
}

func (c *Company) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Company) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Company) Rename(name string) error {
	if strings.TrimSpace(name) == "" {
		return ErrCompanyNameEmpty
	}
	c.name = name
	c.updatedAt = time.Now()
	return nil
}

func (c *Company) UpdateMemo(memo string) {
	c.memo = memo
	c.updatedAt = time.Now()
}
