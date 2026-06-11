package entreapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mcpuc "github.com/karimiku/job-hunting-saas/internal/usecase/mcp"
)

func TestMCPApplication_ListEntriesAndOpenTasks(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/entries", func(w http.ResponseWriter, r *http.Request) {
		assertBearer(t, r)
		writeTestJSON(t, w, map[string]any{
			"entries": []map[string]any{{
				"id":         "entry-1",
				"companyId":  "company-1",
				"route":      "本選考",
				"source":     "マイナビ",
				"sourceUrl":  "https://example.com/job",
				"status":     "in_progress",
				"stageKind":  "document",
				"stageLabel": "ES",
				"memo":       "memo",
				"createdAt":  "2026-06-11T00:00:00Z",
				"updatedAt":  "2026-06-11T00:00:00Z",
			}},
		})
	})
	mux.HandleFunc("/api/v1/companies", func(w http.ResponseWriter, r *http.Request) {
		assertBearer(t, r)
		writeTestJSON(t, w, map[string]any{
			"companies": []map[string]any{{
				"id":   "company-1",
				"name": "Example Inc.",
			}},
		})
	})
	mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		assertBearer(t, r)
		writeTestJSON(t, w, map[string]any{
			"tasks": []map[string]any{
				{
					"id":        "task-1",
					"entryId":   "entry-1",
					"title":     "ES提出",
					"type":      "deadline",
					"dueDate":   nil,
					"status":    "todo",
					"notify":    true,
					"memo":      "",
					"createdAt": "2026-06-11T00:00:00Z",
					"updatedAt": "2026-06-11T00:00:00Z",
				},
				{
					"id":        "task-2",
					"entryId":   "entry-1",
					"title":     "完了済み",
					"type":      "deadline",
					"dueDate":   nil,
					"status":    "done",
					"notify":    false,
					"memo":      "",
					"createdAt": "2026-06-11T00:00:00Z",
					"updatedAt": "2026-06-11T00:00:00Z",
				},
			},
		})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	app, err := NewMCPApplication(server.URL, "test-token", server.Client())
	if err != nil {
		t.Fatalf("NewMCPApplication() failed: %v", err)
	}

	entries, err := app.ListEntries(t.Context())
	if err != nil {
		t.Fatalf("ListEntries() failed: %v", err)
	}
	if len(entries) != 1 || entries[0].Company != "Example Inc." {
		t.Fatalf("entries = %#v", entries)
	}

	tasks, err := app.ListOpenTasks(t.Context())
	if err != nil {
		t.Fatalf("ListOpenTasks() failed: %v", err)
	}
	if len(tasks) != 1 {
		t.Fatalf("len(tasks) = %d, want 1", len(tasks))
	}
	if tasks[0].Company != "Example Inc." {
		t.Errorf("task company = %q", tasks[0].Company)
	}
}

func TestMCPApplication_AppendESMemoWithConfirmation(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/es-memos", func(w http.ResponseWriter, r *http.Request) {
		assertBearer(t, r)
		if r.Method != http.MethodPost {
			t.Fatalf("method = %s, want POST", r.Method)
		}
		var req createESMemoRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if req.Title != "改善経験" || req.Content != "顧客課題を分解した" {
			t.Fatalf("request = %#v", req)
		}
		writeTestJSON(t, w, map[string]any{
			"id":        "memo-1",
			"entryId":   nil,
			"category":  "interview",
			"title":     req.Title,
			"content":   req.Content,
			"source":    "mcp",
			"createdAt": "2026-06-11T00:00:00Z",
			"updatedAt": "2026-06-11T00:00:00Z",
		})
	})
	server := httptest.NewServer(mux)
	defer server.Close()

	app, err := NewMCPApplication(server.URL, "test-token", server.Client())
	if err != nil {
		t.Fatalf("NewMCPApplication() failed: %v", err)
	}
	out, err := app.AppendESMemo(t.Context(), mcpuc.AppendESMemoInput{
		Title:    "改善経験",
		Content:  "顧客課題を分解した",
		Category: "interview",
		Confirm:  true,
	})
	if err != nil {
		t.Fatalf("AppendESMemo() failed: %v", err)
	}
	got := out.(map[string]any)
	if got["created"] != true {
		t.Fatalf("created = %v, want true", got["created"])
	}
}

func assertBearer(t *testing.T, r *http.Request) {
	t.Helper()
	if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
		t.Fatalf("Authorization = %q", got)
	}
}

func writeTestJSON(t *testing.T, w http.ResponseWriter, body any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}
