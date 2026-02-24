package entity

import (
	"errors"
	"strings"
	"time"
)

var (
	ErrAliasEmpty = errors.New("alias must not be empty")
)

type CompanyAlias struct {
	id        CompanyAliasID
	userID    UserID
	companyID CompanyID
	alias     string
	createdAt time.Time
}

func NewCompanyAlias(userID UserID, companyID CompanyID, alias string) (*CompanyAlias, error) {
	if strings.TrimSpace(alias) == "" {
		return nil, ErrAliasEmpty
	}

	return &CompanyAlias{
		id:        NewID(),
		userID:    userID,
		companyID: companyID,
		alias:     alias,
		createdAt: time.Now(),
	}, nil
}

func (a *CompanyAlias) ID() CompanyAliasID  { return a.id }
func (a *CompanyAlias) UserID() UserID      { return a.userID }
func (a *CompanyAlias) CompanyID() CompanyID { return a.companyID }
func (a *CompanyAlias) Alias() string       { return a.alias }
func (a *CompanyAlias) CreatedAt() time.Time { return a.createdAt }
