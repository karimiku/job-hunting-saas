package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// AIAccessToken は Claude / Codex / MCP など外部AIクライアントからAPIを使うための認証トークン。
//
// 平文トークンは保持しない。保存するのは hash と一覧表示用 prefix のみ。
type AIAccessToken struct {
	id         AIAccessTokenID
	userID     UserID
	name       value.AIAccessTokenName
	tokenHash  value.AIAccessTokenHash
	prefix     value.AIAccessTokenPrefix
	createdAt  time.Time
	lastUsedAt *time.Time
	revokedAt  *time.Time
}

// NewAIAccessToken は新しい AI 連携トークンを生成する。
func NewAIAccessToken(
	userID UserID,
	name value.AIAccessTokenName,
	tokenHash value.AIAccessTokenHash,
	prefix value.AIAccessTokenPrefix,
) *AIAccessToken {
	return &AIAccessToken{
		id:        NewAIAccessTokenID(),
		userID:    userID,
		name:      name,
		tokenHash: tokenHash,
		prefix:    prefix,
		createdAt: time.Now(),
	}
}

// ReconstructAIAccessToken はDBから読み取ったデータで AIAccessToken を復元する。
func ReconstructAIAccessToken(
	id AIAccessTokenID,
	userID UserID,
	name value.AIAccessTokenName,
	tokenHash value.AIAccessTokenHash,
	prefix value.AIAccessTokenPrefix,
	createdAt time.Time,
	lastUsedAt *time.Time,
	revokedAt *time.Time,
) *AIAccessToken {
	return &AIAccessToken{
		id:         id,
		userID:     userID,
		name:       name,
		tokenHash:  tokenHash,
		prefix:     prefix,
		createdAt:  createdAt,
		lastUsedAt: lastUsedAt,
		revokedAt:  revokedAt,
	}
}

// ID はトークンIDを返す。
func (t *AIAccessToken) ID() AIAccessTokenID { return t.id }

// UserID は所有ユーザーIDを返す。
func (t *AIAccessToken) UserID() UserID { return t.userID }

// Name は表示名を返す。
func (t *AIAccessToken) Name() value.AIAccessTokenName { return t.name }

// TokenHash は保存用 hash を返す。
func (t *AIAccessToken) TokenHash() value.AIAccessTokenHash { return t.tokenHash }

// Prefix は一覧表示用 prefix を返す。
func (t *AIAccessToken) Prefix() value.AIAccessTokenPrefix { return t.prefix }

// CreatedAt は作成日時を返す。
func (t *AIAccessToken) CreatedAt() time.Time { return t.createdAt }

// LastUsedAt は最後に認証へ使われた日時を返す。
func (t *AIAccessToken) LastUsedAt() *time.Time { return t.lastUsedAt }

// RevokedAt は失効日時を返す。
func (t *AIAccessToken) RevokedAt() *time.Time { return t.revokedAt }

// IsRevoked はトークンが失効済みかを返す。
func (t *AIAccessToken) IsRevoked() bool { return t.revokedAt != nil }
