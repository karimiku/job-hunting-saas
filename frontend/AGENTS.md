<!-- BEGIN:nextjs-agent-rules -->
# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` before writing any code. Heed deprecation notices.
<!-- END:nextjs-agent-rules -->

## Pull Request ルール（必須・厳守）

- **PR 本文は必ず `.github/pull_request_template.md`（リポジトリルート）の見出し構成に従う**。`gh pr create --body` ではテンプレートが自動適用されないため、自分で構成を再現すること。例外なし。
- テンプレのテストチェックリストは Go 向けなので、フロントエンド PR では実際に実行したコマンド（vitest / Playwright / lint / build）に置き換える。
- **PR タイトルは日本語で書く**（prefix `feat:` / `fix:` / `perf:` 等は可。例: `perf: アプリ内ナビゲーションを高速化`）。
