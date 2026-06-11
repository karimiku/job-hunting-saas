package value

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"unicode"
	"unicode/utf8"
)

const (
	// AIAccessTokenNameMaxLength はトークン表示名の最大文字数（rune 数）。
	AIAccessTokenNameMaxLength = 80
	// AIAccessTokenRawPrefix はユーザーが見ても用途が分かるように付ける固定接頭辞。
	AIAccessTokenRawPrefix = "entre_ai_"
	// AIAccessTokenPrefixLength は一覧画面に出す識別用 prefix の長さ。
	AIAccessTokenPrefixLength = 18
)

var (
	// ErrAIAccessTokenNameEmpty はトークン表示名が空のときに返される。
	ErrAIAccessTokenNameEmpty = errors.New("AI access token name must not be empty")
	// ErrAIAccessTokenNameTooLong はトークン表示名が上限長を超えたときに返される。
	ErrAIAccessTokenNameTooLong = errors.New("AI access token name is too long")
	// ErrAIAccessTokenNameInvalid はトークン表示名に制御文字など不正な文字があるときに返される。
	ErrAIAccessTokenNameInvalid = errors.New("AI access token name format is invalid")
	// ErrAIAccessTokenInvalid は平文トークンの形式不正または検証失敗時に返される。
	ErrAIAccessTokenInvalid = errors.New("AI access token is invalid")
)

// AIAccessTokenName は AI / MCP 連携用トークンのユーザー向け表示名。
type AIAccessTokenName struct {
	value string
}

// NewAIAccessTokenName は raw を trim して表示名を生成する。
func NewAIAccessTokenName(raw string) (AIAccessTokenName, error) {
	v := strings.TrimSpace(raw)
	if v == "" {
		return AIAccessTokenName{}, ErrAIAccessTokenNameEmpty
	}
	if utf8.RuneCountInString(v) > AIAccessTokenNameMaxLength {
		return AIAccessTokenName{}, ErrAIAccessTokenNameTooLong
	}
	for _, r := range v {
		if unicode.IsControl(r) {
			return AIAccessTokenName{}, ErrAIAccessTokenNameInvalid
		}
	}
	return AIAccessTokenName{value: v}, nil
}

// String は表示名を文字列で返す。
func (n AIAccessTokenName) String() string { return n.value }

// AIAccessTokenHash は平文トークンを SHA-256 でハッシュ化した保存用値。
type AIAccessTokenHash struct {
	value string
}

// NewAIAccessTokenHash はDB復元用に hash 文字列を検証して生成する。
func NewAIAccessTokenHash(raw string) (AIAccessTokenHash, error) {
	if len(raw) != sha256.Size*2 {
		return AIAccessTokenHash{}, ErrAIAccessTokenInvalid
	}
	if _, err := hex.DecodeString(raw); err != nil {
		return AIAccessTokenHash{}, ErrAIAccessTokenInvalid
	}
	return AIAccessTokenHash{value: strings.ToLower(raw)}, nil
}

// NewAIAccessTokenHashFromRaw は平文トークンから保存用ハッシュを生成する。
func NewAIAccessTokenHashFromRaw(rawToken string) (AIAccessTokenHash, error) {
	if !strings.HasPrefix(rawToken, AIAccessTokenRawPrefix) {
		return AIAccessTokenHash{}, ErrAIAccessTokenInvalid
	}
	if len(rawToken) <= len(AIAccessTokenRawPrefix) {
		return AIAccessTokenHash{}, ErrAIAccessTokenInvalid
	}
	sum := sha256.Sum256([]byte(rawToken))
	return AIAccessTokenHash{value: hex.EncodeToString(sum[:])}, nil
}

// String は hash を文字列で返す。
func (h AIAccessTokenHash) String() string { return h.value }

// AIAccessTokenPrefix は一覧表示で識別するための平文トークン prefix。
type AIAccessTokenPrefix struct {
	value string
}

// NewAIAccessTokenPrefix はDB復元用に prefix を検証して生成する。
func NewAIAccessTokenPrefix(raw string) (AIAccessTokenPrefix, error) {
	if raw == "" || len(raw) > AIAccessTokenPrefixLength {
		return AIAccessTokenPrefix{}, ErrAIAccessTokenInvalid
	}
	if !strings.HasPrefix(raw, AIAccessTokenRawPrefix) {
		return AIAccessTokenPrefix{}, ErrAIAccessTokenInvalid
	}
	return AIAccessTokenPrefix{value: raw}, nil
}

// NewAIAccessTokenPrefixFromRaw は平文トークンから表示用 prefix を切り出す。
func NewAIAccessTokenPrefixFromRaw(rawToken string) (AIAccessTokenPrefix, error) {
	if len(rawToken) < AIAccessTokenPrefixLength {
		return AIAccessTokenPrefix{}, ErrAIAccessTokenInvalid
	}
	return NewAIAccessTokenPrefix(rawToken[:AIAccessTokenPrefixLength])
}

// String は prefix を文字列で返す。
func (p AIAccessTokenPrefix) String() string { return p.value }

// GenerateAIAccessTokenRaw はユーザーに一度だけ表示する平文トークンを生成する。
func GenerateAIAccessTokenRaw() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return AIAccessTokenRawPrefix + base64.RawURLEncoding.EncodeToString(buf), nil
}
