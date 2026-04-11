package value

import "errors"

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

func NewAuthProvider(raw string) (AuthProvider, error) {
	if raw == "" {
		return AuthProvider{}, ErrAuthProviderEmpty
	}
	if !validAuthProviders[raw] {
		return AuthProvider{}, ErrAuthProviderInvalid
	}
	return AuthProvider{value: raw}, nil
}

func (p AuthProvider) String() string {
	return p.value
}

func (p AuthProvider) Equals(other AuthProvider) bool {
	return p.value == other.value
}

// AuthProviderGoogle は Google プロバイダーを返す定数コンストラクタ。
func AuthProviderGoogle() AuthProvider { return AuthProvider{value: authProviderGoogle} }
