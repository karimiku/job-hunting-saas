package value

import (
	"errors"
	"strings"
)

// ErrCompanyNameEmpty は company name が空文字のときに返されるエラー。
// ErrCompanyNameInvalid は company name の形式が不正なときに返されるエラー。
var (
	ErrCompanyNameEmpty   = errors.New("company name must not be empty")
	ErrCompanyNameInvalid = errors.New("company name format is invalid")
)

// CompanyName は応募先企業の正式名称を表す値オブジェクト。
type CompanyName struct {
	value string
}

// NewCompanyName は raw から CompanyName を生成する。空文字や不正値は対応するエラーを返す。
func NewCompanyName(raw string) (CompanyName, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return CompanyName{}, ErrCompanyNameEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return CompanyName{}, ErrCompanyNameInvalid
	}
	return CompanyName{value: raw}, nil
}

// String は company name を文字列で返す。
func (n CompanyName) String() string {
	return n.value
}

// Equals は 2 つの CompanyName が等しいかを判定する。
func (n CompanyName) Equals(other CompanyName) bool {
	return n.value == other.value
}
