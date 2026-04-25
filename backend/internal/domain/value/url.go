package value

import (
	"errors"
	"net/url"
	"strings"
)

// ErrURLEmpty は url が空文字のときに返されるエラー。
// ErrURLInvalid は url の形式が不正なときに返されるエラー。
var (
	ErrURLEmpty   = errors.New("url must not be empty")
	ErrURLInvalid = errors.New("url format is invalid")
)

// URL はHTTPS URLを表す値オブジェクト。
// 企業ページやマイページ等のリンク保存に使用する。
type URL struct {
	value string
}

// NewURL は raw から URL を生成する。空文字や https 以外のスキーム、不正値は対応するエラーを返す。
func NewURL(raw string) (URL, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return URL{}, ErrURLEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return URL{}, ErrURLInvalid
	}
	if !strings.HasPrefix(raw, "https://") {
		return URL{}, ErrURLInvalid
	}
	parsed, err := url.Parse(raw)
	if err != nil {
		return URL{}, ErrURLInvalid
	}
	if parsed.Host == "" {
		return URL{}, ErrURLInvalid
	}
	if strings.Contains(parsed.Host, " ") {
		return URL{}, ErrURLInvalid
	}
	return URL{value: raw}, nil
}

// String は url を文字列で返す。
func (u URL) String() string {
	return u.value
}

// Equals は 2 つの URL が等しいかを判定する。
func (u URL) Equals(other URL) bool {
	return u.value == other.value
}
