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

func (m *ESMemo) ID() ESMemoID { return m.id }

func (m *ESMemo) UserID() UserID { return m.userID }

func (m *ESMemo) EntryID() *EntryID { return m.entryID }

func (m *ESMemo) Category() value.ESMemoCategory { return m.category }

func (m *ESMemo) Title() value.ESMemoTitle { return m.title }

func (m *ESMemo) Content() value.ESMemoContent { return m.content }

func (m *ESMemo) Source() value.ESMemoSource { return m.source }

func (m *ESMemo) CreatedAt() time.Time { return m.createdAt }

func (m *ESMemo) UpdatedAt() time.Time { return m.updatedAt }
