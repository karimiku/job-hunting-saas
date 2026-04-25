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
	companyuc "github.com/karimiku/job-hunting-saas/internal/usecase/company"
)

func setupCompanyHandler() (*CompanyHandler, *inmemory.CompanyRepository) {
	companyRepo := inmemory.NewCompanyRepository()

	h := NewCompanyHandler(
		companyuc.NewCreate(companyRepo),
		companyuc.NewGet(companyRepo),
		companyuc.NewList(companyRepo),
		companyuc.NewUpdate(companyRepo),
		companyuc.NewDelete(companyRepo),
	)
	return h, companyRepo
}

func seedCompany(t *testing.T, companyRepo *inmemory.CompanyRepository, userID entity.UserID, name string) *entity.Company {
	t.Helper()

	companyName, err := value.NewCompanyName(name)
	if err != nil {
		t.Fatalf("NewCompanyName: %v", err)
	}
	company := entity.NewCompany(userID, companyName)
	if err := companyRepo.Save(context.Background(), company); err != nil {
		t.Fatalf("save company: %v", err)
	}
	return company
}

// --- CreateCompany ---

func TestCreateCompany_Success(t *testing.T) {
	h, _ := setupCompanyHandler()
	userID := entity.NewUserID()

	memo := "気になる企業"
	body, _ := json.Marshal(openapi.CreateCompanyRequest{
		Name: "テスト企業",
		Memo: &memo,
	})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.CreateCompany(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201, body = %s", w.Code, w.Body.String())
	}

	var resp openapi.CompanyResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Name != "テスト企業" {
		t.Errorf("Name = %q, want %q", resp.Name, "テスト企業")
	}
	if resp.Memo != "気になる企業" {
		t.Errorf("Memo = %q, want %q", resp.Memo, "気になる企業")
	}
}

func TestCreateCompany_WithoutMemo(t *testing.T) {
	h, _ := setupCompanyHandler()
	body, _ := json.Marshal(openapi.CreateCompanyRequest{Name: "テスト企業"})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.CreateCompany(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201, body = %s", w.Code, w.Body.String())
	}
	var resp openapi.CompanyResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.Memo != "" {
		t.Errorf("Memo = %q, want empty", resp.Memo)
	}
}

func TestCreateCompany_InvalidJSON(t *testing.T) {
	h, _ := setupCompanyHandler()

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("not json")))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.CreateCompany(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestCreateCompany_EmptyName(t *testing.T) {
	h, _ := setupCompanyHandler()
	body, _ := json.Marshal(openapi.CreateCompanyRequest{Name: ""})

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.CreateCompany(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

// --- GetCompany ---

func TestGetCompany_Success(t *testing.T) {
	h, companyRepo := setupCompanyHandler()
	userID := entity.NewUserID()
	company := seedCompany(t, companyRepo, userID, "テスト企業")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.GetCompany(w, req, openapi.CompanyId(company.ID()))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp openapi.CompanyResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.Name != "テスト企業" {
		t.Errorf("Name = %q, want %q", resp.Name, "テスト企業")
	}
}

func TestGetCompany_NotFound(t *testing.T) {
	h, _ := setupCompanyHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.GetCompany(w, req, openapi.CompanyId(entity.NewCompanyID()))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestGetCompany_OtherUser(t *testing.T) {
	h, companyRepo := setupCompanyHandler()
	ownerID := entity.NewUserID()
	otherID := entity.NewUserID()
	company := seedCompany(t, companyRepo, ownerID, "他人の企業")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), otherID))
	w := httptest.NewRecorder()

	h.GetCompany(w, req, openapi.CompanyId(company.ID()))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404 (cross-user access blocked)", w.Code)
	}
}

// --- ListCompanies ---

func TestListCompanies_Success(t *testing.T) {
	h, companyRepo := setupCompanyHandler()
	userID := entity.NewUserID()
	seedCompany(t, companyRepo, userID, "企業A")
	seedCompany(t, companyRepo, userID, "企業B")
	seedCompany(t, companyRepo, entity.NewUserID(), "他人の企業")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.ListCompanies(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp struct {
		Companies []openapi.CompanyResponse `json:"companies"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Companies) != 2 {
		t.Errorf("len = %d, want 2 (other user's company excluded)", len(resp.Companies))
	}
}

func TestListCompanies_Empty(t *testing.T) {
	h, _ := setupCompanyHandler()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.ListCompanies(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp struct {
		Companies []openapi.CompanyResponse `json:"companies"`
	}
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if len(resp.Companies) != 0 {
		t.Errorf("len = %d, want 0", len(resp.Companies))
	}
}

// --- UpdateCompany ---

func TestUpdateCompany_PatchMerge_OnlyNameSent(t *testing.T) {
	h, companyRepo := setupCompanyHandler()
	userID := entity.NewUserID()
	company := seedCompany(t, companyRepo, userID, "旧名")

	newName := "新名"
	body, _ := json.Marshal(openapi.UpdateCompanyRequest{Name: &newName})

	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.UpdateCompany(w, req, openapi.CompanyId(company.ID()))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}
	var resp openapi.CompanyResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.Name != "新名" {
		t.Errorf("Name = %q, want %q", resp.Name, "新名")
	}
}

func TestUpdateCompany_PatchMerge_OnlyMemoSent(t *testing.T) {
	h, companyRepo := setupCompanyHandler()
	userID := entity.NewUserID()
	company := seedCompany(t, companyRepo, userID, "テスト企業")

	newMemo := "更新メモ"
	body, _ := json.Marshal(openapi.UpdateCompanyRequest{Memo: &newMemo})

	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.UpdateCompany(w, req, openapi.CompanyId(company.ID()))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}
	var resp openapi.CompanyResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)
	if resp.Memo != "更新メモ" {
		t.Errorf("Memo = %q, want %q", resp.Memo, "更新メモ")
	}
	// 未送信フィールドは保持
	if resp.Name != "テスト企業" {
		t.Errorf("Name = %q, want %q", resp.Name, "テスト企業")
	}
}

func TestUpdateCompany_InvalidJSON(t *testing.T) {
	h, _ := setupCompanyHandler()

	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader([]byte("not json")))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.UpdateCompany(w, req, openapi.CompanyId(entity.NewCompanyID()))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestUpdateCompany_NotFound(t *testing.T) {
	h, _ := setupCompanyHandler()

	body, _ := json.Marshal(openapi.UpdateCompanyRequest{})
	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.UpdateCompany(w, req, openapi.CompanyId(entity.NewCompanyID()))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

// --- DeleteCompany ---

func TestDeleteCompany_Success(t *testing.T) {
	h, companyRepo := setupCompanyHandler()
	userID := entity.NewUserID()
	company := seedCompany(t, companyRepo, userID, "テスト企業")

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.DeleteCompany(w, req, openapi.CompanyId(company.ID()))

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", w.Code)
	}
	if _, err := companyRepo.FindByID(context.Background(), userID, company.ID()); err == nil {
		t.Error("company should be deleted")
	}
}

func TestDeleteCompany_NotFound(t *testing.T) {
	h, _ := setupCompanyHandler()

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.DeleteCompany(w, req, openapi.CompanyId(entity.NewCompanyID()))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}

func TestDeleteCompany_OtherUser(t *testing.T) {
	h, companyRepo := setupCompanyHandler()
	ownerID := entity.NewUserID()
	otherID := entity.NewUserID()
	company := seedCompany(t, companyRepo, ownerID, "他人の企業")

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), otherID))
	w := httptest.NewRecorder()

	h.DeleteCompany(w, req, openapi.CompanyId(company.ID()))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404 (cross-user delete blocked)", w.Code)
	}
}
