package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	jobemail "github.com/karimiku/job-hunting-saas/internal/usecase/job_email"
	mcpuc "github.com/karimiku/job-hunting-saas/internal/usecase/mcp"
)

type remoteApplication struct {
	baseURL string
	token   string
	client  *http.Client
	extract *jobemail.Extract
	now     func() time.Time
}

type remoteEntry struct {
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

type remoteCompany struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type remoteTask struct {
	ID        string  `json:"id"`
	EntryID   string  `json:"entryId"`
	Title     string  `json:"title"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	DueDate   *string `json:"dueDate"`
	Notify    bool    `json:"notify"`
	Memo      string  `json:"memo"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
}

type remoteInboxClip struct {
	ID         string `json:"id"`
	URL        string `json:"url"`
	Title      string `json:"title"`
	Source     string `json:"source"`
	Guess      string `json:"guess"`
	CapturedAt string `json:"capturedAt"`
}

func usesRemoteAPIEnv() bool {
	return strings.TrimSpace(os.Getenv("ENTRE_API_BASE_URL")) != "" ||
		strings.TrimSpace(os.Getenv("ENTRE_API_TOKEN")) != ""
}

func newRemoteApplicationFromEnv() (*remoteApplication, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(os.Getenv("ENTRE_API_BASE_URL")), "/")
	if baseURL == "" {
		return nil, errors.New("ENTRE_API_BASE_URL is required when ENTRE_API_TOKEN is set")
	}
	token := strings.TrimSpace(os.Getenv("ENTRE_API_TOKEN"))
	if token == "" {
		return nil, errors.New("ENTRE_API_TOKEN is required when ENTRE_API_BASE_URL is set")
	}
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("invalid ENTRE_API_BASE_URL: %w", err)
	}
	return &remoteApplication{
		baseURL: baseURL,
		token:   token,
		client:  &http.Client{Timeout: 20 * time.Second},
		extract: jobemail.NewExtract(),
		now:     time.Now,
	}, nil
}

func (a *remoteApplication) ListEntries(ctx context.Context) ([]mcpuc.EntryDTO, error) {
	var res struct {
		Entries []remoteEntry `json:"entries"`
	}
	if err := a.getJSON(ctx, "/api/v1/entries", &res); err != nil {
		return nil, err
	}
	companies, err := a.companyMap(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]mcpuc.EntryDTO, 0, len(res.Entries))
	for _, entry := range res.Entries {
		out = append(out, entry.remoteDTO(companies[entry.CompanyID]))
	}
	return out, nil
}

func (a *remoteApplication) GetEntryContext(ctx context.Context, rawEntryID string) (*mcpuc.EntryContextDTO, error) {
	entryID := strings.TrimSpace(rawEntryID)
	if entryID == "" {
		return nil, errors.New("entryId is required")
	}
	var entry remoteEntry
	if err := a.getJSON(ctx, "/api/v1/entries/"+url.PathEscape(entryID), &entry); err != nil {
		return nil, err
	}
	company, err := a.getCompany(ctx, entry.CompanyID)
	if err != nil {
		return nil, err
	}
	var tasksRes struct {
		Tasks []remoteTask `json:"tasks"`
	}
	if err := a.getJSON(ctx, "/api/v1/entries/"+url.PathEscape(entryID)+"/tasks", &tasksRes); err != nil {
		return nil, err
	}
	tasks := make([]mcpuc.TaskDTO, 0, len(tasksRes.Tasks))
	for _, task := range tasksRes.Tasks {
		tasks = append(tasks, task.remoteDTO(company.Name))
	}
	return &mcpuc.EntryContextDTO{
		Entry: entry.remoteDTO(company.Name),
		Tasks: tasks,
	}, nil
}

func (a *remoteApplication) ListOpenTasks(ctx context.Context) ([]mcpuc.TaskDTO, error) {
	var tasksRes struct {
		Tasks []remoteTask `json:"tasks"`
	}
	if err := a.getJSON(ctx, "/api/v1/tasks", &tasksRes); err != nil {
		return nil, err
	}
	entryCompanies, err := a.entryCompanyMap(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]mcpuc.TaskDTO, 0, len(tasksRes.Tasks))
	for _, task := range tasksRes.Tasks {
		if task.Status != "todo" {
			continue
		}
		out = append(out, task.remoteDTO(entryCompanies[task.EntryID]))
	}
	return out, nil
}

func (a *remoteApplication) ListInboxClips(ctx context.Context) ([]mcpuc.InboxClipDTO, error) {
	var res struct {
		Clips []remoteInboxClip `json:"clips"`
	}
	if err := a.getJSON(ctx, "/api/v1/inbox/clips", &res); err != nil {
		return nil, err
	}
	out := make([]mcpuc.InboxClipDTO, 0, len(res.Clips))
	for _, clip := range res.Clips {
		out = append(out, mcpuc.InboxClipDTO{
			ID:         clip.ID,
			URL:        clip.URL,
			Title:      clip.Title,
			Source:     clip.Source,
			Guess:      clip.Guess,
			CapturedAt: stringPtr(clip.CapturedAt),
		})
	}
	return out, nil
}

func (a *remoteApplication) AppendESMemo(_ context.Context, input mcpuc.AppendESMemoInput) (any, error) {
	title := strings.TrimSpace(input.Title)
	content := strings.TrimSpace(input.Content)
	if title == "" || content == "" {
		return nil, errors.New("title and content are required")
	}
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "append_es_memo",
		"memo": map[string]any{
			"title":    title,
			"content":  content,
			"category": defaultRemoteString(input.Category, "general"),
			"entryId":  optionalRemoteString(input.EntryID),
			"source":   defaultRemoteString(input.Source, "mcp"),
		},
	}
	if !input.Confirm {
		return preview, nil
	}
	return nil, errors.New("append_es_memo is not supported in ENTRE_API_BASE_URL mode yet")
}

func (a *remoteApplication) CreateTask(ctx context.Context, input mcpuc.CreateTaskInput) (any, error) {
	entryID := strings.TrimSpace(input.EntryID)
	title := strings.TrimSpace(input.Title)
	if entryID == "" || title == "" {
		return nil, errors.New("entryId and title are required")
	}
	taskType := defaultRemoteString(input.Type, "deadline")
	dueDate, err := parseRemoteDueDate(input.DueDate)
	if err != nil {
		return nil, err
	}
	entryCtx, err := a.GetEntryContext(ctx, entryID)
	if err != nil {
		return nil, err
	}
	preview := map[string]any{
		"confirmationRequired": !input.Confirm,
		"action":               "create_task",
		"task": map[string]any{
			"entryId": entryID,
			"company": entryCtx.Entry.Company,
			"title":   title,
			"type":    taskType,
			"dueDate": timeStringPtr(dueDate),
			"memo":    input.Memo,
			"notify":  input.Notify,
		},
	}
	if !input.Confirm {
		return preview, nil
	}
	body := map[string]any{
		"title": title,
		"type":  taskType,
		"memo":  input.Memo,
	}
	if dueDate != nil {
		body["dueDate"] = dueDate.Format(time.RFC3339)
	}
	var created remoteTask
	if err := a.postJSON(ctx, "/api/v1/entries/"+url.PathEscape(entryID)+"/tasks", body, &created); err != nil {
		return nil, err
	}
	if input.Notify && !created.Notify {
		var updated remoteTask
		if err := a.patchJSON(ctx, "/api/v1/tasks/"+url.PathEscape(created.ID), map[string]any{"notify": true}, &updated); err != nil {
			return nil, err
		}
		created = updated
	}
	task := created.remoteDTO(entryCtx.Entry.Company)
	return map[string]any{
		"created": true,
		"task": map[string]any{
			"id":      task.ID,
			"entryId": task.EntryID,
			"company": task.Company,
			"title":   task.Title,
			"type":    task.Type,
			"dueDate": task.DueDate,
			"status":  task.Status,
			"notify":  task.Notify,
			"memo":    task.Memo,
		},
	}, nil
}

func (a *remoteApplication) CaptureJobEmail(input mcpuc.CaptureJobEmailInput) (jobemail.ExtractOutput, error) {
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

func (a *remoteApplication) getJSON(ctx context.Context, path string, out any) error {
	return a.doJSON(ctx, http.MethodGet, path, nil, out)
}

func (a *remoteApplication) postJSON(ctx context.Context, path string, body any, out any) error {
	return a.doJSON(ctx, http.MethodPost, path, body, out)
}

func (a *remoteApplication) patchJSON(ctx context.Context, path string, body any, out any) error {
	return a.doJSON(ctx, http.MethodPatch, path, body, out)
}

func (a *remoteApplication) doJSON(ctx context.Context, method, path string, body any, out any) error {
	var reqBody *bytes.Reader
	if body == nil {
		reqBody = bytes.NewReader(nil)
	} else {
		payload, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewReader(payload)
	}
	req, err := http.NewRequestWithContext(ctx, method, a.baseURL+path, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+a.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errBody struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(resp.Body).Decode(&errBody)
		if errBody.Message == "" {
			errBody.Message = resp.Status
		}
		return fmt.Errorf("%s %s failed: %s", method, path, errBody.Message)
	}
	if out == nil {
		return nil
	}
	return json.NewDecoder(resp.Body).Decode(out)
}

func (a *remoteApplication) companyMap(ctx context.Context) (map[string]string, error) {
	var res struct {
		Companies []remoteCompany `json:"companies"`
	}
	if err := a.getJSON(ctx, "/api/v1/companies", &res); err != nil {
		return nil, err
	}
	companies := make(map[string]string, len(res.Companies))
	for _, company := range res.Companies {
		companies[company.ID] = company.Name
	}
	return companies, nil
}

func (a *remoteApplication) getCompany(ctx context.Context, companyID string) (remoteCompany, error) {
	if strings.TrimSpace(companyID) == "" {
		return remoteCompany{}, errors.New("companyId is required")
	}
	var company remoteCompany
	if err := a.getJSON(ctx, "/api/v1/companies/"+url.PathEscape(companyID), &company); err != nil {
		return remoteCompany{}, err
	}
	return company, nil
}

func (a *remoteApplication) entryCompanyMap(ctx context.Context) (map[string]string, error) {
	var entriesRes struct {
		Entries []remoteEntry `json:"entries"`
	}
	if err := a.getJSON(ctx, "/api/v1/entries", &entriesRes); err != nil {
		return nil, err
	}
	companies, err := a.companyMap(ctx)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string, len(entriesRes.Entries))
	for _, entry := range entriesRes.Entries {
		out[entry.ID] = companies[entry.CompanyID]
	}
	return out, nil
}

func (e remoteEntry) remoteDTO(company string) mcpuc.EntryDTO {
	if company == "" {
		company = e.CompanyID
	}
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
		CreatedAt:  stringPtr(e.CreatedAt),
		UpdatedAt:  stringPtr(e.UpdatedAt),
	}
}

func (t remoteTask) remoteDTO(company string) mcpuc.TaskDTO {
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
		CreatedAt: stringPtr(t.CreatedAt),
		UpdatedAt: stringPtr(t.UpdatedAt),
	}
}

func parseRemoteDueDate(raw string) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t, nil
	}
	if t, err := time.Parse("2006-01-02", raw); err == nil {
		return &t, nil
	}
	return nil, fmt.Errorf("invalid dueDate %q: use YYYY-MM-DD or RFC3339", raw)
}

func defaultRemoteString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func optionalRemoteString(value string) *string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	v := strings.TrimSpace(value)
	return &v
}

func stringPtr(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func timeStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	value := t.Format(time.RFC3339)
	return &value
}
