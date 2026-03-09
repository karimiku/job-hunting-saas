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
	entryuc "github.com/karimiku/job-hunting-saas/internal/usecase/entry"
)

// setupEntryHandler はテスト用のEntryHandlerとリポジトリを初期化する。
func setupEntryHandler() (*EntryHandler, *inmemory.EntryRepository, *inmemory.CompanyRepository) {
	entryRepo := inmemory.NewEntryRepository()
	companyRepo := inmemory.NewCompanyRepository()

	h := NewEntryHandler(
		entryuc.NewCreate(entryRepo, companyRepo),
		entryuc.NewGet(entryRepo),
		entryuc.NewList(entryRepo),
		entryuc.NewUpdate(entryRepo),
		entryuc.NewDelete(entryRepo),
	)
	return h, entryRepo, companyRepo
}

// seedEntry はテスト用のEntryを作成して保存する。
func seedEntry(t *testing.T, entryRepo *inmemory.EntryRepository, companyRepo *inmemory.CompanyRepository, userID entity.UserID) *entity.Entry {
	t.Helper()

	companyName, err := value.NewCompanyName("テスト企業")
	if err != nil {
		t.Fatalf("NewCompanyName: %v", err)
	}
	company := entity.NewCompany(userID, companyName)
	if err := companyRepo.Save(context.Background(), company); err != nil {
		t.Fatalf("save company: %v", err)
	}

	route, err := value.NewRoute("本選考")
	if err != nil {
		t.Fatalf("NewRoute: %v", err)
	}
	source, err := value.NewSource("リクナビ")
	if err != nil {
		t.Fatalf("NewSource: %v", err)
	}
	entry := entity.NewEntry(userID, company.ID(), route, source)
	if err := entryRepo.Save(context.Background(), entry); err != nil {
		t.Fatalf("save entry: %v", err)
	}
	return entry
}

func TestUpdateEntry_PatchMerge_OnlySourceSent(t *testing.T) {
	h, entryRepo, companyRepo := setupEntryHandler()
	userID := entity.NewUserID()
	entry := seedEntry(t, entryRepo, companyRepo, userID)

	newSource := "マイナビ"
	body, _ := json.Marshal(openapi.UpdateEntryRequest{
		Source: &newSource,
	})

	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.UpdateEntry(w, req, openapi.EntryId(entry.ID()))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}

	var resp openapi.EntryResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	// 送信したフィールドだけ更新される
	if resp.Source != "マイナビ" {
		t.Errorf("Source = %q, want %q", resp.Source, "マイナビ")
	}
	// 未送信フィールドは既存値を維持する
	if resp.Status != "in_progress" {
		t.Errorf("Status = %q, want %q", resp.Status, "in_progress")
	}
	if resp.StageKind != "application" {
		t.Errorf("StageKind = %q, want %q", resp.StageKind, "application")
	}
	if resp.StageLabel != "応募" {
		t.Errorf("StageLabel = %q, want %q", resp.StageLabel, "応募")
	}
	if resp.Route != "本選考" {
		t.Errorf("Route = %q, want %q", resp.Route, "本選考")
	}
}

func TestUpdateEntry_PatchMerge_OnlyMemoSent(t *testing.T) {
	h, entryRepo, companyRepo := setupEntryHandler()
	userID := entity.NewUserID()
	entry := seedEntry(t, entryRepo, companyRepo, userID)

	newMemo := "更新メモ"
	body, _ := json.Marshal(openapi.UpdateEntryRequest{
		Memo: &newMemo,
	})

	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.UpdateEntry(w, req, openapi.EntryId(entry.ID()))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var resp openapi.EntryResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Memo != "更新メモ" {
		t.Errorf("Memo = %q, want %q", resp.Memo, "更新メモ")
	}
	// 他のフィールドは既存値を維持
	if resp.Source != "リクナビ" {
		t.Errorf("Source = %q, want %q", resp.Source, "リクナビ")
	}
}

func TestUpdateEntry_InvalidJSON(t *testing.T) {
	h, _, _ := setupEntryHandler()

	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader([]byte("invalid")))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.UpdateEntry(w, req, openapi.EntryId(entity.NewEntryID()))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestUpdateEntry_NotFound(t *testing.T) {
	h, _, _ := setupEntryHandler()

	body, _ := json.Marshal(openapi.UpdateEntryRequest{})
	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.UpdateEntry(w, req, openapi.EntryId(entity.NewEntryID()))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}
