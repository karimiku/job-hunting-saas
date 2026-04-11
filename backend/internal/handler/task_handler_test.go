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
	taskuc "github.com/karimiku/job-hunting-saas/internal/usecase/task"
)

// setupTaskHandler はテスト用のTaskHandlerとリポジトリを初期化する。
func setupTaskHandler() (*TaskHandler, *inmemory.TaskRepository, *inmemory.EntryRepository, *inmemory.CompanyRepository) {
	companyRepo := inmemory.NewCompanyRepository()
	entryRepo := inmemory.NewEntryRepository()
	taskRepo := inmemory.NewTaskRepository(entryRepo)

	h := NewTaskHandler(
		taskuc.NewCreate(taskRepo, entryRepo),
		taskuc.NewGet(taskRepo),
		taskuc.NewList(taskRepo),
		taskuc.NewUpdate(taskRepo),
		taskuc.NewDelete(taskRepo),
	)
	return h, taskRepo, entryRepo, companyRepo
}

// seedTask はテスト用のTaskを作成して保存する。
func seedTask(t *testing.T, taskRepo *inmemory.TaskRepository, entryRepo *inmemory.EntryRepository, companyRepo *inmemory.CompanyRepository, userID entity.UserID) (*entity.Task, *entity.Entry) {
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

	title, err := value.NewTaskTitle("ES提出")
	if err != nil {
		t.Fatalf("NewTaskTitle: %v", err)
	}
	taskType, err := value.NewTaskType("deadline")
	if err != nil {
		t.Fatalf("NewTaskType: %v", err)
	}
	task := entity.NewTask(entry.ID(), title, taskType)
	if err := taskRepo.Save(context.Background(), task); err != nil {
		t.Fatalf("save task: %v", err)
	}

	return task, entry
}

func TestUpdateTask_PatchMerge_OnlyTitleSent(t *testing.T) {
	h, taskRepo, entryRepo, companyRepo := setupTaskHandler()
	userID := entity.NewUserID()
	task, _ := seedTask(t, taskRepo, entryRepo, companyRepo, userID)

	newTitle := "面接準備"
	body, _ := json.Marshal(openapi.UpdateTaskRequest{
		Title: &newTitle,
	})

	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.UpdateTask(w, req, openapi.TaskId(task.ID()))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200, body = %s", w.Code, w.Body.String())
	}

	var resp openapi.TaskResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	// 送信したフィールドだけ更新される
	if resp.Title != "面接準備" {
		t.Errorf("Title = %q, want %q", resp.Title, "面接準備")
	}
	// 未送信フィールドは既存値を維持する
	if resp.Type != "deadline" {
		t.Errorf("Type = %q, want %q", resp.Type, "deadline")
	}
	if resp.Status != "todo" {
		t.Errorf("Status = %q, want %q", resp.Status, "todo")
	}
	if resp.Notify != false {
		t.Errorf("Notify = %v, want false", resp.Notify)
	}
}

func TestUpdateTask_PatchMerge_OnlyMemoSent(t *testing.T) {
	h, taskRepo, entryRepo, companyRepo := setupTaskHandler()
	userID := entity.NewUserID()
	task, _ := seedTask(t, taskRepo, entryRepo, companyRepo, userID)

	newMemo := "更新メモ"
	body, _ := json.Marshal(openapi.UpdateTaskRequest{
		Memo: &newMemo,
	})

	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), userID))
	w := httptest.NewRecorder()

	h.UpdateTask(w, req, openapi.TaskId(task.ID()))

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", w.Code)
	}

	var resp openapi.TaskResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}

	if resp.Memo != "更新メモ" {
		t.Errorf("Memo = %q, want %q", resp.Memo, "更新メモ")
	}
	if resp.Title != "ES提出" {
		t.Errorf("Title = %q, want %q", resp.Title, "ES提出")
	}
}

func TestUpdateTask_InvalidJSON(t *testing.T) {
	h, _, _, _ := setupTaskHandler()

	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader([]byte("invalid")))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.UpdateTask(w, req, openapi.TaskId(entity.NewTaskID()))

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", w.Code)
	}
}

func TestUpdateTask_NotFound(t *testing.T) {
	h, _, _, _ := setupTaskHandler()

	body, _ := json.Marshal(openapi.UpdateTaskRequest{})
	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewReader(body))
	req = req.WithContext(middleware.SetUserID(req.Context(), entity.NewUserID()))
	w := httptest.NewRecorder()

	h.UpdateTask(w, req, openapi.TaskId(entity.NewTaskID()))

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want 404", w.Code)
	}
}
