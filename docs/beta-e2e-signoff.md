# β E2E サインオフ

βの完成条件「Chrome 拡張で保存した求人ページが、ログインユーザーの DB に入り、Web の Inbox から Entry 化して管理できる」を end-to-end で確認するための手順。

## 自動 E2E（Playwright）

`frontend/e2e/` の Playwright スイートが、認証なしで検証できる範囲を自動でカバーする。

```bash
cd frontend
pnpm exec playwright install --with-deps   # 初回のみ
pnpm test:e2e
```

- `auth-guards.spec.ts` — `/inbox` `/entry` `/kanban` `/task` などの認証必須ページが未ログインで `/login` にリダイレクトされる
- `landing.spec.ts` / `onboarding.spec.ts` — LP・オンボーディングの表示

`pnpm dev`（frontend のみ）で起動するため、backend と Firebase 認証を要する保存→Entry 化フローは自動 E2E には含まれない（SSR / Server Action はブラウザ側のネットワークモックが効かないため）。その部分は下記の手動サインオフで確認する。

## 手動サインオフ

[README のβ環境セットアップ](../README.md#β環境セットアップ拡張--inbox--entry-を実ログインで通す) に従い backend / frontend / chrome-extension を起動した状態で実施する。

| # | 手順 | 期待結果 | 確認 |
| --- | --- | --- | --- |
| 1 | `http://localhost:3000/login` で Google ログイン | Dashboard に遷移 | ☐ |
| 2 | ログアウト状態で拡張から保存 | ログイン誘導メッセージ + Web ログインボタン | ☐ |
| 3 | ログイン状態で求人ページを開き拡張から保存 | 保存成功（紙吹雪等） | ☐ |
| 4 | 対応サイト外のページで保存 | URL / タイトル / source のフォールバックで保存できる | ☐ |
| 5 | `http://localhost:3000/inbox` | 保存した clip が表示される | ☐ |
| 6 | 同じ URL をもう一度保存 | clip が二重に増えない（重複 URL 最小対応） | ☐ |
| 7 | Inbox で clip を Entry 化 | Company + Entry が作成され、clip が Inbox から消える | ☐ |
| 8 | `/entry`・`/kanban` | 作成した Entry が会社名つきで表示される | ☐ |
| 9 | Entry 詳細で Task を追加し `/task` で完了トグル | 状態が更新され再読込後も保持される | ☐ |
| 10 | Inbox で clip を削除 | clip が一覧から消える | ☐ |
| 11 | backend 再起動後に `/inbox` を再読込 | PostgreSQL 上の clip が残る（永続化） | ☐ |

## サインオフ記録

| 項目 | 値 |
| --- | --- |
| 実施日 | |
| 実施者 | |
| commit | |
| 結果 | PASS / FAIL |
| 備考 | |
