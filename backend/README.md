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

### 認証 / Chrome拡張向け API 設定

認証ありで起動する場合は `DATABASE_URL` に加えて `SUPABASE_AUTH_ISSUER` または legacy rollback 用の `FIREBASE_PROJECT_ID` が必要。`DATABASE_URL` を設定すると PostgreSQL、未設定だと InMemory リポジトリで起動する。

| 変数 | 用途 | 既定 / 例 |
| --- | --- | --- |
| `PORT` | API ポート | `8080` |
| `DATABASE_URL` | 設定で PostgreSQL、未設定で InMemory | `postgres://postgres:postgres@localhost:15432/job_hunting_dev` |
| `SUPABASE_AUTH_ISSUER` | Supabase Auth issuer。Bearer JWT 検証に使用 | `https://<project-ref>.supabase.co/auth/v1` |
| `SUPABASE_JWKS_URL` | Supabase JWKS URL。未設定なら issuer から導出 | `https://<project-ref>.supabase.co/auth/v1/.well-known/jwks.json` |
| `SUPABASE_JWT_AUDIENCE` | Supabase access token audience | `authenticated` |
| `PGPOOL_MAX_CONNS` | container 1台あたりの pgxpool 最大接続数。autoscale 前提では小さく保つ | `4` |
| `PGPOOL_MIN_CONNS` | 起動時に維持する最小接続数 | `0` |
| `PGPOOL_MAX_CONN_IDLE_TIME` | idle 接続を閉じるまでの時間 | `60s` |
| `PGPOOL_MAX_CONN_LIFETIME` | 接続の最大寿命 | `30m` |
| `PGPOOL_HEALTH_CHECK_PERIOD` | pgxpool health check 間隔 | `30s` |
| `PGX_DEFAULT_QUERY_EXEC_MODE` | Supabase transaction pooler 用に statement cache を避ける | `exec` |
| `PGAPPNAME` | Supabase 側の接続監視に出す application_name | `job-hunting-saas-api` |
| `FIREBASE_PROJECT_ID` | rollback 期間だけ使う Firebase プロジェクト ID | 未設定 |
| `FIREBASE_CREDENTIALS_FILE` | rollback 期間だけ使う service account JSON のパス | 未設定 |
| `CORS_ALLOWED_ORIGINS` | カンマ区切りの許可 origin（Web + 拡張）。Vercel preview は `https://*.vercel.app` を使える | `http://localhost:3000,https://*.vercel.app,chrome-extension://<extension-id>` |
| `COOKIE_SECURE` | 本番 HTTPS は `true`、localhost は `false` | `false` |
| `COOKIE_SAME_SITE` | legacy session cookie / dev auth cookie 用。`lax` / `strict` / `none` | `lax` |
| `DEV_AUTH_ENABLED` | `true` で localhost 専用の `/dev/session` を有効化。Codex内ブラウザ等のローカルUI確認用 | 未設定 |
| `DEV_AUTH_SECRET` | 開発用 session cookie 署名鍵。未設定なら起動ごとに自動生成 | 未設定 |
| `RATE_LIMIT_GLOBAL_REQUESTS_PER_MINUTE` | IPごとの全体レート制限。`0` で無効 | `30` |
| `RATE_LIMIT_AUTH_REQUESTS_PER_MINUTE` | IPごとの `/auth/session` レート制限。`0` で無効 | `5` |
| `RATE_LIMIT_AUTHENTICATED_REQUESTS_PER_MINUTE` | ログインユーザーごとのAPIレート制限。`0` で無効 | `60` |
| `ALLOW_INSECURE_NO_AUTH` | `true` で認証なし起動（ローカル検証用） | 未設定 |

> ⚠️ service account JSON と `.env` はコミットしない（`.gitignore` 済みであることを確認）。secret 値そのものを README に書かないこと。

Chrome拡張から認証付きで API を呼ぶ本番 HTTPS 環境では、拡張の origin を `CORS_ALLOWED_ORIGINS` に追加する。Web frontend は Supabase access token を `Authorization: Bearer` として送る。拡張側の Supabase token 対応は migration issue #228 で別途確認する。

```bash
SUPABASE_AUTH_ISSUER=https://<project-ref>.supabase.co/auth/v1
SUPABASE_JWKS_URL=https://<project-ref>.supabase.co/auth/v1/.well-known/jwks.json
DATABASE_URL=postgres://postgres.<project-ref>:<password>@aws-<region>.pooler.supabase.com:6543/postgres?sslmode=require
PGPOOL_MAX_CONNS=4
PGX_DEFAULT_QUERY_EXEC_MODE=exec
CORS_ALLOWED_ORIGINS=https://your-web-app.example,chrome-extension://<extension-id>
COOKIE_SECURE=true
COOKIE_SAME_SITE=none
```

- **local (HTTP)**: `COOKIE_SECURE=false` / `COOKIE_SAME_SITE=lax`
- **本番 (HTTPS)**: `COOKIE_SECURE=true` / `COOKIE_SAME_SITE=none`。`CORS_ALLOWED_ORIGINS` は CORS と CSRF Origin/Referer 検証の allowlist を兼ねる。

ローカルでGoogleログインやパスキーを使わずUI確認したい場合は、backendに `DEV_AUTH_ENABLED=true` を設定して起動し、frontendの `/dev/login` を開く。`DEV_AUTH_ENABLED=true` は production 系の環境変数 (`APP_ENV=production` 等) では起動時に拒否される。

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
