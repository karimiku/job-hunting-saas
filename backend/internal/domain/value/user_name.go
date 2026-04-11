package value

import (
	"errors"
	"strings"
)

var (
	ErrUserNameEmpty   = errors.New("user name must not be empty")
	ErrUserNameInvalid = errors.New("user name format is invalid")
)

// UserName はサービス利用者の表示名を表す値オブジェクト。
type UserName struct {
	value string
}

func NewUserName(raw string) (UserName, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return UserName{}, ErrUserNameEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return UserName{}, ErrUserNameInvalid
	}
	return UserName{value: raw}, nil
}

func (n UserName) String() string {
	return n.value
}

func (n UserName) Equals(other UserName) bool {
	return n.value == other.value
}
