package value

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
)

const (
	// AIAccessTokenPrefix はAIクライアント連携用トークンの公開prefix。
	AIAccessTokenPrefix  = "entre_ai_"
	aiAccessTokenBytes   = 32
	minAIAccessTokenBody = 32
)

var (
	// ErrAIAccessTokenEmpty はAIアクセストークンが空のときに返されるエラー。
	ErrAIAccessTokenEmpty = errors.New("ai access token must not be empty")
	// ErrAIAccessTokenInvalid はAIアクセストークン形式が不正なときに返されるエラー。
	ErrAIAccessTokenInvalid = errors.New("ai access token is invalid")
)

// AIAccessTokenSecret はAIクライアント連携用トークンの平文を表す値オブジェクト。
// 永続化するときは Hash() の結果だけを保存する。
type AIAccessTokenSecret struct {
	value string
}

// NewAIAccessTokenSecret はAIアクセストークンを検証して生成する。
func NewAIAccessTokenSecret(raw string) (AIAccessTokenSecret, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return AIAccessTokenSecret{}, ErrAIAccessTokenEmpty
	}
	if !strings.HasPrefix(trimmed, AIAccessTokenPrefix) {
		return AIAccessTokenSecret{}, ErrAIAccessTokenInvalid
	}
	body := strings.TrimPrefix(trimmed, AIAccessTokenPrefix)
	if len(body) < minAIAccessTokenBody {
		return AIAccessTokenSecret{}, ErrAIAccessTokenInvalid
	}
	for _, r := range body {
		if !isTokenBodyRune(r) {
			return AIAccessTokenSecret{}, ErrAIAccessTokenInvalid
		}
	}
	return AIAccessTokenSecret{value: trimmed}, nil
}

// GenerateAIAccessTokenSecret は新しいAIアクセストークンを生成する。
func GenerateAIAccessTokenSecret() (AIAccessTokenSecret, error) {
	bytes := make([]byte, aiAccessTokenBytes)
	if _, err := rand.Read(bytes); err != nil {
		return AIAccessTokenSecret{}, err
	}
	raw := AIAccessTokenPrefix + base64.RawURLEncoding.EncodeToString(bytes)
	return NewAIAccessTokenSecret(raw)
}

// String はトークン平文を返す。ログには出力しないこと。
func (t AIAccessTokenSecret) String() string { return t.value }

// Hash はトークン平文のSHA-256 hex digestを返す。
func (t AIAccessTokenSecret) Hash() string {
	sum := sha256.Sum256([]byte(t.value))
	return hex.EncodeToString(sum[:])
}

// Preview はUIやログで識別するための非認証用プレビューを返す。
func (t AIAccessTokenSecret) Preview() string {
	if len(t.value) <= len(AIAccessTokenPrefix)+8 {
		return AIAccessTokenPrefix + "..."
	}
	return t.value[:len(AIAccessTokenPrefix)+4] + "..." + t.value[len(t.value)-4:]
}

func isTokenBodyRune(r rune) bool {
	return (r >= 'a' && r <= 'z') ||
		(r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') ||
		r == '-' ||
		r == '_'
}
