package value

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrPasswordEmpty    = errors.New("password must not be empty")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")
)

const passwordMinLength = 8

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
	if len([]rune(raw)) < passwordMinLength {
		return Password{}, ErrPasswordTooShort
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
func (p Password) Verify(raw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(p.hash), []byte(raw)) == nil
}

// Hash はハッシュ値を返す。永続化時に使用する。
func (p Password) Hash() string {
	return p.hash
}
