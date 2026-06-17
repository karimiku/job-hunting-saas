// Package entreapi provides HTTP adapters for the hosted Entré API.
package entreapi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	jobemail "github.com/karimiku/job-hunting-saas/internal/usecase/job_email"
	mcpuc "github.com/karimiku/job-hunting-saas/internal/usecase/mcp"
)

const defaultBaseURL = "http://localhost:8080"

// MCPApplication implements the MCP application boundary by calling the REST API.
type MCPApplication struct {
	baseURL    string
	token      string
	httpClient *http.Client
	extract    *jobemail.Extract
	now        func() time.Time
}

// NewMCPApplication creates an API-backed MCP application.
func NewMCPApplication(baseURL string, token string, httpClient *http.Client) (*MCPApplication, error) {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	parsed, err := url.Parse(baseURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("invalid ENTRE_API_BASE_URL %q", baseURL)
	}
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, errors.New("ENTRE_API_TOKEN is required")
	}
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &MCPApplication{
		baseURL:    strings.TrimRight(baseURL, "/"),
		token:      token,
		httpClient: httpClient,
		extract:    jobemail.NewExtract(),
		now:        time.Now,
	}, nil
}

// ListEntries returns entries owned by the token user.
func (a *MCPApplication) ListEntries(ctx context.Context) ([]mcpuc.EntryDTO, error) {
	entries, err := a.listEntries(ctx)
	if err != nil {
		return nil, err
	}
	companies, err := a.companyNames(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]mcpuc.EntryDTO, 0, len(entries))
	for _, entry := range entries {
		out = append(out, entry.toMCP(companies[entry.CompanyID]))
	}
	return out, nil
}

// GetEntryContext returns one entry and its tasks.
func (a *MCPApplication) GetEntryContext(ctx context.Context, rawEntryID string) (*mcpuc.EntryContextDTO, error) {
	entryID := strings.TrimSpace(rawEntryID)
	if entryID == "" {
		return nil, errors.New("entryId is required")
	}
	var entry entryResponse
	if err := a.get(ctx, "/api/v1/entries/"+url.PathEscape(entryID), &entry); err != nil {
		return nil, err
	}
	companyName, err := a.companyName(ctx, entry.CompanyID)
	if err != nil {
		return nil, err
	}
	var tasks listTasksResponse
	if err := a.get(ctx, "/api/v1/entries/"+url.PathEscape(entryID)+"/tasks", &tasks); err != nil {
		return nil, err
	}
	taskDTOs := make([]mcpuc.TaskDTO, 0, len(tasks.Tasks))
	for _, task := range tasks.Tasks {
		taskDTOs = append(taskDTOs, task.toMCP(companyName))
	}
	ctxDTO := &mcpuc.EntryContextDTO{
		Entry: entry.toMCP(companyName),
		Tasks: taskDTOs,
	}
	var flow selectionFlowResponse
	if err := a.get(ctx, "/api/v1/entries/"+url.PathEscape(entryID)+"/selection-flow", &flow); err == nil {
		dto := flow.toMCP()
		ctxDTO.SelectionFlow = &dto
	}
	return ctxDTO, nil
}

// ListOpenTasks returns non-done tasks.
func (a *MCPApplication) ListOpenTasks(ctx context.Context) ([]mcpuc.TaskDTO, error) {
	var tasks listTasksResponse
	if err := a.get(ctx, "/api/v1/tasks", &tasks); err != nil {
		return nil, err
	}
	entryCompanies, err := a.entryCompanyNames(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]mcpuc.TaskDTO, 0, len(tasks.Tasks))
	for _, task := range tasks.Tasks {
		if task.Status == "done" {
			continue
		}
		out = append(out, task.toMCP(entryCompanies[task.EntryID]))
	}
	return out, nil
}

// ListInboxClips returns saved page clips.
func (a *MCPApplication) ListInboxClips(ctx context.Context) ([]mcpuc.InboxClipDTO, error) {
	var clips listInboxClipsResponse
	if err := a.get(ctx, "/api/v1/inbox/clips", &clips); err != nil {
		return nil, err
	}
	out := make([]mcpuc.InboxClipDTO, 0, len(clips.Clips))
	for _, clip := range clips.Clips {
		out = append(out, clip.toMCP())
	}
	return out, nil
}

// ListESMemos returns ES/self PR/interview memos.
func (a *MCPApplication) ListESMemos(ctx context.Context, limit int32) ([]mcpuc.ESMemoDTO, error) {
	path := "/api/v1/es-memos"
	if limit > 0 {
		path += "?limit=" + url.QueryEscape(fmt.Sprintf("%d", limit))
	}
	var memos listESMemosResponse
	if err := a.get(ctx, path, &memos); err != nil {
		return nil, err
	}
	entryCompanies := map[string]string{}
	for _, memo := range memos.Memos {
		if memo.EntryID != nil {
			var err error
			entryCompanies, err = a.entryCompanyNames(ctx)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	out := make([]mcpuc.ESMemoDTO, 0, len(memos.Memos))
	for _, memo := range memos.Memos {
		out = append(out, memo.toMCP(entryCompanies))
	}
	return out, nil
}

// AppendESMemo previews or saves an ES memo through the REST API.
func (a *MCPApplication) AppendESMemo(ctx context.Context, input mcpuc.AppendESMemoInput) (any, error) {
	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)
	if title == "" || content == "" {
		return nil, errors.New("title and content are required")
	}
	entryID := strings.TrimSpace(input.EntryID)
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "append_es_memo",
		"memo": map[string]any{
			"title":    title,
			"content":  content,
			"category": defaultString(input.Category, "general"),
			"entryId":  optionalString(entryID),
			"source":   defaultString(input.Source, "mcp"),
		},
	}
	if !input.Confirm {
		return preview, nil
	}
	req := createESMemoRequest{
		EntryID:  optionalString(entryID),
		Category: optionalString(input.Category),
		Title:    title,
		Content:  content,
		Source:   optionalString(input.Source),
	}
	var memo esMemoResponse
	if err := a.post(ctx, "/api/v1/es-memos", req, &memo); err != nil {
		return nil, err
	}
	return map[string]any{
		"created": true,
		"memo":    memo.toMap(),
	}, nil
}

// CreateTask previews or creates a task through the REST API.
func (a *MCPApplication) CreateTask(ctx context.Context, input mcpuc.CreateTaskInput) (any, error) {
	entryID := strings.TrimSpace(input.EntryID)
	if entryID == "" {
		return nil, errors.New("entryId is required")
	}
	title := strings.TrimSpace(input.Title)
	if title == "" {
		return nil, errors.New("title is required")
	}
	dueDate, err := parseOptionalDueDate(input.DueDate)
	if err != nil {
		return nil, err
	}
	entryCtx, err := a.GetEntryContext(ctx, entryID)
	if err != nil {
		return nil, err
	}
	taskType := defaultString(input.Type, "deadline")
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "create_task",
		"task": map[string]any{
			"entryId": entryID,
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
	req := createTaskRequest{
		Title:  title,
		Type:   taskType,
		Memo:   optionalString(input.Memo),
		Notify: &input.Notify,
	}
	if dueDate != nil {
		value := dueDate.Format(time.RFC3339)
		req.DueDate = &value
	}
	var task taskResponse
	if err := a.post(ctx, "/api/v1/entries/"+url.PathEscape(entryID)+"/tasks", req, &task); err != nil {
		return nil, err
	}
	return map[string]any{
		"created": true,
		"task":    task.toMCP(entryCtx.Entry.Company),
	}, nil
}

// DeleteEntry previews or deletes an entry through the REST API.
func (a *MCPApplication) DeleteEntry(ctx context.Context, input mcpuc.DeleteEntryInput) (any, error) {
	entryID := strings.TrimSpace(input.EntryID)
	if entryID == "" {
		return nil, errors.New("entryId is required")
	}
	entryCtx, err := a.GetEntryContext(ctx, entryID)
	if err != nil {
		return nil, err
	}
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "delete_entry",
		"entry":                entryCtx.Entry,
		"relatedTaskCount":     len(entryCtx.Tasks),
	}
	if !input.Confirm {
		return preview, nil
	}
	if err := a.delete(ctx, "/api/v1/entries/"+url.PathEscape(entryID)); err != nil {
		return nil, err
	}
	return map[string]any{
		"deleted":          true,
		"entry":            entryCtx.Entry,
		"relatedTaskCount": len(entryCtx.Tasks),
	}, nil
}

// CaptureJobEmail extracts structured candidates locally. No LLM API is called.
func (a *MCPApplication) CaptureJobEmail(input mcpuc.CaptureJobEmailInput) (jobemail.ExtractOutput, error) {
	if strings.TrimSpace(input.Text) == "" {
		return jobemail.ExtractOutput{}, errors.New("text is required")
	}
	return a.extract.Execute(jobemail.ExtractInput{
		Subject:     input.Subject,
		Text:        input.Text,
		CompanyName: input.CompanyName,
		Now:         a.now(),
	}), nil
}

// UpsertEntrySelectionFlow previews or saves a custom selection flow through the REST API.
func (a *MCPApplication) UpsertEntrySelectionFlow(ctx context.Context, input mcpuc.UpsertEntrySelectionFlowInput) (any, error) {
	entryID := strings.TrimSpace(input.EntryID)
	if entryID == "" {
		return nil, errors.New("entryId is required")
	}
	req := selectionFlowRequest{
		Source:               defaultString(input.Source, "ai_paste"),
		CurrentStagePosition: optionalInt(input.CurrentStagePosition),
		Confidence:           input.Confidence,
		InboxClipID:          optionalString(input.InboxClipID),
		Stages:               selectionStageRequests(input.Stages),
	}
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "upsert_entry_selection_flow",
		"entryId":              entryID,
		"selectionFlow":        req,
	}
	if !input.Confirm {
		return preview, nil
	}
	var flow selectionFlowResponse
	if err := a.put(ctx, "/api/v1/entries/"+url.PathEscape(entryID)+"/selection-flow", req, &flow); err != nil {
		return nil, err
	}
	return map[string]any{
		"updated":       true,
		"selectionFlow": flow.toMCP(),
	}, nil
}

// CreateEntryFromJobPosting previews or creates an Entry and its custom selection flow through the REST API.
func (a *MCPApplication) CreateEntryFromJobPosting(ctx context.Context, input mcpuc.CreateEntryFromJobPostingInput) (any, error) {
	companyName := strings.TrimSpace(input.CompanyName)
	if companyName == "" {
		return nil, errors.New("companyName is required")
	}
	entryReq := createEntryWithCompanyRequest{
		CompanyName: companyName,
		Route:       defaultString(input.Route, "本選考"),
		Source:      defaultString(input.Source, "求人ページ"),
		SourceURL:   optionalString(input.SourceURL),
		Memo:        optionalString(input.Memo),
	}
	flowReq := selectionFlowRequest{
		Source:               defaultString(input.FlowSource, "ai_paste"),
		CurrentStagePosition: optionalInt(input.CurrentStagePosition),
		Confidence:           input.Confidence,
		InboxClipID:          optionalString(input.InboxClipID),
		Stages:               selectionStageRequests(input.Stages),
	}
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "create_entry_from_job_posting",
		"entry":                entryReq,
		"selectionFlow":        flowReq,
	}
	if !input.Confirm {
		return preview, nil
	}
	var entry entryResponse
	if err := a.post(ctx, "/api/v1/entries/with-company", entryReq, &entry); err != nil {
		return nil, err
	}
	var flow selectionFlowResponse
	if err := a.put(ctx, "/api/v1/entries/"+url.PathEscape(entry.ID)+"/selection-flow", flowReq, &flow); err != nil {
		return nil, err
	}
	return map[string]any{
		"created":       true,
		"entry":         entry.toMCP(""),
		"selectionFlow": flow.toMCP(),
	}, nil
}

func (a *MCPApplication) listEntries(ctx context.Context) ([]entryResponse, error) {
	var entries listEntriesResponse
	if err := a.get(ctx, "/api/v1/entries", &entries); err != nil {
		return nil, err
	}
	return entries.Entries, nil
}

func (a *MCPApplication) companyNames(ctx context.Context) (map[string]string, error) {
	var companies listCompaniesResponse
	if err := a.get(ctx, "/api/v1/companies", &companies); err != nil {
		return nil, err
	}
	out := make(map[string]string, len(companies.Companies))
	for _, company := range companies.Companies {
		out[company.ID] = company.Name
	}
	return out, nil
}

func (a *MCPApplication) companyName(ctx context.Context, companyID string) (string, error) {
	if strings.TrimSpace(companyID) == "" {
		return "", nil
	}
	var company companyResponse
	if err := a.get(ctx, "/api/v1/companies/"+url.PathEscape(companyID), &company); err != nil {
		return "", err
	}
	return company.Name, nil
}

func (a *MCPApplication) entryCompanyNames(ctx context.Context) (map[string]string, error) {
	entries, err := a.listEntries(ctx)
	if err != nil {
		return nil, err
	}
	companies, err := a.companyNames(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(entries))
	for _, entry := range entries {
		out[entry.ID] = companies[entry.CompanyID]
	}
	return out, nil
}

func (a *MCPApplication) get(ctx context.Context, path string, out any) error {
	return a.do(ctx, http.MethodGet, path, nil, out)
}

func (a *MCPApplication) post(ctx context.Context, path string, body any, out any) error {
	return a.do(ctx, http.MethodPost, path, body, out)
}

func (a *MCPApplication) put(ctx context.Context, path string, body any, out any) error {
	return a.do(ctx, http.MethodPut, path, body, out)
}

func (a *MCPApplication) delete(ctx context.Context, path string) error {
	return a.do(ctx, http.MethodDelete, path, nil, nil)
}

func (a *MCPApplication) do(ctx context.Context, method string, path string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		var buf bytes.Buffer
		if err := json.NewEncoder(&buf).Encode(body); err != nil {
			return err
		}
		reader = &buf
	}
	req, err := http.NewRequestWithContext(ctx, method, a.baseURL+path, reader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+a.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	res, err := a.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = res.Body.Close()
	}()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return apiError(res)
	}
	if res.StatusCode == http.StatusNoContent || out == nil {
		return nil
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func apiError(res *http.Response) error {
	raw, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
	var body struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(raw, &body); err == nil && body.Message != "" {
		return fmt.Errorf("api %s: %s", res.Status, body.Message)
	}
	if len(raw) > 0 {
		return fmt.Errorf("api %s: %s", res.Status, strings.TrimSpace(string(raw)))
	}
	return fmt.Errorf("api %s", res.Status)
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func optionalString(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
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

type listEntriesResponse struct {
	Entries []entryResponse `json:"entries"`
}

type entryResponse struct {
	ID         string `json:"id"`
	CompanyID  string `json:"companyId"`
	Route      string `json:"route"`
	Source     string `json:"source"`
	SourceURL  string `json:"sourceUrl"`
	Status     string `json:"status"`
	StageKind  string `json:"stageKind"`
	StageLabel string `json:"stageLabel"`
	Memo       string `json:"memo"`
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt"`
}

func (e entryResponse) toMCP(company string) mcpuc.EntryDTO {
	return mcpuc.EntryDTO{
		ID:         e.ID,
		CompanyID:  e.CompanyID,
		Company:    company,
		Route:      e.Route,
		Source:     e.Source,
		SourceURL:  e.SourceURL,
		Status:     e.Status,
		StageKind:  e.StageKind,
		StageLabel: e.StageLabel,
		Memo:       e.Memo,
		CreatedAt:  optionalString(e.CreatedAt),
		UpdatedAt:  optionalString(e.UpdatedAt),
	}
}

type listCompaniesResponse struct {
	Companies []companyResponse `json:"companies"`
}

type companyResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type listTasksResponse struct {
	Tasks []taskResponse `json:"tasks"`
}

type taskResponse struct {
	ID        string  `json:"id"`
	EntryID   string  `json:"entryId"`
	Title     string  `json:"title"`
	Type      string  `json:"type"`
	DueDate   *string `json:"dueDate"`
	Status    string  `json:"status"`
	Notify    bool    `json:"notify"`
	Memo      string  `json:"memo"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

func (t taskResponse) toMCP(company string) mcpuc.TaskDTO {
	return mcpuc.TaskDTO{
		ID:        t.ID,
		EntryID:   t.EntryID,
		Company:   company,
		Title:     t.Title,
		Type:      t.Type,
		DueDate:   t.DueDate,
		Status:    t.Status,
		Notify:    t.Notify,
		Memo:      t.Memo,
		CreatedAt: optionalString(t.CreatedAt),
		UpdatedAt: optionalString(t.UpdatedAt),
	}
}

type listInboxClipsResponse struct {
	Clips []inboxClipResponse `json:"clips"`
}

type inboxClipResponse struct {
	ID          string `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Source      string `json:"source"`
	Guess       string `json:"guess"`
	ContentText string `json:"contentText"`
	CapturedAt  string `json:"capturedAt"`
}

func (c inboxClipResponse) toMCP() mcpuc.InboxClipDTO {
	return mcpuc.InboxClipDTO{
		ID:          c.ID,
		URL:         c.URL,
		Title:       c.Title,
		Source:      c.Source,
		Guess:       c.Guess,
		ContentText: c.ContentText,
		CapturedAt:  optionalString(c.CapturedAt),
	}
}

type listESMemosResponse struct {
	Memos []esMemoResponse `json:"memos"`
}

type createTaskRequest struct {
	Title   string  `json:"title"`
	Type    string  `json:"type"`
	DueDate *string `json:"dueDate,omitempty"`
	Memo    *string `json:"memo,omitempty"`
	Notify  *bool   `json:"notify,omitempty"`
}

type createEntryWithCompanyRequest struct {
	CompanyName string  `json:"companyName"`
	Route       string  `json:"route"`
	Source      string  `json:"source"`
	SourceURL   *string `json:"sourceUrl,omitempty"`
	Memo        *string `json:"memo,omitempty"`
}

type selectionFlowRequest struct {
	Source               string                  `json:"source"`
	CurrentStagePosition *int                    `json:"currentStagePosition,omitempty"`
	Confidence           *int                    `json:"confidence,omitempty"`
	InboxClipID          *string                 `json:"inboxClipId,omitempty"`
	Stages               []selectionStageRequest `json:"stages"`
}

type selectionStageRequest struct {
	StageKind    string  `json:"stageKind"`
	StageLabel   string  `json:"stageLabel"`
	EvidenceText *string `json:"evidenceText,omitempty"`
}

type selectionFlowResponse struct {
	ID                   string                   `json:"id"`
	EntryID              string                   `json:"entryId"`
	Source               string                   `json:"source"`
	CurrentStagePosition int                      `json:"currentStagePosition"`
	Confidence           *int                     `json:"confidence"`
	InboxClipID          *string                  `json:"inboxClipId"`
	Stages               []selectionStageResponse `json:"stages"`
	CreatedAt            string                   `json:"createdAt"`
	UpdatedAt            string                   `json:"updatedAt"`
}

type selectionStageResponse struct {
	ID           string `json:"id"`
	Position     int    `json:"position"`
	StageKind    string `json:"stageKind"`
	StageLabel   string `json:"stageLabel"`
	EvidenceText string `json:"evidenceText"`
}

func (f selectionFlowResponse) toMCP() mcpuc.SelectionFlowDTO {
	stages := make([]mcpuc.SelectionStageDTO, 0, len(f.Stages))
	for _, stage := range f.Stages {
		stages = append(stages, mcpuc.SelectionStageDTO{
			ID:           stage.ID,
			Position:     stage.Position,
			StageKind:    stage.StageKind,
			StageLabel:   stage.StageLabel,
			EvidenceText: stage.EvidenceText,
		})
	}
	return mcpuc.SelectionFlowDTO{
		ID:                   f.ID,
		EntryID:              f.EntryID,
		Source:               f.Source,
		CurrentStagePosition: f.CurrentStagePosition,
		Confidence:           f.Confidence,
		InboxClipID:          f.InboxClipID,
		Stages:               stages,
		CreatedAt:            optionalString(f.CreatedAt),
		UpdatedAt:            optionalString(f.UpdatedAt),
	}
}

func selectionStageRequests(stages []mcpuc.SelectionStageDTO) []selectionStageRequest {
	out := make([]selectionStageRequest, 0, len(stages))
	for _, stage := range stages {
		out = append(out, selectionStageRequest{
			StageKind:    stage.StageKind,
			StageLabel:   stage.StageLabel,
			EvidenceText: optionalString(stage.EvidenceText),
		})
	}
	return out
}

func optionalInt(value int) *int {
	if value == 0 {
		return nil
	}
	return &value
}

type createESMemoRequest struct {
	EntryID  *string `json:"entryId,omitempty"`
	Category *string `json:"category,omitempty"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	Source   *string `json:"source,omitempty"`
}

type esMemoResponse struct {
	ID        string  `json:"id"`
	EntryID   *string `json:"entryId"`
	Category  string  `json:"category"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	Source    string  `json:"source"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

func (m esMemoResponse) toMap() map[string]any {
	return map[string]any{
		"id":        m.ID,
		"entryId":   m.EntryID,
		"category":  m.Category,
		"title":     m.Title,
		"content":   m.Content,
		"source":    m.Source,
		"createdAt": m.CreatedAt,
		"updatedAt": m.UpdatedAt,
	}
}

func (m esMemoResponse) toMCP(entryCompanies map[string]string) mcpuc.ESMemoDTO {
	var company *string
	if m.EntryID != nil {
		if value := strings.TrimSpace(entryCompanies[*m.EntryID]); value != "" {
			company = &value
		}
	}
	return mcpuc.ESMemoDTO{
		ID:        m.ID,
		EntryID:   m.EntryID,
		Company:   company,
		Category:  m.Category,
		Title:     m.Title,
		Content:   m.Content,
		Source:    m.Source,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
