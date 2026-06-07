package value

import (
	"errors"
	"strings"
	"unicode/utf8"
)

// SourceMaxLength は source の最大文字数（rune 数）。
const SourceMaxLength = 128

// ErrSourceEmpty は source が空文字のときに返されるエラー。
// ErrSourceInvalid は source の形式が不正なときに返されるエラー。
// ErrSourceTooLong は source が上限長を超えたときに返されるエラー。
var (
	ErrSourceEmpty   = errors.New("source must not be empty")
	ErrSourceInvalid = errors.New("source format is invalid")
	ErrSourceTooLong = errors.New("source is too long")
)

// Source は応募経由の媒体を表す値オブジェクト。
// 定義済み選択肢（リクナビ・マイナビ・企業HP等）に加え、自由入力にも対応する。
type Source struct {
	value string
}

// NewSource は raw から Source を生成する。空文字や不正値は対応するエラーを返す。
func NewSource(raw string) (Source, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return Source{}, ErrSourceEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return Source{}, ErrSourceInvalid
	}
	if utf8.RuneCountInString(raw) > SourceMaxLength {
		return Source{}, ErrSourceTooLong
	}
	return Source{value: raw}, nil
}

// String は source を文字列で返す。
func (s Source) String() string {
	return s.value
}

// Equals は 2 つの Source が等しいかを判定する。
func (s Source) Equals(other Source) bool {
	return s.value == other.value
}
