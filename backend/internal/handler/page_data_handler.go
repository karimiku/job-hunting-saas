package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/domain/repository"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	companyuc "github.com/karimiku/job-hunting-saas/internal/usecase/company"
	entryuc "github.com/karimiku/job-hunting-saas/internal/usecase/entry"
	inboxclipuc "github.com/karimiku/job-hunting-saas/internal/usecase/inbox_clip"
	taskuc "github.com/karimiku/job-hunting-saas/internal/usecase/task"
)

// PageDataHandler は初期表示に必要な複数リソースを1リクエストで返す。
type PageDataHandler struct {
	userRepo              repository.UserRepository
	listEntriesUseCase    *entryuc.List
	listCompaniesUseCase  *companyuc.List
	listInboxClipsUseCase *inboxclipuc.List
	listAllTasksUseCase   *taskuc.ListAll
}

// NewPageDataHandler は PageDataHandler に必要な依存を DI して新しい PageDataHandler を返す。
func NewPageDataHandler(
	userRepo repository.UserRepository,
	listEntriesUseCase *entryuc.List,
	listCompaniesUseCase *companyuc.List,
	listInboxClipsUseCase *inboxclipuc.List,
	listAllTasksUseCase *taskuc.ListAll,
) *PageDataHandler {
	return &PageDataHandler{
		userRepo:              userRepo,
		listEntriesUseCase:    listEntriesUseCase,
		listCompaniesUseCase:  listCompaniesUseCase,
		listInboxClipsUseCase: listInboxClipsUseCase,
		listAllTasksUseCase:   listAllTasksUseCase,
	}
}

// GetTaskPageData は /task 初期表示に必要な current user / entries / tasks をまとめて返す。
func (h *PageDataHandler) GetTaskPageData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)
	if userID.IsZero() {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		log.Printf("page data: FindByID failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	entriesResult, err := h.listEntriesUseCase.Execute(ctx, entryuc.ListInput{
		UserID: userID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	companyNameByID := map[entity.CompanyID]string{}
	if len(entriesResult.Entries) > 0 {
		companyList, err := h.listCompaniesUseCase.Execute(ctx, companyuc.ListInput{
			UserID: userID,
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

	entries := make([]openapi.EntryResponse, len(entriesResult.Entries))
	for i, entry := range entriesResult.Entries {
		entries[i] = toEntryResponse(entry)
		if companyName, ok := companyNameByID[entry.CompanyID()]; ok {
			entries[i].CompanyName = &companyName
		}
	}

	tasksResult, err := h.listAllTasksUseCase.Execute(ctx, taskuc.ListAllInput{
		UserID: userID,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	tasks := make([]openapi.TaskResponse, len(tasksResult.Tasks))
	for i, task := range tasksResult.Tasks {
		tasks[i] = toTaskResponse(task)
	}

	writeJSON(w, http.StatusOK, openapi.TaskPageDataResponse{
		User: openapi.CurrentUserResponse{
			Id:    user.ID().String(),
			Email: user.Email().String(),
			Name:  user.Name().String(),
		},
		Entries: entries,
		Tasks:   tasks,
	})
}

// GetAppPageData は dashboard / entry / kanban / inbox の初期表示に必要なデータをまとめて返す。
func (h *PageDataHandler) GetAppPageData(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID := middleware.GetUserID(ctx)
	if userID.IsZero() {
		http.Error(w, "unauthenticated", http.StatusUnauthorized)
		return
	}

	user, err := h.userRepo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, "user not found", http.StatusUnauthorized)
			return
		}
		log.Printf("page data: FindByID failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	companyList, err := h.listCompaniesUseCase.Execute(ctx, companyuc.ListInput{
		UserID: userID,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	companyNameByID := make(map[entity.CompanyID]string, len(companyList.Companies))
	companies := make([]openapi.CompanyResponse, len(companyList.Companies))
	for i, company := range companyList.Companies {
		companyNameByID[company.ID()] = company.Name().String()
		companies[i] = toCompanyResponse(company)
	}

	entriesResult, err := h.listEntriesUseCase.Execute(ctx, entryuc.ListInput{
		UserID: userID,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	entries := make([]openapi.EntryResponse, len(entriesResult.Entries))
	for i, entry := range entriesResult.Entries {
		entries[i] = toEntryResponse(entry)
		if companyName, ok := companyNameByID[entry.CompanyID()]; ok {
			entries[i].CompanyName = &companyName
		}
	}

	tasksResult, err := h.listAllTasksUseCase.Execute(ctx, taskuc.ListAllInput{
		UserID: userID,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	tasks := make([]openapi.TaskResponse, len(tasksResult.Tasks))
	for i, task := range tasksResult.Tasks {
		tasks[i] = toTaskResponse(task)
	}

	clipsResult, err := h.listInboxClipsUseCase.Execute(ctx, inboxclipuc.ListInput{
		UserID: userID,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	clips := make([]openapi.InboxClipResponse, len(clipsResult.Clips))
	for i, clip := range clipsResult.Clips {
		clips[i] = toInboxClipResponse(clip)
	}

	writeJSON(w, http.StatusOK, openapi.AppPageDataResponse{
		User: openapi.CurrentUserResponse{
			Id:    user.ID().String(),
			Email: user.Email().String(),
			Name:  user.Name().String(),
		},
		Entries:   entries,
		Tasks:     tasks,
		Clips:     clips,
		Companies: companies,
	})
}
