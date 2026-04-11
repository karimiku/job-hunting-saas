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

## クイックスタート

```bash
cd backend
docker compose up
curl http://localhost:8080/health
```
