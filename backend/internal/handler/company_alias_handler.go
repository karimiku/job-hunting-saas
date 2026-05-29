package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	companyaliasuc "github.com/karimiku/job-hunting-saas/internal/usecase/company_alias"
)

// CompanyAliasHandler は企業別名関連の HTTP リクエストを受ける handler。
// HTTP ↔ UseCase の入出力変換のみを担う adapter 層。
type CompanyAliasHandler struct {
	createUseCase *companyaliasuc.Create
	getUseCase    *companyaliasuc.Get
	listUseCase   *companyaliasuc.List
	deleteUseCase *companyaliasuc.Delete
}

// NewCompanyAliasHandler は CompanyAliasHandler に必要なユースケース群を DI して返す。
func NewCompanyAliasHandler(
	createUC *companyaliasuc.Create,
	getUC *companyaliasuc.Get,
	listUC *companyaliasuc.List,
	deleteUC *companyaliasuc.Delete,
) *CompanyAliasHandler {
	return &CompanyAliasHandler{
		createUseCase: createUC,
		getUseCase:    getUC,
		listUseCase:   listUC,
		deleteUseCase: deleteUC,
	}
}

// CreateCompanyAlias は POST /api/v1/companies/{companyId}/aliases のハンドラ。
func (h *CompanyAliasHandler) CreateCompanyAlias(w http.ResponseWriter, r *http.Request, companyId openapi.CompanyId) {
	var req openapi.CreateCompanyAliasRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
		return
	}

	out, err := h.createUseCase.Execute(r.Context(), companyaliasuc.CreateInput{
		UserID:    middleware.GetUserID(r.Context()),
		CompanyID: entity.CompanyID(companyId),
		Alias:     req.Alias,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toCompanyAliasResponse(out.CompanyAlias))
}

// ListCompanyAliases は GET /api/v1/companies/{companyId}/aliases のハンドラ。
func (h *CompanyAliasHandler) ListCompanyAliases(w http.ResponseWriter, r *http.Request, companyId openapi.CompanyId) {
	out, err := h.listUseCase.Execute(r.Context(), companyaliasuc.ListInput{
		UserID:    middleware.GetUserID(r.Context()),
		CompanyID: entity.CompanyID(companyId),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	items := make([]openapi.CompanyAliasResponse, len(out.CompanyAliases))
	for i, a := range out.CompanyAliases {
		items[i] = toCompanyAliasResponse(a)
	}
	writeJSON(w, http.StatusOK, map[string]any{"aliases": items})
}

// GetCompanyAlias は GET /api/v1/aliases/{aliasId} のハンドラ。
func (h *CompanyAliasHandler) GetCompanyAlias(w http.ResponseWriter, r *http.Request, aliasId openapi.AliasId) {
	out, err := h.getUseCase.Execute(r.Context(), companyaliasuc.GetInput{
		UserID:         middleware.GetUserID(r.Context()),
		CompanyAliasID: entity.CompanyAliasID(aliasId),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toCompanyAliasResponse(out.CompanyAlias))
}

// DeleteCompanyAlias は DELETE /api/v1/aliases/{aliasId} のハンドラ。
func (h *CompanyAliasHandler) DeleteCompanyAlias(w http.ResponseWriter, r *http.Request, aliasId openapi.AliasId) {
	err := h.deleteUseCase.Execute(r.Context(), companyaliasuc.DeleteInput{
		UserID:         middleware.GetUserID(r.Context()),
		CompanyAliasID: entity.CompanyAliasID(aliasId),
	})
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// toCompanyAliasResponse はドメインエンティティを API 応答用の DTO に変換する。
func toCompanyAliasResponse(a *entity.CompanyAlias) openapi.CompanyAliasResponse {
	return openapi.CompanyAliasResponse{
		Id:        uuid.UUID(a.ID()),
		CompanyId: uuid.UUID(a.CompanyID()),
		Alias:     a.Alias().String(),
		CreatedAt: a.CreatedAt(),
	}
}
