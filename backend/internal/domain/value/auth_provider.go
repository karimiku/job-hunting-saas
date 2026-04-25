package value

import "errors"

// ErrAuthProviderEmpty は auth provider が空文字のときに返されるエラー。
// ErrAuthProviderInvalid は auth provider が未定義の値のときに返されるエラー。
var (
	ErrAuthProviderEmpty   = errors.New("auth provider must not be empty")
	ErrAuthProviderInvalid = errors.New("auth provider is invalid")
)

const (
	authProviderGoogle = "google"
)

var validAuthProviders = map[string]bool{
	authProviderGoogle: true,
}

// AuthProvider は外部認証プロバイダーの種別を表す値オブジェクト。
type AuthProvider struct {
	value string
}

// NewAuthProvider は raw から AuthProvider を生成する。空文字や未定義値は対応するエラーを返す。
func NewAuthProvider(raw string) (AuthProvider, error) {
	if raw == "" {
		return AuthProvider{}, ErrAuthProviderEmpty
	}
	if !validAuthProviders[raw] {
		return AuthProvider{}, ErrAuthProviderInvalid
	}
	return AuthProvider{value: raw}, nil
}

// String は auth provider を文字列で返す。
func (p AuthProvider) String() string {
	return p.value
}

// Equals は 2 つの AuthProvider が等しいかを判定する。
func (p AuthProvider) Equals(other AuthProvider) bool {
	return p.value == other.value
}

// AuthProviderGoogle は Google プロバイダーを返す定数コンストラクタ。
func AuthProviderGoogle() AuthProvider { return AuthProvider{value: authProviderGoogle} }
