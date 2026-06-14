# Job Hunting SaaS

就活活動を一元管理するWebアプリケーション。
複数の就活サイトに分散する応募・選考・タスクの情報を横断的に集約し、締切漏れや情報散逸を防ぐ。

## 特徴

- **スキーマ駆動開発** — OpenAPI定義からハンドラーインターフェースを自動生成（oapi-codegen）
- **Clean Architecture** — domain / usecase / handler / infra の4層構造。ビジネスロジックがフレームワークやDBに依存しない設計
- **型安全なドメインモデル** — 値オブジェクト（EntryStatus, Stage, Source等）によるバリデーションをドメイン層で保証
- **TDD** — インメモリRepository注入によりDBなしでユースケースのテストが実行可能

## 技術スタック

| カテゴリ | 技術 |
|---------|------|
| 言語 | Go 1.26 |
| HTTPルーター | Chi v5 |
| APIスキーマ | OpenAPI 3.0.3 + oapi-codegen |
| データベース | PostgreSQL 16 |
| コンテナ | Docker / Docker Compose |
| ホットリロード | Air |
| ID生成 | UUID v4（google/uuid） |

技術選定の詳細な理由は [docs/why-reasons.md](../docs/why-reasons.md) を参照。

## プロジェクト構成

```
.
├── api/
│   ├── openapi.yaml          # API定義（Single Source of Truth）
│   └── oapi-codegen.yaml     # コード生成設定
├── cmd/server/
│   └── main.go               # エントリーポイント・DI配線
├── internal/
│   ├── domain/
│   │   ├── entity/            # エンティティ（Company, Entry, Task, StageHistory, User）
│   │   ├── value/             # 値オブジェクト（EntryStatus, Stage, Source, Route等）
│   │   └── repository/        # Repositoryインターフェース
│   ├── usecase/               # ユースケース（company, entry, task, stage_history, user等）
│   ├── handler/               # HTTPハンドラー（oapi-codegen ServerInterface実装）
│   ├── middleware/            # 認証ミドルウェア
│   ├── infra/inmemory/        # インメモリRepository実装（開発用）
│   └── gen/openapi/           # oapi-codegen自動生成コード
├── docs/
│   ├── requirements.md        # 要件定義書
│   └── why-reasons.md         # 技術選定理由
├── docker-compose.yml
├── Dockerfile                 # 本番用（マルチステージビルド → distroless）
├── Dockerfile.dev             # 開発用（Airによるホットリロード）
└── .air.toml                  # Air設定
```

## セットアップ

### 前提条件

- Docker / Docker Compose

### 起動

```bash
# .envファイルを作成（デフォルト値で動作するが、必要に応じて変更）
cat <<EOF > .env
PORT=8080
POSTGRES_PORT=15432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=job_hunting_dev
EOF

# 起動（Go + PostgreSQL）
docker compose up
```

ホットリロードが有効なため、Goファイルを編集すると自動で再ビルドされる。

### ヘルスチェック

```bash
curl http://localhost:8080/health
```

### MCP server

Claude Desktop / Codex / Gemini CLI などのMCPクライアントから、就活データを読み書きするためのstdio MCP serverを提供する。
設計・resources/tools・配布上の注意は [docs/mcp-server.md](../docs/mcp-server.md) を参照。

通常利用は Webアプリの「アカウント」画面で発行するAI連携トークンを使う。CodexではNode wrapperをstdio serverとして登録する。

```bash
ENTRE_API_BASE_URL=http://localhost:8080 \
ENTRE_API_TOKEN=entre_ai_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
node cmd/mcp-remote/entre-mcp.mjs
```

Go版MCP serverを使う場合は単一バイナリを作る。

```bash
make -C .. build-mcp-server
ENTRE_API_BASE_URL=http://localhost:8080 \
ENTRE_API_TOKEN=entre_ai_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
./bin/mcp-server
```

リポジトリルートからは次でも起動できる。

```bash
ENTRE_API_BASE_URL=http://localhost:8080 \
ENTRE_API_TOKEN=entre_ai_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
make mcp-server
```

開発者向けにDB直結モードも残している。`ENTRE_API_TOKEN` が未設定の場合だけ `DATABASE_URL` と `MCP_USER_EMAIL` または `MCP_USER_ID` で対象ユーザーを指定する。multi-user DBを安全に扱うため、DB直結モードではユーザー指定が必須。

`append_es_memo` と `create_task` は `confirm: true` を渡したときだけDBへ保存する。`list_inbox_clips` と `list_es_memos` で保存箱・ESメモを参照できる。`capture_job_email` はメール本文をルールベースで構造化し、LLM APIは呼ばない。

### 認証 / Chrome拡張向け Cookie 設定

認証ありで起動する場合は `DATABASE_URL` に加えて Firebase 設定が必要。`DATABASE_URL` を設定すると PostgreSQL、未設定だと InMemory リポジトリで起動する。

| 変数 | 用途 | 既定 / 例 |
| --- | --- | --- |
| `PORT` | API ポート | `8080` |
| `DATABASE_URL` | 設定で PostgreSQL、未設定で InMemory | `postgres://postgres:postgres@localhost:15432/job_hunting_dev` |
| `FIREBASE_PROJECT_ID` | Firebase プロジェクト ID（認証に必須） | `your-project-id` |
| `FIREBASE_CREDENTIALS_FILE` | service account JSON のパス | `./secrets/service-account.json` |
| `CORS_ALLOWED_ORIGINS` | カンマ区切りの許可 origin（Web + 拡張） | `http://localhost:3000,chrome-extension://<extension-id>` |
| `COOKIE_SECURE` | 本番 HTTPS は `true`、localhost は `false` | `false` |
| `COOKIE_SAME_SITE` | `lax` / `strict` / `none`。拡張から Cookie を送る本番は `none`。`none` は Origin/Referer CSRF 検証とセットで使う | `lax` |
| `ALLOW_INSECURE_NO_AUTH` | `true` で認証なし起動（ローカル検証用） | 未設定 |

> ⚠️ service account JSON と `.env` はコミットしない（`.gitignore` 済みであることを確認）。secret 値そのものを README に書かないこと。

Chrome拡張から認証付きで API を呼ぶ本番 HTTPS 環境では、拡張の origin を `CORS_ALLOWED_ORIGINS` に追加し、Cookie を cross-site fetch に送れるようにする:

```bash
FIREBASE_PROJECT_ID=your-project-id
CORS_ALLOWED_ORIGINS=https://your-web-app.example,chrome-extension://<extension-id>
COOKIE_SECURE=true
COOKIE_SAME_SITE=none
```

- **local (HTTP)**: `COOKIE_SECURE=false` / `COOKIE_SAME_SITE=lax`
- **本番 (HTTPS)**: `COOKIE_SECURE=true` / `COOKIE_SAME_SITE=none`。`CORS_ALLOWED_ORIGINS` は CORS と CSRF Origin/Referer 検証の allowlist を兼ねる。

backend / frontend / 拡張をまたぐ通し手順と拡張 ID の確認方法は、ルート [README.md](../README.md) の「β環境セットアップ」を参照。

### テスト実行

```bash
go test ./...
```

### コード生成

OpenAPI 定義を変更した場合:

```bash
oapi-codegen -c api/oapi-codegen.yaml api/openapi.yaml
```

## API エンドポイント

APIは OpenAPI 3.0.3 で定義されている。詳細は [api/openapi.yaml](api/openapi.yaml) を参照。

| メソッド | パス | 概要 |
|---------|------|------|
| GET | `/health` | ヘルスチェック |
| POST | `/api/v1/companies` | 企業を登録 |
| GET | `/api/v1/companies` | 企業一覧を取得 |
| GET | `/api/v1/companies/{companyId}` | 企業を取得 |
| PATCH | `/api/v1/companies/{companyId}` | 企業を更新 |
| DELETE | `/api/v1/companies/{companyId}` | 企業を削除 |
| POST | `/api/v1/entries` | エントリーを登録 |
| GET | `/api/v1/entries` | エントリー一覧（status/stageKind/sourceでフィルタ可） |
| GET | `/api/v1/entries/{entryId}` | エントリーを取得 |
| PATCH | `/api/v1/entries/{entryId}` | エントリーを更新 |
| DELETE | `/api/v1/entries/{entryId}` | エントリーを削除 |
| POST | `/api/v1/entries/{entryId}/tasks` | タスクを登録 |
| GET | `/api/v1/entries/{entryId}/tasks` | タスク一覧を取得 |
| GET | `/api/v1/tasks/{taskId}` | タスクを取得 |
| PATCH | `/api/v1/tasks/{taskId}` | タスクを更新 |
| DELETE | `/api/v1/tasks/{taskId}` | タスクを削除 |
| POST | `/api/v1/entries/{entryId}/stage-histories` | 選考フェーズ履歴を追加 |
| GET | `/api/v1/entries/{entryId}/stage-histories` | 選考フェーズ履歴一覧を取得 |

## アーキテクチャ

```
HTTP Request
  → Middleware（認証）
    → Handler（リクエスト変換・レスポンス変換）
      → UseCase（ビジネスロジック）
        → Domain Entity / Value Object（ドメインルール）
        → Repository Interface ← InMemory / PostgreSQL 実装を注入
```

- **外側は自動生成**: openapi.yaml → oapi-codegen（Handler層）
- **中心は手動 + TDD**: Domain層・UseCase層はビジネスロジックを手動で実装しテストで品質を担保

## 開発ステータス

### 実装済み

- [x] Company CRUD
- [x] Entry CRUD（フィルタ付き一覧）
- [x] Task CRUD（deadline / schedule の種別管理）
- [x] StageHistory 作成・一覧（イミュータブルな選考履歴）
- [x] CompanyAlias ユースケース層
- [x] 値オブジェクトによるドメインバリデーション
- [x] OpenAPI スキーマ駆動のハンドラー生成
- [x] Docker Compose 開発環境
- [x] 認証（Firebase セッション Cookie / Google ログイン）
- [x] PostgreSQL Repository 実装（sqlc + pgx）
- [x] Inbox Clip（作成 / 一覧 / 削除、PostgreSQL 永続化）

### 未実装（正式リリース以降）

- [ ] メール通知（#38）
- [ ] CSV エクスポート（#37）
- [ ] アカウント削除 / データ全削除（#35）
- [ ] ダッシュボード用集約API（β では既存 list API から算出）

## ドキュメント

- [要件定義書](../docs/requirements.md) — プロジェクトのビジョン・機能要件・データモデルの全体像
- [技術選定理由](../docs/why-reasons.md) — Go / Clean Architecture / PostgreSQL / sqlc / Chi / oapi-codegen / TDD の選定根拠
