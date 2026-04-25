package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	stagehistoryuc "github.com/karimiku/job-hunting-saas/internal/usecase/stage_history"
)

// StageHistoryHandler は oapi-codegen が生成した ServerInterface のStageHistory関連メソッドを実装する。
type StageHistoryHandler struct {
	createUseCase *stagehistoryuc.Create
	listUseCase   *stagehistoryuc.List
}

// NewStageHistoryHandler は StageHistoryHandler に必要なユースケース群を DI して新しい StageHistoryHandler を返す。
func NewStageHistoryHandler(
	createUseCase *stagehistoryuc.Create,
	listUseCase *stagehistoryuc.List,
) *StageHistoryHandler {
	return &StageHistoryHandler{
		createUseCase: createUseCase,
		listUseCase:   listUseCase,
	}
}

// CreateStageHistory は POST /entries/{entryId}/stage-histories の handler。
func (h *StageHistoryHandler) CreateStageHistory(w http.ResponseWriter, r *http.Request, entryId openapi.EntryId) {
	var req openapi.CreateStageHistoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
		return
	}

	note := ""
	if req.Note != nil {
		note = *req.Note
	}

	created, err := h.createUseCase.Execute(r.Context(), stagehistoryuc.CreateInput{
		UserID:    middleware.GetUserID(r.Context()),
		EntryID:   entity.EntryID(entryId),
		StageKind: string(req.StageKind),
		Label:     req.StageLabel,
		Note:      note,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toStageHistoryResponse(created.StageHistory))
}

// ListStageHistories は GET /entries/{entryId}/stage-histories の handler。
func (h *StageHistoryHandler) ListStageHistories(w http.ResponseWriter, r *http.Request, entryId openapi.EntryId) {
	result, err := h.listUseCase.Execute(r.Context(), stagehistoryuc.ListInput{
		UserID:  middleware.GetUserID(r.Context()),
		EntryID: entity.EntryID(entryId),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	items := make([]openapi.StageHistoryResponse, len(result.StageHistories))
	for i, history := range result.StageHistories {
		items[i] = toStageHistoryResponse(history)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"stageHistories": items,
	})
}

// toStageHistoryResponse はドメインエンティティをAPI応答用のDTOに変換する。
func toStageHistoryResponse(history *entity.StageHistory) openapi.StageHistoryResponse {
	return openapi.StageHistoryResponse{
		Id:         uuid.UUID(history.ID()),
		EntryId:    uuid.UUID(history.EntryID()),
		StageKind:  history.Stage().Kind().String(),
		StageLabel: history.Stage().Label(),
		Note:       history.Note(),
		CreatedAt:  history.CreatedAt(),
	}
}
