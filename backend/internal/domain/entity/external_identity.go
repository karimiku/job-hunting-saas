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

// NewExternalIdentity は ExternalIdentity を新規作成する。各値オブジェクトのバリデーションは呼び出し側で済んでいる前提。
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

// ID は ExternalIdentity の ID を返す。
func (e *ExternalIdentity) ID() ExternalIdentityID { return e.id }

// UserID は ExternalIdentity に紐づくユーザの ID を返す。
func (e *ExternalIdentity) UserID() UserID { return e.userID }

// Provider は外部認証プロバイダーを返す。
func (e *ExternalIdentity) Provider() value.AuthProvider { return e.provider }

// Subject はプロバイダー側のユーザ識別子を返す。
func (e *ExternalIdentity) Subject() string { return e.subject }

// CreatedAt は ExternalIdentity の作成日時を返す。
func (e *ExternalIdentity) CreatedAt() time.Time { return e.createdAt }
