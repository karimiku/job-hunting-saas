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

### UseCase の命名規約

UseCase は「1操作 = 1構造体 + 単一メソッド `Execute`」で統一する（command オブジェクト風）。
現状 `internal/usecase/` 配下の全 25 メソッドが例外なくこの形になっており、その de-facto を標準として明文化したもの。

- **構造体名 = 操作の動詞**: `Create` / `Get` / `List` / `Update` / `Delete`。ドメイン固有の操作は意味のある動詞を使う（例: `user.Authenticate`）。パッケージ名（`entry`, `task` など）が文脈を与えるため、構造体名に名詞を重ねない（`entry.Create` であって `entry.CreateEntry` ではない）。
- **コンストラクタ**: `New<構造体名>`（例: `NewCreate`、`NewAuthenticate`）。依存（Repository インターフェース）を引数で受け取り DI する。
- **実行メソッドは必ず `Execute`**: シグネチャは `func (uc *<構造体名>) Execute(ctx context.Context, input <構造体名>Input) (*<構造体名>Output, error)`。出力を持たない操作（Delete 等）は `(..., error)` のみ。`Handle` / `Run` / `Do` など別名は使わない。
- **入出力型**: `<構造体名>Input` / `<構造体名>Output` を同パッケージに定義（例: `CreateInput`, `CreateOutput`）。
- 別パッケージから参照する際はパッケージエイリアスで衝突を避ける（例: `companyuc.CreateInput`、`useruc.AuthenticateInput`）。

採用理由: 具体名（`CreateEntry` 等）はパッケージ名と冗長になり、Handler から呼ぶ際の呼び出し点が `uc.Execute(...)` に揃って読みやすい。「動詞構造体 + `Execute`」は1操作1責務を型レベルで強制でき、現状のコードベースが既に 100% この形のため、これを唯一の規約とする。

将来そろえる候補（本リポジトリには現状 outlier なし）: 1構造体に複数操作をまとめたくなった場合でも、新たな操作は別構造体に切り出して `Execute` を1つに保つ。
