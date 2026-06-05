// Package mcp は MCP クライアントに公開する就活コンテキスト操作を束ねる。
package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	esmemo "github.com/karimiku/job-hunting-saas/internal/usecase/es_memo"
	jobemail "github.com/karimiku/job-hunting-saas/internal/usecase/job_email"
	taskuc "github.com/karimiku/job-hunting-saas/internal/usecase/task"
)

// ContextQuery はMCPサービスが必要とする参照系クエリを表す。
type ContextQuery interface {
	ListEntries(ctx context.Context, userID entity.UserID) ([]EntryDTO, error)
	GetEntryContext(ctx context.Context, userID entity.UserID, entryID entity.EntryID) (*EntryContextDTO, error)
	ListOpenTasks(ctx context.Context, userID entity.UserID) ([]TaskDTO, error)
	ListInboxClips(ctx context.Context, userID entity.UserID) ([]InboxClipDTO, error)
}

// EntryDTO はMCPクライアントに返すエントリー情報。
type EntryDTO struct {
	ID         string  `json:"id"`
	CompanyID  string  `json:"companyId"`
	Company    string  `json:"company"`
	Route      string  `json:"route"`
	Source     string  `json:"source"`
	SourceURL  string  `json:"sourceUrl"`
	Status     string  `json:"status"`
	StageKind  string  `json:"stageKind"`
	StageLabel string  `json:"stageLabel"`
	Memo       string  `json:"memo"`
	CreatedAt  *string `json:"createdAt"`
	UpdatedAt  *string `json:"updatedAt"`
}

// TaskDTO はMCPクライアントに返すタスク情報。
type TaskDTO struct {
	ID        string  `json:"id"`
	EntryID   string  `json:"entryId"`
	Company   string  `json:"company"`
	Title     string  `json:"title"`
	Type      string  `json:"type"`
	DueDate   *string `json:"dueDate"`
	Status    string  `json:"status"`
	Notify    bool    `json:"notify"`
	Memo      string  `json:"memo"`
	CreatedAt *string `json:"createdAt"`
	UpdatedAt *string `json:"updatedAt"`
}

// InboxClipDTO はMCPクライアントに返すInbox Clip情報。
type InboxClipDTO struct {
	ID         string  `json:"id"`
	URL        string  `json:"url"`
	Title      string  `json:"title"`
	Source     string  `json:"source"`
	Guess      string  `json:"guess"`
	CapturedAt *string `json:"capturedAt"`
}

// EntryContextDTO はエントリーと関連タスクをまとめたMCP向け情報。
type EntryContextDTO struct {
	Entry EntryDTO  `json:"entry"`
	Tasks []TaskDTO `json:"tasks"`
}

// AppendESMemoInput はMCP経由でESメモを追記する入力。
type AppendESMemoInput struct {
	Title    string `json:"title"`
	Content  string `json:"content"`
	Category string `json:"category"`
	EntryID  string `json:"entryId"`
	Source   string `json:"source"`
	Confirm  bool   `json:"confirm"`
}

// CreateTaskInput はMCP経由でタスクを作成する入力。
type CreateTaskInput struct {
	EntryID string `json:"entryId"`
	Title   string `json:"title"`
	Type    string `json:"type"`
	DueDate string `json:"dueDate"`
	Memo    string `json:"memo"`
	Notify  bool   `json:"notify"`
	Confirm bool   `json:"confirm"`
}

// CaptureJobEmailInput はMCP経由で選考メールを解析する入力。
type CaptureJobEmailInput struct {
	Subject     string `json:"subject"`
	Text        string `json:"text"`
	CompanyName string `json:"companyName"`
}

// Service はMCPクライアントに公開する就活コンテキスト操作を提供する。
type Service struct {
	userID     entity.UserID
	query      ContextQuery
	appendMemo *esmemo.Append
	createTask *taskuc.Create
	extract    *jobemail.Extract
	now        func() time.Time
}

// NewService はMCPサービスを生成する。
func NewService(
	userID entity.UserID,
	query ContextQuery,
	appendMemo *esmemo.Append,
	createTask *taskuc.Create,
	extract *jobemail.Extract,
) *Service {
	return &Service{
		userID:     userID,
		query:      query,
		appendMemo: appendMemo,
		createTask: createTask,
		extract:    extract,
		now:        time.Now,
	}
}

// ListEntries はユーザーのエントリー一覧を返す。
func (s *Service) ListEntries(ctx context.Context) ([]EntryDTO, error) {
	return s.query.ListEntries(ctx, s.userID)
}

// GetEntryContext は指定エントリーと関連タスクを返す。
func (s *Service) GetEntryContext(ctx context.Context, rawEntryID string) (*EntryContextDTO, error) {
	entryID, err := parseEntryID(rawEntryID)
	if err != nil {
		return nil, err
	}
	return s.query.GetEntryContext(ctx, s.userID, entryID)
}

// ListOpenTasks は未完了タスク一覧を返す。
func (s *Service) ListOpenTasks(ctx context.Context) ([]TaskDTO, error) {
	return s.query.ListOpenTasks(ctx, s.userID)
}

// ListInboxClips はInbox Clip一覧を返す。
func (s *Service) ListInboxClips(ctx context.Context) ([]InboxClipDTO, error) {
	return s.query.ListInboxClips(ctx, s.userID)
}

// AppendESMemo は確認付きでESメモを追記する。
func (s *Service) AppendESMemo(ctx context.Context, input AppendESMemoInput) (any, error) {
	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)
	if title == "" || content == "" {
		return nil, fmt.Errorf("title and content are required")
	}
	entryID, err := optionalEntryID(input.EntryID)
	if err != nil {
		return nil, err
	}
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "append_es_memo",
		"memo": map[string]any{
			"title":    title,
			"content":  content,
			"category": defaultString(input.Category, "general"),
			"entryId":  optionalEntryIDString(entryID),
			"source":   defaultString(input.Source, "mcp"),
		},
	}
	if !input.Confirm {
		return preview, nil
	}
	out, err := s.appendMemo.Execute(ctx, esmemo.AppendInput{
		UserID:   s.userID,
		EntryID:  entryID,
		Category: input.Category,
		Title:    title,
		Content:  content,
		Source:   input.Source,
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"created": true,
		"memo":    esMemoDTO(out.Memo),
	}, nil
}

// CreateTask は確認付きでタスクを作成する。
func (s *Service) CreateTask(ctx context.Context, input CreateTaskInput) (any, error) {
	entryID, err := parseEntryID(input.EntryID)
	if err != nil {
		return nil, err
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, fmt.Errorf("title is required")
	}
	dueDate, err := parseOptionalDueDate(input.DueDate)
	if err != nil {
		return nil, err
	}
	entryCtx, err := s.query.GetEntryContext(ctx, s.userID, entryID)
	if err != nil {
		return nil, err
	}
	taskType := defaultString(input.Type, "deadline")
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "create_task",
		"task": map[string]any{
			"entryId": entryID.String(),
			"company": entryCtx.Entry.Company,
			"title":   title,
			"type":    taskType,
			"dueDate": formatTimePtr(dueDate),
			"memo":    input.Memo,
			"notify":  input.Notify,
		},
	}
	if !input.Confirm {
		return preview, nil
	}
	out, err := s.createTask.Execute(ctx, taskuc.CreateInput{
		UserID:  s.userID,
		EntryID: entryID,
		Title:   title,
		Type:    taskType,
		DueDate: dueDate,
		Memo:    input.Memo,
		Notify:  input.Notify,
	})
	if err != nil {
		return nil, err
	}
	task := out.Task
	return map[string]any{
		"created": true,
		"task": map[string]any{
			"id":      task.ID().String(),
			"entryId": task.EntryID().String(),
			"company": entryCtx.Entry.Company,
			"title":   task.Title().String(),
			"type":    task.TaskType().String(),
			"dueDate": formatTimePtr(task.DueDate()),
			"status":  task.Status().String(),
			"notify":  task.Notify(),
			"memo":    task.Memo(),
		},
	}, nil
}

// CaptureJobEmail は選考メール本文から候補情報を抽出する。
func (s *Service) CaptureJobEmail(input CaptureJobEmailInput) (jobemail.ExtractOutput, error) {
	if strings.TrimSpace(input.Text) == "" {
		return jobemail.ExtractOutput{}, fmt.Errorf("text is required")
	}
	return s.extract.Execute(jobemail.ExtractInput{
		Subject:     input.Subject,
		Text:        input.Text,
		CompanyName: input.CompanyName,
		Now:         s.now(),
	}), nil
}

func parseEntryID(raw string) (entity.EntryID, error) {
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return entity.EntryID{}, fmt.Errorf("invalid entryId: %w", err)
	}
	return entity.EntryID(id), nil
}

func optionalEntryID(raw string) (*entity.EntryID, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	id, err := parseEntryID(raw)
	if err != nil {
		return nil, err
	}
	return &id, nil
}

func optionalEntryIDString(id *entity.EntryID) *string {
	if id == nil {
		return nil
	}
	value := id.String()
	return &value
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func parseOptionalDueDate(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t, nil
	}
	if t, err := time.ParseInLocation("2006-01-02", raw, time.Local); err == nil {
		return &t, nil
	}
	return nil, fmt.Errorf("invalid dueDate %q: use YYYY-MM-DD or RFC3339", raw)
}

func formatTimePtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	value := t.Format(time.RFC3339)
	return &value
}

func esMemoDTO(memo *entity.ESMemo) map[string]any {
	return map[string]any{
		"id":        memo.ID().String(),
		"entryId":   optionalEntryIDString(memo.EntryID()),
		"category":  memo.Category().String(),
		"title":     memo.Title().String(),
		"content":   memo.Content().String(),
		"source":    memo.Source().String(),
		"createdAt": memo.CreatedAt().Format(time.RFC3339),
		"updatedAt": memo.UpdatedAt().Format(time.RFC3339),
	}
}
