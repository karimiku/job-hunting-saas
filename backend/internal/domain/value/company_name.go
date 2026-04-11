package value

import (
	"errors"
	"strings"
)

var (
	ErrCompanyNameEmpty   = errors.New("company name must not be empty")
	ErrCompanyNameInvalid = errors.New("company name format is invalid")
)

// CompanyName は応募先企業の正式名称を表す値オブジェクト。
type CompanyName struct {
	value string
}

func NewCompanyName(raw string) (CompanyName, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return CompanyName{}, ErrCompanyNameEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return CompanyName{}, ErrCompanyNameInvalid
	}
	return CompanyName{value: raw}, nil
}

func (n CompanyName) String() string {
	return n.value
}

func (n CompanyName) Equals(other CompanyName) bool {
	return n.value == other.value
}
