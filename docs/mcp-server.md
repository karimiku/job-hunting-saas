# MCP server

job-hunting-saas は stdio MCP server を提供し、Claude Desktop / Codex / Gemini CLI などのMCPクライアントから就活データを参照・更新できる。
アプリ本体はLLM APIを呼ばない。料金・モデル選択・APIキー管理はユーザーが使うMCPクライアント側に寄せる。

## 想定ユースケース

- Chrome拡張で保存した求人ページをAIが `inbox://clips` から読み取り、Entry化や調査タスク作成の候補を作る
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

## Run

開発中は backend を起動し、DB migration が適用された状態で MCP server を stdio 起動する。

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

`MCP_USER_EMAIL` または `MCP_USER_ID` のどちらかが必須。multi-user DBの別ユーザー情報を誤って渡さないため、MCP server は必ず1ユーザーにscopeして起動する。

## Resources

| URI | 内容 |
| --- | --- |
| `entries://list` | 応募先一覧 |
| `entries://{entryId}` | 応募先1件と紐づくTask |
| `tasks://open` | 未完了Task一覧 |
| `inbox://clips` | Chrome拡張などで保存したInbox clip一覧 |

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
