# MCP server

job-hunting-saas は stdio MCP server を提供し、Claude Desktop / Codex / Gemini CLI などのMCPクライアントから就活データを参照・更新できる。
アプリ本体はLLM APIを呼ばない。料金・モデル選択・APIキー管理はユーザーが使うMCPクライアント側に寄せる。

## 想定ユースケース

- 保存箱に入った求人ページをAIが `inbox://clips` から読み取り、Entry化や調査タスク作成の候補を作る
- Entry一覧・未完了Task・応募先詳細をAIが引いて、今日やるべき締切や面接準備を整理する
- メール本文を `capture_job_email` に渡して、面接日程・提出締切・持ち物などをTask候補として構造化する
- 自己PR・ガクチカ・面接で出た話を `append_es_memo` でESメモとして蓄積する

## Architecture

MCP server も既存APIと同じ Clean Architecture の境界に合わせる。

| 層 | 実装 | 責務 |
| --- | --- | --- |
| cmd | `backend/cmd/mcp-server` | 環境変数読み込み、API bridge / DB direct のDI |
| Handler | `backend/internal/handler/mcp` | MCP JSON-RPC over stdio、resources/tools の入出力変換 |
| UseCase | `backend/internal/usecase/mcp`, `es_memo`, `job_email`, `task` | user scoped な操作、保存前preview、メール抽出 |
| Domain | `backend/internal/domain/entity`, `value`, `repository` | ESMemo entity/value object、Repository interface |
| Infra | `backend/internal/infra/entreapi`, `backend/internal/infra/postgres` | Hosted API client / sqlc query / repository 実装 |

## Run

通常利用は `ENTRE_API_TOKEN` を使う API bridge mode を使う。Webアプリの「アカウント」画面で AI連携トークンを作成し、MCPクライアント設定の環境変数に渡す。

Codex など公式SDK準拠のstdioを期待するクライアントでは、Node wrapper `backend/cmd/mcp-remote/entre-mcp.mjs` を使う。Go製 `backend/cmd/mcp-server` はDB直結モードや既存検証用として残す。

```bash
ENTRE_API_BASE_URL=https://api.example.com \
ENTRE_API_TOKEN=entre_ai_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
node backend/cmd/mcp-remote/entre-mcp.mjs
```

wrapperは `frontend` の依存に入っている `@modelcontextprotocol/sdk` を使うため、初回は `cd frontend && pnpm install` を実行しておく。

まずローカルMCPサーバーの単一バイナリを作る。

```bash
make build-mcp-server
```

出力先は `backend/bin/mcp-server`。Go版をGUIアプリから起動する場合は、`go run` ではなくこのバイナリの絶対パスを設定に入れる。

```bash
cd backend
ENTRE_API_BASE_URL=http://localhost:8080 \
ENTRE_API_TOKEN=entre_ai_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
./bin/mcp-server
```

`ENTRE_API_BASE_URL` は省略時 `http://localhost:8080`。本番配布では hosted API のURLを指定する。

## Client setup

### Codex CLI / Codex IDE

Codex は `codex mcp add` か `~/.codex/config.toml` で stdio MCP server を設定する。CLI と IDE extension は同じ設定を読む。

```bash
codex mcp add entre \
  --env ENTRE_API_BASE_URL=https://api.example.com \
  --env ENTRE_API_TOKEN=entre_ai_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  -- node /absolute/path/to/backend/cmd/mcp-remote/entre-mcp.mjs
```

`config.toml` に直接書く場合:

```toml
[mcp_servers.entre]
command = "node"
args = ["/absolute/path/to/backend/cmd/mcp-remote/entre-mcp.mjs"]

[mcp_servers.entre.env]
ENTRE_API_BASE_URL = "https://api.example.com"
ENTRE_API_TOKEN = "entre_ai_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

### Claude Code

Claude Code は `claude mcp add` で stdio MCP server を追加する。個人利用のMVP検証では `--scope user` が扱いやすい。

```bash
claude mcp add --transport stdio --scope user entre \
  --env ENTRE_API_BASE_URL=https://api.example.com \
  --env ENTRE_API_TOKEN=entre_ai_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
  -- /absolute/path/to/backend/bin/mcp-server
```

確認は `claude mcp list` または Claude Code 内の `/mcp`。

### Claude Desktop

Claude Desktop などJSON設定型のクライアントでは、次の形で stdio server を登録する。

```json
{
  "mcpServers": {
    "entre": {
      "command": "/absolute/path/to/mcp-server",
      "env": {
        "ENTRE_API_BASE_URL": "https://api.example.com",
        "ENTRE_API_TOKEN": "entre_ai_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
      }
    }
  }
}
```

### References

- Codex MCP: <https://developers.openai.com/codex/mcp>
- Claude Code MCP: <https://docs.anthropic.com/en/docs/claude-code/mcp>

開発者向けにDB直結モードも残している。`ENTRE_API_TOKEN` が未設定で `DATABASE_URL` がある場合だけ direct database mode で起動する。

```bash
cd backend
DATABASE_URL=postgres://postgres:postgres@localhost:15432/job_hunting_dev?sslmode=disable \
MCP_USER_EMAIL=you@example.com \
go run ./cmd/mcp-server
```

`MCP_USER_EMAIL` または `MCP_USER_ID` のどちらかが必須。multi-user DBの別ユーザー情報を誤って渡さないため、MCP server は必ず1ユーザーにscopeして起動する。

## Resources

| URI | 内容 |
| --- | --- |
| `entries://list` | 応募先一覧 |
| `entries://{entryId}` | 応募先1件と紐づくTask |
| `tasks://open` | 未完了Task一覧 |
| `inbox://clips` | 保存箱のInbox clip一覧 |
| `es-memos://list` | ES/自己PR/面接ネタ用メモ一覧 |

## Tools

| Tool | 内容 |
| --- | --- |
| `list_entries` | 応募先一覧を返す |
| `get_entry_context` | Entry詳細と紐づくTaskを返す |
| `list_open_tasks` | 未完了Task一覧を返す |
| `list_inbox_clips` | 保存箱のInbox clip一覧を返す |
| `list_es_memos` | ES/自己PR/面接ネタ用メモ一覧を返す |
| `append_es_memo` | ES/自己PR/面接ネタ用メモを保存する |
| `create_task` | Entryに紐づくTaskを作成する |
| `capture_job_email` | 選考メール本文からTask候補などをルールベースで抽出する |

`append_es_memo` と `create_task` は `confirm: true` のときだけDBに保存する。未指定または `false` の場合は `confirmationRequired: true` と保存予定内容だけを返す。

## Distribution note

Local MCP の配布は API bridge mode を標準にする。ユーザーPCにはDB接続情報を置かず、AI連携トークンだけを置く。

非エンジニア向け配布では、Go toolchain を要求しない単一バイナリ、Docker image、またはWebアプリ側で生成するMCPクライアント設定ファイルが必要になる。
