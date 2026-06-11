package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestRunWithIO_APIBridgeServesListESMemos(t *testing.T) {
	var sawBearer bool
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/es-memos", func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-token" {
			t.Fatalf("Authorization = %q, want Bearer test-token", got)
		}
		sawBearer = true
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want GET", r.Method)
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Fatalf("limit = %q, want 1", got)
		}
		writeJSON(t, w, map[string]any{
			"memos": []map[string]any{{
				"id":        "memo-1",
				"entryId":   nil,
				"category":  "interview",
				"title":     "改善経験",
				"content":   "顧客課題を分解した",
				"source":    "mcp",
				"createdAt": "2026-06-11T00:00:00Z",
				"updatedAt": "2026-06-11T00:00:00Z",
			}},
		})
	})
	api := httptest.NewServer(mux)
	defer api.Close()

	input := newMCPMessage(`{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"list_es_memos","arguments":{"limit":1}}}`)
	var output bytes.Buffer
	err := runWithIO(context.Background(), func(key string) string {
		switch key {
		case "ENTRE_API_BASE_URL":
			return api.URL
		case "ENTRE_API_TOKEN":
			return "test-token"
		default:
			return ""
		}
	}, strings.NewReader(input), &output)
	if err != nil {
		t.Fatalf("runWithIO() failed: %v", err)
	}
	if !sawBearer {
		t.Fatal("mock API was not called")
	}

	resp := decodeMCPResponse(t, output.String())
	result := resp["result"].(map[string]any)
	content := result["content"].([]any)
	text := content[0].(map[string]any)["text"].(string)
	if !strings.Contains(text, "改善経験") {
		t.Fatalf("MCP response text = %s, want memo title", text)
	}
}

func newMCPMessage(payload string) string {
	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len([]byte(payload)), payload)
}

func decodeMCPResponse(t *testing.T, raw string) map[string]any {
	t.Helper()
	header, payload, ok := strings.Cut(raw, "\r\n\r\n")
	if !ok {
		t.Fatalf("response missing header separator: %q", raw)
	}
	var contentLength int
	for _, line := range strings.Split(header, "\r\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok || !strings.EqualFold(strings.TrimSpace(key), "Content-Length") {
			continue
		}
		n, err := strconv.Atoi(strings.TrimSpace(value))
		if err != nil {
			t.Fatalf("invalid Content-Length %q: %v", value, err)
		}
		contentLength = n
	}
	if contentLength == 0 {
		t.Fatal("response missing Content-Length")
	}
	if len([]byte(payload)) != contentLength {
		t.Fatalf("payload length = %d, want %d", len([]byte(payload)), contentLength)
	}
	var resp map[string]any
	if err := json.Unmarshal([]byte(payload), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["error"] != nil {
		t.Fatalf("unexpected MCP error: %v", resp["error"])
	}
	return resp
}

func writeJSON(t *testing.T, w http.ResponseWriter, body any) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Fatalf("encode response: %v", err)
	}
}
