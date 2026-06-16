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

// ID は SelectionStage のIDを返す。
func (s *SelectionStage) ID() SelectionStageID { return s.id }

// FlowID は所属する SelectionFlow のIDを返す。
func (s *SelectionStage) FlowID() SelectionFlowID { return s.flowID }

// Position は選考フロー内の1始まり順序を返す。
func (s *SelectionStage) Position() int { return s.position }

// Stage は選考ステージの種別と表示名を返す。
func (s *SelectionStage) Stage() value.Stage { return s.stage }

// EvidenceText はステージ推定の根拠テキストを返す。
func (s *SelectionStage) EvidenceText() string { return s.evidenceText }

// CreatedAt は SelectionStage の作成日時を返す。
func (s *SelectionStage) CreatedAt() time.Time { return s.createdAt }

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

// ID は SelectionFlow のIDを返す。
func (f *SelectionFlow) ID() SelectionFlowID { return f.id }

// EntryID は紐づくEntryのIDを返す。
func (f *SelectionFlow) EntryID() EntryID { return f.entryID }

// Source は選考フローの作成元を返す。
func (f *SelectionFlow) Source() value.SelectionFlowSource { return f.source }

// CurrentStagePosition は現在ステージの1始まり順序を返す。
func (f *SelectionFlow) CurrentStagePosition() int { return f.currentStagePosition }

// Confidence はAI抽出時の信頼度を返す。
func (f *SelectionFlow) Confidence() *int { return f.confidence }

// InboxClipID はInbox由来の場合のClip IDを返す。
func (f *SelectionFlow) InboxClipID() *InboxClipID { return f.inboxClipID }

// Stages は選考フロー内のステージ一覧を返す。
func (f *SelectionFlow) Stages() []*SelectionStage {
	return append([]*SelectionStage(nil), f.stages...)
}

// CreatedAt は SelectionFlow の作成日時を返す。
func (f *SelectionFlow) CreatedAt() time.Time { return f.createdAt }

// UpdatedAt は SelectionFlow の更新日時を返す。
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
