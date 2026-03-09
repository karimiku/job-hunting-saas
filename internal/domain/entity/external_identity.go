package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// ExternalIdentity は外部認証プロバイダー（Google等）との連携情報を表すエンティティ。
// User とは N:1 の関係で、1ユーザーが複数のプロバイダーと連携できる。
// イミュータブル（作成後に変更しない）。
type ExternalIdentity struct {
	id        ExternalIdentityID
	userID    UserID
	provider  value.AuthProvider
	subject   string // プロバイダー側のユーザー識別子（Google の sub クレーム等）
	createdAt time.Time
}

func NewExternalIdentity(userID UserID, provider value.AuthProvider, subject string) *ExternalIdentity {
	return &ExternalIdentity{
		id:        NewExternalIdentityID(),
		userID:    userID,
		provider:  provider,
		subject:   subject,
		createdAt: time.Now(),
	}
}

// ReconstructExternalIdentity はDBから読み取ったデータでExternalIdentityを復元する。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructExternalIdentity(id ExternalIdentityID, userID UserID, provider value.AuthProvider, subject string, createdAt time.Time) *ExternalIdentity {
	return &ExternalIdentity{
		id:        id,
		userID:    userID,
		provider:  provider,
		subject:   subject,
		createdAt: createdAt,
	}
}

func (e *ExternalIdentity) ID() ExternalIdentityID   { return e.id }
func (e *ExternalIdentity) UserID() UserID            { return e.userID }
func (e *ExternalIdentity) Provider() value.AuthProvider { return e.provider }
func (e *ExternalIdentity) Subject() string           { return e.subject }
func (e *ExternalIdentity) CreatedAt() time.Time      { return e.createdAt }
