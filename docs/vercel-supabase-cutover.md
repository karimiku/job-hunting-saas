# Vercel + Supabase production cutover

実運用を Firebase / Cloud Run 前提から Vercel Services + Supabase Auth/Postgres へ切り替えるための手順。

コード側の準備と、Vercel/Supabase の実プロジェクト設定は別物として扱う。PR が merge 済みでも、下記の preview/prod smoke test が完了するまでは「実運用移行済み」とはみなさない。

## Current target architecture

- Frontend: Vercel Service `frontend` (`frontend/`)
- Backend API: Vercel Service `backend` (`backend/Dockerfile.vercel`)
- Public API path: `https://<app-domain>/backend/*`
- Auth: Supabase Auth Google OAuth
- API auth contract: frontend sends `Authorization: Bearer <Supabase access token>`
- Database: Supabase Postgres through pooler transaction mode
- Rollback path: Firebase Admin SDK / legacy session cookie envs remain unset unless rollback is needed

## Current external status (2026-07-02)

The production cutover is not complete yet. Code and Vercel environment
preparation are partially done, but several external settings still block real
traffic.

- Vercel project: `kamirikus-projects/job-hunting-saas`
  - Current production URL: `https://entre.kamiriku.com`
  - Current Root Directory: `frontend`
  - Impact: root `vercel.json` Services config is not active until Root
    Directory is changed to the repository root.
- Vercel env:
  - Production/Preview configured: Supabase public env, Supabase JWT issuer /
    JWKS / audience, pgxpool tuning, CORS, cookie, `PORT=8080`
  - Not configured: `DATABASE_URL`
  - `DATABASE_URL` is sensitive and must be copied from the selected Supabase
    pooler connection string only after explicit cutover approval.
  - Preview can omit `BACKEND_API_BASE_URL`; the frontend server runtime derives
    `https://$VERCEL_URL/backend` when `VERCEL=1`.
- Supabase project:
  - Project ref: `xmncflhhcebbytvrwoep`
  - Project URL: `https://xmncflhhcebbytvrwoep.supabase.co`
  - Region: Tokyo (`ap-northeast-1`)
  - Dashboard status: `Unhealthy` with recent Postgres errors; investigate
    before routing production traffic.
  - Google Auth provider: disabled
  - Site URL: `http://localhost:3000`
  - Redirect URLs: none
  - JWT current key: asymmetric `ECC (P-256)`; legacy HS256 remains only as a
    previous key for unexpired tokens.
  - Data API: disabled / no schemas can be queried, matching the Go API direct
    DB policy.
- Rollback:
  - Cloud Run service `entre-backend` remains active.
  - Firebase/GCP/Terraform assets remain frozen and available for rollback, not
    deleted during this migration.

## Supabase setup

1. Create or select the Supabase project.
2. Enable Google provider in Supabase Auth.
3. Add redirect URLs:
   - Preview: `https://<preview-domain>/auth/callback`
   - Production: `https://<prod-domain>/auth/callback`
   - Local: `http://localhost:3000/auth/callback`
4. Use asymmetric JWT signing keys so backend can verify access tokens through JWKS.
   - Required JWKS URL: `https://<project-ref>.supabase.co/auth/v1/.well-known/jwks.json`
   - If the project still uses HS256 shared-secret signing, this backend verifier will not accept those tokens.
5. Apply the current Postgres schema baseline to Supabase before routing real traffic.
6. Use the Supabase pooler transaction connection string for autoscaled Vercel runtime:

```text
postgres://postgres.<project-ref>:<password>@aws-<region>.pooler.supabase.com:6543/postgres?sslmode=require
```

Supabase docs: [Google OAuth](https://supabase.com/docs/guides/auth/social-login/auth-google), [JWT/JWKS](https://supabase.com/docs/guides/auth/jwts), [Postgres connection modes](https://supabase.com/docs/guides/database/connecting-to-postgres).

## Vercel setup

1. Change the Vercel project Root Directory to the repository root.
   - Current legacy frontend-only deployment may still be configured as `frontend`.
   - Vercel Services will not read the root `vercel.json` while the project Root Directory is `frontend`.
2. Confirm root `vercel.json` has both services:
   - `frontend` root: `frontend/`
   - `backend` root: `backend/`
3. Configure backend env:

```text
PORT=8080
DATABASE_URL=postgres://postgres.<project-ref>:<password>@aws-<region>.pooler.supabase.com:6543/postgres?sslmode=require
SUPABASE_AUTH_ISSUER=https://<project-ref>.supabase.co/auth/v1
SUPABASE_JWKS_URL=https://<project-ref>.supabase.co/auth/v1/.well-known/jwks.json
SUPABASE_JWT_AUDIENCE=authenticated
PGPOOL_MAX_CONNS=4
PGPOOL_MIN_CONNS=0
PGPOOL_MAX_CONN_IDLE_TIME=60s
PGPOOL_MAX_CONN_LIFETIME=30m
PGPOOL_HEALTH_CHECK_PERIOD=30s
PGX_DEFAULT_QUERY_EXEC_MODE=exec
PGAPPNAME=job-hunting-saas-api
CORS_ALLOWED_ORIGINS=https://<prod-domain>,https://*.vercel.app
COOKIE_SECURE=true
COOKIE_SAME_SITE=none
```

4. Configure frontend env:

```text
NEXT_PUBLIC_SUPABASE_URL=https://<project-ref>.supabase.co
NEXT_PUBLIC_SUPABASE_PUBLISHABLE_KEY=sb_publishable_...
BACKEND_API_BASE_URL=https://<prod-domain>/backend
BACKEND_API_ALLOWED_HOSTS=<prod-domain>,*.vercel.app
```

Preview deployments may omit `BACKEND_API_BASE_URL`; `serverFetch` derives the
same deployment URL from `VERCEL_URL` and appends `/backend`. Production should
keep an explicit `BACKEND_API_BASE_URL` during cutover so rollback can switch the
server-side backend target without code changes.

5. Do not set these in Supabase-only production unless rollback is active:

```text
FIREBASE_PROJECT_ID
FIREBASE_CREDENTIALS_FILE
GOOGLE_APPLICATION_CREDENTIALS
DEV_AUTH_ENABLED
```

Vercel docs: [Services](https://vercel.com/docs/services), [Container Images](https://vercel.com/docs/functions/container-images).

## Preview smoke test

Run this on a Vercel preview deployment before touching production DNS.

1. Open `https://<preview-domain>/backend/health`.
   - Expected: `ok`
2. Open `https://<preview-domain>/login` and sign in with Google.
   - Expected: redirect to `/auth/callback`, then `/dashboard`.
3. Open dashboard/inbox.
   - Expected: frontend API calls include `Authorization: Bearer ...` and backend returns 200.
4. Create one company/entry/task.
   - Expected: data persists after redeploy.
5. Check backend logs.
   - Expected: no Firebase initialization in normal path.
   - Expected: no prepared statement / pooler errors.
6. Check Supabase Observability.
   - Expected: app connections show `application_name=job-hunting-saas-api`.
   - Expected: connection count stays within `PGPOOL_MAX_CONNS * active backend containers` plus Supabase service baseline.

## Production cutover

1. Freeze Cloud Run/Firebase path only after preview smoke test passes.
2. Deploy Vercel production with the same env set.
3. Add production domain to Supabase Auth redirect URLs before login testing.
4. Point DNS/custom domain to Vercel production.
5. Run the same smoke test against production.
6. Keep Cloud Run/Firebase rollback env and deployment available until at least one real usage session completes.
7. After stable production verification, close the migration implementation issues and leave cleanup issues for Firebase/GCP docs and Terraform removal.

## Rollback

Rollback should be DNS/env based, not schema destructive.

1. Restore frontend/backend routing to the previous Cloud Run/Firebase deployment.
2. Re-enable Firebase envs only on the rollback backend.
3. Keep Supabase database unchanged.
4. Do not delete Supabase Auth users or external identities during rollback.
