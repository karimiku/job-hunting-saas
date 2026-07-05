# Entré Frontend

Next.js 16（App Router / SSR / Server Actions）の Web フロントエンド。Supabase Auth で Google ログインし、ホーム / Entry / カンバン / タスク / 保存箱を提供する。backend が OpenAPI で公開する型に合わせて API を呼ぶ。

## コア画面

- `/dashboard`: 次の行動と今日のタスクを確認するホーム
- `/entry`: 応募先を Entry 単位で一覧・追加する
- `/entry/[entryId]`: Entry のステージ、メモ、タスクを確認する
- `/kanban`: Entry を選考ステージ別に俯瞰するカンバン
- `/task`: Entry に紐づく締切・予定を追加/完了する
- `/inbox`: 保存した求人を Entry に変換する保存箱
- `/profile`: アカウント情報とログアウト

`/roadmap` はコア導線から外しており、認証後は `/dashboard` に戻す。

## LP

`/` は Entry とカンバンを中心にしたシンプルなランディングページ。訴求する流れは「保存箱に残す → Entry にまとめる → カンバンで動かす」だけに絞り、画面説明は `frontend/public/marketing/` の実サービススクリーンショットを使う。

## セットアップ

```bash
pnpm install
pnpm dev        # 開発サーバー (http://localhost:3000)
pnpm build      # 本番ビルド
pnpm start      # 本番サーバー
pnpm test       # vitest
pnpm lint       # ESLint
```

## 環境変数 (`frontend/.env.local`)

backend への接続先と Supabase Auth の公開設定が必要。`.env.local` は `.gitignore` 済みなのでコミットしない。

| 変数 | 用途 | 例 |
| --- | --- | --- |
| `BACKEND_API_BASE_URL` | backend API のベース URL（Server Component 用、非公開）。Vercel Services 本番では `/backend` prefix 付き URL を使う。Preview は未設定でも `VERCEL_URL` から同一 deployment の `/backend` を使う | `http://localhost:8080`, `https://entre.kamiriku.com/backend` |
| `BACKEND_API_ALLOWED_HOSTS` | backend proxy の許可 host allowlist | `localhost,127.0.0.1,entre.kamiriku.com,*.vercel.app,api.entre.kamiriku.com,entre-backend-gfsd4pzoxq-an.a.run.app` |
| `NEXT_PUBLIC_SUPABASE_URL` | Supabase project URL | `https://<project-ref>.supabase.co` |
| `NEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY` | Supabase publishable key。ブラウザに公開してよいキーだけを設定する | `sb_publishable_...` |

Supabase の値は Dashboard の Project Settings / API から取得する。`NEXT_PUBLIC_` 接頭辞の変数はクライアントバンドルに埋め込まれる。`service_role` や secret key は絶対に設定しない。backend proxy の接続先は `BACKEND_API_BASE_URL` を使い、公開 env には置かない。

Google OAuth を使うには Supabase Auth の Google provider を有効化し、Redirect URLs に `http://localhost:3000/auth/callback` と本番の `https://<frontend-domain>/auth/callback` を登録する。

## 認証フロー

1. `/login` で Google ログイン（Supabase Auth OAuth）
2. `/auth/callback` で authorization code を Supabase session cookie に交換
3. Client Component / Server Component / Server Action は Supabase session から access token を取得し、Go backend に `Authorization: Bearer <token>` を送る
4. Go backend は Supabase JWKS で JWT を検証し、`external_identities(provider=supabase, subject=sub)` から app user に解決する

backend 側の `SUPABASE_AUTH_ISSUER` / `SUPABASE_JWKS_URL` / `DATABASE_URL` 設定とセットで動く。横断的な手順はルート [README.md](../README.md) の「β環境セットアップ」を参照。

## ディレクトリ

```
src/
├── app/                 # App Router（pages / layouts / server actions）
├── components/entre/     # アプリ UI（EntryListView, KanbanBoard, InboxList 等）
├── components/landing/   # LP
└── lib/api/              # API クライアント（client.ts: Client用 / server.ts: SSR用）
```

## テスト

- ユニット/コンポーネント: `pnpm test`（vitest + Testing Library + MSW）
- E2E: `pnpm test:e2e`（Playwright）
