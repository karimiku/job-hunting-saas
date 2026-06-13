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
| cmd | `backend/cmd/mcp-server` | 環境変数読み込み、DB接続、DI |
| Handler | `backend/internal/handler/mcp` | MCP JSON-RPC over stdio、resources/tools の入出力変換 |
| UseCase | `backend/internal/usecase/mcp`, `es_memo`, `job_email`, `task` | user scoped な操作、保存前preview、メール抽出 |
| Domain | `backend/internal/domain/entity`, `value`, `repository` | ESMemo entity/value object、Repository interface |
| Infra | `backend/internal/infra/postgres` | sqlc query / repository 実装 |

## Token

AI連携トークンは `/profile` から発行する。画面ではtoken全文だけを一度表示し、クライアント別の設定コマンドは出さない。
tokenはユーザーに紐づき、DBには平文ではなくSHA-256 hashだけを保存する。漏れた・失くした場合は同じ画面で失効して作り直す。

CLIで生成・登録する場合は次を使う。既存トークンを渡さずに実行すると、新しい `entre_ai_...` トークンを生成して一度だけ表示する。

```bash
cd backend
AI_ACCESS_TOKEN='entre_ai_...' \
DATABASE_URL=postgres://postgres:postgres@localhost:15432/job_hunting_dev?sslmode=disable \
go run ./cmd/ai-token -email you@example.com -name "AI連携"
```

## Run: local DB MCP server

開発中は backend を起動し、DB migration が適用された状態で MCP server を stdio 起動する。
ローカルDBを直接読む場合は、MCP client 側の環境変数に `MCP_API_KEY` として渡す。

```bash
cd backend
DATABASE_URL=postgres://postgres:postgres@localhost:15432/job_hunting_dev?sslmode=disable \
MCP_API_KEY='entre_ai_...' \
go run ./cmd/mcp-server
```

root からは次でも登録・起動できる。

```bash
AI_ACCESS_TOKEN='entre_ai_...' \
DATABASE_URL=postgres://postgres:postgres@localhost:15432/job_hunting_dev?sslmode=disable \
make ai-token ARGS='-email you@example.com -name AI連携'

DATABASE_URL=postgres://postgres:postgres@localhost:15432/job_hunting_dev?sslmode=disable \
MCP_API_KEY='entre_ai_...' \
make mcp-server
```

開発用には、ユーザーを明示して起動する方法も使える。

```bash
cd backend
DATABASE_URL=postgres://postgres:postgres@localhost:15432/job_hunting_dev?sslmode=disable \
MCP_USER_EMAIL=you@example.com \
go run ./cmd/mcp-server
```

root からは次でも起動できる。

```bash
DATABASE_URL=postgres://postgres:postgres@localhost:15432/job_hunting_dev?sslmode=disable \
MCP_USER_EMAIL=you@example.com \
make mcp-server
```

`MCP_API_KEY`、`MCP_USER_EMAIL`、`MCP_USER_ID` のいずれかが必須。`MCP_API_KEY` がある場合はこれを優先する。multi-user DBの別ユーザー情報を誤って渡さないため、MCP server は必ず1ユーザーにscopeして起動する。

## Run: remote API wrapper

本番APIへ接続するlocal MCP wrapperは `backend/cmd/mcp-remote/entre-mcp.mjs` を使う。
MCP clientからはこのNode scriptをstdio serverとして起動し、環境変数でAPI base URLとtokenを渡す。

| 変数 | 内容 |
| --- | --- |
| `ENTRE_API_BASE_URL` | Entré APIのbase URL |
| `ENTRE_API_TOKEN` | `/profile` 等で発行した `entre_ai_...` token |

remote wrapperはCLIに内部UUIDやtimestampを出さない。`list_entries` は `entry-1` のような一時refを返し、詳細取得やTask作成もそのrefで扱う。

## Resources

| URI | 内容 |
| --- | --- |
| `entries://list` | 応募先一覧 |
| `entries://{entryId}` | 応募先1件と紐づくTask |
| `tasks://open` | 未完了Task一覧 |
| `inbox://clips` | 保存箱のInbox clip一覧 |

## Tools

| Tool | 内容 |
| --- | --- |
| `list_entries` | 応募先一覧を返す |
| `get_entry_context` | Entry詳細と紐づくTaskを返す |
| `list_open_tasks` | 未完了Task一覧を返す |
| `append_es_memo` | ES/自己PR/面接ネタ用メモを保存する |
| `create_task` | Entryに紐づくTaskを作成する |
| `capture_job_email` | 選考メール本文からTask候補などをルールベースで抽出する |

`append_es_memo` と `create_task` は `confirm: true` のときだけDBに保存する。未指定または `false` の場合は `confirmationRequired: true` と保存予定内容だけを返す。

## Distribution note

現状は開発者・セルフホスト向けに `go run ./cmd/mcp-server` を提供する段階。非エンジニア向け配布では、Go toolchain を要求しない単一バイナリ、Docker image、またはWebアプリ側で生成するMCPクライアント設定ファイルが必要になる。
