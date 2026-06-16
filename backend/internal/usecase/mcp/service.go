// Package mcp は MCP クライアントに公開する就活コンテキスト操作を束ねる。
package mcp

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/karimiku/job-hunting-saas/internal/domain/entity"
	entryuc "github.com/karimiku/job-hunting-saas/internal/usecase/entry"
	esmemo "github.com/karimiku/job-hunting-saas/internal/usecase/es_memo"
	jobemail "github.com/karimiku/job-hunting-saas/internal/usecase/job_email"
	selectionflowuc "github.com/karimiku/job-hunting-saas/internal/usecase/selection_flow"
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
	ID          string  `json:"id"`
	URL         string  `json:"url"`
	Title       string  `json:"title"`
	Source      string  `json:"source"`
	Guess       string  `json:"guess"`
	ContentText string  `json:"contentText"`
	CapturedAt  *string `json:"capturedAt"`
}

// ESMemoDTO はMCPクライアントに返すES/自己PR/面接メモ情報。
type ESMemoDTO struct {
	ID        string  `json:"id"`
	EntryID   *string `json:"entryId"`
	Company   *string `json:"company"`
	Category  string  `json:"category"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	Source    string  `json:"source"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

// EntryContextDTO はエントリーと関連タスクをまとめたMCP向け情報。
type EntryContextDTO struct {
	Entry         EntryDTO          `json:"entry"`
	Tasks         []TaskDTO         `json:"tasks"`
	SelectionFlow *SelectionFlowDTO `json:"selectionFlow,omitempty"`
}

// SelectionStageDTO はMCPクライアントに返す選考フロー内ステージ。
type SelectionStageDTO struct {
	ID           string `json:"id,omitempty"`
	Position     int    `json:"position"`
	StageKind    string `json:"stageKind"`
	StageLabel   string `json:"stageLabel"`
	EvidenceText string `json:"evidenceText,omitempty"`
}

// SelectionFlowDTO はMCPクライアントに返すEntryごとの可変選考フロー。
type SelectionFlowDTO struct {
	ID                   string              `json:"id,omitempty"`
	EntryID              string              `json:"entryId"`
	Source               string              `json:"source"`
	CurrentStagePosition int                 `json:"currentStagePosition"`
	Confidence           *int                `json:"confidence,omitempty"`
	InboxClipID          *string             `json:"inboxClipId,omitempty"`
	Stages               []SelectionStageDTO `json:"stages"`
	CreatedAt            *string             `json:"createdAt,omitempty"`
	UpdatedAt            *string             `json:"updatedAt,omitempty"`
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

// UpsertEntrySelectionFlowInput はMCP経由で既存Entryの可変選考フローを保存する入力。
type UpsertEntrySelectionFlowInput struct {
	EntryID              string              `json:"entryId"`
	Source               string              `json:"source"`
	CurrentStagePosition int                 `json:"currentStagePosition"`
	Confidence           *int                `json:"confidence"`
	InboxClipID          string              `json:"inboxClipId"`
	Stages               []SelectionStageDTO `json:"stages"`
	Confirm              bool                `json:"confirm"`
}

// CreateEntryFromJobPostingInput はMCP経由で求人本文由来のEntryと可変選考フローを作成する入力。
type CreateEntryFromJobPostingInput struct {
	CompanyName          string              `json:"companyName"`
	Route                string              `json:"route"`
	Source               string              `json:"source"`
	SourceURL            string              `json:"sourceUrl"`
	Memo                 string              `json:"memo"`
	FlowSource           string              `json:"flowSource"`
	CurrentStagePosition int                 `json:"currentStagePosition"`
	Confidence           *int                `json:"confidence"`
	InboxClipID          string              `json:"inboxClipId"`
	Stages               []SelectionStageDTO `json:"stages"`
	Confirm              bool                `json:"confirm"`
}

// Service はMCPクライアントに公開する就活コンテキスト操作を提供する。
type Service struct {
	userID      entity.UserID
	query       ContextQuery
	appendMemo  *esmemo.Append
	listMemo    *esmemo.List
	createTask  *taskuc.Create
	extract     *jobemail.Extract
	createEntry *entryuc.CreateWithCompany
	upsertFlow  *selectionflowuc.Upsert
	getFlow     *selectionflowuc.Get
	now         func() time.Time
}

// NewService はMCPサービスを生成する。
func NewService(
	userID entity.UserID,
	query ContextQuery,
	appendMemo *esmemo.Append,
	listMemo *esmemo.List,
	createTask *taskuc.Create,
	extract *jobemail.Extract,
	createEntry *entryuc.CreateWithCompany,
	upsertFlow *selectionflowuc.Upsert,
	getFlow *selectionflowuc.Get,
) *Service {
	return &Service{
		userID:      userID,
		query:       query,
		appendMemo:  appendMemo,
		listMemo:    listMemo,
		createTask:  createTask,
		extract:     extract,
		createEntry: createEntry,
		upsertFlow:  upsertFlow,
		getFlow:     getFlow,
		now:         time.Now,
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
	entryCtx, err := s.query.GetEntryContext(ctx, s.userID, entryID)
	if err != nil {
		return nil, err
	}
	if s.getFlow != nil {
		if out, err := s.getFlow.Execute(ctx, selectionflowuc.GetInput{
			UserID:  s.userID,
			EntryID: entryID,
		}); err == nil {
			entryCtx.SelectionFlow = selectionFlowDTO(out.SelectionFlow)
		}
	}
	return entryCtx, nil
}

// ListOpenTasks は未完了タスク一覧を返す。
func (s *Service) ListOpenTasks(ctx context.Context) ([]TaskDTO, error) {
	return s.query.ListOpenTasks(ctx, s.userID)
}

// ListInboxClips はInbox Clip一覧を返す。
func (s *Service) ListInboxClips(ctx context.Context) ([]InboxClipDTO, error) {
	return s.query.ListInboxClips(ctx, s.userID)
}

// ListESMemos はES/自己PR/面接メモ一覧を返す。
func (s *Service) ListESMemos(ctx context.Context, limit int32) ([]ESMemoDTO, error) {
	out, err := s.listMemo.Execute(ctx, esmemo.ListInput{
		UserID: s.userID,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}
	memos := make([]ESMemoDTO, 0, len(out.Memos))
	for _, memo := range out.Memos {
		company, err := s.companyForMemo(ctx, memo)
		if err != nil {
			return nil, err
		}
		memos = append(memos, esMemoDTO(memo, company))
	}
	return memos, nil
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
		"memo":    esMemoDTO(out.Memo, nil),
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

// UpsertEntrySelectionFlow は確認付きで既存Entryの可変選考フローを保存する。
func (s *Service) UpsertEntrySelectionFlow(ctx context.Context, input UpsertEntrySelectionFlowInput) (any, error) {
	entryID, err := parseEntryID(input.EntryID)
	if err != nil {
		return nil, err
	}
	inboxClipID, err := optionalInboxClipID(input.InboxClipID)
	if err != nil {
		return nil, err
	}
	flowInput := selectionflowuc.UpsertInput{
		UserID:               s.userID,
		EntryID:              entryID,
		Source:               defaultString(input.Source, "ai_paste"),
		CurrentStagePosition: input.CurrentStagePosition,
		Confidence:           input.Confidence,
		InboxClipID:          inboxClipID,
		Stages:               stageInputs(input.Stages),
	}
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "upsert_entry_selection_flow",
		"entryId":              entryID.String(),
		"selectionFlow":        selectionFlowPreview(flowInput),
	}
	if !input.Confirm {
		return preview, nil
	}
	if s.upsertFlow == nil {
		return nil, fmt.Errorf("selection flow write is not configured")
	}
	out, err := s.upsertFlow.Execute(ctx, flowInput)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"updated":       true,
		"entry":         entryDTOFromEntry(out.Entry),
		"selectionFlow": selectionFlowDTO(out.SelectionFlow),
	}, nil
}

// CreateEntryFromJobPosting は確認付きで求人本文由来のEntryと可変選考フローを作成する。
func (s *Service) CreateEntryFromJobPosting(ctx context.Context, input CreateEntryFromJobPostingInput) (any, error) {
	companyName := strings.TrimSpace(input.CompanyName)
	if companyName == "" {
		return nil, fmt.Errorf("companyName is required")
	}
	inboxClipID, err := optionalInboxClipID(input.InboxClipID)
	if err != nil {
		return nil, err
	}
	route := defaultString(input.Route, "本選考")
	source := defaultString(input.Source, "求人ページ")
	flowSource := defaultString(input.FlowSource, "ai_paste")
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "create_entry_from_job_posting",
		"entry": map[string]any{
			"companyName": companyName,
			"route":       route,
			"source":      source,
			"sourceUrl":   strings.TrimSpace(input.SourceURL),
			"memo":        strings.TrimSpace(input.Memo),
		},
		"selectionFlow": map[string]any{
			"source":               flowSource,
			"currentStagePosition": input.CurrentStagePosition,
			"confidence":           input.Confidence,
			"inboxClipId":          optionalInboxClipIDString(inboxClipID),
			"stages":               input.Stages,
		},
	}
	if !input.Confirm {
		return preview, nil
	}
	if s.createEntry == nil || s.upsertFlow == nil {
		return nil, fmt.Errorf("entry and selection flow write are not configured")
	}
	created, err := s.createEntry.Execute(ctx, entryuc.CreateWithCompanyInput{
		UserID:      s.userID,
		CompanyName: companyName,
		Route:       route,
		Source:      source,
		SourceURL:   strings.TrimSpace(input.SourceURL),
		Memo:        strings.TrimSpace(input.Memo),
	})
	if err != nil {
		return nil, err
	}
	flowOut, err := s.upsertFlow.Execute(ctx, selectionflowuc.UpsertInput{
		UserID:               s.userID,
		EntryID:              created.Entry.ID(),
		Source:               flowSource,
		CurrentStagePosition: input.CurrentStagePosition,
		Confidence:           input.Confidence,
		InboxClipID:          inboxClipID,
		Stages:               stageInputs(input.Stages),
	})
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"created":       true,
		"entry":         entryDTOFromEntry(created.Entry),
		"selectionFlow": selectionFlowDTO(flowOut.SelectionFlow),
	}, nil
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

func optionalInboxClipID(raw string) (*entity.InboxClipID, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, nil
	}
	id, err := uuid.Parse(strings.TrimSpace(raw))
	if err != nil {
		return nil, fmt.Errorf("invalid inboxClipId: %w", err)
	}
	out := entity.InboxClipID(id)
	return &out, nil
}

func optionalInboxClipIDString(id *entity.InboxClipID) *string {
	if id == nil {
		return nil
	}
	value := id.String()
	return &value
}

func optionalEntryIDString(id *entity.EntryID) *string {
	if id == nil {
		return nil
	}
	value := id.String()
	return &value
}

func (s *Service) companyForMemo(ctx context.Context, memo *entity.ESMemo) (*string, error) {
	if memo.EntryID() == nil {
		return nil, nil
	}
	entryCtx, err := s.query.GetEntryContext(ctx, s.userID, *memo.EntryID())
	if err != nil {
		return nil, err
	}
	if strings.TrimSpace(entryCtx.Entry.Company) == "" {
		return nil, nil
	}
	company := entryCtx.Entry.Company
	return &company, nil
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

func esMemoDTO(memo *entity.ESMemo, company *string) ESMemoDTO {
	return ESMemoDTO{
		ID:        memo.ID().String(),
		EntryID:   optionalEntryIDString(memo.EntryID()),
		Company:   company,
		Category:  memo.Category().String(),
		Title:     memo.Title().String(),
		Content:   memo.Content().String(),
		Source:    memo.Source().String(),
		CreatedAt: memo.CreatedAt().Format(time.RFC3339),
		UpdatedAt: memo.UpdatedAt().Format(time.RFC3339),
	}
}

func stageInputs(stages []SelectionStageDTO) []selectionflowuc.StageInput {
	out := make([]selectionflowuc.StageInput, 0, len(stages))
	for _, stage := range stages {
		out = append(out, selectionflowuc.StageInput{
			StageKind:    stage.StageKind,
			StageLabel:   stage.StageLabel,
			EvidenceText: stage.EvidenceText,
		})
	}
	return out
}

func selectionFlowPreview(input selectionflowuc.UpsertInput) map[string]any {
	return map[string]any{
		"source":               input.Source,
		"currentStagePosition": input.CurrentStagePosition,
		"confidence":           input.Confidence,
		"inboxClipId":          optionalInboxClipIDString(input.InboxClipID),
		"stages":               input.Stages,
	}
}

func selectionFlowDTO(flow *entity.SelectionFlow) *SelectionFlowDTO {
	if flow == nil {
		return nil
	}
	stages := make([]SelectionStageDTO, 0, len(flow.Stages()))
	for _, stage := range flow.Stages() {
		stages = append(stages, SelectionStageDTO{
			ID:           stage.ID().String(),
			Position:     stage.Position(),
			StageKind:    stage.Stage().Kind().String(),
			StageLabel:   stage.Stage().Label(),
			EvidenceText: stage.EvidenceText(),
		})
	}
	return &SelectionFlowDTO{
		ID:                   flow.ID().String(),
		EntryID:              flow.EntryID().String(),
		Source:               flow.Source().String(),
		CurrentStagePosition: flow.CurrentStagePosition(),
		Confidence:           flow.Confidence(),
		InboxClipID:          optionalInboxClipIDString(flow.InboxClipID()),
		Stages:               stages,
		CreatedAt:            stringPtr(flow.CreatedAt().Format(time.RFC3339)),
		UpdatedAt:            stringPtr(flow.UpdatedAt().Format(time.RFC3339)),
	}
}

func entryDTOFromEntry(entry *entity.Entry) EntryDTO {
	sourceURL := ""
	if entry.SourceURL() != nil {
		sourceURL = entry.SourceURL().String()
	}
	createdAt := entry.CreatedAt().Format(time.RFC3339)
	updatedAt := entry.UpdatedAt().Format(time.RFC3339)
	return EntryDTO{
		ID:         entry.ID().String(),
		CompanyID:  entry.CompanyID().String(),
		Route:      entry.Route().String(),
		Source:     entry.Source().String(),
		SourceURL:  sourceURL,
		Status:     entry.Status().String(),
		StageKind:  entry.Stage().Kind().String(),
		StageLabel: entry.Stage().Label(),
		Memo:       entry.Memo(),
		CreatedAt:  &createdAt,
		UpdatedAt:  &updatedAt,
	}
}

func stringPtr(value string) *string {
	return &value
}
