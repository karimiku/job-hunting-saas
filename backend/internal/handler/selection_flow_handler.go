package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	selectionflowuc "github.com/karimiku/job-hunting-saas/internal/usecase/selection_flow"
)

// SelectionFlowHandler はEntryごとの可変選考フロー関連のHTTPリクエストを処理する。
type SelectionFlowHandler struct {
	getUseCase           *selectionflowuc.Get
	upsertUseCase        *selectionflowuc.Upsert
	updateCurrentUseCase *selectionflowuc.UpdateCurrent
}

// NewSelectionFlowHandler は SelectionFlowHandler を DI 構築する。
func NewSelectionFlowHandler(
	getUseCase *selectionflowuc.Get,
	upsertUseCase *selectionflowuc.Upsert,
	updateCurrentUseCase *selectionflowuc.UpdateCurrent,
) *SelectionFlowHandler {
	return &SelectionFlowHandler{
		getUseCase:           getUseCase,
		upsertUseCase:        upsertUseCase,
		updateCurrentUseCase: updateCurrentUseCase,
	}
}

// GetSelectionFlow は GET /api/v1/entries/{entryId}/selection-flow の handler。
func (h *SelectionFlowHandler) GetSelectionFlow(w http.ResponseWriter, r *http.Request, entryId openapi.EntryId) {
	out, err := h.getUseCase.Execute(r.Context(), selectionflowuc.GetInput{
		UserID:  middleware.GetUserID(r.Context()),
		EntryID: entity.EntryID(entryId),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toSelectionFlowResponse(out.SelectionFlow))
}

// UpsertSelectionFlow は PUT /api/v1/entries/{entryId}/selection-flow の handler。
func (h *SelectionFlowHandler) UpsertSelectionFlow(w http.ResponseWriter, r *http.Request, entryId openapi.EntryId) {
	var req openapi.UpsertSelectionFlowRequest
	if !decodeJSONBody(w, r, &req, maxDefaultJSONBodyBytes) {
		return
	}
	stages := make([]selectionflowuc.StageInput, 0, len(req.Stages))
	for _, stage := range req.Stages {
		evidenceText := ""
		if stage.EvidenceText != nil {
			evidenceText = *stage.EvidenceText
		}
		stages = append(stages, selectionflowuc.StageInput{
			StageKind:    string(stage.StageKind),
			StageLabel:   stage.StageLabel,
			EvidenceText: evidenceText,
		})
	}
	var inboxClipID *entity.InboxClipID
	if req.InboxClipId != nil {
		id := entity.InboxClipID(*req.InboxClipId)
		inboxClipID = &id
	}
	currentPosition := 0
	if req.CurrentStagePosition != nil {
		currentPosition = *req.CurrentStagePosition
	}
	out, err := h.upsertUseCase.Execute(r.Context(), selectionflowuc.UpsertInput{
		UserID:               middleware.GetUserID(r.Context()),
		EntryID:              entity.EntryID(entryId),
		Source:               string(req.Source),
		CurrentStagePosition: currentPosition,
		Confidence:           req.Confidence,
		InboxClipID:          inboxClipID,
		Stages:               stages,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toSelectionFlowResponse(out.SelectionFlow))
}

// UpdateSelectionFlowCurrentStage は PATCH /api/v1/entries/{entryId}/selection-flow/current-stage の handler。
func (h *SelectionFlowHandler) UpdateSelectionFlowCurrentStage(w http.ResponseWriter, r *http.Request, entryId openapi.EntryId) {
	var req openapi.UpdateSelectionFlowCurrentStageRequest
	if !decodeJSONBody(w, r, &req, maxDefaultJSONBodyBytes) {
		return
	}
	out, err := h.updateCurrentUseCase.Execute(r.Context(), selectionflowuc.UpdateCurrentInput{
		UserID:   middleware.GetUserID(r.Context()),
		EntryID:  entity.EntryID(entryId),
		Position: req.Position,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, toSelectionFlowResponse(out.SelectionFlow))
}

func toSelectionFlowResponse(flow *entity.SelectionFlow) openapi.SelectionFlowResponse {
	stages := make([]openapi.SelectionStageResponse, 0, len(flow.Stages()))
	for _, stage := range flow.Stages() {
		stages = append(stages, openapi.SelectionStageResponse{
			Id:           uuid.UUID(stage.ID()),
			Position:     stage.Position(),
			StageKind:    stage.Stage().Kind().String(),
			StageLabel:   stage.Stage().Label(),
			EvidenceText: stage.EvidenceText(),
		})
	}
	var inboxClipID *uuid.UUID
	if flow.InboxClipID() != nil {
		id := uuid.UUID(*flow.InboxClipID())
		inboxClipID = &id
	}
	return openapi.SelectionFlowResponse{
		Id:                   uuid.UUID(flow.ID()),
		EntryId:              uuid.UUID(flow.EntryID()),
		Source:               flow.Source().String(),
		CurrentStagePosition: flow.CurrentStagePosition(),
		Confidence:           flow.Confidence(),
		InboxClipId:          inboxClipID,
		Stages:               stages,
		CreatedAt:            flow.CreatedAt(),
		UpdatedAt:            flow.UpdatedAt(),
	}
}
