package value

import (
	"errors"
	"regexp"
	"strings"
)

// ErrEmailEmpty は email が空文字のときに返されるエラー。
// ErrEmailInvalid は email の形式が不正なときに返されるエラー。
var (
	ErrEmailEmpty   = errors.New("email must not be empty")
	ErrEmailInvalid = errors.New("email format is invalid")

	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9]+([.+_-][a-zA-Z0-9]+)*@[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?)*\.[a-zA-Z]{2,}$`)
)

// Email はメールアドレスを表す値オブジェクト。
// Google OAuth 認証で取得したアドレスをユーザー識別に使用する。
type Email struct {
	value string
}

// NewEmail は raw から Email を生成する。空文字や不正値は対応するエラーを返す。
func NewEmail(raw string) (Email, error) {
	if raw == "" {
		return Email{}, ErrEmailEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return Email{}, ErrEmailInvalid
	}
	if !emailRegex.MatchString(raw) {
		return Email{}, ErrEmailInvalid
	}
	return Email{value: strings.ToLower(raw)}, nil
}

// String は email を文字列で返す。
func (e Email) String() string {
	return e.value
}

// Equals は 2 つの Email が等しいかを判定する。
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
