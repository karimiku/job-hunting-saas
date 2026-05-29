package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/value"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	companyaliasuc "github.com/karimiku/job-hunting-saas/internal/usecase/company_alias"
)

func setupCompanyAliasHandler() (*CompanyAliasHandler, *inmemory.CompanyRepository, *inmemory.CompanyAliasRepository) {
	companyRepo := inmemory.NewCompanyRepository()
	aliasRepo := inmemory.NewCompanyAliasRepository()

	h := NewCompanyAliasHandler(
		companyaliasuc.NewCreate(aliasRepo, companyRepo),
		companyaliasuc.NewGet(aliasRepo),
		companyaliasuc.NewList(aliasRepo),
		companyaliasuc.NewDelete(aliasRepo),
	)
	return h, companyRepo, aliasRepo
}

func seedCompanyForAlias(t *testing.T, companyRepo *inmemory.CompanyRepository, userID entity.UserID) *entity.Company {
	t.Helper()
	name, err := value.NewCompanyName("テスト企業")
	if err != nil {
		t.Fatalf("NewCompanyName: %v", err)
	}
	company := entity.NewCompany(userID, name)
	if err := companyRepo.Save(context.Background(), company); err != nil {
		t.Fatalf("save company: %v", err)
	}
	return company
}

func seedAlias(t *testing.T, aliasRepo *inmemory.CompanyAliasRepository, userID entity.UserID, companyID entity.CompanyID, raw string) *entity.CompanyAlias {
	t.Helper()
	alias, err := value.NewAlias(raw)
	if err != nil {
		t.Fatalf("NewAlias: %v", err)
	}
	a := entity.NewCompanyAlias(userID, companyID, alias)
	if err := aliasRepo.Create(context.Background(), a); err != nil {
		t.Fatalf("create alias: %v", err)
	}
	return a
}

func TestCreateCompanyAlias_Success(t *testing.T) {
	h, companyRepo, _ := setupCompanyAliasHandler()
	userID := entity.NewUserID()
	company := seedCompanyForAlias(t, companyRepo, userID)

	body, _ := json.Marshal(openapi.CreateCompanyAliasRequest{Alias: "トヨタ"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.CreateCompanyAlias(w, req, openapi.CompanyId(uuid.UUID(company.ID())))

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want 201, body = %s", w.Code, w.Body.String())
	}
	var resp openapi.CompanyAliasResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Alias != "トヨタ" {
		t.Errorf("Alias = %q, want %q", resp.Alias, "トヨタ")
	}
	if resp.CompanyId != uuid.UUID(company.ID()) {
		t.Errorf("CompanyId = %v, want %v", resp.CompanyId, uuid.UUID(company.ID()))
	}
}

func TestCreateCompanyAlias_CompanyNotFound(t *testing.T) {
	h, _, _ := setupCompanyAliasHandler()
	userID := entity.NewUserID()

	body, _ := json.Marshal(openapi.CreateCompanyAliasRequest{Alias: "トヨタ"})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.CreateCompanyAlias(w, req, openapi.CompanyId(uuid.New()))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404, body = %s", w.Code, w.Body.String())
	}
}

func TestCreateCompanyAlias_EmptyAlias(t *testing.T) {
	h, companyRepo, _ := setupCompanyAliasHandler()
	userID := entity.NewUserID()
	company := seedCompanyForAlias(t, companyRepo, userID)

	body, _ := json.Marshal(openapi.CreateCompanyAliasRequest{Alias: ""})
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.CreateCompanyAlias(w, req, openapi.CompanyId(uuid.UUID(company.ID())))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400, body = %s", w.Code, w.Body.String())
	}
}

func TestListCompanyAliases_Success(t *testing.T) {
	h, companyRepo, aliasRepo := setupCompanyAliasHandler()
	userID := entity.NewUserID()
	company := seedCompanyForAlias(t, companyRepo, userID)
	seedAlias(t, aliasRepo, userID, company.ID(), "トヨタ")
	seedAlias(t, aliasRepo, userID, company.ID(), "TOYOTA")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.ListCompanyAliases(w, req, openapi.CompanyId(uuid.UUID(company.ID())))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}
	var resp struct {
		Aliases []openapi.CompanyAliasResponse `json:"aliases"`
	}
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Aliases) != 2 {
		t.Errorf("len = %d, want 2", len(resp.Aliases))
	}
}

func TestGetCompanyAlias_Success(t *testing.T) {
	h, companyRepo, aliasRepo := setupCompanyAliasHandler()
	userID := entity.NewUserID()
	company := seedCompanyForAlias(t, companyRepo, userID)
	alias := seedAlias(t, aliasRepo, userID, company.ID(), "トヨタ")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.GetCompanyAlias(w, req, openapi.AliasId(uuid.UUID(alias.ID())))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}
	var resp openapi.CompanyAliasResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Id != uuid.UUID(alias.ID()) {
		t.Errorf("Id = %v, want %v", resp.Id, uuid.UUID(alias.ID()))
	}
}

func TestGetCompanyAlias_OtherUser(t *testing.T) {
	h, companyRepo, aliasRepo := setupCompanyAliasHandler()
	owner := entity.NewUserID()
	other := entity.NewUserID()
	company := seedCompanyForAlias(t, companyRepo, owner)
	alias := seedAlias(t, aliasRepo, owner, company.ID(), "トヨタ")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), other))
	w := httptest.NewRecorder()

	h.GetCompanyAlias(w, req, openapi.AliasId(uuid.UUID(alias.ID())))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404, body = %s", w.Code, w.Body.String())
	}
}

func TestDeleteCompanyAlias_Success(t *testing.T) {
	h, companyRepo, aliasRepo := setupCompanyAliasHandler()
	userID := entity.NewUserID()
	company := seedCompanyForAlias(t, companyRepo, userID)
	alias := seedAlias(t, aliasRepo, userID, company.ID(), "トヨタ")

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.DeleteCompanyAlias(w, req, openapi.AliasId(uuid.UUID(alias.ID())))

	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204, body = %s", w.Code, w.Body.String())
	}
}

func TestDeleteCompanyAlias_NotFound(t *testing.T) {
	h, _, _ := setupCompanyAliasHandler()
	userID := entity.NewUserID()

	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.DeleteCompanyAlias(w, req, openapi.AliasId(uuid.New()))

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want 404, body = %s", w.Code, w.Body.String())
	}
}
