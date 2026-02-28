package value

import (
	"errors"
	"net/url"
	"strings"
)

var (
	ErrURLEmpty   = errors.New("url must not be empty")
	ErrURLInvalid = errors.New("url format is invalid")
)

// URL はHTTPS URLを表す値オブジェクト。
// 企業ページやマイページ等のリンク保存に使用する。
type URL struct {
	value string
}

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

func (u URL) String() string {
	return u.value
}

func (u URL) Equals(other URL) bool {
	return u.value == other.value
}
