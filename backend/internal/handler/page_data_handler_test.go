package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/infra/inmemory"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	companyuc "github.com/karimiku/job-hunting-saas/internal/usecase/company"
	entryuc "github.com/karimiku/job-hunting-saas/internal/usecase/entry"
	taskuc "github.com/karimiku/job-hunting-saas/internal/usecase/task"
)

func setupPageDataHandler() (*PageDataHandler, *inmemory.UserRepository, *inmemory.EntryRepository, *inmemory.CompanyRepository, *inmemory.TaskRepository) {
	userRepo := inmemory.NewUserRepository()
	entryRepo := inmemory.NewEntryRepository()
	companyRepo := inmemory.NewCompanyRepository()
	taskRepo := inmemory.NewTaskRepository(entryRepo)

	h := NewPageDataHandler(
		userRepo,
		entryuc.NewList(entryRepo),
		companyuc.NewList(companyRepo),
		taskuc.NewListAll(taskRepo),
	)
	return h, userRepo, entryRepo, companyRepo, taskRepo
}

func TestGetTaskPageData_Success_ReturnsOnlyCurrentUserData(t *testing.T) {
	h, userRepo, entryRepo, companyRepo, taskRepo := setupPageDataHandler()
	user := seedUser(t, userRepo, "task-page@example.com", "Task Page User")
	task, entry := seedTask(t, taskRepo, entryRepo, companyRepo, user.ID())
	seedTask(t, taskRepo, entryRepo, companyRepo, entity.NewUserID())

	req := httptest.NewRequest(http.MethodGet, "/api/v1/page-data/task", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), user.ID()))
	w := httptest.NewRecorder()

	h.GetTaskPageData(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}

	var resp openapi.TaskPageDataResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.User.Id != user.ID().String() {
		t.Errorf("User.Id = %q, want %q", resp.User.Id, user.ID().String())
	}
	if len(resp.Entries) != 1 {
		t.Fatalf("entries len = %d, want 1", len(resp.Entries))
	}
	if resp.Entries[0].Id.String() != entry.ID().String() {
		t.Errorf("entry id = %s, want %s", resp.Entries[0].Id.String(), entry.ID().String())
	}
	if resp.Entries[0].CompanyName == nil || *resp.Entries[0].CompanyName != "テスト企業" {
		t.Errorf("companyName = %v, want テスト企業", resp.Entries[0].CompanyName)
	}
	if len(resp.Tasks) != 1 {
		t.Fatalf("tasks len = %d, want 1", len(resp.Tasks))
	}
	if resp.Tasks[0].Id.String() != task.ID().String() {
		t.Errorf("task id = %s, want %s", resp.Tasks[0].Id.String(), task.ID().String())
	}
}

func TestGetTaskPageData_UserNotFound_ReturnsUnauthorized(t *testing.T) {
	h, _, _, _, _ := setupPageDataHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/page-data/task", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.GetTaskPageData(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401, body = %s", w.Code, w.Body.String())
	}
}
