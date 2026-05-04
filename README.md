# Job Hunting SaaS

就活活動を一元管理する Web サービス。複数の就活サイトに分散する応募・選考・タスクの情報を横断的に集約し、締切漏れや情報散逸を防ぐ。

## モノレポ構成

```
.
├── backend/      # Go バックエンド（API サーバー / DB / OpenAPI）
├── frontend/     # フロントエンド（予定）
└── docs/         # 要件定義 / 技術選定理由 など全体ドキュメント
```

## 各プロジェクト

### backend/
Go 1.25 / Chi / PostgreSQL / sqlc / oapi-codegen による Clean Architecture バックエンド。

- セットアップ・コマンド: [backend/README.md](backend/README.md)
- API 定義: [backend/api/openapi.yaml](backend/api/openapi.yaml)

### frontend/
未着手（予定）。バックエンドが OpenAPI で公開する型を共有して構築する。

## ドキュメント

- [要件定義書](docs/requirements.md) — ビジョン・機能要件・データモデルの全体像
- [技術選定理由](docs/why-reasons.md) — Go / Clean Architecture / PostgreSQL / sqlc / Chi / oapi-codegen / TDD の選定根拠

## 開発環境

Nix を使う場合は、リポジトリルートで dev shell に入る。

```bash
nix develop
```

direnv を使う場合は初回だけ許可する。以後は `cd` するだけで `.envrc` の `use flake` により同じ環境が有効になる。

```bash
direnv allow
```

Nix dev shell は Go / Node.js / pnpm / sqlc / oapi-codegen / Docker Compose などの開発ツールを揃える。Docker daemon は別途 Docker Desktop / OrbStack / Colima などで起動しておく必要がある。Nix を使わない場合も、従来通り各ツールをローカルに用意すれば `make` や `docker compose` の手順は使える。

## クイックスタート

```bash
cd backend
docker compose up
curl http://localhost:8080/health
```
