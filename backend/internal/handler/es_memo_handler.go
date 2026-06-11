package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	esmemo "github.com/karimiku/job-hunting-saas/internal/usecase/es_memo"
)

// ESMemoHandler は ES / 自己PR / 面接ネタ用メモの HTTP リクエストを受ける。
type ESMemoHandler struct {
	appendUseCase *esmemo.Append
	listUseCase   *esmemo.List
}

// NewESMemoHandler は ESMemoHandler を DI 構築する。
func NewESMemoHandler(appendUC *esmemo.Append, listUC *esmemo.List) *ESMemoHandler {
	return &ESMemoHandler{
		appendUseCase: appendUC,
		listUseCase:   listUC,
	}
}

const maxESMemoBodyBytes = 256 * 1024

// CreateEsMemo は POST /api/v1/es-memos のハンドラ。
func (h *ESMemoHandler) CreateEsMemo(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxESMemoBodyBytes)

	var req openapi.CreateEsMemoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		var maxBytesErr *http.MaxBytesError
		if errors.As(err, &maxBytesErr) {
			writeJSON(w, http.StatusRequestEntityTooLarge, openapi.ErrorResponse{Message: "request body too large"})
			return
		}
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
		return
	}

	var entryID *entity.EntryID
	if req.EntryId != nil {
		id := entity.EntryID(*req.EntryId)
		entryID = &id
	}
	category := ""
	if req.Category != nil {
		category = *req.Category
	}
	source := ""
	if req.Source != nil {
		source = *req.Source
	}

	out, err := h.appendUseCase.Execute(r.Context(), esmemo.AppendInput{
		UserID:   middleware.GetUserID(r.Context()),
		EntryID:  entryID,
		Category: category,
		Title:    req.Title,
		Content:  req.Content,
		Source:   source,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toEsMemoResponse(out.Memo))
}

// ListEsMemos は GET /api/v1/es-memos のハンドラ。
func (h *ESMemoHandler) ListEsMemos(w http.ResponseWriter, r *http.Request, params openapi.ListEsMemosParams) {
	limit := int32(0)
	if params.Limit != nil {
		limit = int32(*params.Limit)
	}
	out, err := h.listUseCase.Execute(r.Context(), esmemo.ListInput{
		UserID: middleware.GetUserID(r.Context()),
		Limit:  limit,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	items := make([]openapi.EsMemoResponse, len(out.Memos))
	for i, memo := range out.Memos {
		items[i] = toEsMemoResponse(memo)
	}
	writeJSON(w, http.StatusOK, map[string]any{"memos": items})
}

func toEsMemoResponse(memo *entity.ESMemo) openapi.EsMemoResponse {
	var entryID *uuid.UUID
	if memo.EntryID() != nil {
		id := uuid.UUID(*memo.EntryID())
		entryID = &id
	}
	return openapi.EsMemoResponse{
		Id:        uuid.UUID(memo.ID()),
		EntryId:   entryID,
		Category:  memo.Category().String(),
		Title:     memo.Title().String(),
		Content:   memo.Content().String(),
		Source:    memo.Source().String(),
		CreatedAt: memo.CreatedAt(),
		UpdatedAt: memo.UpdatedAt(),
	}
}
