package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	companyuc "github.com/karimiku/job-hunting-saas/internal/usecase/company"
	entryuc "github.com/karimiku/job-hunting-saas/internal/usecase/entry"
)

// EntryHandler は oapi-codegen が生成した ServerInterface のEntry関連メソッドを実装する。
type EntryHandler struct {
	createUseCase            *entryuc.Create
	createWithCompanyUseCase *entryuc.CreateWithCompany
	getUseCase               *entryuc.Get
	listUseCase              *entryuc.List
	listCompaniesUseCase     *companyuc.List
	updateUseCase            *entryuc.Update
	deleteUseCase            *entryuc.Delete
}

// NewEntryHandler は EntryHandler に必要なユースケース群を DI して新しい EntryHandler を返す。
func NewEntryHandler(
	createUseCase *entryuc.Create,
	createWithCompanyUseCase *entryuc.CreateWithCompany,
	getUseCase *entryuc.Get,
	listUseCase *entryuc.List,
	listCompaniesUseCase *companyuc.List,
	updateUseCase *entryuc.Update,
	deleteUseCase *entryuc.Delete,
) *EntryHandler {
	return &EntryHandler{
		createUseCase:            createUseCase,
		createWithCompanyUseCase: createWithCompanyUseCase,
		getUseCase:               getUseCase,
		listUseCase:              listUseCase,
		listCompaniesUseCase:     listCompaniesUseCase,
		updateUseCase:            updateUseCase,
		deleteUseCase:            deleteUseCase,
	}
}

// CreateEntry は POST /entries の handler。リクエストボディからエントリーを新規作成する。
func (h *EntryHandler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	var req openapi.CreateEntryRequest
	if !decodeJSONBody(w, r, &req, maxDefaultJSONBodyBytes) {
		return
	}

	memo := ""
	if req.Memo != nil {
		memo = *req.Memo
	}
	sourceURL := ""
	if req.SourceUrl != nil {
		sourceURL = *req.SourceUrl
	}

	created, err := h.createUseCase.Execute(r.Context(), entryuc.CreateInput{
		UserID:    middleware.GetUserID(r.Context()),
		CompanyID: entity.CompanyID(req.CompanyId),
		Route:     req.Route,
		Source:    req.Source,
		SourceURL: sourceURL,
		Memo:      memo,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toEntryResponse(created.Entry))
}

// CreateEntryWithCompany は POST /entries/with-company の handler。
// Company と Entry を部分作成なしで登録する。
func (h *EntryHandler) CreateEntryWithCompany(w http.ResponseWriter, r *http.Request) {
	var req openapi.CreateEntryWithCompanyRequest
	if !decodeJSONBody(w, r, &req, maxDefaultJSONBodyBytes) {
		return
	}

	memo := ""
	if req.Memo != nil {
		memo = *req.Memo
	}
	sourceURL := ""
	if req.SourceUrl != nil {
		sourceURL = *req.SourceUrl
	}

	created, err := h.createWithCompanyUseCase.Execute(r.Context(), entryuc.CreateWithCompanyInput{
		UserID:      middleware.GetUserID(r.Context()),
		CompanyName: req.CompanyName,
		Route:       req.Route,
		Source:      req.Source,
		SourceURL:   sourceURL,
		Memo:        memo,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toEntryResponse(created.Entry))
}

// GetEntry は GET /entries/{entryId} の handler。
func (h *EntryHandler) GetEntry(w http.ResponseWriter, r *http.Request, entryId openapi.EntryId) {
	found, err := h.getUseCase.Execute(r.Context(), entryuc.GetInput{
		UserID:  middleware.GetUserID(r.Context()),
		EntryID: entity.EntryID(entryId),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toEntryResponse(found.Entry))
}

// ListEntries は GET /entries の handler。
func (h *EntryHandler) ListEntries(w http.ResponseWriter, r *http.Request, params openapi.ListEntriesParams) {
	input := entryuc.ListInput{
		UserID: middleware.GetUserID(r.Context()),
	}

	if params.Status != nil {
		s := string(*params.Status)
		input.Status = &s
	}
	if params.StageKind != nil {
		k := string(*params.StageKind)
		input.StageKind = &k
	}
	if params.Source != nil {
		input.Source = params.Source
	}

	result, err := h.listUseCase.Execute(r.Context(), input)
	if err != nil {
		writeError(w, err)
		return
	}

	companyNameByID := map[entity.CompanyID]string{}
	if h.listCompaniesUseCase != nil && len(result.Entries) > 0 {
		companyList, err := h.listCompaniesUseCase.Execute(r.Context(), companyuc.ListInput{
			UserID: middleware.GetUserID(r.Context()),
		})
		if err != nil {
			writeError(w, err)
			return
		}
		companyNameByID = make(map[entity.CompanyID]string, len(companyList.Companies))
		for _, company := range companyList.Companies {
			companyNameByID[company.ID()] = company.Name().String()
		}
	}

	items := make([]openapi.EntryResponse, len(result.Entries))
	for i, entry := range result.Entries {
		items[i] = toEntryResponse(entry)
		if companyName, ok := companyNameByID[entry.CompanyID()]; ok {
			items[i].CompanyName = &companyName
		}
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"entries": items,
	})
}

// UpdateEntry は PATCH リクエストを処理する。
// UseCaseは完全な更新入力(PUT相当)を前提とするため、
// HTTP層で現在値を取得し、未送信フィールドを現在値で埋めてから UseCase に渡す。
func (h *EntryHandler) UpdateEntry(w http.ResponseWriter, r *http.Request, entryId openapi.EntryId) {
	var req openapi.UpdateEntryRequest
	if !decodeJSONBody(w, r, &req, maxDefaultJSONBodyBytes) {
		return
	}

	userID := middleware.GetUserID(r.Context())
	entryID := entity.EntryID(entryId)

	existing, err := h.getUseCase.Execute(r.Context(), entryuc.GetInput{
		UserID:  userID,
		EntryID: entryID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	// PATCH: 未送信フィールド(nil)は現在値を維持し、送信されたフィールドのみ上書きする
	resolvedSource := existing.Entry.Source().String()
	if req.Source != nil {
		resolvedSource = *req.Source
	}

	resolvedSourceURL := ""
	if existing.Entry.SourceURL() != nil {
		resolvedSourceURL = existing.Entry.SourceURL().String()
	}
	if req.SourceUrl != nil {
		resolvedSourceURL = *req.SourceUrl
	}

	resolvedStatus := existing.Entry.Status().String()
	if req.Status != nil {
		resolvedStatus = string(*req.Status)
	}

	resolvedStageKind := existing.Entry.Stage().Kind().String()
	if req.StageKind != nil {
		resolvedStageKind = string(*req.StageKind)
	}

	resolvedStageLabel := existing.Entry.Stage().Label()
	if req.StageLabel != nil {
		resolvedStageLabel = *req.StageLabel
	}

	resolvedMemo := existing.Entry.Memo()
	if req.Memo != nil {
		resolvedMemo = *req.Memo
	}

	updated, err := h.updateUseCase.Execute(r.Context(), entryuc.UpdateInput{
		UserID:     userID,
		EntryID:    entryID,
		Source:     resolvedSource,
		SourceURL:  resolvedSourceURL,
		Status:     resolvedStatus,
		StageKind:  resolvedStageKind,
		StageLabel: resolvedStageLabel,
		Memo:       resolvedMemo,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toEntryResponse(updated.Entry))
}

// DeleteEntry は DELETE /entries/{entryId} の handler。
func (h *EntryHandler) DeleteEntry(w http.ResponseWriter, r *http.Request, entryId openapi.EntryId) {
	err := h.deleteUseCase.Execute(r.Context(), entryuc.DeleteInput{
		UserID:  middleware.GetUserID(r.Context()),
		EntryID: entity.EntryID(entryId),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// toEntryResponse はドメインエンティティをAPI応答用のDTOに変換する。
func toEntryResponse(entry *entity.Entry) openapi.EntryResponse {
	resp := openapi.EntryResponse{
		Id:         uuid.UUID(entry.ID()),
		CompanyId:  uuid.UUID(entry.CompanyID()),
		Route:      entry.Route().String(),
		Source:     entry.Source().String(),
		Status:     entry.Status().String(),
		StageKind:  entry.Stage().Kind().String(),
		StageLabel: entry.Stage().Label(),
		Memo:       entry.Memo(),
		CreatedAt:  entry.CreatedAt(),
		UpdatedAt:  entry.UpdatedAt(),
	}
	if entry.SourceURL() != nil {
		sourceURL := entry.SourceURL().String()
		resp.SourceUrl = &sourceURL
	}
	return resp
}
