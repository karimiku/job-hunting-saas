package value

import (
	"errors"
	"strings"
	"unicode/utf8"
)

const (
	// InboxClipTitleMaxLength はクリップ元ページタイトルの最大文字数（rune 数）。
	InboxClipTitleMaxLength = 512
	// InboxClipGuessMaxLength は推定会社名の最大文字数（rune 数）。
	InboxClipGuessMaxLength = 256
	// InboxClipContentTextMaxLength は求人ページ本文スナップショットの最大文字数（rune 数）。
	InboxClipContentTextMaxLength = 20000
)

var (
	// ErrInboxClipTitleEmpty は title が空文字のときに返されるエラー。
	ErrInboxClipTitleEmpty = errors.New("inbox clip title must not be empty")
	// ErrInboxClipTitleInvalid は title の形式が不正なときに返されるエラー。
	ErrInboxClipTitleInvalid = errors.New("inbox clip title format is invalid")
	// ErrInboxClipTitleTooLong は title が上限長を超えたときに返されるエラー。
	ErrInboxClipTitleTooLong = errors.New("inbox clip title is too long")
	// ErrInboxClipGuessInvalid は guess の形式が不正なときに返されるエラー。
	ErrInboxClipGuessInvalid = errors.New("inbox clip guess format is invalid")
	// ErrInboxClipGuessTooLong は guess が上限長を超えたときに返されるエラー。
	ErrInboxClipGuessTooLong = errors.New("inbox clip guess is too long")
	// ErrInboxClipContentTextInvalid は contentText の形式が不正なときに返されるエラー。
	ErrInboxClipContentTextInvalid = errors.New("inbox clip content text format is invalid")
	// ErrInboxClipContentTextTooLong は contentText が上限長を超えたときに返されるエラー。
	ErrInboxClipContentTextTooLong = errors.New("inbox clip content text is too long")
)

// InboxClipTitle は Chrome 拡張等で保存したページタイトル。
type InboxClipTitle struct {
	value string
}

// NewInboxClipTitle は raw から InboxClipTitle を生成する。
func NewInboxClipTitle(raw string) (InboxClipTitle, error) {
	if raw == "" || strings.TrimSpace(raw) == "" {
		return InboxClipTitle{}, ErrInboxClipTitleEmpty
	}
	if raw != strings.TrimSpace(raw) {
		return InboxClipTitle{}, ErrInboxClipTitleInvalid
	}
	if utf8.RuneCountInString(raw) > InboxClipTitleMaxLength {
		return InboxClipTitle{}, ErrInboxClipTitleTooLong
	}
	return InboxClipTitle{value: raw}, nil
}

// String は title を文字列で返す。
func (t InboxClipTitle) String() string { return t.value }

// Equals は 2 つの InboxClipTitle が等しいかを判定する。
func (t InboxClipTitle) Equals(other InboxClipTitle) bool {
	return t.value == other.value
}

// InboxClipGuess はクリップから推定された会社名。空文字を許容する。
type InboxClipGuess struct {
	value string
}

// NewInboxClipGuess は raw から InboxClipGuess を生成する。
func NewInboxClipGuess(raw string) (InboxClipGuess, error) {
	if raw != strings.TrimSpace(raw) {
		return InboxClipGuess{}, ErrInboxClipGuessInvalid
	}
	if utf8.RuneCountInString(raw) > InboxClipGuessMaxLength {
		return InboxClipGuess{}, ErrInboxClipGuessTooLong
	}
	return InboxClipGuess{value: raw}, nil
}

// String は guess を文字列で返す。
func (g InboxClipGuess) String() string { return g.value }

// Equals は 2 つの InboxClipGuess が等しいかを判定する。
func (g InboxClipGuess) Equals(other InboxClipGuess) bool {
	return g.value == other.value
}

// InboxClipContentText は求人ページ本文のスナップショット。空文字を許容する。
type InboxClipContentText struct {
	value string
}

// NewInboxClipContentText は raw から InboxClipContentText を生成する。
func NewInboxClipContentText(raw string) (InboxClipContentText, error) {
	trimmed := strings.TrimSpace(raw)
	if utf8.RuneCountInString(trimmed) > InboxClipContentTextMaxLength {
		return InboxClipContentText{}, ErrInboxClipContentTextTooLong
	}
	if strings.ContainsRune(trimmed, '\x00') {
		return InboxClipContentText{}, ErrInboxClipContentTextInvalid
	}
	return InboxClipContentText{value: trimmed}, nil
}

// String は contentText を文字列で返す。
func (t InboxClipContentText) String() string { return t.value }
