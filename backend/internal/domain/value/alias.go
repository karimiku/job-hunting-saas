package value

import (
	"errors"
	"strings"
)

var (
	ErrAliasEmpty   = errors.New("alias must not be empty")
	ErrAliasInvalid = errors.New("alias format is invalid")
)

// Alias は企業名の表記揺れを吸収するための別名を表す値オブジェクト。
type Alias struct {
	value string
}

func NewAlias(raw string) (Alias, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return Alias{}, ErrAliasEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return Alias{}, ErrAliasInvalid
	}
	return Alias{value: raw}, nil
}

func (a Alias) String() string {
	return a.value
}

func (a Alias) Equals(other Alias) bool {
	return a.value == other.value
}
