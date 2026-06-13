// Package mcp implements the stdio MCP transport for job-hunting context tools.
package mcp

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/textproto"
	"os"
	"strconv"
	"strings"
	"time"

	jobemail "github.com/karimiku/job-hunting-saas/internal/usecase/job_email"
	mcpuc "github.com/karimiku/job-hunting-saas/internal/usecase/mcp"
)

const (
	defaultProtocolVersion = "2024-11-05"
	serverName             = "job-hunting-saas-mcp"
	serverVersion          = "0.1.0"
)

// Application はMCP handlerが呼び出すユースケース境界。
type Application interface {
	ListEntries(ctx context.Context) ([]mcpuc.EntryDTO, error)
	GetEntryContext(ctx context.Context, rawEntryID string) (*mcpuc.EntryContextDTO, error)
	ListOpenTasks(ctx context.Context) ([]mcpuc.TaskDTO, error)
	ListInboxClips(ctx context.Context) ([]mcpuc.InboxClipDTO, error)
	AppendESMemo(ctx context.Context, input mcpuc.AppendESMemoInput) (any, error)
	CreateTask(ctx context.Context, input mcpuc.CreateTaskInput) (any, error)
	CaptureJobEmail(input mcpuc.CaptureJobEmailInput) (jobemail.ExtractOutput, error)
}

// Server はstdio MCPリクエストを処理するサーバー。
type Server struct {
	app Application
}

type rpcRequest struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id,omitempty"`
	Method  string           `json:"method"`
	Params  json.RawMessage  `json:"params,omitempty"`
}

type rpcResponse struct {
	JSONRPC string           `json:"jsonrpc"`
	ID      *json.RawMessage `json:"id"`
	Result  any              `json:"result,omitempty"`
	Error   *rpcError        `json:"error,omitempty"`
}

type rpcError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type contentItem struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type toolResult struct {
	Content []contentItem `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

type resourceContent struct {
	URI      string `json:"uri"`
	MIMEType string `json:"mimeType"`
	Text     string `json:"text"`
}

// NewServer はMCPサーバーを生成する。
func NewServer(app Application) *Server {
	return &Server{app: app}
}

// ServeStdio はstdio上のMCPメッセージを読み書きしてリクエストを処理する。
func ServeStdio(ctx context.Context, in io.Reader, out io.Writer, s *Server) error {
	reader := bufio.NewReader(in)
	writer := bufio.NewWriter(out)
	debugMCP("stdio server started")
	for {
		payload, err := readMCPMessage(reader)
		if err != nil {
			debugMCP("read error: %v", err)
			return err
		}
		var req rpcRequest
		if err := json.Unmarshal(payload, &req); err != nil {
			debugMCP("json unmarshal error: %v", err)
			continue
		}
		if req.ID == nil {
			debugMCP("notification ignored: method=%s", req.Method)
			continue
		}
		debugMCP("request received: method=%s", req.Method)
		resp := s.handle(ctx, req)
		if err := writeMCPMessage(writer, resp); err != nil {
			debugMCP("write error: method=%s error=%v", req.Method, err)
			return err
		}
		if err := writer.Flush(); err != nil {
			debugMCP("flush error: method=%s error=%v", req.Method, err)
			return err
		}
		debugMCP("response sent: method=%s", req.Method)
	}
}

func readMCPMessage(r *bufio.Reader) ([]byte, error) {
	header := textproto.MIMEHeader{}
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return nil, fmt.Errorf("invalid MCP header %q", line)
		}
		header.Add(key, strings.TrimSpace(value))
	}
	lengthRaw := header.Get("Content-Length")
	if lengthRaw == "" {
		return nil, errors.New("missing Content-Length")
	}
	length, err := strconv.Atoi(lengthRaw)
	if err != nil || length < 0 {
		return nil, fmt.Errorf("invalid Content-Length %q", lengthRaw)
	}
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func writeMCPMessage(w io.Writer, msg any) error {
	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(payload)); err != nil {
		return err
	}
	_, err = w.Write(payload)
	return err
}

func (s *Server) handle(ctx context.Context, req rpcRequest) rpcResponse {
	result, err := s.dispatch(ctx, req)
	resp := rpcResponse{JSONRPC: "2.0", ID: req.ID}
	if err != nil {
		resp.Error = &rpcError{Code: -32000, Message: err.Error()}
		return resp
	}
	resp.Result = result
	return resp
}

func (s *Server) dispatch(ctx context.Context, req rpcRequest) (any, error) {
	switch req.Method {
	case "initialize":
		return initialize(req.Params), nil
	case "resources/list":
		return listResources(), nil
	case "resources/templates/list":
		return listResourceTemplates(), nil
	case "resources/read":
		return s.readResource(ctx, req.Params)
	case "tools/list":
		return listTools(), nil
	case "tools/call":
		return s.callTool(ctx, req.Params)
	case "prompts/list":
		return map[string]any{"prompts": []any{}}, nil
	case "ping":
		return map[string]any{}, nil
	default:
		return nil, fmt.Errorf("unsupported method %q", req.Method)
	}
}

func initialize(params json.RawMessage) any {
	var p struct {
		ProtocolVersion string `json:"protocolVersion"`
	}
	_ = json.Unmarshal(params, &p)
	version := p.ProtocolVersion
	if version == "" {
		version = defaultProtocolVersion
	}
	return map[string]any{
		"protocolVersion": version,
		"capabilities": map[string]any{
			"resources": map[string]any{},
			"tools":     map[string]any{},
		},
		"serverInfo": map[string]any{
			"name":    serverName,
			"version": serverVersion,
		},
	}
}

func listResources() any {
	if strings.TrimSpace(os.Getenv("ENTRE_API_BASE_URL")) != "" {
		return map[string]any{"resources": []map[string]any{}}
	}
	return map[string]any{
		"resources": []map[string]any{
			{"uri": "entries://list", "name": "応募先一覧", "description": "MCP対象ユーザーの応募先一覧", "mimeType": "application/json"},
			{"uri": "tasks://open", "name": "未完了Task一覧", "description": "未完了の締切・予定", "mimeType": "application/json"},
			{"uri": "inbox://clips", "name": "Inbox clip一覧", "description": "Chrome拡張等で保存された求人ページ", "mimeType": "application/json"},
		},
	}
}

func listResourceTemplates() any {
	if strings.TrimSpace(os.Getenv("ENTRE_API_BASE_URL")) != "" {
		return map[string]any{"resourceTemplates": []map[string]any{}}
	}
	return map[string]any{
		"resourceTemplates": []map[string]any{
			{"uriTemplate": "entries://{entryId}", "name": "応募先詳細", "description": "応募先1件の文脈", "mimeType": "application/json"},
		},
	}
}

func (s *Server) readResource(ctx context.Context, params json.RawMessage) (any, error) {
	var p struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	switch {
	case p.URI == "entries://list":
		entries, err := s.app.ListEntries(ctx)
		if err != nil {
			return nil, err
		}
		return resourceJSON(p.URI, entries)
	case strings.HasPrefix(p.URI, "entries://"):
		entryID := strings.TrimPrefix(p.URI, "entries://")
		entry, err := s.app.GetEntryContext(ctx, entryID)
		if err != nil {
			return nil, err
		}
		return resourceJSON(p.URI, entry)
	case p.URI == "tasks://open":
		tasks, err := s.app.ListOpenTasks(ctx)
		if err != nil {
			return nil, err
		}
		return resourceJSON(p.URI, tasks)
	case p.URI == "inbox://clips":
		clips, err := s.app.ListInboxClips(ctx)
		if err != nil {
			return nil, err
		}
		return resourceJSON(p.URI, clips)
	default:
		return nil, fmt.Errorf("unknown resource uri %q", p.URI)
	}
}

func resourceJSON(uri string, value any) (any, error) {
	text, err := marshalPretty(value)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"contents": []resourceContent{{
			URI:      uri,
			MIMEType: "application/json",
			Text:     text,
		}},
	}, nil
}

func listTools() any {
	return map[string]any{
		"tools": []map[string]any{
			{
				"name":        "list_entries",
				"description": "応募先一覧を取得します。",
				"inputSchema": objectSchema(nil, nil),
			},
			{
				"name":        "get_entry_context",
				"description": "応募先1件と紐づくTaskを取得します。",
				"inputSchema": objectSchema(map[string]any{
					"entryId": map[string]any{"type": "string", "description": "Entry UUID"},
				}, []string{"entryId"}),
			},
			{
				"name":        "list_open_tasks",
				"description": "未完了Task一覧を取得します。",
				"inputSchema": objectSchema(nil, nil),
			},
			{
				"name":        "append_es_memo",
				"description": "ES/自己PR/ガクチカ/面接ネタ用メモを追記します。confirm=true のときだけDBへ保存します。",
				"inputSchema": objectSchema(map[string]any{
					"title":    map[string]any{"type": "string"},
					"content":  map[string]any{"type": "string"},
					"category": map[string]any{"type": "string", "description": "general, gakuchika, self_pr, motivation, interview など"},
					"entryId":  map[string]any{"type": "string", "description": "任意のEntry UUID"},
					"source":   map[string]any{"type": "string", "description": "memo, interview, mail など"},
					"confirm":  map[string]any{"type": "boolean", "description": "true のときだけ保存"},
				}, []string{"title", "content"}),
			},
			{
				"name":        "create_task",
				"description": "Entryに紐づくTaskを作成します。confirm=true のときだけDBへ保存します。",
				"inputSchema": objectSchema(map[string]any{
					"entryId": map[string]any{"type": "string"},
					"title":   map[string]any{"type": "string"},
					"type":    map[string]any{"type": "string", "enum": []string{"deadline", "schedule"}},
					"dueDate": map[string]any{"type": "string", "description": "YYYY-MM-DD または RFC3339"},
					"memo":    map[string]any{"type": "string"},
					"notify":  map[string]any{"type": "boolean"},
					"confirm": map[string]any{"type": "boolean", "description": "true のときだけ保存"},
				}, []string{"entryId", "title"}),
			},
			{
				"name":        "capture_job_email",
				"description": "選考メール本文を構造化し、Entry更新候補とTask作成候補を返します。LLM APIは呼びません。",
				"inputSchema": objectSchema(map[string]any{
					"text":        map[string]any{"type": "string"},
					"subject":     map[string]any{"type": "string"},
					"companyName": map[string]any{"type": "string", "description": "分かっている場合の会社名"},
				}, []string{"text"}),
			},
		},
	}
}

func objectSchema(properties map[string]any, required []string) map[string]any {
	if properties == nil {
		properties = map[string]any{}
	}
	schema := map[string]any{
		"type":                 "object",
		"properties":           properties,
		"additionalProperties": false,
	}
	if len(required) > 0 {
		schema["required"] = required
	}
	return schema
}

func (s *Server) callTool(ctx context.Context, params json.RawMessage) (any, error) {
	var p struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return nil, err
	}
	var (
		value any
		err   error
	)
	switch p.Name {
	case "list_entries":
		value, err = s.app.ListEntries(ctx)
	case "get_entry_context":
		var args struct {
			EntryID string `json:"entryId"`
		}
		if err = json.Unmarshal(p.Arguments, &args); err == nil {
			value, err = s.app.GetEntryContext(ctx, args.EntryID)
		}
	case "list_open_tasks":
		value, err = s.app.ListOpenTasks(ctx)
	case "append_es_memo":
		var args mcpuc.AppendESMemoInput
		if err = json.Unmarshal(p.Arguments, &args); err == nil {
			value, err = s.app.AppendESMemo(ctx, args)
		}
	case "create_task":
		var args mcpuc.CreateTaskInput
		if err = json.Unmarshal(p.Arguments, &args); err == nil {
			value, err = s.app.CreateTask(ctx, args)
		}
	case "capture_job_email":
		var args mcpuc.CaptureJobEmailInput
		if err = json.Unmarshal(p.Arguments, &args); err == nil {
			value, err = s.app.CaptureJobEmail(args)
		}
	default:
		err = fmt.Errorf("unknown tool %q", p.Name)
	}
	if err != nil {
		return toolText(map[string]any{"error": err.Error()}, true)
	}
	return toolText(value, false)
}

func toolText(value any, isError bool) (toolResult, error) {
	text, err := marshalPretty(value)
	if err != nil {
		return toolResult{}, err
	}
	return toolResult{
		Content: []contentItem{{Type: "text", Text: text}},
		IsError: isError,
	}, nil
}

func marshalPretty(value any) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return "", err
	}
	return strings.TrimRight(buf.String(), "\n"), nil
}

func debugMCP(format string, args ...any) {
	path := strings.TrimSpace(os.Getenv("ENTRE_MCP_DEBUG_LOG"))
	if path == "" {
		return
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return
	}
	defer f.Close()
	message := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(f, "%s %s\n", time.Now().Format(time.RFC3339Nano), message)
}
