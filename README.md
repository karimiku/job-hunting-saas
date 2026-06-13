# Job Hunting SaaS

就活活動を Entry 単位で一元管理する Web サービス。保存箱に残した求人を Entry に変換し、カンバンで選考状況を見ながら、Entry に紐づくタスクを処理する。

## モノレポ構成

```
.
├── backend/           # Go バックエンド（API / DB / OpenAPI / 認証 / MCP server）
├── frontend/          # Next.js 16 フロントエンド（App Router / SSR / Server Actions）
├── chrome-extension/  # 求人ページを保存箱に送る任意の Chrome 拡張（Vite + React）
└── docs/              # 要件定義 / 技術選定理由 など全体ドキュメント
```

## 各プロジェクト

### backend/
Go 1.26 / Chi / PostgreSQL / sqlc / oapi-codegen による Clean Architecture バックエンド。Firebase セッション Cookie による認証付き。

- セットアップ・コマンド: [backend/README.md](backend/README.md)
- API 定義: [backend/api/openapi.yaml](backend/api/openapi.yaml)

### frontend/
Next.js 16（App Router / SSR / Server Actions）。Firebase で Google ログインし、ホーム / Entry / カンバン / タスク / 保存箱を提供する。

- セットアップ・コマンド: [frontend/README.md](frontend/README.md)

### chrome-extension/
求人ページを Web の保存箱に送る任意の Chrome 拡張（Vite + React）。現在のコア体験は Web 単体でも成立する。

- セットアップ・コマンド: [chrome-extension/README.md](chrome-extension/README.md)

## ドキュメント

- [要件定義書](docs/requirements.md) — ビジョン・機能要件・データモデルの全体像
- [技術選定理由](docs/why-reasons.md) — Go / Clean Architecture / PostgreSQL / sqlc / Chi / oapi-codegen / TDD の選定根拠
- [MCP server](docs/mcp-server.md) — AIクライアントから就活データを扱うstdio MCP server
- [βリリースプラン](docs/plans/beta-release-plan.md) — βのスコープ・依存関係・リリースゲート
- [β E2E サインオフ](docs/beta-e2e-signoff.md) — 保存箱→Entry→カンバン/タスク管理フローの動作確認手順

## 開発ステータス（β）

βの完成条件「保存した求人を、ログインユーザーが Web の保存箱から Entry 化し、カンバンとタスクで管理できる」を満たした状態。Chrome 拡張は保存を速くする任意の入口であり、コア体験は Web の Entry / カンバン / タスク / 保存箱に集約する。

LP も同じ方針に合わせ、訴求は「保存箱 → Entry → カンバン」の一本道に絞る。表示画像は実サービス画面の Shots.so 書き出しを `frontend/public/marketing/` に置き、スマホ・PC の見た目が実装とかけ離れないように管理する。

実装済み:

- Google ログイン（Firebase セッション Cookie）
- 保存箱 clip の一覧 / Entry 化 / 削除、重複 URL の最小対応
- Entry 一覧 / カンバン / 詳細（会社名表示）
- タスク画面（実 API データの表示・状態更新）
- ホーム / サイドバーの実データ表示
- Chrome 拡張から保存箱 API への保存（任意の入力補助、未ログイン・保存失敗時の回復導線つき）
- ローカル MCP server（`entre_ai_...` トークンまたはユーザー指定で起動。Entry / Inbox / Task 参照、Task作成、ESメモ蓄積、メール本文のルールベース抽出）
- backend / frontend / chrome-extension の CI ゲート

正式リリースで対応（β対象外）:

- メール通知 / CSV エクスポート / 退会（全データ削除）
- Chrome Web Store 公開（privacy policy・ストア素材）
- Google Calendar / Gmail 連携、ロードマップ、ES自動生成などの拡張機能
- 会社+Entry 同時作成のトランザクション化（[#62](https://github.com/karimiku/job-hunting-saas/issues/62)）
- revoked セッションの即時無効化（[#56](https://github.com/karimiku/job-hunting-saas/issues/56)）

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

## β環境セットアップ（保存箱 → Entry → カンバン/タスク）

βの完成条件は「保存した求人が、ログインユーザーの DB に入り、Web の保存箱から Entry 化してカンバン/タスクで管理できる」こと。Web のコア導線は backend / frontend だけで確認できる。Chrome 拡張まで含めて実ログイン状態で通す場合は、Firebase・CORS・Cookie・拡張 ID の設定も揃える。

> ⚠️ secret 値（Firebase API キー、service account JSON）はコミットしない。`.env` / `.env.local` / service account JSON はいずれも `.gitignore` 済みであることを確認すること。

### 1. backend (`backend/.env`)

詳細は [backend/README.md](backend/README.md) の「認証 / Chrome拡張向け Cookie 設定」。最低限:

| 変数 | 用途 | ローカル例 |
| --- | --- | --- |
| `DATABASE_URL` | 設定すると PostgreSQL、未設定だと InMemory | `postgres://postgres:postgres@localhost:15432/job_hunting_dev` |
| `FIREBASE_PROJECT_ID` | Firebase プロジェクト ID | `your-project-id` |
| `FIREBASE_CREDENTIALS_FILE` | service account JSON のパス（gitignore 対象に置く） | `./secrets/service-account.json` |
| `CORS_ALLOWED_ORIGINS` | カンマ区切りの許可 origin（Web + 任意の拡張） | `http://localhost:3000,chrome-extension://<extension-id>` |
| `COOKIE_SECURE` | 本番 HTTPS は `true`、localhost は未設定/`false` | `false` |
| `COOKIE_SAME_SITE` | 拡張から Cookie を送る本番は `none`、localhost は `lax` | `lax` |

### 2. frontend (`frontend/.env.local`)

Firebase Web SDK 設定と API ベース URL。詳細は [frontend/README.md](frontend/README.md)。

```
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
NEXT_PUBLIC_FIREBASE_API_KEY=...
NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN=...
NEXT_PUBLIC_FIREBASE_PROJECT_ID=...
NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET=...
NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID=...
NEXT_PUBLIC_FIREBASE_APP_ID=...
```

### 3. chrome-extension（任意: 拡張 ID の確認）

接続先は `VITE_API_BASE_URL` / `VITE_WEB_BASE_URL` で切り替え（既定は localhost）。詳細は [chrome-extension/README.md](chrome-extension/README.md)。

1. `cd chrome-extension && pnpm install && pnpm build` で `dist/` を生成
2. `chrome://extensions` →「デベロッパーモード」ON →「パッケージ化されていない拡張機能を読み込む」で `chrome-extension/dist/` を選択
3. 読み込んだ拡張カードに表示される **ID**（例: `abcdefghijklmnopabcdefghijklmnop`）を控える
4. その ID を backend の `CORS_ALLOWED_ORIGINS` に `chrome-extension://<id>` として追加して再起動

### 起動順

1. backend: `cd backend && docker compose up`（PostgreSQL + API, :8080）
2. frontend: `cd frontend && pnpm install && pnpm dev`（:3000）
3. `http://localhost:3000/login` で Google ログイン
4. 任意で chrome-extension: `pnpm build` → `chrome://extensions` で読み込み（再ビルドのたびにリロード）

### 動作確認

| URL | 期待 |
| --- | --- |
| `http://localhost:8080/health` | `ok` |
| `http://localhost:8080/auth/me`（ログイン Cookie 付き） | 現在のユーザー JSON / 未ログインなら 401 |
| `http://localhost:8080/api/v1/inbox/clips`（Cookie 付き） | 自分の保存箱 clip 一覧 |
| `http://localhost:3000/inbox` | 保存した clip が表示される |

保存箱に clip を作成後、backend を再起動しても PostgreSQL 上の clip が残っていれば永続化 OK。

### local と本番 (HTTPS) の Cookie 差

- **local (HTTP)**: `COOKIE_SECURE` 未設定(false) / `COOKIE_SAME_SITE=lax`。
- **本番 (HTTPS)**: `COOKIE_SECURE=true` / `COOKIE_SAME_SITE=none`。これで拡張 origin からも credentials 付き fetch に Cookie が乗る。`CORS_ALLOWED_ORIGINS` に Web と拡張両方の origin を入れる。
