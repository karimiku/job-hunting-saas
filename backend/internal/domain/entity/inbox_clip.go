package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// InboxClip は Chrome 拡張等が保存した求人ページのクリップ。
// イミュータブル — エンティティへ変換される際に削除される設計。
type InboxClip struct {
	id         InboxClipID
	userID     UserID
	url        value.URL
	title      string
	source     value.Source
	guess      string // 推定された会社名（任意、空文字あり得る）
	capturedAt time.Time
}

// NewInboxClip は InboxClip を新規作成する。各値オブジェクトは呼び出し側でバリデーション済み前提。
func NewInboxClip(userID UserID, url value.URL, title string, source value.Source, guess string) *InboxClip {
	return &InboxClip{
		id:         NewInboxClipID(),
		userID:     userID,
		url:        url,
		title:      title,
		source:     source,
		guess:      guess,
		capturedAt: time.Now(),
	}
}

// ReconstructInboxClip はDBから読み取ったデータで InboxClip を復元する。
// Infra 層 (Repository 実装) からのみ呼び出すこと。
func ReconstructInboxClip(id InboxClipID, userID UserID, url value.URL, title string, source value.Source, guess string, capturedAt time.Time) *InboxClip {
	return &InboxClip{
		id:         id,
		userID:     userID,
		url:        url,
		title:      title,
		source:     source,
		guess:      guess,
		capturedAt: capturedAt,
	}
}

// ID は InboxClip の ID を返す。
func (c *InboxClip) ID() InboxClipID { return c.id }

// UserID はクリップを所有するユーザの ID を返す。
func (c *InboxClip) UserID() UserID { return c.userID }

// URL はクリップ元の URL を返す。
func (c *InboxClip) URL() value.URL { return c.url }

// Title はページタイトルを返す。
func (c *InboxClip) Title() string { return c.title }

// Source は応募媒体を返す。
func (c *InboxClip) Source() value.Source { return c.source }

// Guess は推定された会社名を返す（空文字あり）。
func (c *InboxClip) Guess() string { return c.guess }

// CapturedAt はクリップ作成日時を返す。
func (c *InboxClip) CapturedAt() time.Time { return c.capturedAt }
