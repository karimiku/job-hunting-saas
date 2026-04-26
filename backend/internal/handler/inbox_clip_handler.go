package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	inboxclipuc "github.com/karimiku/job-hunting-saas/internal/usecase/inbox_clip"
)

// InboxClipHandler は Inbox クリップ関連の HTTP リクエストを受ける handler。
type InboxClipHandler struct {
	createUseCase *inboxclipuc.Create
	listUseCase   *inboxclipuc.List
	deleteUseCase *inboxclipuc.Delete
}

// NewInboxClipHandler は InboxClipHandler を DI 構築する。
func NewInboxClipHandler(
	createUC *inboxclipuc.Create,
	listUC *inboxclipuc.List,
	deleteUC *inboxclipuc.Delete,
) *InboxClipHandler {
	return &InboxClipHandler{
		createUseCase: createUC,
		listUseCase:   listUC,
		deleteUseCase: deleteUC,
	}
}

// CreateInboxClip は POST /api/v1/inbox/clips のハンドラ。
func (h *InboxClipHandler) CreateInboxClip(w http.ResponseWriter, r *http.Request) {
	var req openapi.CreateInboxClipRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
		return
	}

	guess := ""
	if req.Guess != nil {
		guess = *req.Guess
	}

	out, err := h.createUseCase.Execute(r.Context(), inboxclipuc.CreateInput{
		UserID: middleware.GetUserID(r.Context()),
		URL:    req.Url,
		Title:  req.Title,
		Source: req.Source,
		Guess:  guess,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, toInboxClipResponse(out.Clip))
}

// ListInboxClips は GET /api/v1/inbox/clips のハンドラ。
func (h *InboxClipHandler) ListInboxClips(w http.ResponseWriter, r *http.Request) {
	out, err := h.listUseCase.Execute(r.Context(), inboxclipuc.ListInput{
		UserID: middleware.GetUserID(r.Context()),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	items := make([]openapi.InboxClipResponse, len(out.Clips))
	for i, c := range out.Clips {
		items[i] = toInboxClipResponse(c)
	}
	writeJSON(w, http.StatusOK, map[string]any{"clips": items})
}

// DeleteInboxClip は DELETE /api/v1/inbox/clips/{clipId} のハンドラ。
func (h *InboxClipHandler) DeleteInboxClip(w http.ResponseWriter, r *http.Request, clipId openapi.ClipId) {
	err := h.deleteUseCase.Execute(r.Context(), inboxclipuc.DeleteInput{
		UserID: middleware.GetUserID(r.Context()),
		ClipID: entity.InboxClipID(clipId),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func toInboxClipResponse(c *entity.InboxClip) openapi.InboxClipResponse {
	return openapi.InboxClipResponse{
		Id:         uuid.UUID(c.ID()),
		Url:        c.URL().String(),
		Title:      c.Title(),
		Source:     c.Source().String(),
		Guess:      c.Guess(),
		CapturedAt: c.CapturedAt(),
	}
}
