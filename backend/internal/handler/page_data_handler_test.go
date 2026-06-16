package handler

import (
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
	entryuc "github.com/karimiku/job-hunting-saas/internal/usecase/entry"
	inboxclipuc "github.com/karimiku/job-hunting-saas/internal/usecase/inbox_clip"
	taskuc "github.com/karimiku/job-hunting-saas/internal/usecase/task"
)

func setupPageDataHandler() (*PageDataHandler, *inmemory.UserRepository, *inmemory.EntryRepository, *inmemory.CompanyRepository, *inmemory.TaskRepository, *inmemory.InboxClipRepository) {
	userRepo := inmemory.NewUserRepository()
	entryRepo := inmemory.NewEntryRepository()
	companyRepo := inmemory.NewCompanyRepository()
	taskRepo := inmemory.NewTaskRepository(entryRepo)
	inboxClipRepo := inmemory.NewInboxClipRepository()

	h := NewPageDataHandler(
		userRepo,
		entryuc.NewList(entryRepo),
		companyuc.NewList(companyRepo),
		inboxclipuc.NewList(inboxClipRepo),
		taskuc.NewListAll(taskRepo),
	)
	return h, userRepo, entryRepo, companyRepo, taskRepo, inboxClipRepo
}

func TestGetTaskPageData_Success_ReturnsOnlyCurrentUserData(t *testing.T) {
	h, userRepo, entryRepo, companyRepo, taskRepo, _ := setupPageDataHandler()
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
	h, _, _, _, _, _ := setupPageDataHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/page-data/task", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.GetTaskPageData(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401, body = %s", w.Code, w.Body.String())
	}
}

func TestGetAppPageData_Success_ReturnsOnlyCurrentUserData(t *testing.T) {
	h, userRepo, entryRepo, companyRepo, taskRepo, inboxClipRepo := setupPageDataHandler()
	user := seedUser(t, userRepo, "app-page@example.com", "App Page User")
	task, entry := seedTask(t, taskRepo, entryRepo, companyRepo, user.ID())
	clip := seedInboxClip(t, inboxClipRepo, user.ID(), "https://example.com/jobs/1")
	seedTask(t, taskRepo, entryRepo, companyRepo, entity.NewUserID())
	seedInboxClip(t, inboxClipRepo, entity.NewUserID(), "https://example.com/jobs/other")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/page-data/app", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), user.ID()))
	w := httptest.NewRecorder()

	h.GetAppPageData(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}

	var resp openapi.AppPageDataResponse
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
	if len(resp.Tasks) != 1 || resp.Tasks[0].Id.String() != task.ID().String() {
		t.Fatalf("tasks = %+v, want only %s", resp.Tasks, task.ID().String())
	}
	if len(resp.Clips) != 1 || resp.Clips[0].Id.String() != clip.ID().String() {
		t.Fatalf("clips = %+v, want only %s", resp.Clips, clip.ID().String())
	}
	if len(resp.Companies) != 1 {
		t.Fatalf("companies len = %d, want 1", len(resp.Companies))
	}
}

func TestGetAppPageData_UserNotFound_ReturnsUnauthorized(t *testing.T) {
	h, _, _, _, _, _ := setupPageDataHandler()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/page-data/app", nil)
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.GetAppPageData(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401, body = %s", w.Code, w.Body.String())
	}
}

func seedInboxClip(t *testing.T, repo *inmemory.InboxClipRepository, userID entity.UserID, rawURL string) *entity.InboxClip {
	t.Helper()

	clipURL, err := value.NewURL(rawURL)
	if err != nil {
		t.Fatalf("NewURL: %v", err)
	}
	title, err := value.NewInboxClipTitle("求人ページ")
	if err != nil {
		t.Fatalf("NewInboxClipTitle: %v", err)
	}
	source, err := value.NewSource("リクナビ")
	if err != nil {
		t.Fatalf("NewSource: %v", err)
	}
	guess, err := value.NewInboxClipGuess("テスト企業")
	if err != nil {
		t.Fatalf("NewInboxClipGuess: %v", err)
	}
	clip := entity.NewInboxClip(userID, clipURL, title, source, guess)
	if err := repo.Create(t.Context(), clip); err != nil {
		t.Fatalf("create inbox clip: %v", err)
	}
	return clip
}
