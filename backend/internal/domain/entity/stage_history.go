package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// StageHistory はEntryの選考フェーズ変更履歴。
// イミュータブル（作成後に変更しない）。updatedAt も持たない。
type StageHistory struct {
	id        StageHistoryID
	entryID   EntryID
	stage     value.Stage
	note      string
	createdAt time.Time
}

// NewStageHistory は StageHistory を新規作成する。各値オブジェクトのバリデーションは呼び出し側で済んでいる前提。
func NewStageHistory(entryID EntryID, stage value.Stage, note string) *StageHistory {
	return &StageHistory{
		id:        NewStageHistoryID(),
		entryID:   entryID,
		stage:     stage,
		note:      note,
		createdAt: time.Now(),
	}
}

// ReconstructStageHistory はDBから読み取ったデータでStageHistoryを復元する。
// Infra層（Repository実装）からのみ呼び出すこと。
func ReconstructStageHistory(id StageHistoryID, entryID EntryID, stage value.Stage, note string, createdAt time.Time) *StageHistory {
	return &StageHistory{
		id:        id,
		entryID:   entryID,
		stage:     stage,
		note:      note,
		createdAt: createdAt,
	}
}

// ID は StageHistory の ID を返す。
func (h *StageHistory) ID() StageHistoryID { return h.id }

// EntryID は履歴が紐づく Entry の ID を返す。
func (h *StageHistory) EntryID() EntryID { return h.entryID }

// Stage は遷移先の選考フェーズを返す。
func (h *StageHistory) Stage() value.Stage { return h.stage }

// Note は履歴に紐づくメモを返す。
func (h *StageHistory) Note() string { return h.note }

// CreatedAt は StageHistory の作成日時を返す。
func (h *StageHistory) CreatedAt() time.Time { return h.createdAt }
