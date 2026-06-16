// Package selectionflow はEntryごとの可変選考フローを扱うユースケース群を提供する。
package selectionflow

import (
	"context"
	"fmt"
	"strings"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

const maxSelectionStages = 20

// StageInput は選考フロー内の1ステージ入力。
type StageInput struct {
	StageKind    string
	StageLabel   string
	EvidenceText string
}

// UpsertInput は可変選考フロー保存ユースケースへの入力。
type UpsertInput struct {
	UserID               entity.UserID
	EntryID              entity.EntryID
	Source               string
	CurrentStagePosition int
	Confidence           *int
	InboxClipID          *entity.InboxClipID
	Stages               []StageInput
}

// UpsertOutput は可変選考フロー保存ユースケースの出力。
type UpsertOutput struct {
	SelectionFlow *entity.SelectionFlow
	Entry         *entity.Entry
}

// Upsert はEntryごとの選考フローを作成・置換し、Entryの現在ステージも同期する。
type Upsert struct {
	flowRepo  repository.SelectionFlowRepository
	entryRepo repository.EntryRepository
}

// NewUpsert は Upsert ユースケースを生成する。
func NewUpsert(flowRepo repository.SelectionFlowRepository, entryRepo repository.EntryRepository) *Upsert {
	return &Upsert{flowRepo: flowRepo, entryRepo: entryRepo}
}

// Execute はEntry所有権を確認し、可変フローを保存する。
func (uc *Upsert) Execute(ctx context.Context, input UpsertInput) (*UpsertOutput, error) {
	entry, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID)
	if err != nil {
		return nil, err
	}
	flow, err := buildSelectionFlow(input)
	if err != nil {
		return nil, err
	}
	current := flow.CurrentStage()
	if current == nil {
		return nil, fmt.Errorf("currentStagePosition does not match stages")
	}
	syncEntryStage(entry, current.Stage())
	if err := uc.entryRepo.Save(ctx, entry); err != nil {
		return nil, err
	}
	saved, err := uc.flowRepo.Upsert(ctx, flow)
	if err != nil {
		return nil, err
	}
	return &UpsertOutput{SelectionFlow: saved, Entry: entry}, nil
}

func buildSelectionFlow(input UpsertInput) (*entity.SelectionFlow, error) {
	source, err := value.NewSelectionFlowSource(input.Source)
	if err != nil {
		return nil, err
	}
	confidence, err := value.NewSelectionConfidence(input.Confidence)
	if err != nil {
		return nil, err
	}
	if len(input.Stages) == 0 {
		return nil, fmt.Errorf("stages are required")
	}
	if len(input.Stages) > maxSelectionStages {
		return nil, fmt.Errorf("stages must be %d or fewer", maxSelectionStages)
	}
	currentPosition := input.CurrentStagePosition
	if currentPosition == 0 {
		currentPosition = 1
	}
	if _, err := value.NewSelectionStagePosition(currentPosition); err != nil {
		return nil, err
	}
	stages := make([]*entity.SelectionStage, 0, len(input.Stages))
	hasCurrent := false
	for i, raw := range input.Stages {
		position := i + 1
		if position == currentPosition {
			hasCurrent = true
		}
		kind, err := value.NewStageKind(strings.TrimSpace(raw.StageKind))
		if err != nil {
			return nil, err
		}
		stage, err := value.NewStage(kind, strings.TrimSpace(raw.StageLabel))
		if err != nil {
			return nil, err
		}
		stages = append(stages, entity.NewSelectionStage(
			entity.SelectionFlowID{},
			position,
			stage,
			strings.TrimSpace(raw.EvidenceText),
		))
	}
	if !hasCurrent {
		return nil, fmt.Errorf("currentStagePosition must be within stages")
	}
	return entity.NewSelectionFlow(
		input.EntryID,
		source,
		currentPosition,
		confidence,
		input.InboxClipID,
		stages,
	), nil
}

func syncEntryStage(entry *entity.Entry, stage value.Stage) {
	entry.UpdateStage(stage)
	if stage.Kind().Equals(value.StageKindOffer()) {
		entry.UpdateStatus(value.EntryStatusOffered())
		return
	}
	if entry.Status().IsOpen() {
		entry.UpdateStatus(value.EntryStatusInProgress())
	}
}
