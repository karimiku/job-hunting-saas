package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

type StageHistory struct {
	id        StageHistoryID
	entryID   EntryID
	stage     value.Stage
	note      string
	createdAt time.Time
}

func NewStageHistory(entryID EntryID, stage value.Stage, note string) *StageHistory {
	return &StageHistory{
		id:        NewID(),
		entryID:   entryID,
		stage:     stage,
		note:      note,
		createdAt: time.Now(),
	}
}

func (h *StageHistory) ID() StageHistoryID  { return h.id }
func (h *StageHistory) EntryID() EntryID    { return h.entryID }
func (h *StageHistory) Stage() value.Stage  { return h.stage }
func (h *StageHistory) Note() string        { return h.note }
func (h *StageHistory) CreatedAt() time.Time { return h.createdAt }
