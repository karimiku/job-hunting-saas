package selectionflow

import (
	"context"
	"fmt"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
)

// UpdateCurrentInput は現在ステージ更新ユースケースへの入力。
type UpdateCurrentInput struct {
	UserID   entity.UserID
	EntryID  entity.EntryID
	Position int
}

// UpdateCurrentOutput は現在ステージ更新ユースケースの出力。
type UpdateCurrentOutput struct {
	SelectionFlow *entity.SelectionFlow
	Entry         *entity.Entry
}

// UpdateCurrent は既存フローの現在ステージだけを更新する。
type UpdateCurrent struct {
	flowRepo  repository.SelectionFlowRepository
	entryRepo repository.EntryRepository
}

// NewUpdateCurrent は UpdateCurrent ユースケースを生成する。
func NewUpdateCurrent(flowRepo repository.SelectionFlowRepository, entryRepo repository.EntryRepository) *UpdateCurrent {
	return &UpdateCurrent{flowRepo: flowRepo, entryRepo: entryRepo}
}

// Execute は現在ステージ位置を更新し、Entryの互換ステージも同期する。
func (uc *UpdateCurrent) Execute(ctx context.Context, input UpdateCurrentInput) (*UpdateCurrentOutput, error) {
	if _, err := value.NewSelectionStagePosition(input.Position); err != nil {
		return nil, err
	}
	entry, err := uc.entryRepo.FindByID(ctx, input.UserID, input.EntryID)
	if err != nil {
		return nil, err
	}
	existing, err := uc.flowRepo.FindByEntryID(ctx, input.UserID, input.EntryID)
	if err != nil {
		return nil, err
	}
	stages := existing.Stages()
	stageInputs := make([]StageInput, 0, len(stages))
	for _, stage := range stages {
		stageInputs = append(stageInputs, StageInput{
			StageKind:    stage.Stage().Kind().String(),
			StageLabel:   stage.Stage().Label(),
			EvidenceText: stage.EvidenceText(),
		})
	}
	nextInput := UpsertInput{
		UserID:               input.UserID,
		EntryID:              input.EntryID,
		Source:               existing.Source().String(),
		CurrentStagePosition: input.Position,
		Confidence:           existing.Confidence(),
		InboxClipID:          existing.InboxClipID(),
		Stages:               stageInputs,
	}
	flow, err := buildSelectionFlow(nextInput)
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
	return &UpdateCurrentOutput{SelectionFlow: saved, Entry: entry}, nil
}
