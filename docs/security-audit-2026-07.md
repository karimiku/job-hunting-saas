# セキュリティ監査レポート — 2026-07（β / リリース前）

対象: working tree（未コミットの Supabase Auth / Vercel 移行を含む現状）。backend (Go) / frontend (Next.js 16) / chrome-extension を精読。観点は SQLi・IDOR・JWT 検証・秘密情報・CORS/CSRF・SSRF・依存脆弱性。

## サマリ（深刻度順）

| ID | 深刻度 | 対象 | 概要 | Issue |
| --- | --- | --- | --- | --- |
| V1 | Medium | dev_auth_handler.go:44,92 | dev セッション発行の本番ガードが攻撃者制御の Origin/Referer に依存 | #235 |
| V2 | Medium | rate_limit.go:150 | レート制限キーが偽装可能な X-Forwarded-For 依存で回避可能 | #236 |
| V3 | Medium(要確認) | main.go:484 / csrf.go:79 | CORS `*.vercel.app` + credentials + SameSite=None でクロスオリジン情報窃取 | #237 |
| V4 | Low/Medium | frontend/src/lib/api/server.ts:108-145 | ヘッダ由来ホストへ Bearer/cookie 転送しうる SSRF フォールバック | #238 |
| V5 | Low | frontend/src/app/auth/callback/route.ts:36 | x-forwarded-host によるオープンリダイレクト | #239 |
| V6 | Blocker | supabaseauth/verifier.go:130 | ビルド不能・Supabase 認証パス未完成（enum/provider 未整合） | #240 |
| V7 | Low | supabaseauth/verifier.go:226 | JWKS 負のキャッシュ無し（外部フェッチ増幅） | #241 |
| V8 | Low/Info | frontend/next.config.ts:19 | CSP に script-src 無く XSS 緩和が限定的 | #242 |

Critical / High 相当の即時悪用可能な脆弱性は検出されず。Medium 3 件はいずれも「本番 env の設定次第で成立/しない」性質。

## 良かった点（回帰させないこと）

- **IDOR 対策が堅牢**: entry / task / inbox_clip / es_memo / selection_flow / ai_token の全 sqlc クエリが `WHERE user_id = $1`（または `JOIN entries ... e.user_id = $1`）でスコープ。ハンドラは例外なく `middleware.GetUserID(ctx)` を渡す。生 SQL・`fmt.Sprintf` クエリ組み立てなし。entry Create は `companyRepo.FindByID(userID, companyID)` で所有権検証。
- **秘密情報のコミットなし**（`git ls-files` / `git grep` 確認）。SA 鍵・`.env` 系は `.gitignore` 済み。
- **JWT 検証**: `ValidMethods` を RS256/ES256/EdDSA に固定（HMAC 混同なし）、`alg none` 不可、kid 必須、issuer/audience required、JWKS 1MB 上限・10分キャッシュ。
- **AI アクセストークン**: 32B CSPRNG、DB は SHA-256 ハッシュのみ、平文は一度だけ返却。token 管理系は Bearer を拒否（session 必須）。
- **依存**: `pnpm audit --prod` = No known vulnerabilities。

---

## 詳細

### V1 — [Medium] dev セッション発行の認可がクライアント制御ヘッダに依存
`backend/internal/handler/dev_auth_handler.go:44,92-102` / `backend/cmd/server/main.go`

`CreateSession`（`POST /dev/session`、任意 email で署名済み session cookie を発行＝任意ユーザーなりすまし）の本番ガード `isLocalDevRequest` が、攻撃者が自由に設定できる `Origin` / `Referer` ヘッダのいずれかが localhost なら true を返す。多層防御は `DEV_AUTH_ENABLED` ゲート、`isProductionRuntime()`、route 前段の OriginGuard。しかし `isProductionRuntime()` は `APP_ENV/GO_ENV/ENV/GIN_MODE=="production"` のみ判定で、Vercel/Cloud Run で未設定だと素通り。

**攻撃シナリオ**: (a)`DEV_AUTH_ENABLED=true` が本番に残る、(b)`APP_ENV` 等未設定、(c)`CORS_ALLOWED_ORIGINS` に開発の名残の `http://localhost:3000` が残る、が揃うと `Origin: http://localhost:3000` を付けて `POST /dev/session {"email":"victim@..."}` で任意ユーザーの有効 session を取得できる。

**修正**: dev ログインの本番無効化を、クライアント制御ヘッダではなくサーバ側シグナルに一本化。`isLocalDevRequest` から Origin/Referer を除外し `r.Host`（＋可能なら TCP peer が loopback か）のみで判定する。

### V2 — [Medium] レート制限のクライアントIPが X-Forwarded-For 偽装で回避可能
`backend/internal/middleware/rate_limit.go:150-179`

`X-Forwarded-For` を無条件に信頼し「右端の parse 可能な IP」を採用。信頼プロキシのホップ数を固定していないため、XFF を都度変えると各リクエストが別キー扱いになり 429 に到達しない。userID キーの limiter は偽装不可なので影響は未認証エンドポイントに限定。

**修正**: 実インフラ（Vercel→Cloud Run）の XFF 付与仕様を確定し「右から N 番目」を採用、またはプラットフォーム提供の検証済みクライアント IP を使う。

### V3 — [Medium/要人間確認] CORS `*.vercel.app` + credentials
`backend/cmd/server/main.go` / `backend/internal/middleware/csrf.go`（`OriginAllowed` はワイルドカード対応）

credentialed CORS で `https://*.vercel.app` を許可すると、任意の `*.vercel.app` オリジンが credentials 付きリクエストを許可される。移行期に legacy `session` cookie が有効かつ `COOKIE_SAME_SITE=none` だと、`evil.vercel.app` から被害者データを窃取可能。Supabase 移行後の主経路は Bearer（ambient でない）なので移行完了後は成立しない。

**修正**: credentials 付き CORS では広いワイルドカードを使わず確定 origin のみ列挙。プレビュー URL は都度 allowlist に追加する運用に。`CORS_ALLOWED_ORIGINS` と `COOKIE_SAME_SITE` の本番実値を確認。

### V4 — [Low/Medium] serverFetch のヘッダ由来ホストへのトークン転送
`frontend/src/lib/api/server.ts:31-64,108-145`

`serverFetch` は cookie 全体と Supabase access token を backend へ転送。転送先は `BACKEND_API_BASE_URL` 未設定時に `VERCEL_URL || x-forwarded-host || host` から組み立て、allowlist に `*.vercel.app` を含む。`VERCEL_URL` 空かつ `BACKEND_API_BASE_URL` 未設定だと、`x-forwarded-host` を任意の `*.vercel.app` にして被害者 Bearer を攻撃者ホストへ送出させうる。

**修正**: 資格情報転送 proxy の allowlist から `*.vercel.app` を外し確定ホストのみ。`BACKEND_API_BASE_URL` を本番必須化。

### V5 — [Low] auth コールバックのオープンリダイレクト
`frontend/src/app/auth/callback/route.ts:36-45`

本番で `trustedRedirectOrigin` が `x-forwarded-host` から redirect origin を組み立てる（検証が緩い）。ヘッダ注入可能な構成だと任意 host の `/dashboard` へ誘導（フィッシング）。session cookie は本アプリ domain に載るためトークン自体は漏れない。

**修正**: リダイレクト先 origin は確定した公開 URL の env を使い `x-forwarded-host` に依存しない。

### V6 — [Blocker] working tree がビルド不能（Supabase 認証パス未完成）
`backend/internal/infra/supabaseauth/verifier.go:130` / `backend/internal/domain/value/auth_provider.go`

`value.AuthProviderSupabase()` が未定義で `go build ./...` が失敗。`validAuthProviders` は `"google"` のみで `"supabase"` を含まず、補完しても main.go の UserSync（`Provider:"supabase"`）が `ErrAuthProviderInvalid` で失敗。schema 乖離: `supabase/migrations/...baseline.sql` は enum に `supabase` を追加済みだが `sql/schema.sql` は `('google')` のまま。

**修正**: `AuthProviderSupabase()` 追加・`validAuthProviders` 拡張・`schema.sql` と migration の enum 整合。修正後に verifier.go を再レビュー。

### V7 — [Low] JWKS の負のキャッシュ無し
`backend/internal/infra/supabaseauth/verifier.go:226-259`

未知 kid ごとに `refreshJWKS`（最大5s・1MBの外向き HTTP）を実行。ランダム kid 多投で外部フェッチが増幅。失敗 kid のネガティブキャッシュ/クールダウンが望ましい。

### V8 — [Low/Info] CSP に script-src なし
`frontend/next.config.ts:19-27`

CSP は `base-uri 'self'; object-src 'none'; frame-ancestors 'self'` のみ。`default-src 'self'` ベースの script-src（nonce/strict-dynamic）で XSS 第二防御線を強化する余地。

---

## 要人間確認
1. `.env.example`（root/backend/frontend）に実 secret が紛れていないか（監査時に権限で未読）。
2. `CORS_ALLOWED_ORIGINS` と `COOKIE_SAME_SITE` の本番実値（V3）。
3. Vercel→Cloud Run 間の XFF 付与仕様（V2）、`BACKEND_API_BASE_URL`/`VERCEL_URL` の設定有無（V4）。
