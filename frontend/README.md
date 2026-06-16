# Entré Frontend

Next.js 16（App Router / SSR / Server Actions）の Web フロントエンド。Firebase で Google ログインし、ホーム / Entry / カンバン / タスク / 保存箱を提供する。backend が OpenAPI で公開する型に合わせて API を呼ぶ。

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

backend への接続先と Firebase Web SDK の設定が必要。`.env.local` は `.gitignore` 済みなのでコミットしない。

| 変数 | 用途 | 例 |
| --- | --- | --- |
| `BACKEND_API_BASE_URL` | backend API のベース URL（Next.js rewrite / Server Component 用、非公開） | `http://localhost:8080` |
| `BACKEND_API_ALLOWED_HOSTS` | backend proxy の許可 host allowlist | `localhost,127.0.0.1,api.entre.kamiriku.com,entre-backend-gfsd4pzoxq-an.a.run.app` |
| `NEXT_PUBLIC_FIREBASE_API_KEY` | Firebase Web API キー | — |
| `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN` | Firebase 認証ドメイン。本番はアプリの同一ドメインにする | `entre.kamiriku.com` |
| `NEXT_PUBLIC_FIREBASE_PROJECT_ID` | Firebase プロジェクト ID | `your-project-id` |
| `NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET` | Storage バケット | `your-project.appspot.com` |
| `NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID` | Messaging Sender ID | — |
| `NEXT_PUBLIC_FIREBASE_APP_ID` | Firebase App ID | — |
| `FIREBASE_AUTH_PROXY_HOST` | `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN` を独自ドメインにしたときの Firebase Hosting proxy 先 | `job-hunting-saas.firebaseapp.com` |

値は Firebase Console の「プロジェクトの設定 → マイアプリ（Web アプリ）」から取得する。`NEXT_PUBLIC_` 接頭辞の変数はクライアントバンドルに埋め込まれる前提の公開値（Firebase Web SDK のキーは公開設計）。backend proxy の接続先は `BACKEND_API_BASE_URL` を使い、公開 env には置かない。

本番の Google redirect ログインでは、ブラウザの third-party storage 制限を避けるため `NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN` を `entre.kamiriku.com` にする。Next.js は `/__/auth/*` を `FIREBASE_AUTH_PROXY_HOST` へ rewrite するため、Firebase helper はブラウザから同一originとして見える。

## 認証フロー

1. `/login` で Google ログイン（Firebase Web SDK）
2. 取得した ID トークンを backend `POST /auth/session` に送り、httpOnly セッション Cookie を発行
3. 以降の SSR / Server Action は Cookie を backend に転送して API を呼ぶ（`src/lib/api/server.ts`）

backend 側の Firebase / CORS / Cookie 設定とセットで動く。横断的な手順はルート [README.md](../README.md) の「β環境セットアップ」を参照。

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
