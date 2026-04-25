package handler

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	"github.com/karimiku/job-hunting-saas/internal/gen/openapi"
	"github.com/karimiku/job-hunting-saas/internal/middleware"
	taskuc "github.com/karimiku/job-hunting-saas/internal/usecase/task"
)

// TaskHandler は oapi-codegen が生成した ServerInterface のTask関連メソッドを実装する。
type TaskHandler struct {
	createUseCase *taskuc.Create
	getUseCase    *taskuc.Get
	listUseCase   *taskuc.List
	updateUseCase *taskuc.Update
	deleteUseCase *taskuc.Delete
}

// NewTaskHandler は TaskHandler に必要なユースケース群を DI して新しい TaskHandler を返す。
func NewTaskHandler(
	createUseCase *taskuc.Create,
	getUseCase *taskuc.Get,
	listUseCase *taskuc.List,
	updateUseCase *taskuc.Update,
	deleteUseCase *taskuc.Delete,
) *TaskHandler {
	return &TaskHandler{
		createUseCase: createUseCase,
		getUseCase:    getUseCase,
		listUseCase:   listUseCase,
		updateUseCase: updateUseCase,
		deleteUseCase: deleteUseCase,
	}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request, entryId openapi.EntryId) {
	var req openapi.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
		return
	}

	memo := ""
	if req.Memo != nil {
		memo = *req.Memo
	}

	created, err := h.createUseCase.Execute(r.Context(), taskuc.CreateInput{
		UserID:  middleware.GetUserID(r.Context()),
		EntryID: entity.EntryID(entryId),
		Title:   req.Title,
		Type:    string(req.Type),
		DueDate: req.DueDate,
		Memo:    memo,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, toTaskResponse(created.Task))
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request, taskId openapi.TaskId) {
	found, err := h.getUseCase.Execute(r.Context(), taskuc.GetInput{
		UserID: middleware.GetUserID(r.Context()),
		TaskID: entity.TaskID(taskId),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(found.Task))
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request, entryId openapi.EntryId) {
	result, err := h.listUseCase.Execute(r.Context(), taskuc.ListInput{
		UserID:  middleware.GetUserID(r.Context()),
		EntryID: entity.EntryID(entryId),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	items := make([]openapi.TaskResponse, len(result.Tasks))
	for i, task := range result.Tasks {
		items[i] = toTaskResponse(task)
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"tasks": items,
	})
}

// UpdateTask は PATCH リクエストを処理する。
// UseCaseは完全な更新入力(PUT相当)を前提とするため、
// HTTP層で現在値を取得し、未送信フィールドを現在値で埋めてから UseCase に渡す。
func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request, taskId openapi.TaskId) {
	var req openapi.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, openapi.ErrorResponse{Message: "invalid request body"})
		return
	}

	userID := middleware.GetUserID(r.Context())
	taskID := entity.TaskID(taskId)

	existing, err := h.getUseCase.Execute(r.Context(), taskuc.GetInput{
		UserID: userID,
		TaskID: taskID,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	// PATCH: 未送信フィールド(nil)は現在値を維持し、送信されたフィールドのみ上書きする
	resolvedTitle := existing.Task.Title().String()
	if req.Title != nil {
		resolvedTitle = *req.Title
	}

	resolvedType := existing.Task.TaskType().String()
	if req.Type != nil {
		resolvedType = string(*req.Type)
	}

	resolvedStatus := existing.Task.Status().String()
	if req.Status != nil {
		resolvedStatus = string(*req.Status)
	}

	// DueDate: nilは「フィールド未送信」ではなく「クリア」の可能性もあるが、
	// JSONのomitemptyにより未送信時はnilとなる。明示的にnullを送信した場合もnil。
	// UseCaseはnilでClearDueDate、非nilでSetDueDateを行う。
	resolvedDueDate := existing.Task.DueDate()
	if req.DueDate != nil {
		resolvedDueDate = req.DueDate
	}

	resolvedNotify := existing.Task.Notify()
	if req.Notify != nil {
		resolvedNotify = *req.Notify
	}

	resolvedMemo := existing.Task.Memo()
	if req.Memo != nil {
		resolvedMemo = *req.Memo
	}

	updated, err := h.updateUseCase.Execute(r.Context(), taskuc.UpdateInput{
		UserID:  userID,
		TaskID:  taskID,
		Title:   resolvedTitle,
		Type:    resolvedType,
		Status:  resolvedStatus,
		DueDate: resolvedDueDate,
		Notify:  resolvedNotify,
		Memo:    resolvedMemo,
	})
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, toTaskResponse(updated.Task))
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request, taskId openapi.TaskId) {
	err := h.deleteUseCase.Execute(r.Context(), taskuc.DeleteInput{
		UserID: middleware.GetUserID(r.Context()),
		TaskID: entity.TaskID(taskId),
	})
	if err != nil {
		writeError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// toTaskResponse はドメインエンティティをAPI応答用のDTOに変換する。
func toTaskResponse(task *entity.Task) openapi.TaskResponse {
	return openapi.TaskResponse{
		Id:        uuid.UUID(task.ID()),
		EntryId:   uuid.UUID(task.EntryID()),
		Title:     task.Title().String(),
		Type:      task.TaskType().String(),
		Status:    task.Status().String(),
		DueDate:   task.DueDate(),
		Notify:    task.Notify(),
		Memo:      task.Memo(),
		CreatedAt: task.CreatedAt(),
		UpdatedAt: task.UpdatedAt(),
	}
}
