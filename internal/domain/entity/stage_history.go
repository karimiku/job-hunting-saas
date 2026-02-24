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

func (h *StageHistory) ID() StageHistoryID  { return h.id }
func (h *StageHistory) EntryID() EntryID    { return h.entryID }
func (h *StageHistory) Stage() value.Stage  { return h.stage }
func (h *StageHistory) Note() string        { return h.note }
func (h *StageHistory) CreatedAt() time.Time { return h.createdAt }
