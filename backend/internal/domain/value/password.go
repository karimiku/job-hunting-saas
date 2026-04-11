package value

import (
	"errors"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordEmpty    = errors.New("password must not be empty")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong  = errors.New("password must not exceed 72 bytes")
)

const (
	passwordMinLength = 8
	passwordMaxBytes  = 72 // bcryptの上限と一致させるが、ドメインルールとして定義
)

// Password はハッシュ化済みのパスワードを保持する値オブジェクト。
// 平文パスワードは保持しない。
type Password struct {
	hash string
}

// NewPassword は平文パスワードをバリデーションし、bcryptでハッシュ化して生成する。
func NewPassword(raw string) (Password, error) {
	if raw == "" {
		return Password{}, ErrPasswordEmpty
	}
	if utf8.RuneCountInString(raw) < passwordMinLength {
		return Password{}, ErrPasswordTooShort
	}
	if len(raw) > passwordMaxBytes {
		return Password{}, ErrPasswordTooLong
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, err
	}
	return Password{hash: string(hashed)}, nil
}

// ReconstructPassword はDBから読み取ったハッシュ値でPasswordを復元する。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructPassword(hash string) Password {
	return Password{hash: hash}
}

// Verify は平文パスワードがハッシュと一致するか検証する。
// 生成時と同じ最大長ルールを適用し、72バイト超は常にfalseを返す。
func (p Password) Verify(raw string) bool {
	if len(raw) > passwordMaxBytes {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(raw)) == nil
}

// Hash はハッシュ値を返す。永続化時に使用する。
func (p Password) Hash() string {
	return p.hash
}
