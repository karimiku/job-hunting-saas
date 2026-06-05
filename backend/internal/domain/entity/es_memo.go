package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// ESMemo は ES / 自己PR / ガクチカ / 面接ネタとして再利用する経験メモを表す。
// Entry への紐づけは任意で、ユーザー全体のナレッジとしても扱える。
type ESMemo struct {
	id        ESMemoID
	userID    UserID
	entryID   *EntryID
	category  value.ESMemoCategory
	title     value.ESMemoTitle
	content   value.ESMemoContent
	source    value.ESMemoSource
	createdAt time.Time
	updatedAt time.Time
}

// NewESMemo は新しいESメモを生成する。
func NewESMemo(
	userID UserID,
	entryID *EntryID,
	category value.ESMemoCategory,
	title value.ESMemoTitle,
	content value.ESMemoContent,
	source value.ESMemoSource,
) *ESMemo {
	now := time.Now()
	return &ESMemo{
		id:        NewESMemoID(),
		userID:    userID,
		entryID:   entryID,
		category:  category,
		title:     title,
		content:   content,
		source:    source,
		createdAt: now,
		updatedAt: now,
	}
}

// ReconstructESMemo は永続化済みデータからESメモを再構築する。
func ReconstructESMemo(
	id ESMemoID,
	userID UserID,
	entryID *EntryID,
	category value.ESMemoCategory,
	title value.ESMemoTitle,
	content value.ESMemoContent,
	source value.ESMemoSource,
	createdAt time.Time,
	updatedAt time.Time,
) *ESMemo {
	return &ESMemo{
		id:        id,
		userID:    userID,
		entryID:   entryID,
		category:  category,
		title:     title,
		content:   content,
		source:    source,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

// ID はESメモIDを返す。
func (m *ESMemo) ID() ESMemoID { return m.id }

// UserID は所有ユーザーIDを返す。
func (m *ESMemo) UserID() UserID { return m.userID }

// EntryID は紐づくエントリーIDを返す。未紐づけの場合はnil。
func (m *ESMemo) EntryID() *EntryID { return m.entryID }

// Category はESメモの分類を返す。
func (m *ESMemo) Category() value.ESMemoCategory { return m.category }

// Title はESメモの見出しを返す。
func (m *ESMemo) Title() value.ESMemoTitle { return m.title }

// Content はESメモ本文を返す。
func (m *ESMemo) Content() value.ESMemoContent { return m.content }

// Source はESメモの入力元を返す。
func (m *ESMemo) Source() value.ESMemoSource { return m.source }

// CreatedAt は作成日時を返す。
func (m *ESMemo) CreatedAt() time.Time { return m.createdAt }

// UpdatedAt は更新日時を返す。
func (m *ESMemo) UpdatedAt() time.Time { return m.updatedAt }
