package value

import (
	"errors"
	"strings"
)

var (
	ErrSourceEmpty   = errors.New("source must not be empty")
	ErrSourceInvalid = errors.New("source format is invalid")
)

// Source は応募経由の媒体を表す値オブジェクト。
// 定義済み選択肢（リクナビ・マイナビ・企業HP等）に加え、自由入力にも対応する。
type Source struct {
	value string
}

func NewSource(raw string) (Source, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return Source{}, ErrSourceEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return Source{}, ErrSourceInvalid
	}
	return Source{value: raw}, nil
}

func (s Source) String() string {
	return s.value
}

func (s Source) Equals(other Source) bool {
	return s.value == other.value
}
