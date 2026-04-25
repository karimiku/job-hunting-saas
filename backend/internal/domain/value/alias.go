package value

import (
	"errors"
	"strings"
)

// ErrAliasEmpty は alias が空文字のときに返されるエラー。
// ErrAliasInvalid は alias の形式が不正なときに返されるエラー。
var (
	ErrAliasEmpty   = errors.New("alias must not be empty")
	ErrAliasInvalid = errors.New("alias format is invalid")
)

// Alias は企業名の表記揺れを吸収するための別名を表す値オブジェクト。
type Alias struct {
	value string
}

// NewAlias は raw から Alias を生成する。空文字や不正値は対応するエラーを返す。
func NewAlias(raw string) (Alias, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return Alias{}, ErrAliasEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return Alias{}, ErrAliasInvalid
	}
	return Alias{value: raw}, nil
}

// String は alias を文字列で返す。
func (a Alias) String() string {
	return a.value
}

// Equals は 2 つの Alias が等しいかを判定する。
func (a Alias) Equals(other Alias) bool {
	return a.value == other.value
}
