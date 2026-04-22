# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

就活管理SaaS。Go 1.25 + Clean Architecture のバックエンドAPI。OpenAPI 3.0.3 を Single Source of Truth とし、oapi-codegen でハンドラ層を自動生成。データアクセスは sqlc で型安全なSQLコード生成。

## Commands

```bash
# 開発環境（Go + PostgreSQL 16、Air によるホットリロード）
docker compose up

# テスト（InMemoryリポジトリ使用、DB不要）
cd backend && go test ./...

# 単一パッケージのテスト
cd backend && go test ./internal/domain/entity/...
cd backend && go test ./internal/usecase/entry/...

# ビルド確認
cd backend && go build ./...

# OpenAPI からハンドラコード再生成（openapi.yaml 変更後に必須）
cd backend && oapi-codegen -c api/oapi-codegen.yaml api/openapi.yaml > internal/gen/openapi/server.gen.go

# sqlc コード再生成（sql/queries/*.sql 変更後に必須）
cd backend && sqlc generate

# ヘルスチェック
curl http://localhost:8080/health
```

## Architecture

4層 Clean Architecture。依存方向は Handler → UseCase → Domain ← Infra。

- **Domain層** (`internal/domain/`): Entity, Value Object, Repository インターフェース。フレームワーク・DB に一切依存しない。Value Object がドメインルールを強制する（バリデーションはここ）。
- **UseCase層** (`internal/usecase/`): ビジネスロジックのオーケストレーション。Repository インターフェースを DI で受け取る。
- **Handler層** (`internal/handler/`): oapi-codegen が生成する `ServerInterface` を実装。HTTP ↔ UseCase 入出力の変換のみ。
- **Infra層** (`internal/infra/`): Repository 実装。`inmemory/`（sync.Map、開発・テスト用）と `postgres/`（sqlc + pgx/v5）。

`cmd/server/main.go` で DI を構成。`DATABASE_URL` の有無で InMemory / PostgreSQL を切り替え。

## Code Generation Rules

2つのコード生成パイプラインがある。生成コードは手動編集しない。

1. **OpenAPI → Go**: `api/openapi.yaml` → oapi-codegen → `internal/gen/openapi/server.gen.go`
2. **SQL → Go**: `sql/queries/*.sql` + `sql/schema.sql` → sqlc → `internal/infra/postgres/sqlc/`

API エンドポイントの追加・変更は必ず `openapi.yaml` から始める。DB クエリの追加・変更は `sql/queries/` の `.sql` ファイルから始める。

## Database

PostgreSQL 16。開発環境は docker-compose で起動（ポート 15432）。

- スキーマ: `sql/schema.sql`（8テーブル、ENUM型: entry_status, stage_kind, task_type, task_status, auth_provider）
- 初期化スクリプト: `docker/initdb/01_create_test_db.sql`

## Testing Conventions

- テストは InMemory リポジトリを DI して DB なしで実行
- UseCase テストではモックリポジトリを使用
- Value Object のテスト: 境界値・不正値のバリデーションを網羅
- テストファイルは対象と同じパッケージに配置（`_test.go`）

## Key Design Decisions

- Value Object は `New___()` コンストラクタで生成し、不正値はコンパイル時 or 生成時に弾く
- Entity の ID は UUID v4（`google/uuid`）
- エラーはドメイン固有エラー（`domain/repository/` で定義）を使い、Handler 層で HTTP ステータスに変換
- 詳細な技術選定理由は `docs/why-reasons.md` に記載
