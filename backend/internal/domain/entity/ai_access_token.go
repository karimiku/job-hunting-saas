package entity

import "time"

// AIAccessToken はAIクライアント連携用アクセストークンのメタデータを表す。
// tokenHash は平文トークンではなくSHA-256 hex digestのみを保持する。
type AIAccessToken struct {
	id           AIAccessTokenID
	userID       UserID
	name         string
	tokenHash    string
	tokenPreview string
	lastUsedAt   *time.Time
	revokedAt    *time.Time
	createdAt    time.Time
	updatedAt    time.Time
}

// NewAIAccessToken はAIアクセストークンを新規作成する。
func NewAIAccessToken(userID UserID, name, tokenHash, tokenPreview string) *AIAccessToken {
	now := time.Now()
	return &AIAccessToken{
		id:           NewAIAccessTokenID(),
		userID:       userID,
		name:         name,
		tokenHash:    tokenHash,
		tokenPreview: tokenPreview,
		createdAt:    now,
		updatedAt:    now,
	}
}

// ReconstructAIAccessToken はDBから読み取ったデータでAIAccessTokenを復元する。
func ReconstructAIAccessToken(
	id AIAccessTokenID,
	userID UserID,
	name string,
	tokenHash string,
	tokenPreview string,
	lastUsedAt *time.Time,
	revokedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *AIAccessToken {
	return &AIAccessToken{
		id:           id,
		userID:       userID,
		name:         name,
		tokenHash:    tokenHash,
		tokenPreview: tokenPreview,
		lastUsedAt:   lastUsedAt,
		revokedAt:    revokedAt,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

// ID はAIAccessTokenの識別子を返す。
func (t *AIAccessToken) ID() AIAccessTokenID { return t.id }

// UserID は所有ユーザーIDを返す。
func (t *AIAccessToken) UserID() UserID { return t.userID }

// Name はトークン名を返す。
func (t *AIAccessToken) Name() string { return t.name }

// TokenHash はトークンのハッシュを返す。
func (t *AIAccessToken) TokenHash() string { return t.tokenHash }

// TokenPreview はトークンの非認証用プレビューを返す。
func (t *AIAccessToken) TokenPreview() string { return t.tokenPreview }

// LastUsedAt は最終利用日時を返す。
func (t *AIAccessToken) LastUsedAt() *time.Time { return t.lastUsedAt }

// RevokedAt は失効日時を返す。
func (t *AIAccessToken) RevokedAt() *time.Time { return t.revokedAt }

// CreatedAt は作成日時を返す。
func (t *AIAccessToken) CreatedAt() time.Time { return t.createdAt }

// UpdatedAt は更新日時を返す。
func (t *AIAccessToken) UpdatedAt() time.Time { return t.updatedAt }

// MarkUsed は最終利用日時を更新する。
func (t *AIAccessToken) MarkUsed(at time.Time) {
	t.lastUsedAt = &at
	t.updatedAt = at
}

// Revoke はトークンを失効させる。
func (t *AIAccessToken) Revoke(at time.Time) {
	t.revokedAt = &at
	t.updatedAt = at
}
