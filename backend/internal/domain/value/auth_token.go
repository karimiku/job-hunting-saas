package value

import "errors"

// ErrAuthTokenInvalid は認証トークンが無効なときに返されるエラー。
var ErrAuthTokenInvalid = errors.New("auth token is invalid")
