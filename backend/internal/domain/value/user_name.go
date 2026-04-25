package value

import (
	"errors"
	"strings"
)

// ErrUserNameEmpty は user name が空文字のときに返されるエラー。
// ErrUserNameInvalid は user name の形式が不正なときに返されるエラー。
var (
	ErrUserNameEmpty   = errors.New("user name must not be empty")
	ErrUserNameInvalid = errors.New("user name format is invalid")
)

// UserName はサービス利用者の表示名を表す値オブジェクト。
type UserName struct {
	value string
}

// NewUserName は raw から UserName を生成する。空文字や不正値は対応するエラーを返す。
func NewUserName(raw string) (UserName, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return UserName{}, ErrUserNameEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return UserName{}, ErrUserNameInvalid
	}
	return UserName{value: raw}, nil
}

// String は user name を文字列で返す。
func (n UserName) String() string {
	return n.value
}

// Equals は 2 つの UserName が等しいかを判定する。
func (n UserName) Equals(other UserName) bool {
	return n.value == other.value
}
