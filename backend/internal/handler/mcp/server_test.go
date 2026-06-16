package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	jobemail "github.com/karimiku/job-hunting-saas/internal/usecase/job_email"
	mcpuc "github.com/karimiku/job-hunting-saas/internal/usecase/mcp"
)

func TestServeStdioInitializeIncludesInstructions(t *testing.T) {
	out := runMCPMessages(t, `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2024-11-05"}}`)
	resp := decodeRPCResponse(t, out[0])

	result := resp["result"].(map[string]any)
	if result["instructions"] == "" {
		t.Fatal("initialize result should include server instructions")
	}
	info := result["serverInfo"].(map[string]any)
	if info["name"] != serverName {
		t.Fatalf("serverInfo.name = %v, want %s", info["name"], serverName)
	}
}

func TestServeStdioListsInboxAndESMemoTools(t *testing.T) {
	out := runMCPMessages(t, `{"jsonrpc":"2.0","id":1,"method":"tools/list"}`)
	resp := decodeRPCResponse(t, out[0])
	result := resp["result"].(map[string]any)
	tools := result["tools"].([]any)

	names := map[string]bool{}
	for _, tool := range tools {
		item := tool.(map[string]any)
		names[item["name"].(string)] = true
	}
	for _, name := range []string{"list_inbox_clips", "list_es_memos", "append_es_memo", "create_task"} {
		if !names[name] {
			t.Fatalf("tool %q is missing from tools/list", name)
		}
	}
}

func TestServeStdioCallsListESMemos(t *testing.T) {
	out := runMCPMessages(t, `{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"list_es_memos","arguments":{"limit":10}}}`)
	resp := decodeRPCResponse(t, out[0])
	result := resp["result"].(map[string]any)
	content := result["content"].([]any)
	text := content[0].(map[string]any)["text"].(string)
	if !strings.Contains(text, "改善経験") {
		t.Fatalf("tool text = %s, want ES memo title", text)
	}
}

func runMCPMessages(t *testing.T, payloads ...string) [][]byte {
	t.Helper()
	var input bytes.Buffer
	for _, payload := range payloads {
		_, _ = fmt.Fprintf(&input, "Content-Length: %d\r\n\r\n%s", len([]byte(payload)), payload)
	}

	var output bytes.Buffer
	err := ServeStdio(context.Background(), &input, &output, NewServer(fakeApplication{}))
	if err != nil && !errors.Is(err, io.EOF) {
		t.Fatalf("ServeStdio() failed: %v", err)
	}

	reader := bufio.NewReader(&output)
	messages := make([][]byte, 0, len(payloads))
	for range payloads {
		msg, err := readMCPMessage(reader)
		if err != nil {
			t.Fatalf("readMCPMessage(output) failed: %v", err)
		}
		messages = append(messages, msg)
	}
	return messages
}

func decodeRPCResponse(t *testing.T, raw []byte) map[string]any {
	t.Helper()
	var resp map[string]any
	if err := json.Unmarshal(raw, &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp["error"] != nil {
		t.Fatalf("unexpected rpc error: %v", resp["error"])
	}
	return resp
}

type fakeApplication struct{}

func (fakeApplication) ListEntries(context.Context) ([]mcpuc.EntryDTO, error) {
	return []mcpuc.EntryDTO{{ID: "entry-1", Company: "Example Inc."}}, nil
}

func (fakeApplication) GetEntryContext(context.Context, string) (*mcpuc.EntryContextDTO, error) {
	return &mcpuc.EntryContextDTO{
		Entry: mcpuc.EntryDTO{ID: "entry-1", Company: "Example Inc."},
		Tasks: []mcpuc.TaskDTO{},
	}, nil
}

func (fakeApplication) ListOpenTasks(context.Context) ([]mcpuc.TaskDTO, error) {
	return []mcpuc.TaskDTO{{ID: "task-1", Title: "ES提出"}}, nil
}

func (fakeApplication) ListInboxClips(context.Context) ([]mcpuc.InboxClipDTO, error) {
	return []mcpuc.InboxClipDTO{{ID: "clip-1", URL: "https://example.com/job"}}, nil
}

func (fakeApplication) ListESMemos(context.Context, int32) ([]mcpuc.ESMemoDTO, error) {
	company := "Example Inc."
	return []mcpuc.ESMemoDTO{{
		ID:       "memo-1",
		Company:  &company,
		Category: "interview",
		Title:    "改善経験",
		Content:  "顧客課題を分解した",
		Source:   "mcp",
	}}, nil
}

func (fakeApplication) AppendESMemo(context.Context, mcpuc.AppendESMemoInput) (any, error) {
	return map[string]any{"created": true}, nil
}

func (fakeApplication) CreateTask(context.Context, mcpuc.CreateTaskInput) (any, error) {
	return map[string]any{"created": true}, nil
}

func (fakeApplication) CaptureJobEmail(mcpuc.CaptureJobEmailInput) (jobemail.ExtractOutput, error) {
	return jobemail.ExtractOutput{CompanyName: "Example Inc."}, nil
}

func (fakeApplication) UpsertEntrySelectionFlow(context.Context, mcpuc.UpsertEntrySelectionFlowInput) (any, error) {
	return map[string]any{"updated": true}, nil
}

func (fakeApplication) CreateEntryFromJobPosting(context.Context, mcpuc.CreateEntryFromJobPostingInput) (any, error) {
	return map[string]any{"created": true}, nil
}
