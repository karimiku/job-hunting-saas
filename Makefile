.PHONY: help up down logs dev-fe health \
        test test-be build build-be build-fe \
        lint lint-be lint-fe fmt fmt-be \
        gen gen-api gen-sql \
        install install-fe install-be \
        tidy clean

# ============================================================
# Help
# ============================================================
help: ## このヘルプを表示
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9_-]+:.*?## / {printf "\033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# ============================================================
# 開発サーバ
# ============================================================
up: ## API + DB (docker) とフロント (ローカル pnpm dev) を並行起動。Ctrl+C で両方停止
	@trap 'kill 0' INT TERM EXIT; \
		(cd backend && docker compose up) & \
		(cd frontend && pnpm dev) & \
		wait

up-be: ## バックエンドだけ docker で起動
	cd backend && docker compose up

up-d: ## バックエンドだけ docker バックグラウンド起動
	cd backend && docker compose up -d

down: ## docker compose down（ローカル pnpm dev は対象外）
	cd backend && docker compose down

down-v: ## volumes 含めて全消し（DB も消える）
	cd backend && docker compose down -v

logs: ## app / db のログを追う（SERVICE=app で絞れる）
	cd backend && docker compose logs -f $(SERVICE)

dev-fe: ## フロント dev サーバ単独

health: ## API /health を叩く
	@curl -sf http://localhost:8080/health && echo " OK" || echo "API unreachable"

# ============================================================
# テスト・ビルド
# ============================================================
test: test-be ## すべてのテスト（今はバックのみ）

test-be: ## Go ユニットテスト
	cd backend && go test ./...

build: build-be build-fe ## 両方ビルド

build-be: ## Go ビルド確認
	cd backend && go build ./...

build-fe: ## Next.js 本番ビルド
	cd frontend && pnpm build

# ============================================================
# Lint / Format
# ============================================================
lint: lint-be lint-fe ## 両方 lint

lint-be: ## go vet
	cd backend && go vet ./...

lint-fe: ## pnpm lint (ESLint)
	cd frontend && pnpm lint

fmt: fmt-be ## フォーマット（バック）

fmt-be: ## gofmt -w
	cd backend && gofmt -w .

# ============================================================
# コード生成
# ============================================================
gen: gen-api gen-sql ## OpenAPI + sqlc 両方

gen-api: ## openapi.yaml → server.gen.go
	cd backend && oapi-codegen -c api/oapi-codegen.yaml api/openapi.yaml > internal/gen/openapi/server.gen.go

gen-sql: ## sql/queries/*.sql → sqlc
	cd backend && sqlc generate

# ============================================================
# セットアップ
# ============================================================
install: install-fe install-be ## 依存インストール

install-fe: ## pnpm install
	cd frontend && pnpm install

install-be: ## go mod download
	cd backend && go mod download

tidy: ## go mod tidy
	cd backend && go mod tidy

clean: ## ビルド成果物を消す
	cd backend && rm -rf tmp server
	cd frontend && rm -rf .next
