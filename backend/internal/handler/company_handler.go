package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	companyuc "github.com/karimiku/job-hunting-saas/internal/usecase/company"
)

// CompanyHandler は oapi-codegen が生成した ServerInterface を実装し、
// HTTPリクエストをUseCaseへ変換する adapter 層。
type CompanyHandler struct {
	createUseCase *companyuc.Create
	getUseCase    *companyuc.Get
	listUseCase   *companyuc.List
	updateUseCase *companyuc.Update
	deleteUseCase *companyuc.Delete
}

// NewCompanyHandler は CompanyHandler に必要なユースケース群を DI して新しい CompanyHandler を返す。
func NewCompanyHandler(
	createUseCase *companyuc.Create,
	getUseCase *companyuc.Get,
	listUseCase *companyuc.List,
	updateUseCase *companyuc.Update,
	deleteUseCase *companyuc.Delete,
) *CompanyHandler {
	return &CompanyHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

// CreateCompany は POST /companies の handler。リクエストボディから企業を新規作成する。
func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {
	var createReq openapi.CreateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&createReq); err != nil {
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
		return
	}

	// OpenAPIスキーマで memo は optional。未送信時は空文字をデフォルトとする。
	memo := ""
	if createReq.Memo != nil {
		memo = *createReq.Memo
	}

	createdCompany, err := h.createUseCase.Execute(r.Context(), companyuc.CreateInput{
		UserID: middleware.GetUserID(r.Context()),
		Name:   createReq.Name,
		Memo:   memo,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toCompanyResponse(createdCompany.Company))
}

// GetCompany は GET /companies/{companyId} の handler。
func (h *CompanyHandler) GetCompany(w http.ResponseWriter, r *http.Request, companyId openapi.CompanyId) {
	foundCompany, err := h.getUseCase.Execute(r.Context(), companyuc.GetInput{
		UserID:    middleware.GetUserID(r.Context()),
		CompanyID: entity.CompanyID(companyId),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toCompanyResponse(foundCompany.Company))
}

// ListCompanies は GET /companies の handler。
func (h *CompanyHandler) ListCompanies(w http.ResponseWriter, r *http.Request) {
	companyList, err := h.listUseCase.Execute(r.Context(), companyuc.ListInput{
		UserID: middleware.GetUserID(r.Context()),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	responseItems := make([]openapi.CompanyResponse, len(companyList.Companies))
	for i, company := range companyList.Companies {
		responseItems[i] = toCompanyResponse(company)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"companies": responseItems,
	})
}

// UpdateCompany は PATCH リクエストを処理する。
// UseCaseは完全な更新入力(PUT相当)を前提とするため、
// HTTP層で現在値を取得し、未送信フィールドを現在値で埋めてから UseCase に渡す。
func (h *CompanyHandler) UpdateCompany(w http.ResponseWriter, r *http.Request, companyId openapi.CompanyId) {
	var updateReq openapi.UpdateCompanyRequest
	if err := json.NewDecoder(r.Body).Decode(&updateReq); err != nil {
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
		return
	}

	userID := middleware.GetUserID(r.Context())
	companyID := entity.CompanyID(companyId)

	existingCompany, err := h.getUseCase.Execute(r.Context(), companyuc.GetInput{
		UserID:    userID,
		CompanyID: companyID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	// PATCH: 未送信フィールド(nil)は現在値を維持し、送信されたフィールドのみ上書きする
	resolvedName := existingCompany.Company.Name().String()
	if updateReq.Name != nil {
		resolvedName = *updateReq.Name
	}

	resolvedMemo := existingCompany.Company.Memo()
	if updateReq.Memo != nil {
		resolvedMemo = *updateReq.Memo
	}

	updatedCompany, err := h.updateUseCase.Execute(r.Context(), companyuc.UpdateInput{
		UserID:    userID,
		CompanyID: companyID,
		Name:      resolvedName,
		Memo:      resolvedMemo,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toCompanyResponse(updatedCompany.Company))
}

// DeleteCompany は DELETE /companies/{companyId} の handler。
func (h *CompanyHandler) DeleteCompany(w http.ResponseWriter, r *http.Request, companyId openapi.CompanyId) {
	err := h.deleteUseCase.Execute(r.Context(), companyuc.DeleteInput{
		UserID:    middleware.GetUserID(r.Context()),
		CompanyID: entity.CompanyID(companyId),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// toCompanyResponse はドメインエンティティをAPI応答用のDTOに変換する。
// ドメインの内部表現(値オブジェクト等)をクライアント向けのプリミティブ型に変換する責務を持つ。
func toCompanyResponse(company *entity.Company) openapi.CompanyResponse {
	return openapi.CompanyResponse{
		Id:        uuid.UUID(company.ID()),
		Name:      company.Name().String(),
		Memo:      company.Memo(),
		CreatedAt: company.CreatedAt(),
		UpdatedAt: company.UpdatedAt(),
	}
}
