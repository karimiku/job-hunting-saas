package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	inboxclipuc "github.com/karimiku/job-hunting-saas/internal/usecase/inbox_clip"
)

func setupInboxClipHandler() (*InboxClipHandler, *inmemory.InboxClipRepository) {
	repo := inmemory.NewInboxClipRepository()
	h := NewInboxClipHandler(
		inboxclipuc.NewCreate(repo),
		inboxclipuc.NewList(repo),
		inboxclipuc.NewDelete(repo),
	)
	return h, repo
}

func seedClip(t *testing.T, repo *inmemory.InboxClipRepository, userID entity.UserID) *entity.InboxClip {
	t.Helper()
	url, _ := value.NewURL("https://job.mynavi.jp/26/pc/search/corp123/outline.html")
	src, _ := value.NewSource("マイナビ")
	clip := entity.NewInboxClip(userID, url, "○○商事", src, "○○商事")
	if err := repo.Create(context.Background(), clip); err != nil {
		t.Fatalf("seed: %v", err)
	}
	return clip
}

func TestCreateInboxClip_Success(t *testing.T) {
	h, _ := setupInboxClipHandler()
	userID := entity.NewUserID()

	guess := "○○商事"
	body, _ := json.Marshal(openapi.CreateInboxClipRequest{
		Url:    "https://job.mynavi.jp/26/pc/search/corp123/outline.html",
		Title:  "○○商事 / 総合職",
		Source: "マイナビ",
		Guess:  &guess,
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.CreateInboxClip(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", w.Code, w.Body.String())
	}
	var resp openapi.InboxClipResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.Source != "マイナビ" {
		t.Errorf("Source = %q", resp.Source)
	}
	if resp.Guess != "○○商事" {
		t.Errorf("Guess = %q", resp.Guess)
	}
}

func TestCreateInboxClip_InvalidJSON(t *testing.T) {
	h, _ := setupInboxClipHandler()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("not json")))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.CreateInboxClip(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestCreateInboxClip_InvalidURL(t *testing.T) {
	h, _ := setupInboxClipHandler()
	body, _ := json.Marshal(openapi.CreateInboxClipRequest{
		Url:    "javascript:alert(1)",
		Title:  "x",
		Source: "マイナビ",
	})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.CreateInboxClip(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestListInboxClips_Success(t *testing.T) {
	h, repo := setupInboxClipHandler()
	userID := entity.NewUserID()
	seedClip(t, repo, userID)
	seedClip(t, repo, userID)
	seedClip(t, repo, entity.NewUserID()) // 他人

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.ListInboxClips(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d", w.Code)
	}
	var resp struct {
		Clips []openapi.InboxClipResponse `json:"clips"`
	}
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if len(resp.Clips) != 2 {
		t.Errorf("len = %d, want 2", len(resp.Clips))
	}
}

func TestDeleteInboxClip_Success(t *testing.T) {
	h, repo := setupInboxClipHandler()
	userID := entity.NewUserID()
	clip := seedClip(t, repo, userID)

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.DeleteInboxClip(w, req, openapi.ClipId(clip.ID()))

	if w.Code != http.StatusNoContent {
		t.Errorf("status = %d, want 204", w.Code)
	}
}

func TestDeleteInboxClip_OtherUser_NotFound(t *testing.T) {
	h, repo := setupInboxClipHandler()
	owner := entity.NewUserID()
	clip := seedClip(t, repo, owner)

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.DeleteInboxClip(w, req, openapi.ClipId(clip.ID()))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}
