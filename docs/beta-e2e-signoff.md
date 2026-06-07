# β E2E サインオフ

βの完成条件「保存した求人が、ログインユーザーの DB に入り、Web の保存箱から Entry 化され、カンバンとタスクで管理できる」を end-to-end で確認するための手順。

## 自動 E2E（Playwright）

`frontend/e2e/` の Playwright スイートが、Firebase / PostgreSQL なしで検証できる範囲を自動でカバーする。`pnpm test:e2e` は Playwright 用の軽量 mock API を起動し、SSR / Server Action を含む主要導線をブラウザ操作で確認する。

```bash
cd frontend
pnpm exec playwright install --with-deps   # 初回のみ
pnpm test:e2e
```

- `auth-guards.spec.ts` — `/inbox` `/entry` `/kanban` `/task` などの認証必須ページが未ログインで `/login` にリダイレクトされる
- `beta-core-flow.spec.ts` — mock API に保存済み clip を投入し、`/inbox` で Entry 化 → `/entry` `/kanban` 表示 → `/task` 追加・完了までを desktop / mobile で確認する
- `landing.spec.ts` / `onboarding.spec.ts` — LP・オンボーディングの表示

Chrome 拡張本体・Firebase Google ログイン・PostgreSQL 永続化は実環境依存のため、自動 E2E では代替しない。その部分は下記の手動サインオフで確認する。

## 手動サインオフ

[README のβ環境セットアップ](../README.md#β環境セットアップ保存箱--entry--カンバンタスク) に従い backend / frontend を起動する。拡張連携まで確認する場合のみ chrome-extension も起動する。

| # | 手順 | 期待結果 | 確認 |
| --- | --- | --- | --- |
| 1 | `http://localhost:3000/login` で Google ログイン | Dashboard に遷移 | ☐ |
| 2 | `http://localhost:3000/inbox` を開く | 保存箱が表示される | ☐ |
| 3 | 保存済み clip を用意する（APIまたは拡張） | 保存箱に clip が表示される | ☐ |
| 4 | 同じ URL をもう一度保存 | clip が二重に増えない（重複 URL 最小対応） | ☐ |
| 5 | 保存箱で clip を Entry 化 | Company + Entry が作成され、clip が保存箱から消える | ☐ |
| 6 | `/entry`・`/kanban` | 作成した Entry が会社名つきで表示される | ☐ |
| 7 | Entry 詳細で Task を追加し `/task` で完了トグル | 状態が更新され再読込後も保持される | ☐ |
| 8 | 保存箱で clip を削除 | clip が一覧から消える | ☐ |
| 9 | backend 再起動後に `/inbox` を再読込 | PostgreSQL 上の clip が残る（永続化） | ☐ |
| 10 | 任意: ログアウト状態で拡張から保存 | ログイン誘導メッセージ + Web ログインボタン | ☐ |
| 11 | 任意: 対応サイト外のページで拡張保存 | URL / タイトル / source のフォールバックで保存できる | ☐ |

## サインオフ記録

| 項目 | 値 |
| --- | --- |
| 実施日 | |
| 実施者 | |
| commit | |
| 結果 | PASS / FAIL |
| 備考 | |
