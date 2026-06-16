package entity

import (
	"time"

	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// SelectionStage はEntryごとの選考フローに含まれる1ステージ。
type SelectionStage struct {
	id           SelectionStageID
	flowID       SelectionFlowID
	position     int
	stage        value.Stage
	evidenceText string
	createdAt    time.Time
}

// NewSelectionStage は SelectionStage を新規作成する。
func NewSelectionStage(flowID SelectionFlowID, position int, stage value.Stage, evidenceText string) *SelectionStage {
	return &SelectionStage{
		id:           NewSelectionStageID(),
		flowID:       flowID,
		position:     position,
		stage:        stage,
		evidenceText: evidenceText,
		createdAt:    time.Now(),
	}
}

// ReconstructSelectionStage は永続化データから SelectionStage を復元する。
func ReconstructSelectionStage(
	id SelectionStageID,
	flowID SelectionFlowID,
	position int,
	stage value.Stage,
	evidenceText string,
	createdAt time.Time,
) *SelectionStage {
	return &SelectionStage{
		id:           id,
		flowID:       flowID,
		position:     position,
		stage:        stage,
		evidenceText: evidenceText,
		createdAt:    createdAt,
	}
}

func (s *SelectionStage) ID() SelectionStageID    { return s.id }
func (s *SelectionStage) FlowID() SelectionFlowID { return s.flowID }
func (s *SelectionStage) Position() int           { return s.position }
func (s *SelectionStage) Stage() value.Stage      { return s.stage }
func (s *SelectionStage) EvidenceText() string    { return s.evidenceText }
func (s *SelectionStage) CreatedAt() time.Time    { return s.createdAt }

// SelectionFlow はEntryごとに可変な実選考フローを表す集約。
type SelectionFlow struct {
	id                   SelectionFlowID
	entryID              EntryID
	source               value.SelectionFlowSource
	currentStagePosition int
	confidence           *int
	inboxClipID          *InboxClipID
	stages               []*SelectionStage
	createdAt            time.Time
	updatedAt            time.Time
}

// NewSelectionFlow は SelectionFlow を新規作成する。
func NewSelectionFlow(
	entryID EntryID,
	source value.SelectionFlowSource,
	currentStagePosition int,
	confidence *int,
	inboxClipID *InboxClipID,
	stages []*SelectionStage,
) *SelectionFlow {
	now := time.Now()
	id := NewSelectionFlowID()
	copied := make([]*SelectionStage, 0, len(stages))
	for _, stage := range stages {
		if stage == nil {
			continue
		}
		copied = append(copied, NewSelectionStage(id, stage.Position(), stage.Stage(), stage.EvidenceText()))
	}
	return &SelectionFlow{
		id:                   id,
		entryID:              entryID,
		source:               source,
		currentStagePosition: currentStagePosition,
		confidence:           confidence,
		inboxClipID:          inboxClipID,
		stages:               copied,
		createdAt:            now,
		updatedAt:            now,
	}
}

// ReconstructSelectionFlow は永続化データから SelectionFlow を復元する。
func ReconstructSelectionFlow(
	id SelectionFlowID,
	entryID EntryID,
	source value.SelectionFlowSource,
	currentStagePosition int,
	confidence *int,
	inboxClipID *InboxClipID,
	stages []*SelectionStage,
	createdAt time.Time,
	updatedAt time.Time,
) *SelectionFlow {
	return &SelectionFlow{
		id:                   id,
		entryID:              entryID,
		source:               source,
		currentStagePosition: currentStagePosition,
		confidence:           confidence,
		inboxClipID:          inboxClipID,
		stages:               stages,
		createdAt:            createdAt,
		updatedAt:            updatedAt,
	}
}

func (f *SelectionFlow) ID() SelectionFlowID               { return f.id }
func (f *SelectionFlow) EntryID() EntryID                  { return f.entryID }
func (f *SelectionFlow) Source() value.SelectionFlowSource { return f.source }
func (f *SelectionFlow) CurrentStagePosition() int         { return f.currentStagePosition }
func (f *SelectionFlow) Confidence() *int                  { return f.confidence }
func (f *SelectionFlow) InboxClipID() *InboxClipID         { return f.inboxClipID }
func (f *SelectionFlow) Stages() []*SelectionStage {
	return append([]*SelectionStage(nil), f.stages...)
}
func (f *SelectionFlow) CreatedAt() time.Time { return f.createdAt }
func (f *SelectionFlow) UpdatedAt() time.Time { return f.updatedAt }

// CurrentStage は currentStagePosition に対応するステージを返す。
func (f *SelectionFlow) CurrentStage() *SelectionStage {
	for _, stage := range f.stages {
		if stage.Position() == f.currentStagePosition {
			return stage
		}
	}
	return nil
}
