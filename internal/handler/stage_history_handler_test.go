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
	stagehistoryuc "github.com/karimiku/job-hunting-saas/internal/usecase/stage_history"
)

// setupStageHistoryHandler はテスト用のStageHistoryHandlerとリポジトリを初期化する。
func setupStageHistoryHandler() (*StageHistoryHandler, *inmemory.StageHistoryRepository, *inmemory.EntryRepository, *inmemory.CompanyRepository) {
	companyRepo := inmemory.NewCompanyRepository()
	entryRepo := inmemory.NewEntryRepository()
	historyRepo := inmemory.NewStageHistoryRepository()

	h := NewStageHistoryHandler(
		stagehistoryuc.NewCreate(historyRepo, entryRepo),
		stagehistoryuc.NewList(historyRepo, entryRepo),
	)
	return h, historyRepo, entryRepo, companyRepo
}

func TestCreateStageHistory_Success(t *testing.T) {
	h, _, entryRepo, companyRepo := setupStageHistoryHandler()
	userID := entity.NewUserID()
	entry := seedEntry(t, entryRepo, companyRepo, userID)

	note := "一次面接通過"
	body, _ := json.Marshal(openapi.CreateStageHistoryRequest{
		StageKind:  openapi.CreateStageHistoryRequestStageKindInterview,
		StageLabel: "一次面接",
		Note:       &note,
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.CreateStageHistory(w, req, openapi.EntryId(entry.ID()))

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201, body = %s", w.Code, w.Body.String())
	}

	var resp openapi.StageHistoryResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.StageKind != "interview" {
		t.Errorf("StageKind = %q, want %q", resp.StageKind, "interview")
	}
	if resp.StageLabel != "一次面接" {
		t.Errorf("StageLabel = %q, want %q", resp.StageLabel, "一次面接")
	}
	if resp.Note != "一次面接通過" {
		t.Errorf("Note = %q, want %q", resp.Note, "一次面接通過")
	}
}

func TestCreateStageHistory_WithoutNote(t *testing.T) {
	h, _, entryRepo, companyRepo := setupStageHistoryHandler()
	userID := entity.NewUserID()
	entry := seedEntry(t, entryRepo, companyRepo, userID)

	body, _ := json.Marshal(openapi.CreateStageHistoryRequest{
		StageKind:  openapi.CreateStageHistoryRequestStageKindDocument,
		StageLabel: "ES提出",
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.CreateStageHistory(w, req, openapi.EntryId(entry.ID()))

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201, body = %s", w.Code, w.Body.String())
	}

	var resp openapi.StageHistoryResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Note != "" {
		t.Errorf("Note = %q, want empty", resp.Note)
	}
}

func TestCreateStageHistory_InvalidJSON(t *testing.T) {
	h, _, _, _ := setupStageHistoryHandler()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("invalid")))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.CreateStageHistory(w, req, openapi.EntryId(entity.NewEntryID()))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestCreateStageHistory_EntryNotFound(t *testing.T) {
	h, _, _, _ := setupStageHistoryHandler()

	body, _ := json.Marshal(openapi.CreateStageHistoryRequest{
		StageKind:  openapi.CreateStageHistoryRequestStageKindApplication,
		StageLabel: "応募",
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.CreateStageHistory(w, req, openapi.EntryId(entity.NewEntryID()))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestListStageHistories_Success(t *testing.T) {
	h, historyRepo, entryRepo, companyRepo := setupStageHistoryHandler()
	userID := entity.NewUserID()
	entry := seedEntry(t, entryRepo, companyRepo, userID)

	// 2件の履歴を作成
	stage1 := value.MustNewStage(value.StageKindApplication(), "応募")
	history1 := entity.NewStageHistory(entry.ID(), stage1, "応募完了")
	if err := historyRepo.Create(context.Background(), history1); err != nil {
		t.Fatalf("create history1: %v", err)
	}

	stage2 := value.MustNewStage(value.StageKindDocument(), "ES提出")
	history2 := entity.NewStageHistory(entry.ID(), stage2, "")
	if err := historyRepo.Create(context.Background(), history2); err != nil {
		t.Fatalf("create history2: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.ListStageHistories(w, req, openapi.EntryId(entry.ID()))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}

	var resp struct {
		StageHistories []openapi.StageHistoryResponse `json:"stageHistories"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if len(resp.StageHistories) != 2 {
		t.Errorf("len(stageHistories) = %d, want 2", len(resp.StageHistories))
	}
}

func TestListStageHistories_Empty(t *testing.T) {
	h, _, entryRepo, companyRepo := setupStageHistoryHandler()
	userID := entity.NewUserID()
	entry := seedEntry(t, entryRepo, companyRepo, userID)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.ListStageHistories(w, req, openapi.EntryId(entry.ID()))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var resp struct {
		StageHistories []openapi.StageHistoryResponse `json:"stageHistories"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.StageHistories == nil {
		t.Errorf("stageHistories should be empty array, got nil")
	}
}

func TestListStageHistories_EntryNotFound(t *testing.T) {
	h, _, _, _ := setupStageHistoryHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.ListStageHistories(w, req, openapi.EntryId(entity.NewEntryID()))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}
